package cache

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"yaGoShortURL/internal/entity"
)

type Urls struct {
	urlsMap map[string]string
	mux     sync.RWMutex
	// Для поиска по индексу
	usersUrls map[uint64][]string
	// Мапа с сокращенными URL
	URLsAllInfo map[string]entity.JSONAllInfo
	baseURL     string
}

func (u *Urls) WriteURL(fullURL string, userIDB []byte) (string, error) {
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
	_, err := url.ParseRequestURI(fullURL)
	if err != nil {
		return "", errors.New("передаваемая строка не является URL")
	}

	// Для проверки наличия Url-ов
	strKeyCheck := "fullURL:" + fullURL
	oldURL, found := u.urlsMap[strKeyCheck]

	if found {
		return "", entity.NewErrorURL(errors.New("URL is already in memory"), oldURL)
	}

	//Всегда должнобыть четное число элементов в структуре map
	if numbOfElements%2 != 0 {
		return "", errors.New("cache error: number of elements is invalid")
	}
	//Текущий индекс (для сокращенного URL)
	currentID := numbOfElements / 2

	//Запись в map после проверок
	//Форматирование ключей
	idKey := "id:" + strconv.Itoa(currentID)
	strKey := "fullURL:" + fullURL
	u.urlsMap[idKey] = fullURL
	u.urlsMap[strKey] = strconv.Itoa(currentID)
	//Составим полный адрес сокращенного URL
	baseURL := fmt.Sprintf("%s%s%v", u.baseURL, "/", currentID)
	u.URLsAllInfo[idKey] = entity.JSONAllInfo{
		FullURL: fullURL,
		BaseURL: baseURL,
	}

	// привязка URL к пользователю
	// проверка, что у этого пользователя уже есть URLs
	_, ok := u.usersUrls[userID]

	if !ok {
		// Если URL-ов нет, создаем слайс для их id-ков
		userURLs := make([]string, 0, 10)
		userURLs = append(userURLs, idKey)
		u.usersUrls[userID] = userURLs
	} else {
		userURLs := u.usersUrls[userID]
		userURLs = append(userURLs, idKey)
		u.usersUrls[userID] = userURLs
	}

	return strconv.Itoa(currentID), nil
}

func (u *Urls) FullURL(id string) (string, error) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	fullURL, found := u.urlsMap[id]

	if !found {
		return fullURL, errors.New("ошибка чтения из кеша: такого ID не существует")
	}

	if fullURL == "" {
		return fullURL, errors.New("ошибка чтения из кеша: пустой URL")
	}

	return fullURL, nil
}

func (u *Urls) ReadAllUserURLFromCash(userIDB []byte) ([]entity.JSONAllInfo, error) {
	// Получение id в виде числа
	userID := binary.LittleEndian.Uint64(userIDB)
	// Слайс для результата, в нем все URL от User-а
	userURLs := make([]entity.JSONAllInfo, 0, 10)

	u.mux.Lock()
	defer u.mux.Unlock()

	// проверка, что у этого пользователя уже есть URLs
	_, ok := u.usersUrls[userID]
	if !ok {
		// Если URL-ов нет, создаем слайс для их id-ков
		//userURLs = make([]string, 0, 10)
		return userURLs, nil
	} else {
		userURLsID := u.usersUrls[userID]
		for _, id := range userURLsID {
			currentURLsAllInfo, err := u.URLsAllInfo[id]
			if !err {
				return userURLs, errors.New("ошибка чтения из кеша всех URL пользователя: такого ID не существует")
			}
			if currentURLsAllInfo.FullURL == "" {
				return userURLs, errors.New("ошибка чтения из кеша всех URL пользователя: пустой URL")
			}
			userURLs = append(userURLs, currentURLsAllInfo)
		}
		return userURLs, nil
	}
}

func NewUrls(baseURL string) *Urls {
	return &Urls{
		baseURL:     baseURL,
		urlsMap:     make(map[string]string),
		usersUrls:   make(map[uint64][]string),
		URLsAllInfo: make(map[string]entity.JSONAllInfo),
	}
}
