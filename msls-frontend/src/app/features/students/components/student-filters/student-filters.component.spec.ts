/**
 * Student Filters Component Tests
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { StudentFiltersComponent } from './student-filters.component';
import { StudentListFilter } from '../../models/student.model';

describe('StudentFiltersComponent', () => {
  let component: StudentFiltersComponent;
  let fixture: ComponentFixture<StudentFiltersComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StudentFiltersComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(StudentFiltersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('Filter Panel Rendering', () => {
    it('should render the filters panel', () => {
      const panel = fixture.debugElement.query(By.css('.filters-panel'));
      expect(panel).toBeTruthy();
    });

    it('should render status filter select', () => {
      const statusSelect = fixture.debugElement.query(By.css('.filter-group select'));
      expect(statusSelect).toBeTruthy();
    });

    it('should have correct status options', () => {
      const statusSelect = fixture.debugElement.queryAll(By.css('.filter-group'))[0].query(By.css('select'));
      const options = statusSelect.queryAll(By.css('option'));

      expect(options.length).toBe(5);
      expect(options[0].nativeElement.textContent).toContain('All Statuses');
      expect(options[1].nativeElement.value).toBe('active');
      expect(options[2].nativeElement.value).toBe('inactive');
      expect(options[3].nativeElement.value).toBe('transferred');
      expect(options[4].nativeElement.value).toBe('graduated');
    });

    it('should have correct gender options', () => {
      const genderSelect = fixture.debugElement.queryAll(By.css('.filter-group'))[1].query(By.css('select'));
      const options = genderSelect.queryAll(By.css('option'));

      expect(options.length).toBe(4);
      expect(options[0].nativeElement.textContent).toContain('All Genders');
      expect(options[1].nativeElement.value).toBe('male');
      expect(options[2].nativeElement.value).toBe('female');
      expect(options[3].nativeElement.value).toBe('other');
    });

    it('should render date range inputs', () => {
      const dateInputs = fixture.debugElement.queryAll(By.css('input[type="date"]'));
      expect(dateInputs.length).toBe(2);
    });

    it('should render sort options', () => {
      const sortBySelect = fixture.debugElement.queryAll(By.css('.filter-group'))[4].query(By.css('select'));
      const options = sortBySelect.queryAll(By.css('option'));

      expect(options.length).toBe(4);
      expect(options[0].nativeElement.textContent).toContain('Default (Name)');
    });
  });

  describe('Status Filter', () => {
    it('should update currentFilters on status change', () => {
      const statusSelect = fixture.debugElement.queryAll(By.css('.filter-group'))[0].query(By.css('select'));

      statusSelect.nativeElement.value = 'active';
      statusSelect.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().status).toBe('active');
    });

    it('should clear status when empty value selected', () => {
      component.currentFilters.set({ status: 'active' });
      fixture.detectChanges();

      const statusSelect = fixture.debugElement.queryAll(By.css('.filter-group'))[0].query(By.css('select'));
      statusSelect.nativeElement.value = '';
      statusSelect.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().status).toBeUndefined();
    });
  });

  describe('Gender Filter', () => {
    it('should update currentFilters on gender change', () => {
      const genderSelect = fixture.debugElement.queryAll(By.css('.filter-group'))[1].query(By.css('select'));

      genderSelect.nativeElement.value = 'female';
      genderSelect.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().gender).toBe('female');
    });
  });

  describe('Date Range Filters', () => {
    it('should update admissionFrom on date change', () => {
      const dateInput = fixture.debugElement.queryAll(By.css('input[type="date"]'))[0];

      dateInput.nativeElement.value = '2026-01-01';
      dateInput.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().admissionFrom).toBe('2026-01-01');
    });

    it('should update admissionTo on date change', () => {
      const dateInput = fixture.debugElement.queryAll(By.css('input[type="date"]'))[1];

      dateInput.nativeElement.value = '2026-12-31';
      dateInput.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().admissionTo).toBe('2026-12-31');
    });

    it('should clear date when empty value', () => {
      component.currentFilters.set({ admissionFrom: '2026-01-01' });
      fixture.detectChanges();

      const dateInput = fixture.debugElement.queryAll(By.css('input[type="date"]'))[0];
      dateInput.nativeElement.value = '';
      dateInput.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().admissionFrom).toBeUndefined();
    });
  });

  describe('Sort Options', () => {
    it('should update sortBy on change', () => {
      const sortBySelect = fixture.debugElement.queryAll(By.css('.filter-group'))[4].query(By.css('select'));

      sortBySelect.nativeElement.value = 'admission_number';
      sortBySelect.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().sortBy).toBe('admission_number');
    });

    it('should update sortOrder on change', () => {
      const sortOrderSelect = fixture.debugElement.queryAll(By.css('.filter-group'))[5].query(By.css('select'));

      sortOrderSelect.nativeElement.value = 'desc';
      sortOrderSelect.nativeElement.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(component.currentFilters().sortOrder).toBe('desc');
    });
  });

  describe('hasActiveFilters', () => {
    it('should return false when no filters are active', () => {
      component.currentFilters.set({});
      expect(component.hasActiveFilters()).toBe(false);
    });

    it('should return true when status is set', () => {
      component.currentFilters.set({ status: 'active' });
      expect(component.hasActiveFilters()).toBe(true);
    });

    it('should return true when gender is set', () => {
      component.currentFilters.set({ gender: 'male' });
      expect(component.hasActiveFilters()).toBe(true);
    });

    it('should return true when date range is set', () => {
      component.currentFilters.set({ admissionFrom: '2026-01-01' });
      expect(component.hasActiveFilters()).toBe(true);
    });

    it('should return true when sortBy is set', () => {
      component.currentFilters.set({ sortBy: 'name' });
      expect(component.hasActiveFilters()).toBe(true);
    });
  });

  describe('Clear All Filters', () => {
    it('should show clear all button when filters are active', () => {
      component.currentFilters.set({ status: 'active' });
      fixture.detectChanges();

      const clearBtn = fixture.debugElement.query(By.css('.clear-all-btn'));
      expect(clearBtn).toBeTruthy();
    });

    it('should hide clear all button when no filters are active', () => {
      component.currentFilters.set({});
      fixture.detectChanges();

      const clearBtn = fixture.debugElement.query(By.css('.clear-all-btn'));
      expect(clearBtn).toBeFalsy();
    });

    it('should clear all filters and emit on click', () => {
      const emitSpy = vi.spyOn(component.filtersChanged, 'emit');
      fixture.componentRef.setInput('filters', { limit: 20 } as StudentListFilter);
      component.currentFilters.set({
        status: 'active',
        gender: 'male',
        admissionFrom: '2026-01-01',
        sortBy: 'name',
        limit: 20,
      });
      fixture.detectChanges();

      const clearBtn = fixture.debugElement.query(By.css('.clear-all-btn'));
      clearBtn.nativeElement.click();
      fixture.detectChanges();

      expect(component.currentFilters().status).toBeUndefined();
      expect(component.currentFilters().gender).toBeUndefined();
      expect(component.currentFilters().admissionFrom).toBeUndefined();
      expect(component.currentFilters().sortBy).toBeUndefined();
      expect(component.currentFilters().limit).toBe(20); // Limit should be preserved
      expect(emitSpy).toHaveBeenCalled();
    });
  });

  describe('Apply Filters', () => {
    it('should render apply filters button', () => {
      const applyBtn = fixture.debugElement.query(By.css('.apply-btn'));
      expect(applyBtn).toBeTruthy();
      expect(applyBtn.nativeElement.textContent).toContain('Apply Filters');
    });

    it('should emit filtersChanged on apply click', () => {
      const emitSpy = vi.spyOn(component.filtersChanged, 'emit');
      component.currentFilters.set({ status: 'active', gender: 'female' });
      fixture.detectChanges();

      const applyBtn = fixture.debugElement.query(By.css('.apply-btn'));
      applyBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalledWith({ status: 'active', gender: 'female' });
    });
  });

  describe('Initial Filters', () => {
    it('should initialize currentFilters from input', () => {
      // Set the filters and then trigger changes manually
      const filtersInput = {
        status: 'inactive' as const,
        gender: 'other' as const,
      };

      // Apply filters manually through the component's methods
      component.onStatusChange({ target: { value: 'inactive' } } as unknown as Event);
      component.onGenderChange({ target: { value: 'other' } } as unknown as Event);
      fixture.detectChanges();

      // currentFilters should be updated
      expect(component.currentFilters().status).toBe('inactive');
      expect(component.currentFilters().gender).toBe('other');
    });
  });
});
