package middleware

import (
	"Hygieia/token"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey             = "authorization"
	AuthorizationTypeBearer            = "bearer"
	AuthorizationPayloadKey            = "authorization_payload"
	AuthorizationInvalidAccessTokenKey = "invalid_access_token"
)

type AuthorizationInfo struct {
	Token   string
	Payload *token.Payload
}

func Authorization(tokenMaker token.Maker, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		fmt.Printf("header:%s\n", authorizationHeader)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   err,
				"message": "缺少token授权",
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   err,
				"message": "不正确的token格式",
			})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   err,
				"message": "不支持的token类型",
			})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   err,
				"message": "无效token",
			})
			return
		}

		exists, err := rdb.Exists(ctx, AuthorizationInvalidAccessTokenKey+accessToken).Result()
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   err,
				"message": "出现错误",
			})
			return
		}
		if exists == 1 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   err,
				"message": "无效token",
			})
			return
		}

		ctx.Set(AuthorizationPayloadKey, AuthorizationInfo{
			Token:   accessToken,
			Payload: payload,
		})
		ctx.Next()
	}
}

func checkTokenIsInvalid(ctx context.Context, rdb *redis.Client, token string) (bool, error) {
	exists, err := rdb.Exists(ctx, AuthorizationInvalidAccessTokenKey+token).Result()
	if err != nil {
		return true, fmt.Errorf("failed to query redis %v", err)
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}
