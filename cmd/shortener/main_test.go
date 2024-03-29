package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"yaGoShortURL/internal/components/cache"
	"yaGoShortURL/internal/components/filestorage"
	"yaGoShortURL/internal/components/postgres"
	"yaGoShortURL/internal/controller/handlers"
	"yaGoShortURL/internal/usecase"
)

func TestPingRoute(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code     int
		response string
	}
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name   string
		method string
		body   string
		id     string
		want   want
	}{
		// определяем все тесты
		{
			name:   "positive test #1",
			method: "POST",
			body:   "https://habr.com/ru/company/ruvds/blog/562878/",
			id:     "",
			want: want{
				code:     http.StatusCreated,
				response: "http://localhost:8080/0",
			},
		},
		{
			name:   "positive test #2",
			method: "POST",
			body:   "https://pkg.go.dev/net/http",
			want: want{
				code:     http.StatusCreated,
				response: "http://localhost:8080/1",
			},
		},
		{
			name:   "negative test #3",
			method: "POST",
			body:   "https://pkg.go.dev/net/http",
			want: want{
				code:     http.StatusConflict,
				response: "http://localhost:8080/1",
			},
		},
		{
			name:   "negative test #4",
			method: "POST",
			body:   "",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name:   "positive test #5",
			method: "GET",
			body:   "",
			id:     "0",
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "https://habr.com/ru/company/ruvds/blog/562878/",
			},
		},
		{
			name:   "negative test #6",
			method: "GET",
			body:   "",
			id:     "5",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
	}
	// Конфигурирование сервиса
	cfg := configInit()
	cfg.UnitTestFlag = false
	// Создание экземпляров компоненинтов сервиса
	serverCash := cache.NewCash(cfg.BaseURL)
	pg, err := postgres.New(cfg.PostgresURL, cfg.UsingDB)
	// Ошибка БД
	if err != nil {
		log.Println("app - Run - postgres.New: %w", err)
	}
	defer pg.Close()
	serverFileStorage := filestorage.NewFileStorage(cfg.UnitTestFlag, cfg.FileStoragePath)
	services := usecase.NewService(serverCash, serverFileStorage, pg, cfg.UsingDB)
	myHandlers := handlers.NewHandler(services, cfg.BaseURL)

	router := myHandlers.InitRoutes()

	_, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			url := "http://localhost:8080/" + tt.id
			//Передача параметров в новый запрос
			req, _ := http.NewRequest(tt.method, url, strings.NewReader(tt.body))
			router.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			switch tt.method {
			case "POST":
				// получаем и проверяем тело запроса при запросе POST
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.want.response, string(resBody))
			case "GET":
				// произошёл ли редирект
				assert.Equal(t, tt.want.code, res.StatusCode)
				// получаем и проверяем новый адрес при запросе GET
				assert.Equal(t, tt.want.response, res.Header.Get("Location"))
			}
		})
	}
}
