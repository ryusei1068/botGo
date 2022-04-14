package handler

import (
	"fmt"
	"io/ioutil"
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

type BotGo struct {
	client       *httpClient.IHttpClient
	twitterStrem *tweetStream.ItwitterStream
}

type Organizer interface {
	StartDistribution(string, string)
	StopDistribution(string)
	CmdHandle(*discordgo.Session, *discordgo.MessageCreate)
}

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
	resp, err := httpClient.CreateWebhook(channelId, key)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}

func (b *BotGo) StopDistribution(key string) {

}

func (b *BotGo) CmdHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgArr := strings.Split(m.Content, " ")

	if msgArr[0] == "!stream" {
		fmt.Println(m.Content)
	}

}

func NewBotGo() Organizer {
	httpClient := httpClient.NewHttpClient(GoDotEnvVariable("BOTTOKEN"))
	twitterStreamer := tweetStream.NewTwitterStreamAPI(GoDotEnvVariable("BEARER_TOKEN"))
	return &BotGo{client: &httpClient, twitterStrem: &twitterStreamer}
}
