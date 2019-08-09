/* Package casbin provides middleware to enable ACL, RBAC, ABAC authorization support.

Simple example:

	package main

	import (
		"github.com/casbin/casbin"
		"github.com/labstack/echo/v4"
		"github.com/labstack/echo-contrib/casbin" casbin-mw
	)

	func main() {
		e := echo.New()

		// Mediate the access for every request
		e.Use(casbin-mw.Middleware(casbin.NewEnforcer("auth_model.conf", "auth_policy.csv")))

		e.Logger.Fatal(e.Start(":1323"))
	}

Advanced example:

	package main

	import (
		"github.com/labstack/echo/v4"
	authzmw	"github.com/labstack/echo-contrib/casbin"
	)

	func main() {
		e := echo.New()

		echo.Use(authzmw.Middleware(c))

		e.Logger.Fatal(e.Start(":1323"))
	}
*/

package casbin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	client "github.com/micro/go-micro/client"
	casbinpb "github.com/paysuper/casbin-server/casbinpb"
)

type (
	// Config defines the config for CasbinAuth middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Casbin micro service.
		client casbinpb.CasbinService
	}
)

var (
	// DefaultConfig is the default CasbinAuth middleware config.
	DefaultConfig = &Config{
		Skipper: middleware.DefaultSkipper,
	}
)

// Middleware returns a CasbinAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func Middleware(c client.Client) echo.MiddlewareFunc {
	cfg := DefaultConfig
	cfg.client = casbinpb.NewCasbinService("", c)
	return MiddlewareWithConfig(cfg)
}

// MiddlewareWithConfig returns a CasbinAuth middleware with config.
// See `Middleware()`.
func MiddlewareWithConfig(config *Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) || config.CheckPermission(c) {
				return next(c)
			}

			return echo.ErrForbidden
		}
	}
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (cfg *Config) CheckPermission(c echo.Context) bool {
	// get user from context
	user := ""
	method := c.Request().Method
	path := c.Request().URL.Path
	rsp, err := cfg.client.Enforce(c.Request().Context(), &casbinpb.EnforceRequest{Params: []string{user, path, method}})
	if err != nil || !rsp.Res {
		return false
	}
	return true
}
