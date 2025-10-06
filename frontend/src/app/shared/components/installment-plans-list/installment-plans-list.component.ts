import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal, inject, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatChipsModule } from '@angular/material/chips';
import { MatBadgeModule } from '@angular/material/badge';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDividerModule } from '@angular/material/divider';
import { Observable, Subject, BehaviorSubject, combineLatest } from 'rxjs';
import { takeUntil, switchMap, map, startWith } from 'rxjs/operators';
import { InstallmentService } from '../../../services/installment.service';
import { InstallmentPlan, InstallmentPlanStatus, InstallmentPlansListResponse } from '../../../models';

export interface InstallmentPlanAction {
  type: 'view' | 'pay' | 'cancel';
  planId: string;
  plan: InstallmentPlan;
}

@Component({
  selector: 'app-installment-plans-list',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressBarModule,
    MatChipsModule,
    MatBadgeModule,
    MatMenuModule,
    MatPaginatorModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    MatDividerModule
  ],
  templateUrl: './installment-plans-list.component.html',
  styleUrl: './installment-plans-list.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InstallmentPlansListComponent implements OnInit, OnDestroy {
  private readonly installmentService = inject(InstallmentService);
  private readonly destroy$ = new Subject<void>();

  @Input() cardId: string = '';
  @Input() userId: string = '';
  @Input() showCardInfo: boolean = false;
  @Input() maxPlansToShow: number = 0; // 0 = show all
  @Input() showPagination: boolean = true;
  @Input() statusFilter: InstallmentPlanStatus | '' = '';
  @Input() autoRefresh: boolean = false;

  @Output() planAction = new EventEmitter<InstallmentPlanAction>();
  @Output() plansLoaded = new EventEmitter<InstallmentPlan[]>();

  // Signals for reactive state
  plans = signal<InstallmentPlan[]>([]);
  isLoading = signal(true);
  error = signal<string | null>(null);
  totalPlans = signal(0);
  
  // Pagination
  currentPage = signal(0);
  pageSize = signal(10);
  
  // Filter and refresh subjects
  private refreshTrigger$ = new BehaviorSubject<void>(undefined);
  private pageChange$ = new BehaviorSubject<{ pageIndex: number; pageSize: number }>({ pageIndex: 0, pageSize: 10 });

  // Enum for template access
  readonly InstallmentPlanStatus = InstallmentPlanStatus;

  ngOnInit() {
    this.setupDataLoading();
    
    if (this.autoRefresh) {
      this.setupAutoRefresh();
    }
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private setupDataLoading() {
    // Combine refresh trigger and page changes
    combineLatest([
      this.refreshTrigger$,
      this.pageChange$
    ]).pipe(
      takeUntil(this.destroy$),
      switchMap(([_, pageInfo]) => {
        this.isLoading.set(true);
        this.error.set(null);
        
        // Choose the appropriate service method based on inputs
        if (this.cardId) {
          return this.installmentService.getInstallmentPlansByCard(
            this.cardId, 
            pageInfo.pageIndex + 1, 
            pageInfo.pageSize
          );
        } else if (this.userId) {
          return this.installmentService.getUserInstallmentPlans(
            this.userId,
            this.statusFilter || undefined,
            pageInfo.pageIndex + 1,
            pageInfo.pageSize
          );
        } else {
          throw new Error('Either cardId or userId must be provided');
        }
      })
    ).subscribe({
      next: (response: InstallmentPlansListResponse) => {
        let plansToShow = response.plans || [];
        
        // Apply status filter if specified (for cardId-based requests)
        if (this.statusFilter && this.cardId) {
          plansToShow = plansToShow.filter(plan => plan.status === this.statusFilter);
        }
        
        // Apply maxPlansToShow limit if specified
        if (this.maxPlansToShow > 0 && plansToShow.length > 0) {
          plansToShow = plansToShow.slice(0, this.maxPlansToShow);
        }
        
        this.plans.set(plansToShow);
        // For filtered results, show filtered count
        this.totalPlans.set(this.statusFilter && this.cardId ? plansToShow.length : (response.total || 0));
        this.isLoading.set(false);
        this.plansLoaded.emit(plansToShow);
      },
      error: (error) => {
        console.error('Error loading installment plans:', error);
        this.error.set(error.message || 'Error al cargar los planes de cuotas');
        this.isLoading.set(false);
        this.plans.set([]);
      }
    });
  }

  private setupAutoRefresh() {
    // Auto refresh every 30 seconds
    setInterval(() => {
      if (!this.isLoading()) {
        this.refreshTrigger$.next();
      }
    }, 30000);
  }

  // Public methods
  refresh() {
    this.refreshTrigger$.next();
  }

  onPageChange(event: PageEvent) {
    this.currentPage.set(event.pageIndex);
    this.pageSize.set(event.pageSize);
    this.pageChange$.next({ pageIndex: event.pageIndex, pageSize: event.pageSize });
  }

  onViewPlan(plan: InstallmentPlan) {
    this.planAction.emit({
      type: 'view',
      planId: plan.id,
      plan
    });
  }

  onPayInstallment(plan: InstallmentPlan) {
    this.planAction.emit({
      type: 'pay',
      planId: plan.id,
      plan
    });
  }

  onCancelPlan(plan: InstallmentPlan) {
    this.planAction.emit({
      type: 'cancel',
      planId: plan.id,
      plan
    });
  }

  // Utility methods for template
  getStatusClass(status: InstallmentPlanStatus): string {
    switch (status) {
      case InstallmentPlanStatus.ACTIVE:
        return 'status-active';
      case InstallmentPlanStatus.COMPLETED:
        return 'status-completed';
      case InstallmentPlanStatus.CANCELLED:
        return 'status-cancelled';
      case InstallmentPlanStatus.SUSPENDED:
        return 'status-suspended';
      default:
        return 'status-pending';
    }
  }

  getStatusIcon(status: InstallmentPlanStatus): string {
    switch (status) {
      case InstallmentPlanStatus.ACTIVE:
        return 'play_circle';
      case InstallmentPlanStatus.COMPLETED:
        return 'check_circle';
      case InstallmentPlanStatus.CANCELLED:
        return 'cancel';
      case InstallmentPlanStatus.SUSPENDED:
        return 'warning';
      default:
        return 'schedule';
    }
  }

  getStatusText(status: InstallmentPlanStatus): string {
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

  getProgressPercentage(plan: InstallmentPlan): number {
    if (plan.installmentsCount === 0) return 0;
    return Math.round((plan.paidInstallments / plan.installmentsCount) * 100);
  }

  getProgressColor(plan: InstallmentPlan): string {
    const percentage = this.getProgressPercentage(plan);
    if (percentage === 100) return 'primary';
    if (percentage >= 75) return 'accent';
    if (percentage >= 50) return 'primary';
    if (plan.status === InstallmentPlanStatus.SUSPENDED) return 'warn';
    return 'primary';
  }

  getNextPaymentInfo(plan: InstallmentPlan): { amount: number; date: string } | null {
    // This would typically come from the API, for now we calculate it
    if (plan.status === InstallmentPlanStatus.COMPLETED || plan.status === InstallmentPlanStatus.CANCELLED) {
      return null;
    }
    
    const nextInstallmentNumber = plan.paidInstallments + 1;
    if (nextInstallmentNumber > plan.installmentsCount) {
      return null;
    }

    // Calculate next payment date (this is a simplified calculation)
    const startDate = new Date(plan.startDate);
    const nextPaymentDate = new Date(startDate);
    nextPaymentDate.setMonth(startDate.getMonth() + nextInstallmentNumber - 1);

    return {
      amount: plan.installmentAmount,
      date: nextPaymentDate.toISOString().split('T')[0]
    };
  }

  isSuspended(plan: InstallmentPlan): boolean {
    return plan.status === InstallmentPlanStatus.SUSPENDED;
  }

  canCancel(plan: InstallmentPlan): boolean {
    return plan.status === InstallmentPlanStatus.ACTIVE && plan.paidInstallments < plan.installmentsCount;
  }

  canPay(plan: InstallmentPlan): boolean {
    return plan.status === InstallmentPlanStatus.ACTIVE;
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

  getRemainingInstallments(plan: InstallmentPlan): number {
    return plan.installmentsCount - plan.paidInstallments;
  }

  getDaysUntilNextPayment(plan: InstallmentPlan): number {
    const nextPayment = this.getNextPaymentInfo(plan);
    if (!nextPayment) return 0;
    
    const today = new Date();
    const paymentDate = new Date(nextPayment.date);
    const diffTime = paymentDate.getTime() - today.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    return Math.max(0, diffDays);
  }

  getMerchantDisplay(plan: InstallmentPlan): string {
    return plan.merchantName || plan.description || 'Compra en cuotas';
  }

  trackByPlanId(index: number, plan: InstallmentPlan): string {
    return plan.id;
  }

  // Getters for template
  get currentPlans() {
    return this.plans();
  }

  get isLoadingState() {
    return this.isLoading();
  }

  get hasError() {
    return this.error();
  }

  get totalCount() {
    return this.totalPlans();
  }

  get showPaginationControls() {
    return this.showPagination && this.totalCount > this.pageSize();
  }

  get hasPlans() {
    return this.currentPlans.length > 0;
  }

  get emptyStateMessage() {
    if (this.statusFilter) {
      return `No hay planes de cuotas con estado: ${this.getStatusText(this.statusFilter)}`;
    }
    return this.cardId ? 'No hay planes de cuotas para esta tarjeta' : 'No hay planes de cuotas';
  }
}