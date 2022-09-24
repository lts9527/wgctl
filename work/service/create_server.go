package service

import (
	"context"
	"errors"
	"strconv"
	"time"
	pb "work/api/grpc/v1"
	"work/log"
	"work/model"
	"work/pkg/util"
)

func (s *Service) CreateServer(ctx context.Context, co *model.CreateOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	if _, ok := s.getServerNameMapping(co.Name); ok {
		return reply, errors.New("name already exists")
	}
	Address := s.setIpPool(co.Name, co.Subnet)
	PrivateKey, PublicKey := util.GenerateKeyPair()
	UserConfig := &model.ConfigObjConfig{
		Time:                int32(time.Now().Unix()),
		Name:                co.Name,
		ListenPort:          strconv.Itoa(s.getListenPort()),
		PrivateKey:          PrivateKey,
		PublicKey:           PublicKey,
		Address:             Address,
		DNS:                 co.DNS,
		MTU:                 co.MTU,
		AllowedIPs:          "0.0.0.0/0",
		PersistentKeepalive: "25",
	}
	ServerConfig, err := s.SaveServerConfig(UserConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// 写入服务端缓存
	s.ServerNameMapping[UserConfig.Name] = ServerConfig
	s.startWG(UserConfig.Name)
	reply.UserId = ServerConfig.UserId
	reply.Name = UserConfig.Name
	return
}
