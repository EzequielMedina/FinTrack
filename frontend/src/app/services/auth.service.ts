import { inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap, switchMap } from 'rxjs';
import { environment } from '../../environments/environment';
import { 
  User, 
  AuthResponse, 
  LoginRequest, 
  RegisterRequest,
  UserRole 
} from '../models';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly tokenKey = 'fintrack.token';
  private readonly refreshKey = 'fintrack.refresh';
  private readonly userKey = 'fintrack.user';

  isAuthenticatedSig = signal<boolean>(this.hasToken());
  currentUserSig = signal<User | null>(this.getCurrentUser());

  login(email: string, password: string): Observable<User> {
    return this.http
      .post<AuthResponse>(`${environment.apiUrl}/auth/login`, {
        email,
        password
      })
      .pipe(
        switchMap((res) => {
          // Guardar solo los tokens del login
          this.setToken(res.accessToken, res.refreshToken);
          
          // Llamar a /api/me para obtener información completa del usuario
          return this.loadUserProfile();
        }),
        tap((user) => {
          this.isAuthenticatedSig.set(true);
          this.currentUserSig.set(user);
        })
      );
  }

  register(registerData: RegisterRequest): Observable<User> {
    return this.http
      .post<AuthResponse>(`${environment.apiUrl}/auth/register`, registerData)
      .pipe(
        switchMap((res) => {
          // Guardar solo los tokens del registro
          this.setToken(res.accessToken, res.refreshToken);
          
          // Llamar a /api/me para obtener información completa del usuario
          return this.loadUserProfile();
        }),
        tap((user) => {
          this.isAuthenticatedSig.set(true);
          this.currentUserSig.set(user);
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

  // Métodos adicionales para el sistema de roles
  hasRole(role: UserRole): boolean {
    const currentUser = this.getCurrentUser();
    return currentUser?.role === role;
  }

  isAdmin(): boolean {
    return this.hasRole(UserRole.ADMIN);
  }

  isUser(): boolean {
    return this.hasRole(UserRole.USER);
  }
}
