package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"yaGoShortURL/internal/app/cash"
	"yaGoShortURL/internal/app/handlers"
	"yaGoShortURL/internal/app/service"
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
				code:     http.StatusBadRequest,
				response: "",
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
	}
	serverCash := cash.NewCash()
	services := service.NewService(serverCash)
	handlers := handlers.NewHandler(services)
	router := handlers.InitRoutes()
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
			// получаем и проверяем тело запроса
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}
