import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { MslsButtonComponent, ButtonVariant, ButtonSize } from './button.component';

describe('MslsButtonComponent', () => {
  let component: MslsButtonComponent;
  let fixture: ComponentFixture<MslsButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MslsButtonComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MslsButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('default values', () => {
    it('should have primary variant by default', () => {
      expect(component.variant()).toBe('primary');
    });

    it('should have md size by default', () => {
      expect(component.size()).toBe('md');
    });

    it('should have button type by default', () => {
      expect(component.type()).toBe('button');
    });

    it('should not be loading by default', () => {
      expect(component.loading()).toBe(false);
    });

    it('should not be disabled by default', () => {
      expect(component.disabled()).toBe(false);
    });
  });

  describe('variants', () => {
    const variants: ButtonVariant[] = ['primary', 'secondary', 'danger', 'ghost', 'outline'];

    variants.forEach((variant) => {
      it(`should apply ${variant} variant class`, () => {
        fixture.componentRef.setInput('variant', variant);
        fixture.detectChanges();

        const button = fixture.debugElement.query(By.css('button'));
        expect(button.nativeElement.classList.contains(`msls-button--${variant}`)).toBe(true);
      });
    });
  });

  describe('sizes', () => {
    const sizes: ButtonSize[] = ['sm', 'md', 'lg'];

    sizes.forEach((size) => {
      it(`should apply ${size} size class`, () => {
        fixture.componentRef.setInput('size', size);
        fixture.detectChanges();

        const button = fixture.debugElement.query(By.css('button'));
        expect(button.nativeElement.classList.contains(`msls-button--${size}`)).toBe(true);
      });
    });
  });

  describe('loading state', () => {
    it('should show spinner when loading', () => {
      fixture.componentRef.setInput('loading', true);
      fixture.detectChanges();

      const spinner = fixture.debugElement.query(By.css('.msls-button__spinner'));
      expect(spinner).toBeTruthy();
    });

    it('should hide content when loading', () => {
      fixture.componentRef.setInput('loading', true);
      fixture.detectChanges();

      const content = fixture.debugElement.query(By.css('.msls-button__content'));
      expect(content.nativeElement.classList.contains('msls-button__content--hidden')).toBe(true);
    });

    it('should disable button when loading', () => {
      fixture.componentRef.setInput('loading', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.disabled).toBe(true);
    });

    it('should set aria-busy when loading', () => {
      fixture.componentRef.setInput('loading', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.getAttribute('aria-busy')).toBe('true');
    });

    it('should apply loading class when loading', () => {
      fixture.componentRef.setInput('loading', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.classList.contains('msls-button--loading')).toBe(true);
    });
  });

  describe('disabled state', () => {
    it('should disable button when disabled', () => {
      fixture.componentRef.setInput('disabled', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.disabled).toBe(true);
    });

    it('should set aria-disabled when disabled', () => {
      fixture.componentRef.setInput('disabled', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.getAttribute('aria-disabled')).toBe('true');
    });

    it('should apply disabled class when disabled', () => {
      fixture.componentRef.setInput('disabled', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.classList.contains('msls-button--disabled')).toBe(true);
    });
  });

  describe('button type', () => {
    it('should set type attribute to submit', () => {
      fixture.componentRef.setInput('type', 'submit');
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.type).toBe('submit');
    });

    it('should set type attribute to reset', () => {
      fixture.componentRef.setInput('type', 'reset');
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.type).toBe('reset');
    });
  });

  describe('full width', () => {
    it('should apply full-width class when fullWidth is true', () => {
      fixture.componentRef.setInput('fullWidth', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.classList.contains('msls-button--full-width')).toBe(true);
    });
  });

  describe('icon only', () => {
    it('should apply icon-only class when iconOnly is true', () => {
      fixture.componentRef.setInput('iconOnly', true);
      fixture.detectChanges();

      const button = fixture.debugElement.query(By.css('button'));
      expect(button.nativeElement.classList.contains('msls-button--icon-only')).toBe(true);
    });
  });

  describe('isDisabled computed', () => {
    it('should return true when disabled is true', () => {
      fixture.componentRef.setInput('disabled', true);
      fixture.detectChanges();

      expect(component.isDisabled()).toBe(true);
    });

    it('should return true when loading is true', () => {
      fixture.componentRef.setInput('loading', true);
      fixture.detectChanges();

      expect(component.isDisabled()).toBe(true);
    });

    it('should return false when neither disabled nor loading', () => {
      expect(component.isDisabled()).toBe(false);
    });
  });
});
