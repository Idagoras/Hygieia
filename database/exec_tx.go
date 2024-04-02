package database

import (
	"context"
	"fmt"
)

func (store *MySqlStore) execTx(ctx context.Context, fn func(queries *Queries) (any, error)) (any, error) {
	tx, err := store.db.Begin()
	if err != nil {
		return nil, err
	}

	q := New(tx)
	res, err := fn(q)
	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return nil, fmt.Errorf("tx err: %v,rb err: %v", err, rbError)
		}
		return nil, err
	}
	return res, tx.Commit()
}
