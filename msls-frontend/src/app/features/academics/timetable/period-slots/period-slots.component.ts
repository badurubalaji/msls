import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../../shared/components/modal/modal.component';
import { TimetableService } from '../timetable.service';
import {
  PeriodSlot,
  DayPattern,
  Shift,
  CreatePeriodSlotRequest,
  UpdatePeriodSlotRequest,
  PeriodSlotType,
  PERIOD_SLOT_TYPES,
} from '../timetable.model';
import { ToastService } from '../../../../shared/services/toast.service';
import { BranchService } from '../../../admin/branches/branch.service';
import { Branch } from '../../../admin/branches/branch.model';

@Component({
  selector: 'msls-period-slots',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-table-cells"></i>
          </div>
          <div class="header-text">
            <h1>Period Slots</h1>
            <p>Configure period timings, breaks, and activities for each day pattern</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Period
        </button>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="filter-group">
          <label class="filter-label">Branch</label>
          <select
            class="filter-select"
            [ngModel]="selectedBranchId()"
            (ngModelChange)="onBranchChange($event)"
          >
            <option value="">All Branches</option>
            @for (branch of branches(); track branch.id) {
              <option [value]="branch.id">{{ branch.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Day Pattern</label>
          <select
            class="filter-select"
            [ngModel]="selectedDayPatternId()"
            (ngModelChange)="selectedDayPatternId.set($event)"
          >
            <option value="">All Patterns</option>
            @for (pattern of dayPatterns(); track pattern.id) {
              <option [value]="pattern.id">{{ pattern.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Shift</label>
          <select
            class="filter-select"
            [ngModel]="selectedShiftId()"
            (ngModelChange)="selectedShiftId.set($event)"
          >
            <option value="">All Shifts</option>
            @for (shift of shifts(); track shift.id) {
              <option [value]="shift.id">{{ shift.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Type</label>
          <select
            class="filter-select"
            [ngModel]="selectedSlotType()"
            (ngModelChange)="selectedSlotType.set($event)"
          >
            <option value="">All Types</option>
            @for (type of slotTypes; track type.value) {
              <option [value]="type.value">{{ type.label }}</option>
            }
          </select>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading period slots...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadPeriodSlots()">Retry</button>
          </div>
        } @else {
          <!-- Visual Timeline View -->
          <div class="timeline-container">
            <div class="timeline">
              @for (slot of filteredSlots(); track slot.id) {
                <div
                  class="slot-card"
                  [class]="getSlotTypeClass(slot.slotType)"
                  (click)="editSlot(slot)"
                >
                  <div class="slot-header">
                    <span class="slot-number">
                      @if (slot.periodNumber) {
                        Period {{ slot.periodNumber }}
                      } @else {
                        {{ getSlotTypeLabel(slot.slotType) }}
                      }
                    </span>
                    <span class="slot-type-badge" [class]="getSlotTypeBadgeClass(slot.slotType)">
                      <i class="fa-solid" [class]="getSlotTypeIcon(slot.slotType)"></i>
                    </span>
                  </div>
                  <div class="slot-name">{{ slot.name }}</div>
                  <div class="slot-time">
                    {{ formatTime(slot.startTime) }} - {{ formatTime(slot.endTime) }}
                  </div>
                  <div class="slot-duration">{{ slot.durationMinutes }} min</div>
                  <div class="slot-actions">
                    <button class="slot-action-btn" (click)="editSlot(slot); $event.stopPropagation()" title="Edit">
                      <i class="fa-regular fa-pen-to-square"></i>
                    </button>
                    <button class="slot-action-btn slot-action-btn--danger" (click)="confirmDelete(slot); $event.stopPropagation()" title="Delete">
                      <i class="fa-regular fa-trash-can"></i>
                    </button>
                  </div>
                </div>
              } @empty {
                <div class="empty-timeline">
                  <i class="fa-regular fa-clock"></i>
                  <p>No period slots configured</p>
                  <button class="btn btn-primary btn-sm" (click)="openCreateModal()">Add Period</button>
                </div>
              }
            </div>
          </div>

          <!-- Table View -->
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Order</th>
                  <th>Period</th>
                  <th>Type</th>
                  <th>Time</th>
                  <th>Duration</th>
                  <th>Day Pattern</th>
                  <th>Shift</th>
                  <th>Status</th>
                  <th style="width: 100px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (slot of filteredSlots(); track slot.id) {
                  <tr>
                    <td><span class="order-badge">{{ slot.displayOrder }}</span></td>
                    <td>
                      <div class="name-wrapper">
                        <div class="slot-icon" [class]="getSlotTypeBadgeClass(slot.slotType)">
                          <i class="fa-solid" [class]="getSlotTypeIcon(slot.slotType)"></i>
                        </div>
                        <div class="name-content">
                          <span class="name">{{ slot.name }}</span>
                          @if (slot.periodNumber) {
                            <span class="description">Period {{ slot.periodNumber }}</span>
                          }
                        </div>
                      </div>
                    </td>
                    <td>
                      <span class="type-badge" [class]="getSlotTypeBadgeClass(slot.slotType)">
                        {{ getSlotTypeLabel(slot.slotType) }}
                      </span>
                    </td>
                    <td>
                      <div class="time-display">
                        {{ formatTime(slot.startTime) }} - {{ formatTime(slot.endTime) }}
                      </div>
                    </td>
                    <td>{{ slot.durationMinutes }} min</td>
                    <td>{{ slot.dayPatternName || '-' }}</td>
                    <td>{{ slot.shiftName || '-' }}</td>
                    <td>
                      <span class="badge" [class.badge-green]="slot.isActive" [class.badge-gray]="!slot.isActive">
                        {{ slot.isActive ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button class="action-btn" title="Edit" (click)="editSlot(slot)">
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button class="action-btn action-btn--danger" title="Delete" (click)="confirmDelete(slot)">
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        }
      </div>

      <!-- Period Slot Form Modal -->
      <msls-modal [isOpen]="showFormModal()" [title]="editingSlot() ? 'Edit Period Slot' : 'Create Period Slot'" size="lg" (closed)="closeFormModal()">
        <form class="form" (ngSubmit)="saveSlot()">
          <div class="form-row">
            <div class="form-group">
              <label for="slotName">Name <span class="required">*</span></label>
              <input type="text" id="slotName" [(ngModel)]="formData.name" name="name" placeholder="e.g., Period 1, Morning Break" required />
            </div>
            <div class="form-group">
              <label for="slotType">Type <span class="required">*</span></label>
              <select id="slotType" [(ngModel)]="formData.slotType" name="slotType" required>
                @for (type of slotTypes; track type.value) {
                  <option [value]="type.value">{{ type.label }}</option>
                }
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="startTime">Start Time <span class="required">*</span></label>
              <input type="time" id="startTime" [(ngModel)]="formData.startTime" name="startTime" required (change)="calculateDuration()" />
            </div>
            <div class="form-group">
              <label for="endTime">End Time <span class="required">*</span></label>
              <input type="time" id="endTime" [(ngModel)]="formData.endTime" name="endTime" required (change)="calculateDuration()" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="periodNumber">Period Number</label>
              <input type="number" id="periodNumber" [(ngModel)]="formData.periodNumber" name="periodNumber" min="0" placeholder="For teaching periods" />
            </div>
            <div class="form-group">
              <label for="durationMinutes">Duration (minutes)</label>
              <input type="number" id="durationMinutes" [(ngModel)]="formData.durationMinutes" name="durationMinutes" min="1" readonly />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="branchId">Branch <span class="required">*</span></label>
              <select id="branchId" [(ngModel)]="formData.branchId" name="branchId" required>
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

          <div class="form-row">
            <div class="form-group">
              <label for="dayPatternId">Day Pattern</label>
              <select id="dayPatternId" [(ngModel)]="formData.dayPatternId" name="dayPatternId">
                <option value="">All Patterns (Default)</option>
                @for (pattern of dayPatterns(); track pattern.id) {
                  <option [value]="pattern.id">{{ pattern.name }}</option>
                }
              </select>
            </div>
            <div class="form-group">
              <label for="shiftId">Shift</label>
              <select id="shiftId" [(ngModel)]="formData.shiftId" name="shiftId">
                <option value="">All Shifts (Default)</option>
                @for (shift of shifts(); track shift.id) {
                  <option [value]="shift.id">{{ shift.name }}</option>
                }
              </select>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" (click)="closeFormModal()">Cancel</button>
            <button type="submit" class="btn btn-primary" [disabled]="saving()">
              @if (saving()) {
                <div class="btn-spinner"></div>
                Saving...
              } @else {
                {{ editingSlot() ? 'Update' : 'Create' }}
              }
            </button>
          </div>
        </form>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal [isOpen]="showDeleteModal()" title="Delete Period Slot" size="sm" (closed)="closeDeleteModal()">
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>Are you sure you want to delete <strong>"{{ slotToDelete()?.name }}"</strong>?</p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteSlot()">
              @if (deleting()) { Deleting... } @else { Delete }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1600px; margin: 0 auto; }
    .page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
    .header-content { display: flex; align-items: center; gap: 1rem; }
    .header-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; background: #dcfce7; color: #16a34a; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
    .header-text h1 { margin: 0; font-size: 1.5rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }
    .filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; }
    .filter-group { display: flex; flex-direction: column; gap: 0.25rem; }
    .filter-label { font-size: 0.75rem; font-weight: 500; color: #64748b; }
    .filter-select { padding: 0.5rem 2rem 0.5rem 0.75rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; background: white; min-width: 140px; }
    .content-card { background: white; border: 1px solid #e2e8f0; border-radius: 1rem; overflow: hidden; }
    .loading-container, .error-container { display: flex; align-items: center; justify-content: center; gap: 1rem; padding: 3rem; color: #64748b; }
    .spinner { width: 24px; height: 24px; border: 3px solid #e2e8f0; border-top-color: #4f46e5; border-radius: 50%; animation: spin 0.8s linear infinite; }
    @keyframes spin { to { transform: rotate(360deg); } }

    /* Timeline View */
    .timeline-container { padding: 1.5rem; border-bottom: 1px solid #e2e8f0; }
    .timeline { display: flex; gap: 0.75rem; overflow-x: auto; padding-bottom: 0.5rem; }
    .slot-card { min-width: 140px; padding: 0.75rem; border-radius: 0.75rem; border: 1px solid #e2e8f0; background: white; cursor: pointer; transition: all 0.2s; position: relative; }
    .slot-card:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.08); }
    .slot-card.slot-regular { border-left: 3px solid #3b82f6; }
    .slot-card.slot-short { border-left: 3px solid #8b5cf6; }
    .slot-card.slot-assembly { border-left: 3px solid #f59e0b; }
    .slot-card.slot-break { border-left: 3px solid #22c55e; }
    .slot-card.slot-lunch { border-left: 3px solid #f97316; }
    .slot-card.slot-activity { border-left: 3px solid #14b8a6; }
    .slot-card.slot-zero_period { border-left: 3px solid #6b7280; }
    .slot-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.375rem; }
    .slot-number { font-size: 0.6875rem; font-weight: 600; color: #64748b; text-transform: uppercase; }
    .slot-type-badge { width: 1.5rem; height: 1.5rem; border-radius: 0.375rem; display: flex; align-items: center; justify-content: center; font-size: 0.6875rem; }
    .slot-type-badge.bg-blue-100 { background: #dbeafe; color: #1d4ed8; }
    .slot-type-badge.bg-purple-100 { background: #ede9fe; color: #7c3aed; }
    .slot-type-badge.bg-amber-100 { background: #fef3c7; color: #d97706; }
    .slot-type-badge.bg-green-100 { background: #dcfce7; color: #16a34a; }
    .slot-type-badge.bg-orange-100 { background: #ffedd5; color: #ea580c; }
    .slot-type-badge.bg-teal-100 { background: #ccfbf1; color: #0d9488; }
    .slot-type-badge.bg-gray-100 { background: #f1f5f9; color: #64748b; }
    .slot-name { font-weight: 600; color: #1e293b; font-size: 0.875rem; margin-bottom: 0.25rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    .slot-time { font-size: 0.75rem; color: #64748b; }
    .slot-duration { font-size: 0.6875rem; color: #94a3b8; margin-top: 0.25rem; }
    .slot-actions { position: absolute; top: 0.5rem; right: 0.5rem; display: none; gap: 0.25rem; }
    .slot-card:hover .slot-actions { display: flex; }
    .slot-action-btn { width: 1.5rem; height: 1.5rem; border: none; background: white; color: #64748b; border-radius: 0.25rem; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 0.6875rem; box-shadow: 0 1px 2px rgba(0,0,0,0.1); }
    .slot-action-btn:hover { color: #4f46e5; }
    .slot-action-btn--danger:hover { color: #dc2626; }
    .empty-timeline { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; padding: 2rem; color: #64748b; }

    /* Table View */
    .table-container { overflow-x: auto; }
    .data-table { width: 100%; border-collapse: collapse; }
    .data-table th { text-align: left; padding: 0.75rem 1rem; font-size: 0.6875rem; font-weight: 600; text-transform: uppercase; color: #64748b; background: #f8fafc; border-bottom: 1px solid #e2e8f0; }
    .data-table td { padding: 0.75rem 1rem; border-bottom: 1px solid #f1f5f9; color: #374151; font-size: 0.875rem; }
    .data-table tbody tr:hover { background: #f8fafc; }
    .order-badge { display: inline-flex; align-items: center; justify-content: center; width: 1.5rem; height: 1.5rem; background: #f1f5f9; border-radius: 0.25rem; font-size: 0.75rem; font-weight: 600; color: #64748b; }
    .name-wrapper { display: flex; align-items: center; gap: 0.625rem; }
    .slot-icon { width: 2rem; height: 2rem; border-radius: 0.375rem; display: flex; align-items: center; justify-content: center; font-size: 0.75rem; }
    .name-content { display: flex; flex-direction: column; }
    .name { font-weight: 500; color: #1e293b; }
    .description { font-size: 0.6875rem; color: #64748b; }
    .type-badge { display: inline-flex; padding: 0.25rem 0.5rem; border-radius: 0.25rem; font-size: 0.6875rem; font-weight: 500; }
    .time-display { font-size: 0.8125rem; }
    .badge { display: inline-flex; padding: 0.25rem 0.5rem; border-radius: 9999px; font-size: 0.6875rem; font-weight: 500; }
    .badge-green { background: #dcfce7; color: #166534; }
    .badge-gray { background: #f1f5f9; color: #64748b; }
    .actions-cell { text-align: right; }
    .action-btn { display: inline-flex; align-items: center; justify-content: center; width: 1.75rem; height: 1.75rem; border: none; background: transparent; color: #64748b; border-radius: 0.25rem; cursor: pointer; transition: all 0.2s; }
    .action-btn:hover { background: #f1f5f9; color: #4f46e5; }
    .action-btn--danger:hover { background: #fef2f2; color: #dc2626; }

    .btn { display: inline-flex; align-items: center; gap: 0.5rem; padding: 0.5rem 1rem; border-radius: 0.5rem; font-size: 0.875rem; font-weight: 500; cursor: pointer; border: none; }
    .btn-primary { background: #4f46e5; color: white; }
    .btn-primary:hover:not(:disabled) { background: #4338ca; }
    .btn-secondary { background: #f1f5f9; color: #475569; }
    .btn-danger { background: #dc2626; color: white; }
    .btn-sm { padding: 0.375rem 0.75rem; font-size: 0.8125rem; }
    .btn:disabled { opacity: 0.5; cursor: not-allowed; }
    .btn-spinner { width: 14px; height: 14px; border: 2px solid transparent; border-top-color: currentColor; border-radius: 50%; animation: spin 0.8s linear infinite; }
    .form { display: flex; flex-direction: column; gap: 1rem; }
    .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.25rem; }
    .form-group label { font-size: 0.8125rem; font-weight: 500; color: #374151; }
    .required { color: #dc2626; }
    .form-group input, .form-group select { padding: 0.5rem 0.75rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; }
    .form-group input:focus, .form-group select:focus { outline: none; border-color: #4f46e5; }
    .form-group input[readonly] { background: #f8fafc; color: #64748b; }
    .form-actions { display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 0.5rem; padding-top: 1rem; border-top: 1px solid #e2e8f0; }
    .delete-confirmation { text-align: center; padding: 1rem; }
    .delete-icon { width: 3.5rem; height: 3.5rem; margin: 0 auto 1rem; border-radius: 50%; background: #fef2f2; color: #dc2626; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
    .delete-warning { font-size: 0.8125rem; color: #64748b; }
    .delete-actions { display: flex; gap: 0.75rem; justify-content: center; margin-top: 1.5rem; }
  `],
})
export class PeriodSlotsComponent implements OnInit {
  private timetableService = inject(TimetableService);
  private branchService = inject(BranchService);
  private toastService = inject(ToastService);

  periodSlots = signal<PeriodSlot[]>([]);
  dayPatterns = signal<DayPattern[]>([]);
  shifts = signal<Shift[]>([]);
  branches = signal<Branch[]>([]);

  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);

  selectedBranchId = signal<string>('');
  selectedDayPatternId = signal<string>('');
  selectedShiftId = signal<string>('');
  selectedSlotType = signal<string>('');

  showFormModal = signal(false);
  showDeleteModal = signal(false);
  editingSlot = signal<PeriodSlot | null>(null);
  slotToDelete = signal<PeriodSlot | null>(null);

  slotTypes = PERIOD_SLOT_TYPES;

  formData = {
    name: '',
    slotType: 'regular' as PeriodSlotType,
    startTime: '08:00',
    endTime: '08:45',
    periodNumber: null as number | null,
    durationMinutes: 45,
    branchId: '',
    dayPatternId: '',
    shiftId: '',
    displayOrder: 0,
  };

  filteredSlots = computed(() => {
    let result = this.periodSlots();
    if (this.selectedBranchId()) result = result.filter(s => s.branchId === this.selectedBranchId());
    if (this.selectedDayPatternId()) result = result.filter(s => s.dayPatternId === this.selectedDayPatternId());
    if (this.selectedShiftId()) result = result.filter(s => s.shiftId === this.selectedShiftId());
    if (this.selectedSlotType()) result = result.filter(s => s.slotType === this.selectedSlotType());
    return result.sort((a, b) => a.displayOrder - b.displayOrder);
  });

  ngOnInit(): void {
    this.loadPeriodSlots();
    this.loadDayPatterns();
    this.loadShifts();
    this.loadBranches();
  }

  loadPeriodSlots(): void {
    this.loading.set(true);
    this.timetableService.getPeriodSlots().subscribe({
      next: slots => { this.periodSlots.set(slots); this.loading.set(false); },
      error: () => { this.error.set('Failed to load period slots'); this.loading.set(false); },
    });
  }

  loadDayPatterns(): void {
    this.timetableService.getDayPatterns({ isActive: true }).subscribe({
      next: patterns => this.dayPatterns.set(patterns),
      error: () => console.error('Failed to load day patterns'),
    });
  }

  loadShifts(): void {
    this.timetableService.getShifts({ isActive: true }).subscribe({
      next: shifts => this.shifts.set(shifts),
      error: () => console.error('Failed to load shifts'),
    });
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: allBranches => {
        const branches = allBranches.filter(b => b.isActive);
        this.branches.set(branches);
        if (branches.length === 1) {
          this.selectedBranchId.set(branches[0].id);
          this.formData.branchId = branches[0].id;
        }
      },
      error: () => console.error('Failed to load branches'),
    });
  }

  onBranchChange(branchId: string): void {
    this.selectedBranchId.set(branchId);
  }

  openCreateModal(): void {
    this.editingSlot.set(null);
    const nextOrder = this.periodSlots().length;
    this.formData = {
      name: '',
      slotType: 'regular',
      startTime: '08:00',
      endTime: '08:45',
      periodNumber: null,
      durationMinutes: 45,
      branchId: this.selectedBranchId() || (this.branches().length === 1 ? this.branches()[0].id : ''),
      dayPatternId: this.selectedDayPatternId() || '',
      shiftId: this.selectedShiftId() || '',
      displayOrder: nextOrder,
    };
    this.showFormModal.set(true);
  }

  editSlot(slot: PeriodSlot): void {
    this.editingSlot.set(slot);
    this.formData = {
      name: slot.name,
      slotType: slot.slotType,
      startTime: slot.startTime,
      endTime: slot.endTime,
      periodNumber: slot.periodNumber || null,
      durationMinutes: slot.durationMinutes,
      branchId: slot.branchId,
      dayPatternId: slot.dayPatternId || '',
      shiftId: slot.shiftId || '',
      displayOrder: slot.displayOrder,
    };
    this.showFormModal.set(true);
  }

  closeFormModal(): void {
    this.showFormModal.set(false);
    this.editingSlot.set(null);
  }

  calculateDuration(): void {
    if (this.formData.startTime && this.formData.endTime) {
      const [sh, sm] = this.formData.startTime.split(':').map(Number);
      const [eh, em] = this.formData.endTime.split(':').map(Number);
      let minutes = (eh * 60 + em) - (sh * 60 + sm);
      if (minutes < 0) minutes += 24 * 60;
      this.formData.durationMinutes = minutes;
    }
  }

  saveSlot(): void {
    if (!this.formData.name || !this.formData.branchId || !this.formData.startTime || !this.formData.endTime) {
      this.toastService.error('Please fill in all required fields');
      return;
    }
    this.saving.set(true);
    const editing = this.editingSlot();

    if (editing) {
      const data: UpdatePeriodSlotRequest = {
        name: this.formData.name,
        slotType: this.formData.slotType,
        startTime: this.formData.startTime,
        endTime: this.formData.endTime,
        periodNumber: this.formData.periodNumber || undefined,
        durationMinutes: this.formData.durationMinutes,
        dayPatternId: this.formData.dayPatternId || undefined,
        shiftId: this.formData.shiftId || undefined,
        displayOrder: this.formData.displayOrder,
      };
      this.timetableService.updatePeriodSlot(editing.id, data).subscribe({
        next: () => {
          this.toastService.success('Period slot updated');
          this.closeFormModal();
          this.loadPeriodSlots();
          this.saving.set(false);
        },
        error: () => {
          this.toastService.error('Failed to update period slot');
          this.saving.set(false);
        },
      });
    } else {
      const data: CreatePeriodSlotRequest = {
        branchId: this.formData.branchId,
        name: this.formData.name,
        slotType: this.formData.slotType,
        startTime: this.formData.startTime,
        endTime: this.formData.endTime,
        periodNumber: this.formData.periodNumber || undefined,
        durationMinutes: this.formData.durationMinutes,
        dayPatternId: this.formData.dayPatternId || undefined,
        shiftId: this.formData.shiftId || undefined,
        displayOrder: this.formData.displayOrder,
      };
      this.timetableService.createPeriodSlot(data).subscribe({
        next: () => {
          this.toastService.success('Period slot created');
          this.closeFormModal();
          this.loadPeriodSlots();
          this.saving.set(false);
        },
        error: () => {
          this.toastService.error('Failed to create period slot');
          this.saving.set(false);
        },
      });
    }
  }

  confirmDelete(slot: PeriodSlot): void {
    this.slotToDelete.set(slot);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.slotToDelete.set(null);
  }

  deleteSlot(): void {
    const slot = this.slotToDelete();
    if (!slot) return;
    this.deleting.set(true);
    this.timetableService.deletePeriodSlot(slot.id).subscribe({
      next: () => {
        this.toastService.success('Period slot deleted');
        this.closeDeleteModal();
        this.loadPeriodSlots();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete period slot');
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

  getSlotTypeLabel(type: PeriodSlotType): string {
    return this.slotTypes.find(t => t.value === type)?.label || type;
  }

  getSlotTypeIcon(type: PeriodSlotType): string {
    return this.slotTypes.find(t => t.value === type)?.icon || 'fa-clock';
  }

  getSlotTypeClass(type: PeriodSlotType): string {
    return `slot-${type}`;
  }

  getSlotTypeBadgeClass(type: PeriodSlotType): string {
    return this.slotTypes.find(t => t.value === type)?.color || 'bg-gray-100 text-gray-700';
  }
}
