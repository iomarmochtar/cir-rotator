package http

import (
	"fmt"

	"github.com/imroc/req/v3"
	"golang.org/x/oauth2"
)

//go:generate mockgen -destination mock_http/mock_http.go -source http.go IHttpClient
type IHttpClient interface {
	GetMarshalReturnObj(url string, obj any) error
	DeleteMarshalReturnObj(url string, obj any) error
}

type Option struct {
	BasicAuth struct {
		Username string
		Password string
	}
	Token *oauth2.Token
}

type Client struct {
	req *req.Request
}

func New(o Option) (IHttpClient, error) {
	request := req.C().R()
	if o.Token != nil {
		request.SetHeader("Authorization", fmt.Sprintf("%s %s", o.Token.TokenType, o.Token.AccessToken))
	} else if o.BasicAuth.Username != "" && o.BasicAuth.Password != "" {
		request.SetBasicAuth(o.BasicAuth.Username, o.BasicAuth.Password)
	} else {
		return nil, fmt.Errorf("you must set oauth token or basic auth params (username & password)")
	}

	request.Headers.Set("Content-Type", "application/json")
	return &Client{req: request}, nil
}

func (h Client) GetMarshalReturnObj(url string, obj any) error {
	response, err := h.req.Get(url)
	if err != nil {
		return err
	}

	if err = response.UnmarshalJson(obj); err != nil {
		return err
	}

	return nil
}

func (h Client) DeleteMarshalReturnObj(url string, obj any) error {
	response, err := h.req.Delete(url)
	if err != nil {
		return err
	}

	if err = response.UnmarshalJson(obj); err != nil {
		return err
	}

	return nil
}
