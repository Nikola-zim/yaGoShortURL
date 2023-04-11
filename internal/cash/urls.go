package cash

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
)

type Urls struct {
	urlsMap map[string]string
	mux     sync.RWMutex
}

func (u *Urls) WriteURLInCash(fullURL string) (string, error) {
	u.mux.Lock()
	defer u.mux.Unlock()
	numbOfElements := len(u.urlsMap)
	fmt.Println("smth in cash")
	log.Println(fullURL)
	//Проверка того, что передаваемая строка является URL
	re := regexp.MustCompile(`^http(s)?:\/\/[^\s]+$`)
	if re.MatchString(fullURL) {
		//Проверка наличия элемента
		strKeyCheck := "url:" + fullURL
		_, err := u.urlsMap[strKeyCheck]
		if err {
			return "", errors.New("URL is already in memory")
		}
		//Всегда должнобыть четное число элементов в структуре map
		if numbOfElements%2 != 0 {
			return "", errors.New("ошибка кеша")
		}
		//Запись в map после проверок
		//Форматирование ключей
		idKey := "id:" + strconv.Itoa(numbOfElements/2)
		strKey := "url:" + fullURL
		u.urlsMap[idKey] = fullURL
		u.urlsMap[strKey] = fullURL
		return strconv.Itoa(numbOfElements / 2), nil
	} else {

		return "", errors.New("передаваемая строка не является URL!!!")
	}
}

func (u *Urls) ReadURLFromCash(id string) (string, error) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	fullURL, err := u.urlsMap[id]
	if !err {
		return fullURL, errors.New("ошибка чтения из кеша: такого ID не существует")
	}
	if fullURL == "" {
		return fullURL, errors.New("ошибка чтения из кеша: пустой URL")
	}
	return fullURL, nil
}

func NewUrls() *Urls {
	return &Urls{
		urlsMap: make(map[string]string),
	}
}
