import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Component } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MslsCheckboxComponent } from './checkbox.component';

@Component({
  standalone: true,
  imports: [MslsCheckboxComponent, ReactiveFormsModule],
  template: `
    <msls-checkbox
      [label]="label"
      [disabled]="disabled"
      [indeterminate]="indeterminate"
      [formControl]="control"
    ></msls-checkbox>
  `,
})
class TestHostComponent {
  label = 'Test checkbox';
  disabled = false;
  indeterminate = false;
  control = new FormControl(false);
}

describe('MslsCheckboxComponent', () => {
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

  it('should display label text', () => {
    const labelElement = fixture.nativeElement.querySelector('.checkbox__label');
    expect(labelElement.textContent).toBe('Test checkbox');
  });

  it('should be unchecked by default', () => {
    const inputElement = fixture.nativeElement.querySelector('.checkbox__input');
    expect(inputElement.checked).toBeFalsy();
  });

  it('should toggle checked state when clicked', () => {
    const labelElement = fixture.nativeElement.querySelector('.checkbox');
    labelElement.click();
    fixture.detectChanges();

    expect(component.control.value).toBe(true);

    labelElement.click();
    fixture.detectChanges();

    expect(component.control.value).toBe(false);
  });

  it('should update form control value', () => {
    const inputElement = fixture.nativeElement.querySelector('.checkbox__input');
    inputElement.click();
    fixture.detectChanges();

    expect(component.control.value).toBe(true);
  });

  it('should reflect form control value changes', () => {
    component.control.setValue(true);
    fixture.detectChanges();

    const inputElement = fixture.nativeElement.querySelector('.checkbox__input');
    expect(inputElement.checked).toBe(true);
  });

  it('should be disabled when disabled input is true', () => {
    component.disabled = true;
    fixture.detectChanges();

    const checkboxElement = fixture.nativeElement.querySelector('.checkbox');
    expect(checkboxElement.classList.contains('checkbox--disabled')).toBeTruthy();

    const inputElement = fixture.nativeElement.querySelector('.checkbox__input');
    expect(inputElement.disabled).toBeTruthy();
  });

  it('should be disabled when form control is disabled', () => {
    component.control.disable();
    fixture.detectChanges();

    const checkboxElement = fixture.nativeElement.querySelector('.checkbox');
    expect(checkboxElement.classList.contains('checkbox--disabled')).toBeTruthy();
  });

  it('should show checkmark when checked', () => {
    component.control.setValue(true);
    fixture.detectChanges();

    const boxElement = fixture.nativeElement.querySelector('.checkbox__box');
    expect(boxElement.classList.contains('checkbox__box--checked')).toBeTruthy();
    expect(boxElement.querySelector('svg')).toBeTruthy();
  });

  it('should show indeterminate state', () => {
    component.indeterminate = true;
    fixture.detectChanges();

    const boxElement = fixture.nativeElement.querySelector('.checkbox__box');
    expect(boxElement.classList.contains('checkbox__box--indeterminate')).toBeTruthy();
  });

  it('should not change state when disabled', () => {
    component.disabled = true;
    fixture.detectChanges();

    const labelElement = fixture.nativeElement.querySelector('.checkbox');
    labelElement.click();
    fixture.detectChanges();

    expect(component.control.value).toBe(false);
  });
});
