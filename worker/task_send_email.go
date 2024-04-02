package worker

import (
	"Hygieia/mail"
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

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Email      string        `json:"email"`
	Expiration time.Duration `json:"expiration"`
}

func (distributor *HyTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue eask: %w", err)
	}
	log.Printf("task type: %v,queue: %v,max retry:%d,enqueued task", info.Type, info.Queue, info.MaxRetry)
	return nil
}

func (processor *HyTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload : %w", asynq.SkipRetry)
	}
	code, err := processor.rdb.Get(ctx, mail.RedisMailKey+payload.Email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			code = util.GenerateSecretCode(mail.CodeLen)
			goto Send
		}
		return fmt.Errorf("failed to save code to redis : %v", err)
	}
Send:
	_, err = processor.rdb.Set(ctx, mail.RedisMailKey+payload.Email, code, payload.Expiration*time.Second).Result()
	if err != nil {
		return fmt.Errorf("failed to save code to redis : %v", err)
	}
	subject := "Hygieia 验证邮件"
	content := fmt.Sprintf(`
			<h1>欢迎来到 Hygieia</h1>
			<p>您的验证码为 :%s </p>
			`, code)
	to := []string{payload.Email}
	err = processor.mailSender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email : %v", err)
	}
	return nil

}
