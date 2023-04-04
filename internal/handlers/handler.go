package handlers

import (
	"bytes"
	"compress/flate"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
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

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "deflate" {
			// Читаем тело запроса в сжатом формате
			reader := flate.NewReader(c.Request.Body)
			defer reader.Close()

			// Заменяем тело запроса на распакованные данные
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Fatal(err)
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

			// Устанавливаем заголовок Content-Encoding в значение "identity"
			c.Request.Header.Set("Content-Encoding", "identity")
		}

		c.Next()
	})
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
