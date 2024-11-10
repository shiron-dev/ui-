package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	discord.AddHandler(onMessageCreate)

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	stopBot := make(chan os.Signal, 1)

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-stopBot

	err = discord.Close()
	if err != nil {
		panic(err)
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m == nil {
		panic("message is nil")
	}

	u := m.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)

	if m.Author.Bot {
		return
	}

	if uiCheck(m.Content) {
		_, err := s.ChannelMessageSend(m.ChannelID, "うい！")
		if err != nil {
			panic(err)
		}
	}
}

func uiCheck(message string) bool {
	reg, err := regexp.Compile(`^\s*(うい|ui)[\p{P}\p{S}ー]*$`)
	if err != nil {
		panic(err)
	}

	return reg.MatchString(message)
}
