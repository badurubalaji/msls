import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { MslsIconComponent, IconSize } from './icon.component';

describe('MslsIconComponent', () => {
  let component: MslsIconComponent;
  let fixture: ComponentFixture<MslsIconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MslsIconComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MslsIconComponent);
    component = fixture.componentInstance;
    // Set required input
    fixture.componentRef.setInput('name', 'user');
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('default values', () => {
    it('should have md size by default', () => {
      expect(component.size()).toBe('md');
    });

    it('should have empty label by default', () => {
      expect(component.label()).toBe('');
    });
  });

  describe('sizes', () => {
    const sizes: IconSize[] = ['xs', 'sm', 'md', 'lg', 'xl'];

    sizes.forEach((size) => {
      it(`should apply ${size} size class`, () => {
        fixture.componentRef.setInput('size', size);
        fixture.detectChanges();

        const icon = fixture.debugElement.query(By.css('i'));
        expect(icon.nativeElement.classList.contains(`msls-icon--${size}`)).toBe(true);
      });
    });
  });

  describe('icon classes', () => {
    it('should render user icon with Font Awesome class', () => {
      fixture.componentRef.setInput('name', 'user');
      fixture.detectChanges();

      const icon = fixture.debugElement.query(By.css('i'));
      expect(icon.nativeElement.classList.contains('fa-regular')).toBe(true);
      expect(icon.nativeElement.classList.contains('fa-user')).toBe(true);
    });

    it('should render check icon with Font Awesome class', () => {
      fixture.componentRef.setInput('name', 'check');
      fixture.detectChanges();

      const icon = fixture.debugElement.query(By.css('i'));
      expect(icon.nativeElement.classList.contains('fa-solid')).toBe(true);
      expect(icon.nativeElement.classList.contains('fa-check')).toBe(true);
    });

    it('should render fallback for unknown icon', () => {
      fixture.componentRef.setInput('name', 'unknown-icon');
      fixture.detectChanges();

      const icon = fixture.debugElement.query(By.css('i'));
      expect(icon.nativeElement.classList.contains('fa-solid')).toBe(true);
      expect(icon.nativeElement.classList.contains('fa-question')).toBe(true);
    });
  });

  describe('accessibility', () => {
    it('should have aria-hidden when decorative (no label)', () => {
      const icon = fixture.debugElement.query(By.css('i'));
      expect(icon.nativeElement.getAttribute('aria-hidden')).toBe('true');
    });

    it('should have role="img"', () => {
      const icon = fixture.debugElement.query(By.css('i'));
      expect(icon.nativeElement.getAttribute('role')).toBe('img');
    });

    it('should have aria-label when label is provided', () => {
      fixture.componentRef.setInput('label', 'User profile');
      fixture.detectChanges();

      const icon = fixture.debugElement.query(By.css('i'));
      expect(icon.nativeElement.getAttribute('aria-label')).toBe('User profile');
    });

    it('should not have aria-hidden when label is provided', () => {
      fixture.componentRef.setInput('label', 'User profile');
      fixture.detectChanges();

      const icon = fixture.debugElement.query(By.css('i'));
      expect(icon.nativeElement.getAttribute('aria-hidden')).toBe('false');
    });
  });

  describe('isDecorative computed', () => {
    it('should return true when no label', () => {
      expect(component.isDecorative()).toBe(true);
    });

    it('should return false when label is provided', () => {
      fixture.componentRef.setInput('label', 'Some label');
      fixture.detectChanges();

      expect(component.isDecorative()).toBe(false);
    });
  });

  describe('icon classes computed', () => {
    it('should compute iconClasses with size and icon', () => {
      fixture.componentRef.setInput('name', 'home');
      fixture.componentRef.setInput('size', 'lg');
      fixture.detectChanges();

      expect(component.iconClasses()).toContain('msls-icon--lg');
      expect(component.iconClasses()).toContain('fa-solid fa-house');
    });
  });

  describe('available icons', () => {
    const iconNames = [
      'user',
      'check',
      'x-mark',
      'plus',
      'minus',
      'home',
      'bell',
      'envelope',
      'magnifying-glass',
    ];

    iconNames.forEach((iconName) => {
      it(`should have FA class for ${iconName}`, () => {
        fixture.componentRef.setInput('name', iconName);
        fixture.detectChanges();

        const icon = fixture.debugElement.query(By.css('i'));
        expect(icon.nativeElement.classList.contains('fa-solid') || icon.nativeElement.classList.contains('fa-regular')).toBe(true);
      });
    });
  });
});
