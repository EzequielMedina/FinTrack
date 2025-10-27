package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GroqClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	Model       string        `json:"model"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
	Stream      bool          `json:"stream"`
}

type ChatChoice struct {
	Index   int         `json:"index"`
	Message ChatMessage `json:"message"`
	Reason  string      `json:"finish_reason"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewGroqClient(apiKey, model string) *GroqClient {
	if model == "" {
		model = "llama-3.1-8b-instant" // Modelo por defecto (ultra rápido)
	}

	return &GroqClient{
		apiKey:  apiKey,
		baseURL: "https://api.groq.com/openai/v1",
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *GroqClient) Chat(ctx context.Context, system, user string) (string, error) {
	// Preparar request
	chatReq := ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Model:       g.model,
		Temperature: 0.3, // Respuestas más consistentes para finanzas
		MaxTokens:   1000,
		Stream:      false,
	}

	// Serializar JSON
	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	// Crear HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	// Ejecutar request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Extraer contenido
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	content := chatResp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("empty response content")
	}

	return content, nil
}

// Método para verificar conectividad
func (g *GroqClient) HealthCheck(ctx context.Context) error {
	_, err := g.Chat(ctx, "You are a helpful assistant.", "Say 'OK' if you can respond.")
	return err
}
