package api

import (
	"Hygieia/database"
	"Hygieia/worker"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SendMessageRequest struct {
	Suid        uint64    `json:"suid"`
	RevUid      uint64    `json:"rev_uid"`
	CreatedAt   time.Time `json:"created_at"`
	MessageType uint8     `json:"message_type"`
	Text        string    `json:"text"`
	SubType     uint8     `json:"subType"`
	Attachment  string    `json:"attachment"`
}

func (server *Server) SendMessage(ctx *gin.Context) {
	var req SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	err := server.distributor.DistributeTaskSendMessage(ctx, worker.TaskSendMessagePayload{Message: database.Message{
		ID:         0,
		SendUid:    req.Suid,
		RcvUid:     req.RevUid,
		CreatedAt:  req.CreatedAt,
		HasRead:    false,
		Type:       uint32(req.MessageType),
		Text:       req.Text,
		Subtype:    uint32(req.SubType),
		Attachment: req.Attachment,
	}})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "发送失败")
		return
	}
	ctx.JSON(http.StatusOK, nil)
	return
}

type ListUserReceivedMessageRequest struct {
	Uid    uint64 `json:"uid"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListUserReceivedMessageResponse struct {
	Messages []database.Message `json:"messages"`
}

func (server *Server) ListUserReceivedMessage(ctx *gin.Context) {
	var req ListUserReceivedMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	messages, err := server.store.ListUserReceivedMessages(ctx, database.ListUserReceivedMessagesParams{
		RcvUid: req.Uid,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "发送失败")
		return
	}
	ctx.JSON(http.StatusOK, ListUserReceivedMessageResponse{Messages: messages})
	return
}

type ListUserSentMessageRequest struct {
	Uid    uint64 `json:"uid"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListUserSentMessageResponse struct {
	Messages []database.Message `json:"messages"`
}

func (server *Server) ListUserSentMessage(ctx *gin.Context) {
	var req ListUserSentMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	messages, err := server.store.ListUserSentMessages(ctx, database.ListUserSentMessagesParams{
		SendUid: req.Uid,
		Limit:   req.Limit,
		Offset:  req.Offset,
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "发送失败")
		return
	}
	ctx.JSON(http.StatusOK, ListUserSentMessageResponse{Messages: messages})
	return
}
