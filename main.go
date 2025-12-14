package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/chun37/doppelcord/internal/database"
	"github.com/chun37/doppelcord/internal/handler"
	"github.com/chun37/doppelcord/internal/repository/cached"
	"github.com/chun37/doppelcord/internal/repository/postgres"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "register",
		Description: "ユーザーを登録します",
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set in .env file")
	}

	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		log.Fatal("GUILD_ID is not set in .env file")
	}

	ctx := context.Background()

	dbConfig := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	pool, err := database.NewPostgresPool(ctx, dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	fmt.Println("Connected to database")

	pgUserRepo := postgres.NewUserRepository(pool)
	userRepo := cached.NewCachedUserRepository(pgUserRepo)

	if err := userRepo.LoadAll(ctx); err != nil {
		log.Fatal("Failed to load users into cache:", err)
	}
	fmt.Println("Loaded users into cache")

	msgRepo := postgres.NewMessageRepository(pool)

	msgHandler := handler.NewMessageHandler(userRepo, msgRepo)
	interactionHandler := handler.NewInteractionHandler(userRepo)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(msgHandler.Handle)
	dg.AddHandler(interactionHandler.Handle)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer dg.Close()

	for _, cmd := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, guildID, cmd)
		if err != nil {
			log.Printf("Cannot create '%s' command: %v", cmd.Name, err)
		} else {
			fmt.Printf("Command '%s' registered\n", cmd.Name)
		}
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("\nShutting down gracefully...")
}
