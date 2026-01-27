/**
 * Guardian Section Component
 *
 * Displays and manages guardians and emergency contacts for a student.
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
import { FormsModule } from '@angular/forms';

import {
  MslsBadgeComponent,
  MslsSpinnerComponent,
} from '../../../../shared/components';
import { ToastService } from '../../../../shared/services';
import { GuardianService } from '../../services/guardian.service';
import {
  Guardian,
  GuardianRelation,
  EmergencyContact,
  CreateGuardianRequest,
  CreateEmergencyContactRequest,
  GUARDIAN_RELATION_OPTIONS,
  getGuardianRelationLabel,
} from '../../models/guardian.model';

@Component({
  selector: 'msls-guardian-section',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsBadgeComponent, MslsSpinnerComponent],
  templateUrl: './guardian-section.html',
  styleUrl: './guardian-section.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GuardianSectionComponent implements OnInit, OnChanges {
  private guardianService = inject(GuardianService);
  private toast = inject(ToastService);

  // Input
  studentId = input.required<string>();

  // State
  readonly guardians = this.guardianService.guardians;
  readonly guardiansLoading = this.guardianService.guardiansLoading;
  readonly emergencyContacts = this.guardianService.emergencyContacts;
  readonly contactsLoading = this.guardianService.contactsLoading;

  // UI State
  showGuardianForm = signal(false);
  showContactForm = signal(false);
  editingGuardian = signal<Guardian | null>(null);
  editingContact = signal<EmergencyContact | null>(null);

  // Form data
  guardianForm = signal<CreateGuardianRequest>({
    relation: 'father',
    firstName: '',
    lastName: '',
    phone: '',
  });

  contactForm = signal<CreateEmergencyContactRequest>({
    name: '',
    relation: '',
    phone: '',
  });

  // Options
  readonly relationOptions = GUARDIAN_RELATION_OPTIONS;

  ngOnInit(): void {
    this.loadData();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['studentId'] && !changes['studentId'].firstChange) {
      this.loadData();
    }
  }

  // =========================================================================
  // Data Loading
  // =========================================================================

  private loadData(): void {
    const id = this.studentId();
    if (id) {
      this.guardianService.loadGuardians(id).subscribe();
      this.guardianService.loadEmergencyContacts(id).subscribe();
    }
  }

  // =========================================================================
  // Guardian Actions
  // =========================================================================

  openGuardianForm(guardian?: Guardian): void {
    if (guardian) {
      this.editingGuardian.set(guardian);
      this.guardianForm.set({
        relation: guardian.relation,
        firstName: guardian.firstName,
        lastName: guardian.lastName,
        phone: guardian.phone,
        email: guardian.email,
        occupation: guardian.occupation,
        isPrimary: guardian.isPrimary,
        hasPortalAccess: guardian.hasPortalAccess,
      });
    } else {
      this.editingGuardian.set(null);
      this.guardianForm.set({
        relation: 'father',
        firstName: '',
        lastName: '',
        phone: '',
      });
    }
    this.showGuardianForm.set(true);
  }

  closeGuardianForm(): void {
    this.showGuardianForm.set(false);
    this.editingGuardian.set(null);
  }

  saveGuardian(): void {
    const form = this.guardianForm();
    const editing = this.editingGuardian();

    if (!form.firstName || !form.lastName || !form.phone) {
      this.toast.error('Please fill in required fields');
      return;
    }

    if (editing) {
      this.guardianService.updateGuardian(this.studentId(), editing.id, form).subscribe({
        next: () => {
          this.toast.success('Guardian updated successfully');
          this.closeGuardianForm();
        },
        error: (err) => this.toast.error(`Failed to update guardian: ${err.message}`),
      });
    } else {
      this.guardianService.createGuardian(this.studentId(), form).subscribe({
        next: () => {
          this.toast.success('Guardian added successfully');
          this.closeGuardianForm();
        },
        error: (err) => this.toast.error(`Failed to add guardian: ${err.message}`),
      });
    }
  }

  deleteGuardian(guardian: Guardian): void {
    if (confirm(`Are you sure you want to delete ${guardian.fullName}?`)) {
      this.guardianService.deleteGuardian(this.studentId(), guardian.id).subscribe({
        next: () => this.toast.success('Guardian deleted successfully'),
        error: (err) => this.toast.error(`Failed to delete guardian: ${err.message}`),
      });
    }
  }

  setPrimary(guardian: Guardian): void {
    this.guardianService.setPrimaryGuardian(this.studentId(), guardian.id).subscribe({
      next: () => this.toast.success(`${guardian.fullName} is now the primary guardian`),
      error: (err) => this.toast.error(`Failed to set primary guardian: ${err.message}`),
    });
  }

  // =========================================================================
  // Emergency Contact Actions
  // =========================================================================

  openContactForm(contact?: EmergencyContact): void {
    if (contact) {
      this.editingContact.set(contact);
      this.contactForm.set({
        name: contact.name,
        relation: contact.relation,
        phone: contact.phone,
        alternatePhone: contact.alternatePhone,
        priority: contact.priority,
        notes: contact.notes,
      });
    } else {
      this.editingContact.set(null);
      this.contactForm.set({
        name: '',
        relation: '',
        phone: '',
      });
    }
    this.showContactForm.set(true);
  }

  closeContactForm(): void {
    this.showContactForm.set(false);
    this.editingContact.set(null);
  }

  saveContact(): void {
    const form = this.contactForm();
    const editing = this.editingContact();

    if (!form.name || !form.relation || !form.phone) {
      this.toast.error('Please fill in required fields');
      return;
    }

    if (editing) {
      this.guardianService.updateEmergencyContact(this.studentId(), editing.id, form).subscribe({
        next: () => {
          this.toast.success('Emergency contact updated successfully');
          this.closeContactForm();
        },
        error: (err) => this.toast.error(`Failed to update contact: ${err.message}`),
      });
    } else {
      this.guardianService.createEmergencyContact(this.studentId(), form).subscribe({
        next: () => {
          this.toast.success('Emergency contact added successfully');
          this.closeContactForm();
        },
        error: (err) => this.toast.error(`Failed to add contact: ${err.message}`),
      });
    }
  }

  deleteContact(contact: EmergencyContact): void {
    if (confirm(`Are you sure you want to delete ${contact.name}?`)) {
      this.guardianService.deleteEmergencyContact(this.studentId(), contact.id).subscribe({
        next: () => this.toast.success('Emergency contact deleted successfully'),
        error: (err) => this.toast.error(`Failed to delete contact: ${err.message}`),
      });
    }
  }

  // =========================================================================
  // Helpers
  // =========================================================================

  getRelationLabel(relation: string): string {
    return getGuardianRelationLabel(relation as GuardianRelation) || relation;
  }

  // Type-safe event handlers for guardian form
  onGuardianRelationChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.guardianForm.update(f => ({ ...f, relation: select.value as GuardianRelation }));
  }

  onGuardianFirstNameInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.guardianForm.update(f => ({ ...f, firstName: input.value }));
  }

  onGuardianLastNameInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.guardianForm.update(f => ({ ...f, lastName: input.value }));
  }

  onGuardianPhoneInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.guardianForm.update(f => ({ ...f, phone: input.value }));
  }

  onGuardianEmailInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.guardianForm.update(f => ({ ...f, email: input.value || undefined }));
  }

  onGuardianOccupationInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.guardianForm.update(f => ({ ...f, occupation: input.value || undefined }));
  }

  onGuardianPrimaryChange(event: Event): void {
    const checkbox = event.target as HTMLInputElement;
    this.guardianForm.update(f => ({ ...f, isPrimary: checkbox.checked }));
  }

  onGuardianPortalAccessChange(event: Event): void {
    const checkbox = event.target as HTMLInputElement;
    this.guardianForm.update(f => ({ ...f, hasPortalAccess: checkbox.checked }));
  }

  // Type-safe event handlers for contact form
  onContactNameInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.contactForm.update(f => ({ ...f, name: input.value }));
  }

  onContactRelationInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.contactForm.update(f => ({ ...f, relation: input.value }));
  }

  onContactPhoneInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.contactForm.update(f => ({ ...f, phone: input.value }));
  }

  onContactAlternatePhoneInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.contactForm.update(f => ({ ...f, alternatePhone: input.value || undefined }));
  }

  onContactPriorityInput(event: Event): void {
    const input = event.target as HTMLInputElement;
    const value = parseInt(input.value, 10);
    this.contactForm.update(f => ({ ...f, priority: isNaN(value) ? undefined : value }));
  }

  onContactNotesInput(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.contactForm.update(f => ({ ...f, notes: textarea.value || undefined }));
  }
}
