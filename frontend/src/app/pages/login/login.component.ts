import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { firstValueFrom } from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule, MatDialog } from '@angular/material/dialog';
import { MatToolbarModule } from '@angular/material/toolbar';

@Component({
  selector: 'app-login',
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
    MatDialogModule,
    MatToolbarModule
  ],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LoginComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly dialog = inject(MatDialog);

  loginForm: FormGroup;
  loading = signal(false);
  error = signal('');
  hidePassword = signal(true);

  constructor() {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  get emailControl() {
    return this.loginForm.get('email');
  }

  get passwordControl() {
    return this.loginForm.get('password');
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
      return 'La contraseña debe tener al menos 6 caracteres';
    }
    return '';
  }

  togglePasswordVisibility(): void {
    this.hidePassword.set(!this.hidePassword());
  }

  openAboutDialog(): void {
    this.dialog.open(AboutDialogComponent, {
      width: '800px',
      maxWidth: '90vw',
      maxHeight: '90vh',
      panelClass: 'about-dialog'
    });
  }

  async onSubmit(): Promise<void> {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.loading.set(true);
    this.error.set('');

    try {
      const { email, password } = this.loginForm.value;
      await firstValueFrom(this.auth.login(email, password));
      this.router.navigateByUrl('/');
    } catch (e: any) {
      console.error('Login error:', e);
      if (e.status === 401) {
        this.error.set('Email o contraseña incorrectos');
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

// About Dialog Component
@Component({
  selector: 'app-about-dialog',
  standalone: true,
  imports: [CommonModule, MatDialogModule, MatButtonModule, MatIconModule],
  template: `
    <div class="about-dialog-content">
      <div class="dialog-header">
        <h1 mat-dialog-title>
          <mat-icon class="dialog-logo">account_balance_wallet</mat-icon>
          FinTrack
        </h1>
        <button mat-icon-button mat-dialog-close class="close-button">
          <mat-icon>close</mat-icon>
        </button>
      </div>

      <mat-dialog-content>
        <div class="about-card">
          <h2>🎯 ¿Qué Problema Solucionamos?</h2>
          <p>
            Muchas personas y familias enfrentan dificultades para <strong>gestionar sus finanzas personales</strong>:
          </p>
          <ul>
            <li>📊 <strong>Falta de visibilidad</strong>: No saben en qué gastan su dinero</li>
            <li>💳 <strong>Control de deudas</strong>: Pierden el rastro de pagos de tarjetas y cuotas</li>
            <li>💰 <strong>Sin presupuesto</strong>: Gastan más de lo que ganan sin darse cuenta</li>
            <li>📈 <strong>Objetivos difusos</strong>: No pueden ahorrar para metas específicas</li>
            <li>🔄 <strong>Múltiples cuentas</strong>: Información dispersa en varios bancos y apps</li>
          </ul>
          <p>
            El resultado: <strong>estrés financiero, deudas acumuladas y objetivos no cumplidos</strong>.
          </p>
        </div>

        <div class="about-card">
          <h2>🚀 Nuestra Solución: FinTrack</h2>
          <p>
            FinTrack es una <strong>plataforma integral de gestión financiera personal</strong> que centraliza 
            todas tus finanzas en un solo lugar:
          </p>
          <div class="features-grid">
            <div class="feature">
              <span class="feature-icon">📱</span>
              <h3>Control Total</h3>
              <p>Gestiona cuentas, tarjetas y transacciones en tiempo real</p>
            </div>
            <div class="feature">
              <span class="feature-icon">💡</span>
              <h3>Asistente IA</h3>
              <p>Chatbot inteligente que responde tus dudas financieras</p>
            </div>
            <div class="feature">
              <span class="feature-icon">📊</span>
              <h3>Reportes Visuales</h3>
              <p>Gráficos y análisis de tus gastos e ingresos</p>
            </div>
            <div class="feature">
              <span class="feature-icon">🎯</span>
              <h3>Metas y Ahorro</h3>
              <p>Define objetivos y monitorea tu progreso</p>
            </div>
            <div class="feature">
              <span class="feature-icon">🔔</span>
              <h3>Notificaciones</h3>
              <p>Alertas de vencimientos y recordatorios de pagos</p>
            </div>
            <div class="feature">
              <span class="feature-icon">🌐</span>
              <h3>Multi-divisa</h3>
              <p>Soporte para múltiples monedas con conversión automática</p>
            </div>
          </div>
        </div>

        <div class="about-card">
          <h2>👥 Quiénes Somos</h2>
          <p>
            Somos un equipo de <strong>desarrolladores y especialistas en fintech</strong> comprometidos con 
            democratizar el acceso a herramientas de gestión financiera profesional.
          </p>
          <p>
            <strong>Nuestra misión</strong>: Empoderar a las personas para que tomen el control de sus finanzas 
            y alcancen sus objetivos económicos mediante tecnología accesible e inteligente.
          </p>
          <p class="tech-stack">
            <strong>Tecnología:</strong> Construido con Angular, Go, MySQL y Docker. 
            Arquitectura de microservicios escalable y segura.
          </p>
        </div>
      </mat-dialog-content>

      <mat-dialog-actions align="end">
        <button mat-raised-button color="primary" mat-dialog-close>Cerrar</button>
      </mat-dialog-actions>
    </div>
  `,
  styles: [`
    .about-dialog-content {
      padding: 0;
    }

    .dialog-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px 24px;
      border-bottom: 1px solid #e0e0e0;
    }

    h1[mat-dialog-title] {
      display: flex;
      align-items: center;
      gap: 12px;
      margin: 0;
      font-size: 24px;
      font-weight: 700;
      color: #1a1a1a;
    }

    .dialog-logo {
      width: 32px;
      height: 32px;
      font-size: 32px;
      color: var(--accent-600);
    }

    .close-button {
      margin-right: -8px;
    }

    mat-dialog-content {
      padding: 24px;
      max-height: 70vh;
      overflow-y: auto;
    }

    .about-card {
      background: #f8f9fa;
      border-radius: 12px;
      padding: 24px;
      margin-bottom: 20px;
    }

    .about-card:last-child {
      margin-bottom: 0;
    }

    .about-card h2 {
      font-size: 20px;
      font-weight: 700;
      color: #1a1a1a;
      margin-top: 0;
      margin-bottom: 16px;
    }

    .about-card p {
      font-size: 15px;
      line-height: 1.6;
      color: #4a4a4a;
      margin-bottom: 12px;
    }

    .about-card ul {
      list-style: none;
      padding: 0;
      margin: 16px 0;
    }

    .about-card li {
      font-size: 15px;
      line-height: 1.6;
      color: #4a4a4a;
      padding: 8px 0;
      padding-left: 8px;
    }

    .features-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 16px;
      margin-top: 20px;
    }

    .feature {
      background: white;
      padding: 16px;
      border-radius: 8px;
      text-align: center;
      transition: transform 0.2s;
    }

    .feature:hover {
      transform: translateY(-2px);
    }

    .feature-icon {
      font-size: 32px;
      display: block;
      margin-bottom: 12px;
    }

    .feature h3 {
      font-size: 16px;
      font-weight: 600;
      color: #1a1a1a;
      margin: 0 0 8px 0;
    }

    .feature p {
      font-size: 14px;
      color: #6b6b6b;
      margin: 0;
    }

    .tech-stack {
      background: white;
      padding: 12px 16px;
      border-radius: 8px;
      border-left: 4px solid #3b82f6;
      font-size: 14px;
      margin-top: 16px;
    }

    mat-dialog-actions {
      padding: 16px 24px;
      border-top: 1px solid #e0e0e0;
    }

    @media (max-width: 768px) {
      .features-grid {
        grid-template-columns: 1fr;
      }

      .dialog-header h1 {
        font-size: 20px;
      }

      mat-dialog-content {
        padding: 16px;
      }

      .about-card {
        padding: 16px;
      }
    }
  `]
})
export class AboutDialogComponent {}
