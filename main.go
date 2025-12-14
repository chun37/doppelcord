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

// 登録済みユーザーを保存するマップ
var registeredUsers = make(map[string]bool)

// スラッシュコマンドの定義
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "register",
		Description: "ユーザーを登録します",
	},
}

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

	// 環境変数からGuild IDを取得
	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		log.Fatal("GUILD_ID is not set in .env file")
	}

	// Discordセッションを作成
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	// ハンドラーを登録
	dg.AddHandler(messageCreate)
	dg.AddHandler(interactionCreate)

	// Intentを設定（メッセージ受信に必要）
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// WebSocket接続を開く
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer dg.Close()

	// スラッシュコマンドを登録
	for _, cmd := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, guildID, cmd)
		if err != nil {
			log.Printf("Cannot create '%s' command: %v", cmd.Name, err)
		} else {
			fmt.Printf("Command '%s' registered\n", cmd.Name)
		}
	}

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
	// 登録済みユーザーかどうかでプレフィックスを変える
	if registeredUsers[m.Author.ID] {
		fmt.Printf("[登録済] Channel ID: %s, Author ID: %s, Content: %s\n", m.ChannelID, m.Author.ID, m.Content)
	} else {
		fmt.Printf("Channel ID: %s, Author ID: %s, Content: %s\n", m.ChannelID, m.Author.ID, m.Content)
	}
}

// interactionCreate はインタラクション（スラッシュコマンド等）が発生したときに呼び出される
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ApplicationCommandData().Name == "register" {
		userID := i.Member.User.ID

		if registeredUsers[userID] {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "既に登録済みです。",
				},
			})
			return
		}

		registeredUsers[userID] = true
		fmt.Printf("ユーザーを登録しました: %s\n", userID)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "登録が完了しました！",
			},
		})
	}
}
