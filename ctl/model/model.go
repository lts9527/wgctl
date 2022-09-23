package model

type CreateOptions struct {
	NewServer    bool   `json:"new,omitempty"`
	Time         int32  `json:"time,omitempty"`
	Name         string `json:"name,omitempty"`
	JoinServerId string `json:"join_server_id,omitempty"`
	Subnet       string `json:"subnet,omitempty"`
	ListenPort   string `json:"listen_port,omitempty"`
	Dns          string `json:"dns,omitempty"`
	Mtu          string `json:"mtu,omitempty"`
	PublicIp     string `json:"public_ip,omitempty"`
}

type ConfigObjConfig struct {
	User                string `json:"user,omitempty"`
	JoinServerId        string `json:"join_server_id,omitempty"`
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
