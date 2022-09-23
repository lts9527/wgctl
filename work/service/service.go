package service

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"os/exec"
	"strconv"
	"strings"
	"time"
	api "work/api/grpc/v1"
	"work/config"
	"work/log"
	"work/model"
	"work/pkg/util"
)

type Service struct {
	cache                 *bigcache.BigCache
	ActiveInterface       map[string]string
	NotActivatedInterface map[string]string
	ClientNameMapping     map[string]*model.ConfigObjConfig
	ServerNameMapping     map[string]*model.ConfigObjConfig
	AddressPool           map[string]map[string]interface{}
}

func NewService() *Service {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(1 * time.Minute))
	return &Service{
		cache:                 cache,
		ActiveInterface:       make(map[string]string),
		NotActivatedInterface: make(map[string]string),
		ClientNameMapping:     InspectionClientNameMapping(),
		ServerNameMapping:     InspectionServerNameMapping(),
		AddressPool:           InitializeIpAddressPool(),
	}
}

func (s *Service) getActiveInterface() {
	sl, err := util.FileForEach("/etc/wgctl/wireguard/")
	if err != nil {
		log.Error(err.Error())
	}
	for _, v := range sl {
		configs, err := util.ReadConfigs("/etc/wgctl/wireguard/" + v.Name())
		if err != nil {
			log.Error(err.Error())
			continue
		}
		s.ActiveInterface[v.Name()] = configs.WireguardConfig
	}
	//output, _ := exec.Command("/bin/sh", "-c", "wg | grep \"interface:\" |awk '{print $2}'").Output()
	//s.ActiveInterface
	//return strings.Replace(string(output), "\n", "", -1)
}

func (s *Service) InitializeServerConfiguration() {
	configList, err := util.FileForEach(config.WorkConf.GetString("wireguard.wgctlServerDir"))
	if err != nil {
		log.Error(err.Error())
		return
	}
	if len(configList) == 0 {
		fmt.Println("len(configList) == 0")
		list := make(map[string]*model.ConfigObjConfig)
		configSlice, err := config.WorkConf.UnmarshalKeySliceContainer("wireguard.container")
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, v := range configSlice.Configs {
			list[v.User] = nil
		}
		for _, v := range configSlice.Configs {
			list[v.User] = &model.ConfigObjConfig{
				Name:                v.User,
				Subnet:              v.Subnet,
				Address:             v.Address,
				DNS:                 "8.8.8.8",
				MTU:                 "1350",
				AllowedIPs:          "0.0.0.0/0",
				PersistentKeepalive: "25",
			}
		}
		for k, _ := range list {
			s.buildServerConfig(list[k])
		}
		return
	}
	wgConfigDir, err := util.FileForEach(config.WorkConf.GetString("wireguard.wgConfigDir"))
	if err != nil {
		log.Error(err.Error())
		return
	}
	if len(wgConfigDir) == 0 {
		for _, v := range configList {
			s.readCreateServerTemplateConfig(v.Name())
		}
		return
	}
	createList := make(map[string]bool)
	for _, v := range configList {
		for _, vv := range wgConfigDir {
			if v.Name() == vv.Name() {
				delete(createList, v.Name())
				continue
			}
			createList[v.Name()] = true
		}
	}
	for k, _ := range createList {
		s.readCreateServerTemplateConfig(k)
	}
}

func (s *Service) readCreateServerTemplateConfig(name string) {
	fmt.Println("name", name)
	configs, err := util.ReadConfigs(config.WorkConf.GetString("wireguard.wgctlServerDir") + name)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("configs", configs)
	s.buildServerConfig(configs)
}

func (s *Service) initCreateServerTemplateConfig(confConfigs *api.Configs) {
	//configs := &model.ConfigObjConfig{
	//	Name:                confConfigs.User,
	//	Subnet:              confConfigs.Subnet,
	//	Address:             confConfigs.Address,
	//	DNS:                 "8.8.8.8",
	//	MTU:                 "1350",
	//	AllowedIPs:          "0.0.0.0/0",
	//	PersistentKeepalive: "25",
	//}
	//s.buildServerConfig(configs)
}

func (s *Service) buildServerConfig(configs *model.ConfigObjConfig) {
	var err error
	PrivateKey, PublicKey := util.GenerateKeyPair()
	create := &model.ConfigObjConfig{
		Time:                int32(time.Now().Unix()),
		Name:                configs.Name,
		ListenPort:          strconv.Itoa(s.getListenPort()),
		PrivateKey:          PrivateKey,
		PublicKey:           PublicKey,
		Address:             configs.Address,
		DNS:                 configs.DNS,
		MTU:                 configs.MTU,
		AllowedIPs:          configs.AllowedIPs,
		PersistentKeepalive: configs.PersistentKeepalive,
	}
	marshal, _ := json.Marshal(&create)
	userID := fmt.Sprintf("%x", md5.Sum([]byte(marshal)))
	create.UserId = userID
	marshal, _ = json.Marshal(&create)
	err = util.WriteFile(config.WorkConf.GetString("wireguard.wgctlServerDir")+create.Name, string(marshal))
	if err != nil {
		log.Error(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
		return
	}
	err = util.WriteFile(config.WorkConf.GetString("wireguard.wgConfigDir")+create.Name+".conf", util.BuildAppendWCS(create)+"\n\n")
	if err != nil {
		log.Error(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
		return
	}
	s.ServerNameMapping[create.Name] = create
	s.iptablesCamouflage()
}

func (s *Service) startWG(name string) {
	exec.Command("/bin/sh", "-c", fmt.Sprintf("wg-quick down %s ; wg-quick up %s", name, name)).Run()
}

func (s *Service) stopWG(name string) {
	exec.Command("/bin/sh", "-c", fmt.Sprintf("wg-quick down %s", name)).Run()
}

func (s *Service) iptablesCamouflage() {
	output, err := exec.Command("/bin/sh", "-c", `iptables-save | grep "POSTROUTING -o eth0 -j MASQUERADE"`).Output()
	if err != nil {
		log.Error(err.Error())
	}
	if len(output) == 0 {
		exec.Command("/bin/sh", "-c", "iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE").Run()
	}
}

func (s *Service) getListenPort() int {
	var min, max int
	rule := strings.Split(config.WorkConf.GetString("wireguard.container.port"), "-")
	min, _ = strconv.Atoi(rule[0])
	max, _ = strconv.Atoi(rule[1])
	return util.GenerateRandInt(min, max)
}

func InspectionServerNameMapping() map[string]*model.ConfigObjConfig {
	ServerNameMapping := make(map[string]*model.ConfigObjConfig)
	sl, err := util.FileForEach("/etc/wgctl/server/")
	if err != nil {
		log.Error(err.Error())
	}
	for _, v := range sl {
		configs, err := util.ReadConfigs("/etc/wgctl/server/" + v.Name())
		if err != nil {
			log.Error(err.Error())
			continue
		}
		ServerNameMapping[v.Name()] = configs
	}
	return ServerNameMapping
}

func InspectionClientNameMapping() map[string]*model.ConfigObjConfig {
	ServerNameMapping := make(map[string]*model.ConfigObjConfig)
	sl, err := util.FileForEach("/etc/wgctl/client/")
	if err != nil {
		log.Error(err.Error())
	}
	for _, v := range sl {
		configs, err := util.ReadConfigs("/etc/wgctl/client/" + v.Name())
		if err != nil {
			log.Error(err.Error())
			continue
		}
		ServerNameMapping[configs.Name] = configs
	}
	return ServerNameMapping
}

func (s *Service) getServerNameMapping(name string) (*model.ConfigObjConfig, bool) {
	if k, ok := s.ServerNameMapping[name]; ok {
		return k, true
	}
	return nil, false
}

func (s *Service) getClientNameMapping(name string) (*model.ConfigObjConfig, bool) {
	if k, ok := s.ClientNameMapping[name]; ok {
		return k, true
	}
	return nil, false
}

func (s *Service) getServerNameMappingAll() {
	for k, v := range s.ServerNameMapping {
		fmt.Println("k", k)
		fmt.Println("v", v)
	}
}

func (s *Service) getClientNameMappingAll() {
	for k, _ := range s.ClientNameMapping {
		fmt.Println(k)
	}
}

func (s *Service) getServerListenPort(name string) string {
	return s.ServerNameMapping[name].ListenPort
	//ws, _ := util.ReadConfigs(config.WorkConf.GetString("wireguard.wgctlServerDir") + name)
	//return ws.ListenPort
}

func (s *Service) setIpPool(name, subnet string) string {
	s.AddressPool[name] = make(map[string]interface{})
	ss := strings.FieldsFunc(subnet, util.SplitFunc)
	address := ss[0] + "." + ss[1] + "." + ss[2] + "."
	for i := 1; i <= 255; i++ {
		s.AddressPool[name][address+strconv.Itoa(i)] = nil
	}
	delete(s.AddressPool[name], address+strconv.Itoa(1))
	return address + strconv.Itoa(1)
}

func InitializeIpAddressPool() map[string]map[string]interface{} {
	Address := make(map[string]map[string]interface{})
	configSlice, err := config.WorkConf.UnmarshalKeySliceContainer("wireguard.container")
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	for _, v := range configSlice.Configs {
		ss := strings.FieldsFunc(v.Subnet, util.SplitFunc)
		address := ss[0] + "." + ss[1] + "." + ss[2] + "."
		Address[v.User] = make(map[string]interface{})
		for i := 1; i <= 254; i++ {
			Address[v.User][address+strconv.Itoa(i)] = nil
		}
		delete(Address[v.User], address+strconv.Itoa(1))
	}
	return Address
}

func (s *Service) getAddress(name string) (string, error) {
	for k, v := range s.AddressPool {
		if k == name {
			for kk, _ := range v {
				return kk, nil
			}
		}
	}
	return "", errors.New("the server address does not exist")
}

// SaveServerConfig 保存服务端配置文件
func (s *Service) SaveServerConfig(create *model.ConfigObjConfig) (*model.ConfigObjConfig, error) {
	uc, _ := json.Marshal(create)
	userID := fmt.Sprintf("%x", md5.Sum([]byte(uc)))
	create.UserId = userID
	uc, _ = json.Marshal(create)
	err := util.WriteFile(config.WorkConf.GetString("wireguard.wgctlServerDir")+create.Name, string(uc))
	if err != nil {
		log.Error(err.Error())
		return create, errors.New(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
	}
	err = util.WriteFile(config.WorkConf.GetString("wireguard.wgConfigDir")+create.Name+".conf", util.BuildAppendWCS(create)+"\n\n")
	if err != nil {
		log.Error(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
		return nil, err
	}
	return create, nil
}

func (s *Service) getServerAddress(id string) (string, bool) {
	for k, v := range s.ServerNameMapping {
		if id == k {
			return v.Address, true
		}
	}
	return "", false
}

func (s *Service) getCreateTime(createTime int64) int {
	diff1 := time.Since(time.Unix(createTime, 0))
	str := strconv.FormatFloat(diff1.Seconds(), 'f', 0, 64)
	atoi, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return atoi
}

func (s *Service) formatFileName(name, symbol string) string {
	str := strings.Split(name, symbol)
	return str[0]
}

func (s *Service) formatTimeFormat(atoi int) string {
	switch {
	case atoi > 518400*60:
		return fmt.Sprintf("Create %d year", atoi/(518400*60))
	case atoi > 43200*60:
		return fmt.Sprintf("Create %d month", atoi/(43200*60))
	case atoi > 1440*60:
		return fmt.Sprintf("Create %d days", atoi/(1440*60))
	case atoi > 60*60:
		return fmt.Sprintf("Create %d hours", atoi/(60*60))
	case atoi < 60*60 && atoi > 60:
		return fmt.Sprintf("Create %d minutes", atoi/60)
	case atoi < 60:
		return fmt.Sprintf("Create %d seconds", atoi)
	default:
		return fmt.Sprintf("Create %d minutes", atoi/60)
	}
}
