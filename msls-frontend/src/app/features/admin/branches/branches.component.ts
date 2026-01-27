/**
 * MSLS Branches Management Component
 *
 * Main component for managing branches - displays a list of branches with CRUD operations.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MslsModalComponent } from '../../../shared/components';
import { ToastService } from '../../../shared/services';
import { Branch, CreateBranchRequest } from './branch.model';
import { BranchService } from './branch.service';
import { BranchFormComponent } from './branch-form.component';

@Component({
  selector: 'msls-branches',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsModalComponent,
    BranchFormComponent,
  ],
  template: `
    <div class="branches-page">
      <div class="branches-card">
        <!-- Header -->
        <div class="branches-header">
          <div class="branches-header__left">
            <h1 class="branches-header__title">Branch Management</h1>
            <p class="branches-header__subtitle">
              Manage branches for your multi-campus school setup
            </p>
          </div>
          <div class="branches-header__right">
            <div class="search-input">
              <i class="fa-solid fa-magnifying-glass search-icon"></i>
              <input
                type="text"
                placeholder="Search branches..."
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
              Add Branch
            </button>
          </div>
        </div>

        <!-- Loading State -->
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <p>Loading branches...</p>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <p>{{ error() }}</p>
            <button class="btn btn-secondary" (click)="loadBranches()">
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
                  <th>Code</th>
                  <th>City</th>
                  <th style="width: 100px;">Status</th>
                  <th style="width: 100px;">Primary</th>
                  <th style="width: 200px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (branch of filteredBranches(); track branch.id) {
                  <tr>
                    <td class="name-cell">{{ branch.name }}</td>
                    <td class="code-cell">{{ branch.code }}</td>
                    <td class="city-cell">{{ branch.city || '-' }}</td>
                    <td>
                      <span
                        class="badge"
                        [class.badge-green]="branch.isActive"
                        [class.badge-gray]="!branch.isActive"
                      >
                        {{ branch.isActive ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td>
                      @if (branch.isPrimary) {
                        <span class="badge badge-blue">
                          <i class="fa-solid fa-star badge-icon"></i>
                          Primary
                        </span>
                      } @else {
                        <span class="badge badge-light">-</span>
                      }
                    </td>
                    <td class="actions-cell">
                      <button
                        class="action-btn"
                        (click)="editBranch(branch)"
                        title="Edit branch"
                      >
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      @if (!branch.isPrimary) {
                        <button
                          class="action-btn action-btn--primary"
                          (click)="setPrimary(branch)"
                          title="Set as primary"
                        >
                          <i class="fa-regular fa-star"></i>
                        </button>
                      }
                      <button
                        class="action-btn"
                        [class.action-btn--success]="!branch.isActive"
                        [class.action-btn--warning]="branch.isActive"
                        (click)="toggleStatus(branch)"
                        [title]="branch.isActive ? 'Deactivate' : 'Activate'"
                      >
                        <i
                          class="fa-solid"
                          [class.fa-toggle-on]="branch.isActive"
                          [class.fa-toggle-off]="!branch.isActive"
                        ></i>
                      </button>
                      @if (!branch.isPrimary) {
                        <button
                          class="action-btn action-btn--danger"
                          (click)="confirmDelete(branch)"
                          title="Delete branch"
                        >
                          <i class="fa-regular fa-trash-can"></i>
                        </button>
                      }
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="6" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-building"></i>
                        <p>No branches found</p>
                        @if (!searchTerm()) {
                          <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                            <i class="fa-solid fa-plus"></i>
                            Add your first branch
                          </button>
                        }
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        }
      </div>

      <!-- Create/Edit Branch Modal -->
      <msls-modal
        [isOpen]="showBranchModal()"
        [title]="editingBranch() ? 'Edit Branch' : 'Create Branch'"
        size="lg"
        (closed)="closeBranchModal()"
      >
        <msls-branch-form
          [branch]="editingBranch()"
          [loading]="saving()"
          (save)="saveBranch($event)"
          (cancel)="closeBranchModal()"
        />
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Branch"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete the branch
            <strong>"{{ branchToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">
            This action cannot be undone. All data associated with this branch may be affected.
          </p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deleting()"
              (click)="deleteBranch()"
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

      <!-- Set Primary Confirmation Modal -->
      <msls-modal
        [isOpen]="showPrimaryModal()"
        title="Set Primary Branch"
        size="sm"
        (closed)="closePrimaryModal()"
      >
        <div class="primary-confirmation">
          <div class="primary-icon">
            <i class="fa-solid fa-star"></i>
          </div>
          <p>
            Set <strong>"{{ branchToSetPrimary()?.name }}"</strong> as the primary branch?
          </p>
          <p class="primary-note">
            The current primary branch will be unset automatically.
          </p>
          <div class="primary-actions">
            <button class="btn btn-secondary" (click)="closePrimaryModal()">
              Cancel
            </button>
            <button
              class="btn btn-primary"
              [disabled]="settingPrimary()"
              (click)="confirmSetPrimary()"
            >
              @if (settingPrimary()) {
                <div class="btn-spinner"></div>
                Setting...
              } @else {
                Set as Primary
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .branches-page {
      padding: 1.5rem;
      max-width: 1200px;
      margin: 0 auto;
    }

    .branches-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    /* Header */
    .branches-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
      padding-bottom: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .branches-header__title {
      font-size: 1.5rem;
      font-weight: 700;
      color: #0f172a;
      margin: 0 0 0.375rem 0;
    }

    .branches-header__subtitle {
      font-size: 0.875rem;
      color: #64748b;
      margin: 0;
    }

    .branches-header__right {
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

    .btn-sm {
      padding: 0.5rem 0.875rem;
      font-size: 0.8125rem;
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

    .data-table tbody tr {
      transition: background 0.15s;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    .name-cell {
      font-weight: 500;
      color: #0f172a;
    }

    .code-cell {
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      font-size: 0.8125rem;
      color: #64748b;
      background: #f1f5f9;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
      display: inline-block;
    }

    .city-cell {
      color: #64748b;
    }

    /* Badges */
    .badge {
      display: inline-flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.25rem 0.625rem;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .badge-icon {
      font-size: 0.625rem;
    }

    .badge-green {
      background: #dcfce7;
      color: #166534;
    }

    .badge-gray {
      background: #f1f5f9;
      color: #475569;
    }

    .badge-blue {
      background: #dbeafe;
      color: #1e40af;
    }

    .badge-light {
      background: #f8fafc;
      color: #94a3b8;
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

    .action-btn--success {
      border-color: #bbf7d0;
      color: #16a34a;
    }

    .action-btn--success:hover {
      background: #f0fdf4;
      border-color: #86efac;
      color: #15803d;
    }

    .action-btn--warning {
      border-color: #fed7aa;
      color: #ea580c;
    }

    .action-btn--warning:hover {
      background: #fff7ed;
      border-color: #fdba74;
      color: #c2410c;
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

    /* Primary Confirmation */
    .primary-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .primary-icon {
      width: 3.5rem;
      height: 3.5rem;
      margin: 0 auto 1rem;
      background: #fef9c3;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .primary-icon i {
      font-size: 1.5rem;
      color: #ca8a04;
    }

    .primary-confirmation p {
      color: #475569;
      margin: 0 0 0.5rem 0;
    }

    .primary-confirmation strong {
      color: #0f172a;
    }

    .primary-note {
      font-size: 0.8125rem;
      color: #64748b !important;
      padding: 0.75rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      margin-top: 1rem !important;
    }

    .primary-actions {
      display: flex;
      justify-content: center;
      gap: 0.75rem;
      margin-top: 1.5rem;
    }

    /* Responsive */
    @media (max-width: 768px) {
      .branches-header {
        flex-direction: column;
        align-items: stretch;
      }

      .branches-header__right {
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
        min-width: 700px;
      }
    }
  `],
})
export class BranchesComponent implements OnInit {
  private branchService = inject(BranchService);
  private toastService = inject(ToastService);

  // State signals
  branches = signal<Branch[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  settingPrimary = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');

  // Modal state
  showBranchModal = signal(false);
  showDeleteModal = signal(false);
  showPrimaryModal = signal(false);
  editingBranch = signal<Branch | null>(null);
  branchToDelete = signal<Branch | null>(null);
  branchToSetPrimary = signal<Branch | null>(null);

  // Computed filtered branches
  filteredBranches = computed(() => {
    const term = this.searchTerm().toLowerCase();
    if (!term) return this.branches();

    return this.branches().filter(
      branch =>
        branch.name.toLowerCase().includes(term) ||
        branch.code.toLowerCase().includes(term) ||
        branch.city?.toLowerCase().includes(term)
    );
  });

  ngOnInit(): void {
    this.loadBranches();
  }

  loadBranches(): void {
    this.loading.set(true);
    this.error.set(null);

    this.branchService.getBranches().subscribe({
      next: branches => {
        this.branches.set(branches);
        this.loading.set(false);
      },
      error: err => {
        this.error.set('Failed to load branches. Please try again.');
        this.loading.set(false);
        console.error('Failed to load branches:', err);
      },
    });
  }

  onSearchChange(term: string): void {
    this.searchTerm.set(term);
  }

  // Create/Edit Modal
  openCreateModal(): void {
    this.editingBranch.set(null);
    this.showBranchModal.set(true);
  }

  editBranch(branch: Branch): void {
    this.editingBranch.set(branch);
    this.showBranchModal.set(true);
  }

  closeBranchModal(): void {
    this.showBranchModal.set(false);
    this.editingBranch.set(null);
  }

  saveBranch(data: CreateBranchRequest): void {
    this.saving.set(true);

    const editing = this.editingBranch();
    const operation = editing
      ? this.branchService.updateBranch(editing.id, data)
      : this.branchService.createBranch(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Branch updated successfully' : 'Branch created successfully'
        );
        this.closeBranchModal();
        this.loadBranches();
        this.saving.set(false);
      },
      error: err => {
        this.toastService.error(
          editing ? 'Failed to update branch' : 'Failed to create branch'
        );
        this.saving.set(false);
        console.error('Failed to save branch:', err);
      },
    });
  }

  // Toggle Status
  toggleStatus(branch: Branch): void {
    const newStatus = !branch.isActive;
    this.branchService.setStatus(branch.id, newStatus).subscribe({
      next: () => {
        this.toastService.success(
          newStatus
            ? 'Branch activated successfully'
            : 'Branch deactivated successfully'
        );
        this.loadBranches();
      },
      error: err => {
        this.toastService.error('Failed to update branch status');
        console.error('Failed to toggle status:', err);
      },
    });
  }

  // Set Primary Modal
  setPrimary(branch: Branch): void {
    this.branchToSetPrimary.set(branch);
    this.showPrimaryModal.set(true);
  }

  closePrimaryModal(): void {
    this.showPrimaryModal.set(false);
    this.branchToSetPrimary.set(null);
  }

  confirmSetPrimary(): void {
    const branch = this.branchToSetPrimary();
    if (!branch) return;

    this.settingPrimary.set(true);

    this.branchService.setPrimary(branch.id).subscribe({
      next: () => {
        this.toastService.success('Primary branch updated successfully');
        this.closePrimaryModal();
        this.loadBranches();
        this.settingPrimary.set(false);
      },
      error: err => {
        this.toastService.error('Failed to set primary branch');
        this.settingPrimary.set(false);
        console.error('Failed to set primary:', err);
      },
    });
  }

  // Delete Modal
  confirmDelete(branch: Branch): void {
    this.branchToDelete.set(branch);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.branchToDelete.set(null);
  }

  deleteBranch(): void {
    const branch = this.branchToDelete();
    if (!branch) return;

    this.deleting.set(true);

    this.branchService.deleteBranch(branch.id).subscribe({
      next: () => {
        this.toastService.success('Branch deleted successfully');
        this.closeDeleteModal();
        this.loadBranches();
        this.deleting.set(false);
      },
      error: err => {
        this.toastService.error('Failed to delete branch');
        this.deleting.set(false);
        console.error('Failed to delete branch:', err);
      },
    });
  }
}
