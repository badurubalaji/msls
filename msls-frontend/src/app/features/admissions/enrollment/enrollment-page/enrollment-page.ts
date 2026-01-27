/**
 * MSLS Enrollment Page Component
 *
 * Displays students with accepted offers pending enrollment completion.
 * Allows completing enrollment with section assignment and roll number.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';

import {
  MslsButtonComponent,
  MslsInputComponent,
  MslsFormFieldComponent,
  MslsSelectComponent,
  MslsBadgeComponent,
  MslsCardComponent,
  SelectOption,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services/toast.service';
import { MeritService } from '../../merit/merit.service';
import {
  MeritListEntry,
  EnrollmentRequest,
  EnrollmentResponse,
  getStatusConfig,
  APPLICATION_STATUS_CONFIG,
} from '../../merit/merit.model';

interface EnrollmentCandidate {
  id: string;
  applicationId: string;
  studentName: string;
  parentName: string;
  parentPhone: string;
  className: string;
  sectionAssigned?: string;
  offerAcceptedAt?: string;
  rank?: number;
}

@Component({
  selector: 'app-enrollment-page',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsButtonComponent,
    MslsInputComponent,
    MslsFormFieldComponent,
    MslsSelectComponent,
    MslsBadgeComponent,
    MslsCardComponent,
  ],
  templateUrl: './enrollment-page.html',
  styleUrl: './enrollment-page.scss',
})
export class EnrollmentPage implements OnInit {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private meritService = inject(MeritService);
  private toastService = inject(ToastService);

  // State signals
  loading = signal(false);
  enrolling = signal(false);
  candidates = signal<EnrollmentCandidate[]>([]);
  selectedCandidate = signal<EnrollmentCandidate | null>(null);
  showEnrollDialog = signal(false);
  sectionOptions = signal<SelectOption[]>([]);

  // Computed values
  totalCandidates = computed(() => this.candidates().length);
  enrolledCount = signal(0);

  // Enrollment form
  enrollForm!: FormGroup;

  // Table columns
  tableColumns = [
    { key: 'studentName', label: 'Student Name' },
    { key: 'parentName', label: 'Parent/Guardian' },
    { key: 'parentPhone', label: 'Phone' },
    { key: 'className', label: 'Class' },
    { key: 'sectionAssigned', label: 'Section' },
    { key: 'offerAcceptedAt', label: 'Offer Accepted' },
    { key: 'actions', label: 'Actions', align: 'right' as const },
  ];

  ngOnInit(): void {
    this.initForm();
    this.loadPendingEnrollments();
  }

  private initForm(): void {
    this.enrollForm = this.fb.group({
      sectionId: [''],
      rollNumber: [''],
      admissionDate: [new Date().toISOString().split('T')[0], [Validators.required]],
      remarks: [''],
    });
  }

  private loadPendingEnrollments(): void {
    this.loading.set(true);

    // For now, we use mock data. In production, this would be an API call
    // to get applications with status 'offer_accepted'
    setTimeout(() => {
      const mockCandidates: EnrollmentCandidate[] = [
        {
          id: '1',
          applicationId: 'app-1',
          studentName: 'Aarav Sharma',
          parentName: 'Rajesh Sharma',
          parentPhone: '9876543210',
          className: 'Class 1',
          sectionAssigned: 'Section A',
          offerAcceptedAt: '2026-01-20T10:30:00Z',
          rank: 1,
        },
        {
          id: '2',
          applicationId: 'app-2',
          studentName: 'Priya Patel',
          parentName: 'Vikram Patel',
          parentPhone: '9876543211',
          className: 'Class 1',
          sectionAssigned: 'Section A',
          offerAcceptedAt: '2026-01-21T14:15:00Z',
          rank: 2,
        },
        {
          id: '3',
          applicationId: 'app-3',
          studentName: 'Arjun Singh',
          parentName: 'Manpreet Singh',
          parentPhone: '9876543212',
          className: 'Class 1',
          sectionAssigned: 'Section B',
          offerAcceptedAt: '2026-01-22T09:45:00Z',
          rank: 3,
        },
      ];

      this.candidates.set(mockCandidates);
      this.loading.set(false);
    }, 500);
  }

  private loadSections(className: string): void {
    this.meritService.getSections(className).subscribe({
      next: (sections) => {
        this.sectionOptions.set(
          sections.map(s => ({
            value: s.id,
            label: s.name,
          }))
        );
      },
      error: (err) => {
        console.error('Failed to load sections:', err);
      },
    });
  }

  openEnrollDialog(candidate: EnrollmentCandidate): void {
    this.selectedCandidate.set(candidate);
    this.loadSections(candidate.className);

    // Pre-fill the section if already assigned
    this.enrollForm.patchValue({
      sectionId: candidate.sectionAssigned || '',
      admissionDate: new Date().toISOString().split('T')[0],
      rollNumber: '',
      remarks: '',
    });

    this.showEnrollDialog.set(true);
  }

  closeEnrollDialog(): void {
    this.showEnrollDialog.set(false);
    this.selectedCandidate.set(null);
    this.enrollForm.reset({
      admissionDate: new Date().toISOString().split('T')[0],
    });
  }

  onEnroll(): void {
    const candidate = this.selectedCandidate();
    if (!candidate || !this.enrollForm.valid) return;

    this.enrolling.set(true);
    const formValue = this.enrollForm.value;

    const request: EnrollmentRequest = {
      sectionId: formValue.sectionId || undefined,
      rollNumber: formValue.rollNumber || undefined,
      admissionDate: formValue.admissionDate,
      remarks: formValue.remarks || undefined,
    };

    this.meritService.enroll(candidate.applicationId, request).subscribe({
      next: (response: EnrollmentResponse) => {
        this.enrolling.set(false);
        this.toastService.success(
          `${candidate.studentName} has been enrolled! Enrollment #: ${response.enrollmentNumber || 'N/A'}`
        );

        // Remove from candidates list
        this.candidates.update(list =>
          list.filter(c => c.applicationId !== candidate.applicationId)
        );
        this.enrolledCount.update(c => c + 1);

        this.closeEnrollDialog();
      },
      error: (err) => {
        this.enrolling.set(false);
        console.error('Enrollment failed:', err);
        this.toastService.error('Failed to complete enrollment. Please try again.');
      },
    });
  }

  formatDate(dateString: string): string {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }

  goBack(): void {
    this.router.navigate(['../merit-list'], { relativeTo: this.route });
  }

  getStatusConfig(status: string) {
    return getStatusConfig(status as keyof typeof APPLICATION_STATUS_CONFIG);
  }
}
