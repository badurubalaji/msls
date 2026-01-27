/**
 * Bulk Actions Component Tests
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { BulkActionsComponent, BulkActionType } from './bulk-actions.component';

describe('BulkActionsComponent', () => {
  let component: BulkActionsComponent;
  let fixture: ComponentFixture<BulkActionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BulkActionsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(BulkActionsComponent);
    component = fixture.componentInstance;
    // Set required inputs
    fixture.componentRef.setInput('selectedCount', 0);
    fixture.componentRef.setInput('selectedIds', []);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('Visibility', () => {
    it('should not render action bar when no students selected', () => {
      fixture.componentRef.setInput('selectedCount', 0);
      fixture.componentRef.setInput('selectedIds', []);
      fixture.detectChanges();

      const actionBar = fixture.debugElement.query(By.css('.bulk-actions-bar'));
      expect(actionBar).toBeFalsy();
    });

    it('should render action bar when students are selected', () => {
      fixture.componentRef.setInput('selectedCount', 5);
      fixture.componentRef.setInput('selectedIds', ['id-1', 'id-2', 'id-3', 'id-4', 'id-5']);
      fixture.detectChanges();

      const actionBar = fixture.debugElement.query(By.css('.bulk-actions-bar'));
      expect(actionBar).toBeTruthy();
    });
  });

  describe('Selection Info', () => {
    it('should display correct selection count', () => {
      fixture.componentRef.setInput('selectedCount', 10);
      fixture.componentRef.setInput('selectedIds', Array(10).fill('id'));
      fixture.detectChanges();

      const count = fixture.debugElement.query(By.css('.selection-info .count'));
      expect(count.nativeElement.textContent).toBe('10');
    });

    it('should display "students selected" label', () => {
      fixture.componentRef.setInput('selectedCount', 3);
      fixture.componentRef.setInput('selectedIds', ['id-1', 'id-2', 'id-3']);
      fixture.detectChanges();

      const label = fixture.debugElement.query(By.css('.selection-info .label'));
      expect(label.nativeElement.textContent).toBe('students selected');
    });
  });

  describe('Action Buttons', () => {
    beforeEach(() => {
      fixture.componentRef.setInput('selectedCount', 2);
      fixture.componentRef.setInput('selectedIds', ['id-1', 'id-2']);
      fixture.detectChanges();
    });

    it('should render all action buttons', () => {
      const actionBtns = fixture.debugElement.queryAll(By.css('.action-btn'));
      expect(actionBtns.length).toBe(4);
    });

    it('should have SMS button disabled', () => {
      const smsBtn = fixture.debugElement.queryAll(By.css('.action-btn'))[0];
      expect(smsBtn.nativeElement.disabled).toBe(true);
      expect(smsBtn.nativeElement.classList.contains('action-btn--disabled')).toBe(true);
    });

    it('should have Email button disabled', () => {
      const emailBtn = fixture.debugElement.queryAll(By.css('.action-btn'))[1];
      expect(emailBtn.nativeElement.disabled).toBe(true);
      expect(emailBtn.nativeElement.classList.contains('action-btn--disabled')).toBe(true);
    });

    it('should have Update Status button enabled', () => {
      const statusBtn = fixture.debugElement.queryAll(By.css('.action-btn'))[2];
      expect(statusBtn.nativeElement.disabled).toBeFalsy();
      expect(statusBtn.nativeElement.textContent).toContain('Update Status');
    });

    it('should have Export button with primary style', () => {
      const exportBtn = fixture.debugElement.queryAll(By.css('.action-btn'))[3];
      expect(exportBtn.nativeElement.disabled).toBeFalsy();
      expect(exportBtn.nativeElement.classList.contains('action-btn--primary')).toBe(true);
      expect(exportBtn.nativeElement.textContent).toContain('Export');
    });
  });

  describe('Action Events', () => {
    const testIds = ['student-1', 'student-2', 'student-3'];

    beforeEach(() => {
      fixture.componentRef.setInput('selectedCount', 3);
      fixture.componentRef.setInput('selectedIds', testIds);
      fixture.detectChanges();
    });

    it('should emit status action with correct payload', () => {
      const emitSpy = vi.spyOn(component.action, 'emit');
      const statusBtn = fixture.debugElement.queryAll(By.css('.action-btn'))[2];

      statusBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalledWith({
        type: 'status',
        ids: testIds,
      });
    });

    it('should emit export action with correct payload', () => {
      const emitSpy = vi.spyOn(component.action, 'emit');
      const exportBtn = fixture.debugElement.queryAll(By.css('.action-btn'))[3];

      exportBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalledWith({
        type: 'export',
        ids: testIds,
      });
    });

    it('should emit sms action when button clicked (even if disabled)', () => {
      const emitSpy = vi.spyOn(component.action, 'emit');

      // Call the method directly since button is disabled
      component.onAction('sms');

      expect(emitSpy).toHaveBeenCalledWith({
        type: 'sms',
        ids: testIds,
      });
    });

    it('should emit email action when button clicked (even if disabled)', () => {
      const emitSpy = vi.spyOn(component.action, 'emit');

      // Call the method directly since button is disabled
      component.onAction('email');

      expect(emitSpy).toHaveBeenCalledWith({
        type: 'email',
        ids: testIds,
      });
    });
  });

  describe('Clear Selection', () => {
    it('should render close button', () => {
      fixture.componentRef.setInput('selectedCount', 1);
      fixture.componentRef.setInput('selectedIds', ['id-1']);
      fixture.detectChanges();

      const closeBtn = fixture.debugElement.query(By.css('.close-btn'));
      expect(closeBtn).toBeTruthy();
    });

    it('should emit cleared event on close click', () => {
      const emitSpy = vi.spyOn(component.cleared, 'emit');
      fixture.componentRef.setInput('selectedCount', 1);
      fixture.componentRef.setInput('selectedIds', ['id-1']);
      fixture.detectChanges();

      const closeBtn = fixture.debugElement.query(By.css('.close-btn'));
      closeBtn.nativeElement.click();
      fixture.detectChanges();

      expect(emitSpy).toHaveBeenCalled();
    });
  });

  describe('BulkActionType', () => {
    it('should correctly type action types', () => {
      const validTypes: BulkActionType[] = ['sms', 'email', 'status', 'export'];
      fixture.componentRef.setInput('selectedCount', 1);
      fixture.componentRef.setInput('selectedIds', ['id-1']);
      fixture.detectChanges();

      validTypes.forEach((type) => {
        component.onAction(type);
        // No type error means the type is correct
        expect(true).toBe(true);
      });
    });
  });
});
