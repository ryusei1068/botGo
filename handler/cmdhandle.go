package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	httpClient "github.com/botGo/httpClient"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	webhook = "https://discord.com/api/webhooks/"
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

func CmdHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	channelID := m.ChannelID
	client := httpClient.NewHttpClient(GoDotEnvVariable("BOTTOKEN"))
	resp, err := client.CreateWebhook(channelID, "gobot")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
