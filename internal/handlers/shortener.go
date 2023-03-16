package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"yaGoShortURL/internal/service"
)

type AddAndGetURLHandler struct {
	service service.CashURL
}

func (a *AddAndGetURLHandler) addURL(c *gin.Context) {
	//Читаем Body
	b, err := c.GetRawData()
	if err != nil || string(b) == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	//Запись в память
	id, err := a.service.WriteURLInCash(string(b))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id = "http://localhost:8080/" + id
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

func NewAddAndGetURLHandler(service service.CashURL) *AddAndGetURLHandler {
	return &AddAndGetURLHandler{service: service}
}
