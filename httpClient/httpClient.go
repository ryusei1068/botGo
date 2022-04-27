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
		DeleteWebhookInChannel(string) error
	}
	HttpClient struct {
		token string
	}
)

// Create Webhook
func (c *HttpClient) CreateWebhook(channelId string, hookName string) (JsonData, error) {
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
func (c *HttpClient) GetChannelWebhooks(channelId string) (JsonData, error) {
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

func (c *HttpClient) DeleteWebhookInChannel(webhookid string) error {
	opts := &RequestOpts{
		Method: "DELETE",
		Url:    Endpoint["base"] + fmt.Sprintf("webhooks/%s", webhookid),
	}

	_, err := c.NewHttpRequest(opts)
	if err != nil {
		return err
	}
	return nil
}

// http Request
func (c *HttpClient) NewHttpRequest(opts *RequestOpts) (*http.Response, error) {
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

	h := new(HttpResponseHandler)
	return h.HandleResponse(resp, opts, c.NewHttpRequest)
}

func NewHttpClient(token string) IHttpClient {
	Endpoint["base"] = "https://discord.com/api/v9/"
	return &HttpClient{token: token}
}
