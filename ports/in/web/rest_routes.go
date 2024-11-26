package web

import (
	"contentgit/app"
	"contentgit/appservices"
	"github.com/gin-gonic/gin"
)

type Router struct {
}

func (r Router) MapRoutes(registry *app.ComponentRegistry, routerGroup *gin.RouterGroup) {
	NewContentController(routerGroup, registry.Get("ContentService").(*appservices.ContentService),
		registry.Get("ContentQuery").(*appservices.ContentQuery)).MapRoutes()
}
