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

func getDefaultTemplates() []string {
	return []string{
		"resource/template/base/base.gohtml",
		"resource/template/base/navbar.gohtml",
		"resource/template/base/base_js.gohtml",
		"resource/template/base/base_styles.gohtml",
		"resource/template/base/additional_js.gohtml",
		"resource/template/base/additional_styles.gohtml",
	}
}

type renderer struct {
	StaticFS     *embed.FS
	Templates    []string
	ErrorHandler contract.HttpErrorHandlerInterface
}

type RenderedData struct {
	User interface{} // can be either user.User{} or nil
	Data interface{} // can be anything
}

func NewRenderer(errorHandler contract.HttpErrorHandlerInterface) contract.HttpRenderer {
	r := &renderer{
		StaticFS:     GetStaticFS(),
		ErrorHandler: errorHandler,
	}
	r.Templates = getDefaultTemplates()
	return r
}

func (r *renderer) AddTemplate(file string) contract.HttpRenderer {
	r.Templates = append(r.Templates, file)
	return r
}

func (r *renderer) StartClean() contract.HttpRenderer {
	r.Templates = getDefaultTemplates()
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

func (r *renderer) Render(c *gin.Context, data interface{}, StatusCode int) {
	tmpl := template.New("base")
	var err error
	for _, file := range r.Templates {
		tmpl, err = tmpl.ParseFS(r.StaticFS, file)
		if err != nil {
			log.Println("Error parsing templates from FS. Reason: ", err)
			r.ErrorHandler.RenderTemplate(err, http.StatusInternalServerError, c)
			return
		}
	}

	renderedData := RenderedData{
		Data: data,
	}

	loggedUser, exists := c.Get("user")
	if exists && loggedUser != nil {
		renderedData.User = loggedUser
	}
	c.Status(StatusCode)
	err = tmpl.ExecuteTemplate(c.Writer, "base", renderedData)
	if err != nil {
		log.Println("Could not render template. Reason: ", err)
		r.ErrorHandler.RenderTemplate(err, http.StatusInternalServerError, c)
	}
}

func (r *renderer) RegisterStatic(router *gin.Engine) {
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
