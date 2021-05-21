package handler

import (
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}

			logger.WithFields(logrus.Fields{
				"ip":            c.RealIP(),
				"method":        req.Method,
				"response_time": time.Since(start),
			}).Info(req.RequestURI)
			return err
		}
	}
}
