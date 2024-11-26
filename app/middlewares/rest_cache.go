package middlewares

import (
	"bytes"
	"contentgit/app/cache"
	"contentgit/foundation"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const cachePrefix = "web:"

func HttpResponseCache(cache cache.Cache, maxAge uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := foundation.ContextProvider().GetLogger(c.Request.Context())
		requestKey := fmt.Sprintf("%s:%s:%s", cachePrefix, c.Request.Header.Get("X-ApplicationId"), c.Request.RequestURI)
		if c.Request.Method == http.MethodPost {
			// POST의 경우 body 정보까지 key 로 만든다
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)

			requestKey += ":" + string(body)
		}

		cacheKey := generateMD5Hash(requestKey)
		cachedResponse, found := cache.Get(cacheKey)
		if found {
			logger.Info("web cache hit", zap.String("key", cacheKey))
			c.JSON(http.StatusOK, cachedResponse)
			c.Abort()
			return
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		statusCode := c.Writer.Status()
		if statusCode == 200 {
			responseJson := map[string]any{}
			if err := json.Unmarshal(blw.body.Bytes(), &responseJson); err != nil {
				logger.Error("http response cache set error", zap.Error(err))
			} else {
				cache.Set(cacheKey, responseJson, time.Duration(maxAge)*time.Minute)
			}
		}
	}
}

func generateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
