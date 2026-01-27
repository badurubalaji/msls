/**
 * Staff Form Page Component
 *
 * Handles both creation and editing of staff profiles.
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
} from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  Validators,
  ReactiveFormsModule,
  AbstractControl,
} from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';

import {
  MslsInputComponent,
  MslsSelectComponent,
  MslsFormFieldComponent,
  SelectOption,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services';
import { StaffService } from '../../services/staff.service';
import { BranchService } from '../../../admin/branches/branch.service';
import { DepartmentService } from '../../../admin/departments/department.service';
import { DesignationService } from '../../../admin/designations/designation.service';
import {
  Staff,
  Gender,
  StaffType,
  CreateStaffRequest,
  UpdateStaffRequest,
} from '../../models/staff.model';

@Component({
  selector: 'msls-staff-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MslsInputComponent,
    MslsSelectComponent,
    MslsFormFieldComponent,
  ],
  templateUrl: './staff-form.component.html',
  styleUrl: './staff-form.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StaffFormComponent implements OnInit, OnDestroy {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private staffService = inject(StaffService);
  private branchService = inject(BranchService);
  private departmentService = inject(DepartmentService);
  private designationService = inject(DesignationService);
  private toast = inject(ToastService);
  private cdr = inject(ChangeDetectorRef);

  private formSubscription?: Subscription;

  // State
  readonly isEditMode = signal<boolean>(false);
  readonly staffId = signal<string | null>(null);
  readonly loading = signal<boolean>(false);
  readonly submitting = signal<boolean>(false);
  readonly photoFile = signal<File | null>(null);
  readonly previewPhoto = signal<string | null>(null);
  readonly nextEmployeeId = signal<string>('');

  readonly staff = this.staffService.selectedStaff;

  // Multi-Step Wizard State
  readonly currentStep = signal<number>(1);
  readonly totalSteps = 4;

  readonly steps = [
    { number: 1, label: 'Personal', icon: 'fa-user' },
    { number: 2, label: 'Contact', icon: 'fa-address-book' },
    { number: 3, label: 'Employment', icon: 'fa-briefcase' },
    { number: 4, label: 'Address', icon: 'fa-location-dot' },
  ];

  // Form
  readonly form: FormGroup = this.fb.group({
    // Personal Information
    firstName: ['', [Validators.required, Validators.maxLength(100)]],
    middleName: ['', [Validators.maxLength(100)]],
    lastName: ['', [Validators.required, Validators.maxLength(100)]],
    dateOfBirth: ['', [Validators.required, this.dateNotInFutureValidator]],
    gender: ['', [Validators.required]],
    bloodGroup: [''],
    nationality: ['Indian'],
    religion: [''],
    maritalStatus: [''],

    // Contact Information
    workEmail: ['', [Validators.required, Validators.email]],
    workPhone: ['', [Validators.required, Validators.pattern(/^\d{10}$/)]],
    personalEmail: ['', [Validators.email]],
    personalPhone: ['', [Validators.pattern(/^\d{10}$/)]],
    emergencyContactName: [''],
    emergencyContactPhone: ['', [Validators.pattern(/^\d{10}$/)]],
    emergencyContactRelation: [''],

    // Employment Information
    branchId: ['', [Validators.required]],
    staffType: ['teaching', [Validators.required]],
    departmentId: [''],
    designationId: [''],
    reportingManagerId: [''],
    joinDate: ['', [Validators.required]],
    confirmationDate: [''],
    probationEndDate: [''],
    bio: ['', [Validators.maxLength(1000)]],

    // Address Information
    currentAddress: this.fb.group({
      addressLine1: ['', [Validators.maxLength(255)]],
      addressLine2: ['', [Validators.maxLength(255)]],
      city: ['', [Validators.maxLength(100)]],
      state: ['', [Validators.maxLength(100)]],
      pincode: ['', [Validators.pattern(/^\d{6}$/)]],
      country: ['India', [Validators.maxLength(100)]],
    }),

    sameAsCurrent: [false],

    permanentAddress: this.fb.group({
      addressLine1: ['', [Validators.maxLength(255)]],
      addressLine2: ['', [Validators.maxLength(255)]],
      city: ['', [Validators.maxLength(100)]],
      state: ['', [Validators.maxLength(100)]],
      pincode: ['', [Validators.pattern(/^\d{6}$/)]],
      country: ['India', [Validators.maxLength(100)]],
    }),

    // For optimistic locking on edit
    version: [0],
  });

  // Options
  readonly genderOptions: SelectOption[] = [
    { value: 'male', label: 'Male' },
    { value: 'female', label: 'Female' },
    { value: 'other', label: 'Other' },
  ];

  readonly staffTypeOptions: SelectOption[] = [
    { value: 'teaching', label: 'Teaching' },
    { value: 'non_teaching', label: 'Non-Teaching' },
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

  readonly maritalStatusOptions: SelectOption[] = [
    { value: '', label: 'Select Status' },
    { value: 'single', label: 'Single' },
    { value: 'married', label: 'Married' },
    { value: 'divorced', label: 'Divorced' },
    { value: 'widowed', label: 'Widowed' },
  ];

  // Dynamic options (would be loaded from services)
  readonly branchOptions = signal<SelectOption[]>([
    { value: '', label: 'Select Branch' },
  ]);

  readonly departmentOptions = signal<SelectOption[]>([
    { value: '', label: 'Select Department' },
  ]);

  readonly designationOptions = signal<SelectOption[]>([
    { value: '', label: 'Select Designation' },
  ]);

  // Computed
  readonly pageTitle = computed(() =>
    this.isEditMode() ? 'Edit Staff' : 'Add New Staff'
  );

  readonly submitButtonText = computed(() =>
    this.isEditMode() ? 'Update Staff' : 'Create Staff'
  );

  // Lifecycle
  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id && id !== 'new') {
      this.isEditMode.set(true);
      this.staffId.set(id);
      this.loadStaff(id);
    } else {
      this.loadNextEmployeeId();
      this.setDefaultJoinDate();
    }

    // Load dropdown options
    this.loadBranches();
    this.loadDepartments();
    this.loadDesignations();

    // Listen to "same as current" checkbox
    this.formSubscription = this.form.get('sameAsCurrent')?.valueChanges.subscribe((same) => {
      if (same) {
        this.copyCurrentToPermanent();
      }
    });
  }

  ngOnDestroy(): void {
    this.formSubscription?.unsubscribe();
  }

  // Data Loading
  private loadStaff(id: string): void {
    this.loading.set(true);
    this.staffService.getStaff(id).subscribe({
      next: (staff) => {
        this.populateForm(staff);
        this.loading.set(false);
        this.cdr.markForCheck();
      },
      error: () => {
        this.toast.error('Failed to load staff details');
        this.loading.set(false);
        this.router.navigate(['/staff']);
      },
    });
  }

  private loadNextEmployeeId(): void {
    this.staffService.previewEmployeeId().subscribe({
      next: (result) => {
        this.nextEmployeeId.set(result.employeeId);
        this.cdr.markForCheck();
      },
      error: () => {
        // Non-critical error, ignore
      },
    });
  }

  private loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: (branches) => {
        const options: SelectOption[] = [
          { value: '', label: 'Select Branch' },
          ...branches
            .filter((b) => b.isActive)
            .map((b) => ({
              value: b.id,
              label: b.name,
            })),
        ];
        this.branchOptions.set(options);
        this.cdr.markForCheck();
      },
      error: () => {
        this.toast.error('Failed to load branches');
      },
    });
  }

  private loadDepartments(): void {
    this.departmentService.getDepartmentsDropdown().subscribe({
      next: (departments) => {
        const options: SelectOption[] = [
          { value: '', label: 'Select Department' },
          ...departments.map((d) => ({
            value: d.id,
            label: d.name,
          })),
        ];
        this.departmentOptions.set(options);
        this.cdr.markForCheck();
      },
      error: () => {
        // Non-critical, use empty list
        this.departmentOptions.set([{ value: '', label: 'Select Department' }]);
      },
    });
  }

  private loadDesignations(): void {
    this.designationService.getDesignationsDropdown().subscribe({
      next: (designations) => {
        const options: SelectOption[] = [
          { value: '', label: 'Select Designation' },
          ...designations.map((d) => ({
            value: d.id,
            label: `${d.name} (Level ${d.level})`,
          })),
        ];
        this.designationOptions.set(options);
        this.cdr.markForCheck();
      },
      error: () => {
        // Non-critical, use empty list
        this.designationOptions.set([{ value: '', label: 'Select Designation' }]);
      },
    });
  }

  private setDefaultJoinDate(): void {
    const today = new Date().toISOString().split('T')[0];
    this.form.patchValue({ joinDate: today });
  }

  private populateForm(staff: Staff): void {
    this.form.patchValue({
      firstName: staff.firstName,
      middleName: staff.middleName || '',
      lastName: staff.lastName,
      dateOfBirth: staff.dateOfBirth,
      gender: staff.gender,
      bloodGroup: staff.bloodGroup || '',
      nationality: staff.nationality || 'Indian',
      religion: staff.religion || '',
      maritalStatus: staff.maritalStatus || '',
      workEmail: staff.workEmail,
      workPhone: staff.workPhone,
      personalEmail: staff.personalEmail || '',
      personalPhone: staff.personalPhone || '',
      emergencyContactName: staff.emergencyContactName || '',
      emergencyContactPhone: staff.emergencyContactPhone || '',
      emergencyContactRelation: staff.emergencyContactRelation || '',
      branchId: staff.branchId,
      staffType: staff.staffType,
      departmentId: staff.departmentId || '',
      designationId: staff.designationId || '',
      reportingManagerId: staff.reportingManagerId || '',
      joinDate: staff.joinDate,
      confirmationDate: staff.confirmationDate || '',
      probationEndDate: staff.probationEndDate || '',
      bio: staff.bio || '',
      sameAsCurrent: staff.sameAsCurrent,
      version: staff.version,
    });

    if (staff.currentAddress) {
      this.form.get('currentAddress')?.patchValue(staff.currentAddress);
    }
    if (staff.permanentAddress) {
      this.form.get('permanentAddress')?.patchValue(staff.permanentAddress);
    }
    if (staff.photoUrl) {
      this.previewPhoto.set(staff.photoUrl);
    }
  }

  // Step Navigation
  nextStep(): void {
    if (this.currentStep() < this.totalSteps) {
      this.currentStep.update((s) => s + 1);
    }
  }

  prevStep(): void {
    if (this.currentStep() > 1) {
      this.currentStep.update((s) => s - 1);
    }
  }

  goToStep(step: number): void {
    if (step >= 1 && step <= this.totalSteps) {
      this.currentStep.set(step);
    }
  }

  // Form Submission
  onSubmit(): void {
    if (!this.canSubmit()) {
      this.markAllAsTouched();
      this.toast.error('Please fill in all required fields');
      return;
    }

    this.submitting.set(true);

    if (this.isEditMode()) {
      this.updateStaff();
    } else {
      this.createStaff();
    }
  }

  private createStaff(): void {
    const data = this.buildCreateRequest();
    this.staffService.createStaff(data).subscribe({
      next: (staff) => {
        this.toast.success('Staff member created successfully');
        this.submitting.set(false);

        // Upload photo if selected
        if (this.photoFile()) {
          this.uploadPhoto(staff.id);
        } else {
          this.router.navigate(['/staff', staff.id]);
        }
      },
      error: (error) => {
        this.toast.error(error.message || 'Failed to create staff member');
        this.submitting.set(false);
      },
    });
  }

  private updateStaff(): void {
    const id = this.staffId()!;
    const data = this.buildUpdateRequest();
    this.staffService.updateStaff(id, data).subscribe({
      next: (staff) => {
        this.toast.success('Staff member updated successfully');
        this.submitting.set(false);

        // Upload photo if changed
        if (this.photoFile()) {
          this.uploadPhoto(staff.id);
        } else {
          this.router.navigate(['/staff', staff.id]);
        }
      },
      error: (error) => {
        this.toast.error(error.message || 'Failed to update staff member');
        this.submitting.set(false);
      },
    });
  }

  private uploadPhoto(staffId: string): void {
    const file = this.photoFile();
    if (!file) {
      this.router.navigate(['/staff', staffId]);
      return;
    }

    this.staffService.uploadPhoto(staffId, file).subscribe({
      next: () => {
        this.router.navigate(['/staff', staffId]);
      },
      error: () => {
        this.toast.warning('Staff saved but photo upload failed');
        this.router.navigate(['/staff', staffId]);
      },
    });
  }

  private buildCreateRequest(): CreateStaffRequest {
    const v = this.form.value;
    return {
      branchId: v.branchId,
      firstName: v.firstName,
      middleName: v.middleName || undefined,
      lastName: v.lastName,
      dateOfBirth: v.dateOfBirth,
      gender: v.gender as Gender,
      bloodGroup: v.bloodGroup || undefined,
      nationality: v.nationality || undefined,
      religion: v.religion || undefined,
      maritalStatus: v.maritalStatus || undefined,
      workEmail: v.workEmail,
      workPhone: v.workPhone,
      personalEmail: v.personalEmail || undefined,
      personalPhone: v.personalPhone || undefined,
      emergencyContactName: v.emergencyContactName || undefined,
      emergencyContactPhone: v.emergencyContactPhone || undefined,
      emergencyContactRelation: v.emergencyContactRelation || undefined,
      staffType: v.staffType as StaffType,
      departmentId: v.departmentId || undefined,
      designationId: v.designationId || undefined,
      reportingManagerId: v.reportingManagerId || undefined,
      joinDate: v.joinDate,
      confirmationDate: v.confirmationDate || undefined,
      probationEndDate: v.probationEndDate || undefined,
      bio: v.bio || undefined,
      currentAddress: this.hasAddress(v.currentAddress) ? v.currentAddress : undefined,
      permanentAddress: this.hasAddress(v.permanentAddress) ? v.permanentAddress : undefined,
      sameAsCurrent: v.sameAsCurrent,
    };
  }

  private buildUpdateRequest(): UpdateStaffRequest {
    const v = this.form.value;
    return {
      firstName: v.firstName,
      middleName: v.middleName || undefined,
      lastName: v.lastName,
      dateOfBirth: v.dateOfBirth,
      gender: v.gender as Gender,
      bloodGroup: v.bloodGroup || undefined,
      nationality: v.nationality || undefined,
      religion: v.religion || undefined,
      maritalStatus: v.maritalStatus || undefined,
      workEmail: v.workEmail,
      workPhone: v.workPhone,
      personalEmail: v.personalEmail || undefined,
      personalPhone: v.personalPhone || undefined,
      emergencyContactName: v.emergencyContactName || undefined,
      emergencyContactPhone: v.emergencyContactPhone || undefined,
      emergencyContactRelation: v.emergencyContactRelation || undefined,
      staffType: v.staffType as StaffType,
      departmentId: v.departmentId || undefined,
      designationId: v.designationId || undefined,
      reportingManagerId: v.reportingManagerId || undefined,
      confirmationDate: v.confirmationDate || undefined,
      probationEndDate: v.probationEndDate || undefined,
      bio: v.bio || undefined,
      currentAddress: this.hasAddress(v.currentAddress) ? v.currentAddress : undefined,
      permanentAddress: this.hasAddress(v.permanentAddress) ? v.permanentAddress : undefined,
      sameAsCurrent: v.sameAsCurrent,
      version: v.version,
    };
  }

  // Helpers
  canSubmit(): boolean {
    return this.form.valid && !this.submitting();
  }

  private hasAddress(address: any): boolean {
    return address && address.addressLine1;
  }

  private copyCurrentToPermanent(): void {
    const current = this.form.get('currentAddress')?.value;
    if (current) {
      this.form.get('permanentAddress')?.patchValue(current);
    }
  }

  private markAllAsTouched(): void {
    Object.values(this.form.controls).forEach(control => {
      if (control instanceof FormGroup) {
        Object.values(control.controls).forEach(c => c.markAsTouched());
      } else {
        control.markAsTouched();
      }
    });
  }

  dateNotInFutureValidator(control: AbstractControl): { [key: string]: boolean } | null {
    if (!control.value) return null;
    const date = new Date(control.value);
    if (date > new Date()) {
      return { futureDate: true };
    }
    return null;
  }

  // Photo handling
  onPhotoSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      if (file.size > 5 * 1024 * 1024) {
        this.toast.error('Photo must be less than 5MB');
        return;
      }
      this.photoFile.set(file);

      // Create preview
      const reader = new FileReader();
      reader.onload = () => {
        this.previewPhoto.set(reader.result as string);
        this.cdr.markForCheck();
      };
      reader.readAsDataURL(file);
    }
  }

  removePhoto(): void {
    this.photoFile.set(null);
    this.previewPhoto.set(null);
  }

  // Navigation
  onCancel(): void {
    if (this.isEditMode() && this.staffId()) {
      this.router.navigate(['/staff', this.staffId()]);
    } else {
      this.router.navigate(['/staff']);
    }
  }
}
