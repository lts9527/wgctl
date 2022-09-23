package model

type CreateOptions struct {
	Init         bool   `json:"init,omitempty"`
	Time         int32  `json:"time,omitempty"`
	Name         string `json:"user,omitempty"`
	JoinServerId string `json:"join_server_id,omitempty"`
	Subnet       string `json:"subnet,omitempty"`
	ListenPort   string `json:"port,omitempty"`
	DNS          string `json:"dns,omitempty"`
	MTU          string `json:"MTU,omitempty"`
	PublicIp     string `json:"public_ip,omitempty"`
}

type DeleteOptions struct {
	All  bool     `json:"all,omitempty"`
	Time int32    ` json:"time,omitempty"`
	Id   []string ` json:"id,omitempty"`
}

type ShowOptions struct {
	Server bool   `json:"server,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

type PsOptions struct {
	Server     bool   `json:"server,omitempty"`
	Name       string `json:"name,omitempty"`
	WgConfigId string `json:"wg_config_id"`
}

type ConfigObjConfig struct {
	Activation          bool   `json:"activation,omitempty"`
	Time                int32  `json:"time,omitempty"`
	Name                string `json:"name,omitempty"`
	UserId              string `json:"user_id,omitempty"`
	JoinServerId        string `json:"join_server_id,omitempty"`
	WireguardConfig     string `json:"wireguard_config,omitempty"`
	Subnet              string `json:"subnet,omitempty"`
	ListenPort          string `json:"port,omitempty"`
	PrivateKey          string `json:"private_key,omitempty"`
	PublicKey           string `json:"public_key,omitempty"`
	Address             string `json:"address,omitempty"`
	DNS                 string `json:"dns,omitempty"`
	MTU                 string `json:"MTU,omitempty"`
	AllowedIPs          string `json:"allowedIPs,omitempty"`
	Endpoint            string `json:"endpoint,omitempty"`
	PersistentKeepalive string `json:"persistent_keepalive,omitempty"`
}

// Configs is a map of ConfigObjConfig
type Configs map[string]ConfigObjConfig
