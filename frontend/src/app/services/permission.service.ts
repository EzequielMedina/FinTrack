import { inject, Injectable } from '@angular/core';
import { AuthService } from './auth.service';
import { Permission, ROLE_PERMISSIONS, UserRole } from '../models';

@Injectable({ providedIn: 'root' })
export class PermissionService {
  private readonly authService = inject(AuthService);

  /**
   * Verifica si el usuario actual tiene un permiso específico
   */
  hasPermission(permission: Permission): boolean {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) {
      return false;
    }

    const userPermissions = this.getUserPermissions(currentUser.role);
    return userPermissions.includes(permission);
  }

  /**
   * Verifica si el usuario actual tiene al menos uno de los permisos especificados
   */
  hasAnyPermission(permissions: Permission[]): boolean {
    return permissions.some(permission => this.hasPermission(permission));
  }

  /**
   * Verifica si el usuario actual tiene todos los permisos especificados
   */
  hasAllPermissions(permissions: Permission[]): boolean {
    return permissions.every(permission => this.hasPermission(permission));
  }

  /**
   * Obtiene todos los permisos de un rol específico
   */
  getUserPermissions(role: UserRole): Permission[] {
    return ROLE_PERMISSIONS[role] || [];
  }

  /**
   * Verifica si el usuario puede acceder a una funcionalidad específica
   */
  canAccessAdminPanel(): boolean {
    return this.hasPermission(Permission.ADMIN_PANEL_ACCESS);
  }

  canManageUsers(): boolean {
    return this.hasAnyPermission([
      Permission.USER_CREATE,
      Permission.USER_UPDATE,
      Permission.USER_DELETE
    ]);
  }

  canUpdateUserRoles(): boolean {
    return this.hasPermission(Permission.USER_UPDATE_ROLE);
  }

  canViewReports(): boolean {
    return this.hasPermission(Permission.REPORTS_VIEW);
  }

  canViewAnalytics(): boolean {
    return this.hasPermission(Permission.ANALYTICS_VIEW);
  }

  /**
   * Verifica si el usuario puede editar el perfil de otro usuario
   */
  canEditUserProfile(targetUserId: string): boolean {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) {
      return false;
    }

    // Si es el mismo usuario, puede editar su propio perfil
    if (currentUser.id === targetUserId) {
      return this.hasPermission(Permission.PROFILE_UPDATE_OWN);
    }

    // Si es admin, puede editar cualquier perfil
    return this.hasPermission(Permission.PROFILE_UPDATE_ANY);
  }

  /**
   * Verifica si el usuario puede ver el perfil de otro usuario
   */
  canViewUserProfile(targetUserId: string): boolean {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) {
      return false;
    }

    // Si es el mismo usuario, puede ver su propio perfil
    if (currentUser.id === targetUserId) {
      return this.hasPermission(Permission.PROFILE_READ_OWN);
    }

    // Si es admin, puede ver cualquier perfil
    return this.hasPermission(Permission.PROFILE_READ_ANY);
  }
}