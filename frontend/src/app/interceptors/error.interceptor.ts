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
      // Handle authentication errors
      if (error.status === 401) {
        // Only redirect to login if user is currently authenticated
        // This prevents redirect loops on login/register pages
        if (auth.isAuthenticatedSig()) {
          auth.logout();
          router.navigateByUrl('/login');
        }
      }
      
      // Handle authorization errors
      if (error.status === 403) {
        console.warn('Access forbidden:', error.message);
        // Could redirect to a 403 error page or dashboard
      }

      // Handle server errors
      if (error.status >= 500) {
        console.error('Server error:', error.message);
        // Could show a global error notification
      }

      // Handle network errors
      if (error.status === 0) {
        console.error('Network error:', error.message);
        // Could show a connectivity warning
      }

      return throwError(() => error);
    })
  );
};
