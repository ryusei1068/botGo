package httpclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type JsonData struct {
	WebHooks []Webhook
	WebHook  Webhook
}

func (c *httpClient) ParseJson(opts *RequestOpts, resp *http.Response) (JsonData, error) {
	var webhooks []Webhook // GET
	var webhook Webhook    // POST

	data := JsonData{WebHooks: webhooks, WebHook: webhook}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return JsonData{}, err
	}

	var e error
	if opts.Method == "GET" {
		e = json.Unmarshal(body, &data.WebHooks)
	} else {
		e = json.Unmarshal(body, &data.WebHook)
	}

	if e != nil {
		return JsonData{}, e
	}

	return data, nil
}
