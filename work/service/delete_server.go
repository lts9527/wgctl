package service

import (
	"context"
	"os"
	pb "work/api/grpc/v1"
	"work/config"
	"work/log"
	"work/model"
)

func (s *Service) DeleteServer(ctx context.Context, do *model.DeleteOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	var resp = &pb.DeleteOptions{}
	if do.All {
		for _, v := range s.ServerNameMapping {
			if err = os.Remove("/etc/wireguard/" + v.Name + ".conf"); err != nil {
				log.Error(err.Error())
				continue
			}
			if err = os.Remove("/etc/wgctl/server/" + v.Name); err != nil {
				log.Error(err.Error())
				continue
			}
			resp.Existence = append(resp.Existence, v.Name)
			delete(s.ServerNameMapping, v.Name)
			s.startWG(v.Name)
		}
		reply.Delete = resp
		return
	}
	idList := make(map[string]*deleteEnvMap)
	for _, v := range do.Id {
		idList[v] = &deleteEnvMap{}
	}
	for k, _ := range idList {
		var env = &deleteEnvMap{}
		if kk, v := s.getServerNameMapping(k); v {
			env.Address = kk.Address
			env.Name = kk.Name
			idList[k] = env
			continue
		}
		for _, vv := range s.ServerNameMapping {
			if vv.UserId == k {
				env.Address = vv.Address
				env.Name = vv.Name
				idList[k] = env
				continue
			}
			if vv.UserId[:9] == k {
				env.Address = vv.Address
				env.Name = vv.Name
				idList[k] = env
			}
		}
	}
	for k, _ := range idList {
		if idList[k].Name != "" {
			s.stopWG(idList[k].Name)
			if err = os.Remove(config.WorkConf.GetString("wireguard.wgConfigDir") + idList[k].Name + ".conf"); err != nil {
				log.Error(err.Error())
				continue
			}
			if err = os.Remove(config.WorkConf.GetString("wireguard.wgctlServerDir") + idList[k].Name); err != nil {
				log.Error(err.Error())
				continue
			}
			resp.Existence = append(resp.Existence, k)
			delete(s.ServerNameMapping, idList[k].Name)
			continue
		}
		resp.DoesNotExist = append(resp.DoesNotExist, k)
	}
	reply.Delete = resp
	return
}
