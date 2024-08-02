package db

import "context"

type IDBconnection interface {
	FindUserByEmailAndPassword(ctx context.Context, email, password string) (string, error)
	CreateUser(ctx context.Context, name, email, password string) error
}
