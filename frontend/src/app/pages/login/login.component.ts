import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <section class="login">
      <h2>Iniciar sesión</h2>
      <form (ngSubmit)="onSubmit()" #f="ngForm">
        <label>
          Email
          <input name="email" required [(ngModel)]="email" type="email" />
        </label>
        <label>
          Password
          <input name="password" required [(ngModel)]="password" type="password" />
        </label>
        <button [disabled]="loading" type="submit">Entrar</button>
      </form>
      <p class="error" *ngIf="error">{{ error }}</p>
    </section>
  `,
  styles: [
    `
    .login { max-width: 360px; margin: 3rem auto; display:flex; flex-direction:column; gap:.75rem; }
    label { display:flex; flex-direction:column; gap:.25rem; }
    input { padding:.5rem; border:1px solid #ccc; border-radius:4px; }
    button { padding:.5rem .75rem; }
    .error { color: #b00020; }
    `
  ],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LoginComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  email = '';
  password = '';
  loading = false;
  error = '';

  async onSubmit(): Promise<void> {
    this.loading = true;
    this.error = '';
    try {
      await firstValueFrom(this.auth.login(this.email, this.password));
      this.router.navigateByUrl('/');
    } catch (e: unknown) {
      this.error = 'Credenciales inválidas';
    } finally {
      this.loading = false;
    }
  }
}
