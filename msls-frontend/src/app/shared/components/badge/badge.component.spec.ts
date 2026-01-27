import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { MslsBadgeComponent, BadgeVariant, BadgeSize } from './badge.component';

describe('MslsBadgeComponent', () => {
  let component: MslsBadgeComponent;
  let fixture: ComponentFixture<MslsBadgeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MslsBadgeComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MslsBadgeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('default values', () => {
    it('should have neutral variant by default', () => {
      expect(component.variant()).toBe('neutral');
    });

    it('should have md size by default', () => {
      expect(component.size()).toBe('md');
    });

    it('should not show dot by default', () => {
      expect(component.dot()).toBe(false);
    });
  });

  describe('variants', () => {
    const variants: BadgeVariant[] = ['success', 'warning', 'error', 'info', 'neutral'];

    variants.forEach((variant) => {
      it(`should apply ${variant} variant class`, () => {
        fixture.componentRef.setInput('variant', variant);
        fixture.detectChanges();

        const badge = fixture.debugElement.query(By.css('.msls-badge'));
        expect(badge.nativeElement.classList.contains(`msls-badge--${variant}`)).toBe(true);
      });
    });
  });

  describe('sizes', () => {
    const sizes: BadgeSize[] = ['sm', 'md'];

    sizes.forEach((size) => {
      it(`should apply ${size} size class`, () => {
        fixture.componentRef.setInput('size', size);
        fixture.detectChanges();

        const badge = fixture.debugElement.query(By.css('.msls-badge'));
        expect(badge.nativeElement.classList.contains(`msls-badge--${size}`)).toBe(true);
      });
    });
  });

  describe('dot indicator', () => {
    it('should not show dot by default', () => {
      const dot = fixture.debugElement.query(By.css('.msls-badge__dot'));
      expect(dot).toBeFalsy();
    });

    it('should show dot when dot input is true', () => {
      fixture.componentRef.setInput('dot', true);
      fixture.detectChanges();

      const dot = fixture.debugElement.query(By.css('.msls-badge__dot'));
      expect(dot).toBeTruthy();
    });

    it('should apply dot class when dot is true', () => {
      fixture.componentRef.setInput('dot', true);
      fixture.detectChanges();

      const badge = fixture.debugElement.query(By.css('.msls-badge'));
      expect(badge.nativeElement.classList.contains('msls-badge--dot')).toBe(true);
    });

    it('should have aria-hidden on dot element', () => {
      fixture.componentRef.setInput('dot', true);
      fixture.detectChanges();

      const dot = fixture.debugElement.query(By.css('.msls-badge__dot'));
      expect(dot.nativeElement.getAttribute('aria-hidden')).toBe('true');
    });
  });

  describe('content projection', () => {
    it('should have content container', () => {
      const content = fixture.debugElement.query(By.css('.msls-badge__content'));
      expect(content).toBeTruthy();
    });
  });

  describe('badgeClasses computed', () => {
    it('should include base class', () => {
      expect(component.badgeClasses()).toContain('msls-badge');
    });

    it('should include variant class', () => {
      fixture.componentRef.setInput('variant', 'success');
      fixture.detectChanges();

      expect(component.badgeClasses()).toContain('msls-badge--success');
    });

    it('should include size class', () => {
      fixture.componentRef.setInput('size', 'sm');
      fixture.detectChanges();

      expect(component.badgeClasses()).toContain('msls-badge--sm');
    });

    it('should include dot class when dot is true', () => {
      fixture.componentRef.setInput('dot', true);
      fixture.detectChanges();

      expect(component.badgeClasses()).toContain('msls-badge--dot');
    });

    it('should not include dot class when dot is false', () => {
      expect(component.badgeClasses()).not.toContain('msls-badge--dot');
    });
  });
});
