import { Component, inject } from '@angular/core';
import { RouterOutlet, Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { AuthService } from './services/auth.service';
import { PermissionService } from './services/permission.service';
import { HasPermissionDirective, HasRoleDirective } from './shared';
import { Permission, UserRole } from './models';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    RouterModule,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,
    MatTooltipModule,
    HasPermissionDirective,
    HasRoleDirective
  ],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly permissionService = inject(PermissionService);

  // Expose enums for template
  Permission = Permission;
  UserRole = UserRole;

  get isAuthenticated() {
    return this.auth.isAuthenticatedSig;
  }

  get currentUser() {
    return this.auth.currentUserSig;
  }

  canAccessAdminPanel(): boolean {
    return this.permissionService.canAccessAdminPanel();
  }

  logout(): void {
    this.auth.logout();
    this.router.navigateByUrl('/login');
  }

  navigateTo(route: string): void {
    this.router.navigateByUrl(route);
  }
}
