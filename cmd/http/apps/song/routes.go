package song

import (
	"ffxvi-bard/cmd/http/utils/middleware"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller    *Controller
	JWTMiddleware middleware.AuthMiddleware
}

func NewSongRouter(controller *Controller, JWTMiddleware middleware.AuthMiddleware) contract.RouterInterface {
	return &Router{
		Controller:    controller,
		JWTMiddleware: JWTMiddleware,
	}
}

func (r Router) RegisterRoutes(router *gin.Engine) {
	router.GET("/song/:id", r.Controller.SongDetails)
	router.GET("/song/list", r.Controller.RenderSongList)
	router.GET("/song/add", r.JWTMiddleware.UI(), func(c *gin.Context) {
		r.Controller.RenderAddNewSongForm(c)
	})
	router.POST("/song", r.Controller.HandleAddNewSong)
	router.GET("/song/download/:id", r.Controller.DownloadSong)
	router.POST("/song/:songID/rating", r.JWTMiddleware.UI(), func(c *gin.Context) {
		r.Controller.SubmitSongRating(c)
	})
	router.POST("/song/:songID/comment", r.JWTMiddleware.UI(), func(c *gin.Context) {
		r.Controller.SubmitSongComment(c)
	})
	router.PUT("/song/:songID/comment", r.JWTMiddleware.UI(), func(c *gin.Context) {
		r.Controller.SubmitSongCommentUpdate(c)
	})
}
