import { Routes } from '@angular/router';

/**
 * Attendance feature module routes.
 * Story 5.4: Staff Attendance Marking
 */
export const ATTENDANCE_ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'my',
    pathMatch: 'full',
  },
  {
    path: 'my',
    loadComponent: () =>
      import('./pages/my-attendance/my-attendance.component').then(
        (m) => m.MyAttendanceComponent
      ),
    title: 'My Attendance | MSLS',
  },
  {
    path: 'manage',
    loadComponent: () =>
      import('./pages/attendance-management/attendance-management.component').then(
        (m) => m.AttendanceManagementComponent
      ),
    title: 'Attendance Management | MSLS',
  },
];
