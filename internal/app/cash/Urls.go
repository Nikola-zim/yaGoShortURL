package cash

import (
	"errors"
	"strconv"
	"sync"
)

type Urls struct {
	urlsMap map[string]string
	mux     sync.RWMutex
	wg      sync.WaitGroup
}

func (u *Urls) WriteURLInCash(string2 string) (string, error) {
	defer u.wg.Done()
	u.wg.Add(1)
	u.mux.Lock()
	numbOfElements := len(u.urlsMap)
	//Всегда должнобыть четное число элементов в структуре map
	if numbOfElements%2 == 0 {
		//Проверка наличия элемента
		strKeyCheck := "url:" + string2
		_, err := u.urlsMap[strKeyCheck]
		if err {
			u.mux.Unlock()
			return "0", errors.New("Ошибка записи в кеш: значение уже существует")
		}
		//Запись в map после проверок
		//Форматирование ключей
		idKey := "id:" + strconv.Itoa(numbOfElements/2)
		strKey := "url:" + string2
		u.urlsMap[idKey] = string2
		u.urlsMap[strKey] = string2
		u.mux.Unlock()
	}
	return strconv.Itoa(numbOfElements / 2), nil
}

func (u *Urls) ReadURLFromCash(id string) (string, error) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	fullURL, err := u.urlsMap[id]
	if !err {
		return fullURL, errors.New("ошибка чтения из кеша: такого ID не существует")
	}
	return fullURL, nil
}

func NewUrls() *Urls {
	return &Urls{
		urlsMap: make(map[string]string),
	}
}
