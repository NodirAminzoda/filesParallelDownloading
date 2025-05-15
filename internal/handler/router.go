package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	middlewares "parallel_download_from_many_urls/internal/middleware"
	"parallel_download_from_many_urls/internal/service"

	"parallel_download_from_many_urls/internal/config"
)

//type IHandler interface{
//	DownloadFiles(c *gin.Context)
//}

type Handler struct {
	cfg     *config.Config
	log     *zerolog.Logger
	service service.IService
}

func New(cfg *config.Config, log *zerolog.Logger, service service.IService) *Handler {
	return &Handler{
		cfg:     cfg,
		log:     log,
		service: service,
	}
}

func (h *Handler) InitRoute() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(middlewares.SetRequestID())
	router.POST("/download/files", h.DownloadFiles)
	router.POST("download/file", h.DownloadFile)
	
	return router
}
