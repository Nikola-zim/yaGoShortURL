package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 2
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

// New - .
func New(url string, usingDB bool, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}
	if usingDB {
		// Установка кастомных опций
		for _, opt := range opts {
			opt(pg)
		}

		pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		poolConfig, err := pgxpool.ParseConfig(url)
		if err != nil {
			return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
		}

		poolConfig.MaxConns = int32(pg.maxPoolSize)

		for pg.connAttempts > 0 {
			pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
			if err == nil {
				break
			}

			log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

			time.Sleep(pg.connTimeout)

			pg.connAttempts--
		}

		if err != nil {
			return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
		}
	}
	return pg, nil
}

// PingDB Проверка связи с базой данных
func (pg *Postgres) PingDB() error {
	connAttempts := 5
	for connAttempts > 0 {
		err := pg.Pool.Ping(context.Background())
		if err == nil {
			return err
		}
		connAttempts--
	}
	return errors.New("ошибка подключения к БД")
}

// Close -.
func (pg *Postgres) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
