package utils

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"slices"
)

//go:embed resource/*
var staticFS embed.FS

func GetStaticFS() *embed.FS {
	return &staticFS
}

type inMemoryAssets struct {
	css map[string][]byte
	js  map[string][]byte
	img map[string][]byte
}

func loadAssets(root string) map[string][]byte {
	assets := make(map[string][]byte)
	err := fs.WalkDir(staticFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			content, readErr := fs.ReadFile(staticFS, path)
			if readErr != nil {
				return readErr
			}
			// Adjust the path according to your routing needs
			internalPath := path[len(root):]
			assets[internalPath] = content
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to load assets: %v", err)
	}
	return assets
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

type Renderer struct {
	StaticFS     *embed.FS
	Templates    []string
	ErrorHandler ErrorHandler
}

type RenderedData struct {
	User interface{} // can be either user.User{} or nil
	Data interface{} // can be anything
}

func NewRenderer(errorHandler ErrorHandler) Renderer {
	r := Renderer{
		StaticFS:     GetStaticFS(),
		ErrorHandler: errorHandler,
	}
	r.Templates = getDefaultTemplates()
	return r
}

func (r *Renderer) AddTemplate(file string) contract.HttpRenderer {
	r.Templates = append(r.Templates, file)
	return r
}

func (r *Renderer) StartClean() contract.HttpRenderer {
	r.Templates = getDefaultTemplates()
	return r
}

func (r *Renderer) RemoveTemplate(file string) contract.HttpRenderer {
	for i, f := range r.Templates {
		if f == file {
			r.Templates = slices.Delete(r.Templates, i, i+1)
			break
		}
	}
	return r
}

func (r *Renderer) Render(c *gin.Context, data interface{}, StatusCode int) {
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
		return
	}
}

func (r *Renderer) RegisterStatic(router *gin.Engine) {
	cssAssets := loadAssets("resource/css")
	jsAssets := loadAssets("resource/js")
	imgAssets := loadAssets("resource/img")

	router.GET("/_resource/css/*filepath", func(c *gin.Context) {
		filepathStr := c.Param("filepath")
		content, ok := cssAssets[filepathStr]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/css", content)
	})

	router.GET("/_resource/js/*filepath", func(c *gin.Context) {
		filepathStr := c.Param("filepath")
		content, ok := jsAssets[filepathStr]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "application/javascript", content)
	})

	router.GET("/_resource/img/*filepath", func(c *gin.Context) {
		filepathStr := c.Param("filepath")
		content, ok := imgAssets[filepathStr]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		mimeType := "application/octet-stream"

		switch filepath.Ext(filepathStr) {
		case ".png":
			mimeType = "image/png"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
			// Add more cases as needed for different file types
		}
		c.Data(http.StatusOK, mimeType, content) // Use the correct MIME type
	})
}

//func (r *Renderer) RegisterStatic(router *gin.Engine) {
//	cssFS, err := fs.Sub(r.StaticFS, "resource/css")
//	if err != nil {
//		panic("Cannot parse the css")
//	}
//	jsFS, err := fs.Sub(r.StaticFS, "resource/js")
//	if err != nil {
//		panic("Cannot parse the js")
//	}
//	imgFS, err := fs.Sub(r.StaticFS, "resource/img")
//	if err != nil {
//		panic("Cannot parse the img")
//	}
//
//	router.StaticFS("/_resource/css", http.FS(cssFS))
//	router.StaticFS("/_resource/js", http.FS(jsFS))
//	router.StaticFS("/_resource/img", http.FS(imgFS))
//}
