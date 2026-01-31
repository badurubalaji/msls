import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ExamService } from '../exam.service';
import {
  HallTicketTemplate,
  CreateHallTicketTemplateRequest,
  UpdateHallTicketTemplateRequest,
} from '../exam.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'app-hall-ticket-template',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <button class="back-btn" [routerLink]="['/exams']">
            <i class="fa-solid fa-arrow-left"></i>
          </button>
          <div class="header-icon">
            <i class="fa-solid fa-file-invoice"></i>
          </div>
          <div class="header-text">
            <h1>Hall Ticket Templates</h1>
            <p>Manage hall ticket branding and layout</p>
          </div>
        </div>
        <div class="header-actions">
          <button class="btn btn-primary" (click)="openCreateModal()">
            <i class="fa-solid fa-plus"></i>
            <span class="btn-text">New Template</span>
          </button>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading templates...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadTemplates()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <div class="templates-grid">
            @for (template of templates(); track template.id) {
              <div class="template-card" [class.is-default]="template.isDefault">
                @if (template.isDefault) {
                  <div class="default-badge">
                    <i class="fa-solid fa-star"></i> Default
                  </div>
                }
                <div class="template-header">
                  <div class="template-icon">
                    <i class="fa-solid fa-file-invoice"></i>
                  </div>
                  <h3>{{ template.name }}</h3>
                </div>
                <div class="template-body">
                  <div class="template-info">
                    <p class="school-name">{{ template.schoolName || 'No school name' }}</p>
                    <p class="school-address">{{ template.schoolAddress || 'No address' }}</p>
                  </div>
                  @if (template.instructions) {
                    <div class="instructions-preview">
                      <strong>Instructions:</strong>
                      <p>{{ truncate(template.instructions, 100) }}</p>
                    </div>
                  }
                </div>
                <div class="template-footer">
                  <button
                    class="action-btn"
                    title="Edit"
                    (click)="editTemplate(template)"
                  >
                    <i class="fa-regular fa-pen-to-square"></i>
                  </button>
                  @if (!template.isDefault) {
                    <button
                      class="action-btn action-btn--success"
                      title="Set as Default"
                      (click)="setAsDefault(template)"
                    >
                      <i class="fa-regular fa-star"></i>
                    </button>
                  }
                  <button
                    class="action-btn action-btn--danger"
                    title="Delete"
                    (click)="confirmDelete(template)"
                    [disabled]="template.isDefault"
                  >
                    <i class="fa-regular fa-trash-can"></i>
                  </button>
                </div>
              </div>
            } @empty {
              <div class="empty-state">
                <i class="fa-regular fa-folder-open"></i>
                <p>No templates found</p>
                <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                  <i class="fa-solid fa-plus"></i>
                  Create First Template
                </button>
              </div>
            }
          </div>
        }
      </div>

      <!-- Create/Edit Modal -->
      @if (showFormModal()) {
        <div class="modal-overlay" (click)="closeFormModal()">
          <div class="modal modal--lg" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-file-invoice"></i>
                {{ modalTitle() }}
              </h3>
              <button type="button" class="modal__close" (click)="closeFormModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form class="form">
                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-info-circle"></i>
                    Template Information
                  </h4>
                  <div class="form-group">
                    <label for="name">Template Name <span class="required">*</span></label>
                    <input
                      type="text"
                      id="name"
                      [(ngModel)]="formData.name"
                      name="name"
                      class="form-input"
                      placeholder="e.g., Default Template, Secondary Template"
                      required
                    />
                  </div>
                  <div class="form-group">
                    <label>
                      <input
                        type="checkbox"
                        [(ngModel)]="formData.isDefault"
                        name="isDefault"
                      />
                      Set as default template
                    </label>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-school"></i>
                    School Branding
                  </h4>
                  <div class="form-group">
                    <label for="schoolName">School Name</label>
                    <input
                      type="text"
                      id="schoolName"
                      [(ngModel)]="formData.schoolName"
                      name="schoolName"
                      class="form-input"
                      placeholder="Enter school name"
                    />
                  </div>
                  <div class="form-group">
                    <label for="schoolAddress">School Address</label>
                    <textarea
                      id="schoolAddress"
                      [(ngModel)]="formData.schoolAddress"
                      name="schoolAddress"
                      class="form-textarea"
                      rows="2"
                      placeholder="Enter school address"
                    ></textarea>
                  </div>
                  <div class="form-group">
                    <label for="headerLogoUrl">Header Logo URL</label>
                    <input
                      type="text"
                      id="headerLogoUrl"
                      [(ngModel)]="formData.headerLogoUrl"
                      name="headerLogoUrl"
                      class="form-input"
                      placeholder="https://example.com/logo.png"
                    />
                    <p class="form-hint">URL to school logo image for hall ticket header</p>
                  </div>
                </div>

                <div class="form-section">
                  <h4 class="form-section-title">
                    <i class="fa-solid fa-list-check"></i>
                    Instructions
                  </h4>
                  <div class="form-group">
                    <label for="instructions">Exam Instructions</label>
                    <textarea
                      id="instructions"
                      [(ngModel)]="formData.instructions"
                      name="instructions"
                      class="form-textarea"
                      rows="5"
                      placeholder="Enter instructions to be printed on hall tickets..."
                    ></textarea>
                    <p class="form-hint">These instructions will appear on every hall ticket</p>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--secondary" (click)="closeFormModal()">
                Cancel
              </button>
              <button
                type="button"
                class="btn btn--primary"
                [disabled]="saving() || !isFormValid()"
                (click)="saveTemplate()"
              >
                @if (saving()) {
                  <div class="btn-spinner"></div>
                  {{ editingTemplate() ? 'Saving...' : 'Creating...' }}
                } @else {
                  <i class="fa-solid" [class.fa-check]="editingTemplate()" [class.fa-plus]="!editingTemplate()"></i>
                  {{ editingTemplate() ? 'Save Changes' : 'Create Template' }}
                }
              </button>
            </div>
          </div>
        </div>
      }

      <!-- Delete Confirmation Modal -->
      @if (showDeleteModal()) {
        <div class="modal-overlay" (click)="closeDeleteModal()">
          <div class="modal modal--sm" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-trash"></i>
                Delete Template
              </h3>
              <button type="button" class="modal__close" (click)="closeDeleteModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <div class="delete-confirmation">
                <div class="delete-icon">
                  <i class="fa-solid fa-triangle-exclamation"></i>
                </div>
                <p>
                  Are you sure you want to delete
                  <strong>"{{ templateToDelete()?.name }}"</strong>?
                </p>
                <p class="delete-warning">
                  This action cannot be undone.
                </p>
              </div>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn--secondary" (click)="closeDeleteModal()">
                Cancel
              </button>
              <button
                type="button"
                class="btn btn--danger"
                [disabled]="deleting()"
                (click)="deleteTemplate()"
              >
                @if (deleting()) {
                  <div class="btn-spinner"></div>
                  Deleting...
                } @else {
                  <i class="fa-solid fa-trash"></i>
                  Delete
                }
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
      gap: 1rem;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .back-btn {
      width: 2.5rem;
      height: 2.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f1f5f9;
      border: none;
      border-radius: 0.5rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.2s;
    }

    .back-btn:hover {
      background: #e2e8f0;
      color: #1e293b;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      background: linear-gradient(135deg, #f59e0b, #d97706);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .header-text h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .header-text p {
      margin: 0.25rem 0 0;
      color: #64748b;
      font-size: 0.875rem;
    }

    .header-actions {
      display: flex;
      gap: 0.5rem;
    }

    /* Content Card */
    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    .loading-container, .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-container {
      color: #dc2626;
      flex-direction: column;
    }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #f59e0b;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    /* Templates Grid */
    .templates-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 1.5rem;
    }

    .template-card {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      padding: 1.25rem;
      position: relative;
      transition: all 0.2s;
    }

    .template-card:hover {
      border-color: #f59e0b;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    .template-card.is-default {
      border-color: #f59e0b;
      background: linear-gradient(135deg, #fffbeb, #fef3c7);
    }

    .default-badge {
      position: absolute;
      top: -0.5rem;
      right: 1rem;
      background: #f59e0b;
      color: white;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 600;
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .template-header {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      margin-bottom: 1rem;
    }

    .template-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #fef3c7;
      color: #d97706;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .template-header h3 {
      margin: 0;
      font-size: 1rem;
      font-weight: 600;
      color: #1e293b;
    }

    .template-body {
      margin-bottom: 1rem;
    }

    .template-info {
      margin-bottom: 0.75rem;
    }

    .school-name {
      margin: 0 0 0.25rem;
      font-weight: 500;
      color: #374151;
    }

    .school-address {
      margin: 0;
      font-size: 0.875rem;
      color: #64748b;
    }

    .instructions-preview {
      padding: 0.75rem;
      background: white;
      border-radius: 0.5rem;
      font-size: 0.875rem;
    }

    .instructions-preview strong {
      color: #374151;
      font-size: 0.75rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .instructions-preview p {
      margin: 0.25rem 0 0;
      color: #64748b;
    }

    .template-footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.5rem;
      padding-top: 0.75rem;
      border-top: 1px solid #e2e8f0;
    }

    /* Actions */
    .action-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      border: none;
      background: white;
      color: #64748b;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: all 0.2s;
    }

    .action-btn:hover:not(:disabled) {
      background: #f1f5f9;
      color: #f59e0b;
    }

    .action-btn:disabled {
      opacity: 0.3;
      cursor: not-allowed;
    }

    .action-btn--success:hover:not(:disabled) {
      background: #dcfce7;
      color: #16a34a;
    }

    .action-btn--danger:hover:not(:disabled) {
      background: #fef2f2;
      color: #dc2626;
    }

    /* Empty State */
    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 0.75rem;
      padding: 3rem;
      color: #64748b;
      grid-column: 1 / -1;
    }

    .empty-state i {
      font-size: 2rem;
      color: #cbd5e1;
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
    }

    .btn-primary {
      background: linear-gradient(135deg, #f59e0b, #d97706);
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: linear-gradient(135deg, #d97706, #b45309);
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 16px;
      height: 16px;
      border: 2px solid transparent;
      border-top-color: currentColor;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    /* Modal Styles */
    .modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      padding: 1rem;
    }

    .modal {
      background: white;
      border-radius: 1.25rem;
      width: 100%;
      max-height: 90vh;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
    }

    .modal--sm { max-width: 28rem; }
    .modal--lg { max-width: 48rem; }

    .modal__header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid #f1f5f9;
    }

    .modal__header h3 {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
    }

    .modal__header h3 i {
      color: #f59e0b;
    }

    .modal__close {
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f1f5f9;
      border: none;
      border-radius: 0.5rem;
      cursor: pointer;
      color: #64748b;
      transition: all 0.2s ease;
    }

    .modal__close:hover {
      background: #e2e8f0;
      color: #334155;
    }

    .modal__body {
      padding: 1.5rem;
      overflow-y: auto;
      flex: 1;
    }

    .modal__footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding: 1.25rem 1.5rem;
      background: #f8fafc;
      border-top: 1px solid #f1f5f9;
      border-radius: 0 0 1.25rem 1.25rem;
    }

    .btn--primary {
      background: linear-gradient(135deg, #f59e0b, #d97706);
      color: white;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn--primary:hover:not(:disabled) { background: linear-gradient(135deg, #d97706, #b45309); }
    .btn--primary:disabled { opacity: 0.5; cursor: not-allowed; }

    .btn--secondary {
      background: #f1f5f9;
      color: #475569;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn--secondary:hover { background: #e2e8f0; }

    .btn--danger {
      background: #dc2626;
      color: white;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn--danger:hover:not(:disabled) { background: #b91c1c; }
    .btn--danger:disabled { opacity: 0.5; cursor: not-allowed; }

    /* Form Styles */
    .form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    .form-section {
      padding-bottom: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .form-section:last-of-type {
      border-bottom: none;
      padding-bottom: 0;
    }

    .form-section-title {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin: 0 0 1rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
    }

    .form-section-title i {
      color: #f59e0b;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
      margin-bottom: 1rem;
    }

    .form-group:last-child {
      margin-bottom: 0;
    }

    .form-group label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .required {
      color: #dc2626;
    }

    .form-input,
    .form-textarea {
      width: 100%;
      padding: 0.75rem 1rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      background: #f8fafc;
      color: #0f172a;
      transition: all 0.2s ease;
    }

    .form-input:focus,
    .form-textarea:focus {
      outline: none;
      border-color: #f59e0b;
      background: white;
      box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.1);
    }

    .form-textarea {
      resize: vertical;
      min-height: 80px;
    }

    .form-hint {
      margin: 0.25rem 0 0;
      font-size: 0.75rem;
      color: #64748b;
    }

    /* Delete Confirmation */
    .delete-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .delete-icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      border-radius: 50%;
      background: #fef2f2;
      color: #dc2626;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
    }

    .delete-confirmation p {
      margin: 0 0 0.5rem;
      color: #374151;
    }

    .delete-warning {
      font-size: 0.875rem;
      color: #64748b;
    }

    /* Responsive Styles */
    @media (max-width: 768px) {
      .page {
        padding: 1rem;
      }

      .page-header {
        flex-direction: column;
        align-items: stretch;
      }

      .header-actions {
        justify-content: flex-end;
      }

      .btn-text {
        display: none;
      }

      .templates-grid {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class HallTicketTemplateComponent implements OnInit {
  private examService = inject(ExamService);
  private toastService = inject(ToastService);

  // Data signals
  templates = signal<HallTicketTemplate[]>([]);

  // State signals
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);

  // Modal state
  showFormModal = signal(false);
  showDeleteModal = signal(false);
  editingTemplate = signal<HallTicketTemplate | null>(null);
  templateToDelete = signal<HallTicketTemplate | null>(null);

  // Form data
  formData = this.getEmptyFormData();

  // Computed values
  modalTitle = computed(() => {
    const editing = this.editingTemplate();
    return editing ? `Edit: ${editing.name}` : 'Create New Template';
  });

  ngOnInit(): void {
    this.loadTemplates();
  }

  loadTemplates(): void {
    this.loading.set(true);
    this.error.set(null);

    this.examService.getHallTicketTemplates().subscribe({
      next: templates => {
        this.templates.set(templates);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load templates. Please try again.');
        this.loading.set(false);
      },
    });
  }

  openCreateModal(): void {
    this.editingTemplate.set(null);
    this.formData = this.getEmptyFormData();
    this.showFormModal.set(true);
  }

  editTemplate(template: HallTicketTemplate): void {
    this.editingTemplate.set(template);
    this.formData = {
      name: template.name,
      headerLogoUrl: template.headerLogoUrl || '',
      schoolName: template.schoolName || '',
      schoolAddress: template.schoolAddress || '',
      instructions: template.instructions || '',
      isDefault: template.isDefault,
    };
    this.showFormModal.set(true);
  }

  closeFormModal(): void {
    this.showFormModal.set(false);
    this.editingTemplate.set(null);
  }

  isFormValid(): boolean {
    return !!this.formData.name;
  }

  saveTemplate(): void {
    if (!this.isFormValid()) {
      this.toastService.error('Please enter a template name');
      return;
    }

    this.saving.set(true);
    const editing = this.editingTemplate();

    const operation = editing
      ? this.examService.updateHallTicketTemplate(editing.id, {
          name: this.formData.name,
          headerLogoUrl: this.formData.headerLogoUrl || undefined,
          schoolName: this.formData.schoolName || undefined,
          schoolAddress: this.formData.schoolAddress || undefined,
          instructions: this.formData.instructions || undefined,
          isDefault: this.formData.isDefault,
        })
      : this.examService.createHallTicketTemplate({
          name: this.formData.name,
          headerLogoUrl: this.formData.headerLogoUrl || undefined,
          schoolName: this.formData.schoolName || undefined,
          schoolAddress: this.formData.schoolAddress || undefined,
          instructions: this.formData.instructions || undefined,
          isDefault: this.formData.isDefault,
        });

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Template updated successfully' : 'Template created successfully'
        );
        this.closeFormModal();
        this.loadTemplates();
        this.saving.set(false);
      },
      error: (err) => {
        const message = err?.error?.error || (editing ? 'Failed to update template' : 'Failed to create template');
        this.toastService.error(message);
        this.saving.set(false);
      },
    });
  }

  setAsDefault(template: HallTicketTemplate): void {
    this.examService.updateHallTicketTemplate(template.id, { isDefault: true }).subscribe({
      next: () => {
        this.toastService.success('Template set as default');
        this.loadTemplates();
      },
      error: (err) => {
        const message = err?.error?.error || 'Failed to set default template';
        this.toastService.error(message);
      },
    });
  }

  confirmDelete(template: HallTicketTemplate): void {
    this.templateToDelete.set(template);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.templateToDelete.set(null);
  }

  deleteTemplate(): void {
    const template = this.templateToDelete();
    if (!template) return;

    this.deleting.set(true);

    this.examService.deleteHallTicketTemplate(template.id).subscribe({
      next: () => {
        this.toastService.success('Template deleted successfully');
        this.closeDeleteModal();
        this.loadTemplates();
        this.deleting.set(false);
      },
      error: (err) => {
        const message = err?.error?.error || 'Failed to delete template';
        this.toastService.error(message);
        this.deleting.set(false);
      },
    });
  }

  truncate(text: string, maxLength: number): string {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
  }

  private getEmptyFormData() {
    return {
      name: '',
      headerLogoUrl: '',
      schoolName: '',
      schoolAddress: '',
      instructions: '',
      isDefault: false,
    };
  }
}
