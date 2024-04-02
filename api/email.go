package api

import (
	"Hygieia/worker"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SendEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func (server *Server) SendEmail(ctx *gin.Context) {
	var req SendEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	err := server.distributor.DistributeTaskSendVerifyEmail(ctx, &worker.PayloadSendVerifyEmail{
		Email:      req.Email,
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
