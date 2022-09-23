package cmd

import (
	"context"
	api "ctl/api/grpc/v1"
	http "ctl/api/http/v1"
	"ctl/log"
	req "ctl/service/http"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete 删除指定的wireguard客户端配置 删除多个id或名称使用空格隔开",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		runDelete(ctx, req.NewRequest(), args, &api.DeleteOptions{
			Time:   int32(time.Now().Unix()),
			Server: deletes.Server,
			All:    deletes.All,
		})
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deletes.Server, "server", "s", false, "删除wireguard接口(服务端配置)")
	deleteCmd.Flags().BoolVarP(&deletes.All, "all", "", false, "删除所有wireguard配置")
}

func runDelete(ctx context.Context, req http.Service, args []string, do *api.DeleteOptions) {
	for _, v := range args {
		do.Id = append(do.Id, v)
	}
	response, err := req.Delete(do, "http://127.0.0.1:4000/api/v1/work/delete")
	if err != nil {
		log.Error(err.Error())
		return
	}
	if response.Ret != 1 {
		log.Error(response.Msg)
		return
	}
	marshal, err := json.Marshal(&response.Data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	ss := make(map[string][]string)
	json.Unmarshal(marshal, &ss)
	for _, v := range ss["DoesNotExist"] {
		fmt.Println("Error: No such Name: ", v)
	}
	for _, v := range ss["Existence"] {
		fmt.Println(v)
	}
}
