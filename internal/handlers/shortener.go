package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"yaGoShortURL/internal/static"
)

type AddAndGetURLHandler struct {
	service addAndGetURLService
	baseURL string
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
	// Получение userIdB
	cookie, _ := c.Cookie("user_id")
	data, err := hex.DecodeString(cookie)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	// data[:8] - байты id-шника
	id, err := a.service.WriteURLInCash(string(b), data[:8])
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Получение короткого адреса
	id = fmt.Sprintf("%s%s%s", a.baseURL, "/", id)
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
	log.Println(json.Unmarshal(b, &result))
	if err != nil {
		fmt.Println(err)
	}
	//Запись в кеш
	id, err := a.service.WriteURLInCash(myJSON.URL, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Вывод результата
	result.Res = fmt.Sprintf("%s%s%s", a.baseURL, "/", id)
	c.JSON(http.StatusCreated, result)
}

func NewAddAndGetURLHandler(service addAndGetURLService, baseURL string) *AddAndGetURLHandler {
	return &AddAndGetURLHandler{
		service: service,
		baseURL: baseURL,
	}
}
