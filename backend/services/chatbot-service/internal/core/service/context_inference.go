package service

import (
	"strings"
	"time"

	"github.com/fintrack/chatbot-service/internal/core/ports"
)

// InferContextFromMessage analyzes the user's message to automatically detect:
// 1. Time period (today, yesterday, this month, last week, etc.)
// 2. Context focus (expenses, income, cards, installments, merchants, general)
// 3. Optionally uses previous conversation context for continuity
func InferContextFromMessage(message string, prevContext *ports.InferredContext) ports.InferredContext {
	msgLower := strings.ToLower(message)
	now := time.Now()
	loc := now.Location()

	result := ports.InferredContext{
		ContextFocus:     "general",
		PeriodLabel:      "this month",
		DetectedKeywords: []string{},
	}

	// Use previous context as starting point if available
	if prevContext != nil {
		result.Period = prevContext.Period
		result.ContextFocus = prevContext.ContextFocus
		result.PeriodLabel = prevContext.PeriodLabel
	}

	// === DETECT TIME PERIOD ===
	periodDetected := false

	// Today
	if containsAny(msgLower, "hoy", "today", "día de hoy") {
		result.Period = getPeriodToday(now, loc)
		result.PeriodLabel = "today"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:today")
		periodDetected = true
	}

	// Yesterday
	if !periodDetected && containsAny(msgLower, "ayer", "yesterday") {
		result.Period = getPeriodYesterday(now, loc)
		result.PeriodLabel = "yesterday"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:yesterday")
		periodDetected = true
	}

	// This week
	if !periodDetected && containsAny(msgLower, "esta semana", "this week", "semana actual") {
		result.Period = getPeriodThisWeek(now, loc)
		result.PeriodLabel = "this week"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:this_week")
		periodDetected = true
	}

	// Last week
	if !periodDetected && containsAny(msgLower, "semana pasada", "last week", "última semana") {
		result.Period = getPeriodLastWeek(now, loc)
		result.PeriodLabel = "last week"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:last_week")
		periodDetected = true
	}

	// This month
	if !periodDetected && containsAny(msgLower, "este mes", "this month", "mes actual", "el mes") {
		result.Period = getPeriodThisMonth(now, loc)
		result.PeriodLabel = "this month"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:this_month")
		periodDetected = true
	}

	// Last month
	if !periodDetected && containsAny(msgLower, "mes pasado", "last month", "último mes") {
		result.Period = getPeriodLastMonth(now, loc)
		result.PeriodLabel = "last month"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:last_month")
		periodDetected = true
	}

	// Last 30 days
	if !periodDetected && containsAny(msgLower, "últimos 30", "last 30", "30 días") {
		result.Period = getPeriodLast30Days(now, loc)
		result.PeriodLabel = "last 30 days"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:last_30_days")
		periodDetected = true
	}

	// Last 7 days
	if !periodDetected && containsAny(msgLower, "últimos 7", "last 7", "7 días") {
		result.Period = getPeriodLast7Days(now, loc)
		result.PeriodLabel = "last 7 days"
		result.DetectedKeywords = append(result.DetectedKeywords, "period:last_7_days")
		periodDetected = true
	}

	// If no period detected, default to this month
	if !periodDetected {
		result.Period = getPeriodThisMonth(now, loc)
	}

	// === DETECT CONTEXT FOCUS ===
	contextDetected := false

	// Cards context
	if containsAny(msgLower, "tarjeta", "card", "crédito", "débito", "credit", "debit", "plástico") {
		result.ContextFocus = "cards"
		result.DetectedKeywords = append(result.DetectedKeywords, "context:cards")
		contextDetected = true
	}

	// Installments context
	if !contextDetected && containsAny(msgLower, "cuota", "installment", "plan", "planes", "financiación", "financing") {
		result.ContextFocus = "installments"
		result.DetectedKeywords = append(result.DetectedKeywords, "context:installments")
		contextDetected = true
	}

	// Expenses context
	if !contextDetected && containsAny(msgLower, "gast", "expense", "compra", "purchase", "pago", "payment") {
		result.ContextFocus = "expenses"
		result.DetectedKeywords = append(result.DetectedKeywords, "context:expenses")
		contextDetected = true
	}

	// Income context
	if !contextDetected && containsAny(msgLower, "ingreso", "income", "cobr", "recib", "ganancia", "earning") {
		result.ContextFocus = "income"
		result.DetectedKeywords = append(result.DetectedKeywords, "context:income")
		contextDetected = true
	}

	// Merchants context
	if !contextDetected && containsAny(msgLower, "comercio", "merchant", "tienda", "store", "negocio", "donde gast") {
		result.ContextFocus = "merchants"
		result.DetectedKeywords = append(result.DetectedKeywords, "context:merchants")
		contextDetected = true
	}

	// Accounts context
	if !contextDetected && containsAny(msgLower, "cuenta", "account", "saldo", "balance") {
		result.ContextFocus = "accounts"
		result.DetectedKeywords = append(result.DetectedKeywords, "context:accounts")
		contextDetected = true
	}

	// If no specific context detected but previous context exists, keep it (conversation continuity)
	// Otherwise, it stays as "general"

	return result
}

// Helper: check if text contains any of the given substrings
func containsAny(text string, substrings ...string) bool {
	for _, s := range substrings {
		if strings.Contains(text, s) {
			return true
		}
	}
	return false
}

// === PERIOD CALCULATION FUNCTIONS ===

func getPeriodToday(now time.Time, loc *time.Location) ports.Period {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return ports.Period{From: startOfDay, To: endOfDay}
}

func getPeriodYesterday(now time.Time, loc *time.Location) ports.Period {
	yesterday := now.AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return ports.Period{From: startOfDay, To: endOfDay}
}

func getPeriodThisWeek(now time.Time, loc *time.Location) ports.Period {
	// Week starts on Sunday (0)
	weekday := int(now.Weekday())
	startOfWeek := now.AddDate(0, 0, -weekday)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return ports.Period{From: startOfWeek, To: endOfDay}
}

func getPeriodLastWeek(now time.Time, loc *time.Location) ports.Period {
	weekday := int(now.Weekday())
	startOfLastWeek := now.AddDate(0, 0, -weekday-7)
	endOfLastWeek := startOfLastWeek.AddDate(0, 0, 6)
	startOfLastWeek = time.Date(startOfLastWeek.Year(), startOfLastWeek.Month(), startOfLastWeek.Day(), 0, 0, 0, 0, loc)
	endOfLastWeek = time.Date(endOfLastWeek.Year(), endOfLastWeek.Month(), endOfLastWeek.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return ports.Period{From: startOfLastWeek, To: endOfLastWeek}
}

func getPeriodThisMonth(now time.Time, loc *time.Location) ports.Period {
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return ports.Period{From: firstOfMonth, To: endOfDay}
}

func getPeriodLastMonth(now time.Time, loc *time.Location) ports.Period {
	firstOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, loc)
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	lastOfLastMonth := firstOfThisMonth.Add(-time.Second)
	return ports.Period{From: firstOfLastMonth, To: lastOfLastMonth}
}

func getPeriodLast30Days(now time.Time, loc *time.Location) ports.Period {
	from := now.AddDate(0, 0, -30)
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return ports.Period{From: from, To: to}
}

func getPeriodLast7Days(now time.Time, loc *time.Location) ports.Period {
	from := now.AddDate(0, 0, -7)
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return ports.Period{From: from, To: to}
}

// GenerateQuickSuggestions creates contextual follow-up question suggestions
func GenerateQuickSuggestions(context string, hasData bool) []string {
	if !hasData {
		return []string{
			"¿Cuáles son mis tarjetas?",
			"¿Tengo planes de cuotas?",
			"Muéstrame un resumen",
		}
	}

	switch context {
	case "expenses":
		return []string{
			"¿En qué comercios gasté más?",
			"¿Cuánto gasté con tarjetas?",
			"Ver por tipo de gasto",
		}
	case "cards":
		return []string{
			"¿Cuál tiene más deuda?",
			"¿Cuándo vencen los pagos?",
			"Ver límites disponibles",
		}
	case "installments":
		return []string{
			"¿Cuándo vence la próxima cuota?",
			"¿Puedo cancelar algún plan?",
			"Ver proyección de pagos",
		}
	case "merchants":
		return []string{
			"¿Cuánto gasté en el top comercio?",
			"Ver otros comercios",
			"Analizar gastos por categoría",
		}
	case "income":
		return []string{
			"¿De dónde vienen mis ingresos?",
			"Comparar con gastos",
			"Ver balance",
		}
	default:
		return []string{
			"Ver gastos del mes",
			"Estado de mis tarjetas",
			"Planes de cuotas activos",
		}
	}
}
