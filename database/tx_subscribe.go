package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type SubscribeTxParams struct {
	Suid uint64
	Uid  uint64
	Like bool
}

func (store *MySqlStore) SubscribeTx(ctx context.Context, arg SubscribeTxParams) error {
	_, err := store.execTx(ctx, func(queries *Queries) (any, error) {
		var err error
		pair, err := queries.GetSubscribePair(ctx, GetSubscribePairParams{
			Uid:   arg.Uid,
			SbsID: arg.Suid,
		})
		exist := true
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			exist = false
		}

		if exist {
			if arg.Like != !pair.IsCancel {
				if !arg.Like {
					err := queries.UpdateSubscribePair(ctx, UpdateSubscribePairParams{
						SubscribeTime: sql.NullTime{},
						CancelTime: sql.NullTime{
							Time:  time.Now(),
							Valid: true,
						},
						IsCancel: sql.NullBool{
							Bool:  true,
							Valid: true,
						},
						ID: pair.ID,
					})
					if err != nil {
						return nil, err
					}
				} else {
					err := queries.UpdateSubscribePair(ctx, UpdateSubscribePairParams{
						SubscribeTime: sql.NullTime{
							Time:  time.Now(),
							Valid: true,
						},
						CancelTime: sql.NullTime{},
						IsCancel:   sql.NullBool{},
						ID:         pair.ID,
					})
					if err != nil {
						return nil, err
					}
				}
			} else {
				return nil, nil
			}

		} else {
			err := queries.InsertSubscribePair(ctx, InsertSubscribePairParams{
				Uid:           arg.Uid,
				SbsID:         arg.Suid,
				SubscribeTime: time.Now(),
				CancelTime:    time.Time{},
			})
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	})

	return err
}
