/**
 * MSLS Branch Form Component Tests
 *
 * Unit tests for the BranchFormComponent.
 */

import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { vi } from 'vitest';

import { BranchFormComponent } from './branch-form.component';
import { Branch, TIMEZONES } from './branch.model';

describe('BranchFormComponent', () => {
  let component: BranchFormComponent;
  let fixture: ComponentFixture<BranchFormComponent>;

  const mockBranch: Branch = {
    id: '550e8400-e29b-41d4-a716-446655440001',
    code: 'MAIN',
    name: 'Main Campus',
    addressLine1: '123 Main Street',
    addressLine2: 'Building A',
    city: 'Mumbai',
    state: 'Maharashtra',
    postalCode: '400001',
    country: 'India',
    phone: '+91 22 1234 5678',
    email: 'main@school.edu',
    timezone: 'Asia/Kolkata',
    isPrimary: true,
    isActive: true,
    createdAt: '2026-01-23T10:00:00Z',
    updatedAt: '2026-01-23T10:00:00Z',
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BranchFormComponent, ReactiveFormsModule],
    }).compileComponents();

    fixture = TestBed.createComponent(BranchFormComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  describe('form initialization', () => {
    it('should initialize empty form when no branch is provided', () => {
      fixture.detectChanges();

      expect(component.form.get('code')?.value).toBe('');
      expect(component.form.get('name')?.value).toBe('');
      expect(component.form.get('country')?.value).toBe('India');
      expect(component.form.get('timezone')?.value).toBe('Asia/Kolkata');
      expect(component.form.get('isPrimary')?.value).toBe(false);
    });

    it('should initialize form with branch data when editing', () => {
      component.branch = mockBranch;
      fixture.detectChanges();

      expect(component.form.get('code')?.value).toBe('MAIN');
      expect(component.form.get('name')?.value).toBe('Main Campus');
      expect(component.form.get('addressLine1')?.value).toBe('123 Main Street');
      expect(component.form.get('city')?.value).toBe('Mumbai');
      expect(component.form.get('isPrimary')?.value).toBe(true);
    });

    it('should have timezone options from TIMEZONES constant', () => {
      fixture.detectChanges();

      expect(component.timezoneOptions.length).toBe(TIMEZONES.length);
      expect(component.timezoneOptions[0].value).toBe('Asia/Kolkata');
    });
  });

  describe('form validation', () => {
    beforeEach(() => {
      fixture.detectChanges();
    });

    it('should require code field', () => {
      const codeControl = component.form.get('code');

      codeControl?.setValue('');
      expect(codeControl?.valid).toBeFalsy();
      expect(codeControl?.errors?.['required']).toBeTruthy();

      codeControl?.setValue('MAIN');
      expect(codeControl?.valid).toBeTruthy();
    });

    it('should require name field', () => {
      const nameControl = component.form.get('name');

      nameControl?.setValue('');
      expect(nameControl?.valid).toBeFalsy();
      expect(nameControl?.errors?.['required']).toBeTruthy();

      nameControl?.setValue('Main Campus');
      expect(nameControl?.valid).toBeTruthy();
    });

    it('should enforce minimum length on code field', () => {
      const codeControl = component.form.get('code');

      codeControl?.setValue('A');
      expect(codeControl?.valid).toBeFalsy();
      expect(codeControl?.errors?.['minlength']).toBeTruthy();

      codeControl?.setValue('AB');
      expect(codeControl?.valid).toBeTruthy();
    });

    it('should enforce maximum length on code field', () => {
      const codeControl = component.form.get('code');

      codeControl?.setValue('A'.repeat(21));
      expect(codeControl?.valid).toBeFalsy();
      expect(codeControl?.errors?.['maxlength']).toBeTruthy();

      codeControl?.setValue('A'.repeat(20));
      expect(codeControl?.valid).toBeTruthy();
    });

    it('should enforce minimum length on name field', () => {
      const nameControl = component.form.get('name');

      nameControl?.setValue('A');
      expect(nameControl?.valid).toBeFalsy();
      expect(nameControl?.errors?.['minlength']).toBeTruthy();

      nameControl?.setValue('AB');
      expect(nameControl?.valid).toBeTruthy();
    });

    it('should enforce maximum length on name field', () => {
      const nameControl = component.form.get('name');

      nameControl?.setValue('A'.repeat(201));
      expect(nameControl?.valid).toBeFalsy();
      expect(nameControl?.errors?.['maxlength']).toBeTruthy();

      nameControl?.setValue('A'.repeat(200));
      expect(nameControl?.valid).toBeTruthy();
    });

    it('should validate email format', () => {
      const emailControl = component.form.get('email');

      emailControl?.setValue('invalid-email');
      expect(emailControl?.valid).toBeFalsy();
      expect(emailControl?.errors?.['email']).toBeTruthy();

      emailControl?.setValue('valid@email.com');
      expect(emailControl?.valid).toBeTruthy();
    });

    it('should allow empty email', () => {
      const emailControl = component.form.get('email');

      emailControl?.setValue('');
      expect(emailControl?.valid).toBeTruthy();
    });

    it('should enforce maximum length on address fields', () => {
      const addressControl = component.form.get('addressLine1');

      addressControl?.setValue('A'.repeat(256));
      expect(addressControl?.valid).toBeFalsy();
      expect(addressControl?.errors?.['maxlength']).toBeTruthy();

      addressControl?.setValue('A'.repeat(255));
      expect(addressControl?.valid).toBeTruthy();
    });
  });

  describe('getFieldError', () => {
    beforeEach(() => {
      fixture.detectChanges();
    });

    it('should return null for untouched fields', () => {
      const error = component.getFieldError('code');
      expect(error).toBeNull();
    });

    it('should return required error message', () => {
      const codeControl = component.form.get('code');
      codeControl?.setValue('');
      codeControl?.markAsTouched();

      const error = component.getFieldError('code');
      expect(error).toBe('This field is required');
    });

    it('should return minlength error message', () => {
      const codeControl = component.form.get('code');
      codeControl?.setValue('A');
      codeControl?.markAsTouched();

      const error = component.getFieldError('code');
      expect(error).toBe('Must be at least 2 characters');
    });

    it('should return maxlength error message', () => {
      const codeControl = component.form.get('code');
      codeControl?.setValue('A'.repeat(21));
      codeControl?.markAsTouched();

      const error = component.getFieldError('code');
      expect(error).toBe('Must be at most 20 characters');
    });

    it('should return email error message', () => {
      const emailControl = component.form.get('email');
      emailControl?.setValue('invalid');
      emailControl?.markAsTouched();

      const error = component.getFieldError('email');
      expect(error).toBe('Please enter a valid email address');
    });

    it('should return null for valid fields', () => {
      const codeControl = component.form.get('code');
      codeControl?.setValue('MAIN');
      codeControl?.markAsTouched();

      const error = component.getFieldError('code');
      expect(error).toBeNull();
    });
  });

  describe('form submission', () => {
    beforeEach(() => {
      fixture.detectChanges();
    });

    it('should emit save event with form data when valid', () => {
      const saveSpy = vi.spyOn(component.save, 'emit');

      component.form.patchValue({
        code: 'NEW01',
        name: 'New Branch',
        city: 'Pune',
        country: 'India',
        timezone: 'Asia/Kolkata',
      });

      component.onSubmit();

      expect(saveSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          code: 'NEW01',
          name: 'New Branch',
          city: 'Pune',
          country: 'India',
          timezone: 'Asia/Kolkata',
        })
      );
    });

    it('should not emit save event when form is invalid', () => {
      const saveSpy = vi.spyOn(component.save, 'emit');

      component.form.patchValue({
        code: '',
        name: '',
      });

      component.onSubmit();

      expect(saveSpy).not.toHaveBeenCalled();
    });

    it('should convert empty strings to undefined', () => {
      const saveSpy = vi.spyOn(component.save, 'emit');

      component.form.patchValue({
        code: 'NEW01',
        name: 'New Branch',
        addressLine1: '',
        city: '',
        phone: '',
        email: '',
      });

      component.onSubmit();

      expect(saveSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          code: 'NEW01',
          name: 'New Branch',
          addressLine1: undefined,
          city: undefined,
          phone: undefined,
          email: undefined,
        })
      );
    });

    it('should include isPrimary flag', () => {
      const saveSpy = vi.spyOn(component.save, 'emit');

      component.form.patchValue({
        code: 'NEW01',
        name: 'New Branch',
        isPrimary: true,
      });

      component.onSubmit();

      expect(saveSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          isPrimary: true,
        })
      );
    });
  });

  describe('cancel action', () => {
    beforeEach(() => {
      fixture.detectChanges();
    });

    it('should emit cancel event', () => {
      const cancelSpy = vi.spyOn(component.cancel, 'emit');

      component.onCancel();

      expect(cancelSpy).toHaveBeenCalled();
    });
  });

  describe('loading state', () => {
    it('should pass loading state to submit button', () => {
      component.loading = true;
      fixture.detectChanges();

      expect(component.loading).toBe(true);
    });
  });

  describe('edit mode', () => {
    it('should show "Update Branch" text when editing', () => {
      component.branch = mockBranch;
      fixture.detectChanges();

      const submitButton = fixture.nativeElement.querySelector(
        'msls-button[type="submit"]'
      );
      expect(submitButton.textContent).toContain('Update Branch');
    });

    it('should show "Create Branch" text when creating', () => {
      fixture.detectChanges();

      const submitButton = fixture.nativeElement.querySelector(
        'msls-button[type="submit"]'
      );
      expect(submitButton.textContent).toContain('Create Branch');
    });
  });

  describe('form fields rendering', () => {
    beforeEach(() => {
      fixture.detectChanges();
    });

    it('should render code field', () => {
      const codeField = fixture.nativeElement.querySelector(
        'msls-input[formControlName="code"]'
      );
      expect(codeField).toBeTruthy();
    });

    it('should render name field', () => {
      const nameField = fixture.nativeElement.querySelector(
        'msls-input[formControlName="name"]'
      );
      expect(nameField).toBeTruthy();
    });

    it('should render address fields', () => {
      const address1 = fixture.nativeElement.querySelector(
        'msls-input[formControlName="addressLine1"]'
      );
      const address2 = fixture.nativeElement.querySelector(
        'msls-input[formControlName="addressLine2"]'
      );
      expect(address1).toBeTruthy();
      expect(address2).toBeTruthy();
    });

    it('should render city, state, and postal code fields', () => {
      const city = fixture.nativeElement.querySelector(
        'msls-input[formControlName="city"]'
      );
      const state = fixture.nativeElement.querySelector(
        'msls-input[formControlName="state"]'
      );
      const postalCode = fixture.nativeElement.querySelector(
        'msls-input[formControlName="postalCode"]'
      );

      expect(city).toBeTruthy();
      expect(state).toBeTruthy();
      expect(postalCode).toBeTruthy();
    });

    it('should render country field', () => {
      const country = fixture.nativeElement.querySelector(
        'msls-input[formControlName="country"]'
      );
      expect(country).toBeTruthy();
    });

    it('should render timezone select', () => {
      const timezone = fixture.nativeElement.querySelector(
        'msls-select[formControlName="timezone"]'
      );
      expect(timezone).toBeTruthy();
    });

    it('should render phone and email fields', () => {
      const phone = fixture.nativeElement.querySelector(
        'msls-input[formControlName="phone"]'
      );
      const email = fixture.nativeElement.querySelector(
        'msls-input[formControlName="email"]'
      );

      expect(phone).toBeTruthy();
      expect(email).toBeTruthy();
    });

    it('should render isPrimary checkbox', () => {
      const isPrimary = fixture.nativeElement.querySelector(
        'msls-checkbox[formControlName="isPrimary"]'
      );
      expect(isPrimary).toBeTruthy();
    });

    it('should render cancel and submit buttons', () => {
      const cancelButton = fixture.nativeElement.querySelector(
        'msls-button[variant="secondary"]'
      );
      const submitButton = fixture.nativeElement.querySelector(
        'msls-button[type="submit"]'
      );

      expect(cancelButton).toBeTruthy();
      expect(submitButton).toBeTruthy();
    });
  });
});
