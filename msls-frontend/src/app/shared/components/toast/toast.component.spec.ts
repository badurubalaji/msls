import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { MslsToastComponent } from './toast.component';
import { ToastService } from '../../services/toast.service';

describe('MslsToastComponent', () => {
  let component: MslsToastComponent;
  let fixture: ComponentFixture<MslsToastComponent>;
  let toastService: ToastService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MslsToastComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MslsToastComponent);
    component = fixture.componentInstance;
    toastService = TestBed.inject(ToastService);
    fixture.detectChanges();
  });

  afterEach(() => {
    toastService.dismissAll();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display no toasts initially', () => {
    const toastElements = fixture.nativeElement.querySelectorAll('.toast');
    expect(toastElements.length).toBe(0);
  });

  it('should display a toast when added via service', () => {
    toastService.success('Test message');
    fixture.detectChanges();

    const toastElements = fixture.nativeElement.querySelectorAll('.toast');
    expect(toastElements.length).toBe(1);
    expect(toastElements[0].textContent).toContain('Test message');
  });

  it('should display success variant with correct class', () => {
    toastService.success('Success message');
    fixture.detectChanges();

    const toastElement = fixture.nativeElement.querySelector('.toast');
    expect(toastElement.classList.contains('toast--success')).toBeTruthy();
  });

  it('should display error variant with correct class', () => {
    toastService.error('Error message');
    fixture.detectChanges();

    const toastElement = fixture.nativeElement.querySelector('.toast');
    expect(toastElement.classList.contains('toast--error')).toBeTruthy();
  });

  it('should display warning variant with correct class', () => {
    toastService.warning('Warning message');
    fixture.detectChanges();

    const toastElement = fixture.nativeElement.querySelector('.toast');
    expect(toastElement.classList.contains('toast--warning')).toBeTruthy();
  });

  it('should display info variant with correct class', () => {
    toastService.info('Info message');
    fixture.detectChanges();

    const toastElement = fixture.nativeElement.querySelector('.toast');
    expect(toastElement.classList.contains('toast--info')).toBeTruthy();
  });

  it('should dismiss toast when dismiss button is clicked', () => {
    toastService.success('Test message');
    fixture.detectChanges();

    const dismissButton = fixture.nativeElement.querySelector('.toast__dismiss');
    dismissButton.click();
    fixture.detectChanges();

    const toastElements = fixture.nativeElement.querySelectorAll('.toast');
    expect(toastElements.length).toBe(0);
  });

  it('should display multiple toasts', () => {
    toastService.success('Message 1');
    toastService.error('Message 2');
    toastService.info('Message 3');
    fixture.detectChanges();

    const toastElements = fixture.nativeElement.querySelectorAll('.toast');
    expect(toastElements.length).toBe(3);
  });

  it('should auto-dismiss toast after duration', fakeAsync(() => {
    toastService.success('Test message', 1000);
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelectorAll('.toast').length).toBe(1);

    tick(1000);
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelectorAll('.toast').length).toBe(0);
  }));
});
