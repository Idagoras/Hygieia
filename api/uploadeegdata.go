package api

import (
	"Hygieia/database"
	"Hygieia/pb"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UploadEEGDataRequest struct {
	SessionId  uint64   `json:"session_id" binding:"required"`
	EEGDataArr []string `json:"eeg_data_arr" binding:"required"`
}

type UploadEEGDataResponse struct {
	SessionId   uint64 `json:"session_id"`
	DataBitsMap string `json:"data_bits_map"`
}

func (server *Server) UploadEEGData(ctx *gin.Context) {
	var req UploadEEGDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数有误")
		return
	}

	// 判断 redis 中有没有该 session 的 bitsMap
	var bitsMapString string
	var session database.EegSession
	var err error
	exists, err := server.rdb.Exists(ctx, "eegsession:"+strconv.FormatUint(req.SessionId, 10)+":dataBitMap").Result()
	if err != nil {
		log.Printf("failed to query redis : %v", err)
	}
	if exists == 1 {
		bitsMapString, err = server.rdb.Get(ctx, "eegsession:"+strconv.FormatUint(req.SessionId, 10)+":dataBitMap").Result()
		if err != nil {
			log.Printf("failed to get bitsMap from redis :%v", err)
			goto queryDataBase
		}
		goto insert

	}
queryDataBase:
	session, err = server.store.GetEEGSessionByID(ctx, req.SessionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusNotFound, err, "不存在的会话")
			return
		}
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "服务器内部错误")
		return
	}
	bitsMapString = session.DataBitsMap
insert:
	bitsMapBytes := []byte(bitsMapString)
	var records []*pb.EEGData = make([]*pb.EEGData, len(req.EEGDataArr))
	for _, record := range req.EEGDataArr {
		var eegData pb.EEGData
		if err := proto.Unmarshal([]byte(record), &eegData); err != nil {
			log.Printf("unmarshal load data failed : %v", err)
			continue
		}
		offset := eegData.MessageOffset
		byteIndex := offset / 8
		byteInternalOffset := offset % 8
		if bitsMapBytes[byteIndex]&(1<<byteInternalOffset) > 0 {
			records = append(records, &eegData)
			bitsMapBytes[byteIndex] |= 1 << byteInternalOffset
		}
	}
	err = server.store.InsertEEGDataTx(ctx, database.InsertEEGDataTxParams{
		Records:   records,
		SessionId: req.SessionId,
	})
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "插入数据失败")
		return
	}
	_, err = server.rdb.Set(ctx, "eegsession:"+strconv.FormatUint(req.SessionId, 10)+":dataBitMap", string(bitsMapBytes), 10*time.Minute).Result()
	if err != nil {
		log.Printf("failed to save bitsmap to redis : %v", err)
	}
	ctx.JSON(http.StatusOK, UploadEEGDataResponse{
		SessionId:   req.SessionId,
		DataBitsMap: string(bitsMapBytes),
	})
	return
}
