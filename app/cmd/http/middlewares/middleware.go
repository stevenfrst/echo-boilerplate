package middlewares

import (
	"github.com/coocood/freecache"
	cache "github.com/gitsight/go-echo-cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm/module/apmechov4"
)

func InitMiddleware(e *echo.Echo) *echo.Echo {
	e.Pre(middleware.RemoveTrailingSlash())
	c := freecache.NewCache(1024 * 1024 * 100)
	e.Use(cache.New(&cache.Config{}, c))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(apmechov4.Middleware())
	return e
}
