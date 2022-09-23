package cmd

import (
	"context"
	api "ctl/api/grpc/v1"
	http "ctl/api/http/v1"
	req "ctl/service/http"
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show 查看指定名称的wireguard配置",
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		showUp(ctx, req.NewRequest(), &api.ShowOptions{
			Time:    0,
			UserId:  args[0],
			Picture: false,
			Server:  show.Server,
		})
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVarP(&show.Server, "server", "s", false, "查看服务端接口状态")
	showCmd.Flags().BoolVarP(&show.Picture, "picture", "p", false, "wireguard配置以二维码展示")
}

func showUp(ctx context.Context, req http.Service, so *api.ShowOptions) {
	response, err := req.Show(so, "http://127.0.0.1:4000/api/v1/task/show")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if response.Ret != 1 {
		fmt.Println(response.Msg)
		return
	}
	if show.Picture {
		q, _ := qrcode.New(response.Data.(string), qrcode.Low)
		fmt.Println(q.ToSmallString(true))
		return
	}
	fmt.Println(response.Data)
	return
}
