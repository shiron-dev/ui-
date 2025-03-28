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
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"golang.org/x/exp/rand"
)

const (
	baseURL      = "https://vt.imgs.shiron.dev/ui_shig/stickers/"
	stickersJSON = "stickers.json"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	//nolint:gosec
	rand.Seed(uint64(time.Now().UnixNano()))

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

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-stopBot

	err = discord.Close()
	if err != nil {
		panic(err)
	}
}

func onMessageCreate(discordSession *discordgo.Session, discordMessage *discordgo.MessageCreate) {
	if discordMessage == nil {
		panic("message is nil")
	}

	u := discordMessage.Author

	//nolint:forbidigo
	fmt.Printf("%20s %20s(%20s) > %s\n", discordMessage.ChannelID, u.Username, u.ID, discordMessage.Content)

	if discordMessage.Author.Bot {
		return
	}

	img := ""
	if str, ok := uiSayCheck(discordMessage.Content); ok {
		img = getStickerByMsg(context.Background(), str)
	}

	if img == "" && uiCheck(discordMessage.Content) {
		img = randomSticker(context.Background())
	}

	if img != "" {
		_, err := discordSession.ChannelMessageSend(discordMessage.ChannelID, img)
		if err != nil {
			panic(err)
		}
	}
}

//nolint:gosmopolitan
const uiSayCheckPattern = `^\s*([うぅ憂][いぃ]|憂|(?i)ui)[>＞]`

func uiSayCheck(message string) (string, bool) {
	reg := regexp.MustCompile(uiSayCheckPattern)
	if loc := reg.FindStringIndex(message); loc != nil {
		return strings.TrimSpace(message[loc[1]:]), true
	}

	return "", false
}

//nolint:gosmopolitan
const uiCheckPattern = `^\s*([うぅ憂][いぃ]|憂|(?i)ui)([\p{P}\p{S}ー]|<:.+:\d+>)*$`

func uiCheck(message string) bool {
	reg := regexp.MustCompile(uiCheckPattern)

	return reg.MatchString(message)
}

type response struct {
	Stickers []stickerResp `json:"stickers"`
}

type stickerResp struct {
	Name string `json:"name"`
	Img  string `json:"img"`
}

func randomSticker(ctx context.Context) string {
	//nolint:gosec
	rand.Seed(uint64(time.Now().UnixNano()))

	stickers := getStickers(ctx)

	img := stickers.Stickers[rand.Intn(len(stickers.Stickers))].Img

	imgPath, err := url.JoinPath(baseURL, img)
	if err != nil {
		panic(err)
	}

	return imgPath
}

func getStickerByMsg(ctx context.Context, mgs string) string {
	stickers := getStickers(ctx)

	for _, sticker := range stickers.Stickers {
		if sticker.Name == mgs {
			imgPath, err := url.JoinPath(baseURL, sticker.Img)
			if err != nil {
				panic(err)
			}

			return imgPath
		}
	}

	return ""
}

func getStickers(ctx context.Context) response {
	jsonPath, err := url.JoinPath(baseURL, stickersJSON)
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

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	var stickers response
	if err = json.NewDecoder(resp.Body).Decode(&stickers); err != nil {
		panic(err)
	}

	return stickers
}
