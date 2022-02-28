package middlewares

import (
	"github.com/coocood/freecache"
	cache "github.com/gitsight/go-echo-cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitMiddleware(e *echo.Echo) *echo.Echo {
	e.Pre(middleware.RemoveTrailingSlash())
	c := freecache.NewCache(1024 * 1024)
	e.Use(cache.New(&cache.Config{}, c))
	return e
}
