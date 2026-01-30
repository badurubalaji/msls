/**
 * Payroll feature module routes
 * Story 5.6: Payroll Processing
 */
import { Routes } from '@angular/router';
import { permissionGuard } from '../../core/guards';

export const PAYROLL_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./payroll.component').then((m) => m.PayrollComponent),
    canActivate: [permissionGuard(['payroll.view'])],
  },
  {
    path: 'runs',
    loadComponent: () =>
      import('./pages/pay-runs/pay-runs.component').then(
        (m) => m.PayRunsComponent
      ),
    canActivate: [permissionGuard(['payroll.view'])],
  },
  {
    path: 'runs/:id',
    loadComponent: () =>
      import('./pages/pay-run-detail/pay-run-detail.component').then(
        (m) => m.PayRunDetailComponent
      ),
    canActivate: [permissionGuard(['payroll.view'])],
  },
  {
    path: 'payslip/:id',
    loadComponent: () =>
      import('./pages/payslip-detail/payslip-detail.component').then(
        (m) => m.PayslipDetailComponent
      ),
    canActivate: [permissionGuard(['payroll.view'])],
  },
];
