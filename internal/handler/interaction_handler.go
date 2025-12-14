package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/chun37/doppelcord/internal/domain"
	"github.com/chun37/doppelcord/internal/llm"
	"github.com/chun37/doppelcord/internal/repository"
	"github.com/chun37/doppelcord/internal/repository/postgres"
)

type InteractionHandler struct {
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
	llmClient   *llm.Client
}

func NewInteractionHandler(userRepo repository.UserRepository, messageRepo repository.MessageRepository, llmClient *llm.Client) *InteractionHandler {
	return &InteractionHandler{
		userRepo:    userRepo,
		messageRepo: messageRepo,
		llmClient:   llmClient,
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
	maxDiscordLength = 2000
	truncationSuffix = "\n...(切り詰められました)"
	maxMessages      = 100
	maxPromptChars   = 10000

	systemPromptTemplate = `あなたは以下のメッセージ履歴を持つDiscordユーザーになりきってください。

## このユーザーの発言履歴（新しい順）:
%s

## 指示:
- 上記の発言履歴から、このユーザーの文体、口調、言葉遣い、絵文字の使い方、話題の傾向を分析してください
- このユーザーとして自然にメッセージを送信してください
- 履歴にある特徴的な表現や癖があれば再現してください
- 不自然に履歴を引用したり、なりきりであることを示したりしないでください`

	userPrompt = "何か一言メッセージを送ってください。"
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
	userID := i.Member.User.ID
	channelID := i.ChannelID

	// 1. まずチャンネル指定で履歴取得
	messages, err := h.messageRepo.FindByDiscordIDAndChannelID(ctx, userID, channelID, maxMessages)
	if err != nil {
		log.Printf("Error fetching messages by channel: %v", err)
		h.editResponse(s, i, "メッセージ履歴の取得に失敗しました。")
		return
	}

	// 2. チャンネルに履歴がなければ全チャンネルから取得
	if len(messages) == 0 {
		messages, err = h.messageRepo.FindByDiscordID(ctx, userID, maxMessages, nil)
		if err != nil {
			log.Printf("Error fetching messages: %v", err)
			h.editResponse(s, i, "メッセージ履歴の取得に失敗しました。")
			return
		}
	}

	// 3. 履歴が全くない場合のエラー処理
	if len(messages) == 0 {
		h.editResponse(s, i, "あなたのメッセージ履歴がまだ保存されていません。先に /register で登録してからメッセージを送信してください。")
		return
	}

	// 4. systemプロンプトの生成
	systemPrompt := h.buildSystemPrompt(messages)

	// 5. LLM呼び出し
	response, err := h.llmClient.ChatWithSystem(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("Error calling LLM API: %v", err)
		h.editResponse(s, i, "LLM APIの呼び出しに失敗しました。")
		return
	}

	// 6. 2000文字制限の処理
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

func (h *InteractionHandler) buildSystemPrompt(messages []*domain.Message) string {
	var sb strings.Builder
	charCount := 0

	for _, msg := range messages {
		if charCount+len(msg.Content) > maxPromptChars {
			break
		}
		sb.WriteString(msg.Content)
		sb.WriteString("\n---\n")
		charCount += len(msg.Content) + 5
	}

	return fmt.Sprintf(systemPromptTemplate, sb.String())
}
