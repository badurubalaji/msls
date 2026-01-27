import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { DepartmentFormComponent } from './department-form.component';
import { DepartmentService } from './department.service';
import { Department, CreateDepartmentRequest } from './department.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'msls-departments',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent, DepartmentFormComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-building-user"></i>
          </div>
          <div class="header-text">
            <h1>Departments</h1>
            <p>Manage organizational departments and their structure</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Department
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search departments..."
            [ngModel]="searchTerm()"
            (ngModelChange)="searchTerm.set($event)"
            class="search-input"
          />
          @if (searchTerm()) {
            <button class="clear-search" (click)="searchTerm.set('')">
              <i class="fa-solid fa-xmark"></i>
            </button>
          }
        </div>
        <div class="filter-group">
          <select
            class="filter-select"
            [ngModel]="statusFilter()"
            (ngModelChange)="statusFilter.set($event)"
          >
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading departments...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadDepartments()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Department</th>
                <th>Code</th>
                <th>Branch</th>
                <th style="width: 100px; text-align: center;">Staff</th>
                <th style="width: 100px;">Status</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (dept of filteredDepartments(); track dept.id) {
                <tr>
                  <td class="name-cell">
                    <div class="name-wrapper">
                      <div class="dept-icon">
                        <i class="fa-solid fa-building-user"></i>
                      </div>
                      <div class="name-content">
                        <span class="name">{{ dept.name }}</span>
                        @if (dept.description) {
                          <span class="description">{{ dept.description }}</span>
                        }
                      </div>
                    </div>
                  </td>
                  <td class="code-cell">
                    <span class="code-badge">{{ dept.code }}</span>
                  </td>
                  <td class="branch-cell">{{ dept.branchName || '-' }}</td>
                  <td class="staff-cell">
                    <span class="staff-count">{{ dept.staffCount }}</span>
                  </td>
                  <td>
                    <span
                      class="badge"
                      [class.badge-green]="dept.isActive"
                      [class.badge-gray]="!dept.isActive"
                    >
                      {{ dept.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                  <td class="actions-cell">
                    <button
                      class="action-btn"
                      title="Edit"
                      (click)="editDepartment(dept)"
                    >
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button
                      class="action-btn"
                      title="Toggle Status"
                      (click)="toggleStatus(dept)"
                    >
                      <i
                        class="fa-solid"
                        [class.fa-toggle-on]="dept.isActive"
                        [class.fa-toggle-off]="!dept.isActive"
                      ></i>
                    </button>
                    <button
                      class="action-btn action-btn--danger"
                      title="Delete"
                      (click)="confirmDelete(dept)"
                    >
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="6" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-building"></i>
                      <p>No departments found</p>
                      @if (searchTerm() || statusFilter() !== 'all') {
                        <button class="btn btn-secondary btn-sm" (click)="clearFilters()">
                          Clear Filters
                        </button>
                      }
                    </div>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        }
      </div>

      <!-- Department Modal -->
      <msls-modal
        [isOpen]="showDepartmentModal()"
        [title]="editingDepartment() ? 'Edit Department' : 'Create Department'"
        size="lg"
        (closed)="closeDepartmentModal()"
      >
        <msls-department-form
          [department]="editingDepartment()"
          [loading]="saving()"
          (save)="saveDepartment($event)"
          (cancel)="closeDepartmentModal()"
        />
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Department"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete
            <strong>"{{ departmentToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">
            This action cannot be undone. All associated data will be permanently removed.
          </p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deleting()"
              (click)="deleteDepartment()"
            >
              @if (deleting()) {
                <div class="btn-spinner"></div>
                Deleting...
              } @else {
                <i class="fa-solid fa-trash"></i>
                Delete
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      background: #eef2ff;
      color: #4f46e5;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .header-text h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .header-text p {
      margin: 0.25rem 0 0;
      color: #64748b;
      font-size: 0.875rem;
    }

    .filters-bar {
      display: flex;
      gap: 1rem;
      margin-bottom: 1rem;
    }

    .search-box {
      flex: 1;
      max-width: 400px;
      position: relative;
    }

    .search-icon {
      position: absolute;
      left: 0.875rem;
      top: 50%;
      transform: translateY(-50%);
      color: #9ca3af;
    }

    .search-input {
      width: 100%;
      padding: 0.625rem 2.5rem 0.625rem 2.5rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
    }

    .search-input:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .clear-search {
      position: absolute;
      right: 0.5rem;
      top: 50%;
      transform: translateY(-50%);
      background: none;
      border: none;
      color: #9ca3af;
      cursor: pointer;
      padding: 0.25rem;
    }

    .clear-search:hover {
      color: #6b7280;
    }

    .filter-select {
      padding: 0.625rem 2rem 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
      cursor: pointer;
    }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .loading-container,
    .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-container {
      color: #dc2626;
      flex-direction: column;
    }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to {
        transform: rotate(360deg);
      }
    }

    .data-table {
      width: 100%;
      border-collapse: collapse;
    }

    .data-table th {
      text-align: left;
      padding: 0.875rem 1rem;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .data-table td {
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      color: #374151;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    .name-cell {
      min-width: 250px;
    }

    .name-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .dept-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #f1f5f9;
      color: #64748b;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .name-content {
      display: flex;
      flex-direction: column;
    }

    .name {
      font-weight: 500;
      color: #1e293b;
    }

    .description {
      font-size: 0.75rem;
      color: #64748b;
      max-width: 300px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .code-badge {
      display: inline-flex;
      padding: 0.25rem 0.5rem;
      background: #f1f5f9;
      border-radius: 0.25rem;
      font-family: monospace;
      font-size: 0.75rem;
      color: #475569;
    }

    .staff-cell {
      text-align: center;
    }

    .staff-count {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 2rem;
      height: 1.5rem;
      padding: 0 0.5rem;
      background: #f1f5f9;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
      color: #475569;
    }

    .badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .badge-green {
      background: #dcfce7;
      color: #166534;
    }

    .badge-gray {
      background: #f1f5f9;
      color: #64748b;
    }

    .actions-cell {
      text-align: right;
    }

    .action-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      border: none;
      background: transparent;
      color: #64748b;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: all 0.2s;
    }

    .action-btn:hover {
      background: #f1f5f9;
      color: #4f46e5;
    }

    .action-btn--danger:hover {
      background: #fef2f2;
      color: #dc2626;
    }

    .empty-cell {
      padding: 3rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #64748b;
    }

    .empty-state i {
      font-size: 2rem;
      color: #cbd5e1;
    }

    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    .btn-danger {
      background: #dc2626;
      color: white;
    }

    .btn-danger:hover:not(:disabled) {
      background: #b91c1c;
    }

    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 16px;
      height: 16px;
      border: 2px solid transparent;
      border-top-color: currentColor;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    .delete-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .delete-icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      border-radius: 50%;
      background: #fef2f2;
      color: #dc2626;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
    }

    .delete-confirmation p {
      margin: 0 0 0.5rem;
      color: #374151;
    }

    .delete-warning {
      font-size: 0.875rem;
      color: #64748b;
    }

    .delete-actions {
      display: flex;
      gap: 0.75rem;
      justify-content: center;
      margin-top: 1.5rem;
    }

    @media (max-width: 768px) {
      .page-header {
        flex-direction: column;
        gap: 1rem;
      }

      .filters-bar {
        flex-direction: column;
      }

      .search-box {
        max-width: 100%;
      }
    }
  `],
})
export class DepartmentsComponent implements OnInit {
  private departmentService = inject(DepartmentService);
  private toastService = inject(ToastService);

  // State signals
  departments = signal<Department[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');

  // Modal state
  showDepartmentModal = signal(false);
  showDeleteModal = signal(false);
  editingDepartment = signal<Department | null>(null);
  departmentToDelete = signal<Department | null>(null);

  // Computed filtered list
  filteredDepartments = computed(() => {
    let result = this.departments();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();

    if (term) {
      result = result.filter(
        dept =>
          dept.name.toLowerCase().includes(term) ||
          dept.code.toLowerCase().includes(term)
      );
    }

    if (status === 'active') {
      result = result.filter(dept => dept.isActive);
    } else if (status === 'inactive') {
      result = result.filter(dept => !dept.isActive);
    }

    return result;
  });

  ngOnInit(): void {
    this.loadDepartments();
  }

  loadDepartments(): void {
    this.loading.set(true);
    this.error.set(null);

    this.departmentService.getDepartments().subscribe({
      next: departments => {
        this.departments.set(departments);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load departments. Please try again.');
        this.loading.set(false);
      },
    });
  }

  openCreateModal(): void {
    this.editingDepartment.set(null);
    this.showDepartmentModal.set(true);
  }

  editDepartment(department: Department): void {
    this.editingDepartment.set(department);
    this.showDepartmentModal.set(true);
  }

  closeDepartmentModal(): void {
    this.showDepartmentModal.set(false);
    this.editingDepartment.set(null);
  }

  saveDepartment(data: CreateDepartmentRequest): void {
    this.saving.set(true);
    const editing = this.editingDepartment();

    const operation = editing
      ? this.departmentService.updateDepartment(editing.id, data)
      : this.departmentService.createDepartment(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Department updated successfully' : 'Department created successfully'
        );
        this.closeDepartmentModal();
        this.loadDepartments();
        this.saving.set(false);
      },
      error: () => {
        this.toastService.error(
          editing ? 'Failed to update department' : 'Failed to create department'
        );
        this.saving.set(false);
      },
    });
  }

  toggleStatus(department: Department): void {
    this.departmentService
      .updateDepartment(department.id, { isActive: !department.isActive })
      .subscribe({
        next: () => {
          this.toastService.success(
            `Department ${department.isActive ? 'deactivated' : 'activated'} successfully`
          );
          this.loadDepartments();
        },
        error: () => {
          this.toastService.error('Failed to update department status');
        },
      });
  }

  confirmDelete(department: Department): void {
    this.departmentToDelete.set(department);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.departmentToDelete.set(null);
  }

  deleteDepartment(): void {
    const department = this.departmentToDelete();
    if (!department) return;

    this.deleting.set(true);

    this.departmentService.deleteDepartment(department.id).subscribe({
      next: () => {
        this.toastService.success('Department deleted successfully');
        this.closeDeleteModal();
        this.loadDepartments();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete department. It may be in use.');
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('all');
  }
}
