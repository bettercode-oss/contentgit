package middlewares

import (
	"bytes"
	"contentgit/foundation"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoggingWithZap(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// log request ID
		requestID := uuid.New().String()
		if xRequestId := c.Writer.Header().Get("X-Request-Id"); xRequestId != "" {
			requestID = xRequestId
		}
		c.Request = c.Request.WithContext(foundation.ContextProvider().SetRequestId(c.Request.Context(), requestID))

		reqLogger := logger.With(
			zap.String("requestID", requestID),
		)
		c.Request = c.Request.WithContext(foundation.ContextProvider().SetLogger(c.Request.Context(), reqLogger))

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// log request body
		var body []byte
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ = io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)

		fields := []zapcore.Field{
			zap.String("requestId", requestID),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.String("body", string(body)),
			zap.String("time", end.Format(time.RFC3339)),
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e, fields...)
			}
		} else {
			logger.Info(path, fields...)
		}
	}
}

func RecoveryWithZap(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.String("requestId", foundation.ContextProvider().GetRequestId(c.Request.Context())),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) //nolint: errcheck
					c.Abort()
					return
				}

				logger.Error("[Recovery from panic]",
					zap.String("requestId", foundation.ContextProvider().GetRequestId(c.Request.Context())),
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				// recovery
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
