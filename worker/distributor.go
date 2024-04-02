package worker

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendEEGData(ctx context.Context, eegDataBytes []byte, opts ...asynq.Option) error
	DistributeTaskSendEEGFatigueLevel(ctx context.Context, payload *PayloadSendEEGFatigueLevel, opts ...asynq.Option) error
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail,
		opts ...asynq.Option) error
	DistributeTaskSendMessage(ctx context.Context, payload TaskSendMessagePayload,
		opts ...asynq.Option) error
	DistributeTaskSendShortMessage(ctx context.Context, payload *PayloadSendShortMessage,
		opts ...asynq.Option) error
}

type HyTaskDistributor struct {
	client *asynq.Client
}

func NewHyTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &HyTaskDistributor{client: client}
}
