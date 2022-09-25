package service

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
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
	PortPool              map[int]bool
	ActiveInterface       map[string]string
	NotActivatedInterface map[string]string
	ClientNameMapping     map[string]*model.ConfigObjConfig
	ServerNameMapping     map[string]*model.ConfigObjConfig
	AddressPool           map[string]map[string]bool
}

func NewService() *Service {
	return &Service{
		ActiveInterface:       make(map[string]string),
		NotActivatedInterface: make(map[string]string),
		ClientNameMapping:     InspectionClientNameMapping(),
		ServerNameMapping:     InspectionServerNameMapping(),
		AddressPool:           InitializeIpAddressPool(),
		PortPool:              InitializePortPool(),
	}
}

func (s *Service) Init() {
	s.InitializeServerConfiguration()
	s.InitializeClientConfiguration()
	s.iptablesCamouflage()
	s.startAllWG()
}

func InspectionClientNameMapping() map[string]*model.ConfigObjConfig {
	var err error
	var clientList []fs.FileInfo
	var configs = &model.ConfigObjConfig{}
	ServerNameMapping := make(map[string]*model.ConfigObjConfig)
	if clientList, err = util.FileForEach("/etc/wgctl/client/"); err != nil {
		log.Error(err.Error())
	}
	if len(clientList) == 0 {
		return ServerNameMapping
	}
	for _, v := range clientList {
		if configs, err = util.ReadConfigs("/etc/wgctl/client/" + v.Name()); err != nil {
			log.Error(err.Error())
			continue
		}
		ServerNameMapping[configs.Name] = configs
	}
	return ServerNameMapping
}

func InspectionServerNameMapping() map[string]*model.ConfigObjConfig {
	var err error
	var serverList []fs.FileInfo
	var configs = &model.ConfigObjConfig{}
	ServerNameMapping := make(map[string]*model.ConfigObjConfig)
	if serverList, err = util.FileForEach("/etc/wgctl/server/"); err != nil {
		log.Error(err.Error())
	}
	if len(serverList) == 0 {
		return nil
	}
	for _, v := range serverList {
		if configs, err = util.ReadConfigs("/etc/wgctl/server/" + v.Name()); err != nil {
			log.Error(err.Error())
			continue
		}
		ServerNameMapping[v.Name()] = configs
	}
	return ServerNameMapping
}

func InitializeIpAddressPool() map[string]map[string]bool {
	Address := make(map[string]map[string]bool)
	configSlice, err := config.WorkConf.UnmarshalKeySliceContainer("wireguard.container")
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	for _, v := range configSlice.Configs {
		ss := strings.FieldsFunc(v.Subnet, util.SplitFunc)
		address := ss[0] + "." + ss[1] + "." + ss[2] + "."
		Address[v.User] = make(map[string]bool)
		for i := 1; i <= 254; i++ {
			Address[v.User][address+strconv.Itoa(i)] = true
		}
		delete(Address[v.User], address+strconv.Itoa(1))
	}
	return Address
}

func InitializePortPool() map[int]bool {
	var min, max int
	portPool := make(map[int]bool)
	rule := strings.Split(config.WorkConf.GetString("wireguard.container.port"), "-")
	min, _ = strconv.Atoi(rule[0])
	max, _ = strconv.Atoi(rule[1])
	for min <= max {
		portPool[min] = true
		min++
	}
	return portPool
}

func (s *Service) InitializeServerConfiguration() {
	var err error
	var configList []fs.FileInfo
	var wgConfigDir []fs.FileInfo
	if configList, err = util.FileForEach(config.WorkConf.GetString("wireguard.wgctlServerDir")); err != nil {
		log.Error(err.Error())
		return
	}
	if len(configList) == 0 {
		var configSlice = &api.Container{}
		list := make(map[string]*model.ConfigObjConfig)
		if configSlice, err = config.WorkConf.UnmarshalKeySliceContainer("wireguard.container"); err != nil {
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
	if wgConfigDir, err = util.FileForEach(config.WorkConf.GetString("wireguard.wgConfigDir")); err != nil {
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

func (s *Service) InitializeClientConfiguration() {
	var err error
	var clientList []fs.FileInfo
	var clientConfigs = &model.ConfigObjConfig{}
	if clientList, err = util.FileForEach(config.WorkConf.GetString("wireguard.wgctlClientDir")); err != nil {
		log.Error(err.Error())
		return
	}
	if len(clientList) == 0 {
		return
	}
	for _, v := range clientList {
		if clientConfigs, err = util.ReadConfigs(config.WorkConf.GetString("wireguard.wgctlClientDir") + v.Name()); err != nil {
			log.Error(err.Error())
			continue
		}
		// 重写最新的服务端配置
		if err = util.SaveJoinServerConfig("/etc/wireguard/"+clientConfigs.JoinServerId+".conf", clientConfigs); err != nil {
			log.Error(err.Error())
		}
		delete(s.AddressPool, clientConfigs.Address)
	}
}

func (s *Service) readCreateServerTemplateConfig(name string) {
	var err error
	var configs = &model.ConfigObjConfig{}
	if configs, err = util.ReadConfigs(config.WorkConf.GetString("wireguard.wgctlServerDir") + name); err != nil {
		log.Error(err.Error())
		return
	}
	s.buildServerConfig(configs)
}

func (s *Service) buildServerConfig(configs *model.ConfigObjConfig) {
	var err error
	var ListenPort int
	if ListenPort, err = s.getListenPort(); err != nil {
		log.Error(err.Error())
		return
	}
	if configs.PrivateKey == "" && configs.PublicKey == "" {
		PrivateKey, PublicKey := util.GenerateKeyPair()
		configs.PrivateKey = PrivateKey
		configs.PublicKey = PublicKey
	}
	if configs.ListenPort == "" {
		configs.ListenPort = strconv.Itoa(ListenPort)
	}
	if configs.Time == 0 {
		configs.Time = int32(time.Now().Unix())
	}
	create := &model.ConfigObjConfig{
		Time:                configs.Time,
		Name:                configs.Name,
		ListenPort:          configs.ListenPort,
		PrivateKey:          configs.PrivateKey,
		PublicKey:           configs.PublicKey,
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
	if err = util.WriteFile(config.WorkConf.GetString("wireguard.wgctlServerDir")+create.Name, string(marshal)); err != nil {
		log.Error(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
		return
	}
	if err = util.WriteFile(config.WorkConf.GetString("wireguard.wgConfigDir")+create.Name+".conf", util.BuildAppendWCS(create)+"\n\n"); err != nil {
		log.Error(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
		return
	}
	s.ServerNameMapping[create.Name] = create
}

func (s *Service) startWG(name string) {
	exec.Command("/bin/sh", "-c", fmt.Sprintf("wg-quick down %s ; wg-quick up %s", name, name)).Run()
}

func (s *Service) startAllWG() {
	var err error
	var list []fs.FileInfo
	if list, err = util.FileForEach(config.WorkConf.GetString("wireguard.wgctlServerDir")); err != nil {
		log.Error(err.Error())
	}
	if len(list) == 0 {
		return
	}
	for _, v := range list {
		exec.Command("/bin/sh", "-c", fmt.Sprintf("wg-quick down %s ; wg-quick up %s", v.Name(), v.Name())).Run()
	}
}

func (s *Service) stopWG(name string) {
	exec.Command("/bin/sh", "-c", fmt.Sprintf("wg-quick down %s", name)).Run()
}

func (s *Service) iptablesCamouflage() {
	output, _ := exec.Command("/bin/sh", "-c", `iptables-save | grep "POSTROUTING -o eth0 -j MASQUERADE"`).CombinedOutput()
	if len(output) == 0 {
		exec.Command("/bin/sh", "-c", "iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE").Run()
	}
}

func (s *Service) getLatestHandshake(ip string) string {
	output, _ := exec.Command("/bin/bash", "-c", fmt.Sprintf("wg | grep -A 2 %s | grep 'latest handshake'", ip)).CombinedOutput()
	ss := strings.Split(string(output), ": ")
	if len(ss) < 2 {
		return "Not connected"
	}
	return strings.Replace(ss[1], "\n", "", -1)
}

func (s *Service) getTransfer(ip string) string {
	output, _ := exec.Command("/bin/bash", "-c", fmt.Sprintf("wg | grep -A 2 %s | grep transfer", ip)).CombinedOutput()
	ss := strings.Split(string(output), ": ")
	if len(ss) < 2 {
		return "Not connected"
	}
	return strings.Replace(ss[1], "\n", "", -1)
}

func (s *Service) getListenPort() (int, error) {
	for k, _ := range s.PortPool {
		delete(s.PortPool, k)
		return k, nil
	}
	return -1, errors.New("no ports available")
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
}

func (s *Service) getClientAddress(name string) (string, error) {
	for k, v := range s.AddressPool {
		if k == name {
			for kk, _ := range v {
				return kk, nil
			}
		}
	}
	return "", errors.New("the server address does not exist")
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

func (s *Service) setIpPool(name, subnet string) string {
	s.AddressPool[name] = make(map[string]bool)
	ss := strings.FieldsFunc(subnet, util.SplitFunc)
	address := ss[0] + "." + ss[1] + "." + ss[2] + "."
	for i := 1; i <= 255; i++ {
		s.AddressPool[name][address+strconv.Itoa(i)] = true
	}
	delete(s.AddressPool[name], address+strconv.Itoa(1))
	return address + strconv.Itoa(1)
}

// SaveServerConfig 保存服务端配置文件
func (s *Service) SaveServerConfig(create *model.ConfigObjConfig) (*model.ConfigObjConfig, error) {
	uc, _ := json.Marshal(create)
	userID := fmt.Sprintf("%x", md5.Sum([]byte(uc)))
	create.UserId = userID
	uc, _ = json.Marshal(create)
	if err := util.WriteFile(config.WorkConf.GetString("wireguard.wgctlServerDir")+create.Name, string(uc)); err != nil {
		log.Error(err.Error())
		return create, errors.New(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
	}
	if err := util.WriteFile(config.WorkConf.GetString("wireguard.wgConfigDir")+create.Name+".conf", util.BuildAppendWCS(create)+"\n\n"); err != nil {
		log.Error(fmt.Sprintf("写入服务端wireguard配置失败%s", err.Error()))
		return nil, err
	}
	return create, nil
}

// deleteClientConfig 将服务端配置中的客户端配置删除
func (s *Service) deleteClientConfig(path, address string) error {
	Output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("grep -n -B 2 %s %s | awk '{print $1}' | grep -o -E '\\<[0-9]\\>|\\<[0-9][0-9]\\>|\\<[0-9][0-9][0-9]\\>'", address, path)).Output()
	if err != nil {
		return err
	}
	nums := strings.Split(string(Output), "\n")
	for i := 0; i < len(nums)-1; i++ {
		exec.Command("/bin/sh", "-c", fmt.Sprintf("sed -i \"%sd\" %s", nums[0], path)).Run()
	}
	return nil
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
