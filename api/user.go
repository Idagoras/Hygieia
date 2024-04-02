package api

import (
	"Hygieia/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubscribeRequest struct {
	Suid uint64 `json:"suid"`
	Like int    `json:"like"`
}

type SubscribeResponse struct {
	Success int `json:"success"`
}

func (server *Server) Subscribe(ctx *gin.Context) {
	var req SubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	err := server.store.SubscribeTx(ctx, database.SubscribeTxParams{
		Suid: req.Suid,
		Uid:  req.Suid,
		Like: req.Like > 0,
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "关注失败")
		return
	}
	ctx.JSON(http.StatusOK, SubscribeResponse{Success: 1})
	return
}

type GetUserSubscribersListRequest struct {
	Uid    uint64 `form:"uid"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type GetUserSubscribersListResponse struct {
	Subscribers []database.Subscriber `json:"subscribers"`
}

func (server *Server) GetUserSubscribersList(ctx *gin.Context) {
	var req GetUserSubscribersListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	sbs, err := server.store.ListSubscribePairByUid(ctx, database.ListSubscribePairByUidParams{
		Uid:    req.Uid,
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "查询失败")
		return
	}
	ctx.JSON(http.StatusOK, sbs)
	return

}

type GetSubscribedUsersListRequest struct {
	Uid    uint64 `form:"uid"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type GetSubscribedUsersListResponse struct {
	Subscribers []database.Subscriber `json:"subscribers"`
}

func (server *Server) GetSubscribedUsersList(ctx *gin.Context) {
	var req GetSubscribedUsersListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}
	sbs, err := server.store.ListSubscribePairBySuid(ctx, database.ListSubscribePairBySuidParams{
		SbsID:  req.Uid,
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "查询失败")
		return
	}
	ctx.JSON(http.StatusOK, sbs)
	return
}
