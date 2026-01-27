import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { ToastService } from './toast.service';

describe('ToastService', () => {
  let service: ToastService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ToastService);
  });

  afterEach(() => {
    service.dismissAll();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should start with no toasts', () => {
    expect(service.toasts$().length).toBe(0);
  });

  it('should add a success toast', () => {
    service.success('Success message');
    expect(service.toasts$().length).toBe(1);
    expect(service.toasts$()[0].variant).toBe('success');
    expect(service.toasts$()[0].message).toBe('Success message');
  });

  it('should add an error toast', () => {
    service.error('Error message');
    expect(service.toasts$().length).toBe(1);
    expect(service.toasts$()[0].variant).toBe('error');
    expect(service.toasts$()[0].message).toBe('Error message');
  });

  it('should add a warning toast', () => {
    service.warning('Warning message');
    expect(service.toasts$().length).toBe(1);
    expect(service.toasts$()[0].variant).toBe('warning');
    expect(service.toasts$()[0].message).toBe('Warning message');
  });

  it('should add an info toast', () => {
    service.info('Info message');
    expect(service.toasts$().length).toBe(1);
    expect(service.toasts$()[0].variant).toBe('info');
    expect(service.toasts$()[0].message).toBe('Info message');
  });

  it('should return toast ID when showing toast', () => {
    const id = service.success('Test');
    expect(id).toBeTruthy();
    expect(id.startsWith('toast-')).toBeTruthy();
  });

  it('should dismiss a specific toast by ID', () => {
    const id = service.success('Test message');
    expect(service.toasts$().length).toBe(1);

    service.dismiss(id);
    expect(service.toasts$().length).toBe(0);
  });

  it('should dismiss all toasts', () => {
    service.success('Message 1');
    service.error('Message 2');
    service.info('Message 3');
    expect(service.toasts$().length).toBe(3);

    service.dismissAll();
    expect(service.toasts$().length).toBe(0);
  });

  it('should auto-dismiss toast after duration', fakeAsync(() => {
    service.success('Test message', 1000);
    expect(service.toasts$().length).toBe(1);

    tick(1000);
    expect(service.toasts$().length).toBe(0);
  }));

  it('should allow custom configuration via show()', () => {
    service.show({
      message: 'Custom message',
      variant: 'warning',
      duration: 10000,
      dismissible: false,
    });

    const toast = service.toasts$()[0];
    expect(toast.message).toBe('Custom message');
    expect(toast.variant).toBe('warning');
    expect(toast.duration).toBe(10000);
    expect(toast.dismissible).toBe(false);
  });

  it('should use default values when not specified', () => {
    service.show({ message: 'Test' });

    const toast = service.toasts$()[0];
    expect(toast.variant).toBe('info');
    expect(toast.duration).toBe(5000);
    expect(toast.dismissible).toBe(true);
  });

  it('should handle multiple toasts with different durations', fakeAsync(() => {
    service.success('Short', 500);
    service.error('Long', 2000);

    expect(service.toasts$().length).toBe(2);

    tick(500);
    expect(service.toasts$().length).toBe(1);
    expect(service.toasts$()[0].variant).toBe('error');

    tick(1500);
    expect(service.toasts$().length).toBe(0);
  }));
});
