// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: subscriber.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const getSubscribePair = `-- name: GetSubscribePair :one
SELECT id, uid, sbs_id, subscribe_time, is_cancel, cancel_time FROM ` + "`" + `subscribers` + "`" + ` WHERE uid = ? and sbs_id = ?
`

type GetSubscribePairParams struct {
	Uid   uint64 `json:"uid"`
	SbsID uint64 `json:"sbs_id"`
}

func (q *Queries) GetSubscribePair(ctx context.Context, arg GetSubscribePairParams) (Subscriber, error) {
	row := q.db.QueryRowContext(ctx, getSubscribePair, arg.Uid, arg.SbsID)
	var i Subscriber
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.SbsID,
		&i.SubscribeTime,
		&i.IsCancel,
		&i.CancelTime,
	)
	return i, err
}

const insertSubscribePair = `-- name: InsertSubscribePair :exec
INSERT INTO ` + "`" + `subscribers` + "`" + ` (uid, sbs_id, subscribe_time, cancel_time) VALUES
(
 ?,?,?,?
)
`

type InsertSubscribePairParams struct {
	Uid           uint64    `json:"uid"`
	SbsID         uint64    `json:"sbs_id"`
	SubscribeTime time.Time `json:"subscribe_time"`
	CancelTime    time.Time `json:"cancel_time"`
}

func (q *Queries) InsertSubscribePair(ctx context.Context, arg InsertSubscribePairParams) error {
	_, err := q.db.ExecContext(ctx, insertSubscribePair,
		arg.Uid,
		arg.SbsID,
		arg.SubscribeTime,
		arg.CancelTime,
	)
	return err
}

const listSubscribePairBySuid = `-- name: ListSubscribePairBySuid :many
SELECT id, uid, sbs_id, subscribe_time, is_cancel, cancel_time FROM ` + "`" + `subscribers` + "`" + ` WHERE sbs_id = ? AND is_cancel = false LIMIT ? OFFSET ?
`

type ListSubscribePairBySuidParams struct {
	SbsID  uint64 `json:"sbs_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListSubscribePairBySuid(ctx context.Context, arg ListSubscribePairBySuidParams) ([]Subscriber, error) {
	rows, err := q.db.QueryContext(ctx, listSubscribePairBySuid, arg.SbsID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subscriber{}
	for rows.Next() {
		var i Subscriber
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.SbsID,
			&i.SubscribeTime,
			&i.IsCancel,
			&i.CancelTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSubscribePairByUid = `-- name: ListSubscribePairByUid :many
SELECT id, uid, sbs_id, subscribe_time, is_cancel, cancel_time FROM ` + "`" + `subscribers` + "`" + ` WHERE uid = ? AND is_cancel = false LIMIT ? OFFSET ?
`

type ListSubscribePairByUidParams struct {
	Uid    uint64 `json:"uid"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListSubscribePairByUid(ctx context.Context, arg ListSubscribePairByUidParams) ([]Subscriber, error) {
	rows, err := q.db.QueryContext(ctx, listSubscribePairByUid, arg.Uid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Subscriber{}
	for rows.Next() {
		var i Subscriber
		if err := rows.Scan(
			&i.ID,
			&i.Uid,
			&i.SbsID,
			&i.SubscribeTime,
			&i.IsCancel,
			&i.CancelTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSubscribePair = `-- name: UpdateSubscribePair :exec
UPDATE ` + "`" + `subscribers` + "`" + ` SET
subscribe_time = COALESCE(?,subscribe_time),
cancel_time = COALESCE(?,cancel_time),
is_cancel = COALESCE(?,is_cancel)
WHERE id = ?
`

type UpdateSubscribePairParams struct {
	SubscribeTime sql.NullTime `json:"subscribe_time"`
	CancelTime    sql.NullTime `json:"cancel_time"`
	IsCancel      sql.NullBool `json:"is_cancel"`
	ID            uint64       `json:"id"`
}

func (q *Queries) UpdateSubscribePair(ctx context.Context, arg UpdateSubscribePairParams) error {
	_, err := q.db.ExecContext(ctx, updateSubscribePair,
		arg.SubscribeTime,
		arg.CancelTime,
		arg.IsCancel,
		arg.ID,
	)
	return err
}
