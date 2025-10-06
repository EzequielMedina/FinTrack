import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal, inject, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialogModule } from '@angular/material/dialog';
import { Observable } from 'rxjs';
import { InstallmentService } from '../../../services/installment.service';
import { InstallmentPlan, Installment, InstallmentStatus, InstallmentPlanStatus, InstallmentPlanResponse } from '../../../models';

export interface InstallmentPaymentAction {
  type: 'pay';
  installmentId: string;
  installment: Installment;
  planId: string;
}

@Component({
  selector: 'app-installment-plan-detail',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    MatDialogModule
  ],
  templateUrl: './installment-plan-detail.component.html',
  styleUrl: './installment-plan-detail.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InstallmentPlanDetailComponent implements OnInit, OnChanges {
  private readonly installmentService = inject(InstallmentService);

  @Input() planId: string = '';
  @Input() initialPlan: InstallmentPlan | null = null; // If provided, skip loading
  @Input() showActions: boolean = true;
  @Input() showHeader: boolean = true;

  @Output() installmentAction = new EventEmitter<InstallmentPaymentAction>();
  @Output() planLoaded = new EventEmitter<InstallmentPlan>();
  @Output() error = new EventEmitter<string>();

  // Signals for reactive state
  currentPlan = signal<InstallmentPlan | null>(null);
  installments = signal<Installment[]>([]);
  isLoading = signal(false);
  errorMessage = signal<string | null>(null);

  // Table configuration
  displayedColumns = ['installmentNumber', 'dueDate', 'amount', 'status', 'paidDate', 'actions'];
  readonly InstallmentStatus = InstallmentStatus;
  readonly InstallmentPlanStatus = InstallmentPlanStatus;

  ngOnInit() {
    if (this.initialPlan) {
      this.currentPlan.set(this.initialPlan);
      this.processInstallments(this.initialPlan);
    } else if (this.planId) {
      this.loadPlanDetail();
    }
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['initialPlan'] && this.initialPlan) {
      this.currentPlan.set(this.initialPlan);
      this.processInstallments(this.initialPlan);
    }
    
    if (changes['planId'] && this.planId && !this.initialPlan) {
      this.loadPlanDetail();
    }
  }

  private loadPlanDetail() {
    if (!this.planId) return;

    this.isLoading.set(true);
    this.errorMessage.set(null);

    this.installmentService.getInstallmentPlanDetails(this.planId).subscribe({
      next: (response: InstallmentPlanResponse) => {
        this.currentPlan.set(response.installmentPlan);
        this.processInstallments(response.installmentPlan);
        this.isLoading.set(false);
        this.planLoaded.emit(response.installmentPlan);
      },
      error: (error) => {
        console.error('Error loading plan details:', error);
        const message = error.message || 'Error al cargar el detalle del plan';
        this.errorMessage.set(message);
        this.isLoading.set(false);
        this.error.emit(message);
      }
    });
  }

  private processInstallments(plan: InstallmentPlan) {
    // Generate installments array if not provided
    // In a real implementation, this would come from the API
    const installments: Installment[] = [];
    const startDate = new Date(plan.startDate);

    for (let i = 1; i <= plan.installmentsCount; i++) {
      const dueDate = new Date(startDate);
      dueDate.setMonth(startDate.getMonth() + i - 1);
      
      const isPaid = i <= plan.paidInstallments;
      const isOverdue = !isPaid && dueDate < new Date();
      
      installments.push({
        id: `${plan.id}-${i}`,
        plan_id: plan.id,
        installment_number: i,
        amount: plan.installmentAmount,
        due_date: dueDate.toISOString().split('T')[0],
        status: isPaid ? InstallmentStatus.PAID : 
                isOverdue ? InstallmentStatus.OVERDUE : 
                InstallmentStatus.PENDING,
        paid_date: isPaid ? dueDate.toISOString().split('T')[0] : undefined,
        remaining_amount: isPaid ? 0 : plan.installmentAmount,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        // Legacy properties for backward compatibility
        planId: plan.id,
        installmentNumber: i,
        dueDate: dueDate.toISOString().split('T')[0],
        paidDate: isPaid ? dueDate.toISOString().split('T')[0] : undefined,
        remainingAmount: isPaid ? 0 : plan.installmentAmount,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      });
    }

    this.installments.set(installments);
  }

  // Public methods
  onPayInstallment(installment: Installment) {
    this.installmentAction.emit({
      type: 'pay',
      installmentId: installment.id,
      installment,
      planId: installment.plan_id
    });
  }

  refresh() {
    if (this.planId) {
      this.loadPlanDetail();
    }
  }

  // Utility methods for template
  getStatusClass(status: InstallmentStatus): string {
    switch (status) {
      case InstallmentStatus.PAID:
        return 'status-paid';
      case InstallmentStatus.PENDING:
        return 'status-pending';
      case InstallmentStatus.OVERDUE:
        return 'status-overdue';
      case InstallmentStatus.CANCELLED:
        return 'status-cancelled';
      default:
        return 'status-pending';
    }
  }

  getStatusIcon(status: InstallmentStatus): string {
    switch (status) {
      case InstallmentStatus.PAID:
        return 'check_circle';
      case InstallmentStatus.PENDING:
        return 'schedule';
      case InstallmentStatus.OVERDUE:
        return 'warning';
      case InstallmentStatus.CANCELLED:
        return 'cancel';
      default:
        return 'schedule';
    }
  }

  getStatusText(status: InstallmentStatus): string {
    switch (status) {
      case InstallmentStatus.PAID:
        return 'Pagada';
      case InstallmentStatus.PENDING:
        return 'Pendiente';
      case InstallmentStatus.OVERDUE:
        return 'Vencida';
      case InstallmentStatus.CANCELLED:
        return 'Cancelada';
      default:
        return 'Pendiente';
    }
  }

  getPlanStatusClass(status: InstallmentPlanStatus): string {
    switch (status) {
      case InstallmentPlanStatus.ACTIVE:
        return 'plan-status-active';
      case InstallmentPlanStatus.COMPLETED:
        return 'plan-status-completed';
      case InstallmentPlanStatus.CANCELLED:
        return 'plan-status-cancelled';
      case InstallmentPlanStatus.SUSPENDED:
        return 'plan-status-suspended';
      default:
        return 'plan-status-pending';
    }
  }

  getPlanStatusText(status: InstallmentPlanStatus): string {
    switch (status) {
      case InstallmentPlanStatus.ACTIVE:
        return 'Activo';
      case InstallmentPlanStatus.COMPLETED:
        return 'Completado';
      case InstallmentPlanStatus.CANCELLED:
        return 'Cancelado';
      case InstallmentPlanStatus.SUSPENDED:
        return 'Suspendido';
      default:
        return 'Pendiente';
    }
  }

  canPayInstallment(installment: Installment): boolean {
    const plan = this.currentPlan();
    if (!plan) return false;
    
    return (plan.status === InstallmentPlanStatus.ACTIVE || plan.status === InstallmentPlanStatus.SUSPENDED) &&
           (installment.status === InstallmentStatus.PENDING || installment.status === InstallmentStatus.OVERDUE);
  }

  getProgressPercentage(): number {
    const plan = this.currentPlan();
    if (!plan || plan.installmentsCount === 0) return 0;
    return Math.round((plan.paidInstallments / plan.installmentsCount) * 100);
  }

  getNextDueInstallment(): Installment | null {
    return this.installments().find(inst => 
      inst.status === InstallmentStatus.PENDING || inst.status === InstallmentStatus.OVERDUE
    ) || null;
  }

  getOverdueCount(): number {
    return this.installments().filter(inst => inst.status === InstallmentStatus.OVERDUE).length;
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS'
    }).format(amount);
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('es-AR');
  }

  isOverdue(installment: Installment): boolean {
    return installment.status === InstallmentStatus.OVERDUE;
  }

  isPaid(installment: Installment): boolean {
    return installment.status === InstallmentStatus.PAID;
  }

  // Getters for template
  get currentPlanData() {
    return this.currentPlan();
  }

  get currentInstallments() {
    return this.installments();
  }

  get isLoadingState() {
    return this.isLoading();
  }

  get hasError() {
    return this.errorMessage();
  }

  get hasPlan() {
    return this.currentPlan() !== null;
  }

  get hasInstallments() {
    return this.installments().length > 0;
  }
}