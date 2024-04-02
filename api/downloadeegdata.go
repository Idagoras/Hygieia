package api

import (
	"Hygieia/database"
	"Hygieia/oss"
	"Hygieia/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type DownloadEEGDataRequest struct {
	SessionId uint64 `json:"sessionId" binding:"required"`
}

type DownloadEEGDataResponse struct {
	DownloadURL string `json:"download_url"`
}

func (server *Server) DownloadEEGData(ctx *gin.Context) {
	// 懒上传策略
	var req DownloadEEGDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusBadRequest, err, "参数错误")
		return
	}
	objectName := "eegSession:" + strconv.FormatUint(req.SessionId, 10)
	exist, err := server.ossService.ExistObject(oss.BucketName, objectName)
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "无法上传")
		return
	}
	if exist == false {
		data, err := server.store.GetEEGDataByEEGSessionId(ctx, database.GetEEGDataByEEGSessionIdParams{
			SessionID: req.SessionId,
			Limit:     8000,
		})
		if err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "拉取数据失败")
			return
		}
		var jsonData []byte
		if jsonData, err = json.Marshal(data); err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "归档数据失败")
			return
		}
		filePath := "./Archive/" + "eegSession:" + strconv.FormatUint(req.SessionId, 10) + ".txt"
		zipPath := "./Archive/" + "eegSession:" + strconv.FormatUint(req.SessionId, 10) + ".zip"
		err = util.WriteToFile(filePath, jsonData)
		if err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "归档数据失败")
			return
		}

		err = util.Zip(filePath, zipPath)
		if err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "压缩数据失败")
			return
		}
		err = server.ossService.Upload(oss.BucketName, objectName, 10*time.Second, zipPath, 10)
		if err != nil {
			retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "上传数据失败")
			return
		}
	}
	url, err := server.ossService.GetDownloadURL(oss.BucketName, objectName, 10*time.Second, 10, 120)
	if err != nil {
		retJsonWithSpecificHttpStatusCode(ctx, http.StatusInternalServerError, err, "无法上传")
		return
	}
	ctx.JSON(http.StatusOK, DownloadEEGDataResponse{DownloadURL: url})
	return
}
