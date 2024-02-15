package utils

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"slices"
)

//go:embed resource/*
var staticFS embed.FS

func GetStaticFS() *embed.FS {
	return &staticFS
}

type renderer struct {
	StaticFS     *embed.FS
	Templates    []string
	ErrorHandler contract.HttpErrorHandlerInterface
}

func NewRenderer(errorHandler contract.HttpErrorHandlerInterface) contract.HttpRenderer {
	r := &renderer{
		StaticFS:     GetStaticFS(),
		ErrorHandler: errorHandler,
	}
	defaultTemplates := []string{
		"resource/template/base/base.html",
		"resource/template/base/navbar.html",
		"resource/template/base/base_js.html",
		"resource/template/base/base_styles.html",
		"resource/template/base/additional_js.html",
		"resource/template/base/additional_styles.html",
	}
	r.Templates = defaultTemplates
	return r
}

func (r *renderer) AddTemplate(file string) contract.HttpRenderer {
	r.Templates = append(r.Templates, file)
	return r
}

func (r *renderer) RemoveTemplate(file string) contract.HttpRenderer {
	for i, f := range r.Templates {
		if f == file {
			r.Templates = slices.Delete(r.Templates, i, i+1)
			break
		}
	}
	return r
}

func (r *renderer) Render(c *gin.Context, data interface{}) {
	tmpl := template.New("base")
	var err error
	for _, file := range r.Templates {
		tmpl, err = tmpl.ParseFS(r.StaticFS, file)
		if err != nil {
			log.Println("Error parsing templates from FS. Reason: ", err)
			r.ErrorHandler.RenderTemplate(err, c)
		}
	}
	err = tmpl.ExecuteTemplate(c.Writer, "base", data)
	if err != nil {
		r.ErrorHandler.RenderTemplate(err, c)
		log.Println("Could not render template. Reason: ", err)
	}
}

func (r *renderer) EnableStatic(router *gin.Engine) {
	cssFS, err := fs.Sub(r.StaticFS, "resource/css")
	if err != nil {
		panic("Cannot parse the css")
	}
	jsFS, err := fs.Sub(r.StaticFS, "resource/js")
	if err != nil {
		panic("Cannot parse the js")
	}
	imgFS, err := fs.Sub(r.StaticFS, "resource/img")
	if err != nil {
		panic("Cannot parse the img")
	}

	router.StaticFS("/_resource/css", http.FS(cssFS))
	router.StaticFS("/_resource/js", http.FS(jsFS))
	router.StaticFS("/_resource/img", http.FS(imgFS))
}
