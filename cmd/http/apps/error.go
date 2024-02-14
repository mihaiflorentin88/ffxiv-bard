package httpapps

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"runtime/debug"
	"strings"
)

type HttpError struct {
	Message   string
	Traceback []string
}

func NewHttpError() contract.HttpError {
	return &HttpError{}
}

func (h *HttpError) GetMessage() string {
	return h.Message
}

func (h *HttpError) SetMessage(message string) {
	h.Message = message
}

func (h *HttpError) GetTraceback() []string {
	return h.Traceback
}

func (h *HttpError) SetTraceback(traceback []string) {
	h.Traceback = traceback
}

func (h *HttpError) RenderTemplate(err error, c *gin.Context, staticFS *embed.FS) {
	traceback := strings.Split(string(debug.Stack()), "\n")
	h.SetMessage(err.Error())
	h.SetTraceback(traceback)
	tmpl, err := template.New("base").ParseFS(
		staticFS,
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
