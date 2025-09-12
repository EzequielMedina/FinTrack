import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const authGuard: CanActivateFn = () => {
  const auth = inject(AuthService);
  const router = inject(Router);
  const isAuth = auth.isAuthenticatedSig();
  if (!isAuth) {
    router.navigateByUrl('/login');
    return false;
  }
  return true;
};
