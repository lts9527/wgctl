package service

import (
	"context"
	"errors"
	"io/ioutil"
	pb "work/api/grpc/v1"
	"work/log"
	"work/model"
)

func (s *Service) ShowClient(ctx context.Context, so *model.ShowOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	var configs []byte
	var Existence = false
	// 遍历查询是否存在这个配置名称
	if _, ok := s.getClientNameMapping(so.UserId); !ok {
		for _, v := range s.ClientNameMapping {
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
	configs, err = ioutil.ReadFile("/etc/wgctl/wireguard/" + so.UserId)
	if err != nil {
		log.Error(err.Error())
		return reply, err
	}
	reply.WireguardConfig = string(configs)
	return
}
