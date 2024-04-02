package api

import (
	"Hygieia/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LogoutResponse struct {
	Success int    `json:"success"`
	Message string `json:"message"`
}

func (server *Server) Logout(ctx *gin.Context) {
	authInfo := ctx.MustGet(middleware.AuthorizationPayloadKey).(middleware.AuthorizationInfo)
	_, err := server.rdb.SetEx(ctx, middleware.AuthorizationInvalidAccessTokenKey+authInfo.Token, authInfo.Payload.UID, server.config.AccessTokenDuration).Result()
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "出现错误")
		return
	}
	ctx.JSON(http.StatusOK, LogoutResponse{
		Success: 1,
		Message: "退出登录成功",
	})
	return
}
