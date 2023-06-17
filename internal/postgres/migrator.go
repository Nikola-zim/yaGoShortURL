package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
	"sync"
)

var once sync.Once

// MigratorOptions опции миграции.
type MigratorOptions struct {
	Table            string // Table to store migration info.
	Scheme           string // Scheme is name of a schema that the migration table be referenced.
	ConnectionString string
	RootFS           string
	Driver           string
}

// Migrator is an object responsible for migrations.
type Migrator struct {
	log *zerolog.Logger
	fs  embed.FS
	opt MigratorOptions
}

// NewMigrator creates new migrator.
func NewMigrator(fs embed.FS, log *zerolog.Logger, opt MigratorOptions) *Migrator {
	return &Migrator{log: log, fs: fs, opt: opt}
}

// Run executes embedded SQL scripts.
// For the time being only "up" migrations are supported.
func (m *Migrator) Run(ctx context.Context) error {
	var err error

	once.Do(func() {
		err = m.runImpl(ctx)
	})

	return err
}

func (m *Migrator) runImpl(ctx context.Context) error {
	m.log.Info().Msg("Started migration")

	conn, err := sql.Open("pgx", m.opt.ConnectionString)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			m.log.Err(err).Msg("Failed to close connection.")
		}
	}()

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("can't begin db transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			m.log.Err(err).Msg("Failed to rollback migration transaction.")
		}
	}()

	migrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: m.fs,
		Root:       m.opt.RootFS,
	}

	mg, err := migrationSource.FindMigrations()
	if err != nil {
		return fmt.Errorf("could not find migrations. %w", err)
	}

	m.log.Debug().Int("Count", len(mg)).Msg("Found migrations.")

	if _, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(2)`); err != nil {
		m.log.Err(err).Msg("Could not acquire advisory lock.")

		return fmt.Errorf("can't acquire advisory lock: %w", err)
	}

	migrate.SetTable(m.opt.Table)
	migrate.SetSchema(m.opt.Scheme)

	applied, err := migrate.Exec(conn, m.opt.Driver, migrationSource, migrate.Up)
	if err != nil {
		m.log.Err(err).Msg("Could not apply database migrations.")

		return fmt.Errorf("can't apply database migrations: %w", err)
	}

	m.log.Debug().Int("Count", applied).Msg("Applied migrations")

	if err := tx.Commit(); err != nil {
		m.log.Err(err).Msg("Could not commit db transaction.")

		return fmt.Errorf("can't commit db transaction: %w", err)
	}

	m.log.Info().Msg("Migration finished.")

	return nil
}
