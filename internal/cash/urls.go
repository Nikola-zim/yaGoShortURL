package cash

import (
	"encoding/binary"
	"errors"
	"log"
	"regexp"
	"strconv"
	"sync"
)

type Urls struct {
	urlsMap map[string]string
	mux     sync.RWMutex
	// Для поиска по индексу
	usersUrls map[uint64][]string
}

func (u *Urls) WriteURLInCash(fullURL string, userIDB []byte) (string, error) {
	if len(userIDB) != 8 {
		return "", errors.New("нет userIDB")
	}
	// Получение id в виде числа
	userID := binary.LittleEndian.Uint64(userIDB)

	u.mux.Lock()
	defer u.mux.Unlock()
	//
	numbOfElements := len(u.urlsMap)
	log.Println(fullURL)
	//Проверка того, что передаваемая строка является URL
	re := regexp.MustCompile(`^http(s)?:\/\/[^\s]+$`)
	if re.MatchString(fullURL) {
		// Для проверки наличия Url-ов
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

		// привязка url к пользователю
		// проверка, что у этого пользователя уже есть URLs
		_, ok := u.usersUrls[userID]
		if !ok {
			// Если URL-ов нет, создаем слайс для их id-ков
			userURLsID := make([]string, 0, 10)
			userURLsID = append(userURLsID, idKey)
			u.usersUrls[userID] = userURLsID
		} else {
			userURLs := u.usersUrls[userID]
			userURLs = append(userURLs, idKey)
			u.usersUrls[userID] = userURLs
		}
		return strconv.Itoa(numbOfElements / 2), nil
	} else {

		return "", errors.New("передаваемая строка не является URL")
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

func (u *Urls) ReadAllUserURLFromCash(userIDB []byte) ([]string, error) {
	// Получение id в виде числа
	userID := binary.LittleEndian.Uint64(userIDB)
	// Слайс для результата, в нем все URL от User-а
	userURLs := make([]string, 0, 10)

	u.mux.Lock()
	defer u.mux.Unlock()

	// проверка, что у этого пользователя уже есть URLs
	_, ok := u.usersUrls[userID]
	if !ok {
		// Если URL-ов нет, создаем слайс для их id-ков
		userURLs = make([]string, 0, 10)
		return userURLs, nil
	} else {
		userURLsID := u.usersUrls[userID]
		for _, id := range userURLsID {
			fullURL, err := u.urlsMap[id]
			if !err {
				return userURLs, errors.New("ошибка чтения из кеша всех URL пользователя: такого ID не существует")
			}
			if fullURL == "" {
				return userURLs, errors.New("ошибка чтения из кеша всех URL пользователя: пустой URL")
			}
			userURLs = append(userURLs, fullURL)
		}
		return userURLs, nil
	}
}

func NewUrls() *Urls {
	return &Urls{
		urlsMap:   make(map[string]string),
		usersUrls: make(map[uint64][]string),
	}
}
