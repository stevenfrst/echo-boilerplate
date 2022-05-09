package routes

import (
	"echo-boilerplate/app/cmd/http/middlewares"
	"echo-boilerplate/factory"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	logger "github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	"net/http"
	"strings"
	"time"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func InitHttp() {
	f := factory.InitFactoryHTTP()
	e := echo.New()
	middlewares.InitMiddleware(e)
	e.Validator = &CustomValidator{Validator: validator.New()}

	v1 := e.Group("api/v1")
	// Auth Endpoint
	v1.POST("/users", f.User.Register)
	v1.POST("/users/login", f.User.Login)
	v1.PATCH("/users/change-password", f.User.ChangePassword, middleware.JWTWithConfig(f.ConfigJWT))
	v1.GET("/users/me", f.User.GetDetail, middleware.JWTWithConfig(f.ConfigJWT))
	v1.DELETE("/users", f.User.Delete, middleware.JWTWithConfig(f.ConfigJWT))
	v1.GET("/users", f.User.ListAllUsers)

	v1.GET("/verify", f.User.Verify)
	e.Logger.Fatal(e.Start(":1234"))
}

func traceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			// set childCtx so each API request will creates new serverSpan log
			spanCtx, _ := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(req.Header),
			)

			serverSpan := opentracing.StartSpan(c.Request().URL.Path, ext.RPCServerOption(spanCtx))
			c.Set("serverSpan", serverSpan)

			defer func() {
				serverSpan.Finish()
			}()

			var headers []log.Field
			for k, v := range req.Header {
				headers = append(headers, log.String(k, strings.Join(v, ", ")))
			}

			serverSpan.LogFields(
				headers...,
			)

			traceID := "no-tracer-id"
			if sc, ok := serverSpan.Context().(jaeger.SpanContext); ok {
				traceID = sc.String()
			}

			// inject to response header
			opentracing.GlobalTracer().Inject(
				serverSpan.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(c.Response().Header()),
			)

			serverSpan.SetTag("endpoint", req.RequestURI)
			serverSpan.SetTag("host", req.Host)
			serverSpan.SetTag("clientIP", c.RealIP())
			serverSpan.SetTag("http.status", res.Status)
			serverSpan.SetTag("userAgent", req.UserAgent())

			start := time.Now()
			// continue the request
			if errMiddleware := next(c); errMiddleware != nil {
				c.Error(errMiddleware)
				c.Response().Committed = true
				return errMiddleware
			}

			stop := time.Now()

			logger.Debug().
				Float64("duration", float64(stop.Sub(start).Nanoseconds())/float64(time.Millisecond)).
				Int("status", res.Status).
				Str("protocol", req.Proto).
				Str("endpoint", req.RequestURI).
				Str("host", req.Host).
				Str("clientIP", c.RealIP()).
				Str("method", req.Method).
				Str("tracerID", traceID).
				Msg("handle request")
			//fmt.Println(traceID)
			return
		}
	}
}