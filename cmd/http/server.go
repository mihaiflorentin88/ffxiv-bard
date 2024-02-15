package http

import (
	"ffxvi-bard/container"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

func EnabledRoutes(router *gin.Engine) {
	container.GetHttpRenderer().EnableStatic(router)
	container.GetSongRouter().EnableRoutes(router)
	container.GetMainRouter().EnableRoutes(router)
	container.GetAuthRouter().EnableRoutes(router)
}

func Server(port int, poolSize int) {
	runtime.GOMAXPROCS(poolSize)
	router := container.GetGinRouter()
	authMiddleware := container.GetAuthMiddleware()
	router.Use(authMiddleware.GetLoggedUser())
	EnabledRoutes(router)
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
