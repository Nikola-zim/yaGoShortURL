package postgres

import (
	"context"
	"log"
	"yaGoShortURL/internal/entity"
)

func (pg *Postgres) GetAllURLFromDB(ctx context.Context) ([]entity.DataURL, error) {
	res := make([]entity.DataURL, 0, 100)
	// Start transaction
	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, getAllURLs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataURL entity.DataURL
		err = rows.Scan(&dataURL.URL, &dataURL.UserID)
		if err != nil {
			return nil, err
		}
		res = append(res, dataURL)
	}
	return res, nil
}

func (pg *Postgres) WriteURLInDB(fullURL string, id string, userID uint64) error {
	ctx := context.Background()

	// Start transaction
	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, writeURL, id, fullURL, userID)

	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}

		log.Printf("Failed to insert url in DB")

		return err
	}

	// Подтверждение транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
