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
		CreateWebhook(string, string) (*http.Response, error)
		GetChannelWebhooks(string) (*http.Response, error)
		GetGuildWebhooks(string) (*http.Response, error)
		NewHttpRequest(opts *RequestOpts) (*http.Response, error)
	}
	httpClient struct {
		token string
	}
)

// Create Webhook
func (c *httpClient) CreateWebhook(channelId string, hookName string) (*http.Response, error) {
	res, err := c.NewHttpRequest(&RequestOpts{
		Method: "POST",
		Url:    Endpoint["base"] + fmt.Sprintf("channels/%s/webhooks", channelId),
		Body:   fmt.Sprintf(`{ "name" : "%s" }`, hookName),
	})

	if err != nil {
		return nil, err
	}

	return res, err
}

// Get Channel Webhooks
func (c *httpClient) GetChannelWebhooks(channelId string) (*http.Response, error) {
	res, err := c.NewHttpRequest(&RequestOpts{
		Method: "GET",
		Url:    Endpoint["base"] + fmt.Sprintf("channels/%s/webhooks", channelId),
	})

	if err != nil {
		return nil, err
	}

	return res, err
}

// Get Guild Webhooks
func (c *httpClient) GetGuildWebhooks(guildId string) (*http.Response, error) {
	res, err := c.NewHttpRequest(&RequestOpts{
		Method: "GET",
		Url:    Endpoint["base"] + fmt.Sprintf("guilds/%s/webhooks", guildId),
	})

	if err != nil {
		return nil, err
	}

	return res, err
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
