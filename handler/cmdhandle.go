package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	httpClient "github.com/botGo/httpClient"
	tweetStream "github.com/botGo/twitterStream"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// WebHookId and Key word
var WebHookUrl = make(map[string]string)

const (
	webhook = "https://discord.com/api/webhooks/"
)

type Option struct {
	Keyword string
	Command string
}

type (
	BotGo struct {
		client        *httpClient.IHttpClient
		twitterStream *tweetStream.ItwitterStream
		session       *discordgo.Session
		message       *discordgo.MessageCreate
		opts          *Option
		json          httpClient.JsonData
	}
	Bot interface {
		CmdHandle(*discordgo.Session, *discordgo.MessageCreate)
		execute(string, *Option)
		streaming(string, *Option, httpClient.JsonData)
		streamingTweet()
		stopStreaming()
	}
)

// command list
// !stream [key word]
// !stop [key word]

func GoDotEnvVariable(key string) string {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func (b *BotGo) streamingTweet() {
	url := webhook + fmt.Sprintf("%s/%s", b.json.WebHook.Id, b.json.WebHook.Token)
	fmt.Println(url)
}

func (b *BotGo) stopStreaming() {
	url := webhook + fmt.Sprintf("%s/%s", b.json.WebHooks[0].Id, b.json.WebHooks[0].Token)
	fmt.Println(url)
}

func (b *BotGo) streaming(channelId string, opts *Option, json httpClient.JsonData) {
	b.setOptsAndJson(opts, json)

	if opts.Command == "!stream" {
		b.streamingTweet()
	} else if opts.Command == "!stop" {
		b.stopStreaming()
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
