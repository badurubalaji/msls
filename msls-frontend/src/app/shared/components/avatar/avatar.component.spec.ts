import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { MslsAvatarComponent, AvatarSize } from './avatar.component';

describe('MslsAvatarComponent', () => {
  let component: MslsAvatarComponent;
  let fixture: ComponentFixture<MslsAvatarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MslsAvatarComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MslsAvatarComponent);
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

    it('should have empty src by default', () => {
      expect(component.src()).toBe('');
    });

    it('should have empty name by default', () => {
      expect(component.name()).toBe('');
    });
  });

  describe('sizes', () => {
    const sizes: AvatarSize[] = ['sm', 'md', 'lg', 'xl'];

    sizes.forEach((size) => {
      it(`should apply ${size} size class`, () => {
        fixture.componentRef.setInput('size', size);
        fixture.detectChanges();

        const avatar = fixture.debugElement.query(By.css('.msls-avatar'));
        expect(avatar.nativeElement.classList.contains(`msls-avatar--${size}`)).toBe(true);
      });
    });
  });

  describe('initials generation', () => {
    it('should generate single initial for single name', () => {
      fixture.componentRef.setInput('name', 'John');
      fixture.detectChanges();

      expect(component.initials()).toBe('J');
    });

    it('should generate two initials for full name', () => {
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.detectChanges();

      expect(component.initials()).toBe('JD');
    });

    it('should use first and last name for multiple names', () => {
      fixture.componentRef.setInput('name', 'John Michael Doe');
      fixture.detectChanges();

      expect(component.initials()).toBe('JD');
    });

    it('should handle extra whitespace', () => {
      fixture.componentRef.setInput('name', '  John   Doe  ');
      fixture.detectChanges();

      expect(component.initials()).toBe('JD');
    });

    it('should return empty string for empty name', () => {
      fixture.componentRef.setInput('name', '');
      fixture.detectChanges();

      expect(component.initials()).toBe('');
    });

    it('should uppercase initials', () => {
      fixture.componentRef.setInput('name', 'john doe');
      fixture.detectChanges();

      expect(component.initials()).toBe('JD');
    });
  });

  describe('image display', () => {
    it('should show image when src is provided', () => {
      fixture.componentRef.setInput('src', 'https://example.com/photo.jpg');
      fixture.detectChanges();

      const image = fixture.debugElement.query(By.css('.msls-avatar__image'));
      expect(image).toBeTruthy();
    });

    it('should not show initials when image is displayed', () => {
      fixture.componentRef.setInput('src', 'https://example.com/photo.jpg');
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.detectChanges();

      const initials = fixture.debugElement.query(By.css('.msls-avatar__initials'));
      expect(initials).toBeFalsy();
    });

    it('should show initials when no src', () => {
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.detectChanges();

      const initials = fixture.debugElement.query(By.css('.msls-avatar__initials'));
      expect(initials).toBeTruthy();
      expect(initials.nativeElement.textContent.trim()).toBe('JD');
    });
  });

  describe('image error handling', () => {
    it('should show initials when image fails to load', () => {
      fixture.componentRef.setInput('src', 'https://example.com/invalid.jpg');
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.detectChanges();

      // Trigger error
      component.onImageError();
      fixture.detectChanges();

      const initials = fixture.debugElement.query(By.css('.msls-avatar__initials'));
      expect(initials).toBeTruthy();
    });

    it('should reset error state on successful load', () => {
      component.onImageError();
      expect(component.imageError()).toBe(true);

      component.onImageLoad();
      expect(component.imageError()).toBe(false);
    });
  });

  describe('placeholder icon', () => {
    it('should show placeholder when no name and no image', () => {
      const placeholder = fixture.debugElement.query(By.css('.msls-avatar__placeholder'));
      expect(placeholder).toBeTruthy();
    });

    it('should not show placeholder when name is provided', () => {
      fixture.componentRef.setInput('name', 'John');
      fixture.detectChanges();

      const placeholder = fixture.debugElement.query(By.css('.msls-avatar__placeholder'));
      expect(placeholder).toBeFalsy();
    });
  });

  describe('accessibility', () => {
    it('should have role="img"', () => {
      const avatar = fixture.debugElement.query(By.css('.msls-avatar'));
      expect(avatar.nativeElement.getAttribute('role')).toBe('img');
    });

    it('should have aria-label with name', () => {
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.detectChanges();

      const avatar = fixture.debugElement.query(By.css('.msls-avatar'));
      expect(avatar.nativeElement.getAttribute('aria-label')).toBe('John Doe');
    });

    it('should use alt prop for aria-label when provided', () => {
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.componentRef.setInput('alt', 'Profile picture');
      fixture.detectChanges();

      const avatar = fixture.debugElement.query(By.css('.msls-avatar'));
      expect(avatar.nativeElement.getAttribute('aria-label')).toBe('Profile picture');
    });

    it('should have default aria-label when no name', () => {
      const avatar = fixture.debugElement.query(By.css('.msls-avatar'));
      expect(avatar.nativeElement.getAttribute('aria-label')).toBe('Avatar');
    });

    it('should have aria-hidden on initials', () => {
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.detectChanges();

      const initials = fixture.debugElement.query(By.css('.msls-avatar__initials'));
      expect(initials.nativeElement.getAttribute('aria-hidden')).toBe('true');
    });
  });

  describe('altText computed', () => {
    it('should return alt when provided', () => {
      fixture.componentRef.setInput('alt', 'Custom alt');
      fixture.componentRef.setInput('name', 'John');
      fixture.detectChanges();

      expect(component.altText()).toBe('Custom alt');
    });

    it('should return name when no alt', () => {
      fixture.componentRef.setInput('name', 'John Doe');
      fixture.detectChanges();

      expect(component.altText()).toBe('John Doe');
    });

    it('should return "Avatar" when no alt and no name', () => {
      expect(component.altText()).toBe('Avatar');
    });
  });
});
