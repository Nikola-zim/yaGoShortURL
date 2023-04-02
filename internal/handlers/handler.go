package handlers

import (
	"github.com/gin-gonic/gin"
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

func NewHandler(service addAndGetURLService) *Handler {
	return &Handler{
		addAndGetURL: NewAddAndGetURLHandler(service),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	shortenerURL := router.Group("/")
	{
		shortenerURL.POST("/", h.addURL)
		shortenerURL.GET("/:id", h.getURL)
	}
	shorten := router.Group("/api/")
	{
		shorten.POST("shorten", h.addAndGetJSON)
	}
	return router
}
