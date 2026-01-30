import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink } from '@angular/router';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators, FormArray } from '@angular/forms';
import { SubstitutionService } from '../substitution.service';
import { BranchService } from '../../../../admin/branches/branch.service';
import { StaffService } from '../../../../staff/services/staff.service';
import {
  CreateSubstitutionRequest,
  AvailableTeacher,
  TeacherPeriod,
} from '../substitution.model';
import { Staff } from '../../../../staff/models/staff.model';

@Component({
  selector: 'app-substitution-form',
  standalone: true,
  imports: [CommonModule, RouterLink, ReactiveFormsModule],
  template: `
    <div class="max-w-4xl mx-auto space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-gray-900">New Substitution</h1>
          <p class="mt-1 text-sm text-gray-500">
            Create a teacher substitution for an absent teacher
          </p>
        </div>
        <a
          routerLink="../"
          class="text-sm font-medium text-gray-500 hover:text-gray-700"
        >
          Cancel
        </a>
      </div>

      <!-- Form Steps -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200">
        <!-- Progress Indicator -->
        <div class="border-b border-gray-200 px-6 py-4">
          <nav class="flex justify-center">
            <ol class="flex items-center space-x-8">
              @for (step of steps; track step.id; let i = $index) {
                <li class="flex items-center">
                  <div
                    [class]="currentStep() > i ? 'bg-indigo-600' : currentStep() === i ? 'bg-indigo-600' : 'bg-gray-300'"
                    class="relative flex h-8 w-8 items-center justify-center rounded-full"
                  >
                    @if (currentStep() > i) {
                      <svg class="h-5 w-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                    } @else {
                      <span class="text-sm font-medium text-white">{{ i + 1 }}</span>
                    }
                  </div>
                  <span
                    [class]="currentStep() >= i ? 'text-indigo-600' : 'text-gray-500'"
                    class="ml-3 text-sm font-medium hidden sm:block"
                  >
                    {{ step.name }}
                  </span>
                </li>
              }
            </ol>
          </nav>
        </div>

        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="p-6">
          <!-- Step 1: Select Teacher & Date -->
          @if (currentStep() === 0) {
            <div class="space-y-6">
              <h2 class="text-lg font-medium text-gray-900">Select Absent Teacher</h2>

              <!-- Branch Selection -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Branch <span class="text-red-500">*</span>
                </label>
                <select
                  formControlName="branchId"
                  (change)="onBranchChange()"
                  class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                >
                  <option value="">Select a branch</option>
                  @for (branch of branches(); track branch.id) {
                    <option [value]="branch.id">{{ branch.name }}</option>
                  }
                </select>
              </div>

              <!-- Date Selection -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Substitution Date <span class="text-red-500">*</span>
                </label>
                <input
                  type="date"
                  formControlName="substitutionDate"
                  [min]="today"
                  class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <!-- Teacher Selection -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Absent Teacher <span class="text-red-500">*</span>
                </label>
                <select
                  formControlName="originalStaffId"
                  (change)="onTeacherChange()"
                  [disabled]="!form.get('branchId')?.value"
                  class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm disabled:bg-gray-100"
                >
                  <option value="">Select a teacher</option>
                  @for (teacher of teachers(); track teacher.id) {
                    <option [value]="teacher.id">{{ teacher.firstName }} {{ teacher.lastName }}</option>
                  }
                </select>
              </div>

              <!-- Reason -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Reason for Absence</label>
                <input
                  type="text"
                  formControlName="reason"
                  placeholder="e.g., Sick leave, Personal emergency"
                  class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>
            </div>
          }

          <!-- Step 2: Select Periods -->
          @if (currentStep() === 1) {
            <div class="space-y-6">
              <h2 class="text-lg font-medium text-gray-900">Select Periods to Cover</h2>

              @if (loadingPeriods()) {
                <div class="flex justify-center py-8">
                  <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
                </div>
              } @else if (teacherPeriods().length === 0) {
                <div class="text-center py-8 bg-gray-50 rounded-lg">
                  <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  <p class="mt-2 text-sm text-gray-500">This teacher has no scheduled periods on this day.</p>
                </div>
              } @else {
                <div class="space-y-3">
                  @for (period of teacherPeriods(); track period.id; let i = $index) {
                    <label
                      [class]="isSelected(period) ? 'border-indigo-500 ring-2 ring-indigo-500' : 'border-gray-200'"
                      class="flex items-center p-4 bg-white border rounded-lg cursor-pointer hover:border-indigo-300 transition-colors"
                    >
                      <input
                        type="checkbox"
                        [checked]="isSelected(period)"
                        (change)="togglePeriod(period)"
                        class="h-4 w-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
                      />
                      <div class="ml-4 flex-1">
                        <div class="flex items-center justify-between">
                          <span class="text-sm font-medium text-gray-900">
                            {{ period.periodSlotName }}
                          </span>
                          <span class="text-sm text-gray-500">
                            {{ period.startTime }} - {{ period.endTime }}
                          </span>
                        </div>
                        <div class="mt-1 text-sm text-gray-500">
                          {{ period.className }} {{ period.sectionName }} - {{ period.subjectName }}
                        </div>
                      </div>
                    </label>
                  }
                </div>

                <p class="text-sm text-gray-500 mt-4">
                  Selected: {{ selectedPeriods().length }} period(s)
                </p>
              }
            </div>
          }

          <!-- Step 3: Assign Substitute -->
          @if (currentStep() === 2) {
            <div class="space-y-6">
              <h2 class="text-lg font-medium text-gray-900">Assign Substitute Teacher</h2>

              @if (loadingTeachers()) {
                <div class="flex justify-center py-8">
                  <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
                </div>
              } @else {
                <div class="space-y-3">
                  @for (teacher of availableTeachers(); track teacher.staffId) {
                    <label
                      [class]="form.get('substituteStaffId')?.value === teacher.staffId ? 'border-indigo-500 ring-2 ring-indigo-500' : teacher.hasConflict ? 'border-red-200 bg-red-50' : 'border-gray-200'"
                      class="flex items-center p-4 border rounded-lg cursor-pointer hover:border-indigo-300 transition-colors"
                      [class.cursor-not-allowed]="teacher.hasConflict"
                    >
                      <input
                        type="radio"
                        formControlName="substituteStaffId"
                        [value]="teacher.staffId"
                        [disabled]="teacher.hasConflict"
                        class="h-4 w-4 text-indigo-600 border-gray-300 focus:ring-indigo-500 disabled:text-gray-400"
                      />
                      <div class="ml-4 flex-1">
                        <div class="flex items-center justify-between">
                          <span class="text-sm font-medium" [class]="teacher.hasConflict ? 'text-gray-400' : 'text-gray-900'">
                            {{ teacher.staffName }}
                          </span>
                          @if (teacher.hasConflict) {
                            <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                              Has Conflict
                            </span>
                          } @else {
                            <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                              Available
                            </span>
                          }
                        </div>
                        <div class="mt-1 text-sm" [class]="teacher.hasConflict ? 'text-gray-400' : 'text-gray-500'">
                          {{ teacher.department || 'No department' }} · {{ teacher.freePeriods }} free periods
                        </div>
                      </div>
                    </label>
                  }

                  @if (availableTeachers().length === 0) {
                    <div class="text-center py-8 bg-gray-50 rounded-lg">
                      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                      </svg>
                      <p class="mt-2 text-sm text-gray-500">No available teachers found for this time slot.</p>
                    </div>
                  }
                </div>

                <!-- Notes -->
                <div class="mt-6">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Notes (optional)</label>
                  <textarea
                    formControlName="notes"
                    rows="3"
                    placeholder="Any additional instructions for the substitute..."
                    class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                  ></textarea>
                </div>
              }
            </div>
          }

          <!-- Navigation Buttons -->
          <div class="flex justify-between mt-8 pt-6 border-t border-gray-200">
            <button
              type="button"
              (click)="previousStep()"
              [disabled]="currentStep() === 0"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>

            @if (currentStep() < steps.length - 1) {
              <button
                type="button"
                (click)="nextStep()"
                [disabled]="!canProceed()"
                class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next
              </button>
            } @else {
              <button
                type="submit"
                [disabled]="!form.valid || submitting()"
                class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                @if (submitting()) {
                  <span class="flex items-center">
                    <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Creating...
                  </span>
                } @else {
                  Create Substitution
                }
              </button>
            }
          </div>
        </form>
      </div>
    </div>
  `,
})
export class SubstitutionFormComponent implements OnInit {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private substitutionService = inject(SubstitutionService);
  private branchService = inject(BranchService);
  private staffService = inject(StaffService);

  // Steps
  steps = [
    { id: 0, name: 'Select Teacher' },
    { id: 1, name: 'Choose Periods' },
    { id: 2, name: 'Assign Substitute' },
  ];
  currentStep = signal(0);

  // State
  branches = signal<{ id: string; name: string }[]>([]);
  teachers = signal<Staff[]>([]);
  teacherPeriods = signal<TeacherPeriod[]>([]);
  availableTeachers = signal<AvailableTeacher[]>([]);
  selectedPeriods = signal<TeacherPeriod[]>([]);
  loadingPeriods = signal(false);
  loadingTeachers = signal(false);
  submitting = signal(false);

  // Today's date for min date validation
  today = new Date().toISOString().split('T')[0];

  // Form
  form: FormGroup = this.fb.group({
    branchId: ['', Validators.required],
    originalStaffId: ['', Validators.required],
    substituteStaffId: ['', Validators.required],
    substitutionDate: ['', Validators.required],
    reason: [''],
    notes: [''],
    periods: this.fb.array([]),
  });

  get periodsArray(): FormArray {
    return this.form.get('periods') as FormArray;
  }

  ngOnInit(): void {
    this.branchService.getBranches().subscribe((branches) => this.branches.set(branches));
    // Set default date to today
    this.form.patchValue({ substitutionDate: this.today });
  }

  onBranchChange(): void {
    const branchId = this.form.get('branchId')?.value;
    if (branchId) {
      this.staffService.loadStaff({ branchId, staffType: 'teaching', status: 'active' }).subscribe({
        next: (response) => {
          this.teachers.set(response.staff);
        },
      });
    }
    // Reset downstream selections
    this.form.patchValue({ originalStaffId: '', substituteStaffId: '' });
    this.teacherPeriods.set([]);
    this.selectedPeriods.set([]);
    this.availableTeachers.set([]);
  }

  onTeacherChange(): void {
    this.form.patchValue({ substituteStaffId: '' });
    this.selectedPeriods.set([]);
    this.availableTeachers.set([]);
  }

  loadTeacherPeriods(): void {
    const staffId = this.form.get('originalStaffId')?.value;
    const date = this.form.get('substitutionDate')?.value;

    if (staffId && date) {
      this.loadingPeriods.set(true);
      this.substitutionService.getTeacherPeriods(staffId, date).subscribe({
        next: (periods) => {
          this.teacherPeriods.set(periods);
          this.loadingPeriods.set(false);
        },
        error: () => {
          this.loadingPeriods.set(false);
        },
      });
    }
  }

  isSelected(period: TeacherPeriod): boolean {
    return this.selectedPeriods().some((p) => p.periodSlotId === period.periodSlotId);
  }

  togglePeriod(period: TeacherPeriod): void {
    const selected = this.selectedPeriods();
    const index = selected.findIndex((p) => p.periodSlotId === period.periodSlotId);

    if (index >= 0) {
      this.selectedPeriods.set([...selected.slice(0, index), ...selected.slice(index + 1)]);
    } else {
      this.selectedPeriods.set([...selected, period]);
    }
  }

  loadAvailableTeachers(): void {
    const branchId = this.form.get('branchId')?.value;
    const date = this.form.get('substitutionDate')?.value;
    const originalStaffId = this.form.get('originalStaffId')?.value;
    const periodSlotIds = this.selectedPeriods().map((p) => p.periodSlotId);

    if (branchId && date && periodSlotIds.length > 0) {
      this.loadingTeachers.set(true);
      this.substitutionService.getAvailableTeachers(branchId, date, periodSlotIds, originalStaffId).subscribe({
        next: (response) => {
          this.availableTeachers.set(response.teachers);
          this.loadingTeachers.set(false);
        },
        error: () => {
          this.loadingTeachers.set(false);
        },
      });
    }
  }

  canProceed(): boolean {
    switch (this.currentStep()) {
      case 0:
        return (
          !!this.form.get('branchId')?.value &&
          !!this.form.get('originalStaffId')?.value &&
          !!this.form.get('substitutionDate')?.value
        );
      case 1:
        return this.selectedPeriods().length > 0;
      case 2:
        return !!this.form.get('substituteStaffId')?.value;
      default:
        return false;
    }
  }

  previousStep(): void {
    if (this.currentStep() > 0) {
      this.currentStep.update((s) => s - 1);
    }
  }

  nextStep(): void {
    if (this.currentStep() === 0) {
      this.loadTeacherPeriods();
    } else if (this.currentStep() === 1) {
      this.loadAvailableTeachers();
    }

    if (this.currentStep() < this.steps.length - 1) {
      this.currentStep.update((s) => s + 1);
    }
  }

  onSubmit(): void {
    if (this.form.invalid || this.selectedPeriods().length === 0) {
      return;
    }

    this.submitting.set(true);

    const request: CreateSubstitutionRequest = {
      branchId: this.form.get('branchId')?.value,
      originalStaffId: this.form.get('originalStaffId')?.value,
      substituteStaffId: this.form.get('substituteStaffId')?.value,
      substitutionDate: this.form.get('substitutionDate')?.value,
      reason: this.form.get('reason')?.value || undefined,
      notes: this.form.get('notes')?.value || undefined,
      periods: this.selectedPeriods().map((p) => ({
        periodSlotId: p.periodSlotId,
        timetableEntryId: p.id,
        subjectId: p.subjectId,
        sectionId: p.sectionId,
        roomNumber: p.roomNumber,
      })),
    };

    this.substitutionService.createSubstitution(request).subscribe({
      next: (substitution) => {
        this.submitting.set(false);
        this.router.navigate(['../', substitution.id], { relativeTo: this.router.routerState.root.firstChild });
      },
      error: (err) => {
        this.submitting.set(false);
        alert('Failed to create substitution: ' + err.message);
      },
    });
  }
}
