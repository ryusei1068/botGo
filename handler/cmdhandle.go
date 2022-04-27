package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/botGo/config"
	httpClient "github.com/botGo/httpClient"
	"github.com/botGo/redis"
	tweetStream "github.com/botGo/twitterStream"
	"github.com/bwmarrin/discordgo"
	"github.com/fallenstedt/twitter-stream/rules"
)

const basewebhook = "https://discord.com/api/webhooks/"

type (
	BotGo struct {
		client        *httpClient.IHttpClient
		twitterStream *tweetStream.ItwitterStream
		redisClient   *redis.IRedis
		session       *discordgo.Session
		message       *discordgo.MessageCreate
	}

	Option struct {
		Keyword string
		Command string
	}

	Bot interface {
		CmdHandle(*discordgo.Session, *discordgo.MessageCreate)
		execute(string, *Option)
		streaming(string, *Option, httpClient.JsonData)
		streamingTweet(string, *Option, httpClient.JsonData)
		stopStreaming(string, httpClient.JsonData)
	}
)

func (b *BotGo) findStreamId(rules *rules.TwitterRuleResponse, opts *Option) string {
	for i := range rules.Errors {
		if rules.Errors[i].Value == fmt.Sprintf("(%s -is:retweet -has:mentions -is:reply -is:quote)", opts.Keyword) {
			return rules.Errors[i].Id
		}
	}

	for i := range rules.Data {
		if rules.Data[i].Tag == opts.Keyword {
			return rules.Data[i].Id
		}
	}

	return ""
}

func (b *BotGo) streamingTweet(channelId string, opts *Option, json httpClient.JsonData) {
	t := *b.twitterStream
	r := *b.redisClient
	var streamId string
	// create new rule of twitter stream
	rules, _ := t.AddRules(opts.Keyword)

	streamId = b.findStreamId(rules, opts)
	if len(streamId) > 0 {
		direct := newDirectInfo(json, channelId)
		var info []redis.DirectInfo = []redis.DirectInfo{direct}
		r.SetValues(streamId, info)
		t.InitiateStream()
	} else {
		log.Println(b)
		log.Println(rules)
		b.session.ChannelMessageSend(channelId, "Could not start streaming!")
	}
}

// delete webhook
func (b *BotGo) stopStreaming(channelId string, json httpClient.JsonData) {
	c := *b.client
	for i := range json.WebHooks {
		c.DeleteWebhookInChannel(json.WebHooks[i].Id)
	}
}

func (b *BotGo) streaming(channelId string, opts *Option, json httpClient.JsonData) {
	if opts.Command == "!stream" {
		b.streamingTweet(channelId, opts, json)
	} else if opts.Command == "!stop" {
		b.stopStreaming(channelId, json)
	}
}

func (b *BotGo) execute(channelId string, opts *Option) {
	client := *b.client
	var json httpClient.JsonData
	var err error

	if opts.Command == "!stream" {
		json, err = client.CreateWebhook(channelId, opts.Keyword)
	} else if opts.Command == "!stop" {
		json, err = client.GetChannelWebhooks(channelId)
	}

	if err != nil {
		log.Printf("%s", err)
		b.session.ChannelMessageSend(channelId, "Sorry, failed your request!, "+fmt.Sprint(err))
	} else {
		b.streaming(channelId, opts, json)
	}
}

func (b *BotGo) CmdHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.setSessionAndMsgCreate(s, m)

	msgArr := strings.Split(m.Content, " ")
	channelId := m.ChannelID

	if msgArr[0][0] == '!' && len(msgArr[1:]) > 0 {
		b.execute(channelId, &Option{
			Keyword: strings.Join(msgArr[1:], " "),
			Command: msgArr[0],
		})
	}

}

func (b *BotGo) setSessionAndMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.session = s
	b.message = m
}

func newDirectInfo(json httpClient.JsonData, channelId string) redis.DirectInfo {
	return redis.DirectInfo{
		WebhookUrl:       fmt.Sprintf(basewebhook+"%s/%s", json.WebHook.Id, json.WebHook.Token),
		WebhookId:        json.WebHook.Id,
		WebhookToken:     json.WebHook.Token,
		DiscordChannelId: channelId,
	}
}

func NewBotGo() Bot {
	httpClient := httpClient.NewHttpClient(config.GoDotEnvVariable("BOTTOKEN"))
	twitterStream := tweetStream.NewTwitterStreamAPI(config.GoDotEnvVariable("BEARER_TOKEN"))
	redisClient := redis.NewRedisClient(config.GoDotEnvVariable("ADDR"))
	return &BotGo{client: &httpClient, twitterStream: &twitterStream, redisClient: &redisClient}
}
