package cash

import (
	"errors"
	"log"
	"sync"
)

type Urls struct {
	urlsMap map[int]string
	mux     sync.RWMutex
	wg      sync.WaitGroup
}

func (u *Urls) WriteUrlInCash(string2 string) error {
	defer u.wg.Done()
	u.wg.Add(1)
	u.mux.Lock()
	//Проверка наличия элемента
	_, err := u.urlsMap[len(u.urlsMap)]
	if err {
		log.Println("Ошибка записи в кеш: значение уже существует")
		u.mux.Unlock()
		//TODO вернуть ошибку
		return errors.New("dummy")
	}
	u.urlsMap[len(u.urlsMap)] = string2
	u.mux.Unlock()
	return nil
}

func (u *Urls) ReadUrlFromCash(int2 int) (string, error) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	fullUrl, err := u.urlsMap[int2]
	if err != true {
		log.Println("Ошибка чтения из кеша: такого ID не существует")
		//Todo нормальную ошибку
		return fullUrl, errors.New("dummy")
	}
	return fullUrl, nil
}

func NewUrls() *Urls {
	return &Urls{
		urlsMap: make(map[int]string),
	}
}
