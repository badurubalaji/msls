/**
 * MSLS Admissions Feature Routes
 *
 * Lazy-loaded routes for the admissions module including
 * sessions, enquiries, and application handling.
 */

import { Routes } from '@angular/router';
import { permissionGuard } from '../../core/guards';

export const ADMISSIONS_ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'sessions',
    pathMatch: 'full',
  },
  {
    path: 'sessions',
    loadComponent: () =>
      import('./sessions/sessions.component').then(
        (m) => m.SessionsComponent
      ),
    canActivate: [permissionGuard(['admissions:read'])],
    title: 'Admission Sessions | MSLS',
  },
  {
    path: 'enquiries',
    loadComponent: () =>
      import('./enquiries/enquiries.component').then(
        (m) => m.EnquiriesComponent
      ),
    canActivate: [permissionGuard(['enquiries:read'])],
    title: 'Admission Enquiries | MSLS',
  },
  {
    path: 'applications',
    loadComponent: () =>
      import('./applications/applications.component').then(
        (m) => m.ApplicationsComponent
      ),
    canActivate: [permissionGuard(['applications:read'])],
    title: 'Admission Applications | MSLS',
  },
  {
    path: 'applications/new',
    loadComponent: () =>
      import('./applications/application-form.component').then(
        (m) => m.ApplicationFormComponent
      ),
    canActivate: [permissionGuard(['applications:create'])],
    title: 'New Application | MSLS',
  },
  {
    path: 'applications/:id',
    loadComponent: () =>
      import('./applications/application-view.component').then(
        (m) => m.ApplicationViewComponent
      ),
    canActivate: [permissionGuard(['applications:read'])],
    title: 'Application Details | MSLS',
  },
  {
    path: 'applications/:id/edit',
    loadComponent: () =>
      import('./applications/application-form.component').then(
        (m) => m.ApplicationFormComponent
      ),
    canActivate: [permissionGuard(['applications:update'])],
    title: 'Edit Application | MSLS',
  },
  {
    path: 'review',
    loadComponent: () =>
      import('./review/application-review.component').then(
        (m) => m.ApplicationReviewComponent
      ),
    canActivate: [permissionGuard(['applications:review'])],
    title: 'Application Review | MSLS',
  },
  {
    path: 'review/:id',
    loadComponent: () =>
      import('./review/application-review.component').then(
        (m) => m.ApplicationReviewComponent
      ),
    canActivate: [permissionGuard(['applications:review'])],
    title: 'Application Review | MSLS',
  },
  {
    path: 'tests',
    loadComponent: () =>
      import('./tests/pages/tests/tests').then(
        (m) => m.TestsComponent
      ),
    canActivate: [permissionGuard(['tests:read'])],
    title: 'Entrance Tests | MSLS',
  },
  {
    path: 'tests/:id/results',
    loadComponent: () =>
      import('./tests/components/results-entry/results-entry').then(
        (m) => m.ResultsEntryComponent
      ),
    canActivate: [permissionGuard(['tests:manage'])],
    title: 'Test Results | MSLS',
  },
  {
    path: 'tests/:id/registrations',
    loadComponent: () =>
      import('./tests/components/results-entry/results-entry').then(
        (m) => m.ResultsEntryComponent
      ),
    canActivate: [permissionGuard(['tests:read'])],
    title: 'Test Registrations | MSLS',
  },
  {
    path: 'merit-list',
    loadComponent: () =>
      import('./merit/merit-list.component').then(
        (m) => m.MeritListComponent
      ),
    canActivate: [permissionGuard(['admissions:read'])],
    title: 'Merit List | MSLS',
  },
  {
    path: 'enrollment',
    loadComponent: () =>
      import('./enrollment/enrollment-page/enrollment-page').then(
        (m) => m.EnrollmentPage
      ),
    canActivate: [permissionGuard(['admissions:update'])],
    title: 'Pending Enrollments | MSLS',
  },
  {
    path: 'dashboard',
    loadComponent: () =>
      import('./reports/admission-dashboard.component').then(
        (m) => m.AdmissionDashboardComponent
      ),
    canActivate: [permissionGuard(['admissions:read'])],
    title: 'Admission Dashboard | MSLS',
  },
  {
    path: 'reports',
    loadComponent: () =>
      import('./reports/admission-dashboard.component').then(
        (m) => m.AdmissionDashboardComponent
      ),
    canActivate: [permissionGuard(['admissions:read'])],
    title: 'Admission Reports | MSLS',
  },
];
