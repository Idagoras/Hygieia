package database

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Store interface {
	Querier
	CreateEEGSessionTx(ctx context.Context, arg CreateEEGSessionParams) (EegSession, error)
	InsertEEGDataTx(ctx context.Context, arg InsertEEGDataTxParams) error
	SubscribeTx(ctx context.Context, arg SubscribeTxParams) error
	InsertMessageTx(ctx context.Context, arg InsertMessageTxParams) error
	CreateUserTx(ctx context.Context, arg InsertUserParams) (*User, error)
	UpdateUserTx(ctx context.Context, arg UpdateUserTxParams) (*User, error)
}

type MySqlStore struct {
	db *sql.DB
	*Queries
}

func NewMySqlStore(db *sql.DB) Store {
	return &MySqlStore{
		db,
		New(db),
	}
}
