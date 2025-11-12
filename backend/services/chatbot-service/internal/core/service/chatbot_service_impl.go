package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fintrack/chatbot-service/internal/core/ports"
	"github.com/google/uuid"
)

type ChatbotServiceImpl struct {
	data   ports.DataProvider
	llm    ports.LLMProvider
	report ports.ReportProvider
}

func NewChatbotService(data ports.DataProvider, llm ports.LLMProvider, report ports.ReportProvider) *ChatbotServiceImpl {
	return &ChatbotServiceImpl{data: data, llm: llm, report: report}
}

func (s *ChatbotServiceImpl) HandleQuery(ctx context.Context, req ports.ChatQueryRequest) (ports.ChatQueryResponse, error) {
	// Generate conversation ID if not provided
	if req.ConversationID == "" {
		req.ConversationID = uuid.New().String()
	}

	// Get previous conversation context for continuity
	var prevContext *ports.InferredContext
	history, _ := s.GetConversationHistory(ctx, req.UserID, req.ConversationID, 10)

	// Infer context from message (or use provided period/filters)
	var inferredCtx ports.InferredContext
	if req.Period.From.IsZero() || req.Period.To.IsZero() {
		// Auto-infer from message
		inferredCtx = InferContextFromMessage(req.Message, prevContext)
		req.Period = inferredCtx.Period
	} else {
		// Use provided period but still infer context
		inferredCtx = InferContextFromMessage(req.Message, prevContext)
		inferredCtx.Period = req.Period
	}

	contextFocus := inferredCtx.ContextFocus
	if customContext := getStringFromFilters(req.Filters, "contextFocus", ""); customContext != "" {
		contextFocus = customContext
		inferredCtx.ContextFocus = customContext
	}

	// Save user message to history
	userMsg := ports.ConversationMessage{
		ID:             uuid.New().String(),
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
		Role:           "user",
		Message:        req.Message,
		ContextData: map[string]any{
			"inferredPeriod":  inferredCtx.PeriodLabel,
			"inferredContext": inferredCtx.ContextFocus,
		},
		CreatedAt: time.Now(),
	}
	_ = s.SaveConversationMessage(ctx, userMsg)

	// Build conversational prompt with history
	system := buildConversationalPrompt(history, contextFocus)
	var user string
	var reply string

	// Obtener datos básicos siempre
	totals, err := s.data.GetTotals(ctx, req.UserID, req.Period.From, req.Period.To)
	if err != nil {
		return ports.ChatQueryResponse{}, err
	}

	instSummary, _ := s.data.GetInstallmentsSummary(ctx, req.UserID, req.Period.From, req.Period.To)
	plans, _ := s.data.GetInstallmentPlans(ctx, req.UserID)

	// Contexto específico basado en el enfoque seleccionado
	var ctxText string
	switch contextFocus {
	case "cards":
		byCard, _ := s.data.GetByCard(ctx, req.UserID, req.Period.From, req.Period.To)
		cardsInfo, _ := s.data.GetCardsInfo(ctx, req.UserID)
		byType, _ := s.data.GetByType(ctx, req.UserID, req.Period.From, req.Period.To)

		// Incluir pagos de cuotas como gastos con tarjetas
		installmentPayments := getVal(byType, "installment_payment")
		creditCharges := getVal(byType, "credit_charge")
		debitPurchases := getVal(byType, "debit_purchase")
		totalCardExpenses := getTotalFromByCard(byCard) + installmentPayments + creditCharges + debitPurchases

		ctxText = fmt.Sprintf(`TARJETAS: %s
GASTOS CON TARJETAS PERÍODO: $%.2f
- Consumos directos: $%.2f
- Pagos de cuotas: $%.2f  
- Cargos de crédito: $%.2f
- Compras débito: $%.2f
GASTOS TOTALES: $%.2f | INGRESOS: $%.2f`,
			formatCards(cardsInfo), totalCardExpenses, getTotalFromByCard(byCard),
			installmentPayments, creditCharges, debitPurchases, totals.Expenses, totals.Incomes)
	case "installments":
		byType, _ := s.data.GetByType(ctx, req.UserID, req.Period.From, req.Period.To)
		installmentPayments := getVal(byType, "installment_payment")
		byMonth, _ := s.data.GetInstallmentsByMonth(ctx, req.UserID)

		ctxText = fmt.Sprintf(`CUOTAS ACTIVAS: %d | VENCIDAS: %d | RESTANTE: $%.2f
PAGOS DE CUOTAS EN PERÍODO: $%.2f
PLANES: %s
CUOTAS PENDIENTES POR MES: %s
GASTOS TOTALES: $%.2f | INGRESOS: $%.2f`,
			instSummary.Active, instSummary.Overdue, instSummary.RemainingAmount,
			installmentPayments, formatPlans(plans), formatInstallmentsByMonth(byMonth), totals.Expenses, totals.Incomes)
	case "merchants":
		topMerchants, _ := s.data.GetTopMerchants(ctx, req.UserID, req.Period.From, req.Period.To, 5)
		ctxText = fmt.Sprintf(`TOP COMERCIOS: %s
GASTOS: $%.2f | INGRESOS: $%.2f`,
			formatMerchants(topMerchants), totals.Expenses, totals.Incomes)
	case "expenses":
		byType, _ := s.data.GetByType(ctx, req.UserID, req.Period.From, req.Period.To)
		installmentPayments := getVal(byType, "installment_payment")
		creditCharges := getVal(byType, "credit_charge")
		debitPurchases := getVal(byType, "debit_purchase")
		allExpenses := totals.Expenses + installmentPayments

		ctxText = fmt.Sprintf(`GASTOS DETALLADOS:
- Gastos directos: $%.2f
- Pagos de cuotas: $%.2f
- Cargos tarjetas: $%.2f  
- Compras débito: $%.2f
TOTAL GASTOS: $%.2f | INGRESOS: $%.2f`,
			totals.Expenses, installmentPayments, creditCharges, debitPurchases, allExpenses, totals.Incomes)
	case "income":
		ctxText = fmt.Sprintf(`INGRESOS TOTALES: $%.2f | GASTOS: $%.2f
PLANES ACTIVOS: %d`,
			totals.Incomes, totals.Expenses, instSummary.Active)
	default:
		// Contexto general compacto - incluir cuotas como gastos y planes con fechas
		byType, _ := s.data.GetByType(ctx, req.UserID, req.Period.From, req.Period.To)
		installmentPayments := getVal(byType, "installment_payment")
		allExpenses := totals.Expenses + installmentPayments
		byMonth, _ := s.data.GetInstallmentsByMonth(ctx, req.UserID)

		ctxText = fmt.Sprintf(`HOY: 2025-10-15
GASTOS TOTALES: $%.2f (directos: $%.2f + cuotas: $%.2f)
INGRESOS: $%.2f
PLANES ACTIVOS: %d | VENCIDOS: %d | RESTANTE: $%.2f
CUOTAS PENDIENTES POR MES: %s
DETALLE PLANES: %s`,
			allExpenses, totals.Expenses, installmentPayments, totals.Incomes,
			instSummary.Active, instSummary.Overdue, instSummary.RemainingAmount,
			formatInstallmentsByMonth(byMonth), formatPlans(plans))
	}

	user = fmt.Sprintf("Pregunta: %s\nContexto:\n%s", req.Message, ctxText)

	r, err := s.llm.Chat(ctx, system, user)

	if err != nil || r == "" {
		reply = fmt.Sprintf("Total gastado: %.2f. Total ingresado: %.2f.", totals.Expenses, totals.Incomes)
	} else {
		reply = r
	}

	// Generate contextual quick suggestions
	hasData := totals.Expenses > 0 || totals.Incomes > 0
	quickSuggestions := GenerateQuickSuggestions(contextFocus, hasData)

	// Save assistant response to history
	assistantMsg := ports.ConversationMessage{
		ID:             uuid.New().String(),
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
		Role:           "assistant",
		Message:        reply,
		ContextData: map[string]any{
			"inferredPeriod":  inferredCtx.PeriodLabel,
			"inferredContext": inferredCtx.ContextFocus,
			"totals":          totals,
		},
		CreatedAt: time.Now(),
	}
	_ = s.SaveConversationMessage(ctx, assistantMsg)

	return ports.ChatQueryResponse{
		Reply: reply,
		SuggestedActions: []ports.SuggestedAction{
			{Type: "generate_pdf", Params: map[string]any{"period": map[string]string{"from": req.Period.From.Format("2006-01-02"), "to": req.Period.To.Format("2006-01-02")}}},
			{Type: "show_chart", Params: map[string]any{"chartType": "bar", "groupBy": "account"}},
		},
		Insights:         []string{"Revisa categorías con mayor gasto", "Considera presupuesto semanal"},
		DataRefs:         map[string]any{"totals": totals, "installments": instSummary, "plans": plans},
		ConversationID:   req.ConversationID,
		InferredPeriod:   inferredCtx.PeriodLabel,
		InferredContext:  inferredCtx.ContextFocus,
		QuickSuggestions: quickSuggestions,
	}, nil
}

func (s *ChatbotServiceImpl) GeneratePDF(ctx context.Context, req ports.ReportRequest) ([]byte, error) {
	totals, _ := s.data.GetTotals(ctx, req.UserID, req.Period.From, req.Period.To)
	byType, _ := s.data.GetByType(ctx, req.UserID, req.Period.From, req.Period.To)
	topMerchants, _ := s.data.GetTopMerchants(ctx, req.UserID, req.Period.From, req.Period.To, 10)
	byAccount, _ := s.data.GetByAccountType(ctx, req.UserID, req.Period.From, req.Period.To)
	byCard, _ := s.data.GetByCard(ctx, req.UserID, req.Period.From, req.Period.To)

	data := ports.ReportData{
		Title:         req.Title,
		Period:        req.Period,
		Totals:        totals,
		ByType:        byType,
		TopMerchants:  topMerchants,
		ByAccountType: byAccount,
		ByCard:        byCard,
		Currency:      "ARS",
	}
	return s.report.Generate(ctx, data)
}

func (s *ChatbotServiceImpl) GenerateChartData(ctx context.Context, req ports.ChartRequest) (ports.ChartResponse, error) {
	group := strings.ToLower(req.GroupBy)
	switch group {
	case "account", "accounts":
		byAcc, err := s.data.GetByAccountType(ctx, req.UserID, req.Period.From, req.Period.To)
		if err != nil {
			return ports.ChartResponse{}, err
		}
		labels := make([]string, 0, len(byAcc))
		values := make([]float64, 0, len(byAcc))
		for k, v := range byAcc {
			labels = append(labels, k)
			values = append(values, v)
		}
		return ports.ChartResponse{Labels: labels, Datasets: []ports.ChartDataset{{Label: "Por cuenta", Data: values, BackgroundColor: []string{"#3b82f6"}}}, Meta: map[string]any{"currency": req.Currency}}, nil
	case "card", "cards":
		byCard, err := s.data.GetByCard(ctx, req.UserID, req.Period.From, req.Period.To)
		if err != nil {
			return ports.ChartResponse{}, err
		}
		labels := make([]string, 0, len(byCard))
		values := make([]float64, 0, len(byCard))
		for _, c := range byCard {
			labels = append(labels, fmt.Sprintf("%s •%s", c.Brand, c.LastFour))
			values = append(values, c.Total)
		}
		return ports.ChartResponse{Labels: labels, Datasets: []ports.ChartDataset{{Label: "Por tarjeta", Data: values, BackgroundColor: []string{"#10b981"}}}, Meta: map[string]any{"currency": req.Currency}}, nil
	default:
		return ports.ChartResponse{}, fmt.Errorf("groupBy inválido: use 'account' o 'card'")
	}
}

// Helpers
func getVal(m map[string]float64, k string) float64 {
	if m == nil {
		return 0
	}
	if v, ok := m[k]; ok {
		return v
	}
	return 0
}

func formatDate(t *time.Time) string {
	if t == nil {
		return "(sin fecha)"
	}
	return t.Format("2006-01-02")
}

// Helpers para formato humano
func formatMerchants(ms []ports.MerchantTotal) string {
	if len(ms) == 0 {
		return "(sin datos)"
	}
	n := len(ms)
	if n > 5 {
		n = 5
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		m := ms[i]
		out = append(out, fmt.Sprintf("%s (%.2f)", m.Merchant, m.Total))
	}
	return strings.Join(out, ", ")
}

func formatByCard(cs []ports.CardTotal) string {
	if len(cs) == 0 {
		return "(sin datos)"
	}
	n := len(cs)
	if n > 5 {
		n = 5
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		c := cs[i]
		out = append(out, fmt.Sprintf("%s ****%s (%.2f)", c.Brand, c.LastFour, c.Total))
	}
	return strings.Join(out, ", ")
}

func formatPlans(plans []ports.InstallmentPlanInfo) string {
	if len(plans) == 0 {
		return "(sin planes)"
	}
	n := len(plans)
	if n > 8 {
		n = 8
	} // Mostrar más planes para mejor contexto
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		p := plans[i]
		next := formatDate(p.NextDueDate)
		created := formatDate(p.CreatedAt)
		label := p.MerchantName
		if strings.TrimSpace(label) == "" {
			label = "Plan"
		}

		// Formato más detallado para mejor comprensión del LLM
		status := p.Status
		if p.Status == "active" {
			status = "activo"
		}
		if p.Status == "completed" {
			status = "completado"
		}
		if p.Status == "cancelled" {
			status = "cancelado"
		}

		out = append(out, fmt.Sprintf("[%s] %s '%s': %d cuotas, restante $%.2f, próximo vencimiento %s, creado %s, estado %s",
			p.ID[:8], label, p.Description, p.InstallmentsCount, p.RemainingAmount, next, created, status))
	}
	return strings.Join(out, " | ")
}

func formatTransactions(transactions []ports.TransactionDetail) string {
	if len(transactions) == 0 {
		return "(sin transacciones)"
	}
	n := len(transactions)
	if n > 5 {
		n = 5
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		t := transactions[i]
		out = append(out, fmt.Sprintf("[%s] %s: $%.2f en %s (%s) - %s",
			t.ID[:8], t.Type, t.Amount, t.MerchantName, t.Status, t.CreatedAt.Format("2006-01-02 15:04")))
	}
	return strings.Join(out, " | ")
}

func formatAccounts(accounts []ports.AccountInfo) string {
	if len(accounts) == 0 {
		return "(sin cuentas)"
	}
	out := make([]string, 0, len(accounts))
	for _, a := range accounts {
		out = append(out, fmt.Sprintf("[%s] %s: $%.2f %s (%s)",
			a.ID[:8], a.AccountType, a.Balance, a.Currency, a.Status))
	}
	return strings.Join(out, " | ")
}

func formatCards(cards []ports.CardInfo) string {
	if len(cards) == 0 {
		return "(sin tarjetas)"
	}
	out := make([]string, 0, len(cards))
	for _, c := range cards {
		out = append(out, fmt.Sprintf("[%s] %s ****%s (%s): límite $%.2f, deuda $%.2f (%s)",
			c.ID[:8], c.CardBrand, c.LastFour, c.CardType, c.CreditLimit, c.CurrentDebt, c.Status))
	}
	return strings.Join(out, " | ")
}

func formatExchangeRates(rates []ports.ExchangeRateInfo) string {
	if len(rates) == 0 {
		return "(sin cotizaciones)"
	}
	n := len(rates)
	if n > 3 {
		n = 3
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		r := rates[i]
		out = append(out, fmt.Sprintf("%s→%s: %.4f (%s)",
			r.FromCurrency, r.ToCurrency, r.Rate, r.CreatedAt.Format("2006-01-02")))
	}
	return strings.Join(out, " | ")
}

func formatInstallmentsByMonth(byMonth map[string]ports.InstallmentMonthSummary) string {
	if len(byMonth) == 0 {
		return "(sin cuotas futuras)"
	}

	// Ordenar los meses
	months := make([]string, 0, len(byMonth))
	for month := range byMonth {
		months = append(months, month)
	}
	sort.Strings(months)

	out := make([]string, 0, len(months))
	for _, month := range months {
		summary := byMonth[month]
		// Convertir "2025-11" a "Nov 2025"
		monthName := formatYearMonth(summary.YearMonth)
		out = append(out, fmt.Sprintf("%s: %d cuotas por $%.2f", monthName, summary.Count, summary.Total))
	}
	return strings.Join(out, " | ")
}

func formatYearMonth(yearMonth string) string {
	// Convierte "2025-11" a "Nov 2025"
	parts := strings.Split(yearMonth, "-")
	if len(parts) != 2 {
		return yearMonth
	}
	monthNames := map[string]string{
		"01": "Ene", "02": "Feb", "03": "Mar", "04": "Abr",
		"05": "May", "06": "Jun", "07": "Jul", "08": "Ago",
		"09": "Sep", "10": "Oct", "11": "Nov", "12": "Dic",
	}
	month, ok := monthNames[parts[1]]
	if !ok {
		return yearMonth
	}
	return fmt.Sprintf("%s %s", month, parts[0])
}

// Helper functions for context-aware prompts and filters
func getStringFromFilters(filters map[string]any, key, defaultValue string) string {
	if filters == nil {
		return defaultValue
	}
	if val, ok := filters[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func buildContextualPrompt(contextFocus string) string {
	base := "Asistente financiero FinTrack. ESPAÑOL OBLIGATORIO. Usa TODOS los datos del contexto proporcionado. IMPORTANTE: Los pagos de cuotas SON GASTOS."

	switch contextFocus {
	case "expenses":
		return base + " ENFOQUE: Analiza GASTOS detalladamente (incluye cuotas). Interpreta fechas, montos y categorías. Ejemplo: 'Gastos del día: $X (directos: $Y + cuotas: $Z)'"
	case "income":
		return base + " ENFOQUE: Analiza INGRESOS únicamente. Interpreta todas las fuentes de ingresos. Ejemplo: 'Ingresos del período: $X'"
	case "cards":
		return base + " ENFOQUE: Analiza TARJETAS (incluye cuotas y consumos). Lee información de límites, deudas y próximos vencimientos. Ejemplo: 'Gastos con tarjetas: $X'"
	case "installments":
		return base + " ENFOQUE: Analiza CUOTAS Y PLANES detalladamente. DEBES interpretar fechas de vencimiento, montos pendientes, estados. Si preguntan por vencimientos futuros, analiza las fechas 'próximo vencimiento' de cada plan. Ejemplo: 'Plan X vence el Y con monto Z'"
	case "merchants":
		return base + " ENFOQUE: Analiza COMERCIOS y patrones de gasto. Interpreta nombres de comercios y montos. Ejemplo: 'Gastaste más en: Comercio X ($Y)'"
	default:
		return base + " Interpreta TODA la información del contexto: fechas, montos, comercios, tarjetas, cuotas. Ejemplo: 'Gastos del día: $X (incluye cuotas)'"
	}
}

// Helper function to calculate total from byCard slice
func getTotalFromByCard(byCard []ports.CardTotal) float64 {
	var total float64
	for _, ct := range byCard {
		total += ct.Total
	}
	return total
}

// === CONVERSATIONAL METHODS ===

// GetConversationHistory retrieves conversation messages
func (s *ChatbotServiceImpl) GetConversationHistory(ctx context.Context, userID, conversationID string, limit int) ([]ports.ConversationMessage, error) {
	return s.data.GetConversationHistory(ctx, userID, conversationID, limit)
}

// SaveConversationMessage saves a conversation message
func (s *ChatbotServiceImpl) SaveConversationMessage(ctx context.Context, msg ports.ConversationMessage) error {
	return s.data.SaveConversationMessage(ctx, msg)
}

// buildConversationalPrompt creates a prompt with conversation history context
func buildConversationalPrompt(history []ports.ConversationMessage, contextFocus string) string {
	base := `Eres FinTrack Assistant, un asistente financiero personal amigable y eficiente.

INSTRUCCIONES:
1. Responde en ESPAÑOL de forma natural y conversacional
2. Usa el historial de la conversación para dar respuestas coherentes
3. Cuando el usuario dice "y eso?" o "¿cuál?", refiérete al mensaje anterior
4. Sé conciso pero informativo (máximo 3-4 líneas por defecto)
5. Si el usuario pide más detalles, entonces expándete
6. Usa emojis moderadamente (💰 💳 📊 ✅ ❌)
7. Formatea números con separadores: $322,000.95
8. IMPORTANTE: Los pagos de cuotas SON GASTOS

`

	// Add context focus specific instructions
	switch contextFocus {
	case "expenses":
		base += "ENFOQUE ACTUAL: Analiza GASTOS (incluye cuotas, compras, pagos)\n"
	case "income":
		base += "ENFOQUE ACTUAL: Analiza INGRESOS (sueldos, depósitos, cobros)\n"
	case "cards":
		base += "ENFOQUE ACTUAL: Analiza TARJETAS (límites, deudas, vencimientos)\n"
	case "installments":
		base += "ENFOQUE ACTUAL: Analiza CUOTAS Y PLANES (vencimientos, montos, estados)\n"
	case "merchants":
		base += "ENFOQUE ACTUAL: Analiza COMERCIOS (donde se gasta más)\n"
	default:
		base += "ENFOQUE ACTUAL: Información general financiera\n"
	}

	// Add last 3 messages from history for context
	if len(history) > 0 {
		base += "\nCONTEXTO DE CONVERSACIÓN:\n"
		start := 0
		if len(history) > 3 {
			start = len(history) - 3
		}
		for i := start; i < len(history); i++ {
			msg := history[i]
			role := "Usuario"
			if msg.Role == "assistant" {
				role = "Tú"
			}
			base += fmt.Sprintf("%s: %s\n", role, msg.Message)
		}
	}

	return base
}