/**
 * Teacher Assignment feature module routes
 * Story 5.7: Teacher Subject Assignment
 */
import { Routes } from '@angular/router';
import { permissionGuard } from '../../core/guards';

export const ASSIGNMENT_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/assignment-list/assignment-list.component').then(
        (m) => m.AssignmentListComponent
      ),
    canActivate: [permissionGuard(['assignment.view'])],
  },
  {
    path: 'workload',
    loadComponent: () =>
      import('./pages/workload-report/workload-report.component').then(
        (m) => m.WorkloadReportComponent
      ),
    canActivate: [permissionGuard(['assignment.workload'])],
  },
  {
    path: 'new',
    loadComponent: () =>
      import('./pages/assignment-form/assignment-form.component').then(
        (m) => m.AssignmentFormComponent
      ),
    canActivate: [permissionGuard(['assignment.create'])],
  },
  {
    path: ':id/edit',
    loadComponent: () =>
      import('./pages/assignment-form/assignment-form.component').then(
        (m) => m.AssignmentFormComponent
      ),
    canActivate: [permissionGuard(['assignment.update'])],
  },
];
