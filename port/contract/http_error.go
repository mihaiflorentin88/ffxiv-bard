package contract

import (
	"github.com/gin-gonic/gin"
)

type HttpErrorHandlerInterface interface {
	GetMessage() string
	SetMessage(message string)
	GetTraceback() []string
	SetTraceback(traceback []string)
	RenderTemplate(err error, c *gin.Context)
}
