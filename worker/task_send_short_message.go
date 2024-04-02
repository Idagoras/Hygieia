package worker

import (
	"Hygieia/sms"
	"Hygieia/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const TaskSendShortMessage = "task:send_short_message"

type PayloadSendShortMessage struct {
	Mobile     string        `json:"email"`
	Expiration time.Duration `json:"expiration"`
}

func (distributor *HyTaskDistributor) DistributeTaskSendShortMessage(ctx context.Context, payload *PayloadSendShortMessage,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendShortMessage, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue eask: %w", err)
	}
	log.Printf("task type: %v,queue: %v,max retry:%d,enqueued task", info.Type, info.Queue, info.MaxRetry)
	return nil
}

func (processor *HyTaskProcessor) ProcessTaskSendShortMessage(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendShortMessage
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload : %w", asynq.SkipRetry)
	}
	code, err := processor.rdb.Get(ctx, sms.RedisSmsKey+payload.Mobile).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			code = util.GenerateSecretCode(sms.CodeLen)
			goto Send
		}
		return fmt.Errorf("failed to save code to redis : %v", err)
	}
Send:
	_, err = processor.rdb.Set(ctx, sms.RedisSmsKey+payload.Mobile, code, payload.Expiration*time.Second).Result()
	if err != nil {
		return fmt.Errorf("failed to save code to redis : %v", err)
	}
	err = processor.smsSender.SendShortMessage(ctx, payload.Mobile, code)
	if err != nil {
		return fmt.Errorf("failed to send short message : %v", err)
	}
	return nil

}
