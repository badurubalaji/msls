import { Routes } from '@angular/router';

/**
 * Academics feature module routes.
 * Placeholder routes for Story 1.8 layout integration.
 */
export const ACADEMICS_ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'classes',
    pathMatch: 'full',
  },
  {
    path: 'classes',
    loadComponent: () =>
      import('./classes/classes.component').then((m) => m.ClassesComponent),
  },
  {
    path: 'subjects',
    loadComponent: () =>
      import('./subjects/subjects.component').then((m) => m.SubjectsComponent),
  },
  {
    path: 'timetable',
    loadComponent: () =>
      import('./timetable/timetable.component').then((m) => m.TimetableComponent),
  },
];
