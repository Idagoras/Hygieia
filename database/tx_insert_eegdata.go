package database

import (
	"Hygieia/pb"
	"context"
	"database/sql"
	"encoding/json"
	"log"
)

type InsertEEGDataTxParams struct {
	Records   []*pb.EEGData
	SessionId uint64
}

func (store *MySqlStore) InsertEEGDataTx(ctx context.Context, arg InsertEEGDataTxParams) error {
	_, err := store.execTx(ctx, func(queries *Queries) (any, error) {
		var err error
		session, err := store.GetEEGSessionByID(ctx, arg.SessionId)
		if err != nil {
			return nil, err
		}
		bitsMapBytes := []byte(session.DataBitsMap)
		count := 0
		var sessionID uint64
		for _, record := range arg.Records {
			jsonArr, _ := json.Marshal(record.Raw)
			sessionID = record.EegSessionId
			byteIndex := record.MessageOffset / 8
			byteInternalOffset := record.MessageOffset % 8
			if bitsMapBytes[byteIndex]&(1<<byteInternalOffset) > 0 {
				insertArg := InsertEEGDataParams{
					SessionID:     record.EegSessionId,
					Offset:        record.MessageOffset,
					UserID:        record.UserId,
					CollectedAt:   record.Time.AsTime(),
					Attention:     record.Attention,
					Meditation:    record.Meditation,
					BlinkStrength: record.BlinkStrength,
					Alpha1:        record.Alpha1,
					Alpha2:        record.Alpha2,
					Beta1:         record.Beta1,
					Beta2:         record.Beta2,
					Gamma1:        record.Gamma1,
					Gamma2:        record.Gamma2,
					Delta:         record.Delta,
					Theta:         record.Theta,
					Raw:           jsonArr,
				}
				err = store.InsertEEGData(ctx, insertArg)
				if err != nil {
					continue
				}
				bitsMapBytes[byteIndex] |= 1 << byteInternalOffset
				count++
			}
		}
		ierr := store.UpdateEEGSession(ctx, UpdateEEGSessionParams{
			End:       sql.NullTime{},
			ExpiredAt: sql.NullTime{},
			Finish:    sql.NullBool{},
			DataCount: sql.NullInt32{
				Int32: int32(count),
				Valid: true,
			},
			ID: sessionID,
			DataBitsMap: sql.NullString{
				String: string(bitsMapBytes),
				Valid:  true,
			},
		})
		if ierr != nil {
			log.Printf("failed to update session : %v \n", ierr)
			return nil, ierr
		}
		return nil, err
	})

	return err
}
