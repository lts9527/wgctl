package cmd

import (
	"context"
	api "ctl/api/grpc/v1"
	http "ctl/api/http/v1"
	req "ctl/service/http"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Ps 查看wireguard配置",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		psUp(ctx, req.NewRequest(), ps)
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolVarP(&ps.Server, "server", "s", false, "查看当前的wireguard接口(服务端配置)")
}

func psUp(ctx context.Context, req http.Service, po *api.PsOptions) {
	response, err := req.Ps(&api.PsOptions{
		Server: po.Server,
	}, "http://127.0.0.1:4000/api/v1/work/ps")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if response.Ret != 1 {
		fmt.Println(response.Msg)
		return
	}
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ss := make(map[string][]api.PsOptions)
	err = json.Unmarshal(marshal, &ss)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var ID = "CLIENT ID"
	if po.Server {
		ID = "SERVER ID"
	}
	fmt.Printf("%-15s %-23s %-7s %-20s\n", ID, "STATUS", "PORTS", "NAMES")
	for _, v := range ss["ps"] {
		fmt.Printf("%-15s %-23s %-7s %-20s\n", v.WgConfigId[:9], v.Status, v.Ports, v.Names)
	}
}
