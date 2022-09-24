package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	pb "work/api/grpc/v1"
	"work/config"
	"work/log"
	"work/model"
	"work/pkg/util"
)

func (s *Service) CreateClient(ctx context.Context, co *model.CreateOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	if _, v := s.getClientNameMapping(co.Name); v {
		return reply, errors.New("name already exists")
	}
	if _, bl := s.getServerAddress(co.JoinServerId); !bl {
		return reply, errors.New("the added server does not exist")
	}
	PrivateKey, PublicKey := util.GenerateKeyPair()
	Address, err := s.getClientAddress(co.JoinServerId)
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
	uc, _ = json.Marshal(UserConfig)
	if err = util.WriteFile(config.WorkConf.GetString("wireguard.wgctlClientDir")+UserConfig.Name, string(uc)); err != nil {
		log.Error(err.Error())
		return reply, errors.New(fmt.Sprintf("Failed to write client configuration %s", err.Error()))
	}
	// 重写最新的服务端配置
	if err = util.SaveJoinServerConfig("/etc/wireguard/"+co.JoinServerId+".conf", UserConfig); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	configuration, err := util.GenerateClientConfiguration(UserConfig)
	if err != nil {
		return nil, err
	}
	if err = util.WriteFile("/etc/wgctl/wireguard/"+UserConfig.Name, configuration); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// 新建的客户端写入缓存
	s.ClientNameMapping[UserConfig.Name] = UserConfig
	s.startWG(UserConfig.JoinServerId)
	reply.Name = UserConfig.Name
	reply.UserId = userID
	reply.WireguardConfig = configuration
	return
}
