package twitterstream

import (
	twitterstream "github.com/fallenstedt/twitter-stream"
)

type (
	TwitterStream struct {
		api *twitterstream.TwitterApi
	}
	ItwitterStream interface {
		StreamStweet()
	}
)

func (api *TwitterStream) StreamStweet() {

}

func NewTwitterStreamAPI(bearerToken string) ItwitterStream {
	api := twitterstream.NewTwitterStream(bearerToken)
	return &TwitterStream{api: api}
}
