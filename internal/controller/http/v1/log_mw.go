package v1

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/s4lat/gosavingsbot/internal/log"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"
)

type RequestLog struct {
	UUID string `json:"uuid"`
	Path string `json:"path"`
	Body string `json:"body"`
}

type ResponseLog struct {
	UUID       string `json:"uuid"`
	Body       string `json:"body"`
	StatusCode int    `json:"status_code"`
}

func LogMiddleware() echo.MiddlewareFunc {
	bodyDumpMW := middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, reqBody []byte, respBody []byte) {
			c.Set("reqBody", string(reqBody))
			c.Set("respBody", string(respBody))
		},
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqUUID := uuid.New().String()
			c.Set("reqUUID", reqUUID)
			c.Set("logger", log.Sugar().With("uuid", reqUUID))

			if SkipLog(c) {
				return next(c)
			}

			start := time.Now()

			err := bodyDumpMW(next)(c)

			end := time.Now()
			duration := end.Sub(start)

			query, err := url.QueryUnescape(c.Request().URL.RawQuery)
			if err != nil {
				query = c.Request().URL.RawQuery
			}

			log.Sugar().Infow("req",
				"time", start.Format(time.RFC3339),
				"uuid", reqUUID,
				"path", c.Path(),
				"body", c.Get("reqBody").(string),
				"query", query,
			)

			log.Sugar().Infow("resp",
				"time", end.Format(time.RFC3339),
				"uuid", reqUUID,
				"path", c.Path(),
				"status_code", c.Response().Status,
				"body", c.Get("respBody").(string),
				"duration", duration.Milliseconds(),
			)

			return err
		}
	}
}

func SkipLog(c echo.Context) bool {
	return strings.Contains(c.Request().URL.Path, "/healthz")
}

func getLoggerFromEchoContext(c echo.Context) *zap.SugaredLogger {
	logger := c.Get("logger").(*zap.SugaredLogger)
	return logger
}
