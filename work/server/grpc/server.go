package server

import (
	"context"
	pb "work/api/grpc/v1"
	"work/model"
	"work/service"
)

type Service struct {
	svc *service.Service
}

func NewTaskService() *Service {
	srv := service.NewService()
	srv.InitializeServerConfiguration()
	return &Service{
		svc: srv,
	}
}

// Create 创建配置
func (s *Service) Create(ctx context.Context, req *pb.CreateOptions) (reply *pb.MessageResponse, err error) {
	create := &model.CreateOptions{
		Name:         req.Name,
		JoinServerId: req.JoinServerId,
		Subnet:       req.Subnet,
		Time:         req.Time,
		ListenPort:   req.ListenPort,
		DNS:          req.Dns,
		MTU:          req.Mtu,
		PublicIp:     req.PublicIp,
	}
	if req.NewServer {
		return s.svc.CreateServer(ctx, create)
	}
	return s.svc.CreateClient(ctx, create)
}

// Delete 删除配置
func (s *Service) Delete(ctx context.Context, req *pb.DeleteOptions) (*pb.MessageResponse, error) {
	if req.Server {
		return s.svc.DeleteServer(ctx, &model.DeleteOptions{
			Time: req.Time,
			Id:   req.Id,
			All:  req.All,
		})
	}
	return s.svc.DeleteClient(ctx, &model.DeleteOptions{
		Time: req.Time,
		Id:   req.Id,
		All:  req.All,
	})
}

func (s *Service) Start(ctx context.Context, options *pb.StartOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Restart(ctx context.Context, options *pb.RestartOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Stop(ctx context.Context, options *pb.StopOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Up(ctx context.Context, options *pb.UpOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Logs(ctx context.Context, options *pb.LogOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

// Ps 查看配置列表
func (s *Service) Ps(ctx context.Context, req *pb.PsOptions) (*pb.MessageResponse, error) {
	if req.Server {
		// 查看服务端
		return s.svc.PsServer(ctx, &model.PsOptions{})
	}
	return s.svc.PsClient(ctx, &model.PsOptions{})
}

func (s *Service) Remove(ctx context.Context, options *pb.RemoveOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Exec(ctx context.Context, options *pb.RunOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

// Show 查看指定名称配置
func (s *Service) Show(ctx context.Context, req *pb.ShowOptions) (*pb.MessageResponse, error) {
	if req.Server {
		// 查看服务端
		return s.svc.ShowServer(ctx, &model.ShowOptions{UserId: req.UserId})
	}
	return s.svc.ShowClient(ctx, &model.ShowOptions{UserId: req.UserId})
}

func (s *Service) Edit(ctx context.Context, options *pb.EditOptions) (*pb.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}
