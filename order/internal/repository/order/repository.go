package order

import (
	"database/sql"
	"fmt"
	"time"

	def "github.com/pptkna/rocket-factory/order/internal/repository"
	"github.com/pptkna/rocket-factory/platform/pkg/migrator"

	_ "github.com/lib/pq"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository(dsn, migrationsdir string) (*repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Setting up a connection pool
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of inactive connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

	// Checking the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	migratorRunner := migrator.NewMigrator(db, migrationsdir)

	err = migratorRunner.Up()
	if err != nil {
		return nil, fmt.Errorf("db migration error: %w", err)
	}

	return &repository{db: db}, nil
}

func (r *repository) Close() error {
	return r.db.Close()
}
