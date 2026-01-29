import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../../shared/components/modal/modal.component';
import { TimetableService } from '../timetable.service';
import { Shift, CreateShiftRequest, UpdateShiftRequest } from '../timetable.model';
import { ToastService } from '../../../../shared/services/toast.service';
import { BranchService } from '../../../admin/branches/branch.service';
import { Branch } from '../../../admin/branches/branch.model';

@Component({
  selector: 'msls-shifts',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-clock-rotate-left"></i>
          </div>
          <div class="header-text">
            <h1>Shifts</h1>
            <p>Configure school shift timings for different batches</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Shift
        </button>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search shifts..."
            [ngModel]="searchTerm()"
            (ngModelChange)="searchTerm.set($event)"
            class="search-input"
          />
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
            <span>Loading shifts...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadShifts()">Retry</button>
          </div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Shift</th>
                <th>Code</th>
                <th>Timing</th>
                <th>Duration</th>
                <th>Status</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (shift of filteredShifts(); track shift.id) {
                <tr>
                  <td>
                    <div class="name-wrapper">
                      <div class="shift-icon">
                        <i class="fa-solid fa-sun" [class.afternoon]="isAfternoon(shift.startTime)"></i>
                      </div>
                      <div class="name-content">
                        <span class="name">{{ shift.name }}</span>
                        @if (shift.description) {
                          <span class="description">{{ shift.description }}</span>
                        }
                      </div>
                    </div>
                  </td>
                  <td><span class="code-badge">{{ shift.code }}</span></td>
                  <td>
                    <div class="time-display">
                      <span class="time-start">{{ formatTime(shift.startTime) }}</span>
                      <i class="fa-solid fa-arrow-right time-arrow"></i>
                      <span class="time-end">{{ formatTime(shift.endTime) }}</span>
                    </div>
                  </td>
                  <td>{{ calculateDuration(shift.startTime, shift.endTime) }}</td>
                  <td>
                    <span class="badge" [class.badge-green]="shift.isActive" [class.badge-gray]="!shift.isActive">
                      {{ shift.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                  <td class="actions-cell">
                    <button class="action-btn" title="Edit" (click)="editShift(shift)">
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button class="action-btn" title="Toggle Status" (click)="toggleStatus(shift)">
                      <i class="fa-solid" [class.fa-toggle-on]="shift.isActive" [class.fa-toggle-off]="!shift.isActive"></i>
                    </button>
                    <button class="action-btn action-btn--danger" title="Delete" (click)="confirmDelete(shift)">
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="6" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-folder-open"></i>
                      <p>No shifts found</p>
                      <button class="btn btn-primary btn-sm" (click)="openCreateModal()">Add Shift</button>
                    </div>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        }
      </div>

      <!-- Shift Form Modal -->
      <msls-modal [isOpen]="showFormModal()" [title]="editingShift() ? 'Edit Shift' : 'Create Shift'" size="md" (closed)="closeFormModal()">
        <form class="form" (ngSubmit)="saveShift()">
          <div class="form-row">
            <div class="form-group">
              <label for="shiftName">Shift Name <span class="required">*</span></label>
              <input type="text" id="shiftName" [(ngModel)]="formData.name" name="name" placeholder="e.g., Morning Shift" required />
            </div>
            <div class="form-group">
              <label for="shiftCode">Code <span class="required">*</span></label>
              <input type="text" id="shiftCode" [(ngModel)]="formData.code" name="code" placeholder="e.g., MOR" required />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="startTime">Start Time <span class="required">*</span></label>
              <input type="time" id="startTime" [(ngModel)]="formData.startTime" name="startTime" required />
            </div>
            <div class="form-group">
              <label for="endTime">End Time <span class="required">*</span></label>
              <input type="time" id="endTime" [(ngModel)]="formData.endTime" name="endTime" required />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="branchId">Branch <span class="required">*</span></label>
              <select id="branchId" [(ngModel)]="formData.branchId" name="branchId" required [disabled]="!!editingShift()">
                <option value="">Select Branch</option>
                @for (branch of branches(); track branch.id) {
                  <option [value]="branch.id">{{ branch.name }}</option>
                }
              </select>
            </div>
            <div class="form-group">
              <label for="displayOrder">Display Order</label>
              <input type="number" id="displayOrder" [(ngModel)]="formData.displayOrder" name="displayOrder" min="0" />
            </div>
          </div>

          <div class="form-group">
            <label for="description">Description</label>
            <textarea id="description" [(ngModel)]="formData.description" name="description" rows="2" placeholder="Optional description"></textarea>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" (click)="closeFormModal()">Cancel</button>
            <button type="submit" class="btn btn-primary" [disabled]="saving()">
              @if (saving()) {
                <div class="btn-spinner"></div>
                Saving...
              } @else {
                {{ editingShift() ? 'Update' : 'Create' }}
              }
            </button>
          </div>
        </form>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal [isOpen]="showDeleteModal()" title="Delete Shift" size="sm" (closed)="closeDeleteModal()">
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>Are you sure you want to delete <strong>"{{ shiftToDelete()?.name }}"</strong>?</p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteShift()">
              @if (deleting()) { Deleting... } @else { Delete }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1400px; margin: 0 auto; }
    .page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
    .header-content { display: flex; align-items: center; gap: 1rem; }
    .header-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; background: #dbeafe; color: #2563eb; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
    .header-text h1 { margin: 0; font-size: 1.5rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }
    .filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; }
    .search-box { flex: 1; max-width: 400px; position: relative; }
    .search-icon { position: absolute; left: 0.875rem; top: 50%; transform: translateY(-50%); color: #9ca3af; }
    .search-input { width: 100%; padding: 0.625rem 2.5rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; }
    .filter-select { padding: 0.625rem 2rem 0.625rem 0.875rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; background: white; }
    .content-card { background: white; border: 1px solid #e2e8f0; border-radius: 1rem; overflow: hidden; }
    .loading-container, .error-container { display: flex; align-items: center; justify-content: center; gap: 1rem; padding: 3rem; color: #64748b; }
    .spinner { width: 24px; height: 24px; border: 3px solid #e2e8f0; border-top-color: #4f46e5; border-radius: 50%; animation: spin 0.8s linear infinite; }
    @keyframes spin { to { transform: rotate(360deg); } }
    .data-table { width: 100%; border-collapse: collapse; }
    .data-table th { text-align: left; padding: 0.875rem 1rem; font-size: 0.75rem; font-weight: 600; text-transform: uppercase; color: #64748b; background: #f8fafc; border-bottom: 1px solid #e2e8f0; }
    .data-table td { padding: 1rem; border-bottom: 1px solid #f1f5f9; color: #374151; }
    .data-table tbody tr:hover { background: #f8fafc; }
    .name-wrapper { display: flex; align-items: center; gap: 0.75rem; }
    .shift-icon { width: 2.5rem; height: 2.5rem; border-radius: 0.5rem; background: #fef3c7; color: #d97706; display: flex; align-items: center; justify-content: center; font-size: 1rem; }
    .shift-icon .afternoon { color: #7c3aed; }
    .name-content { display: flex; flex-direction: column; }
    .name { font-weight: 500; color: #1e293b; }
    .description { font-size: 0.75rem; color: #64748b; }
    .code-badge { display: inline-flex; padding: 0.25rem 0.5rem; background: #f1f5f9; border-radius: 0.25rem; font-size: 0.75rem; font-weight: 600; font-family: monospace; }
    .time-display { display: flex; align-items: center; gap: 0.5rem; }
    .time-start, .time-end { font-weight: 500; }
    .time-arrow { font-size: 0.625rem; color: #9ca3af; }
    .badge { display: inline-flex; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 500; }
    .badge-green { background: #dcfce7; color: #166534; }
    .badge-gray { background: #f1f5f9; color: #64748b; }
    .actions-cell { text-align: right; }
    .action-btn { display: inline-flex; align-items: center; justify-content: center; width: 2rem; height: 2rem; border: none; background: transparent; color: #64748b; border-radius: 0.375rem; cursor: pointer; transition: all 0.2s; }
    .action-btn:hover { background: #f1f5f9; color: #4f46e5; }
    .action-btn--danger:hover { background: #fef2f2; color: #dc2626; }
    .empty-cell { padding: 3rem !important; }
    .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; color: #64748b; }
    .btn { display: inline-flex; align-items: center; gap: 0.5rem; padding: 0.625rem 1.25rem; border-radius: 0.5rem; font-size: 0.875rem; font-weight: 500; cursor: pointer; border: none; }
    .btn-primary { background: #4f46e5; color: white; }
    .btn-primary:hover:not(:disabled) { background: #4338ca; }
    .btn-secondary { background: #f1f5f9; color: #475569; }
    .btn-danger { background: #dc2626; color: white; }
    .btn-sm { padding: 0.375rem 0.75rem; font-size: 0.8125rem; }
    .btn:disabled { opacity: 0.5; cursor: not-allowed; }
    .btn-spinner { width: 16px; height: 16px; border: 2px solid transparent; border-top-color: currentColor; border-radius: 50%; animation: spin 0.8s linear infinite; }
    .form { display: flex; flex-direction: column; gap: 1rem; }
    .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-group label { font-size: 0.875rem; font-weight: 500; color: #374151; }
    .required { color: #dc2626; }
    .form-group input, .form-group select, .form-group textarea { padding: 0.625rem 0.875rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; }
    .form-group input:focus, .form-group select:focus, .form-group textarea:focus { outline: none; border-color: #4f46e5; }
    .form-actions { display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 0.5rem; padding-top: 1rem; border-top: 1px solid #e2e8f0; }
    .delete-confirmation { text-align: center; padding: 1rem; }
    .delete-icon { width: 4rem; height: 4rem; margin: 0 auto 1rem; border-radius: 50%; background: #fef2f2; color: #dc2626; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
    .delete-warning { font-size: 0.875rem; color: #64748b; }
    .delete-actions { display: flex; gap: 0.75rem; justify-content: center; margin-top: 1.5rem; }
  `],
})
export class ShiftsComponent implements OnInit {
  private timetableService = inject(TimetableService);
  private branchService = inject(BranchService);
  private toastService = inject(ToastService);

  shifts = signal<Shift[]>([]);
  branches = signal<Branch[]>([]);

  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');

  showFormModal = signal(false);
  showDeleteModal = signal(false);
  editingShift = signal<Shift | null>(null);
  shiftToDelete = signal<Shift | null>(null);

  formData = { name: '', code: '', branchId: '', startTime: '08:00', endTime: '14:00', description: '', displayOrder: 0 };

  filteredShifts = computed(() => {
    let result = this.shifts();
    const term = this.searchTerm().toLowerCase();
    if (term) result = result.filter(s => s.name.toLowerCase().includes(term) || s.code.toLowerCase().includes(term));
    if (this.statusFilter() === 'active') result = result.filter(s => s.isActive);
    else if (this.statusFilter() === 'inactive') result = result.filter(s => !s.isActive);
    return result;
  });

  ngOnInit(): void {
    this.loadShifts();
    this.loadBranches();
  }

  loadShifts(): void {
    this.loading.set(true);
    this.timetableService.getShifts().subscribe({
      next: shifts => { this.shifts.set(shifts); this.loading.set(false); },
      error: () => { this.error.set('Failed to load shifts'); this.loading.set(false); },
    });
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: branches => this.branches.set(branches.filter(b => b.isActive)),
      error: () => console.error('Failed to load branches'),
    });
  }

  openCreateModal(): void {
    this.editingShift.set(null);
    this.formData = { name: '', code: '', branchId: '', startTime: '08:00', endTime: '14:00', description: '', displayOrder: 0 };
    this.showFormModal.set(true);
  }

  editShift(shift: Shift): void {
    this.editingShift.set(shift);
    this.formData = {
      name: shift.name,
      code: shift.code,
      branchId: shift.branchId,
      startTime: shift.startTime,
      endTime: shift.endTime,
      description: shift.description || '',
      displayOrder: shift.displayOrder,
    };
    this.showFormModal.set(true);
  }

  closeFormModal(): void {
    this.showFormModal.set(false);
    this.editingShift.set(null);
  }

  saveShift(): void {
    if (!this.formData.name || !this.formData.code || !this.formData.branchId || !this.formData.startTime || !this.formData.endTime) {
      this.toastService.error('Please fill in all required fields');
      return;
    }
    this.saving.set(true);
    const editing = this.editingShift();

    if (editing) {
      const data: UpdateShiftRequest = {
        name: this.formData.name,
        code: this.formData.code,
        startTime: this.formData.startTime,
        endTime: this.formData.endTime,
        description: this.formData.description || undefined,
        displayOrder: this.formData.displayOrder,
      };
      this.timetableService.updateShift(editing.id, data).subscribe({
        next: () => {
          this.toastService.success('Shift updated');
          this.closeFormModal();
          this.loadShifts();
          this.saving.set(false);
        },
        error: () => {
          this.toastService.error('Failed to update shift');
          this.saving.set(false);
        },
      });
    } else {
      const data: CreateShiftRequest = {
        branchId: this.formData.branchId,
        name: this.formData.name,
        code: this.formData.code,
        startTime: this.formData.startTime,
        endTime: this.formData.endTime,
        description: this.formData.description || undefined,
        displayOrder: this.formData.displayOrder,
      };
      this.timetableService.createShift(data).subscribe({
        next: () => {
          this.toastService.success('Shift created');
          this.closeFormModal();
          this.loadShifts();
          this.saving.set(false);
        },
        error: () => {
          this.toastService.error('Failed to create shift');
          this.saving.set(false);
        },
      });
    }
  }

  toggleStatus(shift: Shift): void {
    this.timetableService.updateShift(shift.id, { isActive: !shift.isActive }).subscribe({
      next: () => { this.toastService.success('Status updated'); this.loadShifts(); },
      error: () => this.toastService.error('Failed to update status'),
    });
  }

  confirmDelete(shift: Shift): void {
    this.shiftToDelete.set(shift);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.shiftToDelete.set(null);
  }

  deleteShift(): void {
    const shift = this.shiftToDelete();
    if (!shift) return;
    this.deleting.set(true);
    this.timetableService.deleteShift(shift.id).subscribe({
      next: () => {
        this.toastService.success('Shift deleted');
        this.closeDeleteModal();
        this.loadShifts();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete shift');
        this.deleting.set(false);
      },
    });
  }

  formatTime(time: string): string {
    if (!time) return '-';
    const [hours, minutes] = time.split(':');
    const h = parseInt(hours, 10);
    const ampm = h >= 12 ? 'PM' : 'AM';
    const h12 = h % 12 || 12;
    return `${h12}:${minutes} ${ampm}`;
  }

  isAfternoon(time: string): boolean {
    if (!time) return false;
    const hours = parseInt(time.split(':')[0], 10);
    return hours >= 12;
  }

  calculateDuration(startTime: string, endTime: string): string {
    if (!startTime || !endTime) return '-';
    const [sh, sm] = startTime.split(':').map(Number);
    const [eh, em] = endTime.split(':').map(Number);
    let minutes = (eh * 60 + em) - (sh * 60 + sm);
    if (minutes < 0) minutes += 24 * 60;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
  }
}
