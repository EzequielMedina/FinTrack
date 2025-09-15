import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatGridListModule } from '@angular/material/grid-list';

import { AuthService } from '../../services/auth.service';
import { PermissionService } from '../../services/permission.service';
import { Permission } from '../../models';

@Component({
  selector: 'app-admin-panel',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatGridListModule
  ],
  template: `
    <div class="admin-panel-container">
      <div class="header-section">
        <h1>
          <mat-icon>admin_panel_settings</mat-icon>
          Panel de Administración
        </h1>
        <p>Gestiona usuarios, roles y permisos del sistema FinTrack</p>
      </div>

      <div class="admin-cards-grid">
        <!-- User Management Card -->
        <mat-card class="admin-card" *ngIf="canManageUsers()">
          <mat-card-header>
            <mat-icon mat-card-avatar>people</mat-icon>
            <mat-card-title>Gestión de Usuarios</mat-card-title>
            <mat-card-subtitle>Administrar usuarios y roles</mat-card-subtitle>
          </mat-card-header>
          <mat-card-content>
            <p>Crear, editar y eliminar usuarios. Asignar roles y gestionar permisos.</p>
            <div class="card-stats">
              <span class="stat">
                <mat-icon>person</mat-icon>
                Total de usuarios
              </span>
            </div>
          </mat-card-content>
          <mat-card-actions>
            <button mat-raised-button color="primary" routerLink="/admin/users">
              <mat-icon>settings</mat-icon>
              Gestionar Usuarios
            </button>
          </mat-card-actions>
        </mat-card>

        <!-- Role Management Card -->
        <mat-card class="admin-card" *ngIf="canUpdateRoles()">
          <mat-card-header>
            <mat-icon mat-card-avatar>security</mat-icon>
            <mat-card-title>Gestión de Roles</mat-card-title>
            <mat-card-subtitle>Configurar roles y permisos</mat-card-subtitle>
          </mat-card-header>
          <mat-card-content>
            <p>Definir roles personalizados y asignar permisos específicos a cada rol.</p>
            <div class="card-stats">
              <span class="stat">
                <mat-icon>verified_user</mat-icon>
                Roles activos
              </span>
            </div>
          </mat-card-content>
          <mat-card-actions>
            <button mat-raised-button color="primary" routerLink="/admin/roles">
              <mat-icon>edit</mat-icon>
              Configurar Roles
            </button>
          </mat-card-actions>
        </mat-card>

        <!-- Reports Card -->
        <mat-card class="admin-card" *ngIf="canViewReports()">
          <mat-card-header>
            <mat-icon mat-card-avatar>assessment</mat-icon>
            <mat-card-title>Reportes</mat-card-title>
            <mat-card-subtitle>Análisis y estadísticas</mat-card-subtitle>
          </mat-card-header>
          <mat-card-content>
            <p>Ver reportes de actividad de usuarios, uso del sistema y estadísticas.</p>
            <div class="card-stats">
              <span class="stat">
                <mat-icon>trending_up</mat-icon>
                Actividad del sistema
              </span>
            </div>
          </mat-card-content>
          <mat-card-actions>
            <button mat-raised-button color="primary" routerLink="/admin/reports">
              <mat-icon>bar_chart</mat-icon>
              Ver Reportes
            </button>
          </mat-card-actions>
        </mat-card>

        <!-- System Settings Card -->
        <mat-card class="admin-card">
          <mat-card-header>
            <mat-icon mat-card-avatar>settings</mat-icon>
            <mat-card-title>Configuración</mat-card-title>
            <mat-card-subtitle>Ajustes del sistema</mat-card-subtitle>
          </mat-card-header>
          <mat-card-content>
            <p>Configurar parámetros del sistema, notificaciones y preferencias globales.</p>
            <div class="card-stats">
              <span class="stat">
                <mat-icon>tune</mat-icon>
                Configuraciones
              </span>
            </div>
          </mat-card-content>
          <mat-card-actions>
            <button mat-raised-button color="primary" routerLink="/admin/settings">
              <mat-icon>settings</mat-icon>
              Configurar Sistema
            </button>
          </mat-card-actions>
        </mat-card>
      </div>

      <!-- Quick Actions Section -->
      <div class="quick-actions-section">
        <h2>Acciones Rápidas</h2>
        <div class="quick-actions">
          <button mat-stroked-button 
                  *ngIf="canManageUsers()" 
                  routerLink="/admin/users"
                  [queryParams]="{action: 'create'}">
            <mat-icon>person_add</mat-icon>
            Crear Usuario
          </button>
          
          <button mat-stroked-button 
                  *ngIf="canViewReports()"
                  routerLink="/admin/reports"
                  [queryParams]="{type: 'activity'}">
            <mat-icon>analytics</mat-icon>
            Ver Actividad
          </button>
          
          <button mat-stroked-button routerLink="/dashboard">
            <mat-icon>dashboard</mat-icon>
            Volver al Dashboard
          </button>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .admin-panel-container {
      padding: 20px;
      max-width: 1200px;
      margin: 0 auto;
    }

    .header-section {
      text-align: center;
      margin-bottom: 40px;
    }

    .header-section h1 {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 12px;
      margin-bottom: 8px;
      color: #1976d2;
    }

    .header-section h1 mat-icon {
      font-size: 32px;
      width: 32px;
      height: 32px;
    }

    .header-section p {
      color: #666;
      font-size: 16px;
    }

    .admin-cards-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 20px;
      margin-bottom: 40px;
    }

    .admin-card {
      transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
    }

    .admin-card:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    }

    .admin-card mat-card-header mat-icon {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border-radius: 50%;
      padding: 8px;
    }

    .card-stats {
      margin-top: 16px;
    }

    .stat {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #666;
      font-size: 14px;
    }

    .stat mat-icon {
      font-size: 16px;
      width: 16px;
      height: 16px;
    }

    .quick-actions-section {
      background: #f5f5f5;
      padding: 20px;
      border-radius: 8px;
    }

    .quick-actions-section h2 {
      margin-bottom: 16px;
      color: #333;
    }

    .quick-actions {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
    }

    .quick-actions button {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    /* Responsive design */
    @media (max-width: 768px) {
      .admin-panel-container {
        padding: 10px;
      }

      .admin-cards-grid {
        grid-template-columns: 1fr;
      }

      .quick-actions {
        flex-direction: column;
      }

      .quick-actions button {
        width: 100%;
      }
    }
  `]
})
export class AdminPanelComponent {
  private readonly authService = inject(AuthService);
  private readonly permissionService = inject(PermissionService);

  // Permission checks for template
  canManageUsers(): boolean {
    return this.permissionService.canManageUsers();
  }

  canUpdateRoles(): boolean {
    return this.permissionService.canUpdateUserRoles();
  }

  canViewReports(): boolean {
    return this.permissionService.canViewReports();
  }
}