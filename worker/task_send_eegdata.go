package worker

import (
	"Hygieia/alogrithm"
	"Hygieia/database"
	"Hygieia/pb"
	"Hygieia/util"
	"context"
	"database/sql"
	"fmt"
	"github.com/hibiken/asynq"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"
	"time"
)

const TaskSendEEGData = "task:send_eeg_data"

func (distributor *HyTaskDistributor) DistributeTaskSendEEGData(ctx context.Context, eegDataBytes []byte, opts ...asynq.Option) error {
	task := asynq.NewTask(TaskSendEEGData, eegDataBytes, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue eask: %w", err)
	}
	log.Printf("task type: %v,queue: %v,max retry:%d,enqueued task", info.Type, info.Queue, info.MaxRetry)
	return nil
}

func (processor *HyTaskProcessor) ProcessTaskSendEEGData(ctx context.Context, task *asynq.Task) error {
	var eegData pb.EEGData
	if err := proto.Unmarshal(task.Payload(), &eegData); err != nil {
		return fmt.Errorf("failed to unmarshal payload : %w", asynq.SkipRetry)
	}
	processor.hub.Data <- &alogrithm.AlgorithmData{
		DataType: alogrithm.EEGData,
		Data:     task.Payload(),
	}
	redisKey := "eegSession:" + strconv.FormatUint(eegData.EegSessionId, 10) + ":data"
	sessionKey := "eegSession:" + strconv.FormatUint(eegData.EegSessionId, 10)

	exist, err := processor.rdb.Exists(ctx, sessionKey).Result()
	if err != nil {
		return err
	}
	if exist == 0 {
		return fmt.Errorf("eegSession not exist : %v", asynq.SkipRetry)
	}

	expiredAt, err := processor.rdb.HGet(ctx, redisKey, "expired_at").Result()
	if err != nil {
		return err
	}
	endAt, err := processor.rdb.HGet(ctx, redisKey, "end").Result()
	if err != nil {
		return err
	}
	expiredTime, err := util.StringToTime(expiredAt)
	if err != nil {
		return err
	}
	if expiredTime.Before(eegData.Time.AsTime()) {

		return nil
	}

	endTime, err := util.StringToTime(endAt)
	if err != nil {
		return err
	}

	if endTime.Before(eegData.Time.AsTime()) {
		return nil
	}

	count, _ := processor.rdb.SCard(ctx, redisKey).Result()
	if count >= 900 || eegData.IsEnd {
		records, err := processor.rdb.SMembers(ctx, redisKey).Result()
		if err != nil {
			log.Printf("failed to get records from redis : %v", err)
			return err
		}
		eegdatas := make([]*pb.EEGData, len(records))
		for i, record := range records {
			var data pb.EEGData
			if err := proto.Unmarshal([]byte(record), &data); err != nil {
				return err
			}
			eegdatas[i] = &data
		}
		err = processor.store.InsertEEGDataTx(ctx, database.InsertEEGDataTxParams{
			Records:   eegdatas,
			SessionId: eegData.EegSessionId,
		})
		if err != nil {
			return err
		}
		_, err = processor.rdb.Del(ctx, redisKey).Result()
		if err != nil {
			return err
		}
		if eegData.IsEnd {
			dataBitsMap, err := processor.rdb.Get(ctx, sessionKey+":dataBitMap").Result()
			if err != nil {
				log.Printf("get dataBitsMap error : %v\n", err)
			}
			_, err = processor.rdb.Del(ctx, sessionKey).Result()
			if err != nil {
				log.Printf("delete session info in redis failed : %v", err)
			}
			err = processor.store.UpdateEEGSession(ctx, database.UpdateEEGSessionParams{
				End: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
				ExpiredAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
				Finish: sql.NullBool{
					Bool:  true,
					Valid: true,
				},
				ID: eegData.EegSessionId,
				DataCount: sql.NullInt32{
					Int32: int32(len(eegdatas)),
					Valid: true,
				},
				DataBitsMap: sql.NullString{
					String: dataBitsMap,
					Valid:  true,
				},
			})
			if err != nil {
				return err
			}
		}
	} else {
		_, err := processor.rdb.SAdd(ctx, redisKey, task.Payload()).Result()
		if err != nil {
			log.Printf("save to redis failed : %v", err)
			return err
		}
	}

	_, err = processor.rdb.SetBit(ctx, sessionKey+":dataBitMap", int64(eegData.MessageOffset), 1).Result()
	if err != nil {
		log.Printf("save to bitmap failed : %v", err)
		return err
	}
	return nil
}
