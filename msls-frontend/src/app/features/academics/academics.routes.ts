import { Routes } from '@angular/router';

/**
 * Academics feature module routes.
 * Routes for class, section, stream, and structure management.
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
    path: 'sections',
    loadComponent: () =>
      import('./sections/sections.component').then((m) => m.SectionsComponent),
  },
  {
    path: 'streams',
    loadComponent: () =>
      import('./streams/streams.component').then((m) => m.StreamsComponent),
  },
  {
    path: 'structure',
    loadComponent: () =>
      import('./structure/structure.component').then((m) => m.StructureComponent),
  },
  {
    path: 'subjects',
    loadComponent: () =>
      import('./subjects/subjects.component').then((m) => m.SubjectsComponent),
  },
  {
    path: 'timetable',
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./timetable/timetable.component').then((m) => m.TimetableComponent),
      },
      {
        path: 'shifts',
        loadComponent: () =>
          import('./timetable/shifts/shifts.component').then((m) => m.ShiftsComponent),
      },
      {
        path: 'day-patterns',
        loadComponent: () =>
          import('./timetable/day-patterns/day-patterns.component').then((m) => m.DayPatternsComponent),
      },
      {
        path: 'period-slots',
        loadComponent: () =>
          import('./timetable/period-slots/period-slots.component').then((m) => m.PeriodSlotsComponent),
      },
      {
        path: 'list',
        loadComponent: () =>
          import('./timetable/timetable-list/timetable-list.component').then((m) => m.TimetableListComponent),
      },
      {
        path: 'builder/:id',
        loadComponent: () =>
          import('./timetable/timetable-builder/timetable-builder.component').then((m) => m.TimetableBuilderComponent),
      },
      {
        path: 'my',
        loadComponent: () =>
          import('./timetable/teacher-timetable/teacher-timetable.component').then((m) => m.TeacherTimetableComponent),
      },
      {
        path: 'substitutions',
        children: [
          {
            path: '',
            loadComponent: () =>
              import('./timetable/substitution/substitution-list/substitution-list.component').then((m) => m.SubstitutionListComponent),
          },
          {
            path: 'new',
            loadComponent: () =>
              import('./timetable/substitution/substitution-form/substitution-form.component').then((m) => m.SubstitutionFormComponent),
          },
          {
            path: ':id',
            loadComponent: () =>
              import('./timetable/substitution/substitution-detail/substitution-detail.component').then((m) => m.SubstitutionDetailComponent),
          },
        ],
      },
    ],
  },
];
