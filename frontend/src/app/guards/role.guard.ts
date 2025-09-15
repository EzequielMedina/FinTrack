import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { UserRole } from '../models';

/**
 * Guard que permite el acceso solo a usuarios con rol de administrador
 */
export const adminGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (authService.isAuthenticatedSig() && authService.isAdmin()) {
    return true;
  }

  // Redirigir al dashboard si está autenticado pero no es admin
  if (authService.isAuthenticatedSig()) {
    router.navigate(['/dashboard']);
    return false;
  }

  // Redirigir al login si no está autenticado
  router.navigate(['/login']);
  return false;
};

/**
 * Guard que verifica si el usuario tiene un rol específico
 */
export const roleGuard = (allowedRoles: UserRole[]): CanActivateFn => {
  return () => {
    const authService = inject(AuthService);
    const router = inject(Router);

    if (!authService.isAuthenticatedSig()) {
      router.navigate(['/login']);
      return false;
    }

    const currentUser = authService.getCurrentUser();
    if (!currentUser || !allowedRoles.includes(currentUser.role)) {
      router.navigate(['/dashboard']);
      return false;
    }

    return true;
  };
};