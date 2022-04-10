package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	handler "github.com/botGo/handler"
	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + handler.GoDotEnvVariable("BOTTOKEN"))
	if err != nil {
		fmt.Println("failed run a bot")
	}

	discord.AddHandler(handler.CmdHandle)

	err = discord.Open()

	stopBot := make(chan os.Signal, 1)

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-stopBot

	err = discord.Close()

	return

}
