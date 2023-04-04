package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmu42/gzip"
	"yaGoShortURL/internal/static"
)

type addAndGetURL interface {
	addURL(c *gin.Context)
	getURL(c *gin.Context)
	addAndGetJSON(c *gin.Context)
}

type addAndGetURLService interface {
	WriteURLInCash(string2 string) (string, error)
	ReadURLFromCash(string string) (string, error)
}

type Handler struct {
	addAndGetURL
}

func NewHandler(service addAndGetURLService, cfg static.ConfigInit) *Handler {
	return &Handler{
		addAndGetURL: NewAddAndGetURLHandler(service, cfg),
	}
}

// Middleware

// InitRoutes Хендлеры
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	shortenerURL := router.Group("/")
	{
		shortenerURL.POST("/", h.addURL)
		shortenerURL.GET("/:id", h.getURL)
	}
	// использование middleware для сжатия запросов
	shortenerURL.Use(gzip.DefaultHandler().Gin)
	shorten := router.Group("/api/")
	{
		shorten.POST("shorten", h.addAndGetJSON)
	}
	// использование middleware для сжатия запросов
	shorten.Use(gzip.DefaultHandler().Gin)
	return router
}
