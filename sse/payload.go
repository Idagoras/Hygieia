package sse

import (
	"Hygieia/database"
	"time"
)

type EventMessageSendFailedPayload struct {
	Suid      uint64    `json:"suid"`
	RevUid    uint64    `json:"rev_uid"`
	CreatedAt time.Time `json:"created_at"`
}

type EventMessageSendSuccessPayload struct {
	Suid      uint64    `json:"suid"`
	RevUid    uint64    `json:"rev_uid"`
	CreatedAt time.Time `json:"created_at"`
}

type EventMessageComingPayload struct {
	Message database.Message
}
