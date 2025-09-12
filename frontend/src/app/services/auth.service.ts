import { inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';

export interface LoginResponse {
  accessToken: string;
  refreshToken?: string;
  user: {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
  };
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface RegisterResponse {
  accessToken: string;
  refreshToken?: string;
  user: {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
  };
}

export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly tokenKey = 'fintrack.token';
  private readonly refreshKey = 'fintrack.refresh';
  private readonly userKey = 'fintrack.user';

  isAuthenticatedSig = signal<boolean>(this.hasToken());
  currentUserSig = signal<User | null>(this.getCurrentUser());

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(`${environment.apiUrl}/auth/login`, {
        email,
        password
      })
      .pipe(
        tap((res) => {
          // Usar los campos en camelCase que envía el backend
          this.setToken(res.accessToken, res.refreshToken);
          this.setUser(res.user);
          this.isAuthenticatedSig.set(true);
          this.currentUserSig.set(res.user);
        })
      );
  }

  register(registerData: RegisterRequest): Observable<RegisterResponse> {
    return this.http
      .post<RegisterResponse>(`${environment.apiUrl}/auth/register`, registerData)
      .pipe(
        tap((res) => {
          this.setToken(res.accessToken, res.refreshToken);
          this.setUser(res.user);
          this.isAuthenticatedSig.set(true);
          this.currentUserSig.set(res.user);
        })
      );
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.refreshKey);
    localStorage.removeItem(this.userKey);
    this.isAuthenticatedSig.set(false);
    this.currentUserSig.set(null);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  getCurrentUser(): User | null {
    const userStr = localStorage.getItem(this.userKey);
    return userStr ? JSON.parse(userStr) : null;
  }

  loadUserProfile(): Observable<User> {
    return this.http.get<User>(`${environment.apiUrl}/me`).pipe(
      tap((user) => {
        this.setUser(user);
        this.currentUserSig.set(user);
      })
    );
  }

  private setToken(access: string, refresh?: string) {
    localStorage.setItem(this.tokenKey, access);
    if (refresh) localStorage.setItem(this.refreshKey, refresh);
  }

  private setUser(user: User) {
    localStorage.setItem(this.userKey, JSON.stringify(user));
  }

  private hasToken(): boolean {
    return !!localStorage.getItem(this.tokenKey);
  }
}
