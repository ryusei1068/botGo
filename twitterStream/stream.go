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
		AddRules(string) (*rules.TwitterRuleResponse, error)
		GetRules() (*rules.TwitterRuleResponse, error)
	}
)

func (t *TwitterStream) AddRules(key string) (*rules.TwitterRuleResponse, error) {
	rules := twitterstream.NewRuleBuilder().AddRule(key, "-is:retweet").Build()

	res, err := t.api.Rules.Create(rules, false)

	if err != nil {
		return nil, err
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		fmt.Printf("Received an error from twitter: %v", res.Errors)
		return res, err
	}

	return res, nil
}

func (t *TwitterStream) GetRules() (*rules.TwitterRuleResponse, error) {
	return t.api.Rules.Get()
}

func NewTwitterStreamAPI(bearerToken string) ItwitterStream {
	api := twitterstream.NewTwitterStream(bearerToken)
	return &TwitterStream{api: api}
}
