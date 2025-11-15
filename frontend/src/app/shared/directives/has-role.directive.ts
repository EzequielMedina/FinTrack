import { Directive, Input, TemplateRef, ViewContainerRef, inject, OnInit, effect } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { UserRole } from '../../models';

/**
 * Directiva estructural que muestra/oculta elementos basado en roles
 * 
 * Uso:
 * *appHasRole="UserRole.ADMIN"
 * *appHasRole="[UserRole.ADMIN, UserRole.USER]"
 */
@Directive({
  selector: '[appHasRole]',
  standalone: true
})
export class HasRoleDirective implements OnInit {
  private readonly templateRef = inject(TemplateRef<any>);
  private readonly viewContainer = inject(ViewContainerRef);
  private readonly authService = inject(AuthService);
  
  private hasView = false;

  @Input('appHasRole') roles!: UserRole | UserRole[];

  constructor() {
    // Effect para reaccionar a cambios en el usuario actual
    effect(() => {
      // Acceder al signal para registrar la dependencia
      this.authService.currentUserSig();
      this.updateView();
    });
  }

  ngOnInit() {
    this.updateView();
  }

  private updateView() {
    const hasRole = this.checkRoles();
    
    if (hasRole && !this.hasView) {
      this.viewContainer.createEmbeddedView(this.templateRef);
      this.hasView = true;
    } else if (!hasRole && this.hasView) {
      this.viewContainer.clear();
      this.hasView = false;
    }
  }

  private checkRoles(): boolean {
    if (!this.roles) {
      return false;
    }

    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) {
      return false;
    }

    if (Array.isArray(this.roles)) {
      return this.roles.includes(currentUser.role);
    } else {
      return currentUser.role === this.roles;
    }
  }
}