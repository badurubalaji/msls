import { Routes } from '@angular/router';
import { permissionGuard } from '../../core/guards';

/**
 * Admin feature module routes.
 * Protected by settings:write permission.
 */
export const ADMIN_ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'feature-flags',
    pathMatch: 'full',
  },
  {
    path: 'feature-flags',
    loadComponent: () =>
      import('./feature-flags/feature-flags.component').then(
        (m) => m.AdminFeatureFlagsComponent
      ),
    canActivate: [permissionGuard(['settings:write'])],
  },
  {
    path: 'roles',
    loadComponent: () =>
      import('./roles/roles.component').then((m) => m.RolesComponent),
    canActivate: [permissionGuard(['roles:read'])],
  },
  {
    path: 'branches',
    loadComponent: () =>
      import('./branches/branches.component').then((m) => m.BranchesComponent),
    canActivate: [permissionGuard(['branches:read'])],
  },
  {
    path: 'academic-years',
    loadComponent: () =>
      import('./academic-years/academic-years.component').then(
        (m) => m.AcademicYearsComponent
      ),
    canActivate: [permissionGuard(['academic-years:read'])],
  },
  {
    path: 'departments',
    loadComponent: () =>
      import('./departments/departments.component').then(
        (m) => m.DepartmentsComponent
      ),
    canActivate: [permissionGuard(['department:read'])],
  },
  {
    path: 'designations',
    loadComponent: () =>
      import('./designations/designations.component').then(
        (m) => m.DesignationsComponent
      ),
    canActivate: [permissionGuard(['designation:read'])],
  },
  {
    path: 'salary-components',
    loadComponent: () =>
      import('./salary/salary-components.component').then(
        (m) => m.SalaryComponentsComponent
      ),
    canActivate: [permissionGuard(['salary:read'])],
  },
  {
    path: 'salary-structures',
    loadComponent: () =>
      import('./salary/salary-structures.component').then(
        (m) => m.SalaryStructuresComponent
      ),
    canActivate: [permissionGuard(['salary:read'])],
  },
];
