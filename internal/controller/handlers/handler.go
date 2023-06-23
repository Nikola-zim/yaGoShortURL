package handlers

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"yaGoShortURL/internal/entity"
)

type addAndGetURL interface {
	addURL(c *gin.Context)
	getURL(c *gin.Context)
	addAndGetJSON(c *gin.Context)
	addAndGetBatchURL(c *gin.Context)
}

type addAndGetURLService interface {
	WriteURL(fullURL string, userIDB []byte) (string, error)
	FullURL(id string) (string, error)
	ReadAllUserURLFromCash(id []byte) ([]entity.JSONAllInfo, error)
}

type authUser interface {
	FindUser(idMsg string) (uint64, bool)
	AddUser() (string, uint64, error)
}

type DBService interface {
	PingDB() error
}

// Cache Собранный интерфейс для кэша
type Cache interface {
	addAndGetURLService
	authUser
	DBService
}

// Добавление id пользователя (запись кук)
type authorizationService interface {
	AddUser() (string, error)
	FindUser()
}

type Handler struct {
	addAndGetURL
	UserInteract
	Postgres
}

func NewHandler(cache Cache, baseURL string) *Handler {
	return &Handler{
		addAndGetURL: NewAddAndGetURLHandler(cache, baseURL),
		UserInteract: *NewUserInteract(cache),
		Postgres:     *NewPostgres(cache),
	}
}

// Middleware для сжатия запросов и ответов
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
			c.Request.Body = reader
			c.Request.Header.Del("Content-Encoding")
			// Передача запроса в handler
			c.Next()
			return
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
			// Передача запроса в handler
			c.Next()
			return
		}
		// Передача запроса в handler
		c.Next()
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
	shortenerURL.Use(h.cookieSetAndGet())
	{
		shortenerURL.POST("/", h.addURL)
		shortenerURL.GET("/:id", h.getURL)
	}

	// Получение сокращенного URL
	shorten := router.Group("/api/")
	// использование middleware для сжатия запросов
	shorten.Use(gzipHandle())
	shorten.Use(h.cookieSetAndGet())
	{
		shorten.POST("shorten", h.addAndGetJSON)
		shorten.GET("user/urls", h.getAllUserURL)
		shorten.POST("shorten/batch", h.addAndGetBatchURL)
	}

	// батчинг
	shortenBatching := router.Group("/ping")
	{
		shortenBatching.GET("", h.pingDB)
	}

	return router
}
