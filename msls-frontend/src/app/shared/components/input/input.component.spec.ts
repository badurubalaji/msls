import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { Component } from '@angular/core';
import { FormsModule, ReactiveFormsModule, FormControl } from '@angular/forms';
import { vi } from 'vitest';
import { MslsInputComponent, InputType } from './input.component';

// Test host component for ngModel testing
@Component({
  template: `
    <msls-input
      [(ngModel)]="value"
      [label]="label"
      [error]="error"
      [disabled]="disabled"
    />
  `,
  standalone: true,
  imports: [MslsInputComponent, FormsModule],
})
class TestHostComponent {
  value = '';
  label = '';
  error = '';
  disabled = false;
}

// Test host component for reactive forms testing
@Component({
  template: `
    <msls-input [formControl]="control" />
  `,
  standalone: true,
  imports: [MslsInputComponent, ReactiveFormsModule],
})
class TestReactiveHostComponent {
  control = new FormControl('');
}

describe('MslsInputComponent', () => {
  let component: MslsInputComponent;
  let fixture: ComponentFixture<MslsInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MslsInputComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MslsInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('default values', () => {
    it('should have text type by default', () => {
      expect(component.type()).toBe('text');
    });

    it('should have empty placeholder by default', () => {
      expect(component.placeholder()).toBe('');
    });

    it('should not be disabled by default', () => {
      expect(component.disabled()).toBe(false);
    });

    it('should not be required by default', () => {
      expect(component.required()).toBe(false);
    });
  });

  describe('input types', () => {
    const types: InputType[] = ['text', 'email', 'password', 'number', 'tel', 'search'];

    types.forEach((type) => {
      it(`should set input type to ${type}`, () => {
        fixture.componentRef.setInput('type', type);
        fixture.detectChanges();

        const input = fixture.debugElement.query(By.css('input'));
        expect(input.nativeElement.type).toBe(type);
      });
    });
  });

  describe('label', () => {
    it('should render label when provided', () => {
      fixture.componentRef.setInput('label', 'Email Address');
      fixture.detectChanges();

      const label = fixture.debugElement.query(By.css('.msls-input__label'));
      expect(label).toBeTruthy();
      expect(label.nativeElement.textContent).toContain('Email Address');
    });

    it('should not render label when empty', () => {
      const label = fixture.debugElement.query(By.css('.msls-input__label'));
      expect(label).toBeFalsy();
    });

    it('should show required indicator when required', () => {
      fixture.componentRef.setInput('label', 'Email');
      fixture.componentRef.setInput('required', true);
      fixture.detectChanges();

      const required = fixture.debugElement.query(By.css('.msls-input__required'));
      expect(required).toBeTruthy();
      expect(required.nativeElement.textContent).toBe('*');
    });
  });

  describe('placeholder', () => {
    it('should set placeholder attribute', () => {
      fixture.componentRef.setInput('placeholder', 'Enter your email');
      fixture.detectChanges();

      const input = fixture.debugElement.query(By.css('input'));
      expect(input.nativeElement.placeholder).toBe('Enter your email');
    });
  });

  describe('error state', () => {
    it('should show error message when error is set', () => {
      fixture.componentRef.setInput('error', 'Invalid email address');
      fixture.detectChanges();

      const error = fixture.debugElement.query(By.css('.msls-input__error'));
      expect(error).toBeTruthy();
      expect(error.nativeElement.textContent).toContain('Invalid email address');
    });

    it('should apply error class when error is set', () => {
      fixture.componentRef.setInput('error', 'Invalid email');
      fixture.detectChanges();

      const container = fixture.debugElement.query(By.css('.msls-input'));
      expect(container.nativeElement.classList.contains('msls-input--error')).toBe(true);
    });

    it('should set aria-invalid when error is set', () => {
      fixture.componentRef.setInput('error', 'Invalid email');
      fixture.detectChanges();

      const input = fixture.debugElement.query(By.css('input'));
      expect(input.nativeElement.getAttribute('aria-invalid')).toBe('true');
    });

    it('should have role="alert" on error message', () => {
      fixture.componentRef.setInput('error', 'Error message');
      fixture.detectChanges();

      const error = fixture.debugElement.query(By.css('.msls-input__error'));
      expect(error.nativeElement.getAttribute('role')).toBe('alert');
    });
  });

  describe('hint', () => {
    it('should show hint when no error', () => {
      fixture.componentRef.setInput('hint', 'We will never share your email');
      fixture.detectChanges();

      const hint = fixture.debugElement.query(By.css('.msls-input__hint'));
      expect(hint).toBeTruthy();
      expect(hint.nativeElement.textContent).toContain('We will never share your email');
    });

    it('should hide hint when error is present', () => {
      fixture.componentRef.setInput('hint', 'Hint text');
      fixture.componentRef.setInput('error', 'Error text');
      fixture.detectChanges();

      const hint = fixture.debugElement.query(By.css('.msls-input__hint'));
      expect(hint).toBeFalsy();
    });
  });

  describe('disabled state', () => {
    it('should disable input when disabled is true', () => {
      fixture.componentRef.setInput('disabled', true);
      fixture.detectChanges();

      const input = fixture.debugElement.query(By.css('input'));
      expect(input.nativeElement.disabled).toBe(true);
    });

    it('should apply disabled class when disabled', () => {
      fixture.componentRef.setInput('disabled', true);
      fixture.detectChanges();

      const container = fixture.debugElement.query(By.css('.msls-input'));
      expect(container.nativeElement.classList.contains('msls-input--disabled')).toBe(true);
    });
  });

  describe('icons', () => {
    it('should show prefix icon when provided', () => {
      fixture.componentRef.setInput('prefixIcon', 'envelope');
      fixture.detectChanges();

      const prefixIcon = fixture.debugElement.query(By.css('.msls-input__icon--prefix'));
      expect(prefixIcon).toBeTruthy();
    });

    it('should show suffix icon when provided', () => {
      fixture.componentRef.setInput('suffixIcon', 'eye');
      fixture.detectChanges();

      const suffixIcon = fixture.debugElement.query(By.css('.msls-input__icon--suffix'));
      expect(suffixIcon).toBeTruthy();
    });

    it('should apply has-prefix class when prefix icon is set', () => {
      fixture.componentRef.setInput('prefixIcon', 'envelope');
      fixture.detectChanges();

      const container = fixture.debugElement.query(By.css('.msls-input'));
      expect(container.nativeElement.classList.contains('msls-input--has-prefix')).toBe(true);
    });

    it('should apply has-suffix class when suffix icon is set', () => {
      fixture.componentRef.setInput('suffixIcon', 'eye');
      fixture.detectChanges();

      const container = fixture.debugElement.query(By.css('.msls-input'));
      expect(container.nativeElement.classList.contains('msls-input--has-suffix')).toBe(true);
    });
  });

  describe('ControlValueAccessor', () => {
    it('should write value', () => {
      component.writeValue('test@example.com');
      expect(component.value()).toBe('test@example.com');
    });

    it('should handle null value', () => {
      component.writeValue(null as any);
      expect(component.value()).toBe('');
    });

    it('should call onChange on input', () => {
      const onChangeSpy = vi.fn();
      component.registerOnChange(onChangeSpy);

      const input = fixture.debugElement.query(By.css('input'));
      input.nativeElement.value = 'new value';
      input.nativeElement.dispatchEvent(new Event('input'));

      expect(onChangeSpy).toHaveBeenCalledWith('new value');
    });

    it('should call onTouched on blur', () => {
      const onTouchedSpy = vi.fn();
      component.registerOnTouched(onTouchedSpy);

      const input = fixture.debugElement.query(By.css('input'));
      input.nativeElement.dispatchEvent(new Event('blur'));

      expect(onTouchedSpy).toHaveBeenCalled();
    });

    it('should disable input when setDisabledState is called', () => {
      component.setDisabledState(true);
      fixture.detectChanges();

      const input = fixture.debugElement.query(By.css('input'));
      expect(input.nativeElement.disabled).toBe(true);
    });
  });
});

describe('MslsInputComponent with ngModel', () => {
  let fixture: ComponentFixture<TestHostComponent>;
  let hostComponent: TestHostComponent;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TestHostComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TestHostComponent);
    hostComponent = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should update ngModel value', async () => {
    const input = fixture.debugElement.query(By.css('input'));
    input.nativeElement.value = 'test value';
    input.nativeElement.dispatchEvent(new Event('input'));

    fixture.detectChanges();
    await fixture.whenStable();

    expect(hostComponent.value).toBe('test value');
  });

  it('should reflect ngModel changes in input', async () => {
    hostComponent.value = 'initial value';
    fixture.detectChanges();
    await fixture.whenStable();

    const input = fixture.debugElement.query(By.css('input'));
    expect(input.nativeElement.value).toBe('initial value');
  });
});

describe('MslsInputComponent with FormControl', () => {
  let fixture: ComponentFixture<TestReactiveHostComponent>;
  let hostComponent: TestReactiveHostComponent;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TestReactiveHostComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TestReactiveHostComponent);
    hostComponent = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should update form control value', () => {
    const input = fixture.debugElement.query(By.css('input'));
    input.nativeElement.value = 'form value';
    input.nativeElement.dispatchEvent(new Event('input'));

    expect(hostComponent.control.value).toBe('form value');
  });

  it('should reflect form control changes in input', () => {
    hostComponent.control.setValue('control value');
    fixture.detectChanges();

    const input = fixture.debugElement.query(By.css('input'));
    expect(input.nativeElement.value).toBe('control value');
  });

  it('should disable input when form control is disabled', () => {
    hostComponent.control.disable();
    fixture.detectChanges();

    const input = fixture.debugElement.query(By.css('input'));
    expect(input.nativeElement.disabled).toBe(true);
  });
});
