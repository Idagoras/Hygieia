package api

import (
	"Hygieia/database"
	"Hygieia/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"net/http"
	"time"
)

type SignUpRequest struct {
	Mobile string `json:"mobile" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type SignUpResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

func (server *Server) SignUp(ctx *gin.Context) {
	var req SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	violations := validateSignUpRequest(&req)
	if violations != nil {
		err := util.InvalidArgumentError(violations)
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数格式不正确")
		return
	}
	smsKey := redisSmsKey + req.Mobile
	var code string
	var err error
	if code, err = server.rdb.Get(ctx, smsKey).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusNotFound, err, "过期验证码")
			return
		}
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "出现错误")
		return
	}
	if code != req.Code {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusNotFound, err, "验证码错误")
		return
	}

	user, err := server.store.CreateUserTx(ctx, database.InsertUserParams{
		Username:          util.RandomString(8),
		Mobile:            req.Mobile,
		Email:             "",
		Avatar:            "",
		HashedPassword:    util.RandomString(20),
		IsEmailVerified:   false,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		UpdatedUsernameAt: time.Now(),
		UpdatedAvatarAt:   time.Now(),
		UpdatedMobileAt:   time.Now(),
		UpdatedEmailAt:    time.Now(),
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "创建用户失败")
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "生成 token 错误")
		return
	}
	ctx.JSON(http.StatusOK, SignUpResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiresAt.Time,
	})
	return
}

func validateSignUpRequest(req *SignUpRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateChineseMobileNumber(req.Mobile); err != nil {
		violations = append(violations, util.FieldViolation("mobile", err))
	}
	if err := util.ValidateSecretCode(req.Code); err != nil {
		violations = append(violations, util.FieldViolation("code", err))
	}

	return
}
