package handlers

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"yaGoShortURL/internal/entity"
)

type AddAndGetURLHandler struct {
	service addAndGetURLService
	baseURL string
}

func NewAddAndGetURLHandler(service addAndGetURLService, baseURL string) *AddAndGetURLHandler {
	return &AddAndGetURLHandler{
		service: service,
		baseURL: baseURL,
	}
}

func (a *AddAndGetURLHandler) getUserID(c *gin.Context) uint64 {
	cookie, err := c.Cookie("user_id")
	data := make([]byte, 8, 39)
	// Ошибка означает что куки не было, и нужно взять ID, который установили в запросе
	if err != nil {
		userID, _ := c.Get("user_ID")
		switch t := userID.(type) {
		case uint64:
			ID := reflect.ValueOf(t).Uint()
			//Если UserID был установлен, т.е. кука была только получена
			if userID != 0 {
				data = make([]byte, 8)
				binary.LittleEndian.PutUint64(data, ID)
			}
		}
	} else {
		data, err = hex.DecodeString(cookie)
		log.Printf("ID:%v", cookie)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	// Получение id в виде числа
	userID := binary.LittleEndian.Uint64(data[:8])

	return userID
}

func (a *AddAndGetURLHandler) addURL(c *gin.Context) {
	//Читаем Body
	b, err := c.GetRawData()
	if err != nil || string(b) == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Получение user id
	userID := a.getUserID(c)

	id, err := a.service.WriteURL(string(b), userID)
	if err != nil {
		log.Println(err)
		var eu *entity.ErrorURL
		if errors.As(err, &eu) {
			id = fmt.Sprintf("%s%s%s", a.baseURL, "/", err.Error())
			c.String(http.StatusConflict, id)
			return
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	// Получение короткого адреса
	id = fmt.Sprintf("%s%s%s", a.baseURL, "/", id)
	c.String(http.StatusCreated, id)
}

func (a *AddAndGetURLHandler) getURL(c *gin.Context) {
	//Получаем
	id := c.Param("id")
	str, err := a.service.FullURL(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, str)
}

func (a *AddAndGetURLHandler) addAndGetJSON(c *gin.Context) {
	var myJSON entity.JSONApi
	var result entity.JSONRes
	err := c.ShouldBindJSON(&myJSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//это для прохождения проверки использования unmarshal
	b, err := c.GetRawData()
	log.Println(json.Unmarshal(b, &result))
	if err != nil {
		fmt.Println(err)
	}
	// Получение user id
	userID := a.getUserID(c)

	//Запись в кеш
	id, err := a.service.WriteURL(myJSON.URL, userID)
	if err != nil {
		log.Println(err)
		var eu *entity.ErrorURL
		if errors.As(err, &eu) {
			result.Res = fmt.Sprintf("%s%s%s", a.baseURL, "/", err.Error())
			c.JSON(http.StatusConflict, result)
			return
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	// Вывод результата
	result.Res = fmt.Sprintf("%s%s%s", a.baseURL, "/", id)
	c.JSON(http.StatusCreated, result)
}

func (a *AddAndGetURLHandler) addAndGetBatchURL(c *gin.Context) {
	var myJSON []entity.BatchAPI
	var result entity.JSONRes
	err := c.ShouldBindJSON(&myJSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//это для прохождения проверки использования unmarshal
	b, err := c.GetRawData()
	log.Println(json.Unmarshal(b, &result))
	if err != nil {
		fmt.Println(err)
	}
	// Получение user id
	userID := a.getUserID(c)

	//Запись в кеш и формирование результата
	res := make([]entity.ResBatchAPI, 0, 100)
	for _, url := range myJSON {
		var ans entity.ResBatchAPI
		id, err := a.service.WriteURL(url.OriginalURL, userID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ans.ShortURL = fmt.Sprintf("%s%s%s", a.baseURL, "/", id)
		ans.CorrelationID = url.CorrelationID
		res = append(res, ans)
	}
	// Вывод результата
	c.JSON(http.StatusCreated, res)
}
