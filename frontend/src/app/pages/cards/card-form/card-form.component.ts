import { Component, Inject, inject, OnInit, signal, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule, AbstractControl, ValidationErrors } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule, MatSnackBar } from '@angular/material/snack-bar';
import { Subject, takeUntil, debounceTime, distinctUntilChanged } from 'rxjs';
import { CardService } from '../../../services/card.service';
import { UserService } from '../../../services/user.service';
import { AuthService } from '../../../services/auth.service';
import { Card, CardType, CardBrand, CreateCardRequest, UpdateCardRequest, CardFormData, CardFormErrors } from '../../../models';

interface DialogData {
  mode: 'create' | 'edit';
  card?: Card;
}

@Component({
  selector: 'app-card-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  templateUrl: './card-form.component.html',
  styleUrls: ['./card-form.component.css']
  // Removemos OnPush para evitar problemas de actualización del modal
})
export class CardFormComponent implements OnInit, OnDestroy {
  private readonly fb = inject(FormBuilder);
  private readonly cardService = inject(CardService);
  private readonly userService = inject(UserService);
  private readonly authService = inject(AuthService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dialogRef = inject(MatDialogRef<CardFormComponent>);
  private readonly destroy$ = new Subject<void>();

  // Signals para el estado del componente
  saving = signal(false);
  detectedBrand = signal<CardBrand | null>(null);
  formErrors = signal<CardFormErrors>({});

  // Formulario reactivo
  cardForm!: FormGroup;

  // Exponer enums para usar en el template
  readonly CardType = CardType;
  readonly CardBrand = CardBrand;

  // Opciones para selects
  readonly cardTypes = [
    { value: CardType.CREDIT, label: 'Tarjeta de Crédito' },
    { value: CardType.DEBIT, label: 'Tarjeta de Débito' }
  ];

  readonly expirationMonths = Array.from({ length: 12 }, (_, i) => ({
    value: i + 1,
    label: (i + 1).toString().padStart(2, '0')
  }));

  readonly expirationYears = Array.from({ length: 20 }, (_, i) => {
    const year = new Date().getFullYear() + i;
    return { value: year, label: year.toString() };
  });

  constructor(@Inject(MAT_DIALOG_DATA) public data: DialogData) {}

  ngOnInit(): void {
    this.initializeForm();
    this.setupFormValidation();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private initializeForm(): void {
    this.cardForm = this.fb.group({
      userId: ['', [Validators.required]],
      cardType: [CardType.DEBIT, [Validators.required]],
      cardNumber: [''],
      holderName: ['', [Validators.required, Validators.minLength(2), Validators.maxLength(50)]],
      expirationMonth: ['', [Validators.required]],
      expirationYear: ['', [Validators.required]],
      cvv: [''],
      nickname: ['', [Validators.maxLength(30)]],
      currency: ['ARS', [Validators.required]]
    });

    // Configurar validadores condicionales después de crear el formulario
    if (this.data.mode === 'create') {
      this.cardForm.get('cardNumber')?.setValidators([Validators.required, this.cardNumberValidator.bind(this)]);
      this.cardForm.get('cvv')?.setValidators([Validators.required, this.cvvValidator.bind(this)]);
    }

    // Si estamos editando, llenar el formulario
    if (this.data.mode === 'edit' && this.data.card) {
      this.populateFormForEdit(this.data.card);
    }

    // Configurar el nombre del titular con el usuario actual si está creando
    if (this.data.mode === 'create') {
      const currentUser = this.authService.currentUserSig();
      if (currentUser) {
        this.cardForm.patchValue({
          userId: currentUser.id,
          holderName: `${currentUser.firstName} ${currentUser.lastName}`.toUpperCase()
        });
      }
    }

    // Actualizar validadores
    this.cardForm.get('cardNumber')?.updateValueAndValidity();
    this.cardForm.get('cvv')?.updateValueAndValidity();
  }

  private populateFormForEdit(card: Card): void {
    this.cardForm.patchValue({
      userId: card.userId,
      cardType: card.cardType,
      holderName: card.holderName,
      expirationMonth: card.expirationMonth,
      expirationYear: card.expirationYear,
      nickname: card.nickname
    });

    // Para edición, no necesitamos validadores de número de tarjeta y CVV
    // ya que estos campos no se muestran en modo edición
  }

  private setupFormValidation(): void {
    // Validación con debounce del número de tarjeta
    this.cardForm.get('cardNumber')?.valueChanges.pipe(
      debounceTime(300), // Esperar 300ms después del último cambio
      distinctUntilChanged(), // Solo procesar si el valor cambió
      takeUntil(this.destroy$)
    ).subscribe(value => {
      if (value && value.length >= 4) {
        const validation = this.cardService.validateCardNumber(value);
        this.detectedBrand.set(validation.brand || null);
        
        // Actualizar validación del CVV según la marca detectada
        const cvvControl = this.cardForm.get('cvv');
        if (cvvControl) {
          cvvControl.updateValueAndValidity();
        }
      } else {
        this.detectedBrand.set(null);
      }
    });

    // Validación con debounce de la fecha de expiración
    this.cardForm.get('expirationMonth')?.valueChanges.pipe(
      debounceTime(200),
      takeUntil(this.destroy$)
    ).subscribe(() => {
      this.validateExpirationDate();
    });

    this.cardForm.get('expirationYear')?.valueChanges.pipe(
      debounceTime(200),
      takeUntil(this.destroy$)
    ).subscribe(() => {
      this.validateExpirationDate();
    });
  }

  private validateExpirationDate(): void {
    const month = this.cardForm.get('expirationMonth')?.value;
    const year = this.cardForm.get('expirationYear')?.value;
    
    if (month && year) {
      const validation = this.cardService.validateExpirationDate(month, year);
      if (!validation.isValid) {
        const errors = this.formErrors();
        validation.errors.forEach(error => {
          errors[error.field as keyof CardFormErrors] = error.message;
        });
        this.formErrors.set({ ...errors });
      } else {
        const errors = this.formErrors();
        delete errors.expirationMonth;
        delete errors.expirationYear;
        this.formErrors.set({ ...errors });
      }
    }
  }

  // Validadores personalizados
  private cardNumberValidator(control: AbstractControl): ValidationErrors | null {
    if (!control.value) return null;
    
    const validation = this.cardService.validateCardNumber(control.value);
    if (!validation.isValid) {
      return { cardNumber: validation.errors[0]?.message || 'Número de tarjeta inválido' };
    }
    return null;
  }

  private cvvValidator(control: AbstractControl): ValidationErrors | null {
    if (!control.value) return null;
    
    const brand = this.detectedBrand();
    const validation = this.cardService.validateCVV(control.value, brand || undefined);
    if (!validation.isValid) {
      return { cvv: validation.errors[0]?.message || 'CVV inválido' };
    }
    return null;
  }

  // Métodos de formateo para el template
  formatCardNumber(event: any): void {
    const input = event.target;
    const value = input.value.replace(/\s/g, '');
    const formatted = this.cardService.formatCardNumber(value);
    input.value = formatted;
    this.cardForm.patchValue({ cardNumber: value });
  }

  onSubmit(): void {
    if (this.cardForm.invalid) {
      this.markFormGroupTouched();
      return;
    }

    this.saving.set(true);
    const formValue = this.cardForm.value as CardFormData;

    if (this.data.mode === 'create') {
      this.createCard(formValue);
    } else {
      this.updateCard(formValue);
    }
  }

  private createCard(formData: CardFormData): void {
    const request: CreateCardRequest = {
      userId: formData.userId,
      cardType: formData.cardType,
      cardNumber: formData.cardNumber,
      holderName: formData.holderName,
      expirationMonth: formData.expirationMonth,
      expirationYear: formData.expirationYear,
      cvv: formData.cvv,
      nickname: formData.nickname,
      currency: formData.currency
    };

    this.cardService.createCard(request).subscribe({
      next: (card) => {
        this.saving.set(false);
        this.dialogRef.close(card);
      },
      error: (error) => {
        console.error('Error creating card:', error);
        this.saving.set(false);
        this.snackBar.open('Error al crear la tarjeta', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  private updateCard(formData: CardFormData): void {
    if (!this.data.card) return;

    const request: UpdateCardRequest = {
      holderName: formData.holderName,
      expirationMonth: formData.expirationMonth,
      expirationYear: formData.expirationYear,
      nickname: formData.nickname
    };

    this.cardService.updateCard(this.data.card.accountId || this.data.card.id, this.data.card.id, request).subscribe({
      next: (card) => {
        this.saving.set(false);
        this.dialogRef.close(card);
      },
      error: (error) => {
        console.error('Error updating card:', error);
        this.saving.set(false);
        this.snackBar.open('Error al actualizar la tarjeta', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  private markFormGroupTouched(): void {
    Object.keys(this.cardForm.controls).forEach(key => {
      const control = this.cardForm.get(key);
      control?.markAsTouched();
    });
  }

  onCancel(): void {
    this.dialogRef.close();
  }

  getFieldError(fieldName: string): string | null {
    const control = this.cardForm.get(fieldName);
    if (control?.invalid && control.touched) {
      const errors = control.errors;
      if (errors) {
        const errorKey = Object.keys(errors)[0];
        return errors[errorKey];
      }
    }
    return null;
  }

  getBrandIcon(): string {
    const brand = this.detectedBrand();
    return brand ? this.cardService.getCardBrandIcon(brand) : '';
  }

  getBrandName(): string {
    const brand = this.detectedBrand();
    if (!brand) return '';
    
    const names = {
      [CardBrand.VISA]: 'Visa',
      [CardBrand.MASTERCARD]: 'Mastercard',
      [CardBrand.AMERICAN_EXPRESS]: 'American Express',
      [CardBrand.DISCOVER]: 'Discover',
      [CardBrand.DINERS]: 'Diners Club',
      [CardBrand.OTHER]: 'Otra'
    };
    
    return names[brand];
  }
}