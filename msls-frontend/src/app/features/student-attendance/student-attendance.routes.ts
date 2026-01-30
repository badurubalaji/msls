import { Routes } from '@angular/router';

export const STUDENT_ATTENDANCE_ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'mark',
    pathMatch: 'full',
  },
  {
    path: 'mark',
    loadComponent: () =>
      import('./pages/mark-attendance/mark-attendance.component').then(
        (m) => m.MarkAttendanceComponent
      ),
    data: {
      title: 'Mark Student Attendance',
      permissions: ['student_attendance:mark_class'],
    },
  },
  {
    path: 'period',
    loadComponent: () =>
      import('./pages/period-attendance/period-attendance.component').then(
        (m) => m.PeriodAttendanceComponent
      ),
    data: {
      title: 'Period-wise Attendance',
      permissions: ['student_attendance:mark_class'],
    },
  },
  {
    path: 'summary',
    loadComponent: () =>
      import('./pages/daily-summary/daily-summary.component').then(
        (m) => m.DailySummaryComponent
      ),
    data: {
      title: 'Daily Attendance Summary',
      permissions: ['student_attendance:view_class'],
    },
  },
];
