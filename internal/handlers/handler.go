package handlers

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
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
func gzipHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the client supports gzip encoding
		if !strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			c.Next()
			return
		}
		// Create a gzip writer
		reader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		defer reader.Close()
		uncompressed, err := io.ReadAll(reader)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(string(uncompressed)))
		c.Request.Header.Del("Content-Encoding")
		c.Request.Header.Set("Content-Type", "application/json")
	}
}

// InitRoutes Хендлеры
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	shortenerURL := router.Group("/")
	{
		shortenerURL.POST("/", h.addURL)
		shortenerURL.GET("/:id", h.getURL)
	}
	// использование middleware для сжатия запросов
	shortenerURL.Use(gzipHandle())
	shorten := router.Group("/api/")
	{
		shorten.POST("shorten", h.addAndGetJSON)
	}
	// использование middleware для сжатия запросов
	shorten.Use(gzipHandle())
	return router
}
