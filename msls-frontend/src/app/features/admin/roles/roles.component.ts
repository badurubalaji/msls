/**
 * MSLS Roles Management Component
 *
 * Main component for managing roles - displays a list of roles with CRUD operations.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import {
  MslsModalComponent,
  TableColumn,
} from '../../../shared/components';
import { ToastService } from '../../../shared/services';
import { RbacService } from '../../../core/services';
import { Role, canDeleteRole, canModifyRole } from '../../../core/models';

import { RoleFormComponent } from './role-form.component';
import { RolePermissionsComponent } from './role-permissions.component';

@Component({
  selector: 'msls-roles',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsModalComponent,
    RoleFormComponent,
    RolePermissionsComponent,
  ],
  template: `
    <div class="roles-page">
      <div class="roles-card">
        <!-- Header -->
        <div class="roles-header">
          <div class="roles-header__left">
            <h1 class="roles-header__title">Role Management</h1>
            <p class="roles-header__subtitle">
              Manage roles and their permissions for your organization
            </p>
          </div>
          <div class="roles-header__right">
            <div class="search-input">
              <i class="fa-solid fa-magnifying-glass search-icon"></i>
              <input
                type="text"
                placeholder="Search roles..."
                [ngModel]="searchTerm()"
                (ngModelChange)="onSearchChange($event)"
                class="search-field"
              />
            </div>
            <button
              class="btn btn-primary"
              (click)="openCreateModal()"
            >
              <i class="fa-solid fa-plus"></i>
              Create Role
            </button>
          </div>
        </div>

        <!-- Loading State -->
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <p>Loading roles...</p>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <p>{{ error() }}</p>
            <button class="btn btn-secondary" (click)="loadRoles()">
              Retry
            </button>
          </div>
        } @else {
          <!-- Table -->
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Description</th>
                  <th style="width: 100px;">Type</th>
                  <th style="width: 150px;">Permissions</th>
                  <th style="width: 220px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (role of filteredRoles(); track role.id) {
                  <tr (click)="viewPermissions(role)" class="clickable-row">
                    <td class="name-cell">{{ role.name }}</td>
                    <td class="description-cell">{{ role.description || '-' }}</td>
                    <td>
                      <span class="badge" [class.badge-blue]="role.isSystem" [class.badge-gray]="!role.isSystem">
                        {{ role.isSystem ? 'System' : 'Custom' }}
                      </span>
                    </td>
                    <td>
                      <span class="permissions-badge">
                        {{ role.permissions?.length || 0 }} permissions
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button
                        class="action-btn"
                        (click)="viewPermissions(role); $event.stopPropagation()"
                        title="View permissions"
                      >
                        <i class="fa-regular fa-eye"></i>
                      </button>
                      @if (!role.isSystem) {
                        <button
                          class="action-btn action-btn--primary"
                          (click)="editPermissions(role); $event.stopPropagation()"
                          title="Edit permissions"
                        >
                          <i class="fa-solid fa-shield-halved"></i>
                        </button>
                        <button
                          class="action-btn"
                          (click)="editRole(role); $event.stopPropagation()"
                          title="Edit role"
                        >
                          <i class="fa-regular fa-pen-to-square"></i>
                        </button>
                        <button
                          class="action-btn action-btn--danger"
                          (click)="confirmDelete(role); $event.stopPropagation()"
                          title="Delete role"
                        >
                          <i class="fa-regular fa-trash-can"></i>
                        </button>
                      }
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="5" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No roles found</p>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        }
      </div>

      <!-- Create/Edit Role Modal -->
      <msls-modal
        [isOpen]="showRoleModal()"
        [title]="editingRole() ? 'Edit Role' : 'Create Role'"
        size="md"
        (closed)="closeRoleModal()"
      >
        <msls-role-form
          [role]="editingRole()"
          [loading]="saving()"
          (save)="saveRole($event)"
          (cancel)="closeRoleModal()"
        />
      </msls-modal>

      <!-- Permissions Modal -->
      <msls-modal
        [isOpen]="showPermissionsModal()"
        [title]="selectedRole()?.name + ' - Permissions'"
        size="lg"
        (closed)="closePermissionsModal()"
      >
        @if (selectedRole()) {
          <msls-role-permissions
            [role]="selectedRole()!"
            [canEdit]="editMode()"
            (permissionsChanged)="onPermissionsChanged()"
            (close)="closePermissionsModal()"
          />
        }
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Role"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete the role
            <strong>"{{ roleToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deleting()"
              (click)="deleteRole()"
            >
              @if (deleting()) {
                <div class="btn-spinner"></div>
                Deleting...
              } @else {
                Delete
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .roles-page {
      padding: 1.5rem;
      max-width: 1200px;
      margin: 0 auto;
    }

    .roles-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    /* Header */
    .roles-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
      padding-bottom: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .roles-header__title {
      font-size: 1.5rem;
      font-weight: 700;
      color: #0f172a;
      margin: 0 0 0.375rem 0;
    }

    .roles-header__subtitle {
      font-size: 0.875rem;
      color: #64748b;
      margin: 0;
    }

    .roles-header__right {
      display: flex;
      gap: 0.75rem;
      align-items: center;
    }

    /* Search Input */
    .search-input {
      position: relative;
    }

    .search-icon {
      position: absolute;
      left: 0.875rem;
      top: 50%;
      transform: translateY(-50%);
      color: #94a3b8;
      font-size: 0.875rem;
    }

    .search-field {
      padding: 0.625rem 0.875rem 0.625rem 2.5rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: white;
      color: #0f172a;
      width: 220px;
      transition: all 0.15s;
    }

    .search-field::placeholder {
      color: #94a3b8;
    }

    .search-field:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      padding: 0.625rem 1rem;
      font-size: 0.875rem;
      font-weight: 500;
      border-radius: 0.5rem;
      border: none;
      cursor: pointer;
      transition: all 0.15s;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    .btn-secondary {
      background: white;
      color: #334155;
      border: 1px solid #e2e8f0;
    }

    .btn-secondary:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
    }

    .btn-danger {
      background: #dc2626;
      color: white;
    }

    .btn-danger:hover:not(:disabled) {
      background: #b91c1c;
    }

    .btn-danger:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 1rem;
      height: 1rem;
      border: 2px solid rgba(255, 255, 255, 0.3);
      border-top-color: white;
      border-radius: 50%;
      animation: spin 0.6s linear infinite;
    }

    /* Loading */
    .loading-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 4rem;
      gap: 1rem;
    }

    .spinner {
      width: 2rem;
      height: 2rem;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .loading-container p {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0;
    }

    /* Error */
    .error-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 4rem;
      gap: 1rem;
    }

    .error-container i {
      font-size: 2.5rem;
      color: #dc2626;
    }

    .error-container p {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0;
    }

    /* Table */
    .table-container {
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      overflow: hidden;
    }

    .data-table {
      width: 100%;
      border-collapse: collapse;
    }

    .data-table thead {
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .data-table th {
      padding: 0.75rem 1rem;
      text-align: left;
      font-size: 0.75rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .data-table td {
      padding: 1rem;
      font-size: 0.875rem;
      border-bottom: 1px solid #f1f5f9;
      color: #334155;
    }

    .data-table tbody tr:last-child td {
      border-bottom: none;
    }

    .clickable-row {
      cursor: pointer;
      transition: background 0.15s;
    }

    .clickable-row:hover {
      background: #f8fafc;
    }

    .name-cell {
      font-weight: 500;
      color: #0f172a;
    }

    .description-cell {
      color: #64748b;
    }

    /* Badges */
    .badge {
      display: inline-flex;
      align-items: center;
      padding: 0.25rem 0.625rem;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .badge-blue {
      background: #dbeafe;
      color: #1e40af;
    }

    .badge-gray {
      background: #f1f5f9;
      color: #475569;
    }

    .permissions-badge {
      display: inline-flex;
      align-items: center;
      padding: 0.25rem 0.75rem;
      background: #eef2ff;
      color: #4338ca;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    /* Actions */
    .actions-cell {
      display: flex;
      justify-content: flex-end;
      gap: 0.5rem;
    }

    .action-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      background: transparent;
      border: 1px solid #e2e8f0;
      border-radius: 0.375rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.15s;
    }

    .action-btn:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
      color: #0f172a;
    }

    .action-btn--primary {
      border-color: #c7d2fe;
      color: #4f46e5;
    }

    .action-btn--primary:hover {
      background: #eef2ff;
      border-color: #a5b4fc;
      color: #4338ca;
    }

    .action-btn--danger:hover {
      background: #fef2f2;
      border-color: #fecaca;
      color: #dc2626;
    }

    /* Empty State */
    .empty-cell {
      padding: 3rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #94a3b8;
    }

    .empty-state i {
      font-size: 2.5rem;
    }

    .empty-state p {
      margin: 0;
      font-size: 0.875rem;
    }

    /* Delete Confirmation */
    .delete-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .delete-icon {
      width: 3.5rem;
      height: 3.5rem;
      margin: 0 auto 1rem;
      background: #fef2f2;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .delete-icon i {
      font-size: 1.5rem;
      color: #dc2626;
    }

    .delete-confirmation p {
      color: #475569;
      margin: 0 0 0.5rem 0;
    }

    .delete-confirmation strong {
      color: #dc2626;
    }

    .delete-warning {
      color: #dc2626 !important;
      font-size: 0.8125rem;
      font-weight: 500;
      padding: 0.75rem;
      background: #fef2f2;
      border-radius: 0.5rem;
      margin-top: 1rem !important;
    }

    .delete-actions {
      display: flex;
      justify-content: center;
      gap: 0.75rem;
      margin-top: 1.5rem;
    }

    /* Responsive */
    @media (max-width: 768px) {
      .roles-header {
        flex-direction: column;
        align-items: stretch;
      }

      .roles-header__right {
        flex-direction: column;
      }

      .search-field {
        width: 100%;
      }

      .btn-primary {
        width: 100%;
        justify-content: center;
      }

      .table-container {
        overflow-x: auto;
      }

      .data-table {
        min-width: 600px;
      }
    }
  `],
})
export class RolesComponent implements OnInit {
  rbacService = inject(RbacService);
  private toastService = inject(ToastService);

  // Table columns (kept for reference but not used with custom table)
  columns: TableColumn[] = [
    { key: 'name', label: 'Name', sortable: true },
    { key: 'description', label: 'Description', sortable: true },
    { key: 'isSystem', label: 'Type', sortable: true, width: '100px' },
    { key: 'permissions', label: 'Permissions', sortable: false, width: '150px' },
    { key: 'actions', label: 'Actions', sortable: false, width: '200px', align: 'right' },
  ];

  // State signals
  roles = signal<Role[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');

  // Modal state
  showRoleModal = signal(false);
  showPermissionsModal = signal(false);
  showDeleteModal = signal(false);
  editingRole = signal<Role | null>(null);
  selectedRole = signal<Role | null>(null);
  roleToDelete = signal<Role | null>(null);
  editMode = signal(false);

  // Computed filtered roles
  filteredRoles = computed(() => {
    const term = this.searchTerm().toLowerCase();
    if (!term) return this.roles();

    return this.roles().filter(
      role =>
        role.name.toLowerCase().includes(term) ||
        role.description?.toLowerCase().includes(term)
    );
  });

  ngOnInit(): void {
    this.loadRoles();
  }

  loadRoles(): void {
    this.loading.set(true);
    this.error.set(null);

    this.rbacService.getRoles(true).subscribe({
      next: roles => {
        this.roles.set(roles);
        this.loading.set(false);
      },
      error: err => {
        this.error.set('Failed to load roles. Please try again.');
        this.loading.set(false);
        console.error('Failed to load roles:', err);
      },
    });
  }

  onSearchChange(term: string): void {
    this.searchTerm.set(term);
  }

  // Create/Edit Modal
  openCreateModal(): void {
    this.editingRole.set(null);
    this.showRoleModal.set(true);
  }

  editRole(role: Role): void {
    if (canModifyRole(role)) {
      this.editingRole.set(role);
      this.showRoleModal.set(true);
    }
  }

  closeRoleModal(): void {
    this.showRoleModal.set(false);
    this.editingRole.set(null);
  }

  saveRole(data: { name: string; description?: string; permissionIds?: string[] }): void {
    this.saving.set(true);

    const editing = this.editingRole();
    const operation = editing
      ? this.rbacService.updateRole(editing.id, data)
      : this.rbacService.createRole(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Role updated successfully' : 'Role created successfully'
        );
        this.closeRoleModal();
        this.loadRoles();
        this.saving.set(false);
      },
      error: err => {
        this.toastService.error(
          editing ? 'Failed to update role' : 'Failed to create role'
        );
        this.saving.set(false);
        console.error('Failed to save role:', err);
      },
    });
  }

  // Permissions Modal
  viewPermissions(role: Role): void {
    this.editMode.set(false); // View only mode
    this.selectedRole.set(role);
    this.showPermissionsModal.set(true);
  }

  editPermissions(role: Role): void {
    this.editMode.set(true); // Enable editing
    this.selectedRole.set(role);
    this.showPermissionsModal.set(true);
  }

  closePermissionsModal(): void {
    this.showPermissionsModal.set(false);
    this.selectedRole.set(null);
    this.editMode.set(false);
  }

  onPermissionsChanged(): void {
    this.loadRoles();
  }

  // Delete Modal
  confirmDelete(role: Role): void {
    if (canDeleteRole(role)) {
      this.roleToDelete.set(role);
      this.showDeleteModal.set(true);
    }
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.roleToDelete.set(null);
  }

  deleteRole(): void {
    const role = this.roleToDelete();
    if (!role) return;

    this.deleting.set(true);

    this.rbacService.deleteRole(role.id).subscribe({
      next: () => {
        this.toastService.success('Role deleted successfully');
        this.closeDeleteModal();
        this.loadRoles();
        this.deleting.set(false);
      },
      error: err => {
        this.toastService.error('Failed to delete role');
        this.deleting.set(false);
        console.error('Failed to delete role:', err);
      },
    });
  }
}
