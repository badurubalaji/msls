import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { Component, signal } from '@angular/core';
import { MslsModalComponent } from './modal.component';

@Component({
  standalone: true,
  imports: [MslsModalComponent],
  template: `
    <msls-modal
      [isOpen]="isOpen()"
      [size]="size"
      [closeOnBackdrop]="closeOnBackdrop"
      [closeOnEscape]="closeOnEscape"
      [showCloseButton]="showCloseButton"
      [title]="title"
      (closed)="onClosed()"
    >
      <ng-container modal-header>Header Content</ng-container>
      <ng-container modal-body>
        <input type="text" id="input1" />
        <button id="button1">Button 1</button>
      </ng-container>
      <ng-container modal-footer>
        <button id="button2">Footer Button</button>
      </ng-container>
    </msls-modal>
  `,
})
class TestHostComponent {
  isOpen = signal(false);
  size: 'sm' | 'md' | 'lg' | 'xl' | 'full' = 'md';
  closeOnBackdrop = true;
  closeOnEscape = true;
  showCloseButton = true;
  title = '';
  closedCalled = false;

  onClosed(): void {
    this.closedCalled = true;
    this.isOpen.set(false);
  }
}

describe('MslsModalComponent', () => {
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

  it('should not display modal when isOpen is false', () => {
    const modalElement = fixture.nativeElement.querySelector('.modal');
    expect(modalElement).toBeFalsy();
  });

  it('should display modal when isOpen is true', () => {
    component.isOpen.set(true);
    fixture.detectChanges();

    const modalElement = fixture.nativeElement.querySelector('.modal');
    expect(modalElement).toBeTruthy();
  });

  it('should apply correct size class', () => {
    component.size = 'lg';
    component.isOpen.set(true);
    fixture.detectChanges();

    const panelElement = fixture.nativeElement.querySelector('.modal__panel');
    expect(panelElement.classList.contains('modal__panel--lg')).toBeTruthy();
  });

  it('should display title when provided', () => {
    component.title = 'Test Title';
    component.isOpen.set(true);
    fixture.detectChanges();

    const titleElement = fixture.nativeElement.querySelector('.modal__title h2');
    expect(titleElement.textContent).toBe('Test Title');
  });

  it('should display close button when showCloseButton is true', () => {
    component.isOpen.set(true);
    fixture.detectChanges();

    const closeButton = fixture.nativeElement.querySelector('.modal__close');
    expect(closeButton).toBeTruthy();
  });

  it('should hide close button when showCloseButton is false', () => {
    component.showCloseButton = false;
    component.isOpen.set(true);
    fixture.detectChanges();

    const closeButton = fixture.nativeElement.querySelector('.modal__close');
    expect(closeButton).toBeFalsy();
  });

  it('should emit closed when close button is clicked', () => {
    component.isOpen.set(true);
    fixture.detectChanges();

    const closeButton = fixture.nativeElement.querySelector('.modal__close');
    closeButton.click();
    fixture.detectChanges();

    expect(component.closedCalled).toBeTruthy();
  });

  it('should emit closed when backdrop is clicked and closeOnBackdrop is true', () => {
    component.isOpen.set(true);
    fixture.detectChanges();

    const modalElement = fixture.nativeElement.querySelector('.modal');
    modalElement.click();
    fixture.detectChanges();

    expect(component.closedCalled).toBeTruthy();
  });

  it('should not emit closed when backdrop is clicked and closeOnBackdrop is false', () => {
    component.closeOnBackdrop = false;
    component.isOpen.set(true);
    fixture.detectChanges();

    const modalElement = fixture.nativeElement.querySelector('.modal');
    modalElement.click();
    fixture.detectChanges();

    expect(component.closedCalled).toBeFalsy();
  });

  it('should project header content', () => {
    component.isOpen.set(true);
    fixture.detectChanges();

    const headerElement = fixture.nativeElement.querySelector('.modal__title');
    expect(headerElement.textContent).toContain('Header Content');
  });

  it('should project body content', () => {
    component.isOpen.set(true);
    fixture.detectChanges();

    const bodyElement = fixture.nativeElement.querySelector('.modal__body');
    expect(bodyElement.querySelector('#input1')).toBeTruthy();
    expect(bodyElement.querySelector('#button1')).toBeTruthy();
  });

  it('should project footer content', () => {
    component.isOpen.set(true);
    fixture.detectChanges();

    const footerElement = fixture.nativeElement.querySelector('.modal__footer');
    expect(footerElement.querySelector('#button2')).toBeTruthy();
  });
});
