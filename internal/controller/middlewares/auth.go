package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
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

// CookieSetAndGet - middleware для cookie.
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
