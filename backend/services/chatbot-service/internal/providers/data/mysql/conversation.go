package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fintrack/chatbot-service/internal/core/ports"
	"github.com/google/uuid"
)

// SaveConversationMessage stores a conversation message in the database
func (p *DataProvider) SaveConversationMessage(ctx context.Context, msg ports.ConversationMessage) error {
	// Generate ID if not provided
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}

	// Serialize context data to JSON
	var contextDataJSON []byte
	var err error
	if msg.ContextData != nil {
		contextDataJSON, err = json.Marshal(msg.ContextData)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO conversation_history 
		(id, user_id, conversation_id, role, message, context_data, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	createdAt := msg.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	_, err = p.db.ExecContext(ctx, query,
		msg.ID,
		msg.UserID,
		msg.ConversationID,
		msg.Role,
		msg.Message,
		contextDataJSON,
		createdAt,
	)

	return err
}

// GetConversationHistory retrieves conversation messages for a specific conversation
func (p *DataProvider) GetConversationHistory(ctx context.Context, userID, conversationID string, limit int) ([]ports.ConversationMessage, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `
		SELECT id, user_id, conversation_id, role, message, context_data, created_at
		FROM conversation_history
		WHERE user_id = ? AND conversation_id = ?
		ORDER BY created_at ASC
		LIMIT ?
	`

	rows, err := p.db.QueryContext(ctx, query, userID, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ports.ConversationMessage
	for rows.Next() {
		var msg ports.ConversationMessage
		var contextDataJSON sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.UserID,
			&msg.ConversationID,
			&msg.Role,
			&msg.Message,
			&contextDataJSON,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Deserialize context data if present
		if contextDataJSON.Valid && contextDataJSON.String != "" {
			var contextData map[string]any
			if err := json.Unmarshal([]byte(contextDataJSON.String), &contextData); err == nil {
				msg.ContextData = contextData
			}
		}

		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// GetUserConversations retrieves all conversation IDs for a user (for listing conversations)
func (p *DataProvider) GetUserConversations(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT DISTINCT conversation_id
		FROM conversation_history
		WHERE user_id = ?
		ORDER BY MAX(created_at) DESC
	`

	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []string
	for rows.Next() {
		var convID string
		if err := rows.Scan(&convID); err != nil {
			return nil, err
		}
		conversations = append(conversations, convID)
	}

	return conversations, rows.Err()
}

// GetLastConversationContext retrieves the last inferred context from a conversation
// This is useful for maintaining context across messages
func (p *DataProvider) GetLastConversationContext(ctx context.Context, userID, conversationID string) (*ports.InferredContext, error) {
	query := `
		SELECT context_data
		FROM conversation_history
		WHERE user_id = ? AND conversation_id = ? AND role = 'assistant'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var contextDataJSON sql.NullString
	err := p.db.QueryRowContext(ctx, query, userID, conversationID).Scan(&contextDataJSON)

	if err == sql.ErrNoRows {
		return nil, nil // No previous context
	}
	if err != nil {
		return nil, err
	}

	if !contextDataJSON.Valid || contextDataJSON.String == "" {
		return nil, nil
	}

	var contextData map[string]any
	if err := json.Unmarshal([]byte(contextDataJSON.String), &contextData); err != nil {
		return nil, err
	}

	// Extract inferred context from context_data if it exists
	// This is a best-effort extraction, might need adjustment based on actual data structure
	inferredCtx := &ports.InferredContext{
		ContextFocus: "general",
		PeriodLabel:  "this month",
	}

	if focus, ok := contextData["inferredContext"].(string); ok {
		inferredCtx.ContextFocus = focus
	}
	if periodLabel, ok := contextData["inferredPeriod"].(string); ok {
		inferredCtx.PeriodLabel = periodLabel
	}

	return inferredCtx, nil
}
