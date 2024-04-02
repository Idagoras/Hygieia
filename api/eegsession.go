package api

import (
	"Hygieia/database"
	"Hygieia/util"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type beginEEGSessionRequest struct {
	SessionId uint64 `json:"sessionId" binding:"required"`
	Reconnect int    `json:"reconnect" binding:"required"`
}

type beginEEGSessionResponse struct {
	SessionId uint64 `json:"session_id" binding:"required"`
}

func retJsonWithSpecificHttpStatusCode(ctx *gin.Context, httpCode int, err error, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": "服务器内部错误",
		"error:":  err.Error(),
	})
}

func (server *Server) BeginEEGSession(ctx *gin.Context) {
	var req beginEEGSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println(ctx.Request.URL)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数有误",
			"error:":  err.Error(),
		})
		return
	}
	var exist int64 = 0
	var err error
	sessionIdStr := strconv.FormatUint(req.SessionId, 10)
	if req.Reconnect == 1 {
		exist, err = server.rdb.Exists(ctx, "eegSession:"+sessionIdStr).Result()
		if err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "内部错误")
			return
		}
	}
	existButExpired := false
	if exist == 1 {
		session, err := server.rdb.HGetAll(ctx, "eegSession:"+sessionIdStr).Result()
		if err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "内部错误")
			return
		}
		// 判断  uid 是否相等，不能开始别人的session
		expiredAt := session["expired_at"]
		expiredTime, err := util.StringToTime(expiredAt)

		fmt.Println(err, expiredTime, time.Now())

		if err != nil {
			fmt.Println(err)
			exist = 0
			existButExpired = true
			goto notExist
		}
		if expiredTime.Before(time.Now()) {
			exist = 0
			existButExpired = true
			goto notExist
		}

		response := beginEEGSessionResponse{SessionId: req.SessionId}

		ctx.JSON(http.StatusOK, response)
		return
	}

notExist:
	if exist == 0 {
		if existButExpired == false {
			session, err := server.store.GetEEGSessionByID(ctx, req.SessionId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					goto createNew
				} else {
					retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "内部错误")
					return
				}
			}
			// 验证 uid
			expiredTime := session.ExpiredAt
			if expiredTime.Before(time.Now()) {
				goto createNew
			}
			sessionStr := strconv.FormatUint(session.ID, 10)
			response := beginEEGSessionResponse{SessionId: req.SessionId}
			_, err = server.rdb.HMSet(ctx, "eegSession:"+sessionStr, "uid", strconv.FormatUint(session.Uid, 10), "begin", session.Begin.String(), "end", session.End.String(), "data_count", session.DataCount, "expired_at", session.ExpiredAt.String()).Result()
			if err != nil {
				fmt.Errorf("failed to add session to redis :%v", err)
			}
			_, err = server.rdb.Expire(ctx, "eegSession:"+sessionStr, time.Hour*4).Result()
			if err != nil {
				fmt.Printf("failed to add expire time to session in redis :%v", err)
			}
			ctx.JSON(http.StatusOK, response)
			return

		}
	createNew:
		session, err := server.store.CreateEEGSessionTx(ctx, database.CreateEEGSessionParams{
			Uid:       1,
			Begin:     time.Now(),
			End:       time.Now().Add(1 * time.Minute),
			DataCount: 0,
			ExpiredAt: time.Now().Add(1 * time.Minute),
		})
		if err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "无法创建新传输")
			return
		}
		sessionStr := strconv.FormatUint(session.ID, 10)

		_, err = server.rdb.HMSet(ctx, "eegSession:"+sessionStr, "uid", strconv.FormatUint(session.Uid, 10), "begin", session.Begin.String(), "end", session.End.String(), "data_count", session.DataCount, "expired_at", session.ExpiredAt.String()).Result()
		if err != nil {
			fmt.Printf("failed to add session to redis :%v", err)
		}
		_, err = server.rdb.Expire(ctx, "eegSession:"+sessionStr, time.Hour*4).Result()
		if err != nil {
			fmt.Printf("failed to add expire time to session in redis :%v", err)
		}
		_, err = server.rdb.SetBit(ctx, "eegSession:"+sessionStr+":dataBitMap", 0, 0).Result()
		if err != nil {
			fmt.Printf("failed to add bitmap to redis : %v", err)
		}
		ctx.JSON(http.StatusOK, beginEEGSessionResponse{SessionId: session.ID})
	}

}
