import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { SalaryStructureFormComponent } from './salary-structure-form.component';
import { SalaryService } from './salary.service';
import { SalaryStructure, CreateStructureRequest } from './salary.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'msls-salary-structures',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent, SalaryStructureFormComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-sitemap"></i>
          </div>
          <div class="header-text">
            <h1>Salary Structures</h1>
            <p>Define salary templates with component breakdowns</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Structure
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search structures..."
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
            <span>Loading structures...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadStructures()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Structure</th>
                <th>Code</th>
                <th>Designation</th>
                <th style="width: 100px; text-align: center;">Components</th>
                <th style="width: 100px; text-align: center;">Staff</th>
                <th style="width: 100px;">Status</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (structure of filteredStructures(); track structure.id) {
                <tr>
                  <td class="name-cell">
                    <div class="name-wrapper">
                      <div class="structure-icon">
                        <i class="fa-solid fa-sitemap"></i>
                      </div>
                      <div class="name-content">
                        <span class="name">{{ structure.name }}</span>
                        @if (structure.description) {
                          <span class="description">{{ structure.description }}</span>
                        }
                      </div>
                    </div>
                  </td>
                  <td class="code-cell">
                    <span class="code-badge">{{ structure.code }}</span>
                  </td>
                  <td class="designation-cell">
                    {{ structure.designationName || '-' }}
                  </td>
                  <td class="count-cell">
                    <span class="count-badge">{{ structure.componentCount }}</span>
                  </td>
                  <td class="count-cell">
                    <span class="count-badge staff">{{ structure.staffCount }}</span>
                  </td>
                  <td>
                    <span
                      class="badge"
                      [class.badge-green]="structure.isActive"
                      [class.badge-gray]="!structure.isActive"
                    >
                      {{ structure.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                  <td class="actions-cell">
                    <button class="action-btn" title="View Details" (click)="viewStructure(structure)">
                      <i class="fa-regular fa-eye"></i>
                    </button>
                    <button class="action-btn" title="Edit" (click)="editStructure(structure)">
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button class="action-btn" title="Toggle Status" (click)="toggleStatus(structure)">
                      <i
                        class="fa-solid"
                        [class.fa-toggle-on]="structure.isActive"
                        [class.fa-toggle-off]="!structure.isActive"
                      ></i>
                    </button>
                    <button
                      class="action-btn action-btn--danger"
                      title="Delete"
                      (click)="confirmDelete(structure)"
                      [disabled]="structure.staffCount > 0"
                    >
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="7" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-sitemap"></i>
                      <p>No salary structures found</p>
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

      <!-- Structure Modal (Create/Edit) -->
      <msls-modal
        [isOpen]="showStructureModal()"
        [title]="editingStructure() ? 'Edit Structure' : 'Create Structure'"
        size="lg"
        (closed)="closeStructureModal()"
      >
        <msls-salary-structure-form
          [structure]="editingStructure()"
          [loading]="saving()"
          (save)="saveStructure($event)"
          (cancel)="closeStructureModal()"
        />
      </msls-modal>

      <!-- View Details Modal -->
      <msls-modal
        [isOpen]="showViewModal()"
        [title]="viewingStructure()?.name || 'Structure Details'"
        size="lg"
        (closed)="closeViewModal()"
      >
        @if (viewingStructure()) {
          <div class="structure-details">
            <div class="details-header">
              <div class="detail-item">
                <span class="detail-label">Code</span>
                <span class="code-badge">{{ viewingStructure()?.code }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Status</span>
                <span
                  class="badge"
                  [class.badge-green]="viewingStructure()?.isActive"
                  [class.badge-gray]="!viewingStructure()?.isActive"
                >
                  {{ viewingStructure()?.isActive ? 'Active' : 'Inactive' }}
                </span>
              </div>
              @if (viewingStructure()?.designationName) {
                <div class="detail-item">
                  <span class="detail-label">Designation</span>
                  <span>{{ viewingStructure()?.designationName }}</span>
                </div>
              }
            </div>

            @if (viewingStructure()?.description) {
              <p class="structure-description">{{ viewingStructure()?.description }}</p>
            }

            <h4>Components</h4>
            <table class="components-table">
              <thead>
                <tr>
                  <th>Component</th>
                  <th>Type</th>
                  <th style="text-align: right;">Amount</th>
                  <th style="text-align: right;">Percentage</th>
                </tr>
              </thead>
              <tbody>
                @for (comp of viewingStructure()?.components || []; track comp.id) {
                  <tr>
                    <td>
                      <strong>{{ comp.componentName }}</strong>
                      <span class="comp-code">({{ comp.componentCode }})</span>
                    </td>
                    <td>
                      <span
                        class="type-badge"
                        [class.type-earning]="comp.componentType === 'earning'"
                        [class.type-deduction]="comp.componentType === 'deduction'"
                      >
                        {{ comp.componentType === 'earning' ? 'Earning' : 'Deduction' }}
                      </span>
                    </td>
                    <td class="amount-cell">{{ comp.amount || '-' }}</td>
                    <td class="amount-cell">{{ comp.percentage ? comp.percentage + '%' : '-' }}</td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="4" class="empty-components">No components defined</td>
                  </tr>
                }
              </tbody>
            </table>

            <div class="details-footer">
              <button class="btn btn-secondary" (click)="closeViewModal()">Close</button>
              <button class="btn btn-primary" (click)="editFromView()">
                <i class="fa-regular fa-pen-to-square"></i>
                Edit
              </button>
            </div>
          </div>
        }
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Structure"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete
            <strong>"{{ structureToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">
            This action cannot be undone. Make sure no staff members are using this structure.
          </p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteStructure()">
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
      background: #e0e7ff;
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
      to { transform: rotate(360deg); }
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

    .structure-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #e0e7ff;
      color: #4f46e5;
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

    .count-cell {
      text-align: center;
    }

    .count-badge {
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

    .count-badge.staff {
      background: #dbeafe;
      color: #1e40af;
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

    .action-btn:hover:not(:disabled) {
      background: #f1f5f9;
      color: #4f46e5;
    }

    .action-btn:disabled {
      opacity: 0.3;
      cursor: not-allowed;
    }

    .action-btn--danger:hover:not(:disabled) {
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

    /* View Modal Styles */
    .structure-details {
      padding: 0.5rem;
    }

    .details-header {
      display: flex;
      gap: 2rem;
      margin-bottom: 1rem;
      padding-bottom: 1rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .detail-item {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .detail-label {
      font-size: 0.75rem;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .structure-description {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0 0 1rem;
    }

    .structure-details h4 {
      margin: 0 0 0.75rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
    }

    .components-table {
      width: 100%;
      border-collapse: collapse;
      margin-bottom: 1rem;
    }

    .components-table th {
      text-align: left;
      padding: 0.625rem 0.75rem;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .components-table td {
      padding: 0.75rem;
      border-bottom: 1px solid #f1f5f9;
      font-size: 0.875rem;
    }

    .comp-code {
      color: #64748b;
      font-size: 0.75rem;
      margin-left: 0.25rem;
    }

    .type-badge {
      display: inline-flex;
      padding: 0.125rem 0.5rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .type-earning {
      background: #dcfce7;
      color: #166534;
    }

    .type-deduction {
      background: #fef2f2;
      color: #991b1b;
    }

    .amount-cell {
      text-align: right;
      font-family: monospace;
    }

    .empty-components {
      text-align: center;
      color: #64748b;
      padding: 1.5rem !important;
    }

    .details-footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
    }

    /* Buttons */
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

      .details-header {
        flex-direction: column;
        gap: 1rem;
      }
    }
  `],
})
export class SalaryStructuresComponent implements OnInit {
  private salaryService = inject(SalaryService);
  private toastService = inject(ToastService);

  // State signals
  structures = signal<SalaryStructure[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');

  // Modal state
  showStructureModal = signal(false);
  showViewModal = signal(false);
  showDeleteModal = signal(false);
  editingStructure = signal<SalaryStructure | null>(null);
  viewingStructure = signal<SalaryStructure | null>(null);
  structureToDelete = signal<SalaryStructure | null>(null);

  // Computed filtered list
  filteredStructures = computed(() => {
    let result = this.structures();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();

    if (term) {
      result = result.filter(
        structure =>
          structure.name.toLowerCase().includes(term) ||
          structure.code.toLowerCase().includes(term)
      );
    }

    if (status === 'active') {
      result = result.filter(structure => structure.isActive);
    } else if (status === 'inactive') {
      result = result.filter(structure => !structure.isActive);
    }

    return result;
  });

  ngOnInit(): void {
    this.loadStructures();
  }

  loadStructures(): void {
    this.loading.set(true);
    this.error.set(null);

    this.salaryService.getStructures().subscribe({
      next: structures => {
        this.structures.set(structures);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load salary structures. Please try again.');
        this.loading.set(false);
      },
    });
  }

  openCreateModal(): void {
    this.editingStructure.set(null);
    this.showStructureModal.set(true);
  }

  editStructure(structure: SalaryStructure): void {
    // Load full structure with components
    this.salaryService.getStructure(structure.id).subscribe({
      next: fullStructure => {
        this.editingStructure.set(fullStructure);
        this.showStructureModal.set(true);
      },
      error: () => {
        this.toastService.error('Failed to load structure details');
      },
    });
  }

  viewStructure(structure: SalaryStructure): void {
    this.salaryService.getStructure(structure.id).subscribe({
      next: fullStructure => {
        this.viewingStructure.set(fullStructure);
        this.showViewModal.set(true);
      },
      error: () => {
        this.toastService.error('Failed to load structure details');
      },
    });
  }

  closeStructureModal(): void {
    this.showStructureModal.set(false);
    this.editingStructure.set(null);
  }

  closeViewModal(): void {
    this.showViewModal.set(false);
    this.viewingStructure.set(null);
  }

  editFromView(): void {
    const structure = this.viewingStructure();
    if (structure) {
      this.closeViewModal();
      this.editingStructure.set(structure);
      this.showStructureModal.set(true);
    }
  }

  saveStructure(data: CreateStructureRequest): void {
    this.saving.set(true);
    const editing = this.editingStructure();

    const operation = editing
      ? this.salaryService.updateStructure(editing.id, data)
      : this.salaryService.createStructure(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Structure updated successfully' : 'Structure created successfully'
        );
        this.closeStructureModal();
        this.loadStructures();
        this.saving.set(false);
      },
      error: () => {
        this.toastService.error(
          editing ? 'Failed to update structure' : 'Failed to create structure'
        );
        this.saving.set(false);
      },
    });
  }

  toggleStatus(structure: SalaryStructure): void {
    this.salaryService
      .updateStructure(structure.id, { isActive: !structure.isActive })
      .subscribe({
        next: () => {
          this.toastService.success(
            `Structure ${structure.isActive ? 'deactivated' : 'activated'} successfully`
          );
          this.loadStructures();
        },
        error: () => {
          this.toastService.error('Failed to update structure status');
        },
      });
  }

  confirmDelete(structure: SalaryStructure): void {
    if (structure.staffCount > 0) {
      this.toastService.error('Cannot delete structure that is assigned to staff');
      return;
    }
    this.structureToDelete.set(structure);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.structureToDelete.set(null);
  }

  deleteStructure(): void {
    const structure = this.structureToDelete();
    if (!structure) return;

    this.deleting.set(true);

    this.salaryService.deleteStructure(structure.id).subscribe({
      next: () => {
        this.toastService.success('Structure deleted successfully');
        this.closeDeleteModal();
        this.loadStructures();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete structure. It may be in use.');
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('all');
  }
}
