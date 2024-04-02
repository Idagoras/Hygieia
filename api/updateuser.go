package api

import (
	"Hygieia/database"
	"Hygieia/mail"
	"Hygieia/middleware"
	"Hygieia/util"
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"net/http"
	"time"
)

type UpdateUserBasicInfoRequest struct {
	UserName *string `json:"userName"`
	Avatar   *string `json:"avatar"`
}

type UpdateUserMobileRequest struct {
	NewMobile string `json:"newMobile" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

type UpdateUserEmailRequest struct {
	NewEmail string `json:"newEmail" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type UpdateUserPasswordRequest struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type UpdateUserBasicInfoResponse struct {
	Success    int       `json:"success"`
	UpdateTime time.Time `json:"update_time"`
}

type UpdateUserMobileResponse struct {
	Success    int       `json:"success"`
	UpdateTime time.Time `json:"updateTime"`
}

type UpdateUserEmailResponse struct {
	Success    int       `json:"success"`
	UpdateTime time.Time `json:"updateTime"`
}

type UpdateUserPasswordResponse struct {
	Success    int       `json:"success"`
	UpdateTime time.Time `json:"updateTime"`
}

func (server *Server) UpdateUserBasicInfo(ctx *gin.Context) {
	authInfo := ctx.MustGet(middleware.AuthorizationPayloadKey).(middleware.AuthorizationInfo)
	var req UpdateUserBasicInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	violations := validateUpdateUserBasicInfoRequest(&req)
	if violations != nil {
		err := util.InvalidArgumentError(violations)
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数格式不正确")
		return
	}
	_, err := server.store.UpdateUserTx(ctx, database.UpdateUserTxParams{
		UpdateUserParams: database.UpdateUserParams{
			Username: sql.NullString{
				String: *req.UserName,
				Valid:  req.UserName != nil,
			},
			Avatar: sql.NullString{
				String: *req.Avatar,
				Valid:  req.Avatar != nil,
			},
			ID: authInfo.Payload.UID,
		},
		AfterCreate: func(user *database.User) error {
			return nil
		},
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "更新失败")
		return
	}
	ctx.JSON(http.StatusOK, UpdateUserBasicInfoResponse{
		Success:    1,
		UpdateTime: time.Now(),
	})
	return
}

func (server *Server) UpdateUserMobile(ctx *gin.Context) {
	authInfo := ctx.MustGet(middleware.AuthorizationPayloadKey).(middleware.AuthorizationInfo)
	var req UpdateUserMobileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	violations := validateUpdateUserMobileRequest(&req)
	if violations != nil {
		err := util.InvalidArgumentError(violations)
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数格式不正确")
		return
	}
	smsKey := redisSmsKey + req.NewMobile
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

	user, err := server.store.UpdateUserTx(ctx, database.UpdateUserTxParams{
		UpdateUserParams: database.UpdateUserParams{
			Mobile: sql.NullString{
				String: req.NewMobile,
				Valid:  true,
			},
			UpdatedMobileAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			ID: authInfo.Payload.UID,
		},
		AfterCreate: func(user *database.User) error {
			return makeAccessTokenInvalid(ctx, server.rdb, authInfo.Token, server.config.AccessTokenDuration)
		},
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "更新失败")
		return
	}
	ctx.JSON(http.StatusOK, UpdateUserMobileResponse{
		Success:    1,
		UpdateTime: user.UpdatedMobileAt,
	})
}

func (server *Server) UpdateUserEmail(ctx *gin.Context) {
	authInfo := ctx.MustGet(middleware.AuthorizationPayloadKey).(middleware.AuthorizationInfo)
	var req UpdateUserEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	violations := validateUpdateUserEmailRequest(&req)
	if violations != nil {
		err := util.InvalidArgumentError(violations)
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数格式不正确")
		return
	}
	mailKey := mail.RedisMailKey + req.NewEmail
	var code string
	var err error
	if code, err = server.rdb.Get(ctx, mailKey).Result(); err != nil {
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

	user, err := server.store.UpdateUserTx(ctx, database.UpdateUserTxParams{
		UpdateUserParams: database.UpdateUserParams{
			ID: authInfo.Payload.UID,
			Email: sql.NullString{
				String: req.NewEmail,
				Valid:  true,
			},
			UpdatedEmailAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
		AfterCreate: func(user *database.User) error {
			return nil
		},
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "更新失败")
		return
	}
	ctx.JSON(http.StatusOK, UpdateUserEmailResponse{
		Success:    1,
		UpdateTime: user.UpdatedEmailAt,
	})
}

func (server *Server) UpdateUserPassword(ctx *gin.Context) {
	authInfo := ctx.MustGet(middleware.AuthorizationPayloadKey).(middleware.AuthorizationInfo)
	var req UpdateUserPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	if req.Password == req.NewPassword {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, nil, "无法修改为相同的密码")
		return
	}
	violations := validateUpdateUserPasswordRequest(&req)
	if violations != nil {
		err := util.InvalidArgumentError(violations)
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数格式不正确")
		return
	}
	hashedPassword, err := util.HashPassword(req.NewPassword)
	user, err := server.store.UpdateUserTx(ctx, database.UpdateUserTxParams{
		UpdateUserParams: database.UpdateUserParams{
			HashedPassword: sql.NullString{
				String: hashedPassword,
				Valid:  true,
			},
			PasswordChangedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			ID: authInfo.Payload.UID,
		},
		AfterCreate: func(user *database.User) error {
			return makeAccessTokenInvalid(ctx, server.rdb, authInfo.Token, server.config.AccessTokenDuration)
		},
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "更新失败")
		return
	}
	ctx.JSON(http.StatusOK, UpdateUserPasswordResponse{
		Success:    1,
		UpdateTime: user.PasswordChangedAt,
	})
}

func makeAccessTokenInvalid(ctx context.Context, rdb *redis.Client, token string, expiration time.Duration) error {
	_, err := rdb.Set(ctx, middleware.AuthorizationInvalidAccessTokenKey+token, token, expiration).Result()
	return err
}

func validateUpdateUserEmailRequest(req *UpdateUserEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateEmail(req.NewEmail); err != nil {
		violations = append(violations, util.FieldViolation("email", err))
	}
	return
}

func validateUpdateUserMobileRequest(req *UpdateUserMobileRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateChineseMobileNumber(req.NewMobile); err != nil {
		violations = append(violations, util.FieldViolation("mobile", err))
	}
	return
}

func validateUpdateUserPasswordRequest(req *UpdateUserPasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidatePassword(req.Password); err != nil {
		violations = append(violations, util.FieldViolation("password", err))
	}
	if err := util.ValidatePassword(req.NewPassword); err != nil {
		violations = append(violations, util.FieldViolation("newPassword", err))
	}
	return
}

func validateUpdateUserBasicInfoRequest(req *UpdateUserBasicInfoRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(*req.UserName); err != nil {
		violations = append(violations, util.FieldViolation("username", err))
	}
	return
}
