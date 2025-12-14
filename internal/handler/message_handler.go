package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/chun37/doppelcord/internal/repository"
)

type MessageHandler struct {
	userRepo repository.UserRepository
}

func NewMessageHandler(userRepo repository.UserRepository) *MessageHandler {
	return &MessageHandler{userRepo: userRepo}
}

func (h *MessageHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	ctx := context.Background()
	isRegistered, err := h.userRepo.IsRegistered(ctx, m.Author.ID)
	if err != nil {
		log.Printf("Error checking user registration: %v", err)
		return
	}

	if isRegistered {
		fmt.Printf("[登録済] Channel ID: %s, Author ID: %s, Content: %s\n", m.ChannelID, m.Author.ID, m.Content)
	} else {
		fmt.Printf("Channel ID: %s, Author ID: %s, Content: %s\n", m.ChannelID, m.Author.ID, m.Content)
	}
}
