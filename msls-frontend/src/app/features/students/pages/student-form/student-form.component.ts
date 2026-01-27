/**
 * Student Form Page Component
 *
 * Handles both creation and editing of student profiles.
 */

import {
  Component,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  inject,
  OnInit,
  OnDestroy,
  signal,
  computed,
  effect,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
  AbstractControl,
} from '@angular/forms';
import { Router, ActivatedRoute, RouterLink } from '@angular/router';
import { Subscription } from 'rxjs';

import {
  MslsInputComponent,
  MslsSelectComponent,
  MslsFormFieldComponent,
  SelectOption,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services';
import { StudentService } from '../../services/student.service';
import { BranchService } from '../../../admin/branches/branch.service';
import {
  Student,
  Gender,
  CreateStudentRequest,
  UpdateStudentRequest,
} from '../../models/student.model';
import { AddressFormComponent } from '../../components/address-form/address-form.component';
import { PhotoUploadComponent } from '../../components/photo-upload/photo-upload.component';

@Component({
  selector: 'msls-student-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterLink,
    MslsInputComponent,
    MslsSelectComponent,
    MslsFormFieldComponent,
    AddressFormComponent,
    PhotoUploadComponent,
  ],
  templateUrl: './student-form.component.html',
  styleUrl: './student-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StudentFormComponent implements OnInit, OnDestroy {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private studentService = inject(StudentService);
  private branchService = inject(BranchService);
  private toast = inject(ToastService);
  private cdr = inject(ChangeDetectorRef);

  private formSubscription?: Subscription;

  // =========================================================================
  // State
  // =========================================================================

  readonly isEditMode = signal<boolean>(false);
  readonly studentId = signal<string | null>(null);
  readonly loading = signal<boolean>(false);
  readonly submitting = signal<boolean>(false);
  readonly photoFile = signal<File | null>(null);
  readonly previewPhoto = signal<string | null>(null);
  readonly nextAdmissionNumber = signal<string>('');

  readonly student = this.studentService.selectedStudent;

  // =========================================================================
  // Multi-Step Wizard State
  // =========================================================================

  readonly currentStep = signal<number>(1);
  readonly totalSteps = 4;

  readonly steps = [
    { number: 1, label: 'Photo', icon: 'fa-camera' },
    { number: 2, label: 'Personal', icon: 'fa-user' },
    { number: 3, label: 'Academic', icon: 'fa-graduation-cap' },
    { number: 4, label: 'Address', icon: 'fa-location-dot' },
  ];

  // =========================================================================
  // Form
  // =========================================================================

  readonly form: FormGroup = this.fb.group({
    // Personal Information
    firstName: ['', [Validators.required, Validators.maxLength(100)]],
    middleName: ['', [Validators.maxLength(100)]],
    lastName: ['', [Validators.required, Validators.maxLength(100)]],
    dateOfBirth: ['', [Validators.required, this.dateNotInFutureValidator]],
    gender: ['', [Validators.required]],
    bloodGroup: [''],
    aadhaarNumber: ['', [Validators.pattern(/^\d{12}$/)]],

    // Academic Information
    branchId: ['', [Validators.required]],
    admissionDate: [''],

    // Address Information
    currentAddress: this.fb.group({
      addressLine1: ['', [Validators.required, Validators.maxLength(255)]],
      addressLine2: ['', [Validators.maxLength(255)]],
      city: ['', [Validators.required, Validators.maxLength(100)]],
      state: ['', [Validators.required, Validators.maxLength(100)]],
      postalCode: ['', [Validators.required, Validators.pattern(/^\d{6}$/)]],
      country: ['India', [Validators.required, Validators.maxLength(100)]],
    }),

    sameAsCurrentAddress: [false],

    permanentAddress: this.fb.group({
      addressLine1: ['', [Validators.maxLength(255)]],
      addressLine2: ['', [Validators.maxLength(255)]],
      city: ['', [Validators.maxLength(100)]],
      state: ['', [Validators.maxLength(100)]],
      postalCode: ['', [Validators.pattern(/^\d{6}$/)]],
      country: ['India', [Validators.maxLength(100)]],
    }),

    // For optimistic locking on edit
    version: [0],
  });

  // =========================================================================
  // Options
  // =========================================================================

  readonly genderOptions: SelectOption[] = [
    { value: 'male', label: 'Male' },
    { value: 'female', label: 'Female' },
    { value: 'other', label: 'Other' },
  ];

  readonly bloodGroupOptions: SelectOption[] = [
    { value: '', label: 'Select Blood Group' },
    { value: 'A+', label: 'A+' },
    { value: 'A-', label: 'A-' },
    { value: 'B+', label: 'B+' },
    { value: 'B-', label: 'B-' },
    { value: 'O+', label: 'O+' },
    { value: 'O-', label: 'O-' },
    { value: 'AB+', label: 'AB+' },
    { value: 'AB-', label: 'AB-' },
  ];

  // Branch options loaded from service
  readonly branchOptions = signal<SelectOption[]>([
    { value: '', label: 'Select Branch' },
  ]);

  // =========================================================================
  // Computed
  // =========================================================================

  readonly pageTitle = computed(() =>
    this.isEditMode() ? 'Edit Student' : 'Add New Student'
  );

  readonly submitButtonText = computed(() =>
    this.isEditMode() ? 'Update Student' : 'Create Student'
  );

  // Note: canSubmit is a method (not computed) because form.valid doesn't trigger signal updates
  canSubmit(): boolean {
    return this.form.valid && !this.submitting();
  }

  // Step validation methods (not computed, because form controls are not signals)
  isStep1Valid(): boolean {
    return true; // Photo is optional
  }

  isStep2Valid(): boolean {
    const firstName = this.form.get('firstName');
    const lastName = this.form.get('lastName');
    const dateOfBirth = this.form.get('dateOfBirth');
    const gender = this.form.get('gender');
    return !!(firstName?.valid && lastName?.valid && dateOfBirth?.valid && gender?.valid);
  }

  isStep3Valid(): boolean {
    const branchId = this.form.get('branchId');
    return !!branchId?.valid;
  }

  isStep4Valid(): boolean {
    const currentAddress = this.form.get('currentAddress') as FormGroup;
    return !!currentAddress?.valid;
  }

  isCurrentStepValid(): boolean {
    switch (this.currentStep()) {
      case 1: return this.isStep1Valid();
      case 2: return this.isStep2Valid();
      case 3: return this.isStep3Valid();
      case 4: return this.isStep4Valid();
      default: return false;
    }
  }

  canGoNext(): boolean {
    return this.currentStep() < this.totalSteps && this.isCurrentStepValid();
  }

  canGoPrevious(): boolean {
    return this.currentStep() > 1;
  }

  isLastStep(): boolean {
    return this.currentStep() === this.totalSteps;
  }

  // =========================================================================
  // Constructor Effect
  // =========================================================================

  constructor() {
    // Watch sameAsCurrentAddress toggle
    effect(() => {
      const sameAs = this.form.get('sameAsCurrentAddress')?.value;
      if (sameAs) {
        this.copyCurrentToPermanent();
      }
    }, { allowSignalWrites: true });
  }

  // =========================================================================
  // Lifecycle
  // =========================================================================

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');

    if (id && id !== 'new') {
      this.isEditMode.set(true);
      this.studentId.set(id);
      this.loadStudent(id);
    } else {
      // For new student, get next admission number preview
      this.loadNextAdmissionNumber();
    }

    // Set default admission date to today
    if (!this.isEditMode()) {
      this.form.patchValue({
        admissionDate: new Date().toISOString().split('T')[0],
      });
    }

    // Load branches for the dropdown
    this.loadBranches();

    // Subscribe to form value changes to trigger change detection for OnPush
    this.formSubscription = this.form.valueChanges.subscribe(() => {
      this.cdr.markForCheck();
    });
  }

  ngOnDestroy(): void {
    this.formSubscription?.unsubscribe();
  }

  // =========================================================================
  // Data Loading
  // =========================================================================

  private loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: (branches) => {
        const options: SelectOption[] = [
          { value: '', label: 'Select Branch' },
          ...branches
            .filter(b => b.isActive)
            .map(b => ({ value: b.id, label: b.name })),
        ];
        this.branchOptions.set(options);
      },
      error: (err) => {
        this.toast.error('Failed to load branches', err.message);
      },
    });
  }

  private loadStudent(id: string): void {
    this.loading.set(true);
    this.studentService.getStudent(id).subscribe({
      next: (student) => {
        this.populateForm(student);
        this.loading.set(false);
      },
      error: (err) => {
        this.toast.error('Failed to load student', err.message);
        this.loading.set(false);
        this.router.navigate(['/students']);
      },
    });
  }

  private loadNextAdmissionNumber(): void {
    const branchId = this.form.get('branchId')?.value;
    if (branchId) {
      this.studentService.getNextAdmissionNumber(branchId).subscribe({
        next: (result) => {
          this.nextAdmissionNumber.set(result.admissionNumber);
        },
        error: () => {
          // Silent fail - admission number will be generated on submit
        },
      });
    }
  }

  private populateForm(student: Student): void {
    this.form.patchValue({
      firstName: student.firstName,
      middleName: student.middleName || '',
      lastName: student.lastName,
      dateOfBirth: student.dateOfBirth,
      gender: student.gender,
      bloodGroup: student.bloodGroup || '',
      aadhaarNumber: student.aadhaarNumber || '',
      branchId: student.branchId,
      version: student.version,
    });

    if (student.currentAddress) {
      this.form.get('currentAddress')?.patchValue({
        addressLine1: student.currentAddress.addressLine1,
        addressLine2: student.currentAddress.addressLine2 || '',
        city: student.currentAddress.city,
        state: student.currentAddress.state,
        postalCode: student.currentAddress.postalCode,
        country: student.currentAddress.country,
      });
    }

    if (student.permanentAddress) {
      this.form.get('permanentAddress')?.patchValue({
        addressLine1: student.permanentAddress.addressLine1,
        addressLine2: student.permanentAddress.addressLine2 || '',
        city: student.permanentAddress.city,
        state: student.permanentAddress.state,
        postalCode: student.permanentAddress.postalCode,
        country: student.permanentAddress.country,
      });
    }

    if (student.photoUrl) {
      this.previewPhoto.set(student.photoUrl);
    }
  }

  // =========================================================================
  // Form Submission
  // =========================================================================

  onSubmit(): void {
    if (!this.form.valid) {
      this.markAllAsTouched();
      return;
    }

    this.submitting.set(true);
    const formValue = this.form.getRawValue();

    if (this.isEditMode()) {
      this.updateStudent(formValue);
    } else {
      this.createStudent(formValue);
    }
  }

  private createStudent(formValue: Record<string, unknown>): void {
    const request: CreateStudentRequest = {
      branchId: formValue['branchId'] as string,
      firstName: formValue['firstName'] as string,
      middleName: formValue['middleName'] as string || undefined,
      lastName: formValue['lastName'] as string,
      dateOfBirth: formValue['dateOfBirth'] as string,
      gender: formValue['gender'] as Gender,
      bloodGroup: formValue['bloodGroup'] as string || undefined,
      aadhaarNumber: formValue['aadhaarNumber'] as string || undefined,
      admissionDate: formValue['admissionDate'] as string || undefined,
      currentAddress: this.extractAddress(formValue['currentAddress'] as Record<string, string>),
      permanentAddress: formValue['sameAsCurrentAddress']
        ? undefined
        : this.extractAddress(formValue['permanentAddress'] as Record<string, string>),
      sameAsCurrentAddress: formValue['sameAsCurrentAddress'] as boolean,
    };

    this.studentService.createStudent(request).subscribe({
      next: (student) => {
        this.handlePhotoUpload(student.id);
        this.toast.success('Student created successfully');
        this.submitting.set(false);
        this.router.navigate(['/students', student.id]);
      },
      error: (err) => {
        this.toast.error('Failed to create student', err.message);
        this.submitting.set(false);
      },
    });
  }

  private updateStudent(formValue: Record<string, unknown>): void {
    const id = this.studentId()!;
    const request: UpdateStudentRequest = {
      firstName: formValue['firstName'] as string,
      middleName: formValue['middleName'] as string || undefined,
      lastName: formValue['lastName'] as string,
      dateOfBirth: formValue['dateOfBirth'] as string,
      gender: formValue['gender'] as Gender,
      bloodGroup: formValue['bloodGroup'] as string || undefined,
      aadhaarNumber: formValue['aadhaarNumber'] as string || undefined,
      currentAddress: this.extractAddress(formValue['currentAddress'] as Record<string, string>),
      permanentAddress: formValue['sameAsCurrentAddress']
        ? undefined
        : this.extractAddress(formValue['permanentAddress'] as Record<string, string>),
      sameAsCurrentAddress: formValue['sameAsCurrentAddress'] as boolean,
      version: formValue['version'] as number,
    };

    this.studentService.updateStudent(id, request).subscribe({
      next: (student) => {
        this.handlePhotoUpload(student.id);
        this.toast.success('Student updated successfully');
        this.submitting.set(false);
        this.router.navigate(['/students', student.id]);
      },
      error: (err) => {
        if (err.message?.includes('conflict') || err.message?.includes('version')) {
          this.toast.error('Conflict detected: Someone else has modified this record. Please refresh and try again.');
        } else {
          this.toast.error(`Failed to update student: ${err.message}`);
        }
        this.submitting.set(false);
      },
    });
  }

  private extractAddress(addr: Record<string, string>): { addressLine1: string; addressLine2?: string; city: string; state: string; postalCode: string; country: string } | undefined {
    if (!addr || !addr['addressLine1']) {
      return undefined;
    }
    return {
      addressLine1: addr['addressLine1'],
      addressLine2: addr['addressLine2'] || undefined,
      city: addr['city'],
      state: addr['state'],
      postalCode: addr['postalCode'],
      country: addr['country'] || 'India',
    };
  }

  private handlePhotoUpload(studentId: string): void {
    const file = this.photoFile();
    if (file) {
      this.studentService.uploadPhoto(studentId, file).subscribe({
        error: (err) => {
          this.toast.warning('Photo upload failed', err.message);
        },
      });
    }
  }

  // =========================================================================
  // Photo Handling
  // =========================================================================

  onPhotoSelected(file: File): void {
    this.photoFile.set(file);

    // Create preview URL
    const reader = new FileReader();
    reader.onload = () => {
      this.previewPhoto.set(reader.result as string);
    };
    reader.readAsDataURL(file);
  }

  onPhotoRemoved(): void {
    this.photoFile.set(null);
    if (!this.isEditMode()) {
      this.previewPhoto.set(null);
    }
  }

  // =========================================================================
  // Address Helpers
  // =========================================================================

  private copyCurrentToPermanent(): void {
    const currentAddress = this.form.get('currentAddress')?.value;
    this.form.get('permanentAddress')?.patchValue(currentAddress);
  }

  onSameAsCurrentChange(checked: boolean): void {
    this.form.patchValue({ sameAsCurrentAddress: checked });
    if (checked) {
      this.copyCurrentToPermanent();
      this.form.get('permanentAddress')?.disable();
    } else {
      this.form.get('permanentAddress')?.enable();
    }
  }

  // =========================================================================
  // Validation Helpers
  // =========================================================================

  private dateNotInFutureValidator(control: AbstractControl): { [key: string]: boolean } | null {
    if (!control.value) return null;
    const date = new Date(control.value);
    if (date > new Date()) {
      return { futureDate: true };
    }
    return null;
  }

  private markAllAsTouched(): void {
    Object.keys(this.form.controls).forEach((key) => {
      const control = this.form.get(key);
      control?.markAsTouched();
      if (control instanceof FormGroup) {
        Object.keys(control.controls).forEach((innerKey) => {
          control.get(innerKey)?.markAsTouched();
        });
      }
    });
  }

  getError(controlName: string): string {
    const control = this.form.get(controlName);
    if (!control || !control.touched || !control.errors) return '';

    if (control.errors['required']) return 'This field is required';
    if (control.errors['maxlength']) return `Maximum ${control.errors['maxlength'].requiredLength} characters allowed`;
    if (control.errors['pattern']) {
      if (controlName === 'aadhaarNumber') return 'Must be 12 digits';
      if (controlName.includes('postalCode')) return 'Must be 6 digits';
      return 'Invalid format';
    }
    if (control.errors['futureDate']) return 'Date cannot be in the future';

    return 'Invalid value';
  }

  // =========================================================================
  // Step Navigation
  // =========================================================================

  nextStep(): void {
    if (this.canGoNext()) {
      // Mark current step fields as touched for validation feedback
      this.markStepFieldsAsTouched(this.currentStep());
      if (this.isCurrentStepValid()) {
        this.currentStep.update(step => step + 1);
      }
    }
  }

  prevStep(): void {
    if (this.canGoPrevious()) {
      this.currentStep.update(step => step - 1);
    }
  }

  goToStep(step: number): void {
    // Only allow going to completed steps or current step
    if (step >= 1 && step <= this.totalSteps) {
      // Can always go back, but can only go forward if previous steps are valid
      if (step < this.currentStep()) {
        this.currentStep.set(step);
      } else if (step === this.currentStep() + 1 && this.isCurrentStepValid()) {
        this.currentStep.set(step);
      }
    }
  }

  isStepCompleted(step: number): boolean {
    if (step >= this.currentStep()) return false;
    switch (step) {
      case 1: return this.isStep1Valid();
      case 2: return this.isStep2Valid();
      case 3: return this.isStep3Valid();
      case 4: return this.isStep4Valid();
      default: return false;
    }
  }

  isStepAccessible(step: number): boolean {
    // Can always access current or previous steps
    if (step <= this.currentStep()) return true;
    // Can access next step if current is valid
    if (step === this.currentStep() + 1 && this.isCurrentStepValid()) return true;
    return false;
  }

  private markStepFieldsAsTouched(step: number): void {
    const stepFields: { [key: number]: string[] } = {
      1: [], // Photo step has no required fields
      2: ['firstName', 'lastName', 'dateOfBirth', 'gender'],
      3: ['branchId'],
      4: ['currentAddress.addressLine1', 'currentAddress.city', 'currentAddress.state', 'currentAddress.postalCode'],
    };

    const fields = stepFields[step] || [];
    fields.forEach(fieldName => {
      const control = this.form.get(fieldName);
      control?.markAsTouched();
    });
  }

  // =========================================================================
  // Page Navigation
  // =========================================================================

  onCancel(): void {
    if (this.isEditMode() && this.studentId()) {
      this.router.navigate(['/students', this.studentId()]);
    } else {
      this.router.navigate(['/students']);
    }
  }
}
