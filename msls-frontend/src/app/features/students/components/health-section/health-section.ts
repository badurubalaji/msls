/**
 * Health Section Component
 *
 * Displays student health records including allergies, conditions,
 * medications, vaccinations, and medical incidents.
 */

import {
  Component,
  ChangeDetectionStrategy,
  inject,
  input,
  OnInit,
  OnChanges,
  SimpleChanges,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';

import {
  MslsBadgeComponent,
  MslsSpinnerComponent,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services';
import { HealthService } from '../../services/health.service';
import {
  HealthSummary,
  Allergy,
  ChronicCondition,
  Medication,
  Vaccination,
  MedicalIncident,
  getAllergyTypeLabel,
  getAllergySeverityLabel,
  getConditionTypeLabel,
  getMedicationFrequencyLabel,
  getIncidentTypeLabel,
  getSeverityColor,
  BLOOD_GROUP_OPTIONS,
  ALLERGY_TYPE_OPTIONS,
  ALLERGY_SEVERITY_OPTIONS,
  CreateAllergyRequest,
  AllergyType,
  AllergySeverity,
} from '../../models/health.model';

@Component({
  selector: 'msls-health-section',
  standalone: true,
  imports: [CommonModule, MslsBadgeComponent, MslsSpinnerComponent],
  templateUrl: './health-section.html',
  styleUrl: './health-section.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HealthSectionComponent implements OnInit, OnChanges {
  private healthService = inject(HealthService);
  private toast = inject(ToastService);

  // Input
  studentId = input.required<string>();

  // State
  readonly healthSummary = this.healthService.healthSummary;
  readonly loading = this.healthService.loading;

  // UI State
  activeTab = signal<'overview' | 'allergies' | 'conditions' | 'medications' | 'vaccinations' | 'incidents'>('overview');
  showAllergyForm = signal(false);
  editingAllergy = signal<Allergy | null>(null);

  // Form data
  allergyForm = signal<CreateAllergyRequest>({
    allergen: '',
    allergyType: 'food',
    severity: 'mild',
  });

  // Options
  readonly bloodGroupOptions = BLOOD_GROUP_OPTIONS;
  readonly allergyTypeOptions = ALLERGY_TYPE_OPTIONS;
  readonly allergySeverityOptions = ALLERGY_SEVERITY_OPTIONS;

  ngOnInit(): void {
    this.loadData();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['studentId'] && !changes['studentId'].firstChange) {
      this.loadData();
    }
  }

  private loadData(): void {
    const id = this.studentId();
    if (id) {
      this.healthService.loadHealthSummary(id).subscribe({
        error: (err) => this.toast.error(`Failed to load health records: ${err.message}`),
      });
    }
  }

  // Tab navigation
  setTab(tab: 'overview' | 'allergies' | 'conditions' | 'medications' | 'vaccinations' | 'incidents'): void {
    this.activeTab.set(tab);
  }

  // Allergy actions
  openAllergyForm(allergy?: Allergy): void {
    if (allergy) {
      this.editingAllergy.set(allergy);
      this.allergyForm.set({
        allergen: allergy.allergen,
        allergyType: allergy.allergyType,
        severity: allergy.severity,
        reactionDescription: allergy.reactionDescription,
        treatmentInstructions: allergy.treatmentInstructions,
        emergencyMedication: allergy.emergencyMedication,
        diagnosedDate: allergy.diagnosedDate,
        notes: allergy.notes,
      });
    } else {
      this.editingAllergy.set(null);
      this.allergyForm.set({
        allergen: '',
        allergyType: 'food',
        severity: 'mild',
      });
    }
    this.showAllergyForm.set(true);
  }

  closeAllergyForm(): void {
    this.showAllergyForm.set(false);
    this.editingAllergy.set(null);
  }

  saveAllergy(): void {
    const form = this.allergyForm();
    const editing = this.editingAllergy();

    if (!form.allergen) {
      this.toast.error('Please enter the allergen');
      return;
    }

    if (editing) {
      this.healthService.updateAllergy(this.studentId(), editing.id, form).subscribe({
        next: () => {
          this.toast.success('Allergy updated successfully');
          this.closeAllergyForm();
        },
        error: (err) => this.toast.error(`Failed to update allergy: ${err.message}`),
      });
    } else {
      this.healthService.createAllergy(this.studentId(), form).subscribe({
        next: () => {
          this.toast.success('Allergy added successfully');
          this.closeAllergyForm();
        },
        error: (err) => this.toast.error(`Failed to add allergy: ${err.message}`),
      });
    }
  }

  deleteAllergy(allergy: Allergy): void {
    if (confirm(`Are you sure you want to delete this allergy record for "${allergy.allergen}"?`)) {
      this.healthService.deleteAllergy(this.studentId(), allergy.id).subscribe({
        next: () => this.toast.success('Allergy deleted successfully'),
        error: (err) => this.toast.error(`Failed to delete allergy: ${err.message}`),
      });
    }
  }

  // Helper methods
  getAllergyTypeLabel(type: string): string {
    return getAllergyTypeLabel(type as AllergyType);
  }

  getAllergySeverityLabel(severity: string): string {
    return getAllergySeverityLabel(severity as AllergySeverity);
  }

  getConditionTypeLabel(type: string): string {
    return getConditionTypeLabel(type as any);
  }

  getMedicationFrequencyLabel(frequency: string): string {
    return getMedicationFrequencyLabel(frequency as any);
  }

  getIncidentTypeLabel(type: string): string {
    return getIncidentTypeLabel(type as any);
  }

  getSeverityVariant(severity: string): 'success' | 'warning' | 'danger' | 'neutral' {
    return getSeverityColor(severity as any) as any;
  }

  formatDate(dateString: string | undefined): string {
    if (!dateString) return '—';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  // Type-safe event handlers
  onAllergenInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.allergyForm.update(f => ({ ...f, allergen: input.value }));
  }

  onAllergyTypeChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.allergyForm.update(f => ({ ...f, allergyType: select.value as AllergyType }));
  }

  onSeverityChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.allergyForm.update(f => ({ ...f, severity: select.value as AllergySeverity }));
  }

  onReactionInput(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.allergyForm.update(f => ({ ...f, reactionDescription: textarea.value || undefined }));
  }

  onTreatmentInput(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.allergyForm.update(f => ({ ...f, treatmentInstructions: textarea.value || undefined }));
  }

  onEmergencyMedInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.allergyForm.update(f => ({ ...f, emergencyMedication: input.value || undefined }));
  }
}
