package contract

import (
	"embed"
	"github.com/gin-gonic/gin"
)

type HttpError interface {
	GetMessage() string
	SetMessage(message string)
	GetTraceback() []string
	SetTraceback(traceback []string)
	RenderTemplate(err error, c *gin.Context, staticFS *embed.FS)
}
