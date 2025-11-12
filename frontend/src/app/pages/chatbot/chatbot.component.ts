import { Component, inject, OnInit, ViewChild, ElementRef, AfterViewChecked } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ChatbotService, ChatMessage, ChatQueryResponse } from '../../services/chatbot.service';

interface ChatBubble {
  role: 'user' | 'assistant';
  message: string;
  timestamp: Date;
  inferredContext?: string;
  inferredPeriod?: string;
  quickSuggestions?: string[];
}

@Component({
  selector: 'app-chatbot',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatTooltipModule
  ],
  templateUrl: './chatbot.component.html',
  styleUrls: ['./chatbot.component.css']
})
export class ChatbotComponent implements OnInit, AfterViewChecked {
  private readonly api = inject(ChatbotService);

  @ViewChild('chatContainer') private chatContainer!: ElementRef;
  
  // Estado del chat
  messages: ChatBubble[] = [];
  currentMessage = '';
  loading = false;
  error?: string;
  
  // Conversación actual
  conversationId?: string;
  
  // Auto-scroll flag
  private shouldScrollToBottom = false;

  ngOnInit(): void {
    // Iniciar con mensaje de bienvenida
    this.messages.push({
      role: 'assistant',
      message: '¡Hola! Soy tu asistente financiero. Pregúntame sobre tus gastos, ingresos, tarjetas o cuotas. Por ejemplo:\n\n- "¿Cuánto gasté hoy?"\n- "Muéstrame mis tarjetas"\n- "Estado de cuotas esta semana"',
      timestamp: new Date()
    });
  }

  ngAfterViewChecked(): void {
    if (this.shouldScrollToBottom) {
      this.scrollToBottom();
      this.shouldScrollToBottom = false;
    }
  }

  sendMessage(): void {
    if (!this.currentMessage.trim() || this.loading) return;

    // Agregar mensaje del usuario
    const userMessage: ChatBubble = {
      role: 'user',
      message: this.currentMessage,
      timestamp: new Date()
    };
    this.messages.push(userMessage);
    this.shouldScrollToBottom = true;

    const messageToSend = this.currentMessage;
    this.currentMessage = '';
    this.loading = true;
    this.error = undefined;

    // Enviar al backend
    this.api.query({ 
      message: messageToSend,
      conversationId: this.conversationId 
    }).subscribe({
      next: (response: ChatQueryResponse) => {
        // Guardar conversationId para continuidad
        if (response.conversationId) {
          this.conversationId = response.conversationId;
          this.api.setCurrentConversation(response.conversationId);
        }

        // Agregar respuesta del asistente
        const assistantMessage: ChatBubble = {
          role: 'assistant',
          message: response.reply,
          timestamp: new Date(),
          inferredContext: response.inferredContext,
          inferredPeriod: response.inferredPeriod,
          quickSuggestions: response.quickSuggestions
        };
        this.messages.push(assistantMessage);
        this.shouldScrollToBottom = true;
        this.loading = false;
      },
      error: (err) => {
        this.error = err?.error?.error || err?.message || 'Error al consultar el chatbot';
        this.loading = false;
        
        // Agregar mensaje de error como asistente
        this.messages.push({
          role: 'assistant',
          message: `Lo siento, ocurrió un error: ${this.error}`,
          timestamp: new Date()
        });
        this.shouldScrollToBottom = true;
      }
    });
  }

  // Usar una sugerencia rápida
  useSuggestion(suggestion: string): void {
    this.currentMessage = suggestion;
    this.sendMessage();
  }

  // Iniciar nueva conversación
  startNewChat(): void {
    if (confirm('¿Estás seguro de que deseas iniciar una nueva conversación?')) {
      this.conversationId = undefined;
      this.api.startNewConversation();
      this.messages = [{
        role: 'assistant',
        message: '¡Nueva conversación iniciada! ¿En qué puedo ayudarte?',
        timestamp: new Date()
      }];
      this.currentMessage = '';
      this.error = undefined;
    }
  }

  // Scroll al final del chat
  private scrollToBottom(): void {
    try {
      if (this.chatContainer) {
        this.chatContainer.nativeElement.scrollTop = this.chatContainer.nativeElement.scrollHeight;
      }
    } catch (err) {
      console.error('Error al hacer scroll:', err);
    }
  }

  // Formatear timestamp
  formatTime(date: Date): string {
    return new Intl.DateTimeFormat('es-AR', {
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  }

  // Manejar Enter para enviar
  onKeyDown(event: KeyboardEvent): void {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      this.sendMessage();
    }
  }
}