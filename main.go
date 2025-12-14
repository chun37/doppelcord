package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 環境変数からDiscord Botトークンを取得
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set in .env file")
	}

	// Discordセッションを作成
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	// メッセージハンドラーを登録
	dg.AddHandler(messageCreate)

	// Intentを設定（メッセージ受信に必要）
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// WebSocket接続を開く
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// グレースフルシャットダウン
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("\nShutting down gracefully...")
}

// messageCreate はメッセージが作成されたときに呼び出される
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ボット自身のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		return
	}

	// // すべてのメッセージをエコー
	// _, err := s.ChannelMessageSend(m.ChannelID, m.Content)
	// if err != nil {
	// 	log.Printf("Error sending message: %v", err)
	// }
	fmt.Printf("Channel ID: %s, Author ID: %s, Content: %s\n", m.ChannelID, m.Author.ID, m.Content) // メッセージ内容をコンソールに出力
}
