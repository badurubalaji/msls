/**
 * Student Search Component Tests
 */

import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { firstValueFrom, take, timeout } from 'rxjs';

import { StudentSearchComponent } from './student-search.component';
import { StudentListFilter } from '../../models/student.model';

describe('StudentSearchComponent', () => {
  let component: StudentSearchComponent;
  let fixture: ComponentFixture<StudentSearchComponent>;

  beforeEach(async () => {
    vi.useFakeTimers();
    await TestBed.configureTestingModule({
      imports: [StudentSearchComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(StudentSearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.useRealTimers();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('Search Input', () => {
    it('should render search input with placeholder', () => {
      const input = fixture.debugElement.query(By.css('.search-input'));
      expect(input).toBeTruthy();
      expect(input.nativeElement.placeholder).toBe('Search by name, admission number, phone...');
    });

    it('should update searchQuery on input', () => {
      const input = fixture.debugElement.query(By.css('.search-input'));
      input.nativeElement.value = 'John';
      input.nativeElement.dispatchEvent(new Event('input'));
      fixture.detectChanges();

      expect(component.searchQuery()).toBe('John');
    });

    it('should emit searchChanged after debounce', async () => {
      const emitSpy = vi.spyOn(component.searchChanged, 'emit');
      const input = fixture.debugElement.query(By.css('.search-input'));

      input.nativeElement.value = 'John Doe';
      input.nativeElement.dispatchEvent(new Event('input'));
      fixture.detectChanges();

      // Should not emit immediately
      expect(emitSpy).not.toHaveBeenCalled();

      // Advance timers by 300ms (debounce time)
      await vi.advanceTimersByTimeAsync(300);

      expect(emitSpy).toHaveBeenCalledWith('John Doe');
    });

    it('should debounce multiple rapid inputs', async () => {
      const emitSpy = vi.spyOn(component.searchChanged, 'emit');
      const input = fixture.debugElement.query(By.css('.search-input'));

      // Simulate rapid typing
      input.nativeElement.value = 'J';
      input.nativeElement.dispatchEvent(new Event('input'));
      await vi.advanceTimersByTimeAsync(100);

      input.nativeElement.value = 'Jo';
      input.nativeElement.dispatchEvent(new Event('input'));
      await vi.advanceTimersByTimeAsync(100);

      input.nativeElement.value = 'John';
      input.nativeElement.dispatchEvent(new Event('input'));
      await vi.advanceTimersByTimeAsync(300);

      // Should only emit once with final value
      expect(emitSpy).toHaveBeenCalledTimes(1);
      expect(emitSpy).toHaveBeenCalledWith('John');
    });

    it('should not emit for same value (distinctUntilChanged)', async () => {
      const emitSpy = vi.spyOn(component.searchChanged, 'emit');
      const input = fixture.debugElement.query(By.css('.search-input'));

      input.nativeElement.value = 'John';
      input.nativeElement.dispatchEvent(new Event('input'));
      await vi.advanceTimersByTimeAsync(300);

      // Type same value again
      input.nativeElement.value = 'John';
      input.nativeElement.dispatchEvent(new Event('input'));
      await vi.advanceTimersByTimeAsync(300);

      // Should only emit once
      expect(emitSpy).toHaveBeenCalledTimes(1);
    });
  });

  describe('Clear Search', () => {
    it('should show clear button when search has value', () => {
      component.searchQuery.set('test');
      fixture.detectChanges();

      const clearBtn = fixture.debugElement.query(By.css('.clear-btn'));
      expect(clearBtn).toBeTruthy();
    });

    it('should hide clear button when search is empty', () => {
      component.searchQuery.set('');
      fixture.detectChanges();

      const clearBtn = fixture.debugElement.query(By.css('.clear-btn'));
      expect(clearBtn).toBeFalsy();
    });

    it('should clear search on button click', async () => {
      const emitSpy = vi.spyOn(component.searchChanged, 'emit');
      component.searchQuery.set('test');
      fixture.detectChanges();

      const clearBtn = fixture.debugElement.query(By.css('.clear-btn'));
      clearBtn.nativeElement.click();
      await vi.advanceTimersByTimeAsync(300);
      fixture.detectChanges();

      expect(component.searchQuery()).toBe('');
      expect(emitSpy).toHaveBeenCalledWith('');
    });
  });

  describe('Filter Toggle', () => {
    it('should render filter button', () => {
      const filterBtn = fixture.debugElement.query(By.css('.filter-btn'));
      expect(filterBtn).toBeTruthy();
      expect(filterBtn.nativeElement.textContent).toContain('Filters');
    });

    it('should emit filtersToggled on button click', () => {
      const emitSpy = vi.spyOn(component.filtersToggled, 'emit');

      const filterBtn = fixture.debugElement.query(By.css('.filter-btn'));
      filterBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalled();
    });

    it('should toggle showFilters state', () => {
      expect(component.showFilters()).toBe(false);

      component.toggleFilters();
      expect(component.showFilters()).toBe(true);

      component.toggleFilters();
      expect(component.showFilters()).toBe(false);
    });
  });

  describe('Active Filter Count', () => {
    it('should show badge when filters are active', () => {
      fixture.componentRef.setInput('filters', {
        status: 'active',
        gender: 'male',
      } as StudentListFilter);
      fixture.detectChanges();

      const badge = fixture.debugElement.query(By.css('msls-badge'));
      expect(badge).toBeTruthy();
      expect(badge.nativeElement.textContent.trim()).toBe('2');
    });

    it('should not show badge when no filters are active', () => {
      fixture.componentRef.setInput('filters', {} as StudentListFilter);
      fixture.detectChanges();

      const badge = fixture.debugElement.query(By.css('msls-badge'));
      expect(badge).toBeFalsy();
    });

    it('should count all active filter types', () => {
      fixture.componentRef.setInput('filters', {
        classId: 'class-1',
        sectionId: 'section-1',
        status: 'active',
        gender: 'male',
        admissionFrom: '2026-01-01',
        admissionTo: '2026-12-31',
      } as StudentListFilter);
      fixture.detectChanges();

      expect(component.activeFilterCount()).toBe(6);
    });
  });

  describe('Initial State', () => {
    it('should initialize search from filters input', async () => {
      // Create a new component with pre-set filters
      const newFixture = TestBed.createComponent(StudentSearchComponent);
      newFixture.componentRef.setInput('filters', { search: 'initial search' } as StudentListFilter);
      newFixture.detectChanges();

      // Need to trigger ngOnInit
      newFixture.componentInstance.ngOnInit();
      newFixture.detectChanges();

      expect(newFixture.componentInstance.searchQuery()).toBe('initial search');
      newFixture.destroy();
    });
  });
});
