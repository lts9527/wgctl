package cmd

import (
	"context"
	api "ctl/api/grpc/v1"
	http "ctl/api/http/v1"
	"ctl/config"
	"ctl/log"
	"ctl/pkg/util"
	req "ctl/service/http"
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init 容器初始化",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		runInit(ctx, req.NewRequest())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	//initCmd.Flags().StringVarP(&Init.Wireguard.InitDir, "config", "c", "/config/config.yaml", "根据配置文件初始化容器")
}

func runInit(ctx context.Context, req http.Service) {
	configSlice, err := ctlConf().UnmarshalKeySliceContainer("wireguard.container")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	CreateFolder()
	Run(&api.InitOptions{Wireguard: &api.Wireguard{
		Container: configSlice,
	}})
}

func CreateFolder() {
	exec.Command("/bin/bash", "-c", "mkdir -p ./config/wgctl/{server,client,wireguard}").Run()
}

func Run(io *api.InitOptions) {
	err := util.WriteFile("docker-compose.yaml", BuildDockerRunWorkCmd(&api.Container{
		Name:   io.Wireguard.Container.Name,
		Port:   io.Wireguard.Container.Port,
		Subnet: io.Wireguard.Container.Subnet,
	}))
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func BuildDockerRunWorkCmd(co *api.Container) string {
	var sb strings.Builder
	sb.Write([]byte(RunDockerfile))
	var envMap = make(map[string]interface{})
	envMap["WorkAddress"] = ctlConf().GetString("server.work.address")
	envMap["ApiGatewayAddress"] = ctlConf().GetString("server.apiGateway.address")
	envMap["ApiGatewayMappingPort"] = ctlConf().GetString("server.apiGateway.port")
	envMap["WorkMappingPort"] = co.Port
	envMap["subnet"] = co.Subnet
	envMap["gateway"] = getGatewayIp(co.Subnet)
	return util.FromTemplateContent(sb.String(), envMap)
}

func getGatewayIp(subnet string) string {
	r := []rune(subnet)
	r[len(subnet)-4] = '1'
	str := strings.Split(string(r), "/")
	return str[0]
}

func getConfig() func() *config.Config {
	f := func() *config.Config {
		return config.NewConfig()
	}
	return f
}

const RunDockerfile = string(`version: '3'
services:
  apiGateway:
    container_name: wgctl-apiGateway
    image: lts9527/api-gateway:test
    restart: unless-stopped
    volumes:
      - ./config/config.yaml:/opt/apiGateway/config/config.yaml
      - /etc/localtime:/etc/localtime
    ports:
      - "{{.ApiGatewayMappingPort}}:{{.ApiGatewayMappingPort}}"
    networks:
      extnetwork:
        ipv4_address: "{{.ApiGatewayAddress}}"
  work:
    container_name: wgctl-work
    image: lts9527/work:test
    privileged: true
    restart: unless-stopped
    volumes:
      - ./config/config.yaml:/opt/work/config/config.yaml
      - ./config/wgctl/:/etc/wgctl/
      - /etc/localtime:/etc/localtime
    ports:
      - "{{.WorkMappingPort}}:{{.WorkMappingPort}}/udp"
    networks:
      extnetwork:
        ipv4_address: "{{.WorkAddress}}"

networks:
  extnetwork:
    name: "wgctl-networks"
    ipam:
      config:
        - subnet: "{{.subnet}}"
          gateway: "{{.gateway}}"`)
