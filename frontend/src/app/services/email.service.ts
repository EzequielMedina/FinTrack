import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export interface SupportEmailData {
  userName: string;
  userEmail: string;
  subject: string;
  message: string;
}

interface SupportEmailRequest {
  name: string;
  email: string;
  subject: string;
  message: string;
}

@Injectable({
  providedIn: 'root'
})
export class EmailService {
  private readonly http = inject(HttpClient);
  private readonly notificationServiceUrl = `${environment.apiUrl}/notifications`;

  /**
   * Envía un email de soporte a través del microservicio de notificaciones
   * Usa el template template_yst8bd2 configurado en el backend
   * @param data Datos del formulario de soporte
   * @returns Observable con el resultado del envío
   */
  sendSupportEmail(data: SupportEmailData): Observable<any> {
    const request: SupportEmailRequest = {
      name: data.userName,
      email: data.userEmail,
      subject: data.subject,
      message: data.message
    };

    return this.http.post(`${this.notificationServiceUrl}/support`, request);
  }
}
