import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { Component } from '@angular/core';
import { MslsDropdownComponent } from './dropdown.component';

@Component({
  standalone: true,
  imports: [MslsDropdownComponent],
  template: `
    <msls-dropdown [trigger]="trigger" [position]="position" (opened)="onOpened()" (closed)="onClosed()">
      <button dropdown-trigger>Open Menu</button>
      <div dropdown-content>
        <a href="#">Item 1</a>
        <a href="#">Item 2</a>
      </div>
    </msls-dropdown>
  `,
})
class TestHostComponent {
  trigger: 'click' | 'hover' = 'click';
  position: 'bottom-start' | 'bottom-end' | 'top-start' | 'top-end' = 'bottom-start';
  openedCalled = false;
  closedCalled = false;

  onOpened(): void {
    this.openedCalled = true;
  }

  onClosed(): void {
    this.closedCalled = true;
  }
}

describe('MslsDropdownComponent', () => {
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

  it('should be closed by default', () => {
    const contentElement = fixture.nativeElement.querySelector('.dropdown__content');
    expect(contentElement).toBeFalsy();
  });

  it('should open on trigger click when trigger is click', () => {
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const contentElement = fixture.nativeElement.querySelector('.dropdown__content');
    expect(contentElement).toBeTruthy();
  });

  it('should close on second trigger click', () => {
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');

    // Open
    triggerElement.click();
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.dropdown__content')).toBeTruthy();

    // Close
    triggerElement.click();
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.dropdown__content')).toBeFalsy();
  });

  it('should emit opened event when opening', () => {
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');
    triggerElement.click();
    fixture.detectChanges();

    expect(component.openedCalled).toBeTruthy();
  });

  it('should emit closed event when closing', () => {
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');

    // Open first
    triggerElement.click();
    fixture.detectChanges();

    // Close
    triggerElement.click();
    fixture.detectChanges();

    expect(component.closedCalled).toBeTruthy();
  });

  it('should apply correct position class', () => {
    component.position = 'bottom-end';
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const contentElement = fixture.nativeElement.querySelector('.dropdown__content');
    expect(contentElement.classList.contains('dropdown__content--bottom-end')).toBeTruthy();
  });

  it('should close on escape key', () => {
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');
    triggerElement.click();
    fixture.detectChanges();

    // Simulate escape key
    const event = new KeyboardEvent('keydown', { key: 'Escape' });
    document.dispatchEvent(event);
    fixture.detectChanges();

    const contentElement = fixture.nativeElement.querySelector('.dropdown__content');
    expect(contentElement).toBeFalsy();
  });

  it('should project trigger content', () => {
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');
    expect(triggerElement.textContent).toContain('Open Menu');
  });

  it('should project dropdown content', () => {
    const triggerElement = fixture.nativeElement.querySelector('.dropdown__trigger');
    triggerElement.click();
    fixture.detectChanges();

    const contentElement = fixture.nativeElement.querySelector('.dropdown__content');
    expect(contentElement.textContent).toContain('Item 1');
    expect(contentElement.textContent).toContain('Item 2');
  });
});
