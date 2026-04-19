package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not load")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	token := os.Getenv("TOKEN")

	if token == "" {
		fmt.Println("cannot load token")
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)

}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	switch text {
	case "/start":
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Bot started",
		})

	default:
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Text saved " + text,
		})
		fileUUID := uuid.New()
		err := os.MkdirAll("../unsorted", 0755)
		path := fmt.Sprintf("../unsorted/%s.md", &fileUUID)
		file, err := os.Create(path)
		if err != nil {
			fmt.Println("unable to create file", err)
			os.Exit(1)
		}
		defer file.Close()
		file.WriteString(text)
		fmt.Println("writen")

		cmd := exec.Command("bash", "-lc", `cd ../unsorted && git add . && git commit -m "feat: add note" && git push`)
		cmd.Run()
	}
}

// git submodule add <url-repo-b> storage/unsorted
