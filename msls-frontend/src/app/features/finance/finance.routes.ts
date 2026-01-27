import { Routes } from '@angular/router';

/**
 * Finance feature module routes.
 * Placeholder routes for Story 1.8 layout integration.
 */
export const FINANCE_ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'fees',
    pathMatch: 'full',
  },
  {
    path: 'fees',
    loadComponent: () =>
      import('./fees/fees.component').then((m) => m.FeesComponent),
  },
  {
    path: 'reports',
    loadComponent: () =>
      import('./reports/reports.component').then((m) => m.ReportsComponent),
  },
];
