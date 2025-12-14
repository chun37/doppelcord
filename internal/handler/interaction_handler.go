package handler

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/chun37/doppelcord/internal/llm"
	"github.com/chun37/doppelcord/internal/repository"
	"github.com/chun37/doppelcord/internal/repository/postgres"
)

type InteractionHandler struct {
	userRepo  repository.UserRepository
	llmClient *llm.Client
}

func NewInteractionHandler(userRepo repository.UserRepository, llmClient *llm.Client) *InteractionHandler {
	return &InteractionHandler{
		userRepo:  userRepo,
		llmClient: llmClient,
	}
}

func (h *InteractionHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "register":
		h.handleRegister(s, i)
	case "test":
		h.handleTest(s, i)
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

const (
	testPrompt       = "こんにちは、簡単に自己紹介してください。"
	maxDiscordLength = 2000
	truncationSuffix = "\n...(切り詰められました)"
)

func (h *InteractionHandler) handleTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 遅延応答（LLM呼び出しは時間がかかるため）
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error deferring response: %v", err)
		return
	}

	ctx := context.Background()
	response, err := h.llmClient.Chat(ctx, testPrompt)
	if err != nil {
		log.Printf("Error calling LLM API: %v", err)
		h.editResponse(s, i, "LLM APIの呼び出しに失敗しました。")
		return
	}

	// 2000文字制限の処理
	if len(response) > maxDiscordLength-len(truncationSuffix) {
		response = response[:maxDiscordLength-len(truncationSuffix)] + truncationSuffix
	}

	h.editResponse(s, i, response)
}

func (h *InteractionHandler) editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		log.Printf("Error editing response: %v", err)
	}
}
