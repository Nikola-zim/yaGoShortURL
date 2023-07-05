package workers

import (
	"context"
	"log"
	"yaGoShortURL/internal/entity"
)

type Eraser struct {
	ctx         context.Context
	cash        CacheURL
	pg          DataBase
	usingDB     bool
	toDeleteMsg chan entity.DeleteMsg
}

func NewEraser(ctx context.Context, cash CacheURL, pg DataBase, usingDB bool, msg chan entity.DeleteMsg) *Eraser {
	return &Eraser{
		ctx:         ctx,
		cash:        cash,
		pg:          pg,
		usingDB:     usingDB,
		toDeleteMsg: msg,
	}
}

func (er *Eraser) Run() {
	var msg entity.DeleteMsg
	for {
		msg = <-er.toDeleteMsg
		err := er.cash.DeleteURLs(msg.UserID, msg.List)
		if err != nil {
			log.Printf("[workers - eraser - Run - cash.DeleteURLs]: %s", err)
		}
		if er.usingDB {
			err = er.pg.DeleteURLsDB(er.ctx, msg.UserID, msg.List)
		}
	}
}

func (er *Eraser) ShutDown() {
	var msg entity.DeleteMsg
	for msg = range er.toDeleteMsg {
		err := er.cash.DeleteURLs(msg.UserID, msg.List)
		if err != nil {
			log.Printf("[workers - eraser - Run - cash.DeleteURLs]: %s", err)
		}
		if er.usingDB {
			err = er.pg.DeleteURLsDB(er.ctx, msg.UserID, msg.List)
		}
	}
}
