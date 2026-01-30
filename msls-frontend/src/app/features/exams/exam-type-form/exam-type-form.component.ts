import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute, RouterModule } from '@angular/router';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';

import { ExamService } from '../exam.service';
import { ExamType, CreateExamTypeRequest, EVALUATION_TYPES, EXAM_TYPE_PRESETS } from '../exam.model';

@Component({
  selector: 'app-exam-type-form',
  standalone: true,
  imports: [CommonModule, RouterModule, ReactiveFormsModule],
  template: `
    <div class="max-w-2xl mx-auto space-y-6">
      <!-- Header -->
      <div class="flex items-center gap-4">
        <a routerLink=".."
           class="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors">
          <i class="fa-solid fa-arrow-left"></i>
        </a>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">
            {{ isEditMode() ? 'Edit Exam Type' : 'Create Exam Type' }}
          </h1>
          <p class="text-gray-500 mt-1">
            {{ isEditMode() ? 'Update exam type configuration' : 'Configure a new type of examination' }}
          </p>
        </div>
      </div>

      <!-- Quick Presets (Create mode only) -->
      @if (!isEditMode()) {
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h3 class="text-sm font-medium text-blue-800 mb-2">
            <i class="fa-solid fa-magic-wand-sparkles mr-2"></i>
            Quick Presets
          </h3>
          <div class="flex flex-wrap gap-2">
            @for (preset of presets; track preset.code) {
              <button type="button"
                      (click)="applyPreset(preset)"
                      class="px-3 py-1 text-sm bg-white border border-blue-200 rounded-full text-blue-700 hover:bg-blue-100 transition-colors">
                {{ preset.name }}
              </button>
            }
          </div>
        </div>
      }

      <!-- Form -->
      <form [formGroup]="form" (ngSubmit)="onSubmit()" class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 space-y-6">
        <!-- Name & Code -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Name <span class="text-red-500">*</span>
            </label>
            <input type="text"
                   formControlName="name"
                   placeholder="e.g., Unit Test"
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                   [class.border-red-500]="form.get('name')?.invalid && form.get('name')?.touched">
            @if (form.get('name')?.errors?.['required'] && form.get('name')?.touched) {
              <p class="text-red-500 text-sm mt-1">Name is required</p>
            }
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Code <span class="text-red-500">*</span>
            </label>
            <input type="text"
                   formControlName="code"
                   placeholder="e.g., UT"
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 uppercase"
                   [class.border-red-500]="form.get('code')?.invalid && form.get('code')?.touched">
            @if (form.get('code')?.errors?.['required'] && form.get('code')?.touched) {
              <p class="text-red-500 text-sm mt-1">Code is required</p>
            }
          </div>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea formControlName="description"
                    rows="2"
                    placeholder="Optional description..."
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          </textarea>
        </div>

        <!-- Evaluation Type -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Evaluation Type <span class="text-red-500">*</span>
          </label>
          <div class="flex gap-4">
            @for (evalType of evaluationTypes; track evalType.value) {
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio"
                       formControlName="evaluationType"
                       [value]="evalType.value"
                       class="w-4 h-4 text-blue-600 focus:ring-blue-500">
                <span [class]="evalType.color + ' px-2 py-1 rounded text-sm'">
                  <i [class]="'fa-solid ' + evalType.icon + ' mr-1'"></i>
                  {{ evalType.label }}
                </span>
              </label>
            }
          </div>
        </div>

        <!-- Max Marks & Passing Marks -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Default Max Marks <span class="text-red-500">*</span>
            </label>
            <input type="number"
                   formControlName="defaultMaxMarks"
                   min="1"
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                   [class.border-red-500]="form.get('defaultMaxMarks')?.invalid && form.get('defaultMaxMarks')?.touched">
            @if (form.get('defaultMaxMarks')?.errors?.['required'] && form.get('defaultMaxMarks')?.touched) {
              <p class="text-red-500 text-sm mt-1">Max marks is required</p>
            }
            @if (form.get('defaultMaxMarks')?.errors?.['min']) {
              <p class="text-red-500 text-sm mt-1">Must be greater than 0</p>
            }
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Default Passing Marks</label>
            <input type="number"
                   formControlName="defaultPassingMarks"
                   min="0"
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          </div>
        </div>

        <!-- Weightage -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Weightage (%) <span class="text-red-500">*</span>
          </label>
          <div class="flex items-center gap-4">
            <input type="range"
                   formControlName="weightage"
                   min="0"
                   max="100"
                   step="5"
                   class="flex-1 h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer">
            <span class="w-16 text-center font-medium text-gray-900">
              {{ form.get('weightage')?.value }}%
            </span>
          </div>
          <p class="text-sm text-gray-500 mt-1">Weightage used for final result calculation</p>
        </div>

        <!-- Submit Buttons -->
        <div class="flex justify-end gap-3 pt-4 border-t border-gray-200">
          <a routerLink=".."
             class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors">
            Cancel
          </a>
          <button type="submit"
                  [disabled]="form.invalid || submitting()"
                  class="px-4 py-2 text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
            @if (submitting()) {
              <i class="fa-solid fa-spinner fa-spin mr-2"></i>
            }
            {{ isEditMode() ? 'Update' : 'Create' }} Exam Type
          </button>
        </div>

        <!-- Error Message -->
        @if (errorMessage()) {
          <div class="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
            <i class="fa-solid fa-exclamation-circle mr-2"></i>
            {{ errorMessage() }}
          </div>
        }
      </form>
    </div>
  `,
})
export class ExamTypeFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly examService = inject(ExamService);

  form!: FormGroup;
  evaluationTypes = EVALUATION_TYPES;
  presets = EXAM_TYPE_PRESETS;
  isEditMode = signal(false);
  submitting = signal(false);
  errorMessage = signal('');
  private examTypeId: string | null = null;

  ngOnInit(): void {
    this.initForm();

    this.examTypeId = this.route.snapshot.paramMap.get('id');
    if (this.examTypeId) {
      this.isEditMode.set(true);
      this.loadExamType(this.examTypeId);
    }
  }

  private initForm(): void {
    this.form = this.fb.group({
      name: ['', [Validators.required, Validators.maxLength(100)]],
      code: ['', [Validators.required, Validators.maxLength(20)]],
      description: [''],
      evaluationType: ['marks', Validators.required],
      defaultMaxMarks: [100, [Validators.required, Validators.min(1)]],
      defaultPassingMarks: [null],
      weightage: [0, [Validators.required, Validators.min(0), Validators.max(100)]],
    });
  }

  private loadExamType(id: string): void {
    this.examService.getExamType(id).subscribe({
      next: (examType) => {
        this.form.patchValue({
          name: examType.name,
          code: examType.code,
          description: examType.description,
          evaluationType: examType.evaluationType,
          defaultMaxMarks: examType.defaultMaxMarks,
          defaultPassingMarks: examType.defaultPassingMarks,
          weightage: examType.weightage,
        });
      },
      error: (err) => {
        console.error('Failed to load exam type:', err);
        this.errorMessage.set('Failed to load exam type');
      },
    });
  }

  applyPreset(preset: Partial<CreateExamTypeRequest>): void {
    this.form.patchValue({
      name: preset.name,
      code: preset.code,
      evaluationType: preset.evaluationType,
      defaultMaxMarks: preset.defaultMaxMarks,
      weightage: preset.weightage,
    });
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.submitting.set(true);
    this.errorMessage.set('');

    const formValue = this.form.value;
    const request: CreateExamTypeRequest = {
      name: formValue.name,
      code: formValue.code.toUpperCase(),
      description: formValue.description || undefined,
      evaluationType: formValue.evaluationType,
      defaultMaxMarks: formValue.defaultMaxMarks,
      defaultPassingMarks: formValue.defaultPassingMarks || undefined,
      weightage: formValue.weightage,
    };

    const operation = this.isEditMode()
      ? this.examService.updateExamType(this.examTypeId!, request)
      : this.examService.createExamType(request);

    operation.subscribe({
      next: () => {
        this.router.navigate(['..'], { relativeTo: this.route });
      },
      error: (err) => {
        console.error('Failed to save exam type:', err);
        this.errorMessage.set(err.error?.detail || 'Failed to save exam type');
        this.submitting.set(false);
      },
    });
  }
}
