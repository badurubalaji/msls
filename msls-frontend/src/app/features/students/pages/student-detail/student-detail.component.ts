/**
 * Student Detail Page Component
 *
 * Displays comprehensive student profile information with quick actions.
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
import { Router, ActivatedRoute } from '@angular/router';

import {
  MslsBadgeComponent,
  MslsAvatarComponent,
  MslsSpinnerComponent,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services';
import { StudentService } from '../../services/student.service';
import { EnrollmentService } from '../../services/enrollment.service';
import { GuardianSectionComponent } from '../../components/guardian-section/guardian-section';
import { HealthSectionComponent } from '../../components/health-section/health-section';
import { BehaviorSectionComponent } from '../../components/behavior-section/behavior-section';
import { DocumentSectionComponent } from '../../components/document-section/document-section.component';
import { EnrollmentHistoryComponent } from '../../components/enrollment-history/enrollment-history.component';
import { EnrollmentFormComponent } from '../../components/enrollment-form/enrollment-form.component';
import { TransferFormComponent, TransferFormMode } from '../../components/transfer-form/transfer-form.component';
import {
  Student,
  getStatusBadgeVariant,
  getStatusLabel,
  getGenderLabel,
  calculateAge,
  formatAddress,
} from '../../models/student.model';
import { Enrollment } from '../../models/enrollment.model';

@Component({
  selector: 'msls-student-detail',
  standalone: true,
  imports: [
    CommonModule,
    MslsBadgeComponent,
    MslsAvatarComponent,
    MslsSpinnerComponent,
    GuardianSectionComponent,
    HealthSectionComponent,
    BehaviorSectionComponent,
    DocumentSectionComponent,
    EnrollmentHistoryComponent,
    EnrollmentFormComponent,
    TransferFormComponent,
  ],
  templateUrl: './student-detail.component.html',
  styleUrl: './student-detail.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StudentDetailComponent implements OnInit {
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private studentService = inject(StudentService);
  private enrollmentService = inject(EnrollmentService);
  private toast = inject(ToastService);

  // =========================================================================
  // State
  // =========================================================================

  readonly studentId = signal<string>('');
  readonly loading = this.studentService.loading;
  readonly error = this.studentService.error;
  readonly student = this.studentService.selectedStudent;

  /** Enrollment form state */
  readonly showEnrollmentForm = signal<boolean>(false);
  readonly editingEnrollment = signal<Enrollment | null>(null);

  /** Transfer/dropout modal state */
  readonly showTransferModal = signal<boolean>(false);
  readonly transferMode = signal<TransferFormMode>('transfer');
  readonly transferEnrollment = signal<Enrollment | null>(null);

  /** Mock academic years (TODO: fetch from API in Epic 7) */
  readonly academicYears = signal<Array<{ id: string; name: string }>>([
    { id: 'ay-2024-2025', name: '2024-2025' },
    { id: 'ay-2023-2024', name: '2023-2024' },
  ]);

  // =========================================================================
  // Computed
  // =========================================================================

  readonly age = computed(() => {
    const s = this.student();
    return s ? calculateAge(s.dateOfBirth) : 0;
  });

  readonly currentAddressFormatted = computed(() => {
    const s = this.student();
    return s ? formatAddress(s.currentAddress) : '';
  });

  readonly permanentAddressFormatted = computed(() => {
    const s = this.student();
    return s ? formatAddress(s.permanentAddress) : '';
  });

  readonly statusVariant = computed(() => {
    const s = this.student();
    return s ? getStatusBadgeVariant(s.status) : 'neutral';
  });

  readonly statusLabel = computed(() => {
    const s = this.student();
    return s ? getStatusLabel(s.status) : '';
  });

  readonly genderLabel = computed(() => {
    const s = this.student();
    return s ? getGenderLabel(s.gender) : '';
  });

  // =========================================================================
  // Lifecycle
  // =========================================================================

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.studentId.set(id);
      this.loadStudent(id);
    } else {
      this.router.navigate(['/students']);
    }
  }

  // =========================================================================
  // Data Loading
  // =========================================================================

  private loadStudent(id: string): void {
    this.studentService.getStudent(id).subscribe({
      error: (err) => {
        this.toast.error('Failed to load student', err.message);
        this.router.navigate(['/students']);
      },
    });
  }

  refresh(): void {
    const id = this.studentId();
    if (id) {
      this.loadStudent(id);
    }
  }

  // =========================================================================
  // Actions
  // =========================================================================

  onEdit(): void {
    this.router.navigate(['/students', this.studentId(), 'edit']);
  }

  onDelete(): void {
    const s = this.student();
    if (!s) return;

    // Use browser confirm for now - can be replaced with modal component later
    const confirmed = confirm(`Are you sure you want to delete ${s.fullName}? This action cannot be undone.`);
    if (confirmed) {
      this.studentService.deleteStudent(s.id).subscribe({
        next: () => {
          this.toast.success('Student deleted successfully');
          this.router.navigate(['/students']);
        },
        error: (err) => {
          this.toast.error(`Failed to delete student: ${err.message}`);
        },
      });
    }
  }

  onViewAttendance(): void {
    // TODO: Navigate to attendance view
    this.toast.info('Coming soon: Attendance view is not yet implemented');
  }

  onViewFees(): void {
    // TODO: Navigate to fees view
    this.toast.info('Coming soon: Fees view is not yet implemented');
  }

  onViewDocuments(): void {
    // TODO: Navigate to documents view
    this.toast.info('Coming soon: Documents view is not yet implemented');
  }

  onViewAuditLog(): void {
    // TODO: Navigate to audit log
    this.toast.info('Coming soon: Audit log is not yet implemented');
  }

  // =========================================================================
  // Enrollment Actions
  // =========================================================================

  onCreateEnrollment(): void {
    this.editingEnrollment.set(null);
    this.showEnrollmentForm.set(true);
  }

  onEditEnrollment(enrollment: Enrollment): void {
    this.editingEnrollment.set(enrollment);
    this.showEnrollmentForm.set(true);
  }

  onEnrollmentSaved(enrollment: Enrollment): void {
    this.showEnrollmentForm.set(false);
    this.editingEnrollment.set(null);
    this.toast.success(
      this.editingEnrollment() ? 'Enrollment updated successfully' : 'Enrollment created successfully'
    );
    // Refresh enrollment history
    this.enrollmentService.loadEnrollmentHistory(this.studentId()).subscribe();
  }

  onEnrollmentFormCancelled(): void {
    this.showEnrollmentForm.set(false);
    this.editingEnrollment.set(null);
  }

  onTransferEnrollment(enrollment: Enrollment): void {
    this.transferEnrollment.set(enrollment);
    this.transferMode.set('transfer');
    this.showTransferModal.set(true);
  }

  onDropoutEnrollment(enrollment: Enrollment): void {
    this.transferEnrollment.set(enrollment);
    this.transferMode.set('dropout');
    this.showTransferModal.set(true);
  }

  onTransferSaved(enrollment: Enrollment): void {
    this.showTransferModal.set(false);
    this.transferEnrollment.set(null);
    this.toast.success(
      this.transferMode() === 'transfer' ? 'Transfer processed successfully' : 'Dropout processed successfully'
    );
    // Refresh enrollment history
    this.enrollmentService.loadEnrollmentHistory(this.studentId()).subscribe();
  }

  onTransferCancelled(): void {
    this.showTransferModal.set(false);
    this.transferEnrollment.set(null);
  }

  // =========================================================================
  // Helpers
  // =========================================================================

  formatDate(dateString: string): string {
    if (!dateString) return '—';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  goBack(): void {
    this.router.navigate(['/students']);
  }
}
