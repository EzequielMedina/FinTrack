import { inject, Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import {
  User,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateProfileRequest,
  ChangeRoleRequest,
  ToggleStatusRequest,
  ChangePasswordRequest,
  UsersListResponse,
  UserRole
} from '../models';

@Injectable({ providedIn: 'root' })
export class UserService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = `${environment.apiUrl}/users`;

  // CRUD Operations
  createUser(userData: CreateUserRequest): Observable<User> {
    return this.http.post<User>(this.apiUrl, userData);
  }

  getAllUsers(page: number = 1, pageSize: number = 20): Observable<UsersListResponse> {
    const params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());
    
    return this.http.get<UsersListResponse>(this.apiUrl, { params });
  }

  getUserById(id: string): Observable<User> {
    return this.http.get<User>(`${this.apiUrl}/${id}`);
  }

  updateUser(id: string, userData: UpdateUserRequest): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/${id}`, userData);
  }

  deleteUser(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  // Profile Management
  updateUserProfile(id: string, profileData: UpdateProfileRequest): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/${id}/profile`, profileData);
  }

  // Role and Status Management
  changeUserRole(id: string, roleData: ChangeRoleRequest): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/${id}/role`, roleData);
  }

  toggleUserStatus(id: string, statusData: ToggleStatusRequest): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/${id}/status`, statusData);
  }

  changePassword(id: string, passwordData: ChangePasswordRequest): Observable<void> {
    return this.http.put<void>(`${this.apiUrl}/${id}/password`, passwordData);
  }

  // Query Operations
  getUsersByRole(role: UserRole): Observable<User[]> {
    return this.http.get<User[]>(`${this.apiUrl}/role/${role}`);
  }

  // Search and Filter
  searchUsers(query: string, page: number = 1, pageSize: number = 20): Observable<UsersListResponse> {
    const params = new HttpParams()
      .set('q', query)
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());
    
    return this.http.get<UsersListResponse>(`${this.apiUrl}/search`, { params });
  }
}