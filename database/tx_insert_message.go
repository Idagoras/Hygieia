package database

import (
	"context"
	"fmt"
)

type InsertMessageTxParams struct {
	Message *Message
}

func (store *MySqlStore) InsertMessageTx(ctx context.Context, arg InsertMessageTxParams) error {
	_, err := store.execTx(ctx, func(queries *Queries) (any, error) {
		var err error
		message, err := store.GetMessage(ctx, GetMessageParams{
			SendUid:   arg.Message.SendUid,
			RcvUid:    arg.Message.RcvUid,
			CreatedAt: arg.Message.CreatedAt,
		})
		if len(message) == 0 {
			err = store.InsertMessage(ctx, InsertMessageParams{
				SendUid:    arg.Message.SendUid,
				RcvUid:     arg.Message.RcvUid,
				CreatedAt:  arg.Message.CreatedAt,
				HasRead:    arg.Message.HasRead,
				Type:       arg.Message.Type,
				Text:       arg.Message.Text,
				Subtype:    arg.Message.Subtype,
				Attachment: arg.Message.Attachment,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to insert message: %v ", err)
			}
		} else {
			return nil, fmt.Errorf("receive duplicate message ")
		}
		return nil, err
	})
	return err
}
