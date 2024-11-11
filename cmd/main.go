package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"golang.org/x/exp/rand"
)

const (
	baseUrl      = "https://vt.imgs.shiron.dev/ui_shig/stickers/"
	stickersJson = "stickers.json"
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
		_, err := s.ChannelMessageSend(m.ChannelID, randomSticker(context.Background()))

		if err != nil {
			panic(err)
		}
	}
}

func uiCheck(message string) bool {
	reg, err := regexp.Compile(`^\s*(うい|うぃ|ぅい|ぅぃ|(?i)ui)[\p{P}\p{S}ー]*$`)
	if err != nil {
		panic(err)
	}

	return reg.MatchString(message)
}

type Response struct {
	Stickers []StickerResp `json:"stickers"`
}

type StickerResp struct {
	Name string `json:"name"`
	Img  string `json:"img"`
}

func randomSticker(ctx context.Context) string {
	jsonPath, err := url.JoinPath(baseUrl, stickersJson)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jsonPath, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var stickers Response
	if err = json.NewDecoder(resp.Body).Decode(&stickers); err != nil {
		panic(err)
	}

	img := stickers.Stickers[rand.Intn(len(stickers.Stickers))].Img

	imgPath, err := url.JoinPath(baseUrl, img)
	if err != nil {
		panic(err)
	}

	return imgPath
}
