package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	httpClient "github.com/botGo/httpClient"
	httpclient "github.com/botGo/httpClient"
	tweetStream "github.com/botGo/twitterStream"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// WebHookUrl and Key word
var WebHookUrl = make(map[string]string)

const (
	webhook = "https://discord.com/api/webhooks/"
)

type Option struct {
	Keyword   string
	Operation string
}

type (
	BotGo struct {
		client       *httpClient.IHttpClient
		twitterStrem *tweetStream.ItwitterStream
		session      *discordgo.Session
		message      *discordgo.MessageCreate
	}
	Bot interface {
		StartDistribution(string, string)
		StopDistribution(string, string)
		CmdHandle(*discordgo.Session, *discordgo.MessageCreate)
		Execute(string, *Option)
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

func (b *BotGo) StartDistribution(channelId string, key string) {

	twitterClient := *b.twitterStrem
	res, _ := twitterClient.AddRules(key)

	fmt.Println(res)
}

func (b *BotGo) StopDistribution(channelId string, key string) {

	twitterClient := *b.twitterStrem
	res, _ := twitterClient.GetRules()

	fmt.Println(res)
}

func (b *BotGo) Execute(channelId string, opts *Option) {
	httpClient := *b.client
	var json httpclient.JsonData
	var err error

	if opts.Operation == "!stream" {
		json, err = httpClient.CreateWebhook(channelId, opts.Keyword)
	} else if opts.Operation == "!stop" {
		json, err = httpClient.GetChannelWebhooks(channelId)
	}

	if err != nil {
		fmt.Println(err)
		b.session.ChannelMessageSend(channelId, "Sorry, failed your request!, "+fmt.Sprint(err))
		return
	}

	fmt.Println(json)
}

func (b *BotGo) CmdHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.setSessionAndMsgCreate(s, m)

	msgArr := strings.Split(m.Content, " ")
	channelId := m.ChannelID

	if msgArr[0][0] == '!' && len(msgArr[1:]) > 0 {
		b.Execute(channelId, &Option{
			Keyword:   strings.Join(msgArr[1:], " "),
			Operation: msgArr[0],
		})
	}

}

func (b *BotGo) setSessionAndMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.session = s
	b.message = m
}

func NewBotGo() Bot {
	httpClient := httpClient.NewHttpClient(GoDotEnvVariable("BOTTOKEN"))
	twitterStreamer := tweetStream.NewTwitterStreamAPI(GoDotEnvVariable("BEARER_TOKEN"))
	return &BotGo{client: &httpClient, twitterStrem: &twitterStreamer}
}
