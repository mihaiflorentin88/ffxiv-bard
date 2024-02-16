package contract

import "github.com/gin-gonic/gin"

type RouterInterface interface {
	RegisterRoutes(router *gin.Engine)
}
