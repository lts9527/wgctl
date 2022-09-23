package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	pb "work/api/grpc/v1"
	"work/log"
	"work/model"
	"work/pkg/util"
)

func (s *Service) CreateClient(ctx context.Context, co *model.CreateOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	if _, v := s.getClientNameMapping(co.Name); v {
		return reply, errors.New("name already exists")
	}
	_, bl := s.getServerAddress(co.JoinServerId)
	if !bl {
		return reply, errors.New("the added server does not exist")
	}
	PrivateKey, PublicKey := util.GenerateKeyPair()
	Address, err := s.getAddress(co.JoinServerId)
	if err != nil {
		log.Error(err.Error())
		return reply, errors.New("the added server does not exist")
	}
	UserConfig := &model.ConfigObjConfig{
		Time:                int32(time.Now().Unix()),
		Name:                co.Name,
		JoinServerId:        co.JoinServerId,
		ListenPort:          util.GetClientPort(),
		PrivateKey:          PrivateKey,
		PublicKey:           PublicKey,
		Address:             Address,
		DNS:                 co.DNS,
		MTU:                 co.MTU,
		AllowedIPs:          "0.0.0.0/0",
		PersistentKeepalive: "25",
		Endpoint:            co.PublicIp + ":" + s.getServerListenPort(co.JoinServerId),
	}
	if co.Name == "" {
		UserConfig.Name = util.RandStringlowercase(6)
	}
	uc, _ := json.Marshal(UserConfig)
	userID := fmt.Sprintf("%x", md5.Sum([]byte(uc)))
	UserConfig.UserId = userID
	s.ClientNameMapping[UserConfig.Name] = UserConfig
	uc, _ = json.Marshal(UserConfig)
	str := "/etc/wgctl/client/" + UserConfig.Name
	err = util.WriteFile(str, string(uc))
	if err != nil {
		log.Error(err.Error())
		return reply, errors.New(fmt.Sprintf("Failed to write client configuration %s", err.Error()))
	}
	// 重写最新的服务端配置
	err = util.SaveJoinServerConfig("/etc/wireguard/"+co.JoinServerId+".conf", UserConfig)
	if err != nil {
		return nil, err
	}
	configuration, err := util.GenerateClientConfiguration(UserConfig)
	if err != nil {
		return nil, err
	}
	err = util.WriteFile("/etc/wgctl/wireguard/"+UserConfig.Name, configuration)
	if err != nil {
		return nil, err
	}
	s.startWG(UserConfig.JoinServerId)
	reply.Name = UserConfig.Name
	reply.UserId = userID
	reply.WireguardConfig = configuration
	return
}
