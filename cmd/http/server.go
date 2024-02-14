package http

import (
	"embed"
	httpapps "ffxvi-bard/cmd/http/apps"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

//go:embed resource/*
var staticFS embed.FS

func Server(port int, poolSize int) {
	runtime.GOMAXPROCS(poolSize)
	router := gin.Default()
	app := httpapps.NewRoot(router, &staticFS)
	app.Initialize()
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
