package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	httpClient "github.com/botGo/httpClient"
	tweetStream "github.com/botGo/twitterStream"
	"github.com/bwmarrin/discordgo"
	"github.com/fallenstedt/twitter-stream/rules"
	"github.com/joho/godotenv"
)

type (
	BotGo struct {
		client        *httpClient.IHttpClient
		twitterStream *tweetStream.ItwitterStream
		session       *discordgo.Session
		message       *discordgo.MessageCreate
		opts          *Option
		json          httpClient.JsonData
	}

	Option struct {
		Keyword string
		Command string
	}

	Bot interface {
		CmdHandle(*discordgo.Session, *discordgo.MessageCreate)
		execute(string, *Option)
		streaming(string, *Option, httpClient.JsonData)
		streamingTweet(string)
		stopStreaming(string)
	}
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func (b *BotGo) findStreamId(rules *rules.TwitterRuleResponse) string {
	for i := range rules.Errors {
		if rules.Errors[i].Value == fmt.Sprintf("(%s -is:retweet -has:mentions -is:reply -is:quote)", b.opts.Keyword) {
			return rules.Errors[i].Id
		}
	}

	for i := range rules.Data {
		if rules.Data[i].Tag == b.opts.Keyword {
			return rules.Data[i].Id
		}
	}

	return ""
}

func (b *BotGo) streamingTweet(channelId string) {
	t := *b.twitterStream
	c := b.client

	var streamId string
	// create new rule of twitter stream
	rules, _ := t.AddRules(b.opts.Keyword)

	streamId = b.findStreamId(rules)
	if len(streamId) > 0 {
		t.SetDirectInfo(b.json, b.opts.Keyword, streamId)
		t.InitiateStream(c)
	} else {
		fmt.Println(b)
		fmt.Println(rules)
		b.session.ChannelMessageSend(channelId, "Could not start streaming!")
	}
}

func (b *BotGo) stopStreaming(channelId string) {

}

func (b *BotGo) streaming(channelId string, opts *Option, json httpClient.JsonData) {
	b.setOptsAndJson(opts, json)

	if opts.Command == "!stream" {
		b.streamingTweet(channelId)
	} else if opts.Command == "!stop" {
		b.stopStreaming(channelId)
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
		fmt.Println(err)
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

func (b *BotGo) setOptsAndJson(opts *Option, json httpClient.JsonData) {
	b.json = json
	b.opts = opts
}

func NewBotGo() Bot {
	httpClient := httpClient.NewHttpClient(GoDotEnvVariable("BOTTOKEN"))
	twitterStream := tweetStream.NewTwitterStreamAPI(GoDotEnvVariable("BEARER_TOKEN"))
	return &BotGo{client: &httpClient, twitterStream: &twitterStream}
}
