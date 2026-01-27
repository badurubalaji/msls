/**
 * MSLS Feature Flag Directive
 *
 * Structural directive that conditionally renders content based on feature flags.
 * Similar to *ngIf but for feature flags.
 */

import {
  Directive,
  Input,
  TemplateRef,
  ViewContainerRef,
  inject,
  effect,
  OnDestroy,
} from '@angular/core';

import { FeatureFlagService, FeatureFlagKey } from '../../core/services/feature-flag.service';

/**
 * FeatureFlagDirective - Conditionally render content based on feature flags.
 *
 * This is a structural directive that shows/hides content based on whether
 * a feature flag is enabled for the current user/tenant.
 *
 * Usage:
 * ```html
 * <!-- Basic usage - show if flag is enabled -->
 * <div *mslsFeatureFlag="'ai_insights'">
 *   AI Insights Content
 * </div>
 *
 * <!-- With else template -->
 * <div *mslsFeatureFlag="'transport_tracking'; else noTracking">
 *   Transport Tracking Enabled
 * </div>
 * <ng-template #noTracking>
 *   Transport Tracking Not Available
 * </ng-template>
 *
 * <!-- Inverted - show if flag is disabled -->
 * <div *mslsFeatureFlag="'beta_feature'; not: true">
 *   Shown when beta_feature is disabled
 * </div>
 * ```
 */
@Directive({
  selector: '[mslsFeatureFlag]',
  standalone: true,
})
export class FeatureFlagDirective implements OnDestroy {
  private templateRef = inject(TemplateRef<FeatureFlagContext>);
  private viewContainer = inject(ViewContainerRef);
  private featureFlagService = inject(FeatureFlagService);

  private flagKey: FeatureFlagKey | null = null;
  private elseTemplateRef: TemplateRef<FeatureFlagContext> | null = null;
  private invert = false;
  private hasView = false;
  private hasElseView = false;

  /**
   * The feature flag key to check
   */
  @Input()
  set mslsFeatureFlag(flagKey: FeatureFlagKey) {
    this.flagKey = flagKey;
    this.updateView();
  }

  /**
   * Optional else template to show when flag is disabled
   */
  @Input()
  set mslsFeatureFlagElse(templateRef: TemplateRef<FeatureFlagContext> | null) {
    this.elseTemplateRef = templateRef;
    this.updateView();
  }

  /**
   * Invert the condition - show content when flag is disabled
   */
  @Input()
  set mslsFeatureFlagNot(invert: boolean) {
    this.invert = invert;
    this.updateView();
  }

  constructor() {
    // React to feature flag changes using effect
    effect(() => {
      // Access the flags signal to track changes
      this.featureFlagService.flags();
      // Update view when flags change
      this.updateView();
    });
  }

  ngOnDestroy(): void {
    this.viewContainer.clear();
  }

  private updateView(): void {
    if (!this.flagKey) {
      return;
    }

    const isEnabled = this.featureFlagService.isEnabled(this.flagKey);
    const shouldShow = this.invert ? !isEnabled : isEnabled;

    if (shouldShow) {
      // Show main content
      if (!this.hasView) {
        this.viewContainer.clear();
        this.viewContainer.createEmbeddedView(this.templateRef, {
          $implicit: isEnabled,
          mslsFeatureFlag: isEnabled,
          flagKey: this.flagKey,
        });
        this.hasView = true;
        this.hasElseView = false;
      }
    } else if (this.elseTemplateRef) {
      // Show else content
      if (!this.hasElseView) {
        this.viewContainer.clear();
        this.viewContainer.createEmbeddedView(this.elseTemplateRef, {
          $implicit: isEnabled,
          mslsFeatureFlag: isEnabled,
          flagKey: this.flagKey,
        });
        this.hasView = false;
        this.hasElseView = true;
      }
    } else {
      // Hide content
      if (this.hasView || this.hasElseView) {
        this.viewContainer.clear();
        this.hasView = false;
        this.hasElseView = false;
      }
    }
  }

  /**
   * Static method for type checking in templates
   */
  static ngTemplateContextGuard(
    _dir: FeatureFlagDirective,
    _ctx: unknown
  ): _ctx is FeatureFlagContext {
    return true;
  }
}

/**
 * Context provided to the template
 */
export interface FeatureFlagContext {
  /** The current enabled state (same as mslsFeatureFlag) */
  $implicit: boolean;
  /** The current enabled state */
  mslsFeatureFlag: boolean;
  /** The flag key being checked */
  flagKey: string;
}

/**
 * FeatureFlagEnabledDirective - Alternative attribute directive for simple show/hide.
 *
 * This is a simpler attribute directive that just adds/removes the element
 * based on the flag state. Does not support else templates.
 *
 * Usage:
 * ```html
 * <div [mslsIfFeature]="'ai_insights'">
 *   Only visible if ai_insights is enabled
 * </div>
 * ```
 */
@Directive({
  selector: '[mslsIfFeature]',
  standalone: true,
})
export class FeatureFlagEnabledDirective {
  private templateRef = inject(TemplateRef<void>);
  private viewContainer = inject(ViewContainerRef);
  private featureFlagService = inject(FeatureFlagService);

  private hasView = false;

  @Input()
  set mslsIfFeature(flagKey: FeatureFlagKey) {
    // Use effect to react to flag changes
    effect(
      () => {
        const isEnabled = this.featureFlagService.isEnabled(flagKey);
        if (isEnabled && !this.hasView) {
          this.viewContainer.createEmbeddedView(this.templateRef);
          this.hasView = true;
        } else if (!isEnabled && this.hasView) {
          this.viewContainer.clear();
          this.hasView = false;
        }
      },
      { allowSignalWrites: true }
    );
  }
}
