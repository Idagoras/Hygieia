package worker

import (
	"Hygieia/alogrithm"
	"Hygieia/database"
	"Hygieia/mail"
	"Hygieia/sms"
	"Hygieia/sse"
	"Hygieia/wb"
	"context"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
}

type HyTaskProcessor struct {
	server     *asynq.Server
	hub        *alogrithm.AlgorithmHub
	wsHub      *wb.Hub
	rdb        *redis.Client
	store      database.Store
	broker     *sse.Broker
	mailSender mail.Sender
	smsSender  sms.Sender
}

func NewHyTaskProcessor(redisOpt asynq.RedisClientOpt, hub *alogrithm.AlgorithmHub) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Printf("process task failed, task type: %v, payload: %v", task.Type(), task.Payload())
		}),
	})
	return &HyTaskProcessor{server: server, hub: hub}
}

func (processor *HyTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	return processor.server.Start(mux)

}
