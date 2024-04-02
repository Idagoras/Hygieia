// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: message.sql

package database

import (
	"context"
	"time"
)

const getMessage = `-- name: GetMessage :many
SELECT id, send_uid, rcv_uid, created_at, has_read, type, text, subtype, attachment FROM ` + "`" + `message` + "`" + ` WHERE send_uid = ? AND rcv_uid = ? AND created_at = ? LIMIT 1
`

type GetMessageParams struct {
	SendUid   uint64    `json:"send_uid"`
	RcvUid    uint64    `json:"rcv_uid"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) GetMessage(ctx context.Context, arg GetMessageParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessage, arg.SendUid, arg.RcvUid, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.SendUid,
			&i.RcvUid,
			&i.CreatedAt,
			&i.HasRead,
			&i.Type,
			&i.Text,
			&i.Subtype,
			&i.Attachment,
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

const insertMessage = `-- name: InsertMessage :exec
INSERT INTO ` + "`" + `message` + "`" + ` (SEND_UID, RCV_UID, CREATED_AT, HAS_READ, TYPE, TEXT, SUBTYPE, ATTACHMENT) VALUES
(
 ?,?,?,?,?,?,?,?
)
`

type InsertMessageParams struct {
	SendUid    uint64    `json:"send_uid"`
	RcvUid     uint64    `json:"rcv_uid"`
	CreatedAt  time.Time `json:"created_at"`
	HasRead    bool      `json:"has_read"`
	Type       uint32    `json:"type"`
	Text       string    `json:"text"`
	Subtype    uint32    `json:"subtype"`
	Attachment string    `json:"attachment"`
}

func (q *Queries) InsertMessage(ctx context.Context, arg InsertMessageParams) error {
	_, err := q.db.ExecContext(ctx, insertMessage,
		arg.SendUid,
		arg.RcvUid,
		arg.CreatedAt,
		arg.HasRead,
		arg.Type,
		arg.Text,
		arg.Subtype,
		arg.Attachment,
	)
	return err
}

const listUserReceivedMessages = `-- name: ListUserReceivedMessages :many
SELECT id, send_uid, rcv_uid, created_at, has_read, type, text, subtype, attachment FROM ` + "`" + `message` + "`" + ` WHERE rcv_uid = ? LIMIT ? OFFSET ?
`

type ListUserReceivedMessagesParams struct {
	RcvUid uint64 `json:"rcv_uid"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListUserReceivedMessages(ctx context.Context, arg ListUserReceivedMessagesParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, listUserReceivedMessages, arg.RcvUid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.SendUid,
			&i.RcvUid,
			&i.CreatedAt,
			&i.HasRead,
			&i.Type,
			&i.Text,
			&i.Subtype,
			&i.Attachment,
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

const listUserSentMessages = `-- name: ListUserSentMessages :many
SELECT id, send_uid, rcv_uid, created_at, has_read, type, text, subtype, attachment FROM ` + "`" + `message` + "`" + ` WHERE send_uid = ? LIMIT ? OFFSET ?
`

type ListUserSentMessagesParams struct {
	SendUid uint64 `json:"send_uid"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

func (q *Queries) ListUserSentMessages(ctx context.Context, arg ListUserSentMessagesParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, listUserSentMessages, arg.SendUid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.SendUid,
			&i.RcvUid,
			&i.CreatedAt,
			&i.HasRead,
			&i.Type,
			&i.Text,
			&i.Subtype,
			&i.Attachment,
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
