import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatDialogModule, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { MatSnackBarModule, MatSnackBar } from '@angular/material/snack-bar';
import { MatChipsModule } from '@angular/material/chips';
import { MatCardModule } from '@angular/material/card';
import { MatTooltipModule } from '@angular/material/tooltip';

import { UserService } from '../../../services/user.service';
import { PermissionService } from '../../../services/permission.service';
import { User, UserRole, UserStatus, UsersListResponse, Permission } from '../../../models';

@Component({
  selector: 'app-user-management',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MatTableModule,
    MatPaginatorModule,
    MatButtonModule,
    MatIconModule,
    MatInputModule,
    MatSelectModule,
    MatDialogModule,
    MatSnackBarModule,
    MatChipsModule,
    MatCardModule,
    MatTooltipModule
  ],
  template: `
    <div class="user-management-container">
      <mat-card class="management-card">
        <mat-card-header>
          <mat-card-title>
            <mat-icon>people</mat-icon>
            Gestión de Usuarios
          </mat-card-title>
          <mat-card-subtitle>
            Administra usuarios, roles y permisos del sistema
          </mat-card-subtitle>
        </mat-card-header>

        <mat-card-content>
          <!-- Filtros y búsqueda -->
          <div class="filters-section">
            <mat-form-field class="search-field">
              <mat-label>Buscar usuarios</mat-label>
              <input matInput 
                     [(ngModel)]="searchQuery" 
                     (keyup.enter)="searchUsers()"
                     placeholder="Nombre, email...">
              <mat-icon matSuffix>search</mat-icon>
            </mat-form-field>

            <mat-form-field class="role-filter">
              <mat-label>Filtrar por rol</mat-label>
              <mat-select [(ngModel)]="selectedRole" (selectionChange)="filterByRole()">
                <mat-option value="">Todos los roles</mat-option>
                <!-- Admin filter is excluded - admins cannot see other admins -->
                <mat-option [value]="UserRole.USER">Usuario</mat-option>
                <mat-option [value]="UserRole.OPERATOR">Operador</mat-option>
                <mat-option [value]="UserRole.TREASURER">Tesorero</mat-option>
              </mat-select>
            </mat-form-field>

            <button mat-raised-button 
                    color="primary" 
                    (click)="openCreateUserDialog()"
                    *ngIf="canCreateUsers()">
              <mat-icon>person_add</mat-icon>
              Crear Usuario
            </button>
          </div>

          <!-- Tabla de usuarios -->
          <div class="table-container">
            <table mat-table [dataSource]="users()" class="users-table">
              <!-- Avatar/Initials Column -->
              <ng-container matColumnDef="avatar">
                <th mat-header-cell *matHeaderCellDef>Avatar</th>
                <td mat-cell *matCellDef="let user">
                  <div class="user-avatar">
                    {{ getUserInitials(user) }}
                  </div>
                </td>
              </ng-container>

              <!-- Name Column -->
              <ng-container matColumnDef="name">
                <th mat-header-cell *matHeaderCellDef>Nombre</th>
                <td mat-cell *matCellDef="let user">
                  <div class="user-name">
                    <strong>{{ user.firstName }} {{ user.lastName }}</strong>
                    <small>{{ user.email }}</small>
                  </div>
                </td>
              </ng-container>

              <!-- Role Column -->
              <ng-container matColumnDef="role">
                <th mat-header-cell *matHeaderCellDef>Rol</th>
                <td mat-cell *matCellDef="let user">
                  <mat-chip [class]="getRoleChipClass(user.role)">
                    {{ getRoleDisplayName(user.role) }}
                  </mat-chip>
                </td>
              </ng-container>

              <!-- Status Column -->
              <ng-container matColumnDef="status">
                <th mat-header-cell *matHeaderCellDef>Estado</th>
                <td mat-cell *matCellDef="let user">
                  <mat-chip [class]="getStatusChipClass(user.isActive)">
                    {{ user.isActive ? 'Activo' : 'Inactivo' }}
                  </mat-chip>
                </td>
              </ng-container>

              <!-- Created At Column -->
              <ng-container matColumnDef="createdAt">
                <th mat-header-cell *matHeaderCellDef>Fecha de Registro</th>
                <td mat-cell *matCellDef="let user">
                  {{ formatDate(user.createdAt) }}
                </td>
              </ng-container>

              <!-- Actions Column -->
              <ng-container matColumnDef="actions">
                <th mat-header-cell *matHeaderCellDef>Acciones</th>
                <td mat-cell *matCellDef="let user">
                  <div class="action-buttons">
                    <button mat-icon-button 
                            (click)="viewUser(user)"
                            matTooltip="Ver perfil">
                      <mat-icon>visibility</mat-icon>
                    </button>
                    
                    <button mat-icon-button 
                            (click)="editUser(user)"
                            *ngIf="canEditUser(user)"
                            matTooltip="Editar usuario">
                      <mat-icon>edit</mat-icon>
                    </button>
                    
                    <button mat-icon-button 
                            (click)="toggleUserStatus(user)"
                            *ngIf="canToggleStatus(user)"
                            [matTooltip]="user.isActive ? 'Desactivar' : 'Activar'">
                      <mat-icon>{{ user.isActive ? 'block' : 'check_circle' }}</mat-icon>
                    </button>
                    
                    <button mat-icon-button 
                            (click)="deleteUser(user)"
                            *ngIf="canDeleteUser(user)"
                            matTooltip="Eliminar usuario"
                            color="warn">
                      <mat-icon>delete</mat-icon>
                    </button>
                  </div>
                </td>
              </ng-container>

              <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
              <tr mat-row *matRowDef="let row; columns: displayedColumns;"></tr>
            </table>
          </div>

          <!-- Paginación -->
          <mat-paginator 
            [length]="totalUsers()"
            [pageSize]="pageSize"
            [pageSizeOptions]="[10, 20, 50]"
            (page)="onPageChange($event)"
            showFirstLastButtons>
          </mat-paginator>

          <!-- Loading state -->
          <div *ngIf="isLoading()" class="loading-container">
            <mat-icon class="spinning">refresh</mat-icon>
            Cargando usuarios...
          </div>

          <!-- Empty state -->
          <div *ngIf="!isLoading() && users().length === 0" class="empty-state">
            <mat-icon>people_outline</mat-icon>
            <h3>No se encontraron usuarios</h3>
            <p>{{ searchQuery ? 'Intenta con otros términos de búsqueda' : 'Comienza creando el primer usuario' }}</p>
          </div>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .user-management-container {
      padding: 20px;
      max-width: 1200px;
      margin: 0 auto;
    }

    .management-card {
      margin-bottom: 20px;
    }

    .filters-section {
      display: flex;
      gap: 16px;
      margin-bottom: 20px;
      flex-wrap: wrap;
      align-items: center;
    }

    .search-field {
      flex: 1;
      min-width: 300px;
    }

    .role-filter {
      min-width: 150px;
    }

    .table-container {
      overflow-x: auto;
      margin-bottom: 20px;
    }

    .users-table {
      width: 100%;
      min-width: 800px;
    }

    .user-avatar {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-weight: bold;
      font-size: 14px;
    }

    .user-name {
      display: flex;
      flex-direction: column;
      gap: 2px;
    }

    .user-name small {
      color: #666;
      font-size: 12px;
    }

    .action-buttons {
      display: flex;
      gap: 8px;
    }

    /* Chip styles */
    .role-admin {
      background-color: #f44336;
      color: white;
    }

    .role-user {
      background-color: #2196f3;
      color: white;
    }

    .status-active {
      background-color: #4caf50;
      color: white;
    }

    .status-inactive {
      background-color: #ff9800;
      color: white;
    }

    .loading-container {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 40px;
      color: #666;
    }

    .spinning {
      animation: spin 1s linear infinite;
      margin-right: 8px;
    }

    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }

    .empty-state {
      text-align: center;
      padding: 40px;
      color: #666;
    }

    .empty-state mat-icon {
      font-size: 48px;
      width: 48px;
      height: 48px;
      margin-bottom: 16px;
      opacity: 0.5;
    }

    /* Responsive design */
    @media (max-width: 768px) {
      .user-management-container {
        padding: 10px;
      }

      .filters-section {
        flex-direction: column;
        align-items: stretch;
      }

      .search-field,
      .role-filter {
        min-width: unset;
      }

      .action-buttons {
        flex-direction: column;
      }
    }
  `]
})
export class UserManagementComponent implements OnInit {
  private readonly userService = inject(UserService);
  private readonly permissionService = inject(PermissionService);
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);

  // Signals for reactive state management
  users = signal<User[]>([]);
  totalUsers = signal<number>(0);
  isLoading = signal<boolean>(false);
  
  // Component state
  displayedColumns = ['avatar', 'name', 'role', 'status', 'createdAt', 'actions'];
  pageSize = 20;
  currentPage = 0;
  searchQuery = '';
  selectedRole = '';

  // Expose enums for template
  UserRole = UserRole;
  UserStatus = UserStatus;

  ngOnInit() {
    this.loadUsers();
  }

  loadUsers() {
    this.isLoading.set(true);
    const page = this.currentPage + 1; // Backend uses 1-based pagination

    this.userService.getAllUsers(page, this.pageSize).subscribe({
      next: (response: UsersListResponse) => {
        this.users.set(response.users);
        this.totalUsers.set(response.total);
        this.isLoading.set(false);
      },
      error: (error) => {
        console.error('Error loading users:', error);
        this.snackBar.open('Error al cargar usuarios', 'Cerrar', { duration: 3000 });
        this.isLoading.set(false);
      }
    });
  }

  searchUsers() {
    if (this.searchQuery.trim()) {
      this.isLoading.set(true);
      this.userService.searchUsers(this.searchQuery, 1, this.pageSize).subscribe({
        next: (response: UsersListResponse) => {
          this.users.set(response.users);
          this.totalUsers.set(response.total);
          this.isLoading.set(false);
        },
        error: (error) => {
          console.error('Error searching users:', error);
          this.snackBar.open('Error en la búsqueda', 'Cerrar', { duration: 3000 });
          this.isLoading.set(false);
        }
      });
    } else {
      this.loadUsers();
    }
  }

  filterByRole() {
    if (this.selectedRole) {
      this.isLoading.set(true);
      this.userService.getUsersByRole(this.selectedRole as UserRole).subscribe({
        next: (users: User[]) => {
          this.users.set(users);
          this.totalUsers.set(users.length);
          this.isLoading.set(false);
        },
        error: (error) => {
          console.error('Error filtering users:', error);
          this.snackBar.open('Error al filtrar usuarios', 'Cerrar', { duration: 3000 });
          this.isLoading.set(false);
        }
      });
    } else {
      this.loadUsers();
    }
  }

  onPageChange(event: PageEvent) {
    this.currentPage = event.pageIndex;
    this.pageSize = event.pageSize;
    this.loadUsers();
  }

  // Permission checks
  canCreateUsers(): boolean {
    return this.permissionService.hasPermission(Permission.USER_CREATE);
  }

  canEditUser(user: User): boolean {
    return this.permissionService.hasPermission(Permission.USER_UPDATE);
  }

  canDeleteUser(user: User): boolean {
    return this.permissionService.hasPermission(Permission.USER_DELETE);
  }

  canToggleStatus(user: User): boolean {
    return this.permissionService.hasPermission(Permission.USER_UPDATE_STATUS);
  }

  // User actions
  openCreateUserDialog() {
    const CreateUserDialog = CreateUserDialogComponent;
    const dialogRef = this.dialog.open(CreateUserDialog, {
      width: '500px',
      data: {}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.createUser(result);
      }
    });
  }

  viewUser(user: User) {
    // TODO: Open user detail dialog
    this.snackBar.open(`Ver perfil de ${user.firstName} ${user.lastName}`, 'Cerrar', { duration: 2000 });
  }

  editUser(user: User) {
    // TODO: Open edit user dialog
    this.snackBar.open(`Editar ${user.firstName} ${user.lastName}`, 'Cerrar', { duration: 2000 });
  }

  toggleUserStatus(user: User) {
    const newStatus = !user.isActive;
    this.userService.toggleUserStatus(user.id, { isActive: newStatus }).subscribe({
      next: () => {
        user.isActive = newStatus;
        this.snackBar.open(
          `Usuario ${newStatus ? 'activado' : 'desactivado'} correctamente`, 
          'Cerrar', 
          { duration: 3000 }
        );
      },
      error: (error) => {
        console.error('Error toggling user status:', error);
        this.snackBar.open('Error al cambiar estado del usuario', 'Cerrar', { duration: 3000 });
      }
    });
  }

  deleteUser(user: User) {
    if (confirm(`¿Estás seguro de que quieres eliminar a ${user.firstName} ${user.lastName}?`)) {
      this.userService.deleteUser(user.id).subscribe({
        next: () => {
          this.loadUsers(); // Reload the list
          this.snackBar.open('Usuario eliminado correctamente', 'Cerrar', { duration: 3000 });
        },
        error: (error) => {
          console.error('Error deleting user:', error);
          this.snackBar.open('Error al eliminar usuario', 'Cerrar', { duration: 3000 });
        }
      });
    }
  }

  createUser(userData: any) {
    this.userService.createUser(userData).subscribe({
      next: (newUser) => {
        this.loadUsers(); // Reload the list
        this.snackBar.open(
          `Usuario ${newUser.firstName} ${newUser.lastName} creado correctamente`, 
          'Cerrar', 
          { duration: 3000 }
        );
      },
      error: (error) => {
        console.error('Error creating user:', error);
        let errorMessage = 'Error al crear usuario';
        if (error.error?.error === 'admins cannot create other admin users') {
          errorMessage = 'Los administradores no pueden crear otros usuarios administradores';
        } else if (error.error?.error === 'email already exists') {
          errorMessage = 'El email ya está en uso';
        }
        this.snackBar.open(errorMessage, 'Cerrar', { duration: 3000 });
      }
    });
  }

  // Utility methods
  getUserInitials(user: User): string {
    return `${user.firstName.charAt(0)}${user.lastName.charAt(0)}`.toUpperCase();
  }

  getRoleDisplayName(role: UserRole): string {
    const roleNames = {
      [UserRole.ADMIN]: 'Administrador',
      [UserRole.USER]: 'Usuario',
      [UserRole.OPERATOR]: 'Operador',
      [UserRole.TREASURER]: 'Tesorero'
    };
    return roleNames[role] || role;
  }

  getRoleChipClass(role: UserRole): string {
    return `role-${role}`;
  }

  getStatusChipClass(isActive: boolean): string {
    return isActive ? 'status-active' : 'status-inactive';
  }

  formatDate(dateString?: string): string {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }
}

// Create User Dialog Component
@Component({
  selector: 'app-create-user-dialog',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatButtonModule,
    MatInputModule,
    MatSelectModule,
    MatIconModule
  ],
  template: `
    <h2 mat-dialog-title>
      <mat-icon>person_add</mat-icon>
      Crear Nuevo Usuario
    </h2>

    <mat-dialog-content>
      <form [formGroup]="createUserForm" class="create-user-form">
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Email</mat-label>
          <input matInput formControlName="email" type="email" placeholder="usuario@fintrack.com">
          <mat-error *ngIf="createUserForm.get('email')?.hasError('required')">
            El email es requerido
          </mat-error>
          <mat-error *ngIf="createUserForm.get('email')?.hasError('email')">
            Ingresa un email válido
          </mat-error>
        </mat-form-field>

        <div class="name-row">
          <mat-form-field appearance="outline">
            <mat-label>Nombre</mat-label>
            <input matInput formControlName="firstName" placeholder="Juan">
            <mat-error *ngIf="createUserForm.get('firstName')?.hasError('required')">
              El nombre es requerido
            </mat-error>
          </mat-form-field>

          <mat-form-field appearance="outline">
            <mat-label>Apellido</mat-label>
            <input matInput formControlName="lastName" placeholder="Pérez">
            <mat-error *ngIf="createUserForm.get('lastName')?.hasError('required')">
              El apellido es requerido
            </mat-error>
          </mat-form-field>
        </div>

        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Contraseña</mat-label>
          <input matInput formControlName="password" type="password" placeholder="Mínimo 8 caracteres">
          <mat-error *ngIf="createUserForm.get('password')?.hasError('required')">
            La contraseña es requerida
          </mat-error>
          <mat-error *ngIf="createUserForm.get('password')?.hasError('minlength')">
            La contraseña debe tener al menos 8 caracteres
          </mat-error>
        </mat-form-field>

        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Rol</mat-label>
          <mat-select formControlName="role">
            <mat-option [value]="UserRole.USER">Usuario</mat-option>
            <mat-option [value]="UserRole.OPERATOR">Operador</mat-option>
            <mat-option [value]="UserRole.TREASURER">Tesorero</mat-option>
            <!-- Admin option is intentionally excluded - admins cannot create other admins -->
          </mat-select>
          <mat-error *ngIf="createUserForm.get('role')?.hasError('required')">
            Selecciona un rol
          </mat-error>
        </mat-form-field>
      </form>
    </mat-dialog-content>

    <mat-dialog-actions align="end">
      <button mat-button (click)="onCancel()">Cancelar</button>
      <button mat-raised-button 
              color="primary" 
              (click)="onCreate()"
              [disabled]="createUserForm.invalid">
        <mat-icon>save</mat-icon>
        Crear Usuario
      </button>
    </mat-dialog-actions>

    <style>
      .create-user-form {
        display: flex;
        flex-direction: column;
        gap: 16px;
        min-width: 400px;
      }

      .name-row {
        display: flex;
        gap: 16px;
      }

      .name-row mat-form-field {
        flex: 1;
      }

      .full-width {
        width: 100%;
      }

      mat-dialog-title {
        display: flex;
        align-items: center;
        gap: 8px;
      }
    </style>
  `
})
export class CreateUserDialogComponent {
  private fb = inject(FormBuilder);
  private dialogRef = inject(MatDialogRef<CreateUserDialogComponent>);

  UserRole = UserRole;

  createUserForm: FormGroup = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    firstName: ['', [Validators.required, Validators.minLength(2)]],
    lastName: ['', [Validators.required, Validators.minLength(2)]],
    password: ['', [Validators.required, Validators.minLength(8)]],
    role: [UserRole.USER, [Validators.required]]
  });

  onCancel(): void {
    this.dialogRef.close();
  }

  onCreate(): void {
    if (this.createUserForm.valid) {
      this.dialogRef.close(this.createUserForm.value);
    }
  }
}