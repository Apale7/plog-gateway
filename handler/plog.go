package handler

import (
	"context"
	"fmt"
	"plog_gateway/MQ/redis"
	plog_gateway "plog_gateway/proto"
)

func UploadLog(ctx context.Context, req *plog_gateway.UploadLogRequest) (resp *plog_gateway.UploadLogResponse, err error) {
	redis.Publish("ropz", req.Log.Message)
	fmt.Println(req.Log.Message)
	resp = &plog_gateway.UploadLogResponse{
		Code: 1,
		Msg:  "",
	}
	return
}
