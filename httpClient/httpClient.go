package httpclient

import (
	"fmt"
	"log"
	"net/http"
)

var Endpoint = make(map[string]string)

type (
	BotClient interface {
		GetWebHook(channelId string) (*http.Response, error)
		NewHttpRequest(opts *RequestOpts) (*http.Response, error)
	}
	httpClient struct {
		token string
	}
)

func (c *httpClient) GetWebHook(channelId string) (*http.Response, error) {
	res, err := c.NewHttpRequest(&RequestOpts{
		Method: "GET",
		Url:    Endpoint["base"] + fmt.Sprintf("channels/%s/webhooks", channelId),
	})

	if err != nil {
		return nil, err
	}

	return res, err
}

func (c *httpClient) NewHttpRequest(opts *RequestOpts) (*http.Response, error) {
	var req *http.Request
	var err error

	if opts.Method == "GET" {
		req, err = http.NewRequest(opts.Method, opts.Url, nil)
	}

	if len(c.token) > 0 {
		req.Header.Set("Authorization", "Bot "+c.token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to perform request for %s: %v", opts.Url, err)
		return nil, err
	}

	return resp, nil
}

func NewHttpClient(token string) BotClient {
	Endpoint["base"] = "https://discord.com/api/v9/"
	return &httpClient{token: token}
}
