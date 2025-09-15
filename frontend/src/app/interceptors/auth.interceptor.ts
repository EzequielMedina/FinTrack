import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { catchError, throwError } from 'rxjs';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const auth = inject(AuthService);
  const router = inject(Router);
  
  // Agregar token de autorización si existe
  const token = auth.getToken();
  if (token) {
    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
  }

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      // Manejar errores de autorización
      if (error.status === 401) {
        // Token expirado o inválido - hacer logout y redirigir
        auth.logout();
        router.navigate(['/login']);
      } else if (error.status === 403) {
        // Sin permisos - redirigir al dashboard
        router.navigate(['/dashboard']);
      }

      return throwError(() => error);
    })
  );
};
