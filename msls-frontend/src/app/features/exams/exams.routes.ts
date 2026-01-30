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
        children: [
          {
            path: '',
            loadComponent: () =>
              import('./exam-type-list/exam-type-list.component').then(m => m.ExamTypeListComponent),
            title: 'Exam Types',
          },
          {
            path: 'new',
            loadComponent: () =>
              import('./exam-type-form/exam-type-form.component').then(m => m.ExamTypeFormComponent),
            title: 'Create Exam Type',
          },
          {
            path: ':id/edit',
            loadComponent: () =>
              import('./exam-type-form/exam-type-form.component').then(m => m.ExamTypeFormComponent),
            title: 'Edit Exam Type',
          },
        ],
      },
    ],
  },
];
