package httpclient

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

var Endpoint = make(map[string]string)

type (
	IHttpClient interface {
		CreateWebhook(string, string) (JsonData, error)
		GetChannelWebhooks(string) (JsonData, error)
		NewHttpRequest(opts *RequestOpts) (*http.Response, error)
	}
	httpClient struct {
		token string
	}
)

// Create Webhook
func (c *httpClient) CreateWebhook(channelId string, hookName string) (JsonData, error) {
	opts := &RequestOpts{
		Method: "POST",
		Url:    Endpoint["base"] + fmt.Sprintf("channels/%s/webhooks", channelId),
		Body:   fmt.Sprintf(`{ "name" : "%s" }`, hookName),
	}

	res, err := c.NewHttpRequest(opts)
	if err != nil {
		return JsonData{}, err
	}

	return c.ParseJson(opts, res)
}

// Get Channel Webhooks
func (c *httpClient) GetChannelWebhooks(channelId string) (JsonData, error) {
	opts := &RequestOpts{
		Method: "GET",
		Url:    Endpoint["base"] + fmt.Sprintf("channels/%s/webhooks", channelId),
	}

	res, err := c.NewHttpRequest(opts)
	if err != nil {
		return JsonData{}, err
	}

	return c.ParseJson(opts, res)
}

// http Request
func (c *httpClient) NewHttpRequest(opts *RequestOpts) (*http.Response, error) {
	var req *http.Request
	var err error

	if opts.Method == "GET" {
		req, err = http.NewRequest(opts.Method, opts.Url, nil)
	} else {
		bufferBody := bytes.NewBuffer([]byte(opts.Body))
		req, err = http.NewRequest(opts.Method, opts.Url, bufferBody)
	}

	req.Header.Set("Content-Type", "application/json")
	if len(opts.Headers) > 0 {
		for _, header := range opts.Headers {
			req.Header.Set(header.Key, header.Value)
		}
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

func NewHttpClient(token string) IHttpClient {
	Endpoint["base"] = "https://discord.com/api/v9/"
	return &httpClient{token: token}
}
