package main

import (
	"context"
	"fmt"
	"plog_gateway/handler"
	plog_gateway "plog_gateway/proto"
)

type PLogGatewayServer struct {
}

func (PLogGatewayServer) UploadLog(ctx context.Context, req *plog_gateway.UploadLogRequest) (resp *plog_gateway.UploadLogResponse, err error) {
	fmt.Println("UploadLog called")
	return handler.UploadLog(ctx, req)
}
