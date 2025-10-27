import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatTableModule } from '@angular/material/table';
import { formatDate } from '@angular/common';
import { ChatbotService, ChatQueryRequest, ReportRequest } from '../../services/chatbot.service';

interface SuggestedAction { type: string; params?: any }

interface QueryTemplate {
  label: string;
  message: string;
  description: string;
  period?: string;
  contextType: string;
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
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatTableModule
  ],
  templateUrl: './chatbot.component.html',
  styleUrls: ['./chatbot.component.css']
})
export class ChatbotComponent {
  private readonly api = inject(ChatbotService);

  // Form fields
  message = '';
  periodFrom: Date | null = null;
  periodTo: Date | null = null;
  selectedType: 'accounts' | 'cards' | 'both' = 'both';
  
  // New enhanced fields
  queryType: string = 'custom';
  quickPeriod: string = 'today';
  contextFocus: string = 'general';

  // Results
  loading = false;
  answer?: string;
  suggestions?: SuggestedAction[];
  chartData?: any;
  error?: string;

  // Query templates for common questions
  queryTemplates: QueryTemplate[] = [
    {
      label: 'Gastos de hoy',
      message: 'decime los gastos de hoy',
      description: 'Ver todos los gastos del día actual',
      period: 'today',
      contextType: 'expenses'
    },
    {
      label: 'Gastos del mes',
      message: 'decime los gastos de este mes',
      description: 'Resumen de gastos mensuales',
      period: 'month',
      contextType: 'expenses'
    },
    {
      label: 'Estado de tarjetas',
      message: 'como están mis tarjetas de crédito',
      description: 'Saldos y estado actual de tarjetas',
      period: 'month',
      contextType: 'cards'
    },
    {
      label: 'Análisis de cuotas',
      message: 'decime sobre mis planes de cuotas',
      description: 'Estado de planes e installments',
      period: 'all',
      contextType: 'installments'
    },
    {
      label: 'Top comercios',
      message: 'en qué comercios gasté más este mes',
      description: 'Ranking de gastos por comercio',
      period: 'month',
      contextType: 'merchants'
    }
  ];

  // Period options
  periodOptions = [
    { value: 'today', label: 'Hoy' },
    { value: 'week', label: 'Esta semana' },
    { value: 'month', label: 'Este mes' },
    { value: 'custom', label: 'Período personalizado' }
  ];

  // Context focus options
  contextOptions = [
    { value: 'general', label: 'Información general' },
    { value: 'expenses', label: 'Enfoque en gastos' },
    { value: 'income', label: 'Enfoque en ingresos' },
    { value: 'cards', label: 'Enfoque en tarjetas' },
    { value: 'installments', label: 'Enfoque en cuotas' },
    { value: 'merchants', label: 'Enfoque en comercios' }
  ];

  runQuery(): void {
    this.error = undefined;
    this.answer = undefined;
    this.suggestions = undefined;
    this.loading = true;

    const req: ChatQueryRequest = {
      message: this.message,
      period: this.buildPeriod(),
      filters: this.buildEnhancedFilters()
    };

    this.api.query(req).subscribe({
      next: (res: any) => {
        this.answer = res?.reply ?? 'Sin respuesta';
        this.suggestions = res?.suggestedActions ?? [];
        this.loading = false;
      },
      error: (err) => {
        this.error = (err?.error?.error) || (err?.message) || 'Error al consultar el chatbot';
        this.loading = false;
      }
    });
  }

  // New method: Use a query template
  useTemplate(template: QueryTemplate): void {
    this.message = template.message;
    this.contextFocus = template.contextType;
    
    // Set period based on template
    if (template.period) {
      this.quickPeriod = template.period;
      this.setPeriodFromQuick();
    }
    
    // Run query immediately
    setTimeout(() => this.runQuery(), 100);
  }

  // New method: Set period from quick selector
  setPeriodFromQuick(): void {
    const today = new Date();
    
    switch (this.quickPeriod) {
      case 'today':
        this.periodFrom = today;
        this.periodTo = today;
        break;
      case 'week':
        const startOfWeek = new Date(today);
        startOfWeek.setDate(today.getDate() - today.getDay());
        this.periodFrom = startOfWeek;
        this.periodTo = today;
        break;
      case 'month':
        const startOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
        this.periodFrom = startOfMonth;
        this.periodTo = today;
        break;
      case 'custom':
        // Keep current dates
        break;
    }
  }

  // Enhanced filters building
  private buildEnhancedFilters() {
    return {
      type: this.selectedType,
      contextFocus: this.contextFocus,
      quickPeriod: this.quickPeriod
    };
  }

  downloadPdf(): void {
    this.error = undefined;
    const req: ReportRequest = {
      title: 'Reporte Chatbot',
      period: this.buildPeriod(),
      groupBy: 'type', // default value
      includeCharts: true,
      filters: this.buildEnhancedFilters()
    };
    this.api.reportPdf(req).subscribe({
      next: (blob) => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'reporte_chatbot.pdf';
        a.click();
        window.URL.revokeObjectURL(url);
      },
      error: (err) => {
        this.error = (err?.error?.message) || 'Error al generar PDF';
      }
    });
  }

  loadChart(): void {
    this.error = undefined;
    const req: ReportRequest = {
      period: this.buildPeriod(),
      groupBy: 'type', // default value
      includeCharts: true,
      filters: this.buildEnhancedFilters()
    };
    this.api.reportChart(req).subscribe({
      next: (res) => {
        this.chartData = res;
      },
      error: (err) => {
        this.error = (err?.error?.message) || 'Error al cargar datos del gráfico';
      }
    });
  }

  runSuggestion(s: SuggestedAction): void {
    if (!s) return;
    const t = s.type?.toLowerCase();
    switch (t) {
      case 'generate_pdf':
        this.downloadPdf();
        break;
      case 'show_chart':
        this.loadChart();
        break;
      default:
        // No-op for unknown suggestions
        break;
    }
  }

  private buildPeriod() {
    const from = this.periodFrom ? formatDate(this.periodFrom, 'yyyy-MM-dd', 'en-US') : undefined;
    const to = this.periodTo ? formatDate(this.periodTo, 'yyyy-MM-dd', 'en-US') : undefined;
    return (from || to) ? { from, to } : undefined;
  }
}