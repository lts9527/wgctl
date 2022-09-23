package cmd

import (
	"context"
	api "ctl/api/grpc/v1"
	http "ctl/api/http/v1"
	"ctl/config"
	"ctl/log"
	"ctl/pkg/util"
	req "ctl/service/http"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create 创建新的wireguard客户端配置",
	PreRun: func(cmd *cobra.Command, args []string) {
		if create.NewServer && create.Name == "" {
			fmt.Println("服务端名称不能为空")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		runUp(ctx, req.NewRequest(), create)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVarP(&create.NewServer, "new", "S", false, "创建服务端配置")
	createCmd.Flags().StringVarP(&create.Name, "name", "n", "", "wireguard配置名称")
	createCmd.Flags().StringVarP(&create.JoinServerId, "join", "j", "root", "要加入的的服务端名称，如果为空表示加入默认")
	createCmd.Flags().StringVarP(&create.Subnet, "subnet", "s", "", "创建配置的网段")
	createCmd.Flags().StringVarP(&create.ListenPort, "port", "p", "", "客户端wireguard监听的端口，为空随机分配")
	createCmd.Flags().StringVarP(&create.Dns, "dns", "d", "8.8.8.8", "配置的dns")
	createCmd.Flags().StringVarP(&create.Mtu, "mtu", "m", "1350", "速率")
	createCmd.Flags().StringVarP(&create.PublicIp, "ip", "i", "", "IP地址")
	createCmd.MarkFlagsRequiredTogether("new", "subnet")
}

func runUp(ctx context.Context, req http.Service, co *api.CreateOptions) {
	//createOptions.Apply(project)

	//err := upOptions.apply(project, services)
	//if err != nil {
	//	return err
	//}

	create := &api.CreateOptions{
		NewServer:    co.NewServer,
		Time:         int32(time.Now().Unix()),
		Name:         co.Name,
		JoinServerId: co.JoinServerId,
		Subnet:       co.Subnet,
		ListenPort:   co.ListenPort,
		Dns:          co.Dns,
		Mtu:          co.Mtu,
		PublicIp:     co.PublicIp,
	}
	if create.PublicIp == "" {
		create.PublicIp, _ = util.GetPublicIp()
	}
	response, err := req.Create(create, fmt.Sprintf("http://127.0.0.1:%s/api/v1/work/create", config.CtlConf.GetString("server.apiGateway.port")))
	if err != nil {
		log.Error(err.Error())
		return
	}
	//fmt.Println("response", response)
	if response.Ret != 1 {
		log.Error(response.Msg)
		return
	}
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	ss := make(map[string]string)
	err = json.Unmarshal(marshal, &ss)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(ss["user_id"])

	//response, err := req.Up(&api.UpOptions{
	//	Time:   0,
	//	UserId: "",
	//}, "")
	//if err != nil {
	//	log.Error(err.Error())
	//}
	//fmt.Println("response", response)
}
