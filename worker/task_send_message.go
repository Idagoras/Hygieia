package worker

import (
	"Hygieia/database"
	"Hygieia/sse"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TaskSendMessage = "task:send_message"

type TaskSendMessagePayload struct {
	Message database.Message `json:"message"`
}

func (distributor *HyTaskDistributor) DistributeTaskSendMessage(ctx context.Context, payload TaskSendMessagePayload,
	opts ...asynq.Option) error {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload : %v", err)
	}
	task := asynq.NewTask(TaskSendMessage, jsonBytes, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue eask: %w", err)
	}
	log.Printf("task type: %v,queue: %v,max retry:%d,enqueued task", info.Type, info.Queue, info.MaxRetry)
	return nil
}

func (processor *HyTaskProcessor) ProcessTaskSendMessage(ctx context.Context, task *asynq.Task) error {
	var err error
	var payload TaskSendMessagePayload
	err = json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload to json : %v", err)
	}
	err = processor.store.InsertMessageTx(ctx, database.InsertMessageTxParams{
		Message: &payload.Message,
	})
	if err != nil {
		processor.broker.Notifier <- sse.NotificationEvent{
			EventName:  sse.EventMessageSendFailed,
			ReceiveUid: []uint64{payload.Message.SendUid},
			Payload: sse.EventMessageSendFailedPayload{
				Suid:      payload.Message.SendUid,
				RevUid:    payload.Message.RcvUid,
				CreatedAt: payload.Message.CreatedAt,
			},
		}
	} else {
		processor.broker.Notifier <- sse.NotificationEvent{
			EventName:  sse.EventMessageSendSuccess,
			ReceiveUid: []uint64{payload.Message.SendUid},
			Payload: sse.EventMessageSendSuccessPayload{
				Suid:      payload.Message.SendUid,
				RevUid:    payload.Message.RcvUid,
				CreatedAt: payload.Message.CreatedAt,
			},
		}
		processor.broker.Notifier <- sse.NotificationEvent{
			EventName:  sse.EventMessageComing,
			ReceiveUid: []uint64{payload.Message.RcvUid},
			Payload:    sse.EventMessageComingPayload{Message: payload.Message},
		}
	}
	return err
}
