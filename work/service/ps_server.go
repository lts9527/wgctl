package service

import (
	"context"
	pb "work/api/grpc/v1"
	"work/config"
	"work/log"
	"work/model"
	"work/pkg/util"
)

func (s *Service) PsServer(ctx context.Context, po *model.PsOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	sl, err := util.FileForEach(config.WorkConf.GetString("wireguard.wgConfigDir"))
	if err != nil {
		log.Error(err.Error())
		return reply, err
	}
	// 遍历当前服务端的wireguard接口
	for _, v := range sl {
		name := s.formatFileName(v.Name(), ".conf")
		reply.Ps = append(reply.Ps, &pb.PsOptions{
			WgConfigId: s.ServerNameMapping[name].UserId,
			Created:    s.formatTimeFormat(s.getCreateTime(int64(s.ServerNameMapping[name].Time))),
			Status:     "Up",
			Ports:      s.ServerNameMapping[name].ListenPort,
			Names:      name,
		})
	}
	return
}
