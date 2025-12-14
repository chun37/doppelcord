package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/chun37/doppelcord/internal/domain"
	"github.com/chun37/doppelcord/internal/repository"
)

type MessageHandler struct {
	userRepo repository.UserRepository
	msgRepo  repository.MessageRepository
}

func NewMessageHandler(userRepo repository.UserRepository, msgRepo repository.MessageRepository) *MessageHandler {
	return &MessageHandler{userRepo: userRepo, msgRepo: msgRepo}
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

		msg := &domain.Message{
			DiscordID: m.Author.ID,
			ChannelID: m.ChannelID,
			MessageID: m.ID,
			Content:   m.Content,
			CreatedAt: m.Timestamp,
		}
		if err := h.msgRepo.Save(ctx, msg); err != nil {
			log.Printf("Error saving message: %v", err)
		}
	} else {
		fmt.Printf("Channel ID: %s, Author ID: %s, Content: %s\n", m.ChannelID, m.Author.ID, m.Content)
	}
}
