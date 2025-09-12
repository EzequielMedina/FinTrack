import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const guestGuard: CanActivateFn = () => {
  const auth = inject(AuthService);
  const router = inject(Router);
  const isAuth = auth.isAuthenticatedSig();
  
  if (isAuth) {
    router.navigateByUrl('/');
    return false;
  }
  return true;
};