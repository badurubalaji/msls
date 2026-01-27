import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { SalaryComponentFormComponent } from './salary-component-form.component';
import { SalaryService } from './salary.service';
import { SalaryComponent, CreateComponentRequest, ComponentType } from './salary.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'msls-salary-components',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent, SalaryComponentFormComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-money-bill-wave"></i>
          </div>
          <div class="header-text">
            <h1>Salary Components</h1>
            <p>Manage earnings and deductions for salary structures</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Component
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search components..."
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
            [ngModel]="typeFilter()"
            (ngModelChange)="typeFilter.set($event)"
          >
            <option value="all">All Types</option>
            <option value="earning">Earnings</option>
            <option value="deduction">Deductions</option>
          </select>
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
            <span>Loading components...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadComponents()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Component</th>
                <th>Code</th>
                <th style="width: 120px;">Type</th>
                <th style="width: 140px;">Calculation</th>
                <th style="width: 100px; text-align: center;">Taxable</th>
                <th style="width: 100px;">Status</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (comp of filteredComponents(); track comp.id) {
                <tr>
                  <td class="name-cell">
                    <div class="name-wrapper">
                      <div
                        class="comp-icon"
                        [class.earning]="comp.componentType === 'earning'"
                        [class.deduction]="comp.componentType === 'deduction'"
                      >
                        <i
                          class="fa-solid"
                          [class.fa-arrow-up]="comp.componentType === 'earning'"
                          [class.fa-arrow-down]="comp.componentType === 'deduction'"
                        ></i>
                      </div>
                      <div class="name-content">
                        <span class="name">{{ comp.name }}</span>
                        @if (comp.description) {
                          <span class="description">{{ comp.description }}</span>
                        }
                      </div>
                    </div>
                  </td>
                  <td class="code-cell">
                    <span class="code-badge">{{ comp.code }}</span>
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
                  <td class="calc-cell">
                    @if (comp.calculationType === 'fixed') {
                      <span class="calc-type">Fixed</span>
                    } @else {
                      <span class="calc-type">
                        % of {{ comp.percentageOfName || 'Base' }}
                      </span>
                    }
                  </td>
                  <td class="taxable-cell">
                    <i
                      class="fa-solid"
                      [class.fa-check]="comp.isTaxable"
                      [class.fa-xmark]="!comp.isTaxable"
                      [style.color]="comp.isTaxable ? '#16a34a' : '#dc2626'"
                    ></i>
                  </td>
                  <td>
                    <span
                      class="badge"
                      [class.badge-green]="comp.isActive"
                      [class.badge-gray]="!comp.isActive"
                    >
                      {{ comp.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                  <td class="actions-cell">
                    <button class="action-btn" title="Edit" (click)="editComponent(comp)">
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button class="action-btn" title="Toggle Status" (click)="toggleStatus(comp)">
                      <i
                        class="fa-solid"
                        [class.fa-toggle-on]="comp.isActive"
                        [class.fa-toggle-off]="!comp.isActive"
                      ></i>
                    </button>
                    <button
                      class="action-btn action-btn--danger"
                      title="Delete"
                      (click)="confirmDelete(comp)"
                    >
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="7" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-money-bill-1"></i>
                      <p>No salary components found</p>
                      @if (searchTerm() || typeFilter() !== 'all' || statusFilter() !== 'all') {
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

      <!-- Component Modal -->
      <msls-modal
        [isOpen]="showComponentModal()"
        [title]="editingComponent() ? 'Edit Component' : 'Create Component'"
        size="lg"
        (closed)="closeComponentModal()"
      >
        <msls-salary-component-form
          [component]="editingComponent()"
          [loading]="saving()"
          (save)="saveComponent($event)"
          (cancel)="closeComponentModal()"
        />
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Component"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete
            <strong>"{{ componentToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">
            This component may be used in salary structures. Deleting it will remove it from all structures.
          </p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteComponent()">
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
      background: #fef3c7;
      color: #d97706;
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

    .comp-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .comp-icon.earning {
      background: #dcfce7;
      color: #16a34a;
    }

    .comp-icon.deduction {
      background: #fef2f2;
      color: #dc2626;
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

    .type-badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
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

    .calc-cell {
      font-size: 0.875rem;
    }

    .calc-type {
      color: #64748b;
    }

    .taxable-cell {
      text-align: center;
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
export class SalaryComponentsComponent implements OnInit {
  private salaryService = inject(SalaryService);
  private toastService = inject(ToastService);

  // State signals
  components = signal<SalaryComponent[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  typeFilter = signal<'all' | ComponentType>('all');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');

  // Modal state
  showComponentModal = signal(false);
  showDeleteModal = signal(false);
  editingComponent = signal<SalaryComponent | null>(null);
  componentToDelete = signal<SalaryComponent | null>(null);

  // Computed filtered list
  filteredComponents = computed(() => {
    let result = this.components();
    const term = this.searchTerm().toLowerCase();
    const type = this.typeFilter();
    const status = this.statusFilter();

    if (term) {
      result = result.filter(
        comp =>
          comp.name.toLowerCase().includes(term) ||
          comp.code.toLowerCase().includes(term)
      );
    }

    if (type !== 'all') {
      result = result.filter(comp => comp.componentType === type);
    }

    if (status === 'active') {
      result = result.filter(comp => comp.isActive);
    } else if (status === 'inactive') {
      result = result.filter(comp => !comp.isActive);
    }

    return result;
  });

  ngOnInit(): void {
    this.loadComponents();
  }

  loadComponents(): void {
    this.loading.set(true);
    this.error.set(null);

    this.salaryService.getComponents().subscribe({
      next: components => {
        this.components.set(components);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load salary components. Please try again.');
        this.loading.set(false);
      },
    });
  }

  openCreateModal(): void {
    this.editingComponent.set(null);
    this.showComponentModal.set(true);
  }

  editComponent(component: SalaryComponent): void {
    this.editingComponent.set(component);
    this.showComponentModal.set(true);
  }

  closeComponentModal(): void {
    this.showComponentModal.set(false);
    this.editingComponent.set(null);
  }

  saveComponent(data: CreateComponentRequest): void {
    this.saving.set(true);
    const editing = this.editingComponent();

    const operation = editing
      ? this.salaryService.updateComponent(editing.id, data)
      : this.salaryService.createComponent(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Component updated successfully' : 'Component created successfully'
        );
        this.closeComponentModal();
        this.loadComponents();
        this.saving.set(false);
      },
      error: () => {
        this.toastService.error(
          editing ? 'Failed to update component' : 'Failed to create component'
        );
        this.saving.set(false);
      },
    });
  }

  toggleStatus(component: SalaryComponent): void {
    this.salaryService
      .updateComponent(component.id, { isActive: !component.isActive })
      .subscribe({
        next: () => {
          this.toastService.success(
            `Component ${component.isActive ? 'deactivated' : 'activated'} successfully`
          );
          this.loadComponents();
        },
        error: () => {
          this.toastService.error('Failed to update component status');
        },
      });
  }

  confirmDelete(component: SalaryComponent): void {
    this.componentToDelete.set(component);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.componentToDelete.set(null);
  }

  deleteComponent(): void {
    const component = this.componentToDelete();
    if (!component) return;

    this.deleting.set(true);

    this.salaryService.deleteComponent(component.id).subscribe({
      next: () => {
        this.toastService.success('Component deleted successfully');
        this.closeDeleteModal();
        this.loadComponents();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete component. It may be in use.');
        this.deleting.set(false);
      },
    });
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.typeFilter.set('all');
    this.statusFilter.set('all');
  }
}
