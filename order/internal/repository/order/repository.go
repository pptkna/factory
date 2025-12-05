package order

import (
	"database/sql"
	"fmt"
	"time"

	def "github.com/pptkna/rocket-factory/order/internal/repository"

	_ "github.com/lib/pq"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository(host, port, user, password, dbname, sslmode string) (*repository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	db.SetMaxIdleConns(5)                  // Максимальное количество неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to database")

	return &repository{db: db}, nil
}
