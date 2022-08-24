package http

import (
	"fmt"
	"sync"

	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
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
	WorkerCount      int
}

type Client struct {
	reqIndex      int
	reqIndexMutex sync.Mutex
	reqs          []*req.Request
}

func New(o Option) (IHttpClient, error) {
	workerCount := o.WorkerCount
	if workerCount <= 0 {
		workerCount = 1
	}
	client := &Client{reqs: make([]*req.Request, workerCount), reqIndex: 0, reqIndexMutex: sync.Mutex{}}
	for i := 0; i < workerCount; i++ {
		log.Debug().Int("worker_index", i).Msg("initialize http worker")
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
		client.reqs[i] = request
	}

	return client, nil
}

// request get current worker
func (h *Client) request() *req.Request {
	h.reqIndexMutex.Lock()
	defer h.reqIndexMutex.Unlock()
	req := h.reqs[h.reqIndex]
	if h.reqIndex+1 == len(h.reqs) {
		h.reqIndex = 0
	} else {
		h.reqIndex++
	}
	return req
}

func (h *Client) GetMarshalReturnObj(url string, obj any) error {
	response, err := h.request().Get(url)
	if err != nil {
		return err
	}

	if err = response.UnmarshalJson(obj); err != nil {
		return err
	}

	return nil
}

func (h *Client) DeleteMarshalReturnObj(url string, obj any) error {
	response, err := h.request().Delete(url)
	if err != nil {
		return err
	}

	if err = response.UnmarshalJson(obj); err != nil {
		return err
	}

	return nil
}
