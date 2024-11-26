package middlewares

import (
	"contentgit/foundation"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GORMDb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		ctx := req.Context()
		c.Request = c.Request.WithContext(foundation.ContextProvider().SetDB(ctx, db))
		c.Next()
	}
}
