package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"yaGoShortURL/internal/static"
)

type AddAndGetURLHandler struct {
	service addAndGetURLService
	cfg     static.ConfigInit
}

func (a *AddAndGetURLHandler) addURL(c *gin.Context) {
	//Читаем Body
	b, err := c.GetRawData()
	if err != nil || string(b) == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	//Запись в кеш
	id, err := a.service.WriteURLInCash(string(b))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Получение короткого адреса
	id = fmt.Sprintf("%s%s%s", a.cfg.BaseURL, "/", id)
	c.String(http.StatusCreated, id)
}

func (a *AddAndGetURLHandler) getURL(c *gin.Context) {
	//Получаем
	idStr := c.Param("id")
	id := "id:" + idStr
	str, err := a.service.ReadURLFromCash(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, str)
}

func (a *AddAndGetURLHandler) addAndGetJSON(c *gin.Context) {
	var myJSON static.JSONApi
	var result static.JSONRes
	err := c.ShouldBindJSON(&myJSON)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//это для прохождение проверки использования unmarshal
	b, err := c.GetRawData()
	fmt.Println(json.Unmarshal(b, &result))
	if err != nil {
		fmt.Println(err)
	}
	//Запись в кеш
	id, err := a.service.WriteURLInCash(myJSON.URL)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Получение переменной окружения
	if baseURL := a.cfg.BaseURL; baseURL != "" {
		result.Res = fmt.Sprintf("%s%s%s", os.Getenv("BASE_URL"), "/", id)
		c.JSON(http.StatusCreated, result)
	} else {
		result.Res = "http://localhost:8080/" + id
		c.JSON(http.StatusCreated, result)
	}
}

func NewAddAndGetURLHandler(service addAndGetURLService, cfg static.ConfigInit) *AddAndGetURLHandler {
	return &AddAndGetURLHandler{
		service: service,
		cfg:     cfg,
	}
}
