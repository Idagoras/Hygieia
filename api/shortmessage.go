package api

import (
	"Hygieia/worker"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SendShortMessageRequest struct {
	Mobile string `json:"mobile" binding:"required"`
}

func (server *Server) SendShortMessage(ctx *gin.Context) {
	var req SendShortMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	err := server.distributor.DistributeTaskSendShortMessage(ctx, &worker.PayloadSendShortMessage{
		Mobile:     req.Mobile,
		Expiration: server.config.AccessTokenDuration,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "发送服务异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, nil)
	return
}
