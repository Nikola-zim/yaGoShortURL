package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Postgres struct {
	service DBService
}

func NewPostgres(service Cash) *Postgres {
	return &Postgres{
		service: service,
	}
}
func (pg *Postgres) pingDB(c *gin.Context) {
	if err := pg.service.PingDB(); err == nil {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}
}
