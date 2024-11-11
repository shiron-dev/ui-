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
	"strings"
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

	img := ""
	if str, ok := uiSayCheck(m.Content); ok {
		img = getStickerByMsg(context.Background(), str)
	}

	if img == "" && uiCheck(m.Content) {
		img = randomSticker(context.Background())
	}

	if img != "" {
		_, err := s.ChannelMessageSend(m.ChannelID, img)
		if err != nil {
			panic(err)
		}
	}
}

func uiSayCheck(message string) (string, bool) {
	re := regexp.MustCompile(`^\s*([うぅ憂][いぃ]|憂|(?i)ui)[>＞]`)
	if loc := re.FindStringIndex(message); loc != nil {
		return strings.TrimSpace(message[loc[1]:]), true
	}

	return "", false
}

func uiCheck(message string) bool {
	reg, err := regexp.Compile(`^\s*([うぅ憂][いぃ]|憂|(?i)ui)([\p{P}\p{S}ー]|<:.+:\d+>)*$`)
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
	stickers := getStickers(ctx)

	img := stickers.Stickers[rand.Intn(len(stickers.Stickers))].Img

	imgPath, err := url.JoinPath(baseUrl, img)
	if err != nil {
		panic(err)
	}

	return imgPath
}

func getStickerByMsg(ctx context.Context, mgs string) string {
	stickers := getStickers(ctx)

	for _, sticker := range stickers.Stickers {
		if sticker.Name == mgs {
			imgPath, err := url.JoinPath(baseUrl, sticker.Img)
			if err != nil {
				panic(err)
			}

			return imgPath
		}
	}

	return ""
}

func getStickers(ctx context.Context) Response {
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

	return stickers
}
