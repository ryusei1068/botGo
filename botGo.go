package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	handler "github.com/botGo/handler"
	twitterstream "github.com/botGo/twitterStream"
	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + twitterstream.GoDotEnvVariable("BOTTOKEN"))

	if err != nil {
		fmt.Println("failed run a bot,", err)
		return
	}

	botGO := handler.NewBotGo()
	discord.AddHandler(botGO.CmdHandle)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	stopBot := make(chan os.Signal, 1)

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-stopBot

	discord.Close()
}
