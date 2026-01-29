/**
 * Staff Document feature module routes
 * Story 5.8: Staff Document Management
 */
import { Routes } from '@angular/router';
import { permissionGuard } from '../../core/guards';

export const STAFF_DOCUMENT_ROUTES: Routes = [
  {
    path: 'types',
    loadComponent: () =>
      import('./pages/document-types/document-types.component').then(
        (m) => m.DocumentTypesComponent
      ),
    canActivate: [permissionGuard(['staff_document.manage_types'])],
  },
  {
    path: 'expiring',
    loadComponent: () =>
      import('./pages/document-list/document-list.component').then(
        (m) => m.DocumentListComponent
      ),
    canActivate: [permissionGuard(['staff_document.view'])],
    data: { mode: 'expiring' },
  },
  {
    path: 'compliance',
    loadComponent: () =>
      import('./pages/compliance-report/compliance-report.component').then(
        (m) => m.ComplianceReportComponent
      ),
    canActivate: [permissionGuard(['staff_document.view'])],
  },
];
