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
	TokenSource      oauth2.TokenSource
	AllowInsecureSSL bool
}

type Client struct {
	req *req.Request
}

func New(o Option) (IHttpClient, error) {
	httpClient := req.C()
	if o.AllowInsecureSSL {
		httpClient.EnableInsecureSkipVerify()
	}
	request := httpClient.R().SetHeader("Content-Type", "application/json")
	if o.TokenSource != nil {
		// injecting authorization header
		httpClient.WrapRoundTripFunc(func(rt req.RoundTripper) req.RoundTripFunc {
			return func(req *req.Request) (resp *req.Response, err error) {
				token, err := o.TokenSource.Token()
				if err != nil {
					return
				}
				req.Headers.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
				resp, err = rt.RoundTrip(req)
				return
			}
		})
	} else if o.BasicAuth.Username != "" && o.BasicAuth.Password != "" {
		request.SetBasicAuth(o.BasicAuth.Username, o.BasicAuth.Password)
	} else {
		return nil, fmt.Errorf("you must set oauth token or basic auth params (username & password)")
	}

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
