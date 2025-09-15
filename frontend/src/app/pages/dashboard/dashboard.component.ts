import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatGridListModule } from '@angular/material/grid-list';
import { HasPermissionDirective } from '../../shared/directives/has-permission.directive';
import { HasRoleDirective } from '../../shared/directives/has-role.directive';
import { Permission, UserRole } from '../../models';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule, 
    RouterModule,
    MatCardModule, 
    MatButtonModule, 
    MatIconModule,
    MatGridListModule,
    HasPermissionDirective,
    HasRoleDirective
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DashboardComponent implements OnInit {
  private readonly auth = inject(AuthService);

  // Exponemos los enums para usar en el template
  readonly Permission = Permission;
  readonly UserRole = UserRole;

  get currentUser() {
    return this.auth.currentUserSig;
  }

  ngOnInit(): void {
    // Cargar el perfil del usuario si no está disponible
    if (!this.currentUser()) {
      this.auth.loadUserProfile().subscribe({
        error: (err) => console.error('Error loading user profile:', err)
      });
    }
  }
}
