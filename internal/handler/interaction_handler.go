package handler

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/chun37/doppelcord/internal/repository"
	"github.com/chun37/doppelcord/internal/repository/postgres"
)

type InteractionHandler struct {
	userRepo repository.UserRepository
}

func NewInteractionHandler(userRepo repository.UserRepository) *InteractionHandler {
	return &InteractionHandler{userRepo: userRepo}
}

func (h *InteractionHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ApplicationCommandData().Name == "register" {
		h.handleRegister(s, i)
	}
}

func (h *InteractionHandler) handleRegister(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	userID := i.Member.User.ID

	isRegistered, err := h.userRepo.IsRegistered(ctx, userID)
	if err != nil {
		log.Printf("Error checking user registration: %v", err)
		h.respondWithError(s, i, "エラーが発生しました。")
		return
	}

	if isRegistered {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "既に登録済みです。",
			},
		})
		return
	}

	_, err = h.userRepo.Register(ctx, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserAlreadyExists) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "既に登録済みです。",
				},
			})
			return
		}
		log.Printf("Error registering user: %v", err)
		h.respondWithError(s, i, "登録中にエラーが発生しました。")
		return
	}

	fmt.Printf("ユーザーを登録しました: %s\n", userID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "登録が完了しました！",
		},
	})
}

func (h *InteractionHandler) respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
