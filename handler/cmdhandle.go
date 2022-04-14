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

const (
	webhook = "https://discord.com/api/webhooks/"
)

type (
	BotGo struct {
		client       *httpClient.IHttpClient
		twitterStrem *tweetStream.ItwitterStream
	}
	Bot interface {
		StartDistribution(string, string)
		StopDistribution(string, string)
		CmdHandle(*discordgo.Session, *discordgo.MessageCreate)
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
	httpClient := *b.client
	webhook, err := httpClient.CreateWebhook(channelId, key)

	if err != nil {
		return
	}

	fmt.Println(webhook)
}

func (b *BotGo) StopDistribution(channelId string, key string) {
	httpClient := *b.client
	webhooks, err := httpClient.GetChannelWebhooks(channelId)

	if err != nil {
		return
	}

	for i := range webhooks {
		fmt.Println(webhooks[i].GuildId)
	}
}

func (b *BotGo) CmdHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgArr := strings.Split(m.Content, " ")
	channelId := m.ChannelID

	if msgArr[0] == "!stream" && len(msgArr[1:]) > 0 {
		b.StartDistribution(channelId, strings.Join(msgArr[1:], " "))
	} else if msgArr[0] == "!stop" && len(msgArr[1:]) > 0 {
		b.StopDistribution(channelId, strings.Join(msgArr[1:], " "))
	}

}

func NewBotGo() Bot {
	httpClient := httpClient.NewHttpClient(GoDotEnvVariable("BOTTOKEN"))
	twitterStreamer := tweetStream.NewTwitterStreamAPI(GoDotEnvVariable("BEARER_TOKEN"))
	return &BotGo{client: &httpClient, twitterStrem: &twitterStreamer}
}
