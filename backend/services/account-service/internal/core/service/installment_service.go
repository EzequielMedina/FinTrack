package service

import (
	"fmt"
	"time"

	"github.com/fintrack/account-service/internal/core/domain/entities"
	"github.com/fintrack/account-service/internal/core/ports"
	"github.com/fintrack/account-service/internal/infrastructure/clients"
	carddto "github.com/fintrack/account-service/internal/infrastructure/entrypoints/handlers/card/dto"
	"github.com/google/uuid"
)

type InstallmentService struct {
	installmentRepo      ports.InstallmentRepositoryInterface
	installmentPlanRepo  ports.InstallmentPlanRepositoryInterface
	installmentAuditRepo ports.InstallmentPlanAuditRepositoryInterface
	cardRepo             ports.CardRepositoryInterface
	transactionClient    *clients.TransactionClient
}

func NewInstallmentService(installmentRepo ports.InstallmentRepositoryInterface, installmentPlanRepo ports.InstallmentPlanRepositoryInterface, installmentAuditRepo ports.InstallmentPlanAuditRepositoryInterface, cardRepo ports.CardRepositoryInterface) *InstallmentService {
	return &InstallmentService{
		installmentRepo:      installmentRepo,
		installmentPlanRepo:  installmentPlanRepo,
		installmentAuditRepo: installmentAuditRepo,
		cardRepo:             cardRepo,
		transactionClient:    clients.NewTransactionClient(),
	}
}

// CalculateInstallmentPlan calcula un plan de cuotas sin persistirlo
func (s *InstallmentService) CalculateInstallmentPlan(amount float64, installmentsCount int, startDate time.Time, interestRate float64) (*carddto.InstallmentPreviewResponse, error) {
	if installmentsCount <= 0 {
		return nil, fmt.Errorf("number of installments must be greater than 0")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// Calcular monto total con interés
	totalAmount := amount * (1 + interestRate/100)
	installmentAmount := totalAmount / float64(installmentsCount)

	// Crear response de preview (implementación básica)
	response := &carddto.InstallmentPreviewResponse{
		TotalAmount:       totalAmount,
		InstallmentAmount: installmentAmount,
		InstallmentsCount: installmentsCount,
		InterestRate:      interestRate,
	}

	return response, nil
}

// CreateInstallmentPlan crea y persiste un plan de cuotas
func (s *InstallmentService) CreateInstallmentPlan(req *carddto.CreateInstallmentPlanRequest) (*entities.InstallmentPlan, error) {
	// Verificar que la tarjeta existe y es de crédito
	card, err := s.cardRepo.GetByID(req.CardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	if card.CardType != "credit" {
		return nil, fmt.Errorf("installment plans are only available for credit cards")
	}

	// Crear plan básico
	plan := &entities.InstallmentPlan{
		ID:                uuid.New().String(),
		CardID:            req.CardID,
		UserID:            req.UserID,
		TransactionID:     uuid.New().String(),
		TotalAmount:       req.TotalAmount,
		InstallmentsCount: req.InstallmentsCount,
		InstallmentAmount: req.TotalAmount / float64(req.InstallmentsCount),
		StartDate:         req.StartDate,
		Status:            "active",
		RemainingAmount:   req.TotalAmount,
		Description:       req.Description,
		MerchantName:      req.MerchantName,
		MerchantID:        req.MerchantID,
		InterestRate:      req.InterestRate,
		AdminFee:          req.AdminFee,
	}

	// Crear plan en base de datos
	createdPlan, err := s.installmentPlanRepo.Create(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to create installment plan: %w", err)
	}

	// Crear cuotas individuales
	fmt.Printf("🟢🟢🟢 INSTALLMENT_SERVICE - About to create %d individual installments for plan %s 🟢🟢🟢\n", req.InstallmentsCount, createdPlan.ID)
	fmt.Printf("🟢🟢🟢 INSTALLMENT_SERVICE - TotalAmount: %.2f, UserID: %s 🟢🟢🟢\n", req.TotalAmount, req.UserID)
	installmentAmount := req.TotalAmount / float64(req.InstallmentsCount)

	for i := 1; i <= req.InstallmentsCount; i++ {
		// Calcular fecha de vencimiento (mensual)
		dueDate := req.StartDate.AddDate(0, i-1, 0)

		installment := &entities.Installment{
			ID:                uuid.New().String(),
			PlanID:            createdPlan.ID,
			InstallmentNumber: i,
			Amount:            installmentAmount,
			DueDate:           dueDate,
			Status:            entities.InstallmentStatusPending,
			PaidAmount:        0.0,
			RemainingAmount:   installmentAmount,
			LateFee:           0.0,
			PenaltyAmount:     0.0,
			GracePeriodDays:   7,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		fmt.Printf("DEBUG - Creating installment %d: Amount=%.2f, DueDate=%v\n",
			i, installmentAmount, dueDate)

		_, err := s.installmentRepo.Create(installment)
		if err != nil {
			fmt.Printf("ERROR - Failed to create installment %d: %v\n", i, err)
			return nil, fmt.Errorf("failed to create installment %d: %w", i, err)
		}
		fmt.Printf("DEBUG - Successfully created installment %d\n", i)
	}

	fmt.Printf("🔥🔥🔥 DEBUG - BEFORE TRANSACTION CLIENT CALL 🔥🔥🔥\n")

	// Registrar transacción en transaction service (async)
	go func() {
		// Obtener información de la cuenta de la tarjeta
		cardWithAccount, err := s.cardRepo.GetByIDWithAccount(req.CardID)
		if err != nil {
			fmt.Printf("Warning: Failed to get card account for transaction recording: %v\n", err)
			return
		}

		_, err = s.transactionClient.CreateInstallmentTransaction(
			req.UserID,
			cardWithAccount.Account.ID,
			req.CardID,
			req.TotalAmount,
			req.InstallmentsCount,
			createdPlan.ID,
			req.Description,
			req.MerchantName,
			req.Reference,
		)
		if err != nil {
			fmt.Printf("Warning: Failed to record installment transaction: %v\n", err)
		}
	}()

	// El plan de cuotas se ha creado exitosamente
	// La carga de la tarjeta se manejará en el CardService
	return createdPlan, nil
}

// GetInstallmentPlansByCard obtiene todos los planes de una tarjeta
func (s *InstallmentService) GetInstallmentPlansByCard(cardID string, page, pageSize int) ([]*entities.InstallmentPlan, int64, error) {
	offset := (page - 1) * pageSize
	return s.installmentPlanRepo.GetByCard(cardID, "", pageSize, offset)
}

// GetInstallmentPlan obtiene un plan específico
func (s *InstallmentService) GetInstallmentPlan(planID string) (*entities.InstallmentPlan, error) {
	return s.installmentPlanRepo.GetByIDWithInstallments(planID)
}

// PayInstallment procesa el pago de una cuota
func (s *InstallmentService) PayInstallment(req *carddto.PayInstallmentRequest) (*entities.Installment, error) {
	// Obtener la cuota
	installment, err := s.installmentRepo.GetByID(req.InstallmentID)
	if err != nil {
		return nil, fmt.Errorf("installment not found: %w", err)
	}

	if installment.Status == "paid" {
		return nil, fmt.Errorf("installment already paid")
	}

	if req.Amount < installment.Amount {
		return nil, fmt.Errorf("payment amount insufficient. Required: %.2f, provided: %.2f", installment.Amount, req.Amount)
	}

	// Marcar como pagada
	installment.Status = "paid"
	now := time.Now()
	installment.PaidDate = &now

	// Actualizar en base de datos
	updatedInstallment, err := s.installmentRepo.Update(installment)
	if err != nil {
		return nil, fmt.Errorf("failed to update installment: %w", err)
	}

	// Registrar pago en transaction service (async)
	go func() {
		// Obtener el plan para metadata
		plan, err := s.installmentPlanRepo.GetByID(installment.PlanID)
		if err != nil {
			fmt.Printf("Warning: Failed to get installment plan for transaction recording: %v\n", err)
			return
		}

		// Obtener información de la tarjeta y cuenta
		cardWithAccount, err := s.cardRepo.GetByIDWithAccount(plan.CardID)
		if err != nil {
			fmt.Printf("Warning: Failed to get card account for transaction recording: %v\n", err)
			return
		}

		_, err = s.transactionClient.CreateInstallmentPaymentTransaction(
			plan.UserID,
			cardWithAccount.Account.ID,
			plan.CardID,
			req.Amount,
			req.InstallmentID,
			plan.ID,
			fmt.Sprintf("%d", installment.InstallmentNumber),
			plan.Description,
		)
		if err != nil {
			fmt.Printf("Warning: Failed to record installment payment transaction: %v\n", err)
		}
	}()

	return updatedInstallment, nil
}

// GetInstallmentPlansByUser obtiene planes por usuario
func (s *InstallmentService) GetInstallmentPlansByUser(userID string, status string, page, pageSize int) ([]*entities.InstallmentPlan, int64, error) {
	offset := (page - 1) * pageSize
	return s.installmentPlanRepo.GetByUser(userID, status, pageSize, offset)
}

// CancelInstallmentPlan cancela un plan de cuotas activo
func (s *InstallmentService) CancelInstallmentPlan(planID, reason string, cancelledBy string) (*entities.InstallmentPlan, error) {
	// Obtener el plan
	plan, err := s.installmentPlanRepo.GetByID(planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	if plan.Status != "active" {
		return nil, fmt.Errorf("only active plans can be cancelled")
	}

	// Cancelar plan
	plan.Status = "cancelled"
	now := time.Now()
	plan.CancelledAt = &now

	updatedPlan, err := s.installmentPlanRepo.Update(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel plan: %w", err)
	}

	// Registrar cancelación en transaction service (async)
	go func() {
		// Obtener información de la tarjeta y cuenta
		cardWithAccount, err := s.cardRepo.GetByIDWithAccount(plan.CardID)
		if err != nil {
			fmt.Printf("Warning: Failed to get card account for cancellation transaction recording: %v\n", err)
			return
		}

		_, err = s.transactionClient.CreateInstallmentCancellationTransaction(
			plan.UserID,
			cardWithAccount.Account.ID,
			plan.CardID,
			plan.RemainingAmount,
			plan.ID,
			reason,
		)
		if err != nil {
			fmt.Printf("Warning: Failed to record installment cancellation transaction: %v\n", err)
		}
	}()

	return updatedPlan, nil
}

// SuspendInstallmentPlan suspende un plan de cuotas
func (s *InstallmentService) SuspendInstallmentPlan(planID, reason string, suspendedBy string) (*entities.InstallmentPlan, error) {
	plan, err := s.installmentPlanRepo.GetByID(planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	plan.Status = "suspended"
	return s.installmentPlanRepo.Update(plan)
}

// ReactivateInstallmentPlan reactiva un plan de cuotas suspendido
func (s *InstallmentService) ReactivateInstallmentPlan(planID, reason string, reactivatedBy string) (*entities.InstallmentPlan, error) {
	plan, err := s.installmentPlanRepo.GetByID(planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	plan.Status = "active"
	return s.installmentPlanRepo.Update(plan)
}

// GetInstallment obtiene una cuota específica
func (s *InstallmentService) GetInstallment(installmentID string) (*entities.Installment, error) {
	return s.installmentRepo.GetByID(installmentID)
}

// GetInstallmentsByPlan obtiene todas las cuotas de un plan
func (s *InstallmentService) GetInstallmentsByPlan(planID string) ([]*entities.Installment, error) {
	return s.installmentRepo.GetByPlan(planID)
}

// GetOverdueInstallments obtiene cuotas vencidas por usuario
func (s *InstallmentService) GetOverdueInstallments(userID string, limit, offset int) ([]*entities.Installment, int64, error) {
	return s.installmentRepo.GetOverdue(userID, limit, offset)
}

// GetUpcomingInstallments obtiene cuotas próximas a vencer
func (s *InstallmentService) GetUpcomingInstallments(userID string, days int, limit, offset int) ([]*entities.Installment, int64, error) {
	return s.installmentRepo.GetUpcoming(userID, days, limit, offset)
}

// GetInstallmentHistory obtiene el historial de auditoría de una cuota
func (s *InstallmentService) GetInstallmentHistory(installmentID string) ([]*entities.InstallmentPlanAudit, error) {
	if s.installmentAuditRepo == nil {
		return nil, fmt.Errorf("audit repository not available")
	}
	return s.installmentAuditRepo.GetByInstallment(installmentID)
}

// GetInstallmentSummary obtiene resumen de cuotas por usuario
func (s *InstallmentService) GetInstallmentSummary(userID string) (map[string]interface{}, error) {
	// Implementación básica
	summary := map[string]interface{}{
		"user_id": userID,
		"message": "Summary functionality to be implemented",
	}
	return summary, nil
}

// GetMonthlyInstallmentLoad obtiene carga mensual de cuotas por usuario
func (s *InstallmentService) GetMonthlyInstallmentLoad(userID string, year, month int) (map[string]interface{}, error) {
	// Implementación básica
	load := map[string]interface{}{
		"user_id": userID,
		"year":    year,
		"month":   month,
		"message": "Monthly load functionality to be implemented",
	}
	return load, nil
}
