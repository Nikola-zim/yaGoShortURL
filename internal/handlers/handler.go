package handlers

import (
	"github.com/gin-gonic/gin"
	"yaGoShortURL/internal/service"
)

type addAndGetURL interface {
	addURL(c *gin.Context)
	getURL(c *gin.Context)
}

type Handler struct {
	addAndGetURL
}

func NewHandler(addAndGetURL service.CashURL) *Handler {
	return &Handler{addAndGetURL: NewAddAndGetURLHandler(addAndGetURL)}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	shortenerURL := router.Group("/")
	{
		shortenerURL.POST("/", h.addURL)
		shortenerURL.GET("/:id", h.getURL)
	}
	return router
}
