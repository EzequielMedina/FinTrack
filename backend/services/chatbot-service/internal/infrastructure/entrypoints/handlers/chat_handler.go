package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/fintrack/chatbot-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	svc ports.ChatbotService
}

func NewChatHandler(svc ports.ChatbotService) *ChatHandler { return &ChatHandler{svc: svc} }

// POST /api/chat/query
func (h *ChatHandler) Query(c *gin.Context) {
	var req struct {
		UserID         string `json:"userId"`
		Message        string `json:"message"`
		ConversationID string `json:"conversationId"` // New: for conversational context
		Period         struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"period"`
		Filters map[string]any `json:"filters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body", "details": err.Error()})
		return
	}
	// Fallback: tomar el userId del header si no viene en el body
	if req.UserID == "" {
		hID := c.GetHeader("X-User-ID")
		if hID != "" {
			req.UserID = hID
		}
	}
	from, _ := time.Parse("2006-01-02", req.Period.From)
	to, _ := time.Parse("2006-01-02", req.Period.To)
	// Period is now optional - will be inferred if not provided
	// Normalizar periodo: start-of-day para from, end-of-day para to
	if !from.IsZero() {
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	}
	if !to.IsZero() {
		to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), to.Location())
	}

	resp, err := h.svc.HandleQuery(c, ports.ChatQueryRequest{
		UserID:         req.UserID,
		Message:        req.Message,
		ConversationID: req.ConversationID,
		Filters:        req.Filters,
		Period:         ports.Period{From: from, To: to},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GET /api/chat/history/:conversationId
func (h *ChatHandler) GetHistory(c *gin.Context) {
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversationId is required"})
		return
	}

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-ID header is required"})
		return
	}

	history, err := h.svc.GetConversationHistory(c, userID, conversationID, 50) // Get last 50 messages
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversationId": conversationID,
		"messages":       history,
		"total":          len(history),
	})
}

// POST /api/chat/report/pdf
func (h *ChatHandler) ReportPDF(c *gin.Context) {
	var req struct {
		UserID string `json:"userId"`
		Title  string `json:"title"`
		Period struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"period"`
		GroupBy       string `json:"groupBy"`
		IncludeCharts bool   `json:"includeCharts"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body", "details": err.Error()})
		return
	}
	// Fallback de userId por header
	if req.UserID == "" {
		hID := c.GetHeader("X-User-ID")
		if hID != "" {
			req.UserID = hID
		}
	}
	from, _ := time.Parse("2006-01-02", req.Period.From)
	to, _ := time.Parse("2006-01-02", req.Period.To)
	if from.IsZero() && to.IsZero() {
		now := time.Now()
		firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		from = firstOfMonth
		to = now
	} else {
		if !from.IsZero() {
			from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
		}
		if !to.IsZero() {
			to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), to.Location())
		}
	}
	pdfBytes, err := h.svc.GeneratePDF(c, ports.ReportRequest{
		UserID: req.UserID, Title: req.Title, Period: ports.Period{From: from, To: to}, GroupBy: req.GroupBy, IncludeCharts: req.IncludeCharts,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// POST /api/chat/report/chart
func (h *ChatHandler) ReportChart(c *gin.Context) {
	var req struct {
		UserID string `json:"userId"`
		Period struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"period"`
		GroupBy  string `json:"groupBy"`
		Currency string `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body", "details": err.Error()})
		return
	}
	// Fallback de userId por header
	if req.UserID == "" {
		hID := c.GetHeader("X-User-ID")
		if hID != "" {
			req.UserID = hID
		}
	}
	from, _ := time.Parse("2006-01-02", req.Period.From)
	to, _ := time.Parse("2006-01-02", req.Period.To)
	if from.IsZero() && to.IsZero() {
		now := time.Now()
		firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		from = firstOfMonth
		to = now
	} else {
		if !from.IsZero() {
			from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
		}
		if !to.IsZero() {
			to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), to.Location())
		}
	}
	// Validar groupBy permitido
	gb := strings.ToLower(req.GroupBy)
	if gb != "account" && gb != "accounts" && gb != "card" && gb != "cards" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupBy inválido", "details": "Use 'account' o 'card'"})
		return
	}
	chart, err := h.svc.GenerateChartData(c, ports.ChartRequest{
		UserID: req.UserID, Period: ports.Period{From: from, To: to}, GroupBy: req.GroupBy, Currency: req.Currency,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chart)
}

// GET /health
func (h *ChatHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "chatbot-service"})
}