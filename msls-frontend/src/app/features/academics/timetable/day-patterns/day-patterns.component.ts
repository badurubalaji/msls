import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../../shared/components/modal/modal.component';
import { TimetableService } from '../timetable.service';
import { DayPattern, DayPatternAssignment, CreateDayPatternRequest, UpdateDayPatternRequest, DAYS_OF_WEEK } from '../timetable.model';
import { ToastService } from '../../../../shared/services/toast.service';
import { BranchService } from '../../../admin/branches/branch.service';
import { Branch } from '../../../admin/branches/branch.model';

@Component({
  selector: 'msls-day-patterns',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-calendar-days"></i>
          </div>
          <div class="header-text">
            <h1>Day Patterns</h1>
            <p>Configure different day schedules and assign them to weekdays</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Pattern
        </button>
      </div>

      <!-- Tabs -->
      <div class="tabs">
        <button class="tab" [class.active]="activeTab() === 'patterns'" (click)="activeTab.set('patterns')">
          <i class="fa-solid fa-layer-group"></i>
          Day Patterns
        </button>
        <button class="tab" [class.active]="activeTab() === 'assignments'" (click)="activeTab.set('assignments')">
          <i class="fa-solid fa-calendar-week"></i>
          Week Schedule
        </button>
      </div>

      <!-- Patterns Tab -->
      @if (activeTab() === 'patterns') {
        <!-- Filters -->
        <div class="filters-bar">
          <div class="search-box">
            <i class="fa-solid fa-search search-icon"></i>
            <input
              type="text"
              placeholder="Search patterns..."
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

        <!-- Patterns Table -->
        <div class="content-card">
          @if (loading()) {
            <div class="loading-container">
              <div class="spinner"></div>
              <span>Loading patterns...</span>
            </div>
          } @else if (error()) {
            <div class="error-container">
              <i class="fa-solid fa-circle-exclamation"></i>
              <span>{{ error() }}</span>
              <button class="btn btn-secondary btn-sm" (click)="loadPatterns()">Retry</button>
            </div>
          } @else {
            <table class="data-table">
              <thead>
                <tr>
                  <th>Pattern Name</th>
                  <th>Code</th>
                  <th>Total Periods</th>
                  <th>Description</th>
                  <th>Status</th>
                  <th style="width: 140px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (pattern of filteredPatterns(); track pattern.id) {
                  <tr>
                    <td>
                      <div class="name-wrapper">
                        <div class="pattern-icon">
                          <i class="fa-solid fa-calendar-day"></i>
                        </div>
                        <span class="name">{{ pattern.name }}</span>
                      </div>
                    </td>
                    <td><span class="code-badge">{{ pattern.code }}</span></td>
                    <td>
                      <span class="periods-badge">{{ pattern.totalPeriods }} periods</span>
                    </td>
                    <td class="description-cell">{{ pattern.description || '-' }}</td>
                    <td>
                      <span class="badge" [class.badge-green]="pattern.isActive" [class.badge-gray]="!pattern.isActive">
                        {{ pattern.isActive ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button class="action-btn" title="Edit" (click)="editPattern(pattern)">
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button class="action-btn" title="Toggle Status" (click)="toggleStatus(pattern)">
                        <i class="fa-solid" [class.fa-toggle-on]="pattern.isActive" [class.fa-toggle-off]="!pattern.isActive"></i>
                      </button>
                      <button class="action-btn action-btn--danger" title="Delete" (click)="confirmDelete(pattern)">
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="6" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No day patterns found</p>
                        <button class="btn btn-primary btn-sm" (click)="openCreateModal()">Add Pattern</button>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          }
        </div>
      }

      <!-- Assignments Tab -->
      @if (activeTab() === 'assignments') {
        <div class="filters-bar">
          <div class="filter-group">
            <label class="filter-label">Branch</label>
            <select
              class="filter-select"
              [ngModel]="selectedBranchId()"
              (ngModelChange)="onBranchChange($event)"
            >
              <option value="">Select Branch</option>
              @for (branch of branches(); track branch.id) {
                <option [value]="branch.id">{{ branch.name }}</option>
              }
            </select>
          </div>
        </div>

        <div class="content-card">
          @if (!selectedBranchId()) {
            <div class="empty-state-large">
              <i class="fa-solid fa-building"></i>
              <p>Select a branch to view and manage day assignments</p>
            </div>
          } @else if (loadingAssignments()) {
            <div class="loading-container">
              <div class="spinner"></div>
              <span>Loading assignments...</span>
            </div>
          } @else {
            <div class="week-grid">
              @for (day of daysOfWeek; track day.value) {
                <div class="day-card" [class.day-card--off]="!getDayAssignment(day.value)?.isWorkingDay">
                  <div class="day-header">
                    <span class="day-name">{{ day.label }}</span>
                    <label class="toggle-switch">
                      <input
                        type="checkbox"
                        [checked]="getDayAssignment(day.value)?.isWorkingDay ?? true"
                        (change)="toggleWorkingDay(day.value, $event)"
                      />
                      <span class="toggle-slider"></span>
                    </label>
                  </div>
                  <div class="day-content">
                    @if (getDayAssignment(day.value)?.isWorkingDay !== false) {
                      <select
                        class="day-pattern-select"
                        [ngModel]="getDayAssignment(day.value)?.dayPatternId || ''"
                        (ngModelChange)="updateDayPattern(day.value, $event)"
                      >
                        <option value="">No Pattern</option>
                        @for (pattern of patterns(); track pattern.id) {
                          @if (pattern.isActive) {
                            <option [value]="pattern.id">{{ pattern.name }}</option>
                          }
                        }
                      </select>
                      @if (getDayAssignment(day.value)?.dayPatternName) {
                        <span class="pattern-info">{{ getPatternPeriods(getDayAssignment(day.value)?.dayPatternId) }} periods</span>
                      }
                    } @else {
                      <span class="day-off-label">Non-working day</span>
                    }
                  </div>
                </div>
              }
            </div>
          }
        </div>
      }

      <!-- Pattern Form Modal -->
      <msls-modal [isOpen]="showFormModal()" [title]="editingPattern() ? 'Edit Day Pattern' : 'Create Day Pattern'" size="md" (closed)="closeFormModal()">
        <form class="form" (ngSubmit)="savePattern()">
          <div class="form-row">
            <div class="form-group">
              <label for="patternName">Pattern Name <span class="required">*</span></label>
              <input type="text" id="patternName" [(ngModel)]="formData.name" name="name" placeholder="e.g., Regular Day" required />
            </div>
            <div class="form-group">
              <label for="patternCode">Code <span class="required">*</span></label>
              <input type="text" id="patternCode" [(ngModel)]="formData.code" name="code" placeholder="e.g., REG" required />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="totalPeriods">Total Periods</label>
              <input type="number" id="totalPeriods" [(ngModel)]="formData.totalPeriods" name="totalPeriods" min="1" max="20" />
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
                {{ editingPattern() ? 'Update' : 'Create' }}
              }
            </button>
          </div>
        </form>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal [isOpen]="showDeleteModal()" title="Delete Day Pattern" size="sm" (closed)="closeDeleteModal()">
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>Are you sure you want to delete <strong>"{{ patternToDelete()?.name }}"</strong>?</p>
          <p class="delete-warning">Patterns in use cannot be deleted.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deletePattern()">
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
    .header-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; background: #e0e7ff; color: #4f46e5; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
    .header-text h1 { margin: 0; font-size: 1.5rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }
    .tabs { display: flex; gap: 0.5rem; margin-bottom: 1rem; }
    .tab { display: inline-flex; align-items: center; gap: 0.5rem; padding: 0.625rem 1rem; border: 1px solid #e2e8f0; background: white; border-radius: 0.5rem; font-size: 0.875rem; font-weight: 500; color: #64748b; cursor: pointer; transition: all 0.2s; }
    .tab:hover { background: #f8fafc; }
    .tab.active { background: #4f46e5; color: white; border-color: #4f46e5; }
    .filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; align-items: flex-end; }
    .filter-label { display: block; font-size: 0.75rem; font-weight: 500; color: #64748b; margin-bottom: 0.25rem; }
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
    .pattern-icon { width: 2.5rem; height: 2.5rem; border-radius: 0.5rem; background: #e0e7ff; color: #4f46e5; display: flex; align-items: center; justify-content: center; font-size: 1rem; }
    .name { font-weight: 500; color: #1e293b; }
    .code-badge { display: inline-flex; padding: 0.25rem 0.5rem; background: #f1f5f9; border-radius: 0.25rem; font-size: 0.75rem; font-weight: 600; font-family: monospace; }
    .periods-badge { display: inline-flex; padding: 0.25rem 0.5rem; background: #dbeafe; color: #1d4ed8; border-radius: 0.25rem; font-size: 0.75rem; font-weight: 500; }
    .description-cell { max-width: 200px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: #64748b; }
    .badge { display: inline-flex; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 500; }
    .badge-green { background: #dcfce7; color: #166534; }
    .badge-gray { background: #f1f5f9; color: #64748b; }
    .actions-cell { text-align: right; }
    .action-btn { display: inline-flex; align-items: center; justify-content: center; width: 2rem; height: 2rem; border: none; background: transparent; color: #64748b; border-radius: 0.375rem; cursor: pointer; transition: all 0.2s; }
    .action-btn:hover { background: #f1f5f9; color: #4f46e5; }
    .action-btn--danger:hover { background: #fef2f2; color: #dc2626; }
    .empty-cell { padding: 3rem !important; }
    .empty-state, .empty-state-large { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; color: #64748b; padding: 3rem; }
    .empty-state-large i { font-size: 2rem; }
    .week-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 0.75rem; padding: 1.5rem; }
    @media (max-width: 1200px) { .week-grid { grid-template-columns: repeat(4, 1fr); } }
    @media (max-width: 768px) { .week-grid { grid-template-columns: repeat(2, 1fr); } }
    .day-card { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 0.75rem; padding: 1rem; }
    .day-card--off { background: #fef2f2; border-color: #fecaca; }
    .day-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem; }
    .day-name { font-weight: 600; color: #1e293b; }
    .toggle-switch { position: relative; display: inline-block; width: 40px; height: 22px; }
    .toggle-switch input { opacity: 0; width: 0; height: 0; }
    .toggle-slider { position: absolute; cursor: pointer; inset: 0; background: #cbd5e1; border-radius: 22px; transition: 0.2s; }
    .toggle-slider::before { content: ''; position: absolute; height: 16px; width: 16px; left: 3px; bottom: 3px; background: white; border-radius: 50%; transition: 0.2s; }
    .toggle-switch input:checked + .toggle-slider { background: #4f46e5; }
    .toggle-switch input:checked + .toggle-slider::before { transform: translateX(18px); }
    .day-content { display: flex; flex-direction: column; gap: 0.5rem; }
    .day-pattern-select { width: 100%; padding: 0.5rem; border: 1px solid #e2e8f0; border-radius: 0.375rem; font-size: 0.8125rem; }
    .pattern-info { font-size: 0.75rem; color: #64748b; }
    .day-off-label { font-size: 0.8125rem; color: #dc2626; font-style: italic; }
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
export class DayPatternsComponent implements OnInit {
  private timetableService = inject(TimetableService);
  private branchService = inject(BranchService);
  private toastService = inject(ToastService);

  patterns = signal<DayPattern[]>([]);
  assignments = signal<DayPatternAssignment[]>([]);
  branches = signal<Branch[]>([]);

  loading = signal(true);
  loadingAssignments = signal(false);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');
  activeTab = signal<'patterns' | 'assignments'>('patterns');
  selectedBranchId = signal<string>('');

  showFormModal = signal(false);
  showDeleteModal = signal(false);
  editingPattern = signal<DayPattern | null>(null);
  patternToDelete = signal<DayPattern | null>(null);

  formData = { name: '', code: '', description: '', totalPeriods: 8, displayOrder: 0 };

  daysOfWeek = DAYS_OF_WEEK;

  filteredPatterns = computed(() => {
    let result = this.patterns();
    const term = this.searchTerm().toLowerCase();
    if (term) result = result.filter(p => p.name.toLowerCase().includes(term) || p.code.toLowerCase().includes(term));
    if (this.statusFilter() === 'active') result = result.filter(p => p.isActive);
    else if (this.statusFilter() === 'inactive') result = result.filter(p => !p.isActive);
    return result;
  });

  ngOnInit(): void {
    this.loadPatterns();
    this.loadBranches();
  }

  loadPatterns(): void {
    this.loading.set(true);
    this.timetableService.getDayPatterns().subscribe({
      next: patterns => { this.patterns.set(patterns); this.loading.set(false); },
      error: () => { this.error.set('Failed to load day patterns'); this.loading.set(false); },
    });
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: allBranches => {
        const branches = allBranches.filter(b => b.isActive);
        this.branches.set(branches);
        if (branches.length === 1) {
          this.selectedBranchId.set(branches[0].id);
          this.loadAssignments(branches[0].id);
        }
      },
      error: () => console.error('Failed to load branches'),
    });
  }

  onBranchChange(branchId: string): void {
    this.selectedBranchId.set(branchId);
    if (branchId) this.loadAssignments(branchId);
  }

  loadAssignments(branchId: string): void {
    this.loadingAssignments.set(true);
    this.timetableService.getDayPatternAssignments(branchId).subscribe({
      next: assignments => { this.assignments.set(assignments); this.loadingAssignments.set(false); },
      error: () => { this.toastService.error('Failed to load assignments'); this.loadingAssignments.set(false); },
    });
  }

  getDayAssignment(dayOfWeek: number): DayPatternAssignment | undefined {
    return this.assignments().find(a => a.dayOfWeek === dayOfWeek);
  }

  getPatternPeriods(patternId: string | undefined): number {
    if (!patternId) return 0;
    const pattern = this.patterns().find(p => p.id === patternId);
    return pattern?.totalPeriods || 0;
  }

  toggleWorkingDay(dayOfWeek: number, event: Event): void {
    const isWorkingDay = (event.target as HTMLInputElement).checked;
    const branchId = this.selectedBranchId();
    if (!branchId) return;

    this.timetableService.updateDayPatternAssignment(branchId, dayOfWeek, { isWorkingDay }).subscribe({
      next: () => {
        this.toastService.success('Updated');
        this.loadAssignments(branchId);
      },
      error: () => this.toastService.error('Failed to update'),
    });
  }

  updateDayPattern(dayOfWeek: number, patternId: string): void {
    const branchId = this.selectedBranchId();
    if (!branchId) return;

    this.timetableService.updateDayPatternAssignment(branchId, dayOfWeek, {
      dayPatternId: patternId || null,
    }).subscribe({
      next: () => {
        this.toastService.success('Updated');
        this.loadAssignments(branchId);
      },
      error: () => this.toastService.error('Failed to update'),
    });
  }

  openCreateModal(): void {
    this.editingPattern.set(null);
    this.formData = { name: '', code: '', description: '', totalPeriods: 8, displayOrder: 0 };
    this.showFormModal.set(true);
  }

  editPattern(pattern: DayPattern): void {
    this.editingPattern.set(pattern);
    this.formData = {
      name: pattern.name,
      code: pattern.code,
      description: pattern.description || '',
      totalPeriods: pattern.totalPeriods,
      displayOrder: pattern.displayOrder,
    };
    this.showFormModal.set(true);
  }

  closeFormModal(): void {
    this.showFormModal.set(false);
    this.editingPattern.set(null);
  }

  savePattern(): void {
    if (!this.formData.name || !this.formData.code) {
      this.toastService.error('Please fill in all required fields');
      return;
    }
    this.saving.set(true);
    const editing = this.editingPattern();

    if (editing) {
      const data: UpdateDayPatternRequest = {
        name: this.formData.name,
        code: this.formData.code,
        description: this.formData.description || undefined,
        totalPeriods: this.formData.totalPeriods,
        displayOrder: this.formData.displayOrder,
      };
      this.timetableService.updateDayPattern(editing.id, data).subscribe({
        next: () => {
          this.toastService.success('Pattern updated');
          this.closeFormModal();
          this.loadPatterns();
          this.saving.set(false);
        },
        error: () => {
          this.toastService.error('Failed to update pattern');
          this.saving.set(false);
        },
      });
    } else {
      const data: CreateDayPatternRequest = {
        name: this.formData.name,
        code: this.formData.code,
        description: this.formData.description || undefined,
        totalPeriods: this.formData.totalPeriods,
        displayOrder: this.formData.displayOrder,
      };
      this.timetableService.createDayPattern(data).subscribe({
        next: () => {
          this.toastService.success('Pattern created');
          this.closeFormModal();
          this.loadPatterns();
          this.saving.set(false);
        },
        error: () => {
          this.toastService.error('Failed to create pattern');
          this.saving.set(false);
        },
      });
    }
  }

  toggleStatus(pattern: DayPattern): void {
    this.timetableService.updateDayPattern(pattern.id, { isActive: !pattern.isActive }).subscribe({
      next: () => { this.toastService.success('Status updated'); this.loadPatterns(); },
      error: () => this.toastService.error('Failed to update status'),
    });
  }

  confirmDelete(pattern: DayPattern): void {
    this.patternToDelete.set(pattern);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.patternToDelete.set(null);
  }

  deletePattern(): void {
    const pattern = this.patternToDelete();
    if (!pattern) return;
    this.deleting.set(true);
    this.timetableService.deleteDayPattern(pattern.id).subscribe({
      next: () => {
        this.toastService.success('Pattern deleted');
        this.closeDeleteModal();
        this.loadPatterns();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete pattern (may be in use)');
        this.deleting.set(false);
      },
    });
  }
}
