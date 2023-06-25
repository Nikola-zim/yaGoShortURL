package handlers

import (
	"github.com/gin-gonic/gin"
	"yaGoShortURL/internal/controller/middlewares"
	"yaGoShortURL/internal/entity"
)

type addAndGetURL interface {
	addURL(c *gin.Context)
	getURL(c *gin.Context)
	addAndGetJSON(c *gin.Context)
	addAndGetBatchURL(c *gin.Context)
}

type addAndGetURLService interface {
	WriteURL(fullURL string, id uint64) (string, error)
	FullURL(id string) (string, error)
	ReadAllUserURL(id uint64) ([]entity.JSONAllInfo, error)
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
	middlewares.UserInteract
	Postgres
}

func NewHandler(cache Cache, baseURL string) *Handler {
	return &Handler{
		addAndGetURL: NewAddAndGetURLHandler(cache, baseURL),
		UserInteract: *middlewares.NewUserInteract(cache),
		Postgres:     *NewPostgres(cache),
	}
}

// InitRoutes Хендлеры
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	shortenerURL := router.Group("/")
	// использование middleware для сжатия запросов
	shortenerURL.Use(middlewares.GzipHandle())
	shortenerURL.Use(h.UserInteract.CookieSetAndGet())
	{
		shortenerURL.POST("/", h.addURL)
		shortenerURL.GET("/:id", h.getURL)
	}

	// Получение сокращенного URL
	shorten := router.Group("/api/")
	// использование middleware для сжатия запросов
	shorten.Use(middlewares.GzipHandle())
	shorten.Use(h.UserInteract.CookieSetAndGet())
	{
		shorten.POST("shorten", h.addAndGetJSON)
		shorten.GET("user/urls", h.UserInteract.GetAllUserURL)
		shorten.POST("shorten/batch", h.addAndGetBatchURL)
	}

	// Батчинг
	shortenBatching := router.Group("/ping")
	{
		shortenBatching.GET("", h.pingDB)
	}

	return router
}
