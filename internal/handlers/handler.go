package handlers

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
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

func NewHandler(service addAndGetURLService, baseURL string) *Handler {
	return &Handler{
		addAndGetURL: NewAddAndGetURLHandler(service, baseURL),
	}
}

// Middleware
func gzipHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content-Encoding
		// проверяем, что запрос отправлен в формате gzip
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			// Create a gzip writer
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer func(reader *gzip.Reader) {
				err := reader.Close()
				if err != nil {
					log.Println(err)
				}
			}(reader)
			uncompressed, err := io.ReadAll(reader)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Request.Body = io.NopCloser(strings.NewReader(string(uncompressed)))
			c.Request.Header.Del("Content-Encoding")
		}
		// Accept-Encoding
		// проверяем, что клиент поддерживает gzip-сжатие
		if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			// создаём gzip.Writer поверх текущего w
			gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer func(gz *gzip.Writer) {
				err := gz.Close()
				if err != nil {
					log.Println(err)
				}
			}(gz)

			//Set the Content-Encoding header
			c.Header("Content-Encoding", "gzip")

			//Replace the writer with gzip writer
			c.Writer = &gzipWriter{Writer: gz, ResponseWriter: c.Writer}
		}
		// Передача запроса в handler
		c.Next()
		return
	}

}

// опишем тип gzipWriter, поддерживающий интерфейс http.ResponseWriter, и реализуем недостающие методы.
type gzipWriter struct {
	gin.ResponseWriter
	*gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func (g *gzipWriter) Flush() {
	err := g.Writer.Flush()
	if err != nil {
		log.Println(err)
		return
	}
	g.ResponseWriter.Flush()
}

func (g *gzipWriter) Header() http.Header {
	return g.ResponseWriter.Header()
}

// InitRoutes Хендлеры
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	shortenerURL := router.Group("/")
	// использование middleware для сжатия запросов
	shortenerURL.Use(gzipHandle())
	{
		shortenerURL.POST("/", h.addURL)
		shortenerURL.GET("/:id", h.getURL)
	}

	shorten := router.Group("/api/")
	// использование middleware для сжатия запросов
	shorten.Use(gzipHandle())
	{
		shorten.POST("shorten", h.addAndGetJSON)
	}
	return router
}
