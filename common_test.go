package casbin

import (
	"bufio"
	"context"
	"flag"
	casbinServer "github.com/paysuper/casbin-server"
	"github.com/paysuper/casbin-server/casbinpb"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/micro/go-micro/client"
)

var (
	reqs          []*http.Request
	getHandler    = func(ctx echo.Context) error { return ctx.String(http.StatusOK, "GET\n") }
	postHandler   = func(ctx echo.Context) error { return ctx.String(http.StatusCreated, "POST\n") }
	deleteHandler = func(ctx echo.Context) error { return ctx.NoContent(http.StatusNoContent) }
	patchHandler  = func(ctx echo.Context) error { return ctx.String(http.StatusOK, "PATCH\n") }
	putHandler    = func(ctx echo.Context) error { return ctx.String(http.StatusOK, "PUT\n") }
)

type (
	routerFn func(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
)

func TestMain(m *testing.M) {
	var err error

	// rewrite args to disable test flags
	os.Args = []string{os.Args[0]}
	// Echo instance
	e := echo.New()
	e.HideBanner = true
	c := client.NewClient()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(Middleware(c))

	mapfn := make(map[string]routerFn)
	mapfn["GET"] = e.GET
	mapfn["POST"] = e.POST
	mapfn["PATCH"] = e.PATCH
	mapfn["DELETE"] = e.DELETE
	mapfn["PUT"] = e.PUT

	if err := setupHandlers(e, mapfn); err != nil {
		log.Fatal(err)
	}

	// Route => handler
	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World!\n")
	})

	port, err := getFreeTestPort()
	if err != nil {
		log.Fatal(err)
	}
	httpAddr := "localhost:" + port

	go func(err error) {
		// Start server
		err = e.Start(httpAddr)
	}(err)

	// close server
	defer e.Close()

	flag.Parse()
	if err = os.Setenv("CASBIN_ADAPTER", "sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err = os.Setenv("CASBIN_DSN", "./database.sqlite3"); err != nil {
		log.Fatal(err)
	}

	app := casbinServer.NewApplication()
	app.Init()
	go func() {
		app.Run()
	}()

	req := &casbinpb.UserRoleRequest{User: "12345678-1234-1234-1234-123456789012", Role: "system_admin"}
	rsp := &casbinpb.BoolReply{}
	if err = app.csvc.AddRoleForUser(context.Background(), req, rsp); err != nil {
		log.Fatal(err)
	}
	req = &casbinpb.UserRoleRequest{User: "12345678-1234-1234-1234-123456789012", Role: "system_view_only"}
	if err = app.csvc.AddRoleForUser(context.Background(), req, rsp); err != nil {
		log.Fatal(err)
	}

	sl := &casbinpb.ArrayReply{}
	if err := app.csvc.GetRolesForUser(context.Background(), &casbinpb.UserRoleRequest{
		User: "12345678-1234-1234-1234-123456789012"}, sl); err != nil {
		log.Fatal(err)
	}
	ecode := m.Run()

	os.Exit(ecode)
}

func TestRequests(t *testing.T) {
	var rsp *http.Response
	var err error

	c := &http.Client{}
	for _, req := range reqs {
		rsp, err = client.Do(req)
	}
}

func setupHandlers(e *echo.Echo, mapfn map[string]routerFn) error {
	f, err := os.OpenFile("policy.conf", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	var req *http.Request

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		fields := strings.Split(line, ",")
		if len(fields) < 2 {
			continue
		}

		for i, f := range fields {
			fields[i] = strings.TrimSpace(f)
		}

		method := fields[len(fields)-1]
		path := fields[2]

		fn, ok := mapfn[method]
		if !ok {
			continue
		}

		if strings.Count(path, "*") > 0 {
			rid, _ := uuid.NewRandom()
			path = strings.ReplaceAll(path, "*", rid)
		}

		req = http.NewRequest(method, path, nil)
		req.Header.Add("X-USER", `W/"wyzzy"`)

		switch method {
		case "GET":
			fn(fields[2], getHandler)
		case "POST":
			fn(fields[2], postHandler)
		case "PUT":
			fn(fields[2], putHandler)
		case "PATCH":
			fn(fields[2], patchHandler)
		case "DELETE":
			fn(fields[2], deleteHandler)
		}

		reqs = append(reqs, req)
	}

	return nil
}

func getFreeTestPort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = l.Close()
	}()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}
