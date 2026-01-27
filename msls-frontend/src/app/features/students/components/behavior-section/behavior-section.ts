/**
 * Behavior Section Component
 *
 * Displays student behavioral incidents and summary.
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
  computed,
} from '@angular/core';
import { CommonModule } from '@angular/common';

import {
  MslsBadgeComponent,
  MslsSpinnerComponent,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services';
import { BehavioralService } from '../../services/behavioral.service';
import {
  BehavioralIncident,
  BehaviorSummary,
  BehavioralIncidentType,
  BehavioralSeverity,
  CreateBehavioralIncidentRequest,
  getBehavioralIncidentTypeLabel,
  getBehavioralSeverityLabel,
  getBehavioralIncidentTypeVariant,
  getBehavioralSeverityVariant,
  getTrendIcon,
  getTrendVariant,
  BEHAVIORAL_INCIDENT_TYPE_OPTIONS,
  BEHAVIORAL_SEVERITY_OPTIONS,
} from '../../models/behavioral.model';

@Component({
  selector: 'msls-behavior-section',
  standalone: true,
  imports: [CommonModule, MslsBadgeComponent, MslsSpinnerComponent],
  templateUrl: './behavior-section.html',
  styleUrl: './behavior-section.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BehaviorSectionComponent implements OnInit, OnChanges {
  private behavioralService = inject(BehavioralService);
  private toast = inject(ToastService);

  // Input
  studentId = input.required<string>();

  // State
  readonly incidents = this.behavioralService.incidents;
  readonly summary = this.behavioralService.summary;
  readonly loading = this.behavioralService.loading;

  // UI State
  showIncidentForm = signal(false);
  editingIncident = signal<BehavioralIncident | null>(null);

  // Form data
  incidentForm = signal<CreateBehavioralIncidentRequest>({
    incidentType: 'minor_infraction',
    severity: 'medium',
    incidentDate: new Date().toISOString().split('T')[0],
    incidentTime: new Date().toTimeString().slice(0, 5),
    description: '',
    actionTaken: '',
  });

  // Options
  readonly incidentTypeOptions = BEHAVIORAL_INCIDENT_TYPE_OPTIONS;
  readonly severityOptions = BEHAVIORAL_SEVERITY_OPTIONS;

  // Computed
  readonly recentIncidents = computed(() => this.incidents().slice(0, 10));
  readonly hasIncidents = computed(() => this.incidents().length > 0);

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
      this.behavioralService.loadSummary(id).subscribe({
        error: (err) => this.toast.error(`Failed to load behavior summary: ${err.message}`),
      });
      this.behavioralService.loadIncidents(id, { limit: 10 }).subscribe({
        error: (err) => this.toast.error(`Failed to load incidents: ${err.message}`),
      });
    }
  }

  // Incident actions
  openIncidentForm(incident?: BehavioralIncident): void {
    if (incident) {
      this.editingIncident.set(incident);
      this.incidentForm.set({
        incidentType: incident.incidentType,
        severity: incident.severity,
        incidentDate: incident.incidentDate,
        incidentTime: incident.incidentTime,
        location: incident.location,
        description: incident.description,
        witnesses: incident.witnesses,
        studentResponse: incident.studentResponse,
        actionTaken: incident.actionTaken,
        parentMeetingRequired: incident.parentMeetingRequired,
      });
    } else {
      this.editingIncident.set(null);
      this.incidentForm.set({
        incidentType: 'minor_infraction',
        severity: 'medium',
        incidentDate: new Date().toISOString().split('T')[0],
        incidentTime: new Date().toTimeString().slice(0, 5),
        description: '',
        actionTaken: '',
      });
    }
    this.showIncidentForm.set(true);
  }

  closeIncidentForm(): void {
    this.showIncidentForm.set(false);
    this.editingIncident.set(null);
  }

  saveIncident(): void {
    const form = this.incidentForm();
    const editing = this.editingIncident();

    if (!form.description) {
      this.toast.error('Please enter a description');
      return;
    }
    if (!form.actionTaken) {
      this.toast.error('Please enter the action taken');
      return;
    }

    if (editing) {
      this.behavioralService.updateIncident(this.studentId(), editing.id, form).subscribe({
        next: () => {
          this.toast.success('Incident updated successfully');
          this.closeIncidentForm();
        },
        error: (err) => this.toast.error(`Failed to update incident: ${err.message}`),
      });
    } else {
      this.behavioralService.createIncident(this.studentId(), form).subscribe({
        next: () => {
          this.toast.success('Incident recorded successfully');
          this.closeIncidentForm();
        },
        error: (err) => this.toast.error(`Failed to record incident: ${err.message}`),
      });
    }
  }

  deleteIncident(incident: BehavioralIncident): void {
    if (confirm('Are you sure you want to delete this incident record?')) {
      this.behavioralService.deleteIncident(this.studentId(), incident.id).subscribe({
        next: () => this.toast.success('Incident deleted successfully'),
        error: (err) => this.toast.error(`Failed to delete incident: ${err.message}`),
      });
    }
  }

  // Helper methods
  getIncidentTypeLabel(type: BehavioralIncidentType): string {
    return getBehavioralIncidentTypeLabel(type);
  }

  getSeverityLabel(severity: BehavioralSeverity): string {
    return getBehavioralSeverityLabel(severity);
  }

  getIncidentTypeVariant(type: BehavioralIncidentType): 'success' | 'warning' | 'danger' {
    return getBehavioralIncidentTypeVariant(type);
  }

  getSeverityVariant(severity: BehavioralSeverity): 'success' | 'warning' | 'danger' | 'neutral' {
    return getBehavioralSeverityVariant(severity);
  }

  getTrendIcon(trend: string): string {
    return getTrendIcon(trend);
  }

  getTrendVariant(trend: string): 'success' | 'warning' | 'danger' | 'neutral' {
    return getTrendVariant(trend);
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

  formatTime(timeString: string | undefined): string {
    if (!timeString) return '';
    return timeString.slice(0, 5);
  }

  // Type-safe event handlers
  onIncidentTypeChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.incidentForm.update(f => ({ ...f, incidentType: select.value as BehavioralIncidentType }));
  }

  onSeverityChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.incidentForm.update(f => ({ ...f, severity: select.value as BehavioralSeverity }));
  }

  onDateChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.incidentForm.update(f => ({ ...f, incidentDate: input.value }));
  }

  onTimeChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.incidentForm.update(f => ({ ...f, incidentTime: input.value }));
  }

  onLocationInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.incidentForm.update(f => ({ ...f, location: input.value || undefined }));
  }

  onDescriptionInput(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.incidentForm.update(f => ({ ...f, description: textarea.value }));
  }

  onWitnessesInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    const witnesses = input.value ? input.value.split(',').map(w => w.trim()).filter(w => w) : undefined;
    this.incidentForm.update(f => ({ ...f, witnesses }));
  }

  onStudentResponseInput(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.incidentForm.update(f => ({ ...f, studentResponse: textarea.value || undefined }));
  }

  onActionTakenInput(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.incidentForm.update(f => ({ ...f, actionTaken: textarea.value }));
  }

  onParentMeetingChange(event: Event): void {
    const checkbox = event.target as HTMLInputElement;
    this.incidentForm.update(f => ({ ...f, parentMeetingRequired: checkbox.checked }));
  }
}
