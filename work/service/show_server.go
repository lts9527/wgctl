package service

import (
	"context"
	"errors"
	pb "work/api/grpc/v1"
	"work/config"
	"work/log"
	"work/model"
	"work/pkg/util"
)

func (s *Service) ShowServer(ctx context.Context, so *model.ShowOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	var Existence = false
	// 遍历查询是否存在这个配置名称
	if _, ok := s.getServerNameMapping(so.UserId); !ok {
		for _, v := range s.ServerNameMapping {
			if v.UserId[:9] == so.UserId {
				so.UserId = v.Name
				Existence = true
				break
			}
			if v.UserId == so.UserId {
				so.UserId = v.Name
				Existence = true
				break
			}
		}
	} else {
		Existence = true
	}
	if !Existence {
		return reply, errors.New("Error: No such Name: " + so.UserId)
	}
	data, err := util.ReadFile(config.WorkConf.GetString("wireguard.wgConfigDir") + so.UserId + ".conf")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	reply.WireguardConfig = data
	return
}
