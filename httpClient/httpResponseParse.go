package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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

func (c *httpClient) HandleResponse(resp *http.Response, opts *RequestOpts, fn func(opts *RequestOpts) (*http.Response, error)) (*http.Response, error) {
	if resp.StatusCode == 429 {
		log.Printf("Retrying network request %s with backoff", opts.Url)

		var msg string
		if resp.Body != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			msg = "Network request failed: " + string(body)
		} else {
			msg = "Network request failed with status: " + fmt.Sprint(resp.StatusCode)
		}
		log.Printf(msg)

		var delay time.Duration = 30
		var retrySec int
		retry := resp.Header.Get("retry-after")
		if len(retry) > 0 {
			intVar, err := strconv.Atoi(retry)
			if err != nil {
				retrySec = 30
			} else {
				retrySec = intVar
			}
			delay = time.Duration(retrySec)
		}

		log.Printf("Sleeping for %v seconds", delay)
		time.Sleep(delay)

		return fn(opts)
	}

	if resp.StatusCode >= 400 {
		log.Printf("Network Request at %s failed: %v", opts.Url, resp.StatusCode)

		var msg string
		if resp.Body != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			msg = "Network request failed: " + string(body)
		} else {
			msg = "Network request failed with status" + fmt.Sprint(resp.StatusCode)
		}

		return nil, errors.New(msg)
	}

	return resp, nil
}
