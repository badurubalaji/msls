import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { TimetableService } from '../timetable.service';
import { Timetable, TimetableEntry, PeriodSlot, DAYS_OF_WEEK, CreateTimetableEntryRequest, PERIOD_SLOT_TYPES } from '../timetable.model';
import { ToastService } from '../../../../shared/services/toast.service';
import { SubjectService } from '../../services/subject.service';
import { StaffService } from '../../../staff/services/staff.service';
import { Subject } from '../../academic.model';
import { Staff } from '../../../staff/models/staff.model';

@Component({
  selector: 'msls-timetable-builder',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <!-- Header -->
      <div class="page-header">
        <div class="header-left">
          <button class="back-btn" (click)="goBack()">
            <i class="fa-solid fa-arrow-left"></i>
          </button>
          <div class="header-content">
            <h1>{{ timetable()?.name || 'Timetable Builder' }}</h1>
            <p>{{ timetable()?.className }} - {{ timetable()?.sectionName }} | {{ timetable()?.academicYearName }}</p>
          </div>
        </div>
        <div class="header-actions">
          @if (timetable()?.status === 'draft') {
            <button class="btn btn-success" (click)="publishTimetable()">
              <i class="fa-solid fa-check-circle"></i> Publish
            </button>
          }
          <span class="status-badge" [class]="timetable()?.status">
            {{ timetable()?.status | titlecase }}
          </span>
        </div>
      </div>

      <!-- Loading/Error States -->
      @if (loading()) {
        <div class="loading-container">
          <div class="spinner"></div>
          <span>Loading timetable...</span>
        </div>
      } @else if (error()) {
        <div class="error-container">
          <i class="fa-solid fa-circle-exclamation"></i>
          <span>{{ error() }}</span>
          <button class="btn btn-secondary" (click)="loadTimetable()">Retry</button>
        </div>
      } @else {
        <!-- Timetable Grid -->
        <div class="grid-container">
          <div class="timetable-grid">
            <!-- Header Row -->
            <div class="grid-header">
              <div class="grid-cell header-cell time-col">Time</div>
              @for (day of workingDays(); track day.value) {
                <div class="grid-cell header-cell">{{ day.label }}</div>
              }
            </div>

            <!-- Period Rows -->
            @for (slot of periodSlots(); track slot.id) {
              <div class="grid-row" [class.break-row]="!isTeachingSlot(slot)">
                <div class="grid-cell time-cell time-col">
                  <div class="slot-info">
                    <span class="slot-name">{{ slot.name }}</span>
                    <span class="slot-time">{{ slot.startTime }} - {{ slot.endTime }}</span>
                  </div>
                </div>
                @for (day of workingDays(); track day.value) {
                  <div
                    class="grid-cell entry-cell"
                    [class.break-cell]="!isTeachingSlot(slot)"
                    [class.has-entry]="getEntry(day.value, slot.id)"
                    [class.editable]="timetable()?.status === 'draft' && isTeachingSlot(slot)"
                    (click)="openEntryModal(day.value, slot)"
                  >
                    @if (getEntry(day.value, slot.id); as entry) {
                      <div class="entry-content" [style.background-color]="getSubjectColor(entry.subjectId)">
                        <span class="entry-subject">{{ entry.subjectName || entry.subjectCode || 'Free' }}</span>
                        @if (entry.staffName) {
                          <span class="entry-teacher">{{ entry.staffName }}</span>
                        }
                        @if (entry.roomNumber) {
                          <span class="entry-room">{{ entry.roomNumber }}</span>
                        }
                      </div>
                    } @else if (isTeachingSlot(slot)) {
                      <div class="empty-slot">
                        <i class="fa-solid fa-plus"></i>
                      </div>
                    } @else {
                      <div class="break-label">{{ getSlotTypeLabel(slot.slotType) }}</div>
                    }
                  </div>
                }
              </div>
            }
          </div>
        </div>

        <!-- Legend -->
        <div class="legend">
          <span class="legend-title">Legend:</span>
          @for (type of slotTypes; track type.value) {
            <span class="legend-item">
              <span class="legend-color" [class]="'slot-' + type.value"></span>
              {{ type.label }}
            </span>
          }
        </div>
      }

      <!-- Entry Modal -->
      @if (showEntryModal()) {
        <div class="modal-overlay" (click)="closeEntryModal()">
          <div class="modal" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-clock"></i>
                {{ selectedSlot()?.name }} - {{ selectedDayName() }}
              </h3>
              <button class="modal__close" (click)="closeEntryModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form class="form">
                <div class="form-group">
                  <label>Subject</label>
                  <select [(ngModel)]="entryForm.subjectId" name="subjectId" class="form-input">
                    <option value="">-- Free Period --</option>
                    @for (subject of subjects(); track subject.id) {
                      <option [value]="subject.id">{{ subject.name }} ({{ subject.code }})</option>
                    }
                  </select>
                </div>
                <div class="form-group">
                  <label>Teacher</label>
                  <select [(ngModel)]="entryForm.staffId" name="staffId" class="form-input">
                    <option value="">-- Select Teacher --</option>
                    @for (staff of teachers(); track staff.id) {
                      <option [value]="staff.id">{{ staff.firstName }} {{ staff.lastName }}</option>
                    }
                  </select>
                </div>
                <div class="form-group">
                  <label>Room Number</label>
                  <input type="text" [(ngModel)]="entryForm.roomNumber" name="roomNumber" class="form-input" placeholder="e.g., Room 101" />
                </div>
                <div class="form-group">
                  <label>Notes</label>
                  <textarea [(ngModel)]="entryForm.notes" name="notes" class="form-input" rows="2" placeholder="Optional notes..."></textarea>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              @if (selectedEntry()) {
                <button class="btn btn-danger" (click)="deleteEntry()">
                  <i class="fa-solid fa-trash"></i> Remove
                </button>
              }
              <div class="spacer"></div>
              <button class="btn btn-secondary" (click)="closeEntryModal()">Cancel</button>
              <button class="btn btn-primary" [disabled]="savingEntry()" (click)="saveEntry()">
                @if (savingEntry()) {
                  <div class="btn-spinner"></div> Saving...
                } @else {
                  <i class="fa-solid fa-check"></i> Save
                }
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1600px; margin: 0 auto; }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1.5rem;
      gap: 1rem;
    }

    .header-left { display: flex; align-items: center; gap: 1rem; }

    .back-btn {
      width: 2.5rem;
      height: 2.5rem;
      border: 1px solid #e2e8f0;
      background: white;
      border-radius: 0.5rem;
      cursor: pointer;
      color: #64748b;
      transition: all 0.2s;
    }

    .back-btn:hover { background: #f8fafc; color: #4f46e5; }

    .header-content h1 { margin: 0; font-size: 1.25rem; font-weight: 600; color: #1e293b; }
    .header-content p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }

    .header-actions { display: flex; align-items: center; gap: 1rem; }

    .status-badge {
      display: inline-flex;
      padding: 0.375rem 1rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
    }

    .status-badge.draft { background: #fef3c7; color: #92400e; }
    .status-badge.published { background: #dcfce7; color: #166534; }
    .status-badge.archived { background: #f1f5f9; color: #64748b; }

    .loading-container, .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 4rem;
      background: white;
      border-radius: 1rem;
      border: 1px solid #e2e8f0;
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

    .grid-container {
      background: white;
      border-radius: 1rem;
      border: 1px solid #e2e8f0;
      overflow: hidden;
    }

    .timetable-grid {
      display: flex;
      flex-direction: column;
      overflow-x: auto;
    }

    .grid-header, .grid-row {
      display: flex;
      min-width: max-content;
    }

    .grid-cell {
      flex: 1;
      min-width: 140px;
      padding: 0.75rem;
      border-right: 1px solid #e2e8f0;
      border-bottom: 1px solid #e2e8f0;
    }

    .grid-cell:last-child { border-right: none; }

    .time-col {
      flex: 0 0 120px;
      min-width: 120px;
      background: #f8fafc;
    }

    .header-cell {
      background: #4f46e5;
      color: white;
      font-weight: 600;
      font-size: 0.875rem;
      text-align: center;
    }

    .header-cell.time-col {
      background: #4338ca;
    }

    .time-cell {
      display: flex;
      align-items: center;
    }

    .slot-info {
      display: flex;
      flex-direction: column;
      gap: 0.125rem;
    }

    .slot-name { font-weight: 500; font-size: 0.875rem; color: #1e293b; }
    .slot-time { font-size: 0.75rem; color: #64748b; }

    .entry-cell {
      min-height: 70px;
      cursor: default;
      transition: all 0.2s;
    }

    .entry-cell.editable {
      cursor: pointer;
    }

    .entry-cell.editable:hover {
      background: #f8fafc;
    }

    .entry-cell.has-entry {
      padding: 0.25rem;
    }

    .break-row .grid-cell {
      background: #fef3c7;
    }

    .break-cell {
      cursor: default !important;
    }

    .break-label {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      font-size: 0.75rem;
      color: #92400e;
      font-weight: 500;
    }

    .entry-content {
      height: 100%;
      padding: 0.5rem;
      border-radius: 0.5rem;
      background: #eef2ff;
      display: flex;
      flex-direction: column;
      gap: 0.125rem;
    }

    .entry-subject {
      font-weight: 600;
      font-size: 0.8rem;
      color: #1e293b;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .entry-teacher {
      font-size: 0.7rem;
      color: #475569;
    }

    .entry-room {
      font-size: 0.65rem;
      color: #64748b;
    }

    .empty-slot {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #cbd5e1;
      font-size: 1rem;
    }

    .legend {
      display: flex;
      align-items: center;
      gap: 1.5rem;
      padding: 1rem;
      margin-top: 1rem;
      background: white;
      border-radius: 0.5rem;
      border: 1px solid #e2e8f0;
      flex-wrap: wrap;
    }

    .legend-title { font-weight: 600; color: #374151; font-size: 0.875rem; }

    .legend-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.75rem;
      color: #64748b;
    }

    .legend-color {
      width: 1rem;
      height: 1rem;
      border-radius: 0.25rem;
    }

    .slot-regular { background: #eef2ff; }
    .slot-short { background: #f3e8ff; }
    .slot-assembly { background: #fef3c7; }
    .slot-break { background: #dcfce7; }
    .slot-lunch { background: #ffedd5; }
    .slot-activity { background: #ccfbf1; }
    .slot-zero_period { background: #f1f5f9; }

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

    .btn-primary { background: #4f46e5; color: white; }
    .btn-primary:hover:not(:disabled) { background: #4338ca; }
    .btn-secondary { background: #f1f5f9; color: #475569; }
    .btn-secondary:hover { background: #e2e8f0; }
    .btn-success { background: #16a34a; color: white; }
    .btn-success:hover:not(:disabled) { background: #15803d; }
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
      max-width: 28rem;
      max-height: 90vh;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
    }

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
        font-size: 1rem;
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
      align-items: center;
      gap: 0.75rem;
      padding: 1.25rem 1.5rem;
      background: #f8fafc;
      border-top: 1px solid #f1f5f9;
    }

    .spacer { flex: 1; }

    .form { display: flex; flex-direction: column; gap: 1rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-group label { font-size: 0.875rem; font-weight: 500; color: #374151; }

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

    @media (max-width: 768px) {
      .page { padding: 1rem; }
      .page-header { flex-direction: column; align-items: flex-start; }
      .header-actions { width: 100%; justify-content: flex-end; }
    }
  `]
})
export class TimetableBuilderComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly timetableService = inject(TimetableService);
  private readonly subjectService = inject(SubjectService);
  private readonly staffService = inject(StaffService);
  private readonly toastService = inject(ToastService);

  slotTypes = PERIOD_SLOT_TYPES;

  // Data
  timetable = signal<Timetable | null>(null);
  periodSlots = signal<PeriodSlot[]>([]);
  entries = signal<TimetableEntry[]>([]);
  subjects = signal<Subject[]>([]);
  teachers = signal<Staff[]>([]);

  // State
  loading = signal(true);
  error = signal<string | null>(null);
  savingEntry = signal(false);

  // Modal
  showEntryModal = signal(false);
  selectedDay = signal<number>(0);
  selectedSlot = signal<PeriodSlot | null>(null);
  selectedEntry = signal<TimetableEntry | null>(null);

  entryForm: CreateTimetableEntryRequest = {
    dayOfWeek: 0,
    periodSlotId: '',
    subjectId: undefined,
    staffId: undefined,
    roomNumber: '',
    notes: '',
    isFreePeriod: false
  };

  // Entry map for quick lookup
  private entryMap = new Map<string, TimetableEntry>();

  // Subject colors
  private subjectColors = new Map<string, string>();
  private colorPalette = [
    '#eef2ff', '#f3e8ff', '#fce7f3', '#fee2e2', '#ffedd5',
    '#fef3c7', '#ecfccb', '#dcfce7', '#ccfbf1', '#e0f2fe',
    '#e0e7ff', '#ede9fe', '#fae8ff', '#fecdd3', '#fed7aa'
  ];

  workingDays = computed(() => {
    // Default working days (Monday-Saturday)
    return DAYS_OF_WEEK.filter(d => d.value >= 1 && d.value <= 6);
  });

  selectedDayName = computed(() => {
    const day = this.selectedDay();
    return DAYS_OF_WEEK.find(d => d.value === day)?.label || '';
  });

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadTimetable(id);
      this.loadSubjects();
      this.loadTeachers();
    } else {
      this.error.set('Timetable ID not provided');
      this.loading.set(false);
    }
  }

  loadTimetable(id?: string): void {
    const timetableId = id || this.timetable()?.id;
    if (!timetableId) return;

    this.loading.set(true);
    this.error.set(null);

    this.timetableService.getTimetable(timetableId).subscribe({
      next: tt => {
        this.timetable.set(tt);
        this.entries.set(tt.entries || []);
        this.buildEntryMap(tt.entries || []);
        this.assignSubjectColors(tt.entries || []);
        this.loadPeriodSlots(tt.branchId);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load timetable');
        this.loading.set(false);
      }
    });
  }

  loadPeriodSlots(branchId: string): void {
    this.timetableService.getPeriodSlots({ branchId, isActive: true }).subscribe({
      next: slots => this.periodSlots.set(slots),
      error: () => console.error('Failed to load period slots')
    });
  }

  loadSubjects(): void {
    this.subjectService.getSubjects({ isActive: true }).subscribe({
      next: subjects => this.subjects.set(subjects),
      error: () => console.error('Failed to load subjects')
    });
  }

  loadTeachers(): void {
    this.staffService.loadStaff({ staffType: 'teaching', status: 'active' }).subscribe({
      next: response => this.teachers.set(response.staff || []),
      error: () => console.error('Failed to load staff')
    });
  }

  private buildEntryMap(entries: TimetableEntry[]): void {
    this.entryMap.clear();
    entries.forEach(e => {
      this.entryMap.set(`${e.dayOfWeek}-${e.periodSlotId}`, e);
    });
  }

  private assignSubjectColors(entries: TimetableEntry[]): void {
    let colorIndex = 0;
    entries.forEach(e => {
      if (e.subjectId && !this.subjectColors.has(e.subjectId)) {
        this.subjectColors.set(e.subjectId, this.colorPalette[colorIndex % this.colorPalette.length]);
        colorIndex++;
      }
    });
  }

  getEntry(dayOfWeek: number, periodSlotId: string): TimetableEntry | undefined {
    return this.entryMap.get(`${dayOfWeek}-${periodSlotId}`);
  }

  getSubjectColor(subjectId?: string): string {
    if (!subjectId) return '#f1f5f9';
    return this.subjectColors.get(subjectId) || '#eef2ff';
  }

  isTeachingSlot(slot: PeriodSlot): boolean {
    return slot.slotType === 'regular' || slot.slotType === 'short';
  }

  getSlotTypeLabel(slotType: string): string {
    return PERIOD_SLOT_TYPES.find(t => t.value === slotType)?.label || slotType;
  }

  openEntryModal(dayOfWeek: number, slot: PeriodSlot): void {
    if (this.timetable()?.status !== 'draft' || !this.isTeachingSlot(slot)) return;

    this.selectedDay.set(dayOfWeek);
    this.selectedSlot.set(slot);

    const existingEntry = this.getEntry(dayOfWeek, slot.id);
    this.selectedEntry.set(existingEntry || null);

    this.entryForm = {
      dayOfWeek,
      periodSlotId: slot.id,
      subjectId: existingEntry?.subjectId || undefined,
      staffId: existingEntry?.staffId || undefined,
      roomNumber: existingEntry?.roomNumber || '',
      notes: existingEntry?.notes || '',
      isFreePeriod: existingEntry?.isFreePeriod || false
    };

    this.showEntryModal.set(true);
  }

  closeEntryModal(): void {
    this.showEntryModal.set(false);
    this.selectedEntry.set(null);
  }

  saveEntry(): void {
    const timetableId = this.timetable()?.id;
    if (!timetableId) return;

    this.savingEntry.set(true);

    const data: CreateTimetableEntryRequest = {
      dayOfWeek: this.entryForm.dayOfWeek,
      periodSlotId: this.entryForm.periodSlotId,
      subjectId: this.entryForm.subjectId || undefined,
      staffId: this.entryForm.staffId || undefined,
      roomNumber: this.entryForm.roomNumber || undefined,
      notes: this.entryForm.notes || undefined,
      isFreePeriod: !this.entryForm.subjectId
    };

    this.timetableService.upsertTimetableEntry(timetableId, data).subscribe({
      next: () => {
        this.toastService.success('Entry saved successfully');
        this.closeEntryModal();
        this.savingEntry.set(false);
        this.loadTimetable();
      },
      error: () => {
        this.toastService.error('Failed to save entry');
        this.savingEntry.set(false);
      }
    });
  }

  deleteEntry(): void {
    const timetableId = this.timetable()?.id;
    const entry = this.selectedEntry();
    if (!timetableId || !entry) return;

    this.timetableService.deleteTimetableEntry(timetableId, entry.id).subscribe({
      next: () => {
        this.toastService.success('Entry removed');
        this.closeEntryModal();
        this.loadTimetable();
      },
      error: () => this.toastService.error('Failed to remove entry')
    });
  }

  publishTimetable(): void {
    const timetableId = this.timetable()?.id;
    if (!timetableId) return;

    this.timetableService.publishTimetable(timetableId).subscribe({
      next: () => {
        this.toastService.success('Timetable published successfully');
        this.loadTimetable();
      },
      error: () => this.toastService.error('Failed to publish timetable')
    });
  }

  goBack(): void {
    this.router.navigate(['/academics/timetable/list']);
  }
}
