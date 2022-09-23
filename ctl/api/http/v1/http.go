package v1

import (
	api "ctl/api/grpc/v1"
)

type Service interface {
	Create(req *api.CreateOptions, url string) (*Response, error)
	Show(req *api.ShowOptions, url string) (*Response, error)
	Up(req *api.UpOptions, url string) (*Response, error)
	Ps(req *api.PsOptions, url string) (*Response, error)
	Delete(req *api.DeleteOptions, url string) (*Response, error)
}

type Response struct {
	Ret   uint        `json:"ret,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg,omitempty"`
	Error string      `json:"error,omitempty"`
}
