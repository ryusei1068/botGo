package twitterstream

import (
	"fmt"

	twitterstream "github.com/fallenstedt/twitter-stream"
	"github.com/fallenstedt/twitter-stream/rules"
)

type (
	TwitterStream struct {
		api *twitterstream.TwitterApi
	}
	ItwitterStream interface {
		StreamStweet()
		AddRules(string) (*rules.TwitterRuleResponse, error)
	}
)

func (api *TwitterStream) StreamStweet() {

}

func (api *TwitterStream) AddRules(key string) (*rules.TwitterRuleResponse, error) {
	rules := twitterstream.NewRuleBuilder().AddRule(key, "-is:retweet").Build()

	res, err := api.api.Rules.Create(rules, false)

	if err != nil {
		return nil, err
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		fmt.Printf("Received an error from twitter: %v", res.Errors)
		return nil, err
	}

	return res, nil
}

func NewTwitterStreamAPI(bearerToken string) ItwitterStream {
	api := twitterstream.NewTwitterStream(bearerToken)
	return &TwitterStream{api: api}
}
