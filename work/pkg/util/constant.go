package util

const WCS string = `[Interface]
PrivateKey = {{.PrivateKey}}
ListenPort = {{.ListenPort}}
Address = {{.Address}}
DNS = {{.DNS}}
MTU = {{.MTU}}`

const WCC string = `[Peer]
PublicKey = {{.PublicKey}}
AllowedIPs = {{.AllowedIPs}}
Endpoint = {{.Endpoint}}
PersistentKeepalive = {{.PersistentKeepalive}}`

const APPENDSERVERCONFIG string = `[Peer]
PublicKey = {{.PublicKey}}
AllowedIPs = {{.AllowedIPs}}`

const APPENDSERVERCONFIGS string = `[Interface]
PrivateKey = {{.PrivateKey}}
ListenPort = {{.ListenPort}}`

const SERVERCONFIGTEMPLATE = string(`{"time":1663248935,"name":"{{.Name}}","port":"{{.Port}}","private_key":"{{.PrivateKey}}","public_key":"{{.PublicKey}}","address":"{{.Address}}","dns":"{{.DNS}}","MTU":"{{.MTU}}","allowedIPs":"{{.AllowedIPs}}","persistent_keepalive":"{{.PersistentKeepalive}}"}`)
