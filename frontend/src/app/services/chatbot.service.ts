import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from './auth.service';

export interface Period {
  from?: string;
  to?: string;
}

export interface ChatQueryRequest {
  message: string;
  period?: Period;
  filters?: Record<string, any>;
}

export interface ReportRequest {
  title?: string;
  period?: Period;
  groupBy?: 'type' | 'merchant' | 'account' | 'card';
  includeCharts?: boolean;
  filters?: Record<string, any>;
}

@Injectable({ providedIn: 'root' })
export class ChatbotService {
  private readonly base = '/api/chat';
  private readonly http = inject(HttpClient);
  private readonly auth = inject(AuthService);

  query(req: ChatQueryRequest): Observable<any> {
    const headers: Record<string, string> = {};
    const user = this.auth.getCurrentUser();
    if (user?.id) headers['X-User-ID'] = user.id;
    return this.http.post(`${this.base}/query`, req, { headers });
  }

  reportPdf(req: ReportRequest): Observable<Blob> {
    const headers: Record<string, string> = {};
    const user = this.auth.getCurrentUser();
    if (user?.id) headers['X-User-ID'] = user.id;
    return this.http.post(`${this.base}/report/pdf`, req, { responseType: 'blob', headers });
  }

  reportChart(req: ReportRequest): Observable<any> {
    const headers: Record<string, string> = {};
    const user = this.auth.getCurrentUser();
    if (user?.id) headers['X-User-ID'] = user.id;
    return this.http.post(`${this.base}/report/chart`, req, { headers });
  }
}