package middlewares

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"yaGoShortURL/internal/usecase"
)

type UserInteract struct {
	service usecase.Cache
}

func NewUserInteract(service usecase.Cache) *UserInteract {
	return &UserInteract{
		service: service,
	}
}

// Middleware для cookie
func (uI *UserInteract) CookieSetAndGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user_id")
		if err != nil {
			// Добавим юзера
			cookie, id, err := uI.service.AddUser()
			if err != nil {
				return
			}
			//Устанавливаем куку
			c.SetCookie("user_id", cookie, 3600, "/", "localhost", false, true)
			// переброска данных далее в запрос
			c.Set("user_ID", id)
			c.Next()
			return
		}
		// Нахождение пользователя и проверка куки
		id, ok := uI.service.FindUser(cookie)
		if !ok {
			// Добавим юзера
			cookie, id, err = uI.service.AddUser()
			if err != nil {
				return
			}
			//Устанавливаем куку
			c.SetCookie("user_id", cookie, 3600, "/", "localhost", false, true)
			c.Set("user_ID", id)
			c.Next()
			return
		}

		//log.Printf("user_ID: %s \n", cookie)
		log.Printf("user_ID: %v \n", id)
		// Передача запроса в handler
		c.Next()
	}
}

func (uI *UserInteract) GetAllUserURL(c *gin.Context) {
	/// Получение userIdB
	cookie, err := c.Cookie("user_id")
	data := make([]byte, 8, 39)
	// Ошибка означает, что куки не было, и нужно взять ID, который установили в запросе
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
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	userID := binary.LittleEndian.Uint64(data[:8])
	userURLs, err := uI.service.ReadAllUserURL(userID)
	if err != nil {
		log.Println("Ошибка во время получения всех URL юзера")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if len(userURLs) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.JSON(http.StatusOK, userURLs)
}
