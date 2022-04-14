package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var Endpoint = make(map[string]string)

type (
	IHttpClient interface {
		CreateWebhook(string, string) (Webhook, error)
		GetChannelWebhooks(string) ([]Webhook, error)
		NewHttpRequest(opts *RequestOpts) (*http.Response, error)
	}
	httpClient struct {
		token string
	}
)

// Create Webhook
func (c *httpClient) CreateWebhook(channelId string, hookName string) (Webhook, error) {
	res, err := c.NewHttpRequest(&RequestOpts{
		Method: "POST",
		Url:    Endpoint["base"] + fmt.Sprintf("channels/%s/webhooks", channelId),
		Body:   fmt.Sprintf(`{ "name" : "%s" }`, hookName),
	})

	if err != nil {
		return Webhook{}, err
	}
	var webhook Webhook
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &webhook)

	return webhook, err
}

// Get Channel Webhooks
func (c *httpClient) GetChannelWebhooks(channelId string) ([]Webhook, error) {
	res, err := c.NewHttpRequest(&RequestOpts{
		Method: "GET",
		Url:    Endpoint["base"] + fmt.Sprintf("channels/%s/webhooks", channelId),
	})

	if err != nil {
		return nil, err
	}

	var webhooks []Webhook
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &webhooks)

	return webhooks, err
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
