package http

import (
	"ffxvi-bard/container"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

func RegisterRoutes(router *gin.Engine) {
	httpRenderer := container.Load.HttpRenderer()
	httpRenderer.RegisterStatic(router)
	container.Load.MainRouter().RegisterRoutes(router)
	container.Load.AuthRouter().RegisterRoutes(router)
	container.Load.SongRouter().RegisterRoutes(router)
}

func StartServer(port int, poolSize int) {
	runtime.GOMAXPROCS(poolSize)
	router := container.Load.GinRouter()
	authMiddleware := container.Load.AuthMiddleware()
	router.Use(authMiddleware.GetLoggedUser())
	RegisterRoutes(router)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	err := s.ListenAndServe()
	if err != nil {
		panic("Cannot start the http server. Reason: " + err.Error())
	}
}
