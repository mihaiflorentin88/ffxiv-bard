package httpapps

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"slices"
)

type renderer struct {
	StaticFS     *embed.FS
	Templates    []string
	ErrorHandler contract.HttpError
}

func NewRenderer(errorHandler contract.HttpError, staticFS *embed.FS) contract.HttpRenderer {
	r := &renderer{
		StaticFS:     staticFS,
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
			r.ErrorHandler.RenderTemplate(err, c, r.StaticFS)
		}
	}
	err = tmpl.ExecuteTemplate(c.Writer, "base", data)
	if err != nil {
		r.ErrorHandler.RenderTemplate(err, c, r.StaticFS)
		log.Println("Could not render template. Reason: ", err)
	}
}
