package database

import (
	"context"
	"fmt"
)

func (store *MySqlStore) CreateEEGSessionTx(ctx context.Context, arg CreateEEGSessionParams) (EegSession, error) {
	res, err := store.execTx(ctx, func(queries *Queries) (any, error) {
		var err error
		err = queries.CreateEEGSession(ctx, arg)
		if err != nil {
			return nil, err
		}
		session, err := queries.GetLatestEEGSession(ctx)
		if err != nil {
			return nil, err
		}
		return session, err
	})

	if err != nil {
		return EegSession{}, err
	}
	session, ok := res.(EegSession)
	if ok {
		return session, nil
	}
	return EegSession{}, fmt.Errorf("cannot convert to eegsession")
}
