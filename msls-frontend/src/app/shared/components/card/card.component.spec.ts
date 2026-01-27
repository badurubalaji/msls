import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Component } from '@angular/core';
import { MslsCardComponent } from './card.component';

@Component({
  standalone: true,
  imports: [MslsCardComponent],
  template: `
    <msls-card [variant]="variant" [padding]="padding" [hoverable]="hoverable" [clickable]="clickable">
      <ng-container card-header>Header Content</ng-container>
      <ng-container card-body>Body Content</ng-container>
      <ng-container card-footer>Footer Content</ng-container>
    </msls-card>
  `,
})
class TestHostComponent {
  variant: 'default' | 'elevated' | 'outlined' = 'default';
  padding: 'none' | 'sm' | 'md' | 'lg' = 'md';
  hoverable = false;
  clickable = false;
}

describe('MslsCardComponent', () => {
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

  it('should apply default variant class', () => {
    const cardElement = fixture.nativeElement.querySelector('.card');
    expect(cardElement.classList.contains('card--default')).toBeTruthy();
  });

  it('should apply elevated variant class', () => {
    component.variant = 'elevated';
    fixture.detectChanges();
    const cardElement = fixture.nativeElement.querySelector('.card');
    expect(cardElement.classList.contains('card--elevated')).toBeTruthy();
  });

  it('should apply outlined variant class', () => {
    component.variant = 'outlined';
    fixture.detectChanges();
    const cardElement = fixture.nativeElement.querySelector('.card');
    expect(cardElement.classList.contains('card--outlined')).toBeTruthy();
  });

  it('should apply padding class', () => {
    component.padding = 'lg';
    fixture.detectChanges();
    const cardElement = fixture.nativeElement.querySelector('.card');
    expect(cardElement.classList.contains('card--padding-lg')).toBeTruthy();
  });

  it('should apply hoverable class when hoverable is true', () => {
    component.hoverable = true;
    fixture.detectChanges();
    const cardElement = fixture.nativeElement.querySelector('.card');
    expect(cardElement.classList.contains('card--hoverable')).toBeTruthy();
  });

  it('should apply clickable class when clickable is true', () => {
    component.clickable = true;
    fixture.detectChanges();
    const cardElement = fixture.nativeElement.querySelector('.card');
    expect(cardElement.classList.contains('card--clickable')).toBeTruthy();
  });

  it('should project header content', () => {
    const headerElement = fixture.nativeElement.querySelector('.card__header');
    expect(headerElement.textContent).toContain('Header Content');
  });

  it('should project body content', () => {
    const bodyElement = fixture.nativeElement.querySelector('.card__body');
    expect(bodyElement.textContent).toContain('Body Content');
  });

  it('should project footer content', () => {
    const footerElement = fixture.nativeElement.querySelector('.card__footer');
    expect(footerElement.textContent).toContain('Footer Content');
  });
});
