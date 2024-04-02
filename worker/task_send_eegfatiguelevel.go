package worker

import (
	"Hygieia/pb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TaskSendEEGFatigueLevel = "task:send_eeg_fatigue_level"

type PayloadSendEEGFatigueLevel struct {
	FatigueLevel pb.EEGFatigueLevel `json:"fatigue_level"`
	SessionId    uint64             `json:"session_id"`
}

func (distributor *HyTaskDistributor) DistributeTaskSendEEGFatigueLevel(ctx context.Context, payload *PayloadSendEEGFatigueLevel, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendEEGFatigueLevel, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task : %w", err)
	}
	log.Printf("task type: %v,queue: %v,max retry:%d,enqueued task", info.Type, info.Queue, info.MaxRetry)
	return nil
}

func (processor *HyTaskProcessor) ProcessTaskSendEEGFatigueLevel(ctx context.Context, task *asynq.Task) error {
	processor.wsHub.FatigueLevel <- task.Payload()
	return nil
}
