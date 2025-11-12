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
  conversationId?: string;  // Nuevo: para continuidad conversacional
  period?: Period;          // Ahora opcional - se infiere automáticamente
  filters?: Record<string, any>;
}

export interface ChatMessage {
  id: string;
  userId: string;
  conversationId: string;
  role: 'user' | 'assistant';
  message: string;
  contextData?: Record<string, any>;
  createdAt: string;
}

export interface ChatHistoryResponse {
  conversationId: string;
  messages: ChatMessage[];
  total: number;
}

export interface ChatQueryResponse {
  reply: string;
  conversationId: string;        // Nuevo: ID de la conversación
  inferredPeriod?: string;       // Nuevo: período inferido
  inferredContext?: string;      // Nuevo: contexto inferido
  quickSuggestions?: string[];   // Nuevo: sugerencias rápidas
  suggestedActions?: any[];
  insights?: string[];
  dataRefs?: Record<string, any>;
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

  // Conversación actual en memoria
  private currentConversationId: string | null = null;

  query(req: ChatQueryRequest): Observable<ChatQueryResponse> {
    const headers: Record<string, string> = {};
    const user = this.auth.getCurrentUser();
    if (user?.id) headers['X-User-ID'] = user.id;
    
    // Agregar conversationId actual si existe
    const body = {
      ...req,
      conversationId: req.conversationId || this.currentConversationId || undefined
    };
    
    return this.http.post<ChatQueryResponse>(`${this.base}/query`, body, { headers });
  }

  // Nuevo: Obtener historial de conversación
  getHistory(conversationId: string): Observable<ChatHistoryResponse> {
    const headers: Record<string, string> = {};
    const user = this.auth.getCurrentUser();
    if (user?.id) headers['X-User-ID'] = user.id;
    
    return this.http.get<ChatHistoryResponse>(`${this.base}/history/${conversationId}`, { headers });
  }

  // Nuevo: Establecer conversación actual
  setCurrentConversation(conversationId: string | null): void {
    this.currentConversationId = conversationId;
  }

  // Nuevo: Obtener conversación actual
  getCurrentConversation(): string | null {
    return this.currentConversationId;
  }

  // Nuevo: Iniciar nueva conversación
  startNewConversation(): void {
    this.currentConversationId = null;
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