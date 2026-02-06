import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ReportService } from '../../services/report.service';
import { AuthService } from '../../services/auth.service';
import { UserRole } from '../../models';

@Component({
  selector: 'app-reports',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './reports.component.html',
  styleUrls: ['./reports.component.css']
})
export class ReportsComponent implements OnInit {
  userId: string = '';
  isLoading: boolean = false;
  error: string = '';

  // Opciones de navegación de reportes
  reports = [
    {
      id: 'transactions',
      title: 'Transacciones',
      description: 'Análisis detallado de transacciones por período',
      icon: '📊',
      route: '/reports/transactions'
    },
    {
      id: 'installments',
      title: 'Cuotas',
      description: 'Estado de planes de cuotas y pagos pendientes',
      icon: '💳',
      route: '/reports/installments'
    },
    {
      id: 'accounts',
      title: 'Cuentas',
      description: 'Resumen de cuentas y tarjetas',
      icon: '🏦',
      route: '/reports/accounts'
    },
    {
      id: 'expenses-income',
      title: 'Gastos vs Ingresos',
      description: 'Análisis de flujo de efectivo y tendencias',
      icon: '💰',
      route: '/reports/expenses-income'
    }
    // Reporte de notificaciones temporalmente deshabilitado
    // {
    //   id: 'notifications',
    //   title: 'Notificaciones',
    //   description: 'Estadísticas de notificaciones del sistema',
    //   icon: '🔔',
    //   route: '/reports/notifications',
    //   adminOnly: true
    // }
  ];

  constructor(
    private reportService: ReportService,
    private authService: AuthService,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Obtener el ID del usuario actual
    const currentUser = this.authService.getCurrentUser();
    if (currentUser) {
      this.userId = currentUser.id;
    } else {
      this.router.navigate(['/login']);
    }
  }

  navigateToReport(route: string): void {
    this.router.navigate([route]);
  }

  isAdmin(): boolean {
    const currentUser = this.authService.getCurrentUser();
    return currentUser?.role === UserRole.ADMIN;
  }

  getVisibleReports() {
    return this.reports.filter(report => !(report as any).adminOnly || this.isAdmin());
  }

  getReportIcon(emoji: string): string {
    const iconMap: { [key: string]: string } = {
      '📊': 'assets/icons/chart-up.svg',
      '💳': 'assets/icons/card.svg',
      '🏦': 'assets/icons/account.svg',
      '💰': 'assets/icons/wallet.svg',
      '🔔': 'assets/icons/alert-circle.svg'
    };
    return iconMap[emoji] || 'assets/icons/report.svg';
  }
}
