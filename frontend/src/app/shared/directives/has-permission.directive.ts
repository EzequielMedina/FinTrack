import { Directive, Input, TemplateRef, ViewContainerRef, inject, OnInit, OnDestroy, effect } from '@angular/core';
import { PermissionService } from '../../services/permission.service';
import { AuthService } from '../../services/auth.service';
import { Permission } from '../../models';

/**
 * Directiva estructural que muestra/oculta elementos basado en permisos
 * 
 * Uso:
 * *appHasPermission="Permission.USER_CREATE"
 * *appHasPermission="[Permission.USER_CREATE, Permission.USER_UPDATE]; operator: 'OR'"
 */
@Directive({
  selector: '[appHasPermission]',
  standalone: true
})
export class HasPermissionDirective implements OnInit {
  private readonly templateRef = inject(TemplateRef<any>);
  private readonly viewContainer = inject(ViewContainerRef);
  private readonly permissionService = inject(PermissionService);
  private readonly authService = inject(AuthService);
  
  private hasView = false;

  @Input('appHasPermission') permissions!: Permission | Permission[];
  @Input('appHasPermissionOperator') operator: 'AND' | 'OR' = 'AND';

  constructor() {
    // Effect para reaccionar a cambios en el estado de autenticación
    effect(() => {
      // Acceder al signal para registrar la dependencia
      this.authService.isAuthenticatedSig();
      this.updateView();
    });
  }

  ngOnInit() {
    this.updateView();
  }

  private updateView() {
    const hasPermission = this.checkPermissions();
    
    if (hasPermission && !this.hasView) {
      this.viewContainer.createEmbeddedView(this.templateRef);
      this.hasView = true;
    } else if (!hasPermission && this.hasView) {
      this.viewContainer.clear();
      this.hasView = false;
    }
  }

  private checkPermissions(): boolean {
    if (!this.permissions) {
      return false;
    }

    if (Array.isArray(this.permissions)) {
      if (this.operator === 'OR') {
        return this.permissionService.hasAnyPermission(this.permissions);
      } else {
        return this.permissionService.hasAllPermissions(this.permissions);
      }
    } else {
      return this.permissionService.hasPermission(this.permissions);
    }
  }
}