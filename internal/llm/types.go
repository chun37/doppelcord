package llm

// ChatMessage はOpenAI Chat APIのメッセージ形式
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest はOpenAI Chat APIリクエスト形式
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse はOpenAI Chat APIレスポンス形式
type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error *APIError `json:"error,omitempty"`
}

// APIError はAPIエラーレスポンス
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}
