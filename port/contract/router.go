package contract

import "github.com/gin-gonic/gin"

type RouterInterface interface {
	EnableRoutes(router *gin.Engine)
}
