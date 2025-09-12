import { inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';

export interface LoginResponse {
  accessToken: string;
  refreshToken?: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly tokenKey = 'fintrack.token';
  private readonly refreshKey = 'fintrack.refresh';

  isAuthenticatedSig = signal<boolean>(this.hasToken());

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(`${environment.apiUrl}/auth/login`, {
        email,
        password
      })
      .pipe(
        tap((res) => {
          this.setToken(res.accessToken, res.refreshToken);
          this.isAuthenticatedSig.set(true);
        })
      );
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.refreshKey);
    this.isAuthenticatedSig.set(false);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  private setToken(access: string, refresh?: string) {
    localStorage.setItem(this.tokenKey, access);
    if (refresh) localStorage.setItem(this.refreshKey, refresh);
  }

  private hasToken(): boolean {
    return !!localStorage.getItem(this.tokenKey);
  }
}
