package main

import (
	service "api-gateway/api/grpc/v1"
	"api-gateway/config"
	"api-gateway/routes"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 服务发现
	//etcdaddress := []string{viper.GetString("etcd.address")}
	//etcdRegister := discovery.NewRegister(etcdaddress, logrus.New())
	//resolver.Register(etcdRegister)
	go startListen()
	{
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignal
		fmt.Println("exit!", s)
	}
	fmt.Println("网关监听4000")
}

func startListen() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	//userConn, err := grpc.Dial(viper.GetString("domain.user"), opts...)
	//if err != nil {
	//	panic(err)
	//}
	//userService := service.NewUserServiceClient(userConn)

	taskrConn, err := grpc.Dial(config.ApiGatewayConf.GetString("domain.work"), opts...)
	if err != nil {
		panic(err)
	}
	taskService := service.NewServiceClient(taskrConn)

	r := routes.NewRouter(taskService)
	server := &http.Server{
		Addr:           ":" + config.ApiGatewayConf.GetString("server.apiGateway.port"),
		Handler:        r,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
