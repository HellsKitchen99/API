package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(dns string) (*sql.DB, error) {
	connStr := dns
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД")
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("ошибка пинга БД")
	}
	return db, nil
}
