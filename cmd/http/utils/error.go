package utils

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"runtime/debug"
	"strings"
)

type ErrorHandler struct {
	Message   string
	Traceback []string
	StaticFS  *embed.FS
}

func NewHttpErrorHandler() contract.HttpErrorHandlerInterface {
	return &ErrorHandler{StaticFS: &staticFS}
}

func (h *ErrorHandler) GetMessage() string {
	return h.Message
}

func (h *ErrorHandler) SetMessage(message string) {
	h.Message = message
}

func (h *ErrorHandler) GetTraceback() []string {
	return h.Traceback
}

func (h *ErrorHandler) SetTraceback(traceback []string) {
	h.Traceback = traceback
}

func (h *ErrorHandler) RenderTemplate(err error, c *gin.Context) {
	traceback := strings.Split(string(debug.Stack()), "\n")
	h.SetMessage(err.Error())
	h.SetTraceback(traceback)
	tmpl, err := template.New("base").ParseFS(
		h.StaticFS,
		"resource/template/base/base.html",
		"resource/template/base/navbar.html",
		"resource/template/error/error.html",
		"resource/template/base/base_js.html",
		"resource/template/base/base_styles.html",
		"resource/template/base/additional_js.html",
		"resource/template/base/additional_styles.html",
	)
	if err != nil {
		println("error parsing templates from FS: %s", err)
	}
	err = tmpl.ExecuteTemplate(c.Writer, "base", h)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
