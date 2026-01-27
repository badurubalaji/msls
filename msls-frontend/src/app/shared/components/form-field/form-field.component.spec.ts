import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Component } from '@angular/core';
import { MslsFormFieldComponent } from './form-field.component';

@Component({
  standalone: true,
  imports: [MslsFormFieldComponent],
  template: `
    <msls-form-field
      [label]="label"
      [hint]="hint"
      [error]="error"
      [required]="required"
    >
      <input class="input" type="text" />
    </msls-form-field>
  `,
})
class TestHostComponent {
  label = 'Test Label';
  hint = 'Test hint text';
  error = '';
  required = false;
}

describe('MslsFormFieldComponent', () => {
  let component: TestHostComponent;
  let fixture: ComponentFixture<TestHostComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TestHostComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TestHostComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display label', () => {
    const labelElement = fixture.nativeElement.querySelector('.form-field__label');
    expect(labelElement.textContent).toContain('Test Label');
  });

  it('should display hint text when no error', () => {
    const hintElement = fixture.nativeElement.querySelector('.form-field__hint');
    expect(hintElement.textContent).toContain('Test hint text');
  });

  it('should show required indicator when required is true', () => {
    component.required = true;
    fixture.detectChanges();
    const requiredElement = fixture.nativeElement.querySelector('.form-field__required');
    expect(requiredElement).toBeTruthy();
    expect(requiredElement.textContent).toBe('*');
  });

  it('should not show required indicator when required is false', () => {
    component.required = false;
    fixture.detectChanges();
    const requiredElement = fixture.nativeElement.querySelector('.form-field__required');
    expect(requiredElement).toBeFalsy();
  });

  it('should display error message when error is provided', () => {
    component.error = 'This field is required';
    fixture.detectChanges();
    const errorElement = fixture.nativeElement.querySelector('.form-field__error');
    expect(errorElement.textContent).toContain('This field is required');
  });

  it('should hide hint text when error is displayed', () => {
    component.error = 'This field is required';
    fixture.detectChanges();
    const hintElement = fixture.nativeElement.querySelector('.form-field__hint');
    expect(hintElement).toBeFalsy();
  });

  it('should project input content', () => {
    const inputElement = fixture.nativeElement.querySelector('input');
    expect(inputElement).toBeTruthy();
  });
});
