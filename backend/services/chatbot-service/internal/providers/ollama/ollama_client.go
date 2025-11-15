package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fintrack/chatbot-service/internal/core/ports"
)

type Client struct {
	host    string
	model   string
	http    *http.Client
	timeout time.Duration
}

func NewClient(host, model string) *Client {
	return &Client{
		host:    host,
		model:   model,
		timeout: 30 * time.Second,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Chat llama a Ollama /api/chat optimizada para modelos livianos
func (c *Client) Chat(ctx context.Context, systemPrompt string, userPrompt string) (string, error) {
	// Ollama chat API optimizada
	// POST {host}/api/chat { model, messages, stream:false, options }
	payload := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"stream": false,
		"options": map[string]any{
			"temperature":    0.5,  // Reducido para más velocidad
			"top_p":          0.8,  // Menos diversidad, más velocidad
			"num_ctx":        512,  // Contexto muy reducido
			"num_predict":    128,  // Respuestas muy cortas
			"repeat_penalty": 1.05, // Menor penalización
			"top_k":          20,   // Opciones muy limitadas
		},
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", c.host), bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var out struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Response string `json:"response"`
		Error    string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if out.Error != "" {
		return "", fmt.Errorf("ollama error: %s", out.Error)
	}

	result := out.Message.Content
	if result == "" {
		result = out.Response
	}

	return result, nil
}

var _ ports.LLMProvider = (*Client)(nil)
