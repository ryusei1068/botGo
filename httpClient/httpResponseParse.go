package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type JsonData struct {
	WebHooks []Webhook
	WebHook  Webhook
}

func (c *httpClient) ParseJson(opts *RequestOpts, resp *http.Response) (JsonData, error) {
	if resp.StatusCode >= 400 {
		log.Printf("Network Request at %s failed: %v", opts.Url, resp.StatusCode)

		var msg string
		if resp.Body != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			msg = "Network request failed: " + string(body)
		} else {
			msg = "Network request failed with status" + fmt.Sprint(resp.StatusCode)
		}

		return JsonData{}, errors.New(msg)
	}

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
