package api

import (
	"Hygieia/util"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"net/http"
	"time"
)

type LoginWithPasswordRequest struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginWithSecretCodeRequest struct {
	Mobile string `json:"mobile" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type LoginResponse struct {
	Uid                  uint64
	UserName             string
	Mobile               string
	Email                string
	AccessToken          string
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) LoginWithPassword(ctx *gin.Context) {
	var req LoginWithPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	violations := validateLoginWithPasswordRequest(&req)
	if violations != nil {
		err := util.InvalidArgumentError(violations)
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数格式不正确")
		return
	}
	user, err := server.store.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusForbidden, err, "用户不存在")
			return
		}
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "出现错误")
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "密码错误")
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "生成 token 错误")
		return
	}
	ctx.JSON(http.StatusOK, LoginResponse{
		Uid:                  user.ID,
		UserName:             user.Username,
		Mobile:               user.Mobile,
		Email:                user.Email,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiresAt.Time,
	})
	return
}

func (server *Server) LoginWithSecretCode(ctx *gin.Context) {
	var req LoginWithSecretCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	violations := validateLoginWithSecretCodeRequest(&req)
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
	user, err := server.store.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusForbidden, err, "用户不存在")
			return
		}
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "出现错误")
		return
	}
	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "生成 token 错误")
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Uid:                  user.ID,
		UserName:             user.Username,
		Mobile:               user.Mobile,
		Email:                user.Email,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiresAt.Time,
	})
	return
}

func validateLoginWithPasswordRequest(req *LoginWithPasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateChineseMobileNumber(req.Mobile); err != nil {
		violations = append(violations, util.FieldViolation("mobile", err))
	}
	if err := util.ValidatePassword(req.Password); err != nil {
		violations = append(violations, util.FieldViolation("password", err))
	}
	return
}

func validateLoginWithSecretCodeRequest(req *LoginWithSecretCodeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateChineseMobileNumber(req.Mobile); err != nil {
		violations = append(violations, util.FieldViolation("mobile", err))
	}
	if err := util.ValidateSecretCode(req.Code); err != nil {
		violations = append(violations, util.FieldViolation("code", err))
	}
	return
}
