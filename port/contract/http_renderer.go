package contract

import "github.com/gin-gonic/gin"

type HttpRenderer interface {
	AddTemplate(file string) HttpRenderer
	RemoveTemplate(file string) HttpRenderer
	Render(c *gin.Context, data interface{}, StatusCode int)
	RegisterStatic(router *gin.Engine)
	StartClean() HttpRenderer
}
