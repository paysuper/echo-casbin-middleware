/* Package casbin provides middleware to enable ACL, RBAC, ABAC authorization support. */

package casbin

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	client "github.com/micro/go-micro/client"
	casbinpb "github.com/paysuper/casbin-server/casbinpb"
)

type Logger interface {
	Printf(string, ...interface{})
}

type logger struct{}

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
	DefaultConfig = &Config{
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
	cfg.client = casbinpb.NewCasbinService("", c)
	if mode != EnforceModeUnknown {
		cfg.Mode = mode
	}
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
	// EnforceModeDisabled return true (permission granted) early
	if cfg.Mode == EnforceModeDisabled {
		return true
	}

	// get user from context
	user := ""
	method := c.Request().Method
	path := c.Request().URL.Path

	// check permisions
	_, err := cfg.client.Enforce(c.Request().Context(), &casbinpb.EnforceRequest{Params: []string{user, path, method}})
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
