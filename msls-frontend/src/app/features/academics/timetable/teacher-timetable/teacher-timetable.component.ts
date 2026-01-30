import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { TimetableService } from '../timetable.service';
import { TimetableEntry, PeriodSlot, DAYS_OF_WEEK, PERIOD_SLOT_TYPES } from '../timetable.model';
import { AcademicYearService } from '../../../admin/academic-years/academic-year.service';
import { AcademicYear } from '../../../admin/academic-years/academic-year.model';

interface SubjectStats {
  subjectName: string;
  subjectCode: string;
  periodCount: number;
  color: string;
}

interface DayStats {
  dayOfWeek: number;
  dayName: string;
  periodCount: number;
  freePeriods: number;
}

@Component({
  selector: 'msls-teacher-timetable',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-calendar-check"></i>
          </div>
          <div class="header-text">
            <h1>My Timetable</h1>
            <p>View your teaching schedule across all classes</p>
          </div>
        </div>
        <div class="header-actions">
          <select class="filter-select" [(ngModel)]="selectedAcademicYearId" (ngModelChange)="loadSchedule()">
            @for (year of academicYears(); track year.id) {
              <option [value]="year.id">{{ year.name }}</option>
            }
          </select>
        </div>
      </div>

      <!-- Statistics Cards -->
      <div class="stats-grid">
        <div class="stat-card total">
          <div class="stat-icon"><i class="fa-solid fa-clock"></i></div>
          <div class="stat-content">
            <span class="stat-value">{{ totalPeriods() }}</span>
            <span class="stat-label">Total Periods</span>
          </div>
        </div>
        <div class="stat-card free">
          <div class="stat-icon"><i class="fa-solid fa-coffee"></i></div>
          <div class="stat-content">
            <span class="stat-value">{{ totalFreePeriods() }}</span>
            <span class="stat-label">Free Periods</span>
          </div>
        </div>
        <div class="stat-card subjects">
          <div class="stat-icon"><i class="fa-solid fa-book"></i></div>
          <div class="stat-content">
            <span class="stat-value">{{ uniqueSubjects().size }}</span>
            <span class="stat-label">Subjects</span>
          </div>
        </div>
        <div class="stat-card classes">
          <div class="stat-icon"><i class="fa-solid fa-users"></i></div>
          <div class="stat-content">
            <span class="stat-value">{{ uniqueSections().size }}</span>
            <span class="stat-label">Sections</span>
          </div>
        </div>
      </div>

      <!-- View Tabs -->
      <div class="view-tabs">
        <button
          class="tab-btn"
          [class.active]="viewMode() === 'week'"
          (click)="viewMode.set('week')"
        >
          <i class="fa-solid fa-calendar-week"></i> Weekly View
        </button>
        <button
          class="tab-btn"
          [class.active]="viewMode() === 'day'"
          (click)="viewMode.set('day')"
        >
          <i class="fa-solid fa-calendar-day"></i> Day View
        </button>
      </div>

      <!-- Loading/Error States -->
      @if (loading()) {
        <div class="loading-container">
          <div class="spinner"></div>
          <span>Loading your schedule...</span>
        </div>
      } @else if (error()) {
        <div class="error-container">
          <i class="fa-solid fa-circle-exclamation"></i>
          <span>{{ error() }}</span>
          <button class="btn btn-secondary" (click)="loadSchedule()">Retry</button>
        </div>
      } @else if (entries().length === 0) {
        <div class="empty-state">
          <i class="fa-regular fa-calendar-xmark"></i>
          <h3>No Schedule Found</h3>
          <p>You don't have any classes assigned for this academic year.</p>
        </div>
      } @else {
        <!-- Weekly View -->
        @if (viewMode() === 'week') {
          <div class="grid-container">
            <div class="timetable-grid">
              <!-- Header Row -->
              <div class="grid-header">
                <div class="grid-cell header-cell time-col">Time</div>
                @for (day of workingDays; track day.value) {
                  <div class="grid-cell header-cell" [class.today]="isToday(day.value)">
                    {{ day.label }}
                    @if (isToday(day.value)) {
                      <span class="today-badge">Today</span>
                    }
                  </div>
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
                  @for (day of workingDays; track day.value) {
                    <div
                      class="grid-cell entry-cell"
                      [class.break-cell]="!isTeachingSlot(slot)"
                      [class.current-period]="isCurrentPeriod(day.value, slot)"
                      [class.today-cell]="isToday(day.value)"
                    >
                      @if (getEntry(day.value, slot.id); as entry) {
                        <div class="entry-content" [style.background-color]="getSubjectColor(entry.subjectId)">
                          <span class="entry-subject">{{ entry.subjectName || 'Free Period' }}</span>
                          <span class="entry-section">{{ entry.timetable?.className }} - {{ entry.timetable?.sectionName }}</span>
                          @if (entry.roomNumber) {
                            <span class="entry-room">{{ entry.roomNumber }}</span>
                          }
                        </div>
                      } @else if (isTeachingSlot(slot)) {
                        <div class="free-slot">
                          <i class="fa-regular fa-face-smile"></i>
                          <span>Free</span>
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
        }

        <!-- Day View -->
        @if (viewMode() === 'day') {
          <div class="day-view">
            <div class="day-selector">
              @for (day of workingDays; track day.value) {
                <button
                  class="day-btn"
                  [class.active]="selectedDay() === day.value"
                  [class.today]="isToday(day.value)"
                  (click)="selectedDay.set(day.value)"
                >
                  <span class="day-short">{{ day.short }}</span>
                  <span class="day-label">{{ day.label }}</span>
                </button>
              }
            </div>

            <div class="day-schedule">
              @for (slot of periodSlots(); track slot.id) {
                <div
                  class="schedule-item"
                  [class.break-item]="!isTeachingSlot(slot)"
                  [class.current]="isCurrentPeriod(selectedDay(), slot)"
                >
                  <div class="schedule-time">
                    <span class="time-start">{{ slot.startTime }}</span>
                    <span class="time-end">{{ slot.endTime }}</span>
                  </div>
                  <div class="schedule-content">
                    @if (getEntry(selectedDay(), slot.id); as entry) {
                      <div class="schedule-entry" [style.border-left-color]="getSubjectColor(entry.subjectId)">
                        <h4>{{ entry.subjectName || 'Free Period' }}</h4>
                        <p class="entry-details">
                          <i class="fa-solid fa-users"></i>
                          {{ entry.timetable?.className }} - {{ entry.timetable?.sectionName }}
                        </p>
                        @if (entry.roomNumber) {
                          <p class="entry-room">
                            <i class="fa-solid fa-location-dot"></i>
                            {{ entry.roomNumber }}
                          </p>
                        }
                      </div>
                    } @else if (isTeachingSlot(slot)) {
                      <div class="schedule-free">
                        <i class="fa-regular fa-face-smile"></i>
                        <span>Free Period</span>
                      </div>
                    } @else {
                      <div class="schedule-break">
                        <i [class]="getSlotTypeIcon(slot.slotType)"></i>
                        <span>{{ getSlotTypeLabel(slot.slotType) }}</span>
                      </div>
                    }
                  </div>
                </div>
              }
            </div>
          </div>
        }

        <!-- Subject Breakdown -->
        <div class="subject-breakdown">
          <h3>Subject Distribution</h3>
          <div class="subject-list">
            @for (stat of subjectStats(); track stat.subjectName) {
              <div class="subject-item">
                <div class="subject-color" [style.background-color]="stat.color"></div>
                <span class="subject-name">{{ stat.subjectName }}</span>
                <span class="subject-code">{{ stat.subjectCode }}</span>
                <span class="period-count">{{ stat.periodCount }} periods/week</span>
              </div>
            }
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
      align-items: center;
      margin-bottom: 1.5rem;
      gap: 1rem;
    }

    .header-content { display: flex; align-items: center; gap: 1rem; }
    .header-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; background: linear-gradient(135deg, #4f46e5, #7c3aed); color: white; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
    .header-text h1 { margin: 0; font-size: 1.5rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }

    .filter-select {
      padding: 0.625rem 2rem 0.625rem 1rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
      cursor: pointer;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .stat-card {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1.25rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
    }

    .stat-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .stat-card.total .stat-icon { background: #dbeafe; color: #2563eb; }
    .stat-card.free .stat-icon { background: #dcfce7; color: #16a34a; }
    .stat-card.subjects .stat-icon { background: #fef3c7; color: #d97706; }
    .stat-card.classes .stat-icon { background: #e0e7ff; color: #4f46e5; }

    .stat-content { display: flex; flex-direction: column; }
    .stat-value { font-size: 1.5rem; font-weight: 700; color: #1e293b; }
    .stat-label { font-size: 0.75rem; color: #64748b; }

    .view-tabs {
      display: flex;
      gap: 0.5rem;
      margin-bottom: 1rem;
    }

    .tab-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.2s;
    }

    .tab-btn:hover { border-color: #c7d2fe; color: #4f46e5; }
    .tab-btn.active { background: #4f46e5; border-color: #4f46e5; color: white; }

    .loading-container, .error-container, .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 4rem;
      background: white;
      border-radius: 1rem;
      border: 1px solid #e2e8f0;
      color: #64748b;
    }

    .error-container { color: #dc2626; }
    .empty-state i { font-size: 3rem; color: #cbd5e1; }
    .empty-state h3 { margin: 0; color: #1e293b; }
    .empty-state p { margin: 0; }

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
      margin-bottom: 1.5rem;
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
      flex: 0 0 110px;
      min-width: 110px;
      background: #f8fafc;
    }

    .header-cell {
      background: #4f46e5;
      color: white;
      font-weight: 600;
      font-size: 0.875rem;
      text-align: center;
      position: relative;
    }

    .header-cell.today {
      background: #16a34a;
    }

    .header-cell.time-col {
      background: #4338ca;
    }

    .today-badge {
      display: block;
      font-size: 0.625rem;
      font-weight: 400;
      text-transform: uppercase;
      opacity: 0.9;
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

    .slot-name { font-weight: 500; font-size: 0.8rem; color: #1e293b; }
    .slot-time { font-size: 0.7rem; color: #64748b; }

    .entry-cell {
      min-height: 65px;
      transition: all 0.2s;
    }

    .entry-cell.today-cell {
      background: #f0fdf4;
    }

    .entry-cell.current-period {
      background: #fef3c7;
      box-shadow: inset 0 0 0 2px #f59e0b;
    }

    .break-row .grid-cell {
      background: #f8fafc;
    }

    .break-label {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      font-size: 0.75rem;
      color: #64748b;
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
      font-size: 0.75rem;
      color: #1e293b;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .entry-section {
      font-size: 0.65rem;
      color: #475569;
    }

    .entry-room {
      font-size: 0.6rem;
      color: #64748b;
    }

    .free-slot {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #94a3b8;
      font-size: 0.75rem;
      gap: 0.25rem;
    }

    /* Day View Styles */
    .day-view {
      background: white;
      border-radius: 1rem;
      border: 1px solid #e2e8f0;
      overflow: hidden;
      margin-bottom: 1.5rem;
    }

    .day-selector {
      display: flex;
      border-bottom: 1px solid #e2e8f0;
      overflow-x: auto;
    }

    .day-btn {
      flex: 1;
      min-width: 100px;
      padding: 1rem;
      background: transparent;
      border: none;
      cursor: pointer;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.25rem;
      transition: all 0.2s;
      border-bottom: 3px solid transparent;
    }

    .day-btn:hover { background: #f8fafc; }
    .day-btn.active { border-bottom-color: #4f46e5; background: #f8fafc; }
    .day-btn.today .day-short { background: #16a34a; color: white; }

    .day-short {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 50%;
      background: #f1f5f9;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      font-size: 0.875rem;
      color: #475569;
    }

    .day-label {
      font-size: 0.75rem;
      color: #64748b;
    }

    .day-schedule {
      padding: 1rem;
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }

    .schedule-item {
      display: flex;
      gap: 1rem;
      padding: 0.75rem;
      border-radius: 0.75rem;
      background: #f8fafc;
      transition: all 0.2s;
    }

    .schedule-item.current {
      background: #fef3c7;
      box-shadow: 0 0 0 2px #f59e0b;
    }

    .schedule-item.break-item {
      background: #f1f5f9;
    }

    .schedule-time {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      min-width: 60px;
      padding: 0.5rem;
      background: white;
      border-radius: 0.5rem;
      font-size: 0.75rem;
      color: #64748b;
    }

    .time-start { font-weight: 600; color: #1e293b; }

    .schedule-content {
      flex: 1;
      display: flex;
      align-items: center;
    }

    .schedule-entry {
      padding-left: 0.75rem;
      border-left: 3px solid #4f46e5;
    }

    .schedule-entry h4 { margin: 0 0 0.25rem; font-size: 0.9375rem; color: #1e293b; }
    .entry-details, .entry-room { margin: 0; font-size: 0.8125rem; color: #64748b; display: flex; align-items: center; gap: 0.375rem; }

    .schedule-free, .schedule-break {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      color: #64748b;
      font-size: 0.875rem;
    }

    /* Subject Breakdown */
    .subject-breakdown {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.25rem;
    }

    .subject-breakdown h3 {
      margin: 0 0 1rem;
      font-size: 1rem;
      font-weight: 600;
      color: #1e293b;
    }

    .subject-list {
      display: flex;
      flex-wrap: wrap;
      gap: 0.75rem;
    }

    .subject-item {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.75rem 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      font-size: 0.875rem;
    }

    .subject-color {
      width: 0.75rem;
      height: 0.75rem;
      border-radius: 0.25rem;
    }

    .subject-name { font-weight: 500; color: #1e293b; }
    .subject-code { color: #64748b; font-size: 0.75rem; }
    .period-count { margin-left: auto; color: #4f46e5; font-weight: 500; font-size: 0.75rem; }

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

    .btn-secondary { background: #f1f5f9; color: #475569; }
    .btn-secondary:hover { background: #e2e8f0; }

    @media (max-width: 768px) {
      .page { padding: 1rem; }
      .page-header { flex-direction: column; align-items: flex-start; }
      .stats-grid { grid-template-columns: repeat(2, 1fr); }
      .view-tabs { width: 100%; }
      .tab-btn { flex: 1; justify-content: center; }
    }
  `]
})
export class TeacherTimetableComponent implements OnInit {
  private readonly timetableService = inject(TimetableService);
  private readonly academicYearService = inject(AcademicYearService);

  // Data
  entries = signal<TimetableEntry[]>([]);
  periodSlots = signal<PeriodSlot[]>([]);
  academicYears = signal<AcademicYear[]>([]);

  // State
  loading = signal(true);
  error = signal<string | null>(null);
  viewMode = signal<'week' | 'day'>('week');
  selectedDay = signal(new Date().getDay() || 1); // Default to today or Monday
  selectedAcademicYearId = '';

  // Entry map for quick lookup
  private entryMap = new Map<string, TimetableEntry>();

  // Subject colors
  private subjectColors = new Map<string, string>();
  private colorPalette = [
    '#eef2ff', '#f3e8ff', '#fce7f3', '#fee2e2', '#ffedd5',
    '#fef3c7', '#ecfccb', '#dcfce7', '#ccfbf1', '#e0f2fe',
    '#e0e7ff', '#ede9fe', '#fae8ff', '#fecdd3', '#fed7aa'
  ];

  workingDays = DAYS_OF_WEEK.filter(d => d.value >= 1 && d.value <= 6);

  // Computed stats
  totalPeriods = computed(() => {
    return this.entries().filter(e => e.subjectId && !e.isFreePeriod).length;
  });

  totalFreePeriods = computed(() => {
    const teachingSlots = this.periodSlots().filter(s => this.isTeachingSlot(s));
    const totalSlots = teachingSlots.length * this.workingDays.length;
    return totalSlots - this.totalPeriods();
  });

  uniqueSubjects = computed(() => {
    const subjects = new Set<string>();
    this.entries().forEach(e => {
      if (e.subjectId) subjects.add(e.subjectId);
    });
    return subjects;
  });

  uniqueSections = computed(() => {
    const sections = new Set<string>();
    this.entries().forEach(e => {
      if (e.timetable?.sectionId) sections.add(e.timetable.sectionId);
    });
    return sections;
  });

  subjectStats = computed(() => {
    const stats = new Map<string, SubjectStats>();

    this.entries().forEach(e => {
      if (e.subjectId && e.subjectName) {
        const existing = stats.get(e.subjectId);
        if (existing) {
          existing.periodCount++;
        } else {
          stats.set(e.subjectId, {
            subjectName: e.subjectName,
            subjectCode: e.subjectCode || '',
            periodCount: 1,
            color: this.getSubjectColor(e.subjectId)
          });
        }
      }
    });

    return Array.from(stats.values()).sort((a, b) => b.periodCount - a.periodCount);
  });

  ngOnInit(): void {
    this.loadAcademicYears();
  }

  loadAcademicYears(): void {
    this.academicYearService.getAcademicYears().subscribe({
      next: years => {
        this.academicYears.set(years);
        // Select current academic year
        const current = years.find(y => y.isCurrent);
        if (current) {
          this.selectedAcademicYearId = current.id;
          this.loadSchedule();
        } else if (years.length > 0) {
          this.selectedAcademicYearId = years[0].id;
          this.loadSchedule();
        } else {
          this.loading.set(false);
        }
      },
      error: () => {
        this.error.set('Failed to load academic years');
        this.loading.set(false);
      }
    });
  }

  loadSchedule(): void {
    if (!this.selectedAcademicYearId) return;

    this.loading.set(true);
    this.error.set(null);

    this.timetableService.getMySchedule(this.selectedAcademicYearId).subscribe({
      next: response => {
        this.entries.set(response.entries || []);
        this.buildEntryMap(response.entries || []);
        this.assignSubjectColors(response.entries || []);
        this.loadPeriodSlots();
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load your schedule');
        this.loading.set(false);
      }
    });
  }

  loadPeriodSlots(): void {
    // Get branch ID from first entry if available
    const firstEntry = this.entries()[0];
    if (firstEntry?.timetable?.branchId) {
      this.timetableService.getPeriodSlots({ branchId: firstEntry.timetable.branchId, isActive: true }).subscribe({
        next: slots => this.periodSlots.set(slots),
        error: () => console.error('Failed to load period slots')
      });
    }
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

  getSlotTypeIcon(slotType: string): string {
    const type = PERIOD_SLOT_TYPES.find(t => t.value === slotType);
    return type ? `fa-solid ${type.icon}` : 'fa-solid fa-clock';
  }

  isToday(dayOfWeek: number): boolean {
    return new Date().getDay() === dayOfWeek;
  }

  isCurrentPeriod(dayOfWeek: number, slot: PeriodSlot): boolean {
    if (!this.isToday(dayOfWeek)) return false;

    const now = new Date();
    const currentTime = now.getHours() * 60 + now.getMinutes();

    const [startHour, startMin] = slot.startTime.split(':').map(Number);
    const [endHour, endMin] = slot.endTime.split(':').map(Number);

    const slotStart = startHour * 60 + startMin;
    const slotEnd = endHour * 60 + endMin;

    return currentTime >= slotStart && currentTime < slotEnd;
  }
}
