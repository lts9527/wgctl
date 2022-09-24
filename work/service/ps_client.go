package service

import (
	"context"
	"strings"
	pb "work/api/grpc/v1"
	"work/log"
	"work/model"
	"work/pkg/util"
)

func (s *Service) PsClient(ctx context.Context, po *model.PsOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	sl, err := util.FileForEach("/etc/wgctl/wireguard/")
	if err != nil {
		log.Error(err.Error())
		return reply, err
	}
	// 遍历当前的客户端配置
	for _, v := range sl {
		if !strings.Contains(v.Name(), ".gitignore") {
			reply.Ps = append(reply.Ps, &pb.PsOptions{
				WgConfigId: s.ClientNameMapping[v.Name()].UserId,
				Status:     s.formatTimeFormat(s.getCreateTime(int64(s.ClientNameMapping[v.Name()].Time))),
				Ports:      s.ClientNameMapping[v.Name()].ListenPort,
				Names:      v.Name(),
			})
		}
	}
	return
}
