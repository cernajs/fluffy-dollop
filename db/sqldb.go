package db

import (
	"context"
	"database/sql"
)

type MySQLDbConnection struct {
	db *sql.DB
}

func NewMySQLDBconnection(db *sql.DB) *MySQLDbConnection {
	return &MySQLDbConnection{db: db}
}

func (store *MySQLDbConnection) FindUserByEmailAndPassword(ctx context.Context, email, password string) (string, error) {
	var name string
	err := store.db.QueryRowContext(ctx, "SELECT name FROM users WHERE email = ? AND password = ?", email, password).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (store *MySQLDbConnection) CreateUser(ctx context.Context, name, email, password string) error {
	_, err := store.db.ExecContext(ctx, "INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, password)
	return err
}
