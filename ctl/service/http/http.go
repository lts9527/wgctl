package http

import (
	api "ctl/api/grpc/v1"
	http "ctl/api/http/v1"
	"encoding/json"
	"github.com/nahid/gohttp"
	"strconv"
	"time"
)

type Request struct {
}

func NewRequest() http.Service {
	return &Request{}
}

func (r *Request) Create(req *api.CreateOptions, url string) (reply *http.Response, err error) {
	marshal, err := json.Marshal(&api.CreateOptions{
		Time:         req.Time,
		Name:         req.Name,
		JoinServerId: req.JoinServerId,
		Subnet:       req.Subnet,
		ListenPort:   req.ListenPort,
		Dns:          req.Dns,
		Mtu:          req.Mtu,
		PublicIp:     req.PublicIp,
		NewServer:    req.NewServer,
	})
	if err != nil {
		return nil, err
	}
	resp, err := gohttp.NewRequest().Body(marshal).Post(url)
	if err != nil || resp == nil {
		return reply, err
	}
	body, err := resp.GetBodyAsByte()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (r *Request) Show(req *api.ShowOptions, url string) (reply *http.Response, err error) {
	marshal, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := gohttp.NewRequest().Body(marshal).Post(url)
	if err != nil || resp == nil {
		return nil, err
	}
	body, err := resp.GetBodyAsString()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &reply)
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (r *Request) Ps(req *api.PsOptions, url string) (reply *http.Response, err error) {
	marshal, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := gohttp.NewRequest().Body(marshal).Post(url)
	if err != nil || resp == nil {
		return nil, err
	}
	body, err := resp.GetBodyAsString()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (r *Request) Up(req *api.UpOptions, url string) (reply *http.Response, err error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := gohttp.NewRequest().JSON(map[string]interface{}{
		"req":       req,
		"timestamp": timestamp,
	}).Post(url)
	if err != nil || resp == nil {
		return nil, err
	}
	body, err := resp.GetBodyAsString()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (r *Request) Delete(req *api.DeleteOptions, url string) (reply *http.Response, err error) {
	marshal, err := json.Marshal(req)
	resp, err := gohttp.NewRequest().Body(marshal).Post(url)
	if err != nil || resp == nil {
		return nil, err
	}
	body, err := resp.GetBodyAsString()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
