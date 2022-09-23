package main

import (
	"google.golang.org/grpc"
	"net"
	service "work/api/grpc/v1"
	"work/config"
	srv "work/server/grpc"
)

func main() {
	grpcAddress := config.WorkConf.GetString("server.work.address") + ":" + config.WorkConf.GetString("server.work.port")
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定服务
	service.RegisterServiceServer(server, srv.NewTaskService())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
