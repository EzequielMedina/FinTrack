package report

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fintrack/report-service/internal/core/domain/dto"
	"github.com/fintrack/report-service/internal/core/service"
	"github.com/fintrack/report-service/pkg/pdf"
	"github.com/gin-gonic/gin"
)

// ReportHandler maneja las peticiones HTTP de reportes
type ReportHandler struct {
	reportService service.ReportService
}

// NewReportHandler crea un nuevo handler de reportes
func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetTransactionReport obtiene el reporte de transacciones
// @Summary Obtener reporte de transacciones
// @Description Obtiene un reporte detallado de transacciones por período
// @Tags reports
// @Accept json
// @Produce json
// @Param user_id query string true "ID del usuario"
// @Param start_date query string false "Fecha inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha fin (YYYY-MM-DD)"
// @Param type query string false "Tipo de transacción"
// @Success 200 {object} dto.TransactionReportResponse
// @Router /api/v1/reports/transactions [get]
func (h *ReportHandler) GetTransactionReport(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de start_date inválido (usar YYYY-MM-DD)"})
			return
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de end_date inválido (usar YYYY-MM-DD)"})
			return
		}
	}

	req := &dto.TransactionReportRequest{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Type:      c.Query("type"),
		GroupBy:   c.Query("group_by"),
	}

	report, err := h.reportService.GetTransactionReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetInstallmentReport obtiene el reporte de cuotas
// @Summary Obtener reporte de cuotas
// @Description Obtiene un reporte detallado de planes de cuotas y pagos
// @Tags reports
// @Accept json
// @Produce json
// @Param user_id query string true "ID del usuario"
// @Param status query string false "Estado de los planes (active, completed, overdue)"
// @Success 200 {object} dto.InstallmentReportResponse
// @Router /api/v1/reports/installments [get]
func (h *ReportHandler) GetInstallmentReport(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	req := &dto.InstallmentReportRequest{
		UserID: userID,
		Status: c.Query("status"),
	}

	report, err := h.reportService.GetInstallmentReport(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Error en GetInstallmentReport: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetAccountReport obtiene el reporte de cuentas
// @Summary Obtener reporte de cuentas
// @Description Obtiene un resumen de todas las cuentas y tarjetas del usuario
// @Tags reports
// @Accept json
// @Produce json
// @Param user_id query string true "ID del usuario"
// @Success 200 {object} dto.AccountReportResponse
// @Router /api/v1/reports/accounts [get]
func (h *ReportHandler) GetAccountReport(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	req := &dto.AccountReportRequest{
		UserID: userID,
	}

	report, err := h.reportService.GetAccountReport(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Error en GetAccountReport: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetExpenseIncomeReport obtiene el reporte de gastos vs ingresos
// @Summary Obtener reporte de gastos vs ingresos
// @Description Obtiene un análisis de gastos e ingresos por período
// @Tags reports
// @Accept json
// @Produce json
// @Param user_id query string true "ID del usuario"
// @Param start_date query string true "Fecha inicio (YYYY-MM-DD)"
// @Param end_date query string true "Fecha fin (YYYY-MM-DD)"
// @Param group_by query string false "Agrupar por (day, week, month)"
// @Success 200 {object} dto.ExpenseIncomeReportResponse
// @Router /api/v1/reports/expenses-income [get]
func (h *ReportHandler) GetExpenseIncomeReport(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date y end_date son requeridos"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de start_date inválido (usar YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de end_date inválido (usar YYYY-MM-DD)"})
		return
	}

	req := &dto.ExpenseIncomeReportRequest{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   c.Query("group_by"),
	}

	report, err := h.reportService.GetExpenseIncomeReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetNotificationReport obtiene el reporte de notificaciones
// @Summary Obtener reporte de notificaciones
// @Description Obtiene estadísticas de envío de notificaciones
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string false "Fecha inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha fin (YYYY-MM-DD)"
// @Success 200 {object} dto.NotificationReportResponse
// @Router /api/v1/reports/notifications [get]
func (h *ReportHandler) GetNotificationReport(c *gin.Context) {
	var startDate, endDate time.Time
	var err error

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de start_date inválido (usar YYYY-MM-DD)"})
			return
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de end_date inválido (usar YYYY-MM-DD)"})
			return
		}
	}

	req := &dto.NotificationReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	report, err := h.reportService.GetNotificationReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetInstallmentReportPDF obtiene el reporte de cuotas en PDF
// @Summary Obtener reporte de cuotas en PDF
// @Description Genera y descarga un PDF del reporte de cuotas
// @Tags reports
// @Produce application/pdf
// @Param user_id query string true "ID del usuario"
// @Param status query string false "Estado de los planes (active, completed, overdue)"
// @Success 200 {file} binary
// @Router /api/v1/reports/installments/pdf [get]
func (h *ReportHandler) GetInstallmentReportPDF(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	req := &dto.InstallmentReportRequest{
		UserID: userID,
		Status: c.Query("status"),
	}

	report, err := h.reportService.GetInstallmentReport(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Error en GetInstallmentReport (PDF): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := pdf.GenerateInstallmentReportPDF(report)
	if err != nil {
		log.Printf("❌ Error generando PDF de cuotas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando PDF"})
		return
	}

	filename := fmt.Sprintf("reporte-cuotas-%s.pdf", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetAccountReportPDF obtiene el reporte de cuentas en PDF
// @Summary Obtener reporte de cuentas en PDF
// @Description Genera y descarga un PDF del reporte de cuentas
// @Tags reports
// @Produce application/pdf
// @Param user_id query string true "ID del usuario"
// @Success 200 {file} binary
// @Router /api/v1/reports/accounts/pdf [get]
func (h *ReportHandler) GetAccountReportPDF(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	req := &dto.AccountReportRequest{
		UserID: userID,
	}

	report, err := h.reportService.GetAccountReport(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Error en GetAccountReport (PDF): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := pdf.GenerateAccountReportPDF(report)
	if err != nil {
		log.Printf("❌ Error generando PDF de cuentas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando PDF"})
		return
	}

	filename := fmt.Sprintf("reporte-cuentas-%s.pdf", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetTransactionReportPDF obtiene el reporte de transacciones en PDF
// @Summary Obtener reporte de transacciones en PDF
// @Description Genera y descarga un PDF del reporte de transacciones
// @Tags reports
// @Produce application/pdf
// @Param user_id query string true "ID del usuario"
// @Param start_date query string false "Fecha inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha fin (YYYY-MM-DD)"
// @Param type query string false "Tipo de transacción"
// @Success 200 {file} binary
// @Router /api/v1/reports/transactions/pdf [get]
func (h *ReportHandler) GetTransactionReportPDF(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de start_date inválido (usar YYYY-MM-DD)"})
			return
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de end_date inválido (usar YYYY-MM-DD)"})
			return
		}
	}

	req := &dto.TransactionReportRequest{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Type:      c.Query("type"),
		GroupBy:   c.Query("group_by"),
	}

	report, err := h.reportService.GetTransactionReport(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Error en GetTransactionReport (PDF): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := pdf.GenerateTransactionReportPDF(report)
	if err != nil {
		log.Printf("❌ Error generando PDF de transacciones: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando PDF"})
		return
	}

	filename := fmt.Sprintf("reporte-transacciones-%s.pdf", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetExpenseIncomeReportPDF obtiene el reporte de gastos vs ingresos en PDF
// @Summary Obtener reporte de gastos vs ingresos en PDF
// @Description Genera y descarga un PDF del análisis de gastos e ingresos
// @Tags reports
// @Produce application/pdf
// @Param user_id query string true "ID del usuario"
// @Param start_date query string true "Fecha inicio (YYYY-MM-DD)"
// @Param end_date query string true "Fecha fin (YYYY-MM-DD)"
// @Param group_by query string false "Agrupar por (day, week, month)"
// @Success 200 {file} binary
// @Router /api/v1/reports/expenses-income/pdf [get]
func (h *ReportHandler) GetExpenseIncomeReportPDF(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id es requerido"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date y end_date son requeridos"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de start_date inválido (usar YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de end_date inválido (usar YYYY-MM-DD)"})
		return
	}

	req := &dto.ExpenseIncomeReportRequest{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   c.Query("group_by"),
	}

	report, err := h.reportService.GetExpenseIncomeReport(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ Error en GetExpenseIncomeReport (PDF): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := pdf.GenerateExpenseIncomeReportPDF(report)
	if err != nil {
		log.Printf("❌ Error generando PDF de gastos vs ingresos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando PDF"})
		return
	}

	filename := fmt.Sprintf("reporte-gastos-ingresos-%s.pdf", time.Now().Format("2006-01-02"))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
