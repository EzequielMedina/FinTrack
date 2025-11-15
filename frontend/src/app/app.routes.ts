import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';
import { guestGuard } from './guards/guest.guard';
import { adminPanelGuard } from './guards/permission.guard';
import { roleGuard } from './guards/role.guard';
import { UserRole } from './models';

export const routes: Routes = [
  {
    path: 'login',
    canActivate: [guestGuard],
    loadComponent: () =>
      import('./pages/login/login.component').then((m) => m.LoginComponent)
  },
  {
    path: 'register',
    canActivate: [guestGuard],
    loadComponent: () =>
      import('./pages/register/register.component').then((m) => m.RegisterComponent)
  },
  {
    path: 'dashboard',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/dashboard/dashboard.component').then(
        (m) => m.DashboardComponent
      )
  },
  {
    path: 'cards',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/cards/cards.component').then(
        (m) => m.CardsComponent
      )
  },
  {
    path: 'accounts',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/accounts/accounts.component').then(
        (m) => m.AccountsComponent
      )
  },
  {
    path: 'transactions',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/transactions/transactions.component').then(
        (m) => m.TransactionsComponent
      )
  },
  {
    path: 'chatbot',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/chatbot/chatbot.component').then(
        (m) => m.ChatbotComponent
      )
  },
  {
    path: 'faq',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/faq/faq.component').then(
        (m) => m.FaqComponent
      )
  },
  {
    path: 'reports',
    canActivate: [authGuard],
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./pages/reports/reports.component').then(
            (m) => m.ReportsComponent
          )
      },
      {
        path: 'transactions',
        loadComponent: () =>
          import('./pages/reports/transaction-report/transaction-report.component').then(
            (m) => m.TransactionReportComponent
          )
      },
      {
        path: 'installments',
        loadComponent: () =>
          import('./pages/reports/installment-report/installment-report.component').then(
            (m) => m.InstallmentReportComponent
          )
      },
      {
        path: 'accounts',
        loadComponent: () =>
          import('./pages/reports/account-report/account-report.component').then(
            (m) => m.AccountReportComponent
          )
      },
      {
        path: 'expenses-income',
        loadComponent: () =>
          import('./pages/reports/expense-income-report/expense-income-report.component').then(
            (m) => m.ExpenseIncomeReportComponent
          )
      },
      {
        path: 'notifications',
        loadComponent: () =>
          import('./pages/reports/notification-report/notification-report.component').then(
            (m) => m.NotificationReportComponent
          )
      }
    ]
  },
  {
    path: 'admin',
    canActivate: [adminPanelGuard],
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./pages/admin/admin-panel.component').then(
            (m) => m.AdminPanelComponent
          )
      },
      {
        path: 'users',
        loadComponent: () =>
          import('./pages/admin/user-management/user-management.component').then(
            (m) => m.UserManagementComponent
          )
      },
      {
        path: 'roles',
        loadComponent: () =>
          import('./pages/admin/admin-panel.component').then(
            (m) => m.AdminPanelComponent
          )
      },
      {
        path: 'reports',
        loadComponent: () =>
          import('./pages/admin/admin-panel.component').then(
            (m) => m.AdminPanelComponent
          )
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./pages/admin/admin-panel.component').then(
            (m) => m.AdminPanelComponent
          )
      }
    ]
  },
  {
    path: '',
    redirectTo: '/dashboard',
    pathMatch: 'full'
  },
  {
    path: '**',
    loadComponent: () =>
      import('./pages/not-found/not-found.component').then(
        (m) => m.NotFoundComponent
      )
  }
];
