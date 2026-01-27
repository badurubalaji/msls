/**
 * MSLS Results Entry Component
 *
 * Component for entering and managing test results for candidates.
 * Also handles candidate registration when accessed via the registrations route.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { ToastService } from '../../../../../shared/services';
import { MslsModalComponent } from '../../../../../shared/components';
import {
  EntranceTest,
  TestRegistration,
  SubjectMarks,
  REGISTRATION_STATUS_CONFIG,
  calculateTotalMarks,
  calculateMaxMarks,
  calculatePercentage,
} from '../../entrance-test.model';
import { EntranceTestService } from '../../entrance-test.service';
import { ApplicationService } from '../../../applications/application.service';
import { AdmissionApplication } from '../../../applications/application.model';

interface ResultEntry {
  registrationId: string;
  rollNumber: string;
  studentName: string;
  marks: { [subject: string]: number };
  remarks: string;
  saved: boolean;
}

@Component({
  selector: 'msls-results-entry',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent],
  templateUrl: './results-entry.html',
  styleUrl: './results-entry.scss',
})
export class ResultsEntryComponent implements OnInit {
  private readonly testService = inject(EntranceTestService);
  private readonly applicationService = inject(ApplicationService);
  private readonly toastService = inject(ToastService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  // Mode detection (results vs registrations)
  isRegistrationsMode = signal(false);

  // Test data
  test = signal<EntranceTest | null>(null);
  registrations = signal<TestRegistration[]>([]);
  loading = signal(true);
  saving = signal(false);
  error = signal<string | null>(null);

  // Registration modal state
  showRegisterModal = signal(false);
  availableApplications = signal<AdmissionApplication[]>([]);
  loadingApplications = signal(false);
  selectedApplicationId = signal<string>('');
  registering = signal(false);

  // Results data
  resultEntries = signal<ResultEntry[]>([]);
  searchTerm = signal('');

  // Filter computed
  filteredEntries = computed(() => {
    const term = this.searchTerm().toLowerCase();
    if (!term) return this.resultEntries();
    return this.resultEntries().filter(
      entry =>
        entry.studentName.toLowerCase().includes(term) ||
        entry.rollNumber.toLowerCase().includes(term)
    );
  });

  // Stats computed
  stats = computed(() => {
    const entries = this.resultEntries();
    const saved = entries.filter(e => e.saved).length;
    const pending = entries.length - saved;
    return { total: entries.length, saved, pending };
  });

  ngOnInit(): void {
    const testId = this.route.snapshot.paramMap.get('id');
    // Detect mode based on route
    const url = this.route.snapshot.url.map(s => s.path).join('/');
    this.isRegistrationsMode.set(url.includes('registrations'));

    if (testId) {
      this.loadTestData(testId);
    }
  }

  loadTestData(testId: string): void {
    this.loading.set(true);
    this.error.set(null);

    // Load test details
    this.testService.getTest(testId).subscribe({
      next: test => {
        this.test.set(test);
        // Load registrations
        this.testService.getRegistrations(testId).subscribe({
          next: registrations => {
            this.registrations.set(registrations);
            this.initializeResultEntries(test, registrations);
            this.loading.set(false);
          },
          error: err => {
            this.error.set('Failed to load registrations');
            this.loading.set(false);
            console.error('Failed to load registrations:', err);
          },
        });
      },
      error: err => {
        this.error.set('Failed to load test details');
        this.loading.set(false);
        console.error('Failed to load test:', err);
      },
    });
  }

  initializeResultEntries(test: EntranceTest, registrations: TestRegistration[]): void {
    const entries: ResultEntry[] = registrations.map(reg => {
      const marks: { [subject: string]: number } = {};

      // Initialize marks for each subject
      test.subjects.forEach(subject => {
        // Check if registration already has marks
        const existingMark = reg.marks.find(m => m.subjectName === subject.name);
        marks[subject.name] = existingMark ? existingMark.obtainedMarks : 0;
      });

      return {
        registrationId: reg.id,
        rollNumber: reg.rollNumber || '',
        studentName: reg.application?.studentName || 'Unknown',
        marks,
        remarks: reg.remarks || '',
        saved: reg.result !== undefined && reg.result !== null,
      };
    });

    this.resultEntries.set(entries);
  }

  updateMark(entry: ResultEntry, subjectName: string, value: number): void {
    this.resultEntries.update(entries =>
      entries.map(e => {
        if (e.registrationId === entry.registrationId) {
          return {
            ...e,
            marks: { ...e.marks, [subjectName]: value },
            saved: false,
          };
        }
        return e;
      })
    );
  }

  updateRemarks(entry: ResultEntry, remarks: string): void {
    this.resultEntries.update(entries =>
      entries.map(e => {
        if (e.registrationId === entry.registrationId) {
          return { ...e, remarks, saved: false };
        }
        return e;
      })
    );
  }

  getTotal(entry: ResultEntry): number {
    return Object.values(entry.marks).reduce((sum, m) => sum + (m || 0), 0);
  }

  getMaxTotal(): number {
    const test = this.test();
    if (!test) return 0;
    return test.subjects.reduce((sum, s) => sum + s.maxMarks, 0);
  }

  getPercentage(entry: ResultEntry): number {
    const total = this.getTotal(entry);
    const max = this.getMaxTotal();
    if (max === 0) return 0;
    return Math.round((total / max) * 100 * 100) / 100;
  }

  getResult(entry: ResultEntry): { verdict: string; class: string } {
    const percentage = this.getPercentage(entry);
    if (percentage >= 75) {
      return { verdict: 'Distinction', class: 'result--distinction' };
    } else if (percentage >= 60) {
      return { verdict: 'Merit', class: 'result--merit' };
    } else if (percentage >= 40) {
      return { verdict: 'Pass', class: 'result--pass' };
    } else {
      return { verdict: 'Fail', class: 'result--fail' };
    }
  }

  saveResult(entry: ResultEntry): void {
    this.saving.set(true);
    const test = this.test();
    if (!test) return;

    const marks: SubjectMarks[] = test.subjects.map(subject => ({
      subjectName: subject.name,
      maxMarks: subject.maxMarks,
      obtainedMarks: entry.marks[subject.name] || 0,
    }));

    this.testService
      .submitResult(test.id, {
        registrationId: entry.registrationId,
        marks,
        remarks: entry.remarks,
      })
      .subscribe({
        next: () => {
          this.resultEntries.update(entries =>
            entries.map(e => {
              if (e.registrationId === entry.registrationId) {
                return { ...e, saved: true };
              }
              return e;
            })
          );
          this.toastService.success('Result saved successfully');
          this.saving.set(false);
        },
        error: err => {
          this.toastService.error('Failed to save result');
          this.saving.set(false);
          console.error('Failed to save result:', err);
        },
      });
  }

  saveAllResults(): void {
    const unsavedEntries = this.resultEntries().filter(e => !e.saved);
    if (unsavedEntries.length === 0) {
      this.toastService.info('All results are already saved');
      return;
    }

    this.saving.set(true);
    const test = this.test();
    if (!test) return;

    const results = unsavedEntries.map(entry => ({
      registrationId: entry.registrationId,
      marks: test.subjects.map(subject => ({
        subjectName: subject.name,
        maxMarks: subject.maxMarks,
        obtainedMarks: entry.marks[subject.name] || 0,
      })),
      remarks: entry.remarks,
    }));

    this.testService.submitBulkResults(test.id, { results }).subscribe({
      next: () => {
        this.resultEntries.update(entries =>
          entries.map(e => ({ ...e, saved: true }))
        );
        this.toastService.success(`${unsavedEntries.length} results saved successfully`);
        this.saving.set(false);
      },
      error: err => {
        this.toastService.error('Failed to save results');
        this.saving.set(false);
        console.error('Failed to save bulk results:', err);
      },
    });
  }

  goBack(): void {
    this.router.navigate(['/admissions/tests']);
  }

  // Registration methods
  openRegisterModal(): void {
    this.showRegisterModal.set(true);
    this.selectedApplicationId.set('');
    this.loadAvailableApplications();
  }

  closeRegisterModal(): void {
    this.showRegisterModal.set(false);
    this.selectedApplicationId.set('');
  }

  loadAvailableApplications(): void {
    const test = this.test();
    if (!test) return;

    this.loadingApplications.set(true);

    // Get applications with documents_verified status for the class that this test is for
    this.applicationService.getApplications({
      stage: 'documents_verified',
    }).subscribe({
      next: applications => {
        // Filter by class name and exclude already registered
        const registeredAppIds = new Set(this.registrations().map(r => r.applicationId));
        const filtered = applications.filter(app =>
          test.classNames.includes(app.classApplying) &&
          !registeredAppIds.has(app.id)
        );
        this.availableApplications.set(filtered);
        this.loadingApplications.set(false);
      },
      error: err => {
        console.error('Failed to load applications:', err);
        this.toastService.error('Failed to load eligible applications');
        this.loadingApplications.set(false);
      },
    });
  }

  registerCandidate(): void {
    const test = this.test();
    const appId = this.selectedApplicationId();
    if (!test || !appId) return;

    this.registering.set(true);

    this.testService.registerCandidate(test.id, { applicationId: appId }).subscribe({
      next: registration => {
        this.toastService.success('Candidate registered successfully');
        // Reload registrations to get updated list
        this.testService.getRegistrations(test.id).subscribe({
          next: registrations => {
            this.registrations.set(registrations);
            this.initializeResultEntries(test, registrations);
          },
        });
        this.closeRegisterModal();
        this.registering.set(false);
      },
      error: err => {
        console.error('Failed to register candidate:', err);
        this.toastService.error('Failed to register candidate');
        this.registering.set(false);
      },
    });
  }

  // Helper to get full name from application
  getFullName(app: AdmissionApplication): string {
    return [app.firstName, app.middleName, app.lastName]
      .filter(Boolean)
      .join(' ');
  }

  // Download hall tickets
  downloadHallTickets(): void {
    const test = this.test();
    if (!test) return;

    this.testService.generateHallTickets(test.id).subscribe({
      next: blob => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `hall-tickets-${test.testName.replace(/\s+/g, '-')}.pdf`;
        link.click();
        window.URL.revokeObjectURL(url);
        this.toastService.success('Hall tickets downloaded');
      },
      error: err => {
        console.error('Failed to download hall tickets:', err);
        this.toastService.error('Failed to download hall tickets');
      },
    });
  }
}
