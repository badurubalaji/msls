import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { MslsSpinnerComponent, SpinnerSize, SpinnerVariant } from './spinner.component';

describe('MslsSpinnerComponent', () => {
  let component: MslsSpinnerComponent;
  let fixture: ComponentFixture<MslsSpinnerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MslsSpinnerComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MslsSpinnerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('default values', () => {
    it('should have md size by default', () => {
      expect(component.size()).toBe('md');
    });

    it('should have primary variant by default', () => {
      expect(component.variant()).toBe('primary');
    });

    it('should have "Loading" label by default', () => {
      expect(component.label()).toBe('Loading');
    });
  });

  describe('sizes', () => {
    const sizes: SpinnerSize[] = ['xs', 'sm', 'md', 'lg', 'xl'];

    sizes.forEach((size) => {
      it(`should apply ${size} size class`, () => {
        fixture.componentRef.setInput('size', size);
        fixture.detectChanges();

        const spinner = fixture.debugElement.query(By.css('.msls-spinner'));
        expect(spinner.nativeElement.classList.contains(`msls-spinner--${size}`)).toBe(true);
      });
    });
  });

  describe('variants', () => {
    const variants: SpinnerVariant[] = ['primary', 'secondary', 'white'];

    variants.forEach((variant) => {
      it(`should apply ${variant} variant class`, () => {
        fixture.componentRef.setInput('variant', variant);
        fixture.detectChanges();

        const spinner = fixture.debugElement.query(By.css('.msls-spinner'));
        expect(spinner.nativeElement.classList.contains(`msls-spinner--${variant}`)).toBe(true);
      });
    });
  });

  describe('accessibility', () => {
    it('should have role="status"', () => {
      const spinner = fixture.debugElement.query(By.css('.msls-spinner'));
      expect(spinner.nativeElement.getAttribute('role')).toBe('status');
    });

    it('should have aria-label', () => {
      const spinner = fixture.debugElement.query(By.css('.msls-spinner'));
      expect(spinner.nativeElement.getAttribute('aria-label')).toBe('Loading');
    });

    it('should use custom label for aria-label', () => {
      fixture.componentRef.setInput('label', 'Processing request');
      fixture.detectChanges();

      const spinner = fixture.debugElement.query(By.css('.msls-spinner'));
      expect(spinner.nativeElement.getAttribute('aria-label')).toBe('Processing request');
    });

    it('should have aria-hidden on SVG', () => {
      const svg = fixture.debugElement.query(By.css('svg'));
      expect(svg.nativeElement.getAttribute('aria-hidden')).toBe('true');
    });

    it('should have screen reader text', () => {
      const srText = fixture.debugElement.query(By.css('.sr-only'));
      expect(srText).toBeTruthy();
      expect(srText.nativeElement.textContent).toContain('Loading');
    });
  });

  describe('SVG structure', () => {
    it('should have track circle', () => {
      const track = fixture.debugElement.query(By.css('.msls-spinner__track'));
      expect(track).toBeTruthy();
    });

    it('should have head path', () => {
      const head = fixture.debugElement.query(By.css('.msls-spinner__head'));
      expect(head).toBeTruthy();
    });

    it('should have correct viewBox', () => {
      const svg = fixture.debugElement.query(By.css('svg'));
      expect(svg.nativeElement.getAttribute('viewBox')).toBe('0 0 24 24');
    });
  });

  describe('spinnerClasses computed', () => {
    it('should include base class', () => {
      expect(component.spinnerClasses()).toContain('msls-spinner');
    });

    it('should include size class', () => {
      fixture.componentRef.setInput('size', 'lg');
      fixture.detectChanges();

      expect(component.spinnerClasses()).toContain('msls-spinner--lg');
    });

    it('should include variant class', () => {
      fixture.componentRef.setInput('variant', 'secondary');
      fixture.detectChanges();

      expect(component.spinnerClasses()).toContain('msls-spinner--secondary');
    });
  });
});
