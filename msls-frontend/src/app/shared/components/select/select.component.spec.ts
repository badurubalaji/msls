import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Component, signal } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MslsSelectComponent, SelectOption } from './select.component';

@Component({
  standalone: true,
  imports: [MslsSelectComponent, ReactiveFormsModule],
  template: `
    <msls-select
      [options]="options"
      [placeholder]="placeholder"
      [searchable]="searchable"
      [multiple]="multiple"
      [disabled]="disabled"
      [formControl]="control"
    ></msls-select>
  `,
})
class TestHostComponent {
  options: SelectOption[] = [
    { value: '1', label: 'Option 1' },
    { value: '2', label: 'Option 2' },
    { value: '3', label: 'Option 3' },
    { value: '4', label: 'Disabled Option', disabled: true },
  ];
  placeholder = 'Select an option';
  searchable = false;
  multiple = false;
  disabled = false;
  control = new FormControl<string | string[] | null>(null);
}

describe('MslsSelectComponent', () => {
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

  it('should display placeholder when no value is selected', () => {
    const triggerElement = fixture.nativeElement.querySelector('.select__value');
    expect(triggerElement.textContent.trim()).toBe('Select an option');
  });

  it('should open dropdown when trigger is clicked', () => {
    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const dropdownElement = fixture.nativeElement.querySelector('.select__dropdown');
    expect(dropdownElement).toBeTruthy();
  });

  it('should display all options', () => {
    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const optionElements = fixture.nativeElement.querySelectorAll('.select__option');
    expect(optionElements.length).toBe(4);
  });

  it('should select option when clicked', () => {
    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const optionElements = fixture.nativeElement.querySelectorAll('.select__option');
    optionElements[1].click();
    fixture.detectChanges();

    expect(component.control.value).toBe('2');
    const valueElement = fixture.nativeElement.querySelector('.select__value');
    expect(valueElement.textContent.trim()).toBe('Option 2');
  });

  it('should close dropdown after selecting (single mode)', () => {
    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const optionElements = fixture.nativeElement.querySelectorAll('.select__option');
    optionElements[0].click();
    fixture.detectChanges();

    const dropdownElement = fixture.nativeElement.querySelector('.select__dropdown');
    expect(dropdownElement).toBeFalsy();
  });

  it('should not select disabled option', () => {
    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const optionElements = fixture.nativeElement.querySelectorAll('.select__option');
    optionElements[3].click(); // Disabled option
    fixture.detectChanges();

    expect(component.control.value).toBeNull();
  });

  it('should clear selection when clear button is clicked', () => {
    component.control.setValue('1');
    fixture.detectChanges();

    const clearButton = fixture.nativeElement.querySelector('.select__clear');
    clearButton.click();
    fixture.detectChanges();

    expect(component.control.value).toBeNull();
  });

  it('should filter options when searchable', () => {
    component.searchable = true;
    fixture.detectChanges();

    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const searchInput = fixture.nativeElement.querySelector('.select__search-input');
    searchInput.value = 'Option 2';
    searchInput.dispatchEvent(new Event('input'));
    fixture.detectChanges();

    const optionElements = fixture.nativeElement.querySelectorAll('.select__option');
    expect(optionElements.length).toBe(1);
    expect(optionElements[0].textContent).toContain('Option 2');
  });

  it('should show empty message when no options match search', () => {
    component.searchable = true;
    fixture.detectChanges();

    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const searchInput = fixture.nativeElement.querySelector('.select__search-input');
    searchInput.value = 'xyz';
    searchInput.dispatchEvent(new Event('input'));
    fixture.detectChanges();

    const emptyElement = fixture.nativeElement.querySelector('.select__empty');
    expect(emptyElement).toBeTruthy();
    expect(emptyElement.textContent).toContain('No options found');
  });

  it('should be disabled when disabled input is true', () => {
    component.disabled = true;
    fixture.detectChanges();

    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    expect(triggerElement.disabled).toBeTruthy();

    const selectElement = fixture.nativeElement.querySelector('.select');
    expect(selectElement.classList.contains('select--disabled')).toBeTruthy();
  });

  it('should handle multiple selection', () => {
    component.multiple = true;
    fixture.detectChanges();

    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const optionElements = fixture.nativeElement.querySelectorAll('.select__option');
    optionElements[0].click();
    fixture.detectChanges();

    optionElements[1].click();
    fixture.detectChanges();

    expect(component.control.value).toEqual(['1', '2']);
  });

  it('should toggle selection in multiple mode', () => {
    component.multiple = true;
    component.control.setValue(['1', '2']);
    fixture.detectChanges();

    const triggerElement = fixture.nativeElement.querySelector('.select__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const optionElements = fixture.nativeElement.querySelectorAll('.select__option');
    optionElements[0].click(); // Deselect Option 1
    fixture.detectChanges();

    expect(component.control.value).toEqual(['2']);
  });
});
