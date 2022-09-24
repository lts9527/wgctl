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
