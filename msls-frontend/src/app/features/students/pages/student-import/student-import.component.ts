/**
 * Student Bulk Import Component
 *
 * Allows bulk uploading of students via Excel or CSV file.
 * Provides template download and import result display.
 */

import {
  Component,
  ChangeDetectionStrategy,
  inject,
  signal,
  computed,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';

import { StudentService } from '../../services/student.service';
import { ImportResult, ImportError } from '../../models/student.model';
import { BranchService } from '../../../admin/branches/branch.service';
import { AcademicYearService } from '../../../admin/academic-years/academic-year.service';

interface Branch {
  id: string;
  name: string;
  code: string;
}

interface AcademicYear {
  id: string;
  name: string;
  isCurrent: boolean;
}

@Component({
  selector: 'msls-student-import',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <a routerLink="/students" class="back-link">
            <i class="fa-solid fa-arrow-left"></i>
          </a>
          <div class="header-text">
            <h1>Bulk Import Students</h1>
            <p>Upload an Excel or CSV file to add multiple students at once</p>
          </div>
        </div>
      </div>

      <!-- Steps Guide -->
      <div class="steps-guide">
        <div class="step" [class.active]="currentStep() === 1" [class.completed]="currentStep() > 1">
          <div class="step-number">1</div>
          <div class="step-content">
            <h4>Download Template</h4>
            <p>Get the Excel template with all required columns</p>
          </div>
        </div>
        <div class="step-connector"></div>
        <div class="step" [class.active]="currentStep() === 2" [class.completed]="currentStep() > 2">
          <div class="step-number">2</div>
          <div class="step-content">
            <h4>Fill Data</h4>
            <p>Add student information to the template</p>
          </div>
        </div>
        <div class="step-connector"></div>
        <div class="step" [class.active]="currentStep() === 3" [class.completed]="currentStep() > 3">
          <div class="step-number">3</div>
          <div class="step-content">
            <h4>Upload File</h4>
            <p>Select branch and upload the filled template</p>
          </div>
        </div>
        <div class="step-connector"></div>
        <div class="step" [class.active]="currentStep() === 4">
          <div class="step-number">4</div>
          <div class="step-content">
            <h4>Review Results</h4>
            <p>Check import status and fix any errors</p>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="content-grid">
        <!-- Template Download Card -->
        <div class="card template-card">
          <div class="card-icon">
            <i class="fa-solid fa-file-excel"></i>
          </div>
          <h3>Download Template</h3>
          <p>Start by downloading the Excel template. It includes all required columns and sample data.</p>
          <ul class="template-features">
            <li><i class="fa-solid fa-check"></i> Pre-formatted columns</li>
            <li><i class="fa-solid fa-check"></i> Sample data row</li>
            <li><i class="fa-solid fa-check"></i> Available classes &amp; sections reference</li>
            <li><i class="fa-solid fa-check"></i> Instructions sheet</li>
          </ul>
          <button class="btn btn-primary" (click)="downloadTemplate()">
            <i class="fa-solid fa-download"></i>
            Download Template
          </button>
        </div>

        <!-- Upload Card -->
        <div class="card upload-card">
          <h3>Upload Student File</h3>

          <!-- Configuration Form -->
          <div class="form-section">
            <div class="form-row">
              <div class="form-group">
                <label for="branch">Branch <span class="required">*</span></label>
                <select id="branch" [(ngModel)]="selectedBranch" class="form-control">
                  <option value="">Select Branch</option>
                  @for (branch of branches(); track branch.id) {
                    <option [value]="branch.id">{{ branch.name }}</option>
                  }
                </select>
              </div>
              <div class="form-group">
                <label for="academicYear">Academic Year <span class="required">*</span></label>
                <select id="academicYear" [(ngModel)]="selectedAcademicYear" class="form-control">
                  <option value="">Select Academic Year</option>
                  @for (year of academicYears(); track year.id) {
                    <option [value]="year.id">{{ year.name }} {{ year.isCurrent ? '(Current)' : '' }}</option>
                  }
                </select>
              </div>
            </div>
          </div>

          <!-- File Upload Area -->
          <div
            class="upload-area"
            [class.dragover]="isDragover()"
            [class.has-file]="selectedFile()"
            (dragover)="onDragOver($event)"
            (dragleave)="onDragLeave($event)"
            (drop)="onDrop($event)"
          >
            @if (selectedFile()) {
              <div class="file-preview">
                <i class="fa-solid fa-file-excel"></i>
                <div class="file-info">
                  <span class="file-name">{{ selectedFile()?.name }}</span>
                  <span class="file-size">{{ formatFileSize(selectedFile()?.size || 0) }}</span>
                </div>
                <button class="btn-remove" (click)="clearFile()" type="button">
                  <i class="fa-solid fa-times"></i>
                </button>
              </div>
            } @else {
              <i class="fa-solid fa-cloud-upload-alt upload-icon"></i>
              <p>Drag and drop your file here</p>
              <span class="divider">or</span>
              <label class="btn btn-secondary">
                <i class="fa-solid fa-folder-open"></i>
                Browse Files
                <input
                  type="file"
                  accept=".xlsx,.xls,.csv"
                  (change)="onFileSelect($event)"
                  hidden
                />
              </label>
              <span class="file-hint">Supports Excel (.xlsx, .xls) and CSV files (max 500 rows)</span>
            }
          </div>

          <!-- Upload Button -->
          <button
            class="btn btn-primary btn-upload"
            [disabled]="!canUpload()"
            (click)="uploadFile()"
          >
            @if (uploading()) {
              <i class="fa-solid fa-spinner fa-spin"></i>
              Importing Students...
            } @else {
              <i class="fa-solid fa-upload"></i>
              Import Students
            }
          </button>
        </div>
      </div>

      <!-- Import Results -->
      @if (importResult()) {
        <div class="results-section">
          <h2>Import Results</h2>
          <div class="results-summary">
            <div class="result-stat success">
              <i class="fa-solid fa-check-circle"></i>
              <div class="stat-content">
                <span class="stat-value">{{ importResult()!.successCount }}</span>
                <span class="stat-label">Successful</span>
              </div>
            </div>
            <div class="result-stat failed">
              <i class="fa-solid fa-times-circle"></i>
              <div class="stat-content">
                <span class="stat-value">{{ importResult()!.failedCount }}</span>
                <span class="stat-label">Failed</span>
              </div>
            </div>
            <div class="result-stat total">
              <i class="fa-solid fa-list"></i>
              <div class="stat-content">
                <span class="stat-value">{{ importResult()!.totalRows }}</span>
                <span class="stat-label">Total Rows</span>
              </div>
            </div>
          </div>

          @if (importResult()!.errors && importResult()!.errors!.length > 0) {
            <div class="errors-table">
              <h3>Errors</h3>
              <table>
                <thead>
                  <tr>
                    <th>Row</th>
                    <th>Column</th>
                    <th>Error</th>
                  </tr>
                </thead>
                <tbody>
                  @for (error of importResult()!.errors; track error.row + error.message) {
                    <tr>
                      <td>{{ error.row }}</td>
                      <td>{{ error.column || '-' }}</td>
                      <td>{{ error.message }}</td>
                    </tr>
                  }
                </tbody>
              </table>
            </div>
          }

          <div class="results-actions">
            @if (importResult()!.successCount > 0) {
              <a routerLink="/students" class="btn btn-primary">
                <i class="fa-solid fa-users"></i>
                View Students
              </a>
            }
            <button class="btn btn-secondary" (click)="resetImport()">
              <i class="fa-solid fa-redo"></i>
              Import More
            </button>
          </div>
        </div>
      }

      <!-- Instructions -->
      <div class="instructions-section">
        <h2>Import Instructions</h2>
        <div class="instructions-grid">
          <div class="instruction-item">
            <div class="instruction-icon"><i class="fa-solid fa-asterisk"></i></div>
            <div class="instruction-content">
              <h4>Required Fields</h4>
              <p>Admission Number, First Name, Last Name, Date of Birth, Gender, Class, Section, Guardian Name, and Guardian Phone are mandatory.</p>
            </div>
          </div>
          <div class="instruction-item">
            <div class="instruction-icon"><i class="fa-solid fa-calendar"></i></div>
            <div class="instruction-content">
              <h4>Date Format</h4>
              <p>Use YYYY-MM-DD format for dates (e.g., 2015-06-15 for June 15, 2015).</p>
            </div>
          </div>
          <div class="instruction-item">
            <div class="instruction-icon"><i class="fa-solid fa-venus-mars"></i></div>
            <div class="instruction-content">
              <h4>Gender Values</h4>
              <p>Use lowercase values: male, female, or other.</p>
            </div>
          </div>
          <div class="instruction-item">
            <div class="instruction-icon"><i class="fa-solid fa-school"></i></div>
            <div class="instruction-content">
              <h4>Class &amp; Section</h4>
              <p>Class and Section names must exactly match existing records in the system. Check the Reference sheet.</p>
            </div>
          </div>
          <div class="instruction-item">
            <div class="instruction-icon"><i class="fa-solid fa-id-card"></i></div>
            <div class="instruction-content">
              <h4>Admission Number</h4>
              <p>Each admission number must be unique. Duplicate numbers will be rejected.</p>
            </div>
          </div>
          <div class="instruction-item">
            <div class="instruction-icon"><i class="fa-solid fa-layer-group"></i></div>
            <div class="instruction-content">
              <h4>Maximum Rows</h4>
              <p>You can import up to 500 students at a time. Split larger imports into multiple files.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1200px; margin: 0 auto; }

    .page-header { margin-bottom: 1.5rem; }
    .header-content { display: flex; align-items: center; gap: 1rem; }
    .back-link { width: 2.5rem; height: 2.5rem; border-radius: 0.5rem; background: #f1f5f9; color: #64748b; display: flex; align-items: center; justify-content: center; text-decoration: none; transition: all 0.2s; }
    .back-link:hover { background: #e2e8f0; color: #1e293b; }
    .header-text h1 { margin: 0; font-size: 1.5rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }

    .steps-guide { display: flex; align-items: flex-start; gap: 0.5rem; padding: 1.25rem; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 1rem; margin-bottom: 1.5rem; overflow-x: auto; }
    .step { display: flex; align-items: flex-start; gap: 0.75rem; padding: 0.75rem; border-radius: 0.5rem; min-width: 180px; }
    .step.active { background: white; border: 1px solid #3b82f6; }
    .step.completed .step-number { background: #22c55e; }
    .step-number { width: 1.75rem; height: 1.75rem; border-radius: 50%; background: #e2e8f0; color: #64748b; font-weight: 600; font-size: 0.75rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
    .step.active .step-number { background: #3b82f6; color: white; }
    .step-content h4 { margin: 0; font-size: 0.8125rem; font-weight: 600; color: #1e293b; }
    .step-content p { margin: 0.125rem 0 0; font-size: 0.6875rem; color: #64748b; }
    .step-connector { width: 2rem; height: 2px; background: #e2e8f0; margin-top: 0.875rem; flex-shrink: 0; }

    .content-grid { display: grid; grid-template-columns: 1fr 1.5fr; gap: 1.5rem; margin-bottom: 1.5rem; }
    @media (max-width: 900px) { .content-grid { grid-template-columns: 1fr; } }

    .card { background: white; border: 1px solid #e2e8f0; border-radius: 1rem; padding: 1.5rem; }
    .template-card { display: flex; flex-direction: column; align-items: center; text-align: center; }
    .card-icon { width: 4rem; height: 4rem; border-radius: 1rem; background: linear-gradient(135deg, #22c55e, #16a34a); color: white; display: flex; align-items: center; justify-content: center; font-size: 1.75rem; margin-bottom: 1rem; }
    .card h3 { margin: 0 0 0.5rem; font-size: 1.125rem; font-weight: 600; color: #1e293b; }
    .card > p { margin: 0 0 1rem; color: #64748b; font-size: 0.875rem; line-height: 1.5; }
    .template-features { list-style: none; padding: 0; margin: 0 0 1.5rem; text-align: left; width: 100%; }
    .template-features li { display: flex; align-items: center; gap: 0.5rem; padding: 0.5rem 0; font-size: 0.8125rem; color: #475569; border-bottom: 1px solid #f1f5f9; }
    .template-features li:last-child { border-bottom: none; }
    .template-features i { color: #22c55e; font-size: 0.75rem; }

    .upload-card h3 { margin: 0 0 1rem; }
    .form-section { margin-bottom: 1.25rem; }
    .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    @media (max-width: 600px) { .form-row { grid-template-columns: 1fr; } }
    .form-group { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-group label { font-size: 0.8125rem; font-weight: 500; color: #374151; }
    .required { color: #ef4444; }
    .form-control { padding: 0.625rem 0.75rem; border: 1px solid #d1d5db; border-radius: 0.5rem; font-size: 0.875rem; background: white; transition: border-color 0.2s, box-shadow 0.2s; }
    .form-control:focus { outline: none; border-color: #3b82f6; box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1); }

    .upload-area { border: 2px dashed #d1d5db; border-radius: 0.75rem; padding: 2rem; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 0.75rem; text-align: center; transition: all 0.2s; cursor: pointer; min-height: 180px; }
    .upload-area:hover { border-color: #3b82f6; background: #f8fafc; }
    .upload-area.dragover { border-color: #3b82f6; background: #eff6ff; }
    .upload-area.has-file { border-style: solid; border-color: #22c55e; background: #f0fdf4; }
    .upload-icon { font-size: 2.5rem; color: #94a3b8; }
    .upload-area p { margin: 0; font-size: 0.9375rem; color: #64748b; }
    .divider { font-size: 0.75rem; color: #94a3b8; }
    .file-hint { font-size: 0.75rem; color: #94a3b8; margin-top: 0.5rem; }
    .file-preview { display: flex; align-items: center; gap: 1rem; padding: 0.75rem 1rem; background: white; border-radius: 0.5rem; width: 100%; max-width: 300px; }
    .file-preview i { font-size: 2rem; color: #22c55e; }
    .file-info { flex: 1; text-align: left; }
    .file-name { display: block; font-size: 0.875rem; font-weight: 500; color: #1e293b; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 180px; }
    .file-size { font-size: 0.75rem; color: #64748b; }
    .btn-remove { background: none; border: none; color: #ef4444; cursor: pointer; padding: 0.25rem; border-radius: 0.25rem; }
    .btn-remove:hover { background: #fef2f2; }

    .btn { display: inline-flex; align-items: center; gap: 0.5rem; padding: 0.625rem 1.25rem; border-radius: 0.5rem; font-size: 0.875rem; font-weight: 500; cursor: pointer; transition: all 0.2s; border: none; text-decoration: none; }
    .btn-primary { background: #3b82f6; color: white; }
    .btn-primary:hover:not(:disabled) { background: #2563eb; }
    .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
    .btn-secondary { background: #f1f5f9; color: #475569; border: 1px solid #e2e8f0; }
    .btn-secondary:hover { background: #e2e8f0; }
    .btn-upload { width: 100%; justify-content: center; margin-top: 1rem; padding: 0.875rem; }

    .results-section { background: white; border: 1px solid #e2e8f0; border-radius: 1rem; padding: 1.5rem; margin-bottom: 1.5rem; }
    .results-section h2 { margin: 0 0 1rem; font-size: 1.125rem; font-weight: 600; color: #1e293b; }
    .results-summary { display: flex; gap: 1rem; margin-bottom: 1.5rem; flex-wrap: wrap; }
    .result-stat { display: flex; align-items: center; gap: 0.75rem; padding: 1rem 1.25rem; border-radius: 0.75rem; min-width: 150px; }
    .result-stat.success { background: #f0fdf4; }
    .result-stat.success i { color: #22c55e; font-size: 1.5rem; }
    .result-stat.failed { background: #fef2f2; }
    .result-stat.failed i { color: #ef4444; font-size: 1.5rem; }
    .result-stat.total { background: #f1f5f9; }
    .result-stat.total i { color: #64748b; font-size: 1.5rem; }
    .stat-content { display: flex; flex-direction: column; }
    .stat-value { font-size: 1.5rem; font-weight: 700; color: #1e293b; }
    .stat-label { font-size: 0.75rem; color: #64748b; }

    .errors-table { margin-bottom: 1.5rem; }
    .errors-table h3 { margin: 0 0 0.75rem; font-size: 0.9375rem; font-weight: 600; color: #1e293b; }
    .errors-table table { width: 100%; border-collapse: collapse; font-size: 0.8125rem; }
    .errors-table th { text-align: left; padding: 0.625rem 0.75rem; background: #f8fafc; color: #64748b; font-weight: 500; border-bottom: 1px solid #e2e8f0; }
    .errors-table td { padding: 0.625rem 0.75rem; border-bottom: 1px solid #f1f5f9; color: #1e293b; }
    .errors-table tr:hover td { background: #fef2f2; }

    .results-actions { display: flex; gap: 0.75rem; }

    .instructions-section { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 1rem; padding: 1.5rem; }
    .instructions-section h2 { margin: 0 0 1.25rem; font-size: 1.125rem; font-weight: 600; color: #1e293b; }
    .instructions-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 1rem; }
    .instruction-item { display: flex; gap: 0.75rem; padding: 1rem; background: white; border-radius: 0.75rem; border: 1px solid #e2e8f0; }
    .instruction-icon { width: 2rem; height: 2rem; border-radius: 0.5rem; background: #eff6ff; color: #3b82f6; display: flex; align-items: center; justify-content: center; font-size: 0.875rem; flex-shrink: 0; }
    .instruction-content h4 { margin: 0 0 0.25rem; font-size: 0.875rem; font-weight: 600; color: #1e293b; }
    .instruction-content p { margin: 0; font-size: 0.75rem; color: #64748b; line-height: 1.4; }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StudentImportComponent {
  private studentService = inject(StudentService);
  private branchService = inject(BranchService);
  private academicYearService = inject(AcademicYearService);
  private router = inject(Router);

  // State
  readonly branches = signal<Branch[]>([]);
  readonly academicYears = signal<AcademicYear[]>([]);
  readonly selectedBranch = signal<string>('');
  readonly selectedAcademicYear = signal<string>('');
  readonly selectedFile = signal<File | null>(null);
  readonly isDragover = signal(false);
  readonly uploading = signal(false);
  readonly importResult = signal<ImportResult | null>(null);

  readonly currentStep = computed((): number => {
    if (this.importResult()) return 4;
    if (this.selectedFile()) return 3;
    if (this.selectedBranch() && this.selectedAcademicYear()) return 2;
    return 1;
  });

  readonly canUpload = computed(() => {
    return (
      this.selectedBranch() &&
      this.selectedAcademicYear() &&
      this.selectedFile() &&
      !this.uploading()
    );
  });

  constructor() {
    this.loadBranches();
    this.loadAcademicYears();
  }

  private loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: (branches) => this.branches.set(branches),
      error: () => console.error('Failed to load branches'),
    });
  }

  private loadAcademicYears(): void {
    this.academicYearService.getAcademicYears().subscribe({
      next: (years) => {
        this.academicYears.set(years);
        const current = years.find((y) => y.isCurrent);
        if (current) {
          this.selectedAcademicYear.set(current.id);
        }
      },
      error: () => console.error('Failed to load academic years'),
    });
  }

  downloadTemplate(): void {
    this.studentService.downloadImportTemplate();
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragover.set(true);
  }

  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragover.set(false);
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragover.set(false);

    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      this.handleFile(files[0]);
    }
  }

  onFileSelect(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.handleFile(input.files[0]);
      input.value = ''; // Reset to allow selecting same file again
    }
  }

  private handleFile(file: File): void {
    const validTypes = [
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'application/vnd.ms-excel',
      'text/csv',
    ];
    const validExtensions = ['.xlsx', '.xls', '.csv'];

    const hasValidType = validTypes.includes(file.type);
    const hasValidExtension = validExtensions.some((ext) =>
      file.name.toLowerCase().endsWith(ext)
    );

    if (!hasValidType && !hasValidExtension) {
      alert('Please select an Excel (.xlsx, .xls) or CSV file');
      return;
    }

    this.selectedFile.set(file);
    this.importResult.set(null);
  }

  clearFile(): void {
    this.selectedFile.set(null);
    this.importResult.set(null);
  }

  formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  uploadFile(): void {
    const file = this.selectedFile();
    const branchId = this.selectedBranch();
    const academicYearId = this.selectedAcademicYear();

    if (!file || !branchId || !academicYearId) return;

    this.uploading.set(true);
    this.importResult.set(null);

    this.studentService.importStudents(file, branchId, academicYearId).subscribe({
      next: (result) => {
        this.importResult.set(result);
        this.uploading.set(false);
      },
      error: (err) => {
        this.uploading.set(false);
        alert(err.message || 'Failed to import students');
      },
    });
  }

  resetImport(): void {
    this.selectedFile.set(null);
    this.importResult.set(null);
  }
}
