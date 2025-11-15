import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { EmailService } from '../../../services/email.service';
import { AuthService } from '../../../services/auth.service';

@Component({
  selector: 'app-support-dialog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './support-dialog.component.html',
  styleUrls: ['./support-dialog.component.css']
})
export class SupportDialogComponent {
  private readonly fb = inject(FormBuilder);
  private readonly emailService = inject(EmailService);
  private readonly authService = inject(AuthService);
  private readonly dialogRef = inject(MatDialogRef<SupportDialogComponent>);

  supportForm: FormGroup;
  sending = false;
  sent = false;
  error = false;
  errorMessage = '';

  constructor() {
    const currentUser = this.authService.currentUserSig();
    
    this.supportForm = this.fb.group({
      name: [
        `${currentUser?.firstName || ''} ${currentUser?.lastName || ''}`.trim() || '',
        [Validators.required, Validators.minLength(3)]
      ],
      email: [
        currentUser?.email || '',
        [Validators.required, Validators.email]
      ],
      subject: ['', [Validators.required, Validators.minLength(5)]],
      message: ['', [Validators.required, Validators.minLength(20)]]
    });
  }

  async onSubmit(): Promise<void> {
    if (this.supportForm.invalid) {
      this.supportForm.markAllAsTouched();
      return;
    }

    this.sending = true;
    this.error = false;
    this.errorMessage = '';

    const formValue = this.supportForm.value;
    
    this.emailService.sendSupportEmail({
      userName: formValue.name,
      userEmail: formValue.email,
      subject: formValue.subject,
      message: formValue.message
    }).subscribe({
      next: () => {
        console.log('✅ Email de soporte enviado exitosamente');
        this.sent = true;
        this.sending = false;
        
        // Cerrar el diálogo después de 2 segundos
        setTimeout(() => {
          this.dialogRef.close(true);
        }, 2000);
      },
      error: (error) => {
        console.error('❌ Error al enviar email de soporte:', error);
        this.sending = false;
        this.error = true;
        this.errorMessage = error.error?.details || 'No se pudo enviar el mensaje. Por favor, intenta nuevamente.';
      }
    });
  }

  close(): void {
    this.dialogRef.close(false);
  }

  getErrorMessage(fieldName: string): string {
    const field = this.supportForm.get(fieldName);
    
    if (!field || !field.errors || !field.touched) {
      return '';
    }

    if (field.errors['required']) {
      return 'Este campo es obligatorio';
    }

    if (field.errors['email']) {
      return 'Email inválido';
    }

    if (field.errors['minlength']) {
      const minLength = field.errors['minlength'].requiredLength;
      return `Debe tener al menos ${minLength} caracteres`;
    }

    return '';
  }
}
