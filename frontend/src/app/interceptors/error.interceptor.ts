import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const auth = inject(AuthService);
  
  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      // Manejar errores de autenticación
      if (error.status === 401) {
        // Solo redirigir al login si el usuario está autenticado actualmente
        // Esto previene bucles de redirección en páginas de login/registro
        if (auth.isAuthenticatedSig()) {
          auth.logout();
          router.navigateByUrl('/login');
        }
      }
      
      // Manejar errores de autorización
      if (error.status === 403) {
        console.warn('Acceso prohibido:', error.message);
        // Podría redirigir a una página de error 403 o al dashboard
      }

      // Manejar errores del servidor
      if (error.status >= 500) {
        console.error('Error del servidor:', error.message);
        // Podría mostrar una notificación de error global
      }

      // Manejar errores de red
      if (error.status === 0) {
        console.error('Error de red:', error.message);
        // Podría mostrar una advertencia de conectividad
      }

      return throwError(() => error);
    })
  );
};
