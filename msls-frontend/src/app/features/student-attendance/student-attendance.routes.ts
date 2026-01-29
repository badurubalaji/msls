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
];
