package database

import (
	"context"
	"fmt"
)

type UpdateUserTxParams struct {
	UpdateUserParams
	AfterCreate func(user *User) error
}

func (store *MySqlStore) UpdateUserTx(ctx context.Context, arg UpdateUserTxParams) (*User, error) {
	res, err := store.execTx(ctx, func(queries *Queries) (any, error) {
		var err error
		err = queries.UpdateUser(ctx, arg.UpdateUserParams)
		if err != nil {
			return nil, err
		}
		user, err := queries.GetUserByUid(ctx, arg.ID)
		if err != nil {
			return nil, err
		}
		return user, arg.AfterCreate(&user)
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
