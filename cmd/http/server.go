package http

import (
	"ffxvi-bard/container"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

func RegisterRoutes(router *gin.Engine, serviceContainer *container.ServiceContainer) {
	httpRenderer := serviceContainer.GetHttpRenderer()
	httpRenderer.RegisterStatic(router)
	serviceContainer.GetMainRouter().RegisterRoutes(router)
	serviceContainer.GetAuthRouter().RegisterRoutes(router)
	serviceContainer.GetSongRouter().RegisterRoutes(router)
}

func StartServer(port int, poolSize int, serviceContainer *container.ServiceContainer) {
	runtime.GOMAXPROCS(poolSize)
	router := serviceContainer.GetGinRouter()
	authMiddleware := serviceContainer.GetAuthMiddleware()
	router.Use(authMiddleware.GetLoggedUser())
	RegisterRoutes(router, serviceContainer)
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
