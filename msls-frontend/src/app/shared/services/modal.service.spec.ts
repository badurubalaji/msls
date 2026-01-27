import { TestBed } from '@angular/core/testing';
import { Component } from '@angular/core';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { firstValueFrom } from 'rxjs';
import { ModalService } from './modal.service';

@Component({
  standalone: true,
  template: '<div>Test Modal Content</div>',
})
class TestModalComponent {
  data?: string;
}

describe('ModalService', () => {
  let service: ModalService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModalService);
  });

  afterEach(() => {
    service.closeAll();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should start with no open modals', () => {
    expect(service.hasOpenModals()).toBeFalsy();
    expect(service.openModals$().length).toBe(0);
  });

  it('should open a modal and return a ModalRef', () => {
    const modalRef = service.open(TestModalComponent);

    expect(modalRef).toBeTruthy();
    expect(modalRef.instance).toBeTruthy();
    expect(service.hasOpenModals()).toBeTruthy();
    expect(service.openModals$().length).toBe(1);
  });

  it('should pass data to the modal component', () => {
    const modalRef = service.open(TestModalComponent, {
      data: { data: 'test data' },
    });

    expect((modalRef.instance as TestModalComponent).data).toBe('test data');
  });

  it('should close a modal when close() is called on ModalRef', async () => {
    const modalRef = service.open(TestModalComponent);

    const resultPromise = firstValueFrom(modalRef.afterClosed$);
    modalRef.close('test result');

    const result = await resultPromise;
    expect(result).toBe('test result');
    expect(service.hasOpenModals()).toBeFalsy();
  });

  it('should close all modals when closeAll() is called', () => {
    service.open(TestModalComponent);
    service.open(TestModalComponent);
    service.open(TestModalComponent);

    expect(service.openModals$().length).toBe(3);

    service.closeAll();

    expect(service.openModals$().length).toBe(0);
    expect(service.hasOpenModals()).toBeFalsy();
  });

  it('should get the topmost modal', () => {
    service.open(TestModalComponent);
    const secondModalRef = service.open(TestModalComponent);

    const topModal = service.getTopModal();
    expect(topModal).toBeTruthy();
    // The second opened modal should be on top
    expect(topModal?.componentRef.instance).toBe(secondModalRef.instance);
  });

  it('should use default config values', () => {
    service.open(TestModalComponent);

    const config = service.getModalConfig(service.openModals$()[0].id);
    expect(config?.size).toBe('md');
    expect(config?.closeOnBackdrop).toBe(true);
    expect(config?.closeOnEscape).toBe(true);
    expect(config?.showCloseButton).toBe(true);
  });

  it('should merge custom config with defaults', () => {
    service.open(TestModalComponent, {
      size: 'lg',
      closeOnBackdrop: false,
    });

    const config = service.getModalConfig(service.openModals$()[0].id);
    expect(config?.size).toBe('lg');
    expect(config?.closeOnBackdrop).toBe(false);
    expect(config?.closeOnEscape).toBe(true); // default
  });
});
