package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	idList := make(map[string]*deleteEnvMap)
	if do.All {
		for _, v := range s.ClientNameMapping {
			err = os.Remove("/etc/wgctl/wireguard/" + v.Name)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			err = os.Remove("/etc/wgctl/client/" + v.Name)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			err := DeleteClientConfig(fmt.Sprintf("/etc/wireguard/%s", v.JoinServerId+".conf"), v.Address)
			if err != nil {
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
			err = os.Remove("/etc/wgctl/wireguard/" + idList[k].Name)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			err = os.Remove("/etc/wgctl/client/" + idList[k].Name)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			err := DeleteClientConfig(fmt.Sprintf("/etc/wireguard/%s", idList[k].Join+".conf"), idList[k].Address)
			if err != nil {
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
	s.getClientNameMappingAll()
	reply.Delete = resp
	return
}

// DeleteClientConfig 将服务端配置中的客户端配置删除
func DeleteClientConfig(path, address string) error {
	Output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("grep -n -B 2 %s %s | awk '{print $1}' | grep -o -E '\\<[0-9]\\>|\\<[0-9][0-9]\\>|\\<[0-9][0-9][0-9]\\>'", address, path)).Output()
	if err != nil {
		return err
	}
	nums := strings.Split(string(Output), "\n")
	for i := 0; i < len(nums)-1; i++ {
		err = exec.Command("/bin/sh", "-c", fmt.Sprintf("sed -i \"%sd\" %s", nums[0], path)).Run()
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
	return nil
}
