package app

import (
	"github.com/gin-gonic/gin"
)

type GinRoute interface {
  MapRoutes(registry *ComponentRegistry, routerGroup *gin.RouterGroup)
}
