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
