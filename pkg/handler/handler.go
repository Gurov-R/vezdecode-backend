package handler

import (
	"Gurov-R/vezdecode-backend/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(CORSMiddleware())

	api := router.Group("/api")
	{
		memes := api.Group("/memes")
		{
			memes.GET("/", h.GetAllMemes)
			memes.POST("/load-vezdekod", h.LoadVezdekod)
			memes.POST("/load-group", h.LoadGroup)
			memes.POST("/feed", h.Feed)
			memes.POST("/like", h.Like)
			memes.POST("/promote", h.Promote)
		}
	}

	return router
}
