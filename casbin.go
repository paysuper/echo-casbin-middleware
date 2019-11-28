/* Package casbin provides middleware to enable ACL, RBAC, ABAC authorization support. */

package casbin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/micro/go-micro/client"
	"github.com/paysuper/casbin-server/pkg"
	"github.com/paysuper/casbin-server/pkg/generated/api/proto/casbinpb"
	"log"
	"strings"
)

type Logger interface {
	Printf(string, ...interface{})
}

type logger struct{}

type CtxUserExtractor func(c echo.Context) string

type (
	EnforceMode int
	// Config defines the config for CasbinAuth middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		// Enforce mode
		Mode EnforceMode
		// Logger
		Logger Logger
		// CtxUserExtractor
		CtxUserExtractor CtxUserExtractor
		// Casbin micro service.
		client casbinpb.CasbinService
	}
)

const (
	EnforceModeUnknown EnforceMode = iota
	EnforceModePermissive
	EnforceModeEnforcing
	EnforceModeDisabled
)

var (
	DefaultLogger = &logger{}
	// DefaultConfig is the default CasbinAuth middleware config.
	DefaultConfig = Config{
		Skipper: middleware.DefaultSkipper,
		Mode:    EnforceModeEnforcing,
		Logger:  DefaultLogger,
	}
)

func (*logger) Printf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

// Middleware returns a CasbinAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func Middleware(c client.Client, mode EnforceMode) echo.MiddlewareFunc {
	cfg := DefaultConfig
	cfg.client = casbinpb.NewCasbinService(pkg.ServiceName, c)
	if mode != EnforceModeUnknown {
		cfg.Mode = mode
	}
	return MiddlewareWithConfig(c, cfg)
}

// MiddlewareWithConfig returns a CasbinAuth middleware with config.
// See `Middleware()`.
func MiddlewareWithConfig(c client.Client, config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	if config.CtxUserExtractor == nil {
		panic("CtxUserExtractor callback function required")
	}
	config.client = casbinpb.NewCasbinService(pkg.ServiceName, c)
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
	// EnforceModeDisabled return true (permission granted) early
	if cfg.Mode == EnforceModeDisabled {
		return true
	}
	//
	user := cfg.CtxUserExtractor(c)
	method := c.Request().Method
	path := c.Request().URL.Path
	var routeType string
	switch {
	case len(c.ParamNames()) > 0:
		routeType = "param"
	case strings.HasSuffix(c.Path(), "*"):
		routeType = "wildcard"
	default:
		routeType = "static"
	}
	// Check permissions
	_, err := cfg.client.Enforce(c.Request().Context(), &casbinpb.EnforceRequest{Params: []string{user, path, method, routeType}})
	if err == nil {
		return true
	}
	// EnforceModePermissive log and return true (permission granted)
	if cfg.Mode == EnforceModePermissive {
		cfg.Logger.Printf("casbin enforce user:%v path:%v method:%v err:%v", user, path, method, err)
		return true
	}
	// EnforceModeEnforcing return false (permission forbidden)
	return false
}
