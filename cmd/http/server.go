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
	container.GetHttpRenderer().RegisterStatic(router)
	container.GetMainRouter().RegisterRoutes(router)
	container.GetAuthRouter().RegisterRoutes(router)
	container.GetSongRouter().RegisterRoutes(router)
}

func Server(port int, poolSize int) {
	runtime.GOMAXPROCS(poolSize)
	router := container.GetGinRouter()
	authMiddleware := container.GetAuthMiddleware()
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
