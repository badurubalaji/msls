import { Routes } from '@angular/router';
import { MainLayoutComponent } from './layouts/main-layout/main-layout.component';
import { authGuard, guestGuard } from './core/guards';

/**
 * Application routes configuration.
 *
 * Structure:
 * - Public routes (login, etc.) are at the root level with guestGuard
 * - Authenticated routes are wrapped with MainLayoutComponent and authGuard
 * - Lazy loading is used for feature modules
 */
export const routes: Routes = [
  // Public routes (no layout, protected by guestGuard)
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/login/login.component').then((m) => m.LoginComponent),
    canActivate: [guestGuard],
  },

  // Public status check route (no authentication required)
  {
    path: 'check-status',
    loadComponent: () =>
      import('./features/admissions/applications/pages/public-status-check/public-status-check').then(
        (m) => m.PublicStatusCheck
      ),
    title: 'Check Application Status | MSLS',
  },

  // Authenticated routes (with MainLayoutComponent)
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [authGuard],
    children: [
      // Dashboard (default route)
      {
        path: '',
        redirectTo: 'dashboard',
        pathMatch: 'full',
      },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./features/dashboard/dashboard.component').then(
            (m) => m.DashboardComponent
          ),
      },

      // Students
      {
        path: 'students',
        loadChildren: () =>
          import('./features/students/students.routes').then((m) => m.STUDENTS_ROUTES),
      },

      // Staff
      {
        path: 'staff',
        loadChildren: () =>
          import('./features/staff/staff.routes').then((m) => m.STAFF_ROUTES),
      },

      // Attendance
      {
        path: 'attendance',
        loadChildren: () =>
          import('./features/attendance/attendance.routes').then((m) => m.ATTENDANCE_ROUTES),
      },

      // Academics
      {
        path: 'academics',
        loadChildren: () =>
          import('./features/academics/academics.routes').then((m) => m.ACADEMICS_ROUTES),
      },

      // Finance
      {
        path: 'finance',
        loadChildren: () =>
          import('./features/finance/finance.routes').then((m) => m.FINANCE_ROUTES),
      },

      // Admissions
      {
        path: 'admissions',
        loadChildren: () =>
          import('./features/admissions/admissions.routes').then((m) => m.ADMISSIONS_ROUTES),
      },

      // Settings
      {
        path: 'settings',
        loadComponent: () =>
          import('./features/settings/settings.component').then(
            (m) => m.SettingsComponent
          ),
      },

      // Profile
      {
        path: 'profile',
        loadComponent: () =>
          import('./features/profile/profile.component').then((m) => m.ProfileComponent),
      },

      // Admin
      {
        path: 'admin',
        loadChildren: () =>
          import('./features/admin/admin.routes').then((m) => m.ADMIN_ROUTES),
      },
    ],
  },

  // Wildcard redirect
  {
    path: '**',
    redirectTo: 'login',
  },
];
