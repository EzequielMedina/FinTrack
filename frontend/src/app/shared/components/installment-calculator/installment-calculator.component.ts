import { ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output, signal, inject, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { Observable, debounceTime, distinctUntilChanged, switchMap, of, BehaviorSubject, combineLatest } from 'rxjs';
import { InstallmentService } from '../../../services/installment.service';
import { InstallmentPreview, InstallmentPreviewRequest } from '../../../models';

export interface InstallmentCalculatorResult {
  installmentsCount: number;
  installmentAmount: number;
  totalToPay: number;
  totalInterest: number;
  adminFee: number;
  startDate: string;
  preview: InstallmentPreview;
}

@Component({
  selector: 'app-installment-calculator',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatTableModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatProgressSpinnerModule,
    MatCheckboxModule
  ],
  templateUrl: './installment-calculator.component.html',
  styleUrl: './installment-calculator.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InstallmentCalculatorComponent implements OnInit, OnChanges {
  private readonly fb = inject(FormBuilder);
  private readonly installmentService = inject(InstallmentService);

  @Input() amount: number = 0;
  @Input() cardId: string = '';
  @Input() disabled: boolean = false;
  @Input() showPreview: boolean = true;
  @Input() defaultInstallments: number = 3;
  
  @Output() calculationChanged = new EventEmitter<InstallmentCalculatorResult | null>();
  @Output() installmentsSelected = new EventEmitter<InstallmentCalculatorResult>();

  // Signals for reactive state
  calculatorForm = signal<FormGroup>(this.createForm());
  preview = signal<InstallmentPreview | null>(null);
  isCalculating = signal(false);
  error = signal<string | null>(null);
  availableInstallments = signal([3, 6, 9, 12, 15, 18, 21, 24]);

  // Table columns for preview
  displayedColumns = ['installmentNumber', 'dueDate', 'amount', 'principal', 'interest', 'fee'];

  // Subjects for reactive calculation
  private formChanges$ = new BehaviorSubject<any>(null);

  ngOnInit() {
    this.setupForm();
    this.setupReactiveCalculation();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['amount'] && this.amount > 0) {
      this.calculatorForm().patchValue({ amount: this.amount });
    }
    if (changes['disabled']) {
      if (this.disabled) {
        this.calculatorForm().disable();
      } else {
        this.calculatorForm().enable();
      }
    }
  }

  private createForm(): FormGroup {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    
    return this.fb.group({
      amount: [this.amount, [Validators.required, Validators.min(1000), Validators.max(500000)]],
      installmentsCount: [this.defaultInstallments, [Validators.required, Validators.min(1), Validators.max(24)]],
      startDate: [tomorrow, [Validators.required]],
      calculateNow: [true]
    });
  }

  private setupForm() {
    const form = this.createForm();
    this.calculatorForm.set(form);

    // Update amount if provided as input
    if (this.amount > 0) {
      form.patchValue({ amount: this.amount });
    }

    // Disable form if needed
    if (this.disabled) {
      form.disable();
    }
  }

  private setupReactiveCalculation() {
    // Listen to form changes
    this.calculatorForm().valueChanges.pipe(
      debounceTime(300),
      distinctUntilChanged()
    ).subscribe(formValue => {
      this.formChanges$.next(formValue);
    });

    // Reactive calculation
    this.formChanges$.pipe(
      debounceTime(500),
      distinctUntilChanged(),
      switchMap(formValue => {
        if (!formValue || !this.calculatorForm().valid || !formValue.calculateNow) {
          return of(null);
        }
        return this.calculatePreview(formValue);
      })
    ).subscribe({
      next: (preview) => {
        this.preview.set(preview);
        this.isCalculating.set(false);
        this.error.set(null);
        
        if (preview) {
          this.emitCalculationResult(preview);
        } else {
          this.calculationChanged.emit(null);
        }
      },
      error: (error) => {
        console.error('Error calculating installments:', error);
        this.error.set(error.message || 'Error al calcular las cuotas');
        this.isCalculating.set(false);
        this.preview.set(null);
        this.calculationChanged.emit(null);
      }
    });

    // Trigger initial calculation if amount is provided
    if (this.amount > 0) {
      setTimeout(() => this.formChanges$.next(this.calculatorForm().value), 100);
    }
  }

  private calculatePreview(formValue: any): Observable<InstallmentPreview | null> {
    if (!this.cardId) {
      // Mock calculation for preview without cardId
      return this.installmentService.calculateInstallmentOptions(
        formValue.amount,
        [formValue.installmentsCount]
      ).pipe(
        switchMap(options => of(options[0] || null))
      );
    }

    const request: InstallmentPreviewRequest = {
      amount: formValue.amount,
      installmentsCount: formValue.installmentsCount,
      startDate: this.formatDateForAPI(formValue.startDate)
    };

    this.isCalculating.set(true);
    return this.installmentService.previewInstallments(this.cardId, request);
  }

  private formatDateForAPI(date: Date): string {
    return date.toISOString().split('T')[0];
  }

  private emitCalculationResult(preview: InstallmentPreview) {
    const result: InstallmentCalculatorResult = {
      installmentsCount: preview.installmentsCount,
      installmentAmount: preview.installmentAmount,
      totalToPay: preview.totalToPay,
      totalInterest: preview.totalInterest,
      adminFee: preview.adminFee,
      startDate: preview.startDate,
      preview
    };

    this.calculationChanged.emit(result);
  }

  // Public methods
  onCalculateNowToggle() {
    const calculateNow = this.calculatorForm().get('calculateNow')?.value;
    if (calculateNow) {
      this.formChanges$.next(this.calculatorForm().value);
    } else {
      this.preview.set(null);
      this.calculationChanged.emit(null);
    }
  }

  onInstallmentsCountChange() {
    // Trigger recalculation when installments count changes
    if (this.calculatorForm().get('calculateNow')?.value) {
      this.formChanges$.next(this.calculatorForm().value);
    }
  }

  onStartDateChange() {
    // Trigger recalculation when start date changes
    if (this.calculatorForm().get('calculateNow')?.value) {
      this.formChanges$.next(this.calculatorForm().value);
    }
  }

  onSelectInstallments() {
    const currentPreview = this.preview();
    if (currentPreview) {
      const result: InstallmentCalculatorResult = {
        installmentsCount: currentPreview.installmentsCount,
        installmentAmount: currentPreview.installmentAmount,
        totalToPay: currentPreview.totalToPay,
        totalInterest: currentPreview.totalInterest,
        adminFee: currentPreview.adminFee,
        startDate: currentPreview.startDate,
        preview: currentPreview
      };
      
      this.installmentsSelected.emit(result);
    }
  }

  // Getters for template
  get form() {
    return this.calculatorForm();
  }

  get currentPreview() {
    return this.preview();
  }

  get isLoading() {
    return this.isCalculating();
  }

  get hasError() {
    return this.error();
  }

  get canSelectInstallments() {
    return this.preview() !== null && this.calculatorForm().valid;
  }

  // Utility methods for template
  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS'
    }).format(amount);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('es-AR');
  }

  getInterestRateDisplay(): string {
    const preview = this.preview();
    if (preview && preview.interestRate) {
      return `${preview.interestRate.toFixed(2)}% anual`;
    }
    return '';
  }

  getTotalSavings(): number {
    const preview = this.preview();
    if (preview) {
      return preview.totalToPay - preview.totalAmount;
    }
    return 0;
  }
}