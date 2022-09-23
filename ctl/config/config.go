package config

import (
	api "ctl/api/grpc/v1"
	"fmt"
	"github.com/spf13/viper"
)

var (
	CtlConf *config
)

type config struct {
	viper           *viper.Viper
	ContainerConfig *api.Container
}

type configInfo struct {
	ConfName string
	ConfType string
	ConfPath string
}

func init() {
	c1 := &configInfo{
		ConfName: "config",
		ConfType: "yaml",
		ConfPath: "config/",
	}
	CtlConf = &config{
		viper: getConf(c1),
	}
}

func getConf(c1 *configInfo) *viper.Viper {
	v := viper.New()
	v.SetConfigName(c1.ConfName)
	v.SetConfigType(c1.ConfType)
	v.AddConfigPath(c1.ConfPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return v
}

func (c *config) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *config) UnmarshalKeySliceContainer(key string) (*api.Container, error) {
	err := c.viper.UnmarshalKey(key, &c.ContainerConfig)
	if err != nil {
		return nil, err
	}
	return c.ContainerConfig, nil
}

func (c *config) GetInt(key string) int {
	return c.viper.GetInt(key)
}
