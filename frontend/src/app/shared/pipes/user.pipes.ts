import { Pipe, PipeTransform } from '@angular/core';
import { UserRole, UserStatus } from '../../models';

@Pipe({
  name: 'roleDisplay',
  standalone: true
})
export class RoleDisplayPipe implements PipeTransform {
  transform(role: UserRole): string {
    const roleNames: Record<UserRole, string> = {
      [UserRole.ADMIN]: 'Administrador',
      [UserRole.USER]: 'Usuario',
      [UserRole.OPERATOR]: 'Operador',
      [UserRole.TREASURER]: 'Tesorero'
    };
    
    return roleNames[role] || role;
  }
}

@Pipe({
  name: 'statusDisplay',
  standalone: true
})
export class StatusDisplayPipe implements PipeTransform {
  transform(status: UserStatus | boolean): string {
    if (typeof status === 'boolean') {
      return status ? 'Activo' : 'Inactivo';
    }

    const statusNames: Record<UserStatus, string> = {
      [UserStatus.ACTIVE]: 'Activo',
      [UserStatus.INACTIVE]: 'Inactivo',
      [UserStatus.SUSPENDED]: 'Suspendido'
    };
    
    return statusNames[status] || status;
  }
}

@Pipe({
  name: 'userInitials',
  standalone: true
})
export class UserInitialsPipe implements PipeTransform {
  transform(firstName: string, lastName: string): string {
    if (!firstName || !lastName) {
      return '??';
    }
    
    return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();
  }
}