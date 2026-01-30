import { Routes } from '@angular/router';

/**
 * Students feature module routes.
 *
 * Routes:
 * - /students              - Student list page
 * - /students/new          - Create new student
 * - /students/import       - Bulk import students
 * - /students/promotion    - Promotion wizard
 * - /students/:id          - View student details
 * - /students/:id/edit     - Edit student
 */
export const STUDENTS_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/student-list/student-list.component').then(
        (m) => m.StudentListComponent
      ),
    title: 'Students',
  },
  {
    path: 'new',
    loadComponent: () =>
      import('./pages/student-form/student-form.component').then(
        (m) => m.StudentFormComponent
      ),
    title: 'Add Student',
  },
  {
    path: 'import',
    loadComponent: () =>
      import('./pages/student-import/student-import.component').then(
        (m) => m.StudentImportComponent
      ),
    title: 'Bulk Import Students',
  },
  {
    path: 'promotion',
    loadComponent: () =>
      import('./pages/promotion-wizard/promotion-wizard.component').then(
        (m) => m.PromotionWizardComponent
      ),
    title: 'Student Promotion',
  },
  {
    path: ':id',
    loadComponent: () =>
      import('./pages/student-detail/student-detail.component').then(
        (m) => m.StudentDetailComponent
      ),
    title: 'Student Details',
  },
  {
    path: ':id/edit',
    loadComponent: () =>
      import('./pages/student-form/student-form.component').then(
        (m) => m.StudentFormComponent
      ),
    title: 'Edit Student',
  },
];
