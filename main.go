package main

import (
	"context"
	"net"
	"os"
	plog_gateway "plog_gateway/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	conn   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
	err    error
	addr   = ":9999"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("failed to listen: %+v", err)
	}
	server := grpc.NewServer()
	plog_gateway.RegisterPLogGatewayServer(server, &PLogGatewayServer{})
	logrus.Infof("server listen at %s", lis.Addr().String())
	server.Serve(lis)
}

func init() {
	tmpAddr := os.Getenv("user_center_addr")
	if tmpAddr != "" {
		addr = tmpAddr
	}
}
