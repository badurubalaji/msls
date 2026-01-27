import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { DesignationFormComponent } from './designation-form.component';
import { DesignationService } from './designation.service';
import { Designation, CreateDesignationRequest, DESIGNATION_LEVELS } from './designation.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'msls-designations',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent, DesignationFormComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-user-tie"></i>
          </div>
          <div class="header-text">
            <h1>Designations</h1>
            <p>Manage position titles and hierarchy levels</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Designation
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search designations..."
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
            <span>Loading designations...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadDesignations()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Designation</th>
                <th style="width: 140px;">Level</th>
                <th>Department</th>
                <th style="width: 100px; text-align: center;">Staff</th>
                <th style="width: 100px;">Status</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (desig of filteredDesignations(); track desig.id) {
                <tr>
                  <td class="name-cell">
                    <div class="name-wrapper">
                      <div class="desig-icon" [class]="'level-' + desig.level">
                        <i class="fa-solid fa-user-tie"></i>
                      </div>
                      <span class="name">{{ desig.name }}</span>
                    </div>
                  </td>
                  <td class="level-cell">
                    <span class="level-badge" [class]="'level-badge--' + getLevelColor(desig.level)">
                      Level {{ desig.level }}
                    </span>
                  </td>
                  <td class="department-cell">
                    @if (desig.departmentName) {
                      {{ desig.departmentName }}
                    } @else {
                      <span class="text-muted">All Departments</span>
                    }
                  </td>
                  <td class="staff-cell">
                    <span class="staff-count">{{ desig.staffCount }}</span>
                  </td>
                  <td>
                    <span
                      class="badge"
                      [class.badge-green]="desig.isActive"
                      [class.badge-gray]="!desig.isActive"
                    >
                      {{ desig.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                  <td class="actions-cell">
                    <button
                      class="action-btn"
                      title="Edit"
                      (click)="editDesignation(desig)"
                    >
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button
                      class="action-btn"
                      title="Toggle Status"
                      (click)="toggleStatus(desig)"
                    >
                      <i
                        class="fa-solid"
                        [class.fa-toggle-on]="desig.isActive"
                        [class.fa-toggle-off]="!desig.isActive"
                      ></i>
                    </button>
                    <button
                      class="action-btn action-btn--danger"
                      title="Delete"
                      (click)="confirmDelete(desig)"
                    >
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="6" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-id-badge"></i>
                      <p>No designations found</p>
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

      <!-- Designation Modal -->
      <msls-modal
        [isOpen]="showDesignationModal()"
        [title]="editingDesignation() ? 'Edit Designation' : 'Create Designation'"
        size="lg"
        (closed)="closeDesignationModal()"
      >
        <msls-designation-form
          [designation]="editingDesignation()"
          [loading]="saving()"
          (save)="saveDesignation($event)"
          (cancel)="closeDesignationModal()"
        />
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Designation"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete
            <strong>"{{ designationToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">
            This action cannot be undone. Staff members with this designation will need to be reassigned.
          </p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deleting()"
              (click)="deleteDesignation()"
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
      min-width: 200px;
    }

    .name-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .desig-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
      font-size: 0.875rem;
    }

    .desig-icon.level-1,
    .desig-icon.level-2 {
      background: #fef3c7;
      color: #b45309;
    }

    .desig-icon.level-3,
    .desig-icon.level-4 {
      background: #dbeafe;
      color: #1e40af;
    }

    .desig-icon.level-5,
    .desig-icon.level-6 {
      background: #dcfce7;
      color: #166534;
    }

    .desig-icon.level-7,
    .desig-icon.level-8,
    .desig-icon.level-9,
    .desig-icon.level-10 {
      background: #f1f5f9;
      color: #64748b;
    }

    .name {
      font-weight: 500;
      color: #1e293b;
    }

    .level-badge {
      display: inline-flex;
      padding: 0.25rem 0.625rem;
      border-radius: 0.375rem;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .level-badge--gold {
      background: #fef3c7;
      color: #92400e;
    }

    .level-badge--blue {
      background: #dbeafe;
      color: #1e40af;
    }

    .level-badge--green {
      background: #dcfce7;
      color: #166534;
    }

    .level-badge--gray {
      background: #f1f5f9;
      color: #64748b;
    }

    .text-muted {
      color: #9ca3af;
      font-style: italic;
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
export class DesignationsComponent implements OnInit {
  private designationService = inject(DesignationService);
  private toastService = inject(ToastService);

  // State signals
  designations = signal<Designation[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');

  // Modal state
  showDesignationModal = signal(false);
  showDeleteModal = signal(false);
  editingDesignation = signal<Designation | null>(null);
  designationToDelete = signal<Designation | null>(null);

  // Computed filtered list
  filteredDesignations = computed(() => {
    let result = this.designations();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();

    if (term) {
      result = result.filter(desig => desig.name.toLowerCase().includes(term));
    }

    if (status === 'active') {
      result = result.filter(desig => desig.isActive);
    } else if (status === 'inactive') {
      result = result.filter(desig => !desig.isActive);
    }

    return result;
  });

  ngOnInit(): void {
    this.loadDesignations();
  }

  loadDesignations(): void {
    this.loading.set(true);
    this.error.set(null);

    this.designationService.getDesignations().subscribe({
      next: designations => {
        this.designations.set(designations);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load designations. Please try again.');
        this.loading.set(false);
      },
    });
  }

  getLevelColor(level: number): string {
    if (level <= 2) return 'gold';
    if (level <= 4) return 'blue';
    if (level <= 6) return 'green';
    return 'gray';
  }

  openCreateModal(): void {
    this.editingDesignation.set(null);
    this.showDesignationModal.set(true);
  }

  editDesignation(designation: Designation): void {
    this.editingDesignation.set(designation);
    this.showDesignationModal.set(true);
  }

  closeDesignationModal(): void {
    this.showDesignationModal.set(false);
    this.editingDesignation.set(null);
  }

  saveDesignation(data: CreateDesignationRequest): void {
    this.saving.set(true);
    const editing = this.editingDesignation();

    const operation = editing
      ? this.designationService.updateDesignation(editing.id, data)
      : this.designationService.createDesignation(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Designation updated successfully' : 'Designation created successfully'
        );
        this.closeDesignationModal();
        this.loadDesignations();
        this.saving.set(false);
      },
      error: () => {
        this.toastService.error(
          editing ? 'Failed to update designation' : 'Failed to create designation'
        );
        this.saving.set(false);
      },
    });
  }

  toggleStatus(designation: Designation): void {
    this.designationService
      .updateDesignation(designation.id, { isActive: !designation.isActive })
      .subscribe({
        next: () => {
          this.toastService.success(
            `Designation ${designation.isActive ? 'deactivated' : 'activated'} successfully`
          );
          this.loadDesignations();
        },
        error: () => {
          this.toastService.error('Failed to update designation status');
        },
      });
  }

  confirmDelete(designation: Designation): void {
    this.designationToDelete.set(designation);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.designationToDelete.set(null);
  }

  deleteDesignation(): void {
    const designation = this.designationToDelete();
    if (!designation) return;

    this.deleting.set(true);

    this.designationService.deleteDesignation(designation.id).subscribe({
      next: () => {
        this.toastService.success('Designation deleted successfully');
        this.closeDeleteModal();
        this.loadDesignations();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete designation. It may be in use.');
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('all');
  }
}
