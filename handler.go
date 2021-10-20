package main

import (
	"context"
	plog_gateway "plog_gateway/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PLogGatewayServer struct {
}

func (PLogGatewayServer) UploadLog(ctx context.Context, req *plog_gateway.UploadLogRequest) (resp *plog_gateway.UploadLogResponse, err error) {

	return nil, status.Errorf(codes.Unimplemented, "method UploadLog not implemented")
}
