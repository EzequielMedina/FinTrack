package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fintrack/notification-service/internal/config"
	"github.com/fintrack/notification-service/internal/core/domain/entities"
)

// EmailJSClient cliente para la API de EmailJS
type EmailJSClient struct {
	config     *config.EmailJSConfig
	httpClient *http.Client
	baseURL    string
}

// EmailJSRequest estructura para la API de EmailJS
type EmailJSRequest struct {
	ServiceID      string            `json:"service_id"`
	TemplateID     string            `json:"template_id"`
	UserID         string            `json:"user_id"`
	TemplateParams map[string]string `json:"template_params"`
	// Eliminamos AccessToken ya que no se usa en el ejemplo oficial
}

// NewEmailJSClient crea un nuevo cliente de EmailJS
func NewEmailJSClient(cfg *config.EmailJSConfig) *EmailJSClient {
	return &EmailJSClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.emailjs.com/api/v1.0/email/send",
	}
}

// SendCardDueNotification envía una notificación de vencimiento de tarjeta
func (c *EmailJSClient) SendCardDueNotification(notification *entities.CardDueNotification) error {
	htmlContent := c.buildEmailHTML(notification)

	templateParams := map[string]string{
		"from_name":         c.config.FromName,
		"subject":           fmt.Sprintf("Recordatorio: Tu tarjeta %s vence mañana 📅", notification.CardName),
		"to_email":          notification.UserEmail,
		"reply_to":          c.config.ReplyTo,
		"html_content":      htmlContent,
		"user_name":         notification.UserName,
		"card_name":         notification.CardName,
		"bank_name":         notification.BankName,
		"last_four":         notification.LastFour,
		"due_date":          notification.DueDate.Format("02/01/2006"),
		"total_amount":      fmt.Sprintf("$%.2f", notification.TotalPendingAmount),
		"installment_count": fmt.Sprintf("%d", notification.PendingInstallments),
	}

	request := EmailJSRequest{
		ServiceID:      c.config.ServiceID,
		TemplateID:     c.config.TemplateID,
		UserID:         c.config.PublicKey,
		TemplateParams: templateParams,
		// Eliminamos AccessToken siguiendo el ejemplo oficial de JS
	}

	return c.sendRequest(request)
}

// buildEmailHTML construye el HTML del email con los datos de la notificación
func (c *EmailJSClient) buildEmailHTML(notification *entities.CardDueNotification) string {
	installmentsHTML := c.buildInstallmentsHTML(notification.InstallmentDetails)

	// Mensaje diferente dependiendo si hay cuotas pendientes o no
	var alertMessage string

	if notification.PendingInstallments == 0 {
		alertMessage = `
		<div style="background: #d1ecf1; padding: 15px; border-radius: 5px; border-left: 4px solid #17a2b8; margin: 20px 0;">
			<h4 style="color: #0c5460; margin: 0 0 10px 0;">ℹ️ Sin cuotas pendientes</h4>
			<p style="color: #0c5460; margin: 0;">Tu tarjeta vence mañana pero no tienes cuotas pendientes para esta fecha. ¡Todo al día! 🎉</p>
		</div>`
	} else {
		alertMessage = fmt.Sprintf(`
		<div style="background: #fff3cd; padding: 15px; border-radius: 5px; border-left: 4px solid #ffc107; margin: 20px 0;">
			<h4 style="color: #856404; margin: 0 0 10px 0;">💰 Total a pagar: %s</h4>
			<p style="color: #856404; margin: 0;">Tienes %d cuotas pendientes que vencen mañana o antes:</p>
		</div>`, fmt.Sprintf("$%.2f", notification.TotalPendingAmount), notification.PendingInstallments)
	}

	// Solo el contenido interno para que el template de EmailJS maneje la estructura externa
	html := fmt.Sprintf(`
		<h2 style="color: #333;">Hola %s, tu tarjeta vence mañana 📅</h2>
		<div style="background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
			<h3 style="color: #667eea; margin-top: 0;">%s - %s (****%s)</h3>
			<p style="font-size: 16px; color: #333; margin: 15px 0;">
				<strong>Fecha de vencimiento:</strong> %s
			</p>
			%s
			%s
			<div style="text-align: center; margin-top: 30px;">
				<p style="color: #666; font-size: 14px;">
					⏰ Recuerda revisar tu tarjeta antes del vencimiento.
				</p>
			</div>
		</div>`,
		notification.UserName,
		notification.CardName,
		notification.BankName,
		notification.LastFour,
		notification.DueDate.Format("02/01/2006"),
		alertMessage,
		installmentsHTML,
	)

	return html
}

// buildInstallmentsHTML construye el HTML de la lista de cuotas
func (c *EmailJSClient) buildInstallmentsHTML(installments []entities.InstallmentSummary) string {
	if len(installments) == 0 {
		return `<p style="color: #666; font-style: italic;">No hay cuotas pendientes.</p>`
	}

	html := `<div style="margin: 20px 0;"><h4 style="color: #333; margin-bottom: 15px;">📋 Detalle de cuotas:</h4><ul style="list-style: none; padding: 0;">`

	for _, installment := range installments {
		html += fmt.Sprintf(`
		<li style="background: #f8f9fa; padding: 12px; margin: 8px 0; border-radius: 5px; border-left: 3px solid #667eea;">
			<div style="display: flex; justify-content: space-between; align-items: center;">
				<div>
					<strong style="color: #333;">%s</strong><br>
					<small style="color: #666;">%s - Cuota %d/%d</small>
				</div>
				<div style="text-align: right;">
					<strong style="color: #667eea; font-size: 16px;">$%.2f</strong><br>
					<small style="color: #666;">Vence: %s</small>
				</div>
			</div>
		</li>`,
			installment.GetInstallmentText(),
			installment.MerchantName,
			installment.InstallmentNum,
			installment.TotalInstallments,
			installment.Amount,
			installment.DueDate.Format("02/01/2006"),
		)
	}

	html += `</ul></div>`
	return html
}

// sendRequest envía la solicitud HTTP a EmailJS
func (c *EmailJSClient) sendRequest(request EmailJSRequest) error {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Headers para simular browser y evitar bloqueo de EmailJS
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
	req.Header.Set("Origin", "http://localhost:4200")
	req.Header.Set("Referer", "http://localhost:4200/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("EmailJS API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
