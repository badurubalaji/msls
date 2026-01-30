import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

import { TimetableService } from '../timetable.service';
import { Timetable, TimetableFilter, TIMETABLE_STATUSES, TimetableStatus } from '../timetable.model';
import { ToastService } from '../../../../shared/services/toast.service';
import { BranchService } from '../../../admin/branches/branch.service';
import { Branch } from '../../../admin/branches/branch.model';
import { SectionService } from '../../services/section.service';
import { Section } from '../../academic.model';
import { AcademicYearService } from '../../../admin/academic-years/academic-year.service';
import { AcademicYear } from '../../../admin/academic-years/academic-year.model';

@Component({
  selector: 'msls-timetable-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-calendar-days"></i>
          </div>
          <div class="header-text">
            <h1>Timetables</h1>
            <p>Manage section timetables and schedules</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          <span>New Timetable</span>
        </button>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="filter-group">
          <select class="filter-select" [ngModel]="branchFilter()" (ngModelChange)="branchFilter.set($event); loadSections()">
            <option value="">All Branches</option>
            @for (branch of branches(); track branch.id) {
              <option [value]="branch.id">{{ branch.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <select class="filter-select" [ngModel]="sectionFilter()" (ngModelChange)="sectionFilter.set($event)">
            <option value="">All Sections</option>
            @for (section of filteredSections(); track section.id) {
              <option [value]="section.id">{{ section.className }} - {{ section.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <select class="filter-select" [ngModel]="academicYearFilter()" (ngModelChange)="academicYearFilter.set($event)">
            <option value="">All Academic Years</option>
            @for (year of academicYears(); track year.id) {
              <option [value]="year.id">{{ year.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <select class="filter-select" [ngModel]="statusFilter()" (ngModelChange)="statusFilter.set($event)">
            <option value="">All Status</option>
            @for (status of statuses; track status.value) {
              <option [value]="status.value">{{ status.label }}</option>
            }
          </select>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading timetables...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadTimetables()">Retry</button>
          </div>
        } @else {
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Timetable</th>
                  <th>Section</th>
                  <th>Academic Year</th>
                  <th style="width: 120px;">Status</th>
                  <th style="width: 140px;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (tt of filteredTimetables(); track tt.id) {
                  <tr>
                    <td>
                      <div class="name-wrapper">
                        <div class="tt-icon" [class]="getStatusClass(tt.status)">
                          <i class="fa-solid fa-calendar-days"></i>
                        </div>
                        <div class="name-content">
                          <span class="name">{{ tt.name }}</span>
                          @if (tt.description) {
                            <span class="description">{{ tt.description }}</span>
                          }
                        </div>
                      </div>
                    </td>
                    <td>
                      <span class="section-badge">{{ tt.className }} - {{ tt.sectionName }}</span>
                    </td>
                    <td>{{ tt.academicYearName }}</td>
                    <td>
                      <span class="status-badge" [class]="getStatusBadgeClass(tt.status)">
                        {{ getStatusLabel(tt.status) }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button class="action-btn" title="Edit" (click)="openBuilder(tt)">
                        <i class="fa-solid fa-pen-to-square"></i>
                      </button>
                      @if (tt.status === 'draft') {
                        <button class="action-btn action-btn--success" title="Publish" (click)="publishTimetable(tt)">
                          <i class="fa-solid fa-check-circle"></i>
                        </button>
                        <button class="action-btn action-btn--danger" title="Delete" (click)="confirmDelete(tt)">
                          <i class="fa-solid fa-trash"></i>
                        </button>
                      }
                      @if (tt.status === 'published') {
                        <button class="action-btn" title="Archive" (click)="archiveTimetable(tt)">
                          <i class="fa-solid fa-archive"></i>
                        </button>
                      }
                    </td>
                  </tr>
                } @empty {
                  <tr>
                    <td colspan="5" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-calendar"></i>
                        <p>No timetables found</p>
                        <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                          Create Timetable
                        </button>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        }
      </div>

      <!-- Create Modal -->
      @if (showCreateModal()) {
        <div class="modal-overlay" (click)="closeCreateModal()">
          <div class="modal" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3><i class="fa-solid fa-plus"></i> Create Timetable</h3>
              <button class="modal__close" (click)="closeCreateModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form class="form">
                <div class="form-group">
                  <label>Name <span class="required">*</span></label>
                  <input type="text" [(ngModel)]="formData.name" name="name" class="form-input" placeholder="e.g., Class 10A - Term 1 Timetable" />
                </div>
                <div class="form-row">
                  <div class="form-group">
                    <label>Branch <span class="required">*</span></label>
                    <select [(ngModel)]="formData.branchId" name="branchId" class="form-input" (ngModelChange)="onBranchChange()">
                      <option value="">Select Branch</option>
                      @for (branch of branches(); track branch.id) {
                        <option [value]="branch.id">{{ branch.name }}</option>
                      }
                    </select>
                  </div>
                  <div class="form-group">
                    <label>Section <span class="required">*</span></label>
                    <select [(ngModel)]="formData.sectionId" name="sectionId" class="form-input">
                      <option value="">Select Section</option>
                      @for (section of modalSections(); track section.id) {
                        <option [value]="section.id">{{ section.className }} - {{ section.name }}</option>
                      }
                    </select>
                  </div>
                </div>
                <div class="form-group">
                  <label>Academic Year <span class="required">*</span></label>
                  <select [(ngModel)]="formData.academicYearId" name="academicYearId" class="form-input">
                    <option value="">Select Academic Year</option>
                    @for (year of academicYears(); track year.id) {
                      <option [value]="year.id">{{ year.name }}</option>
                    }
                  </select>
                </div>
                <div class="form-group">
                  <label>Description</label>
                  <textarea [(ngModel)]="formData.description" name="description" class="form-input" rows="2" placeholder="Optional description..."></textarea>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              <button class="btn btn-secondary" (click)="closeCreateModal()">Cancel</button>
              <button class="btn btn-primary" [disabled]="saving()" (click)="createTimetable()">
                @if (saving()) {
                  <div class="btn-spinner"></div> Creating...
                } @else {
                  <i class="fa-solid fa-plus"></i> Create
                }
              </button>
            </div>
          </div>
        </div>
      }

      <!-- Delete Confirmation Modal -->
      @if (showDeleteModal()) {
        <div class="modal-overlay" (click)="closeDeleteModal()">
          <div class="modal modal--sm" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3><i class="fa-solid fa-trash"></i> Delete Timetable</h3>
              <button class="modal__close" (click)="closeDeleteModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <div class="delete-confirmation">
                <div class="delete-icon">
                  <i class="fa-solid fa-triangle-exclamation"></i>
                </div>
                <p>Are you sure you want to delete <strong>"{{ timetableToDelete()?.name }}"</strong>?</p>
                <p class="delete-warning">This action cannot be undone.</p>
              </div>
            </div>
            <div class="modal__footer">
              <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
              <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteTimetable()">
                @if (deleting()) {
                  <div class="btn-spinner"></div> Deleting...
                } @else {
                  <i class="fa-solid fa-trash"></i> Delete
                }
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1400px; margin: 0 auto; }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
      gap: 1rem;
    }

    .header-content { display: flex; align-items: center; gap: 1rem; }

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

    .header-text h1 { margin: 0; font-size: 1.5rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }

    .filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; }
    .filter-group { min-width: 180px; }

    .filter-select {
      width: 100%;
      padding: 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
    }

    .filter-select:focus { outline: none; border-color: #4f46e5; }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .table-container { overflow-x: auto; }

    .loading-container, .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-container { color: #dc2626; flex-direction: column; }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin { to { transform: rotate(360deg); } }

    .data-table { width: 100%; border-collapse: collapse; }

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

    .data-table tbody tr:hover { background: #f8fafc; }

    .name-wrapper { display: flex; align-items: center; gap: 0.75rem; }

    .tt-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .tt-icon.draft { background: #fef3c7; color: #d97706; }
    .tt-icon.published { background: #dcfce7; color: #16a34a; }
    .tt-icon.archived { background: #f1f5f9; color: #64748b; }

    .name-content { display: flex; flex-direction: column; }
    .name { font-weight: 500; color: #1e293b; }
    .description { font-size: 0.75rem; color: #64748b; }

    .section-badge {
      display: inline-flex;
      padding: 0.25rem 0.5rem;
      background: #eef2ff;
      color: #4f46e5;
      border-radius: 0.25rem;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .status-badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .status-badge.draft { background: #fef3c7; color: #92400e; }
    .status-badge.published { background: #dcfce7; color: #166534; }
    .status-badge.archived { background: #f1f5f9; color: #64748b; }

    .actions-cell { text-align: right; white-space: nowrap; }

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

    .action-btn:hover { background: #f1f5f9; color: #4f46e5; }
    .action-btn--success:hover { background: #dcfce7; color: #16a34a; }
    .action-btn--danger:hover { background: #fef2f2; color: #dc2626; }

    .empty-cell { padding: 3rem !important; }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #64748b;
    }

    .empty-state i { font-size: 2rem; color: #cbd5e1; }

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

    .btn-sm { padding: 0.375rem 0.75rem; font-size: 0.75rem; }
    .btn-primary { background: #4f46e5; color: white; }
    .btn-primary:hover:not(:disabled) { background: #4338ca; }
    .btn-secondary { background: #f1f5f9; color: #475569; }
    .btn-secondary:hover { background: #e2e8f0; }
    .btn-danger { background: #dc2626; color: white; }
    .btn-danger:hover:not(:disabled) { background: #b91c1c; }
    .btn:disabled { opacity: 0.5; cursor: not-allowed; }

    .btn-spinner {
      width: 16px;
      height: 16px;
      border: 2px solid transparent;
      border-top-color: currentColor;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    /* Modal Styles */
    .modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      padding: 1rem;
    }

    .modal {
      background: white;
      border-radius: 1rem;
      width: 100%;
      max-width: 32rem;
      max-height: 90vh;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
    }

    .modal--sm { max-width: 24rem; }

    .modal__header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid #f1f5f9;

      h3 {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        font-size: 1.125rem;
        font-weight: 600;
        color: #0f172a;
        margin: 0;

        i { color: #4f46e5; }
      }
    }

    .modal__close {
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f1f5f9;
      border: none;
      border-radius: 0.5rem;
      cursor: pointer;
      color: #64748b;
      transition: all 0.2s;

      &:hover { background: #e2e8f0; color: #334155; }
    }

    .modal__body { padding: 1.5rem; overflow-y: auto; flex: 1; }

    .modal__footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding: 1.25rem 1.5rem;
      background: #f8fafc;
      border-top: 1px solid #f1f5f9;
    }

    .form { display: flex; flex-direction: column; gap: 1rem; }
    .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-group label { font-size: 0.875rem; font-weight: 500; color: #374151; }
    .required { color: #dc2626; }

    .form-input {
      padding: 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      transition: border-color 0.2s, box-shadow 0.2s;
    }

    .form-input:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .delete-confirmation { text-align: center; padding: 1rem; }

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

    .delete-confirmation p { margin: 0 0 0.5rem; color: #374151; }
    .delete-warning { font-size: 0.875rem; color: #64748b; }

    @media (max-width: 768px) {
      .page { padding: 1rem; }
      .page-header { flex-direction: column; align-items: stretch; }
      .filters-bar { flex-direction: column; }
      .filter-group { min-width: 100%; }
      .form-row { grid-template-columns: 1fr; }
    }
  `]
})
export class TimetableListComponent implements OnInit {
  private readonly timetableService = inject(TimetableService);
  private readonly branchService = inject(BranchService);
  private readonly sectionService = inject(SectionService);
  private readonly academicYearService = inject(AcademicYearService);
  private readonly toastService = inject(ToastService);
  private readonly router = inject(Router);

  statuses = TIMETABLE_STATUSES;

  // Data
  timetables = signal<Timetable[]>([]);
  branches = signal<Branch[]>([]);
  sections = signal<Section[]>([]);
  academicYears = signal<AcademicYear[]>([]);
  modalSections = signal<Section[]>([]);

  // State
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);

  // Filters
  branchFilter = signal('');
  sectionFilter = signal('');
  academicYearFilter = signal('');
  statusFilter = signal('');

  // Modals
  showCreateModal = signal(false);
  showDeleteModal = signal(false);
  timetableToDelete = signal<Timetable | null>(null);

  // Form
  formData = {
    name: '',
    branchId: '',
    sectionId: '',
    academicYearId: '',
    description: ''
  };

  filteredTimetables = computed(() => {
    let result = this.timetables();
    const branch = this.branchFilter();
    const section = this.sectionFilter();
    const year = this.academicYearFilter();
    const status = this.statusFilter();

    if (branch) result = result.filter(t => t.branchId === branch);
    if (section) result = result.filter(t => t.sectionId === section);
    if (year) result = result.filter(t => t.academicYearId === year);
    if (status) result = result.filter(t => t.status === status);

    return result;
  });

  filteredSections = computed(() => {
    const branch = this.branchFilter();
    if (!branch) return this.sections();
    // Note: sections might need branch filtering if available
    return this.sections();
  });

  ngOnInit(): void {
    this.loadBranches();
    this.loadSections();
    this.loadAcademicYears();
    this.loadTimetables();
  }

  loadTimetables(): void {
    this.loading.set(true);
    this.error.set(null);

    this.timetableService.getTimetables().subscribe({
      next: timetables => {
        this.timetables.set(timetables);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load timetables');
        this.loading.set(false);
      }
    });
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: branches => this.branches.set(branches.filter(b => b.isActive)),
      error: () => console.error('Failed to load branches')
    });
  }

  loadSections(): void {
    this.sectionService.getSections({ isActive: true }).subscribe({
      next: sections => this.sections.set(sections),
      error: () => console.error('Failed to load sections')
    });
  }

  loadAcademicYears(): void {
    this.academicYearService.getAcademicYears().subscribe({
      next: years => this.academicYears.set(years),
      error: () => console.error('Failed to load academic years')
    });
  }

  getStatusClass(status: TimetableStatus): string {
    return status;
  }

  getStatusBadgeClass(status: TimetableStatus): string {
    return status;
  }

  getStatusLabel(status: TimetableStatus): string {
    return TIMETABLE_STATUSES.find(s => s.value === status)?.label || status;
  }

  openCreateModal(): void {
    this.formData = { name: '', branchId: '', sectionId: '', academicYearId: '', description: '' };
    this.modalSections.set([]);
    this.showCreateModal.set(true);
  }

  closeCreateModal(): void {
    this.showCreateModal.set(false);
  }

  onBranchChange(): void {
    const branchId = this.formData.branchId;
    if (!branchId) {
      this.modalSections.set([]);
      return;
    }
    // Load sections for the selected branch
    this.sectionService.getSections({ isActive: true }).subscribe({
      next: sections => this.modalSections.set(sections),
      error: () => this.modalSections.set([])
    });
  }

  createTimetable(): void {
    if (!this.formData.name || !this.formData.branchId || !this.formData.sectionId || !this.formData.academicYearId) {
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.saving.set(true);
    this.timetableService.createTimetable({
      name: this.formData.name,
      branchId: this.formData.branchId,
      sectionId: this.formData.sectionId,
      academicYearId: this.formData.academicYearId,
      description: this.formData.description || undefined
    }).subscribe({
      next: timetable => {
        this.toastService.success('Timetable created successfully');
        this.closeCreateModal();
        this.saving.set(false);
        // Navigate to builder
        this.router.navigate(['/academics/timetable/builder', timetable.id]);
      },
      error: () => {
        this.toastService.error('Failed to create timetable');
        this.saving.set(false);
      }
    });
  }

  openBuilder(timetable: Timetable): void {
    this.router.navigate(['/academics/timetable/builder', timetable.id]);
  }

  publishTimetable(timetable: Timetable): void {
    this.timetableService.publishTimetable(timetable.id).subscribe({
      next: () => {
        this.toastService.success('Timetable published successfully');
        this.loadTimetables();
      },
      error: () => this.toastService.error('Failed to publish timetable')
    });
  }

  archiveTimetable(timetable: Timetable): void {
    this.timetableService.archiveTimetable(timetable.id).subscribe({
      next: () => {
        this.toastService.success('Timetable archived successfully');
        this.loadTimetables();
      },
      error: () => this.toastService.error('Failed to archive timetable')
    });
  }

  confirmDelete(timetable: Timetable): void {
    this.timetableToDelete.set(timetable);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.timetableToDelete.set(null);
  }

  deleteTimetable(): void {
    const timetable = this.timetableToDelete();
    if (!timetable) return;

    this.deleting.set(true);
    this.timetableService.deleteTimetable(timetable.id).subscribe({
      next: () => {
        this.toastService.success('Timetable deleted successfully');
        this.closeDeleteModal();
        this.loadTimetables();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete timetable');
        this.deleting.set(false);
      }
    });
  }
}
