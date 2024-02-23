package utils

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

type ErrorHandler struct {
	User      interface{}
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

func (h *ErrorHandler) RenderTemplate(err error, statusCode int, c *gin.Context) {
	traceback := strings.Split(string(debug.Stack()), "\n")
	loggedUser, exists := c.Get("user")
	if exists && loggedUser != nil {
		h.User = loggedUser
	}
	h.SetMessage(err.Error())
	h.SetTraceback(traceback)
	tmpl, err := template.New("base").ParseFS(
		h.StaticFS,
		"resource/template/base/base.gohtml",
		"resource/template/base/navbar.gohtml",
		"resource/template/error/error.gohtml",
		"resource/template/base/base_js.gohtml",
		"resource/template/base/base_styles.gohtml",
		"resource/template/base/additional_js.gohtml",
		"resource/template/error/error_css.gohtml",
	)
	if err != nil {
		statusCode = http.StatusInternalServerError
		println("error parsing templates from FS: %s", err)
	}
	c.Status(statusCode)
	c.Header("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(c.Writer, "base", h)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
