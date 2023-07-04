package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Cache interface {
	FindUser(idMsg string) (uint64, bool)
	AddUser() (string, uint64, error)
}

type UserInteract struct {
	service Cache
}

func NewUserInteract(service Cache) *UserInteract {
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
