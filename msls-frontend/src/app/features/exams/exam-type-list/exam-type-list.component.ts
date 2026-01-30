import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';

import { ExamService } from '../exam.service';
import { ExamType, EVALUATION_TYPES } from '../exam.model';

@Component({
  selector: 'app-exam-type-list',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  template: `
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Exam Types</h1>
          <p class="text-gray-500 mt-1">Configure different types of examinations</p>
        </div>
        <div class="flex items-center gap-2">
          <button (click)="showHelp.set(!showHelp())"
                  class="inline-flex items-center px-3 py-2 text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors">
            <i class="fa-solid fa-circle-question mr-2"></i>
            How it Works
          </button>
          <a routerLink="new"
             class="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
            <i class="fa-solid fa-plus mr-2"></i>
            Add Exam Type
          </a>
        </div>
      </div>

      <!-- How it Works Section -->
      @if (showHelp()) {
        <div class="bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-lg p-6">
          <div class="flex items-start justify-between">
            <h2 class="text-lg font-semibold text-blue-900 mb-4">
              <i class="fa-solid fa-graduation-cap mr-2"></i>
              Understanding Exam Types
            </h2>
            <button (click)="showHelp.set(false)" class="text-blue-400 hover:text-blue-600">
              <i class="fa-solid fa-times"></i>
            </button>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- What are Exam Types -->
            <div class="bg-white rounded-lg p-4 shadow-sm">
              <h3 class="font-medium text-gray-900 mb-2">
                <i class="fa-solid fa-clipboard-list text-blue-500 mr-2"></i>
                What are Exam Types?
              </h3>
              <p class="text-sm text-gray-600">
                Exam types define the different categories of assessments in your school, such as Unit Tests,
                Mid-Term Exams, Final Exams, Practicals, and Projects. Each type can have its own evaluation
                method and contribution to the final grade.
              </p>
            </div>

            <!-- Evaluation Types -->
            <div class="bg-white rounded-lg p-4 shadow-sm">
              <h3 class="font-medium text-gray-900 mb-2">
                <i class="fa-solid fa-star text-yellow-500 mr-2"></i>
                Evaluation Types
              </h3>
              <ul class="text-sm text-gray-600 space-y-2">
                <li>
                  <span class="font-medium text-blue-700">Marks-based:</span>
                  Students receive numeric scores (e.g., 85/100). Best for traditional exams with quantifiable answers.
                </li>
                <li>
                  <span class="font-medium text-green-700">Grade-based:</span>
                  Students receive letter grades (A+, A, B, etc.). Best for subjective assessments like projects or presentations.
                </li>
              </ul>
            </div>

            <!-- Weightage -->
            <div class="bg-white rounded-lg p-4 shadow-sm">
              <h3 class="font-medium text-gray-900 mb-2">
                <i class="fa-solid fa-scale-balanced text-purple-500 mr-2"></i>
                Weightage Explained
              </h3>
              <p class="text-sm text-gray-600 mb-2">
                Weightage determines how much each exam type contributes to the final result calculation.
              </p>
              <div class="bg-gray-50 rounded p-2 text-xs text-gray-500">
                <strong>Example:</strong> If Unit Tests = 20%, Mid-Term = 30%, Final = 50%, a student scoring
                80 in Unit Test, 70 in Mid-Term, and 90 in Final gets: (80×0.20) + (70×0.30) + (90×0.50) = 82
              </div>
            </div>

            <!-- Display Order & Status -->
            <div class="bg-white rounded-lg p-4 shadow-sm">
              <h3 class="font-medium text-gray-900 mb-2">
                <i class="fa-solid fa-list-ol text-orange-500 mr-2"></i>
                Order & Status
              </h3>
              <ul class="text-sm text-gray-600 space-y-2">
                <li>
                  <span class="font-medium">Display Order:</span>
                  Use the up/down arrows to arrange exam types in the order they appear in reports and schedules.
                </li>
                <li>
                  <span class="font-medium">Active/Inactive:</span>
                  Toggle to hide exam types from selection without deleting them. Useful for seasonal exams.
                </li>
              </ul>
            </div>
          </div>

          <!-- Quick Tips -->
          <div class="mt-4 p-3 bg-amber-50 border border-amber-200 rounded-lg">
            <h4 class="font-medium text-amber-800 text-sm mb-2">
              <i class="fa-solid fa-lightbulb mr-1"></i>
              Quick Tips for Administrators
            </h4>
            <ul class="text-xs text-amber-700 space-y-1">
              <li><i class="fa-solid fa-check mr-1"></i> Ensure total weightage across all active exam types equals 100% for accurate result calculation</li>
              <li><i class="fa-solid fa-check mr-1"></i> Use unique codes (UT, MT, FE) for easy identification in reports and data exports</li>
              <li><i class="fa-solid fa-check mr-1"></i> Set appropriate default max marks to avoid teachers entering incorrect values</li>
              <li><i class="fa-solid fa-check mr-1"></i> Deactivate rather than delete exam types if historical data needs to be preserved</li>
            </ul>
          </div>
        </div>
      }

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
            <input type="text"
                   [(ngModel)]="searchQuery"
                   (ngModelChange)="onSearchChange()"
                   placeholder="Search by name or code..."
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          </div>
          <div class="w-full sm:w-48">
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select [(ngModel)]="activeFilter"
                    (ngModelChange)="onFilterChange()"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
              <option [ngValue]="null">All</option>
              <option [ngValue]="true">Active</option>
              <option [ngValue]="false">Inactive</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      @if (loading()) {
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      }

      <!-- Empty State -->
      @if (!loading() && examTypes().length === 0) {
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
          <i class="fa-solid fa-clipboard-list text-4xl text-gray-300 mb-4"></i>
          <h3 class="text-lg font-medium text-gray-900 mb-2">No exam types found</h3>
          <p class="text-gray-500 mb-4">Create your first exam type to get started</p>
          <a routerLink="new"
             class="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
            <i class="fa-solid fa-plus mr-2"></i>
            Add Exam Type
          </a>
        </div>
      }

      <!-- Exam Types List -->
      @if (!loading() && examTypes().length > 0) {
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <div class="px-4 py-3 bg-gray-50 border-b border-gray-200">
            <p class="text-sm text-gray-600">
              <i class="fa-solid fa-list-ol mr-2 text-gray-400"></i>
              {{ examTypes().length }} exam types configured
            </p>
          </div>

          <div class="divide-y divide-gray-200">
            @for (examType of examTypes(); track examType.id; let i = $index; let first = $first; let last = $last) {
              <div class="flex items-center px-4 py-4 hover:bg-gray-50 transition-colors">
                <!-- Move Buttons -->
                <div class="flex flex-col mr-4">
                  <button (click)="moveUp(i)"
                          [disabled]="first"
                          class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed">
                    <i class="fa-solid fa-chevron-up text-xs"></i>
                  </button>
                  <button (click)="moveDown(i)"
                          [disabled]="last"
                          class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed">
                    <i class="fa-solid fa-chevron-down text-xs"></i>
                  </button>
                </div>

                <!-- Exam Type Info -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <h3 class="text-sm font-medium text-gray-900 truncate">{{ examType.name }}</h3>
                    <span class="px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-700 rounded">
                      {{ examType.code }}
                    </span>
                    @if (!examType.isActive) {
                      <span class="px-2 py-0.5 text-xs font-medium bg-red-100 text-red-700 rounded">
                        Inactive
                      </span>
                    }
                  </div>
                  <p class="text-sm text-gray-500 mt-1">
                    <span [class]="getEvaluationTypeColor(examType.evaluationType)">
                      {{ getEvaluationTypeLabel(examType.evaluationType) }}
                    </span>
                    <span class="mx-2">•</span>
                    Max Marks: {{ examType.defaultMaxMarks }}
                    <span class="mx-2">•</span>
                    Weightage: {{ examType.weightage }}%
                  </p>
                </div>

                <!-- Actions -->
                <div class="flex items-center gap-2 ml-4">
                  <button (click)="toggleActive(examType)"
                          [class]="examType.isActive ? 'text-green-600 hover:text-green-700' : 'text-gray-400 hover:text-gray-600'"
                          class="p-2 rounded-lg hover:bg-gray-100 transition-colors"
                          [title]="examType.isActive ? 'Deactivate' : 'Activate'">
                    <i [class]="examType.isActive ? 'fa-solid fa-toggle-on text-xl' : 'fa-solid fa-toggle-off text-xl'"></i>
                  </button>
                  <a [routerLink]="[examType.id, 'edit']"
                     class="p-2 text-gray-400 hover:text-blue-600 rounded-lg hover:bg-gray-100 transition-colors"
                     title="Edit">
                    <i class="fa-solid fa-pen"></i>
                  </a>
                  <button (click)="confirmDelete(examType)"
                          class="p-2 text-gray-400 hover:text-red-600 rounded-lg hover:bg-gray-100 transition-colors"
                          title="Delete">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </div>
              </div>
            }
          </div>
        </div>
      }

      <!-- Delete Confirmation Modal -->
      @if (deleteTarget()) {
        <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" (click)="cancelDelete()">
          <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6" (click)="$event.stopPropagation()">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Delete Exam Type</h3>
            <p class="text-gray-500 mb-4">
              Are you sure you want to delete <strong>{{ deleteTarget()!.name }}</strong>?
              This action cannot be undone.
            </p>
            <div class="flex justify-end gap-3">
              <button (click)="cancelDelete()"
                      class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors">
                Cancel
              </button>
              <button (click)="deleteExamType()"
                      class="px-4 py-2 text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors">
                Delete
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
})
export class ExamTypeListComponent implements OnInit {
  private readonly examService = inject(ExamService);

  examTypes = signal<ExamType[]>([]);
  loading = signal(true);
  searchQuery = '';
  activeFilter: boolean | null = null;
  deleteTarget = signal<ExamType | null>(null);
  showHelp = signal(false);

  ngOnInit(): void {
    this.loadExamTypes();
  }

  loadExamTypes(): void {
    this.loading.set(true);
    this.examService.getExamTypes({
      search: this.searchQuery || undefined,
      isActive: this.activeFilter ?? undefined,
    }).subscribe({
      next: (types) => {
        this.examTypes.set(types);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load exam types:', err);
        this.loading.set(false);
      },
    });
  }

  onSearchChange(): void {
    this.loadExamTypes();
  }

  onFilterChange(): void {
    this.loadExamTypes();
  }

  moveUp(index: number): void {
    if (index <= 0) return;
    this.swapItems(index, index - 1);
  }

  moveDown(index: number): void {
    if (index >= this.examTypes().length - 1) return;
    this.swapItems(index, index + 1);
  }

  private swapItems(fromIndex: number, toIndex: number): void {
    const items = [...this.examTypes()];
    [items[fromIndex], items[toIndex]] = [items[toIndex], items[fromIndex]];

    // Update display order
    const orderUpdates = items.map((item, index) => ({
      id: item.id,
      displayOrder: index + 1,
    }));

    this.examTypes.set(items);

    this.examService.updateDisplayOrder({ items: orderUpdates }).subscribe({
      error: (err) => {
        console.error('Failed to update display order:', err);
        this.loadExamTypes(); // Reload on error
      },
    });
  }

  toggleActive(examType: ExamType): void {
    this.examService.toggleExamTypeActive(examType.id, !examType.isActive).subscribe({
      next: () => {
        const types = this.examTypes().map(t =>
          t.id === examType.id ? { ...t, isActive: !t.isActive } : t
        );
        this.examTypes.set(types);
      },
      error: (err) => {
        console.error('Failed to toggle status:', err);
      },
    });
  }

  confirmDelete(examType: ExamType): void {
    this.deleteTarget.set(examType);
  }

  cancelDelete(): void {
    this.deleteTarget.set(null);
  }

  deleteExamType(): void {
    const target = this.deleteTarget();
    if (!target) return;

    this.examService.deleteExamType(target.id).subscribe({
      next: () => {
        this.examTypes.set(this.examTypes().filter(t => t.id !== target.id));
        this.deleteTarget.set(null);
      },
      error: (err) => {
        console.error('Failed to delete exam type:', err);
        this.deleteTarget.set(null);
      },
    });
  }

  getEvaluationTypeLabel(type: string): string {
    const found = EVALUATION_TYPES.find(t => t.value === type);
    return found?.label || type;
  }

  getEvaluationTypeColor(type: string): string {
    const found = EVALUATION_TYPES.find(t => t.value === type);
    return found?.color || 'bg-gray-100 text-gray-700';
  }
}
