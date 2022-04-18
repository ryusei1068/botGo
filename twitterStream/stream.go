package twitterstream

import (
	"encoding/json"
	"fmt"
	"time"

	httpclient "github.com/botGo/httpClient"
	twitterstream "github.com/fallenstedt/twitter-stream"
	"github.com/fallenstedt/twitter-stream/rules"
	"github.com/fallenstedt/twitter-stream/stream"
)

type (
	TwitterStream struct {
		api        *twitterstream.TwitterApi
		directInfo map[string]Association
	}
	ItwitterStream interface {
		AddRules(string) (*rules.TwitterRuleResponse, error)
		GetRules() (*rules.TwitterRuleResponse, error)
		InitiateStream()
	}
)

type StreamData struct {
	Data struct {
		Text      string    `json:"text"`
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		AuthorID  string    `json:"author_id"`
	} `json:"data"`
	Includes struct {
		Users []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"users"`
	} `json:"includes"`
	MatchingRules []struct {
		ID  string `json:"id"`
		Tag string `json:"tag"`
	} `json:"matching_rules"`
}

func (t *TwitterStream) fetchTweets() stream.IStream {
	api := t.api.Stream

	var err error
	api.SetUnmarshalHook(func(bytes []byte) (interface{}, error) {
		data := StreamData{}

		if err = json.Unmarshal(bytes, &data); err != nil {
			fmt.Printf("failed to unmarshal bytes: %v", err)
		}
		return data, err
	})

	streamExpansions := twitterstream.NewStreamQueryParamsBuilder().
		AddExpansion("author_id").
		AddTweetField("created_at").
		Build()

	err = api.StartStream(streamExpansions)

	if err != nil {
		panic(err)
	}

	return api
}

func (t *TwitterStream) InitiateStream() {
	fmt.Println("Starting Stream")

	api := t.fetchTweets()

	defer t.InitiateStream()

	for tweet := range api.GetMessages() {
		if tweet.Err != nil {
			fmt.Printf("got error from twitter: %v", tweet.Err)
			api.StopStream()
			continue
		}

		result := tweet.Data.(StreamData)
		fmt.Println(result.Data)
	}

	fmt.Println("Stopped Stream")
}

func (t *TwitterStream) AddRules(key string) (*rules.TwitterRuleResponse, error) {
	rules := twitterstream.NewRuleBuilder().AddRule(key, "-is:retweet").Build()

	res, err := t.api.Rules.Create(rules, false)

	if err != nil {
		return nil, err
	}

	if res.Errors != nil && len(res.Errors) > 0 {
		fmt.Printf("Received an error from twitter: %v", res.Errors)
	}

	return res, nil
}

func (t *TwitterStream) GetRules() (*rules.TwitterRuleResponse, error) {
	return t.api.Rules.Get()
}

func (t *TwitterStream) sendAMsgToDiscord() {

}

func (t *TwitterStream) SetDirectInfo(json httpclient.JsonData, keyword string, streamId string) {
	url := webhook + fmt.Sprintf("%s/%s", json.WebHook.Id, json.WebHook.Token)
	DirectInfo[json.WebHook.Id] = Association{word: keyword, url: url, streamId: streamId}
	t.directInfo = DirectInfo
}

func NewTwitterStreamAPI(bearerToken string) ItwitterStream {
	api := twitterstream.NewTwitterStream(bearerToken)
	return &TwitterStream{api: api}
}
