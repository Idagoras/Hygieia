package database

import (
	"context"
	"fmt"
)

func (store *MySqlStore) CreateUserTx(ctx context.Context, arg InsertUserParams) (*User, error) {
	res, err := store.execTx(ctx, func(queries *Queries) (any, error) {
		var err error
		err = queries.InsertUser(ctx, arg)
		if err != nil {
			return nil, err
		}
		user, err := queries.GetLatestUser(ctx)
		if err != nil {
			return nil, err
		}
		return user, err
	})
	if err != nil {
		return nil, err
	}
	user, ok := res.(User)
	if ok {
		return &user, nil
	}
	return nil, fmt.Errorf("cannot convert to user")
}
