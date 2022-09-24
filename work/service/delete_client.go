package service

import (
	"context"
	"fmt"
	"os"
	pb "work/api/grpc/v1"
	"work/log"
	"work/model"
)

type deleteEnvMap struct {
	Name    string
	Join    string
	Address string
}

func (s *Service) DeleteClient(ctx context.Context, do *model.DeleteOptions) (reply *pb.MessageResponse, err error) {
	reply = new(pb.MessageResponse)
	var resp = &pb.DeleteOptions{}
	if do.All {
		for _, v := range s.ClientNameMapping {
			if err = os.Remove("/etc/wgctl/wireguard/" + v.Name); err != nil {
				log.Error(err.Error())
				continue
			}
			if err = os.Remove("/etc/wgctl/client/" + v.Name); err != nil {
				log.Error(err.Error())
				continue
			}
			if err = s.deleteClientConfig(fmt.Sprintf("/etc/wireguard/%s", v.JoinServerId+".conf"), v.Address); err != nil {
				log.Error(err.Error())
				continue
			}
			resp.Existence = append(resp.Existence, v.Name)
			delete(s.ClientNameMapping, v.Name)
			s.startWG(v.JoinServerId)
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
		if v, ok := s.getClientNameMapping(k); ok {
			env.Address = v.Address
			env.Name = v.Name
			env.Join = v.JoinServerId
			idList[k] = env
			continue
		}
		for _, v := range s.ClientNameMapping {
			if v.UserId == k {
				env.Address = v.Address
				env.Name = v.Name
				env.Join = v.JoinServerId
				idList[k] = env
				continue
			}
			if v.UserId[:9] == k {
				env.Address = v.Address
				env.Name = v.Name
				env.Join = v.JoinServerId
				idList[k] = env
			}
		}
	}
	for k, _ := range idList {
		if idList[k].Name != "" {
			if err = os.Remove("/etc/wgctl/wireguard/" + idList[k].Name); err != nil {
				log.Error(err.Error())
				continue
			}
			if err = os.Remove("/etc/wgctl/client/" + idList[k].Name); err != nil {
				log.Error(err.Error())
				continue
			}
			if err = s.deleteClientConfig(fmt.Sprintf("/etc/wireguard/%s", idList[k].Join+".conf"), idList[k].Address); err != nil {
				log.Error(err.Error())
				continue
			}
			resp.Existence = append(resp.Existence, k)
			delete(s.ClientNameMapping, idList[k].Name)
			s.startWG(idList[k].Join)
			continue
		}
		resp.DoesNotExist = append(resp.DoesNotExist, k)
	}
	reply.Delete = resp
	return
}
