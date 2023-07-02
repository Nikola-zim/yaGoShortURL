package cache

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"yaGoShortURL/internal/entity"
)

const (
	defaultURLsNumber = 15
)

type Urls struct {
	urlsMap map[string]string
	mux     sync.RWMutex
	// Для поиска по индексу
	usersUrls map[uint64][]string
	// Мапа с сокращенными URL
	qURLsAllInfo map[string]entity.JSONAllInfo
	URLs         URLsAllInfo
	baseURL      string
}

func NewUrls(baseURL string) *Urls {
	return &Urls{
		baseURL:   baseURL,
		urlsMap:   make(map[string]string),
		usersUrls: make(map[uint64][]string),
		URLs: struct {
			IDKey   map[string]entity.JSONAllInfo
			URLKey  map[string]string
			Counter int
		}{
			IDKey:   make(map[string]entity.JSONAllInfo, defaultURLsNumber),
			URLKey:  make(map[string]string, defaultURLsNumber),
			Counter: 0,
		},
	}
}

func (u *Urls) WriteURL(fullURL string, userID uint64) (string, error) {
	u.mux.Lock()
	defer u.mux.Unlock()

	//Проверка того, что передаваемая строка является URL
	_, err := url.ParseRequestURI(fullURL)
	if err != nil {
		return "", errors.New("передаваемая строка не является URL")
	}

	// Для проверки наличия Url-ов
	oldURL, found := u.URLs.URLKey[fullURL]

	if found {
		return "", entity.NewErrorURL(errors.New("URL is already in memory"), oldURL)
	}

	// Запись в map.
	// Составим полный адрес сокращенного URL.
	ShortURL := fmt.Sprintf("%s%s%v", u.baseURL, "/", u.URLs.Counter)
	URLAllInfo := entity.JSONAllInfo{
		ShortURL:  ShortURL,
		FullURL:   fullURL,
		IsDeleted: false,
	}
	idKey := strconv.Itoa(u.URLs.Counter)
	u.URLs.Counter++

	u.URLs.IDKey[idKey] = URLAllInfo
	u.URLs.URLKey[fullURL] = idKey

	// Привязка URL к пользователю.
	// Проверка, что у этого пользователя уже есть URLs.
	_, ok := u.usersUrls[userID]

	if !ok {
		// Если URL-ов нет, создаем слайс для их id-ков
		userURLs := make([]string, 0, defaultURLsNumber)
		userURLs = append(userURLs, idKey)
		u.usersUrls[userID] = userURLs
	} else {
		userURLs := u.usersUrls[userID]
		userURLs = append(userURLs, idKey)
		u.usersUrls[userID] = userURLs
	}

	return idKey, nil
}

func (u *Urls) FullURL(id string) (string, bool, error) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	URLInfo, found := u.URLs.IDKey[id]

	if !found {
		return URLInfo.FullURL, false, errors.New("ошибка чтения из кеша: такого ID не существует")
	}

	return URLInfo.FullURL, URLInfo.IsDeleted, nil
}

func (u *Urls) ReadAllUserURL(userID uint64) ([]entity.JSONAllInfo, error) {
	// Слайс для результата, в нем все URL от User-а
	userURLs := make([]entity.JSONAllInfo, 0, defaultURLsNumber)

	u.mux.Lock()
	defer u.mux.Unlock()

	// проверка, что у этого пользователя уже есть URLs
	_, ok := u.usersUrls[userID]
	if !ok {
		return userURLs, nil
	}
	userURLsID := u.usersUrls[userID]
	for _, id := range userURLsID {
		currentURLsAllInfo, err := u.URLs.IDKey[id]
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

func (u *Urls) DeleteURLs(userID uint64, IDs []string) error {
	u.mux.Lock()
	defer u.mux.Unlock()
	userURLsID := u.usersUrls[userID]

	for _, id := range IDs {
		for _, idUser := range userURLsID {
			if id == idUser {
				currentURLsAllInfo, _ := u.URLs.IDKey[idUser]
				URLAllInfo := entity.JSONAllInfo{
					ShortURL:  currentURLsAllInfo.ShortURL,
					FullURL:   currentURLsAllInfo.FullURL,
					IsDeleted: true,
				}
				u.URLs.IDKey[idUser] = URLAllInfo
				break
			}
		}
	}

	return nil
}
