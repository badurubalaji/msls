import { Routes } from '@angular/router';

export const EXAMS_ROUTES: Routes = [
  {
    path: '',
    children: [
      {
        path: '',
        redirectTo: 'list',
        pathMatch: 'full',
      },
      {
        path: 'list',
        loadComponent: () =>
          import('./examination-list/examination-list.component').then(m => m.ExaminationListComponent),
        title: 'Examinations',
      },
      {
        path: ':id/schedules',
        loadComponent: () =>
          import('./examination-schedule/examination-schedule').then(m => m.ExaminationSchedule),
        title: 'Exam Schedules',
      },
      {
        path: 'types',
        loadComponent: () =>
          import('./exam-type-list/exam-type-list.component').then(m => m.ExamTypeListComponent),
        title: 'Exam Types',
      },
    ],
  },
];
