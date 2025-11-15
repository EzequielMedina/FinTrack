import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators, AbstractControl } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { firstValueFrom } from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { TermsDialogComponent } from '../../components/terms-dialog/terms-dialog.component';
import { PrivacyDialogComponent } from '../../components/privacy-dialog/privacy-dialog.component';

// Custom validator for password confirmation
function passwordMatchValidator(control: AbstractControl): { [key: string]: any } | null {
  const password = control.get('password');
  const confirmPassword = control.get('confirmPassword');
  
  if (password && confirmPassword && password.value !== confirmPassword.value) {
    return { passwordMismatch: true };
  }
  return null;
}

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterLink,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCheckboxModule,
    MatDialogModule
  ],
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RegisterComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly dialog = inject(MatDialog);

  registerForm: FormGroup;
  loading = signal(false);
  error = signal('');
  hidePassword = signal(true);
  hideConfirmPassword = signal(true);

  constructor() {
    this.registerForm = this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [
        Validators.required, 
        Validators.minLength(8),
        Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&]).{8,}$/)
      ]],
      confirmPassword: ['', [Validators.required]],
      acceptTerms: [false, [Validators.requiredTrue]]
    }, { validators: passwordMatchValidator });
  }

  get firstNameControl() { return this.registerForm.get('firstName'); }
  get lastNameControl() { return this.registerForm.get('lastName'); }
  get emailControl() { return this.registerForm.get('email'); }
  get passwordControl() { return this.registerForm.get('password'); }
  get confirmPasswordControl() { return this.registerForm.get('confirmPassword'); }
  get acceptTermsControl() { return this.registerForm.get('acceptTerms'); }

  getFirstNameErrorMessage(): string {
    if (this.firstNameControl?.hasError('required')) {
      return 'El nombre es requerido';
    }
    if (this.firstNameControl?.hasError('minlength')) {
      return 'El nombre debe tener al menos 2 caracteres';
    }
    return '';
  }

  getLastNameErrorMessage(): string {
    if (this.lastNameControl?.hasError('required')) {
      return 'El apellido es requerido';
    }
    if (this.lastNameControl?.hasError('minlength')) {
      return 'El apellido debe tener al menos 2 caracteres';
    }
    return '';
  }

  getEmailErrorMessage(): string {
    if (this.emailControl?.hasError('required')) {
      return 'El email es requerido';
    }
    if (this.emailControl?.hasError('email')) {
      return 'Ingresa un email válido';
    }
    return '';
  }

  getPasswordErrorMessage(): string {
    if (this.passwordControl?.hasError('required')) {
      return 'La contraseña es requerida';
    }
    if (this.passwordControl?.hasError('minlength')) {
      return 'La contraseña debe tener al menos 8 caracteres';
    }
    if (this.passwordControl?.hasError('pattern')) {
      return 'Debe incluir: mayúscula, minúscula, número y símbolo (@$!%*?&)';
    }
    return '';
  }

  getConfirmPasswordErrorMessage(): string {
    if (this.confirmPasswordControl?.hasError('required')) {
      return 'Confirma tu contraseña';
    }
    if (this.registerForm.hasError('passwordMismatch') && this.confirmPasswordControl?.touched) {
      return 'Las contraseñas no coinciden';
    }
    return '';
  }

  togglePasswordVisibility(): void {
    this.hidePassword.set(!this.hidePassword());
  }

  toggleConfirmPasswordVisibility(): void {
    this.hideConfirmPassword.set(!this.hideConfirmPassword());
  }

  openTermsDialog(event: Event): void {
    event.preventDefault();
    this.dialog.open(TermsDialogComponent, {
      width: '800px',
      maxWidth: '95vw',
      maxHeight: '90vh',
      autoFocus: false,
      restoreFocus: false
    });
  }

  openPrivacyDialog(event: Event): void {
    event.preventDefault();
    this.dialog.open(PrivacyDialogComponent, {
      width: '800px',
      maxWidth: '95vw',
      maxHeight: '90vh',
      autoFocus: false,
      restoreFocus: false
    });
  }

  async onSubmit(): Promise<void> {
    if (this.registerForm.invalid) {
      this.registerForm.markAllAsTouched();
      return;
    }

    this.loading.set(true);
    this.error.set('');

    try {
      const { firstName, lastName, email, password } = this.registerForm.value;
      await firstValueFrom(this.auth.register({
        firstName,
        lastName,
        email,
        password
      }));
      this.router.navigateByUrl('/');
    } catch (e: any) {
      console.error('Register error:', e);
      if (e.status === 409) {
        this.error.set('Ya existe una cuenta con este email');
      } else if (e.status === 400) {
        this.error.set('Datos inválidos. Verifica la información ingresada');
      } else if (e.status === 0) {
        this.error.set('Error de conexión. Verifica tu conexión a internet');
      } else {
        this.error.set('Error inesperado. Intenta nuevamente');
      }
    } finally {
      this.loading.set(false);
    }
  }
}