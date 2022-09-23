package cmd

import api "ctl/api/grpc/v1"

var (
	Init   = &api.InitOptions{}
	create = &api.CreateOptions{}
	//up     = &api.UpOptions{}
	show    = &api.ShowOptions{}
	ps      = &api.PsOptions{}
	deletes = &api.DeleteOptions{}
)
