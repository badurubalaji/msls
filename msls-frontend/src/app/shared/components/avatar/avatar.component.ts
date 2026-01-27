import { Component, input, signal, computed, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/** Avatar size types */
export type AvatarSize = 'sm' | 'md' | 'lg' | 'xl';

/**
 * MslsAvatarComponent - Reusable avatar component with image and fallback initials.
 *
 * @example
 * ```html
 * <msls-avatar src="https://example.com/photo.jpg" name="John Doe" size="md" />
 * <msls-avatar name="Jane Smith" size="lg" />
 * ```
 */
@Component({
  selector: 'msls-avatar',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './avatar.component.html',
  styleUrl: './avatar.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MslsAvatarComponent {
  /** Image source URL */
  readonly src = input<string>('');

  /** Name for generating initials and alt text */
  readonly name = input<string>('');

  /** Avatar size: sm, md, lg, xl */
  readonly size = input<AvatarSize>('md');

  /** Alt text override (defaults to name) */
  readonly alt = input<string>('');

  /** Internal state for image load error */
  readonly imageError = signal<boolean>(false);

  /** Whether to show the image */
  readonly showImage = computed(() => {
    return this.src() && !this.imageError();
  });

  /** Generated initials from name */
  readonly initials = computed(() => {
    const name = this.name();
    if (!name) return '';

    const parts = name.trim().split(/\s+/);
    if (parts.length === 1) {
      return parts[0].charAt(0).toUpperCase();
    }

    // Take first letter of first and last parts
    return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
  });

  /** Alt text for the image */
  readonly altText = computed(() => {
    return this.alt() || this.name() || 'Avatar';
  });

  /** Computed CSS classes based on size */
  readonly avatarClasses = computed(() => {
    const classes: string[] = ['msls-avatar'];

    // Size classes
    classes.push(`msls-avatar--${this.size()}`);

    return classes.join(' ');
  });

  /** Handle image load error - show fallback initials */
  onImageError(): void {
    this.imageError.set(true);
  }

  /** Reset error state when src changes */
  onImageLoad(): void {
    this.imageError.set(false);
  }
}
