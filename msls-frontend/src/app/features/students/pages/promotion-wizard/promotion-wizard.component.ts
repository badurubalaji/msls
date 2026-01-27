/**
 * Promotion Wizard Component
 *
 * Multi-step wizard for processing student promotions/retentions.
 */

import {
  Component,
  ChangeDetectionStrategy,
  inject,
  OnInit,
  signal,
  computed,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';

import {
  MslsBadgeComponent,
  MslsSpinnerComponent,
  MslsAvatarComponent,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services';
import { PromotionService } from '../../services/promotion.service';
import {
  PromotionBatch,
  PromotionRecord,
  PromotionDecision,
  CreateBatchRequest,
  getDecisionBadgeVariant,
  getDecisionLabel,
  getBatchStatusBadgeVariant,
  getBatchStatusLabel,
} from '../../models/promotion.model';

interface AcademicYear {
  id: string;
  name: string;
}

interface ClassOption {
  id: string;
  name: string;
}

interface SectionOption {
  id: string;
  name: string;
}

@Component({
  selector: 'msls-promotion-wizard',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsBadgeComponent,
    MslsSpinnerComponent,
    MslsAvatarComponent,
  ],
  templateUrl: './promotion-wizard.component.html',
  styleUrl: './promotion-wizard.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PromotionWizardComponent implements OnInit {
  private router = inject(Router);
  private promotionService = inject(PromotionService);
  private toast = inject(ToastService);

  // =========================================================================
  // State
  // =========================================================================

  readonly loading = this.promotionService.loading;
  readonly error = this.promotionService.error;
  readonly batch = this.promotionService.selectedBatch;
  readonly records = this.promotionService.records;
  readonly summary = this.promotionService.summary;

  /** Current wizard step */
  readonly currentStep = signal<number>(0);

  /** Wizard steps */
  readonly steps = [
    { label: 'Select Source', icon: 'fa-search' },
    { label: 'Review Students', icon: 'fa-users' },
    { label: 'Assign Sections', icon: 'fa-sitemap' },
    { label: 'Confirm', icon: 'fa-check-circle' },
    { label: 'Complete', icon: 'fa-flag-checkered' },
  ];

  /** Source form */
  readonly sourceForm = signal({
    fromAcademicYearId: '',
    toAcademicYearId: '',
    fromClassId: '',
    fromSectionId: '',
    toClassId: '',
    notes: '',
  });

  /** Selected records for bulk operations */
  readonly selectedRecords = signal<Set<string>>(new Set());

  /** Mock data - TODO: Replace with API calls */
  readonly academicYears = signal<AcademicYear[]>([
    { id: 'ay-2024-2025', name: '2024-2025' },
    { id: 'ay-2025-2026', name: '2025-2026' },
  ]);

  readonly classes = signal<ClassOption[]>([
    { id: 'class-1', name: 'Class 1' },
    { id: 'class-2', name: 'Class 2' },
    { id: 'class-3', name: 'Class 3' },
    { id: 'class-4', name: 'Class 4' },
    { id: 'class-5', name: 'Class 5' },
  ]);

  readonly sections = signal<SectionOption[]>([
    { id: 'section-a', name: 'Section A' },
    { id: 'section-b', name: 'Section B' },
    { id: 'section-c', name: 'Section C' },
  ]);

  /** Generate roll numbers option */
  readonly generateRollNumbers = signal<boolean>(true);

  // =========================================================================
  // Computed
  // =========================================================================

  readonly canProceedFromSource = computed(() => {
    const form = this.sourceForm();
    return form.fromAcademicYearId && form.toAcademicYearId && form.fromClassId;
  });

  readonly canProceedFromReview = computed(() => {
    const pendingCount = this.summary()?.pendingCount ?? 0;
    return pendingCount === 0;
  });

  readonly allSelected = computed(() => {
    const records = this.records();
    const selected = this.selectedRecords();
    return records.length > 0 && records.every((r) => selected.has(r.id));
  });

  readonly someSelected = computed(() => {
    const selected = this.selectedRecords();
    return selected.size > 0;
  });

  readonly targetClass = computed(() => {
    const form = this.sourceForm();
    return this.classes().find((c) => c.id === form.toClassId);
  });

  readonly fromClass = computed(() => {
    const form = this.sourceForm();
    return this.classes().find((c) => c.id === form.fromClassId);
  });

  // =========================================================================
  // Lifecycle
  // =========================================================================

  ngOnInit(): void {
    // Clear any previously selected batch
    this.promotionService.clearSelected();
  }

  // =========================================================================
  // Navigation
  // =========================================================================

  nextStep(): void {
    if (this.currentStep() < this.steps.length - 1) {
      this.currentStep.update((s) => s + 1);
    }
  }

  prevStep(): void {
    if (this.currentStep() > 0) {
      this.currentStep.update((s) => s - 1);
    }
  }

  goToStep(step: number): void {
    if (step <= this.currentStep()) {
      this.currentStep.set(step);
    }
  }

  goBack(): void {
    this.router.navigate(['/students']);
  }

  // =========================================================================
  // Step 1: Source Selection
  // =========================================================================

  onSourceFormChange(field: string, value: string): void {
    this.sourceForm.update((f) => ({ ...f, [field]: value }));
  }

  async createBatch(): Promise<void> {
    const form = this.sourceForm();
    const request: CreateBatchRequest = {
      fromAcademicYearId: form.fromAcademicYearId,
      toAcademicYearId: form.toAcademicYearId,
      fromClassId: form.fromClassId,
      fromSectionId: form.fromSectionId || undefined,
      toClassId: form.toClassId || undefined,
      notes: form.notes || undefined,
    };

    this.promotionService.createBatch(request).subscribe({
      next: (batch) => {
        this.toast.success(`Batch created with ${batch.totalStudents} students`);
        // Load records
        this.promotionService.loadRecords(batch.id).subscribe({
          next: () => {
            this.nextStep();
          },
        });
      },
      error: (err) => {
        this.toast.error(err.error?.error || 'Failed to create batch');
      },
    });
  }

  // =========================================================================
  // Step 2: Review Students
  // =========================================================================

  toggleSelectAll(): void {
    const records = this.records();
    if (this.allSelected()) {
      this.selectedRecords.set(new Set());
    } else {
      this.selectedRecords.set(new Set(records.map((r) => r.id)));
    }
  }

  toggleRecord(recordId: string): void {
    const selected = new Set(this.selectedRecords());
    if (selected.has(recordId)) {
      selected.delete(recordId);
    } else {
      selected.add(recordId);
    }
    this.selectedRecords.set(selected);
  }

  isRecordSelected(recordId: string): boolean {
    return this.selectedRecords().has(recordId);
  }

  autoDecide(): void {
    const batch = this.batch();
    if (!batch) return;

    this.promotionService.autoDecide(batch.id).subscribe({
      next: () => {
        this.toast.success('Auto-decision applied based on promotion rules');
      },
      error: (err) => {
        this.toast.error(err.error?.error || 'Failed to auto-decide');
      },
    });
  }

  bulkSetDecision(decision: PromotionDecision): void {
    const batch = this.batch();
    if (!batch) return;

    const selected = Array.from(this.selectedRecords());
    if (selected.length === 0) {
      this.toast.warning('Please select students first');
      return;
    }

    this.promotionService
      .bulkUpdateRecords(batch.id, {
        recordIds: selected,
        decision: decision,
        toClassId: decision === 'promote' ? this.sourceForm().toClassId : undefined,
      })
      .subscribe({
        next: () => {
          this.toast.success(`${selected.length} students marked as ${getDecisionLabel(decision)}`);
          this.selectedRecords.set(new Set());
        },
        error: (err) => {
          this.toast.error(err.error?.error || 'Failed to update records');
        },
      });
  }

  setRecordDecision(record: PromotionRecord, decision: PromotionDecision): void {
    const batch = this.batch();
    if (!batch) return;

    this.promotionService
      .updateRecord(batch.id, record.id, {
        decision: decision,
        toClassId: decision === 'promote' ? this.sourceForm().toClassId : undefined,
      })
      .subscribe({
        next: () => {
          this.toast.success(`Student marked as ${getDecisionLabel(decision)}`);
        },
        error: (err) => {
          this.toast.error(err.error?.error || 'Failed to update record');
        },
      });
  }

  // =========================================================================
  // Step 3: Assign Sections
  // =========================================================================

  setRecordSection(record: PromotionRecord, sectionId: string): void {
    const batch = this.batch();
    if (!batch) return;

    this.promotionService
      .updateRecord(batch.id, record.id, {
        toSectionId: sectionId || undefined,
      })
      .subscribe({
        error: (err) => {
          this.toast.error(err.error?.error || 'Failed to update section');
        },
      });
  }

  // =========================================================================
  // Step 4: Confirm & Process
  // =========================================================================

  processBatch(): void {
    const batch = this.batch();
    if (!batch) return;

    this.promotionService
      .processBatch(batch.id, {
        generateRollNumbers: this.generateRollNumbers(),
      })
      .subscribe({
        next: (updatedBatch) => {
          this.toast.success(
            `Promotion completed! ${updatedBatch.promotedCount} promoted, ${updatedBatch.retainedCount} retained`
          );
          this.nextStep();
        },
        error: (err) => {
          this.toast.error(err.error?.error || 'Failed to process batch');
        },
      });
  }

  // =========================================================================
  // Step 5: Complete
  // =========================================================================

  viewReport(): void {
    const batch = this.batch();
    if (!batch) return;

    // TODO: Implement report download
    this.promotionService.getReport(batch.id).subscribe({
      next: (report) => {
        console.log('Report:', report);
        this.toast.info('Report downloaded');
      },
      error: (err) => {
        this.toast.error(err.error?.error || 'Failed to get report');
      },
    });
  }

  startNewBatch(): void {
    this.promotionService.clearSelected();
    this.currentStep.set(0);
    this.sourceForm.set({
      fromAcademicYearId: '',
      toAcademicYearId: '',
      fromClassId: '',
      fromSectionId: '',
      toClassId: '',
      notes: '',
    });
    this.selectedRecords.set(new Set());
  }

  // =========================================================================
  // Template Helpers
  // =========================================================================

  getDecisionVariant = getDecisionBadgeVariant;
  getDecisionLabel = getDecisionLabel;
  getBatchStatusVariant = getBatchStatusBadgeVariant;
  getBatchStatusLabel = getBatchStatusLabel;

  trackByRecordId(index: number, record: PromotionRecord): string {
    return record.id;
  }
}
