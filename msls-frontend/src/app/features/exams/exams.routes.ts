import { Routes } from '@angular/router';

export const EXAMS_ROUTES: Routes = [
  {
    path: '',
    children: [
      {
        path: '',
        redirectTo: 'types',
        pathMatch: 'full',
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
