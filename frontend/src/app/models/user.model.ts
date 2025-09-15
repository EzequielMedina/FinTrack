export enum UserRole {
  ADMIN = 'admin',
  USER = 'user',
  OPERATOR = 'operator',
  TREASURER = 'treasurer'
}

export enum UserStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended'
}

export interface Address {
  street?: string;
  city?: string;
  state?: string;
  postalCode?: string;  // Cambiar zipCode por postalCode para coincidir con backend
  country?: string;
}

export interface Preferences {
  language?: string;
  timezone?: string;
  notificationEmail?: boolean;  // Coincidir con backend
  notificationSMS?: boolean;    // Coincidir con backend
}

export interface UserProfile {
  phone?: string;              // Cambiar phoneNumber por phone
  dateOfBirth?: string;        // Backend devuelve Date, pero JSON será string
  address?: Address;
  profilePicture?: string;     // Agregar campo del backend
  preferences?: Preferences;
}

export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  fullName: string;           // Requerido en backend
  role: UserRole;
  isActive: boolean;
  emailVerified: boolean;     // Requerido en backend
  profile: UserProfile;       // Requerido en backend  
  createdAt: string;          // Requerido en backend como timestamp
  updatedAt: string;          // Requerido en backend como timestamp
  lastLoginAt?: string;       // Opcional
}

export interface CreateUserRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role: string; // Backend espera string, no enum
}

export interface UpdateUserRequest {
  email?: string;
  firstName?: string;
  lastName?: string;
  role?: string;           // Backend espera string
  isActive?: boolean;      // Agregar campo del backend
  emailVerified?: boolean; // Agregar campo del backend
}

export interface UpdateProfileRequest {
  phone?: string;          // Cambiar phoneNumber por phone
  dateOfBirth?: string;    // Backend maneja como Date en JSON string
  address?: {
    street?: string;
    city?: string;
    state?: string;
    postalCode?: string;   // Cambiar zipCode por postalCode
    country?: string;
  };
  preferences?: {
    language?: string;
    timezone?: string;
    notificationEmail?: boolean;
    notificationSMS?: boolean;
  };
}

export interface ChangeRoleRequest {
  role: string; // Backend espera string
}

export interface ChangePasswordRequest {
  oldPassword: string;     // Backend usa oldPassword
  newPassword: string;
}

export interface ToggleStatusRequest {
  isActive: boolean;
}

export interface UsersListResponse {
  users: User[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}