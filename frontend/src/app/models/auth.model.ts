export interface AuthResponse {
  accessToken: string;
  refreshToken?: string;
  user: {
    email: string;
    firstName: string;
    lastName: string;
  };
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface ApiError {
  error: string;
  message?: string;
  statusCode?: number;
}