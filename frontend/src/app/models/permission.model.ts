export enum Permission {
  // User Management
  USER_CREATE = 'user:create',
  USER_READ = 'user:read',
  USER_UPDATE = 'user:update',
  USER_DELETE = 'user:delete',
  USER_UPDATE_ROLE = 'user:update:role',
  USER_UPDATE_STATUS = 'user:update:status',
  
  // Profile Management
  PROFILE_READ_OWN = 'profile:read:own',
  PROFILE_UPDATE_OWN = 'profile:update:own',
  PROFILE_READ_ANY = 'profile:read:any',
  PROFILE_UPDATE_ANY = 'profile:update:any',
  
  // Admin Panel
  ADMIN_PANEL_ACCESS = 'admin:panel:access',
  
  // Reports and Analytics
  REPORTS_VIEW = 'reports:view',
  ANALYTICS_VIEW = 'analytics:view'
}

export interface RolePermissions {
  role: string;
  permissions: Permission[];
}

// Configuración de permisos por rol
export const ROLE_PERMISSIONS: Record<string, Permission[]> = {
  admin: [
    Permission.USER_CREATE,
    Permission.USER_READ,
    Permission.USER_UPDATE,
    Permission.USER_DELETE,
    Permission.USER_UPDATE_ROLE,
    Permission.USER_UPDATE_STATUS,
    Permission.PROFILE_READ_OWN,
    Permission.PROFILE_UPDATE_OWN,
    Permission.PROFILE_READ_ANY,
    Permission.PROFILE_UPDATE_ANY,
    Permission.ADMIN_PANEL_ACCESS,
    Permission.REPORTS_VIEW,
    Permission.ANALYTICS_VIEW
  ],
  operator: [
    Permission.USER_READ,
    Permission.USER_UPDATE,
    Permission.USER_UPDATE_STATUS,
    Permission.PROFILE_READ_OWN,
    Permission.PROFILE_UPDATE_OWN,
    Permission.PROFILE_READ_ANY,
    Permission.REPORTS_VIEW
  ],
  treasurer: [
    Permission.USER_READ,
    Permission.PROFILE_READ_OWN,
    Permission.PROFILE_UPDATE_OWN,
    Permission.REPORTS_VIEW,
    Permission.ANALYTICS_VIEW
  ],
  user: [
    Permission.PROFILE_READ_OWN,
    Permission.PROFILE_UPDATE_OWN
  ]
};