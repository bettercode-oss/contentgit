package foundation

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	ginErrorHandlerOnce     sync.Once
	ginErrorHandlerInstance *ginErrorHandler
)

func GinErrorHandler() *ginErrorHandler {
	ginErrorHandlerOnce.Do(func() {
		ginErrorHandlerInstance = &ginErrorHandler{}
	})

	return ginErrorHandlerInstance
}

type ginErrorHandler struct {
}

func (ginErrorHandler) InternalServerError(ctx *gin.Context, err error) {
	ctx.Error(err)
	ctx.Status(http.StatusInternalServerError)
}
