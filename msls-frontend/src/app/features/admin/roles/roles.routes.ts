/**
 * MSLS Roles Feature Routes
 *
 * Routes for role management functionality.
 */

import { Routes } from '@angular/router';
import { permissionGuard } from '../../../core/guards';

export const ROLES_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./roles.component').then(m => m.RolesComponent),
    canActivate: [permissionGuard(['roles:read'])],
  },
];
