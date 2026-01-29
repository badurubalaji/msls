import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { SectionService } from '../services/section.service';
import { ClassService } from '../services/class.service';
import { StreamService } from '../services/stream.service';
import { Section, CreateSectionRequest, UpdateSectionRequest, Class, Stream } from '../academic.model';
import { ToastService } from '../../../shared/services/toast.service';
import { AcademicYearService } from '../../admin/academic-years/academic-year.service';
import { AcademicYear } from '../../admin/academic-years/academic-year.model';

@Component({
  selector: 'msls-sections',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-layer-group"></i>
          </div>
          <div class="header-text">
            <h1>Sections</h1>
            <p>Manage class sections and assign class teachers</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Section
        </button>
      </div>

      <!-- Search & Filters -->
      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input
            type="text"
            placeholder="Search sections..."
            [ngModel]="searchTerm()"
            (ngModelChange)="searchTerm.set($event)"
            class="search-input"
          />
        </div>
        <div class="filter-group">
          <select
            class="filter-select"
            [ngModel]="classFilter()"
            (ngModelChange)="classFilter.set($event)"
          >
            <option value="">All Classes</option>
            @for (cls of classes(); track cls.id) {
              <option [value]="cls.id">{{ cls.name }}</option>
            }
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
            <span>Loading sections...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadSections()">Retry</button>
          </div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Section</th>
                <th>Class</th>
                <th>Capacity</th>
                <th>Students</th>
                <th>Class Teacher</th>
                <th>Status</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (section of filteredSections(); track section.id) {
                <tr>
                  <td>
                    <div class="name-wrapper">
                      <div class="section-icon">{{ section.code }}</div>
                      <div class="name-content">
                        <span class="name">{{ section.name }}</span>
                        @if (section.roomNumber) {
                          <span class="description">Room: {{ section.roomNumber }}</span>
                        }
                      </div>
                    </div>
                  </td>
                  <td>{{ section.className || '-' }}</td>
                  <td>
                    <div class="capacity-bar">
                      <div class="capacity-fill" [style.width.%]="getCapacityPercent(section)"></div>
                    </div>
                    <span class="capacity-text">{{ section.studentCount }}/{{ section.capacity }}</span>
                  </td>
                  <td class="students-cell">
                    <span class="count-badge">{{ section.studentCount }}</span>
                  </td>
                  <td>{{ section.classTeacherName || '-' }}</td>
                  <td>
                    <span class="badge" [class.badge-green]="section.isActive" [class.badge-gray]="!section.isActive">
                      {{ section.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                  <td class="actions-cell">
                    <button class="action-btn" title="Edit" (click)="editSection(section)">
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button class="action-btn" title="Toggle Status" (click)="toggleStatus(section)">
                      <i class="fa-solid" [class.fa-toggle-on]="section.isActive" [class.fa-toggle-off]="!section.isActive"></i>
                    </button>
                    <button class="action-btn action-btn--danger" title="Delete" (click)="confirmDelete(section)">
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </td>
                </tr>
              } @empty {
                <tr>
                  <td colspan="7" class="empty-cell">
                    <div class="empty-state">
                      <i class="fa-regular fa-folder-open"></i>
                      <p>No sections found</p>
                    </div>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        }
      </div>

      <!-- Section Form Modal -->
      <msls-modal [isOpen]="showSectionModal()" [title]="editingSection() ? 'Edit Section' : 'Create Section'" size="lg" (closed)="closeSectionModal()">
        <form class="form" (ngSubmit)="saveSection()">
          <div class="form-row">
            <div class="form-group">
              <label for="sectionName">Section Name <span class="required">*</span></label>
              <input type="text" id="sectionName" [(ngModel)]="formData.name" name="name" placeholder="e.g., A" required />
            </div>
            <div class="form-group">
              <label for="sectionCode">Code <span class="required">*</span></label>
              <input type="text" id="sectionCode" [(ngModel)]="formData.code" name="code" placeholder="e.g., XA" required />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="classId">Class <span class="required">*</span></label>
              <select id="classId" [(ngModel)]="formData.classId" name="classId" required [disabled]="!!editingSection()">
                <option value="">Select Class</option>
                @for (cls of classes(); track cls.id) {
                  <option [value]="cls.id">{{ cls.name }}</option>
                }
              </select>
            </div>
            <div class="form-group">
              <label for="capacity">Capacity</label>
              <input type="number" id="capacity" [(ngModel)]="formData.capacity" name="capacity" min="1" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="academicYearId">Academic Year</label>
              <select id="academicYearId" [(ngModel)]="formData.academicYearId" name="academicYearId">
                <option value="">Select Academic Year</option>
                @for (year of academicYears(); track year.id) {
                  <option [value]="year.id">{{ year.name }}</option>
                }
              </select>
            </div>
            <div class="form-group">
              <label for="roomNumber">Room Number</label>
              <input type="text" id="roomNumber" [(ngModel)]="formData.roomNumber" name="roomNumber" placeholder="e.g., 101" />
            </div>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" (click)="closeSectionModal()">Cancel</button>
            <button type="submit" class="btn btn-primary" [disabled]="saving()">
              @if (saving()) {
                <div class="btn-spinner"></div>
                Saving...
              } @else {
                {{ editingSection() ? 'Update' : 'Create' }}
              }
            </button>
          </div>
        </form>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal [isOpen]="showDeleteModal()" title="Delete Section" size="sm" (closed)="closeDeleteModal()">
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>Are you sure you want to delete <strong>"{{ sectionToDelete()?.name }}"</strong>?</p>
          <p class="delete-warning">Sections with students cannot be deleted.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteSection()">
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
    .header-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; background: #fef3c7; color: #d97706; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
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
    .section-icon { width: 2.5rem; height: 2.5rem; border-radius: 0.5rem; background: #fef3c7; color: #d97706; display: flex; align-items: center; justify-content: center; font-weight: 600; }
    .name-content { display: flex; flex-direction: column; }
    .name { font-weight: 500; color: #1e293b; }
    .description { font-size: 0.75rem; color: #64748b; }
    .capacity-bar { width: 60px; height: 6px; background: #e2e8f0; border-radius: 3px; overflow: hidden; }
    .capacity-fill { height: 100%; background: #4f46e5; transition: width 0.3s; }
    .capacity-text { font-size: 0.75rem; color: #64748b; margin-left: 0.5rem; }
    .count-badge { display: inline-flex; padding: 0.25rem 0.5rem; background: #f1f5f9; border-radius: 9999px; font-size: 0.75rem; font-weight: 500; }
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
    .btn:disabled { opacity: 0.5; cursor: not-allowed; }
    .btn-spinner { width: 16px; height: 16px; border: 2px solid transparent; border-top-color: currentColor; border-radius: 50%; animation: spin 0.8s linear infinite; }
    .form { display: flex; flex-direction: column; gap: 1rem; }
    .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-group label { font-size: 0.875rem; font-weight: 500; color: #374151; }
    .required { color: #dc2626; }
    .form-group input, .form-group select { padding: 0.625rem 0.875rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; }
    .form-group input:focus, .form-group select:focus { outline: none; border-color: #4f46e5; }
    .form-actions { display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 0.5rem; padding-top: 1rem; border-top: 1px solid #e2e8f0; }
    .delete-confirmation { text-align: center; padding: 1rem; }
    .delete-icon { width: 4rem; height: 4rem; margin: 0 auto 1rem; border-radius: 50%; background: #fef2f2; color: #dc2626; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
    .delete-warning { font-size: 0.875rem; color: #64748b; }
    .delete-actions { display: flex; gap: 0.75rem; justify-content: center; margin-top: 1.5rem; }
  `],
})
export class SectionsComponent implements OnInit {
  private sectionService = inject(SectionService);
  private classService = inject(ClassService);
  private academicYearService = inject(AcademicYearService);
  private toastService = inject(ToastService);

  sections = signal<Section[]>([]);
  classes = signal<Class[]>([]);
  academicYears = signal<AcademicYear[]>([]);

  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');
  classFilter = signal<string>('');

  showSectionModal = signal(false);
  showDeleteModal = signal(false);
  editingSection = signal<Section | null>(null);
  sectionToDelete = signal<Section | null>(null);

  formData = { name: '', code: '', classId: '', capacity: 40, academicYearId: '', roomNumber: '' };

  filteredSections = computed(() => {
    let result = this.sections();
    const term = this.searchTerm().toLowerCase();
    if (term) result = result.filter(s => s.name.toLowerCase().includes(term) || s.code.toLowerCase().includes(term));
    if (this.statusFilter() === 'active') result = result.filter(s => s.isActive);
    else if (this.statusFilter() === 'inactive') result = result.filter(s => !s.isActive);
    if (this.classFilter()) result = result.filter(s => s.classId === this.classFilter());
    return result;
  });

  ngOnInit(): void {
    this.loadSections();
    this.loadClasses();
    this.loadAcademicYears();
  }

  loadSections(): void {
    this.loading.set(true);
    this.sectionService.getSections().subscribe({
      next: sections => { this.sections.set(sections); this.loading.set(false); },
      error: () => { this.error.set('Failed to load sections'); this.loading.set(false); },
    });
  }

  loadClasses(): void {
    this.classService.getClasses({ isActive: true }).subscribe({
      next: classes => this.classes.set(classes),
      error: () => console.error('Failed to load classes'),
    });
  }

  loadAcademicYears(): void {
    this.academicYearService.getAcademicYears().subscribe({
      next: years => this.academicYears.set(years),
      error: () => console.error('Failed to load academic years'),
    });
  }

  openCreateModal(): void {
    this.editingSection.set(null);
    this.formData = { name: '', code: '', classId: '', capacity: 40, academicYearId: '', roomNumber: '' };
    this.showSectionModal.set(true);
  }

  editSection(section: Section): void {
    this.editingSection.set(section);
    this.formData = {
      name: section.name,
      code: section.code,
      classId: section.classId,
      capacity: section.capacity,
      academicYearId: section.academicYearId || '',
      roomNumber: section.roomNumber || '',
    };
    this.showSectionModal.set(true);
  }

  closeSectionModal(): void {
    this.showSectionModal.set(false);
    this.editingSection.set(null);
  }

  saveSection(): void {
    if (!this.formData.name || !this.formData.code || !this.formData.classId) {
      this.toastService.error('Please fill in all required fields');
      return;
    }
    this.saving.set(true);
    const editing = this.editingSection();
    const data: CreateSectionRequest | UpdateSectionRequest = {
      name: this.formData.name,
      code: this.formData.code,
      capacity: this.formData.capacity,
      roomNumber: this.formData.roomNumber || undefined,
      academicYearId: this.formData.academicYearId || undefined,
    };
    if (!editing) (data as CreateSectionRequest).classId = this.formData.classId;

    const operation = editing
      ? this.sectionService.updateSection(editing.id, data as UpdateSectionRequest)
      : this.sectionService.createSection(data as CreateSectionRequest);

    operation.subscribe({
      next: () => {
        this.toastService.success(editing ? 'Section updated' : 'Section created');
        this.closeSectionModal();
        this.loadSections();
        this.saving.set(false);
      },
      error: () => {
        this.toastService.error(editing ? 'Failed to update section' : 'Failed to create section');
        this.saving.set(false);
      },
    });
  }

  toggleStatus(section: Section): void {
    this.sectionService.updateSection(section.id, { isActive: !section.isActive }).subscribe({
      next: () => { this.toastService.success('Status updated'); this.loadSections(); },
      error: () => this.toastService.error('Failed to update status'),
    });
  }

  confirmDelete(section: Section): void {
    this.sectionToDelete.set(section);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.sectionToDelete.set(null);
  }

  deleteSection(): void {
    const section = this.sectionToDelete();
    if (!section) return;
    this.deleting.set(true);
    this.sectionService.deleteSection(section.id).subscribe({
      next: () => {
        this.toastService.success('Section deleted');
        this.closeDeleteModal();
        this.loadSections();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete section');
        this.deleting.set(false);
      },
    });
  }

  getCapacityPercent(section: Section): number {
    return Math.min(100, (section.studentCount / section.capacity) * 100);
  }
}
