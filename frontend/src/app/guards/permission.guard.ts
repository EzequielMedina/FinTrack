import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { PermissionService } from '../services/permission.service';
import { AuthService } from '../services/auth.service';
import { Permission } from '../models';

/**
 * Guard que verifica si el usuario tiene un permiso específico
 */
export const permissionGuard = (requiredPermission: Permission): CanActivateFn => {
  return () => {
    const authService = inject(AuthService);
    const permissionService = inject(PermissionService);
    const router = inject(Router);

    if (!authService.isAuthenticatedSig()) {
      router.navigate(['/login']);
      return false;
    }

    if (!permissionService.hasPermission(requiredPermission)) {
      router.navigate(['/dashboard']);
      return false;
    }

    return true;
  };
};

/**
 * Guard que verifica si el usuario tiene al menos uno de los permisos especificados
 */
export const anyPermissionGuard = (requiredPermissions: Permission[]): CanActivateFn => {
  return () => {
    const authService = inject(AuthService);
    const permissionService = inject(PermissionService);
    const router = inject(Router);

    if (!authService.isAuthenticatedSig()) {
      router.navigate(['/login']);
      return false;
    }

    if (!permissionService.hasAnyPermission(requiredPermissions)) {
      router.navigate(['/dashboard']);
      return false;
    }

    return true;
  };
};

/**
 * Guard específico para el panel de administración
 */
export const adminPanelGuard: CanActivateFn = permissionGuard(Permission.ADMIN_PANEL_ACCESS);