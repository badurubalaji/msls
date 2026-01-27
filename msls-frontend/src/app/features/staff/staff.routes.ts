import { Routes } from '@angular/router';

/**
 * Staff feature module routes.
 * Story 5.1: Staff Profile Management
 */
export const STAFF_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/staff-list/staff-list.component').then(
        (m) => m.StaffListComponent
      ),
  },
  {
    path: 'new',
    loadComponent: () =>
      import('./pages/staff-form/staff-form.component').then(
        (m) => m.StaffFormComponent
      ),
  },
  {
    path: ':id',
    loadComponent: () =>
      import('./pages/staff-detail/staff-detail.component').then(
        (m) => m.StaffDetailComponent
      ),
  },
  {
    path: ':id/edit',
    loadComponent: () =>
      import('./pages/staff-form/staff-form.component').then(
        (m) => m.StaffFormComponent
      ),
  },
];
