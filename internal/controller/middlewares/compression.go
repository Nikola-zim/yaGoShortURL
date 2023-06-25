package middlewares

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// GzipHandle для сжатия запросов и ответов.
func GzipHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content-Encoding
		// проверяем, что запрос отправлен в формате gzip
		if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			// Create a gzip writer
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer func(reader *gzip.Reader) {
				err := reader.Close()
				if err != nil {
					log.Println(err)
				}
			}(reader)
			c.Request.Body = reader
			c.Request.Header.Del("Content-Encoding")
			// Передача запроса в handler
			c.Next()
			return
		}
		// Accept-Encoding
		// проверяем, что клиент поддерживает gzip-сжатие
		if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			// создаём gzip.Writer поверх текущего w
			gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer func(gz *gzip.Writer) {
				err := gz.Close()
				if err != nil {
					log.Println(err)
				}
			}(gz)
			//Set the Content-Encoding header
			c.Header("Content-Encoding", "gzip")

			//Replace the writer with gzip writer
			c.Writer = &gzipWriter{Writer: gz, ResponseWriter: c.Writer}
			// Передача запроса в handler
			c.Next()
			return
		}
		// Передача запроса в handler
		c.Next()
	}
}

// Опишем тип gzipWriter, поддерживающий интерфейс http.ResponseWriter, и реализуем недостающие методы.
type gzipWriter struct {
	gin.ResponseWriter
	*gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func (g *gzipWriter) Flush() {
	err := g.Writer.Flush()
	if err != nil {
		log.Println(err)
		return
	}
	g.ResponseWriter.Flush()
}

func (g *gzipWriter) Header() http.Header {
	return g.ResponseWriter.Header()
}
