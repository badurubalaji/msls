/**
 * Document Types Management Component
 * Story 5.8: Staff Document Management
 */
import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { StaffDocumentService } from '../../staff-document.service';
import {
  DocumentType,
  CreateDocumentTypeRequest,
  UpdateDocumentTypeRequest,
  DOCUMENT_CATEGORIES,
  STAFF_TYPES,
  DocumentCategory,
  StaffType,
  getCategoryLabel,
} from '../../staff-document.model';
import { MslsIconComponent } from '../../../../shared/components/icon/icon.component';

@Component({
  selector: 'app-document-types',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsIconComponent],
  template: `
    <div class="container mx-auto p-6">
      <!-- Header -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Document Types</h1>
          <p class="text-gray-600">Configure document types for staff documents</p>
        </div>
        <button
          (click)="openCreateModal()"
          class="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
        >
          <msls-icon name="plus" class="w-5 h-5" />
          Add Document Type
        </button>
      </div>

      <!-- Loading State -->
      @if (loading()) {
        <div class="flex items-center justify-center h-64">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        </div>
      }

      <!-- Document Types List -->
      @if (!loading() && documentTypes().length > 0) {
        <div class="bg-white shadow rounded-lg overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Code</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Mandatory</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Has Expiry</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              @for (docType of documentTypes(); track docType.id) {
                <tr class="hover:bg-gray-50">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="font-medium text-gray-900">{{ docType.name }}</div>
                    @if (docType.description) {
                      <div class="text-sm text-gray-500">{{ docType.description }}</div>
                    }
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <code class="bg-gray-100 px-2 py-1 rounded">{{ docType.code }}</code>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <span class="px-2 py-1 text-xs font-medium rounded-full"
                          [class]="getCategoryClass(docType.category)">
                      {{ getCategoryLabel(docType.category) }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    @if (docType.is_mandatory) {
                      <span class="text-green-600">Yes</span>
                    } @else {
                      <span class="text-gray-400">No</span>
                    }
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    @if (docType.has_expiry) {
                      <span class="text-amber-600">
                        Yes
                        @if (docType.default_validity_months) {
                          ({{ docType.default_validity_months }} months)
                        }
                      </span>
                    } @else {
                      <span class="text-gray-400">No</span>
                    }
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    @if (docType.is_active) {
                      <span class="px-2 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800">Active</span>
                    } @else {
                      <span class="px-2 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-800">Inactive</span>
                    }
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button (click)="openEditModal(docType)" class="text-indigo-600 hover:text-indigo-900 mr-3">
                      Edit
                    </button>
                    <button (click)="confirmDelete(docType)" class="text-red-600 hover:text-red-900">
                      Delete
                    </button>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
      }

      <!-- Empty State -->
      @if (!loading() && documentTypes().length === 0) {
        <div class="bg-white shadow rounded-lg p-12 text-center">
          <msls-icon name="document-text" class="w-12 h-12 mx-auto text-gray-400 mb-4" />
          <h3 class="text-lg font-medium text-gray-900 mb-2">No document types</h3>
          <p class="text-gray-500 mb-4">Get started by creating a new document type.</p>
          <button
            (click)="openCreateModal()"
            class="inline-flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
          >
            <msls-icon name="plus" class="w-5 h-5" />
            Add Document Type
          </button>
        </div>
      }

      <!-- Create/Edit Modal -->
      @if (showModal()) {
        <div class="fixed inset-0 z-50 overflow-y-auto">
          <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
            <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" (click)="closeModal()"></div>
            <div class="relative bg-white rounded-lg shadow-xl max-w-lg w-full mx-auto z-10">
              <div class="px-6 py-4 border-b">
                <h3 class="text-lg font-medium text-gray-900">
                  {{ editingType() ? 'Edit Document Type' : 'Create Document Type' }}
                </h3>
              </div>
              <form (ngSubmit)="saveDocumentType()" class="px-6 py-4 space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Name *</label>
                  <input
                    type="text"
                    [(ngModel)]="formData.name"
                    name="name"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
                    placeholder="e.g., Aadhaar Card"
                  />
                </div>
                @if (!editingType()) {
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Code *</label>
                    <input
                      type="text"
                      [(ngModel)]="formData.code"
                      name="code"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
                      placeholder="e.g., aadhaar"
                    />
                  </div>
                }
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Category *</label>
                  <select
                    [(ngModel)]="formData.category"
                    name="category"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
                  >
                    @for (cat of categories; track cat.value) {
                      <option [value]="cat.value">{{ cat.label }}</option>
                    }
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                  <textarea
                    [(ngModel)]="formData.description"
                    name="description"
                    rows="2"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
                    placeholder="Optional description"
                  ></textarea>
                </div>
                <div class="grid grid-cols-2 gap-4">
                  <div class="flex items-center">
                    <input
                      type="checkbox"
                      [(ngModel)]="formData.is_mandatory"
                      name="is_mandatory"
                      id="is_mandatory"
                      class="h-4 w-4 text-indigo-600 border-gray-300 rounded"
                    />
                    <label for="is_mandatory" class="ml-2 text-sm text-gray-700">Mandatory</label>
                  </div>
                  <div class="flex items-center">
                    <input
                      type="checkbox"
                      [(ngModel)]="formData.has_expiry"
                      name="has_expiry"
                      id="has_expiry"
                      class="h-4 w-4 text-indigo-600 border-gray-300 rounded"
                    />
                    <label for="has_expiry" class="ml-2 text-sm text-gray-700">Has Expiry</label>
                  </div>
                </div>
                @if (formData.has_expiry) {
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Default Validity (months)</label>
                    <input
                      type="number"
                      [(ngModel)]="formData.default_validity_months"
                      name="default_validity_months"
                      min="1"
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
                      placeholder="e.g., 12"
                    />
                  </div>
                }
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Applicable To</label>
                  <div class="flex gap-4">
                    @for (type of staffTypes; track type.value) {
                      <label class="flex items-center">
                        <input
                          type="checkbox"
                          [checked]="formData.applicable_to?.includes(type.value)"
                          (change)="toggleApplicableTo(type.value)"
                          class="h-4 w-4 text-indigo-600 border-gray-300 rounded"
                        />
                        <span class="ml-2 text-sm text-gray-700">{{ type.label }}</span>
                      </label>
                    }
                  </div>
                </div>
                @if (editingType()) {
                  <div class="flex items-center">
                    <input
                      type="checkbox"
                      [(ngModel)]="formData.is_active"
                      name="is_active"
                      id="is_active"
                      class="h-4 w-4 text-indigo-600 border-gray-300 rounded"
                    />
                    <label for="is_active" class="ml-2 text-sm text-gray-700">Active</label>
                  </div>
                }
              </form>
              <div class="px-6 py-4 border-t flex justify-end gap-3">
                <button
                  type="button"
                  (click)="closeModal()"
                  class="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  (click)="saveDocumentType()"
                  [disabled]="saving()"
                  class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50"
                >
                  {{ saving() ? 'Saving...' : 'Save' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      }

      <!-- Delete Confirmation Modal -->
      @if (deleteTarget()) {
        <div class="fixed inset-0 z-50 overflow-y-auto">
          <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
            <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" (click)="cancelDelete()"></div>
            <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-auto z-10 p-6">
              <msls-icon name="exclamation-triangle" class="w-12 h-12 text-red-600 mx-auto mb-4" />
              <h3 class="text-lg font-medium text-gray-900 mb-2">Delete Document Type</h3>
              <p class="text-gray-500 mb-6">
                Are you sure you want to delete "{{ deleteTarget()?.name }}"? This action cannot be undone.
              </p>
              <div class="flex justify-center gap-3">
                <button
                  (click)="cancelDelete()"
                  class="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  (click)="executeDelete()"
                  [disabled]="deleting()"
                  class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
                >
                  {{ deleting() ? 'Deleting...' : 'Delete' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      }
    </div>
  `,
})
export class DocumentTypesComponent implements OnInit {
  private readonly service = inject(StaffDocumentService);

  // State
  documentTypes = signal<DocumentType[]>([]);
  loading = signal(false);
  saving = signal(false);
  deleting = signal(false);
  showModal = signal(false);
  editingType = signal<DocumentType | null>(null);
  deleteTarget = signal<DocumentType | null>(null);

  // Form data
  formData: CreateDocumentTypeRequest & { is_active?: boolean } = {
    name: '',
    code: '',
    category: 'identity',
    description: '',
    is_mandatory: false,
    has_expiry: false,
    default_validity_months: undefined,
    applicable_to: ['teaching', 'non_teaching'],
    is_active: true,
  };

  // Constants
  categories = DOCUMENT_CATEGORIES;
  staffTypes = STAFF_TYPES;

  ngOnInit(): void {
    this.loadDocumentTypes();
  }

  loadDocumentTypes(): void {
    this.loading.set(true);
    this.service.getDocumentTypes().subscribe({
      next: (response) => {
        this.documentTypes.set(response.document_types);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      },
    });
  }

  getCategoryLabel(category: DocumentCategory): string {
    return getCategoryLabel(category);
  }

  getCategoryClass(category: DocumentCategory): string {
    const classes: Record<DocumentCategory, string> = {
      identity: 'bg-blue-100 text-blue-800',
      education: 'bg-purple-100 text-purple-800',
      employment: 'bg-green-100 text-green-800',
      compliance: 'bg-amber-100 text-amber-800',
      other: 'bg-gray-100 text-gray-800',
    };
    return classes[category] || classes.other;
  }

  openCreateModal(): void {
    this.editingType.set(null);
    this.formData = {
      name: '',
      code: '',
      category: 'identity',
      description: '',
      is_mandatory: false,
      has_expiry: false,
      default_validity_months: undefined,
      applicable_to: ['teaching', 'non_teaching'],
      is_active: true,
    };
    this.showModal.set(true);
  }

  openEditModal(docType: DocumentType): void {
    this.editingType.set(docType);
    this.formData = {
      name: docType.name,
      code: docType.code,
      category: docType.category,
      description: docType.description || '',
      is_mandatory: docType.is_mandatory,
      has_expiry: docType.has_expiry,
      default_validity_months: docType.default_validity_months,
      applicable_to: [...docType.applicable_to],
      is_active: docType.is_active,
    };
    this.showModal.set(true);
  }

  closeModal(): void {
    this.showModal.set(false);
    this.editingType.set(null);
  }

  toggleApplicableTo(type: StaffType): void {
    const current = this.formData.applicable_to || [];
    const index = current.indexOf(type);
    if (index === -1) {
      this.formData.applicable_to = [...current, type];
    } else {
      this.formData.applicable_to = current.filter((t) => t !== type);
    }
  }

  saveDocumentType(): void {
    if (!this.formData.name || !this.formData.category) return;

    this.saving.set(true);
    const editing = this.editingType();

    if (editing) {
      const updateData: UpdateDocumentTypeRequest = {
        name: this.formData.name,
        description: this.formData.description || undefined,
        is_mandatory: this.formData.is_mandatory,
        has_expiry: this.formData.has_expiry,
        default_validity_months: this.formData.has_expiry ? this.formData.default_validity_months : undefined,
        applicable_to: this.formData.applicable_to,
        is_active: this.formData.is_active,
      };

      this.service.updateDocumentType(editing.id, updateData).subscribe({
        next: () => {
          this.saving.set(false);
          this.closeModal();
          this.loadDocumentTypes();
        },
        error: () => {
          this.saving.set(false);
        },
      });
    } else {
      const createData: CreateDocumentTypeRequest = {
        name: this.formData.name,
        code: this.formData.code,
        category: this.formData.category,
        description: this.formData.description || undefined,
        is_mandatory: this.formData.is_mandatory,
        has_expiry: this.formData.has_expiry,
        default_validity_months: this.formData.has_expiry ? this.formData.default_validity_months : undefined,
        applicable_to: this.formData.applicable_to,
      };

      this.service.createDocumentType(createData).subscribe({
        next: () => {
          this.saving.set(false);
          this.closeModal();
          this.loadDocumentTypes();
        },
        error: () => {
          this.saving.set(false);
        },
      });
    }
  }

  confirmDelete(docType: DocumentType): void {
    this.deleteTarget.set(docType);
  }

  cancelDelete(): void {
    this.deleteTarget.set(null);
  }

  executeDelete(): void {
    const target = this.deleteTarget();
    if (!target) return;

    this.deleting.set(true);
    this.service.deleteDocumentType(target.id).subscribe({
      next: () => {
        this.deleting.set(false);
        this.deleteTarget.set(null);
        this.loadDocumentTypes();
      },
      error: () => {
        this.deleting.set(false);
      },
    });
  }
}
