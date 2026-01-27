/**
 * Export Dialog Component Tests
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { ExportDialogComponent } from './export-dialog.component';
import { DEFAULT_EXPORT_COLUMNS } from '../../models/student.model';

describe('ExportDialogComponent', () => {
  let component: ExportDialogComponent;
  let fixture: ComponentFixture<ExportDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ExportDialogComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ExportDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('Dialog Visibility', () => {
    it('should not render dialog when isOpen is false', () => {
      fixture.componentRef.setInput('isOpen', false);
      fixture.detectChanges();

      const overlay = fixture.debugElement.query(By.css('.modal-overlay'));
      expect(overlay).toBeFalsy();
    });

    it('should render dialog when isOpen is true', () => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1']);
      fixture.detectChanges();

      const overlay = fixture.debugElement.query(By.css('.modal-overlay'));
      expect(overlay).toBeTruthy();
    });

    it('should render modal with header', () => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1']);
      fixture.detectChanges();

      const header = fixture.debugElement.query(By.css('.modal__header'));
      expect(header).toBeTruthy();
      expect(header.nativeElement.textContent).toContain('Export Students');
    });
  });

  describe('Student Count Display', () => {
    it('should display correct student count (plural)', () => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1', 'id-2', 'id-3']);
      fixture.detectChanges();

      const exportInfo = fixture.debugElement.query(By.css('.export-info'));
      expect(exportInfo.nativeElement.textContent).toContain('3');
      expect(exportInfo.nativeElement.textContent).toContain('students');
    });

    it('should display correct student count (singular)', () => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1']);
      fixture.detectChanges();

      const exportInfo = fixture.debugElement.query(By.css('.export-info'));
      expect(exportInfo.nativeElement.textContent).toContain('1');
      expect(exportInfo.nativeElement.textContent).toContain('student');
    });

    it('should compute studentCount from studentIds length', () => {
      fixture.componentRef.setInput('studentIds', ['a', 'b', 'c', 'd', 'e']);
      fixture.detectChanges();

      expect(component.studentCount()).toBe(5);
    });
  });

  describe('Format Selection', () => {
    beforeEach(() => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1']);
      fixture.detectChanges();
    });

    it('should default to xlsx format', () => {
      expect(component.selectedFormat()).toBe('xlsx');
    });

    it('should render format options', () => {
      const formatOptions = fixture.debugElement.queryAll(By.css('.format-option'));
      expect(formatOptions.length).toBe(2);
    });

    it('should show xlsx option as selected by default', () => {
      const xlsxOption = fixture.debugElement.queryAll(By.css('.format-option'))[0];
      expect(xlsxOption.nativeElement.classList.contains('selected')).toBe(true);
    });

    it('should change format to csv on selection', () => {
      const csvRadio = fixture.debugElement.queryAll(By.css('.format-option input'))[1];
      csvRadio.nativeElement.click();
      fixture.detectChanges();

      expect(component.selectedFormat()).toBe('csv');
    });

    it('should update selected class when format changes', () => {
      component.selectedFormat.set('csv');
      fixture.detectChanges();

      const xlsxOption = fixture.debugElement.queryAll(By.css('.format-option'))[0];
      const csvOption = fixture.debugElement.queryAll(By.css('.format-option'))[1];

      expect(xlsxOption.nativeElement.classList.contains('selected')).toBe(false);
      expect(csvOption.nativeElement.classList.contains('selected')).toBe(true);
    });
  });

  describe('Column Selection', () => {
    beforeEach(() => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1']);
      fixture.detectChanges();
    });

    it('should render all columns from DEFAULT_EXPORT_COLUMNS', () => {
      const columnOptions = fixture.debugElement.queryAll(By.css('.column-option'));
      expect(columnOptions.length).toBe(DEFAULT_EXPORT_COLUMNS.length);
    });

    it('should show default selected columns as selected', () => {
      const defaultSelectedCount = DEFAULT_EXPORT_COLUMNS.filter((c) => c.selected).length;
      const selectedOptions = fixture.debugElement.queryAll(By.css('.column-option.selected'));
      expect(selectedOptions.length).toBe(defaultSelectedCount);
    });

    it('should toggle column selection on click', () => {
      const firstColumn = component.columns()[0];
      const initialSelected = firstColumn.selected;

      component.toggleColumn(firstColumn.key);
      fixture.detectChanges();

      expect(component.columns()[0].selected).toBe(!initialSelected);
    });

    it('should render select all button', () => {
      const selectAllBtn = fixture.debugElement.query(By.css('.select-all-btn'));
      expect(selectAllBtn).toBeTruthy();
    });

    it('should select all columns on select all click', () => {
      // First deselect some
      component.columns.update((cols) => cols.map((c, i) => ({ ...c, selected: i < 3 })));
      fixture.detectChanges();

      const selectAllBtn = fixture.debugElement.query(By.css('.select-all-btn'));
      selectAllBtn.nativeElement.click();
      fixture.detectChanges();

      expect(component.columns().every((c) => c.selected)).toBe(true);
    });

    it('should deselect all columns on deselect all click', () => {
      // First select all
      component.columns.update((cols) => cols.map((c) => ({ ...c, selected: true })));
      fixture.detectChanges();

      const selectAllBtn = fixture.debugElement.query(By.css('.select-all-btn'));
      expect(selectAllBtn.nativeElement.textContent).toContain('Deselect All');

      selectAllBtn.nativeElement.click();
      fixture.detectChanges();

      expect(component.columns().every((c) => !c.selected)).toBe(true);
    });
  });

  describe('allColumnsSelected', () => {
    it('should return true when all columns are selected', () => {
      component.columns.update((cols) => cols.map((c) => ({ ...c, selected: true })));
      expect(component.allColumnsSelected()).toBe(true);
    });

    it('should return false when some columns are not selected', () => {
      component.columns.update((cols) =>
        cols.map((c, i) => ({ ...c, selected: i !== 0 }))
      );
      expect(component.allColumnsSelected()).toBe(false);
    });

    it('should return false when no columns are selected', () => {
      component.columns.update((cols) => cols.map((c) => ({ ...c, selected: false })));
      expect(component.allColumnsSelected()).toBe(false);
    });
  });

  describe('canExport', () => {
    it('should return true when at least one column is selected', () => {
      component.columns.update((cols) =>
        cols.map((c, i) => ({ ...c, selected: i === 0 }))
      );
      expect(component.canExport()).toBe(true);
    });

    it('should return false when no columns are selected', () => {
      component.columns.update((cols) => cols.map((c) => ({ ...c, selected: false })));
      expect(component.canExport()).toBe(false);
    });
  });

  describe('Cancel', () => {
    beforeEach(() => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1']);
      fixture.detectChanges();
    });

    it('should emit cancelled on cancel button click', () => {
      const emitSpy = vi.spyOn(component.cancelled, 'emit');

      const cancelBtn = fixture.debugElement.query(By.css('.btn--secondary'));
      cancelBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalled();
    });

    it('should emit cancelled on close button click', () => {
      const emitSpy = vi.spyOn(component.cancelled, 'emit');

      const closeBtn = fixture.debugElement.query(By.css('.modal__close'));
      closeBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalled();
    });

    it('should emit cancelled on overlay click', () => {
      const emitSpy = vi.spyOn(component.cancelled, 'emit');

      const overlay = fixture.debugElement.query(By.css('.modal-overlay'));
      overlay.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalled();
    });

    it('should not emit cancelled on modal click (stopPropagation)', () => {
      const emitSpy = vi.spyOn(component.cancelled, 'emit');

      const modal = fixture.debugElement.query(By.css('.modal'));
      modal.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).not.toHaveBeenCalled();
    });
  });

  describe('Export', () => {
    const testStudentIds = ['student-1', 'student-2'];

    beforeEach(() => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', testStudentIds);
      fixture.detectChanges();
    });

    it('should render export button', () => {
      const exportBtn = fixture.debugElement.query(By.css('.btn--primary'));
      expect(exportBtn).toBeTruthy();
      expect(exportBtn.nativeElement.textContent).toContain('Export');
    });

    it('should disable export button when no columns selected', () => {
      component.columns.update((cols) => cols.map((c) => ({ ...c, selected: false })));
      fixture.detectChanges();

      const exportBtn = fixture.debugElement.query(By.css('.btn--primary'));
      expect(exportBtn.nativeElement.disabled).toBe(true);
    });

    it('should enable export button when columns are selected', () => {
      component.columns.update((cols) =>
        cols.map((c, i) => ({ ...c, selected: i === 0 }))
      );
      fixture.detectChanges();

      const exportBtn = fixture.debugElement.query(By.css('.btn--primary'));
      expect(exportBtn.nativeElement.disabled).toBe(false);
    });

    it('should emit exported with correct request on export click', () => {
      const emitSpy = vi.spyOn(component.exported, 'emit');

      // Select specific columns
      component.columns.update((cols) =>
        cols.map((c) => ({
          ...c,
          selected: ['admission_number', 'first_name', 'last_name'].includes(c.key),
        }))
      );
      component.selectedFormat.set('csv');
      fixture.detectChanges();

      const exportBtn = fixture.debugElement.query(By.css('.btn--primary'));
      exportBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalledWith({
        studentIds: testStudentIds,
        format: 'csv',
        columns: ['admission_number', 'first_name', 'last_name'],
      });
    });

    it('should include only selected columns in export request', () => {
      const emitSpy = vi.spyOn(component.exported, 'emit');

      // Select only first column
      component.columns.update((cols) =>
        cols.map((c, i) => ({ ...c, selected: i === 0 }))
      );
      fixture.detectChanges();

      component.onExport();

      expect(emitSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          columns: [DEFAULT_EXPORT_COLUMNS[0].key],
        })
      );
    });
  });

  describe('Export Button Label', () => {
    it('should show singular label for one student', () => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1']);
      fixture.detectChanges();

      const exportBtn = fixture.debugElement.query(By.css('.btn--primary'));
      expect(exportBtn.nativeElement.textContent).toContain('Export 1 Student');
    });

    it('should show plural label for multiple students', () => {
      fixture.componentRef.setInput('isOpen', true);
      fixture.componentRef.setInput('studentIds', ['id-1', 'id-2', 'id-3']);
      fixture.detectChanges();

      const exportBtn = fixture.debugElement.query(By.css('.btn--primary'));
      expect(exportBtn.nativeElement.textContent).toContain('Export 3 Students');
    });
  });
});
