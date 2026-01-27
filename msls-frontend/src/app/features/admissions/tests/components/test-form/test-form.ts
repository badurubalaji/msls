/**
 * MSLS Test Form Component
 *
 * Form component for creating and editing entrance tests.
 */

import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { EntranceTest, TestSubject, CreateTestDto } from '../../entrance-test.model';
import { AdmissionSessionService } from '../../../sessions/admission-session.service';
import { AdmissionSession } from '../../../sessions/admission-session.model';

/** Session option for dropdown */
interface SessionOption {
  id: string;
  name: string;
  status: string;
}

@Component({
  selector: 'msls-test-form',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './test-form.html',
  styleUrl: './test-form.scss',
})
export class TestFormComponent implements OnInit, OnChanges {
  private readonly sessionService = inject(AdmissionSessionService);

  @Input() test: EntranceTest | null = null;
  @Input() loading = false;
  @Output() save = new EventEmitter<CreateTestDto>();
  @Output() cancel = new EventEmitter<void>();

  /** Whether we are in edit mode */
  isEditMode = signal(false);

  /** Form title based on mode */
  formTitle = computed(() => this.isEditMode() ? 'Edit Entrance Test' : 'Create Entrance Test');

  // Sessions for dropdown
  sessions = signal<SessionOption[]>([]);
  loadingSessions = signal(false);

  // Form data
  formData = signal({
    testName: '',
    sessionId: '', // Will be set from sessions dropdown (UUID)
    testDate: '',
    startTime: '',
    durationMinutes: 60,
    venue: '',
    maxCandidates: 50,
    classNames: [] as string[],
    subjects: [] as TestSubject[],
  });

  // Available classes
  availableClasses = [
    'LKG', 'UKG',
    'Class 1', 'Class 2', 'Class 3', 'Class 4', 'Class 5',
    'Class 6', 'Class 7', 'Class 8', 'Class 9', 'Class 10',
    'Class 11', 'Class 12',
  ];

  // New subject input
  newSubjectName = signal('');
  newSubjectMaxMarks = signal(25);

  // Validation
  errors = signal<Record<string, string>>({});

  ngOnInit(): void {
    this.loadSessions();
    this.initializeForm();
  }

  ngOnChanges(changes: SimpleChanges): void {
    // Re-initialize form when test input changes (switching between create/edit)
    if (changes['test'] && !changes['test'].firstChange) {
      this.initializeForm();
    }
  }

  /** Initialize or reset the form based on test input */
  private initializeForm(): void {
    if (this.test) {
      // Edit mode - populate form with existing test data
      this.isEditMode.set(true);
      this.formData.set({
        testName: this.test.testName,
        sessionId: this.test.sessionId,
        testDate: this.test.testDate,
        startTime: this.test.startTime,
        durationMinutes: this.test.durationMinutes,
        venue: this.test.venue || '',
        maxCandidates: this.test.maxCandidates,
        classNames: [...this.test.classNames],
        subjects: this.test.subjects.map(s => ({ ...s })),
      });
    } else {
      // Create mode - reset form to defaults
      this.isEditMode.set(false);
      this.formData.set({
        testName: '',
        sessionId: '',
        testDate: '',
        startTime: '',
        durationMinutes: 60,
        venue: '',
        maxCandidates: 50,
        classNames: [],
        subjects: [],
      });
      // Auto-select session will happen in loadSessions callback
    }
    // Clear any validation errors
    this.errors.set({});
  }

  private loadSessions(): void {
    this.loadingSessions.set(true);
    this.sessionService.getSessions().subscribe({
      next: (sessions) => {
        const sessionOptions: SessionOption[] = sessions.map(s => ({
          id: s.id,
          name: s.name,
          status: s.status,
        }));
        this.sessions.set(sessionOptions);
        this.loadingSessions.set(false);

        // Auto-select first open session if creating new test
        if (!this.test && sessionOptions.length > 0) {
          const openSession = sessionOptions.find(s => s.status === 'open');
          const defaultSession = openSession || sessionOptions[0];
          this.updateField('sessionId', defaultSession.id);
        }
      },
      error: (err) => {
        console.error('Failed to load sessions:', err);
        this.loadingSessions.set(false);
      },
    });
  }

  updateField(field: string, value: any): void {
    this.formData.update(data => ({ ...data, [field]: value }));
    // Clear error for field
    this.errors.update(errs => {
      const { [field]: _, ...rest } = errs;
      return rest;
    });
  }

  toggleClass(className: string): void {
    const current = this.formData().classNames;
    if (current.includes(className)) {
      this.updateField('classNames', current.filter(c => c !== className));
    } else {
      this.updateField('classNames', [...current, className]);
    }
  }

  isClassSelected(className: string): boolean {
    return this.formData().classNames.includes(className);
  }

  addSubject(): void {
    const name = this.newSubjectName().trim();
    const maxMarks = this.newSubjectMaxMarks();

    if (!name) return;
    if (maxMarks <= 0) return;

    // Check for duplicate
    if (this.formData().subjects.some(s => s.name.toLowerCase() === name.toLowerCase())) {
      return;
    }

    const subjects = [...this.formData().subjects, { name, maxMarks }];
    this.updateField('subjects', subjects);
    this.newSubjectName.set('');
    this.newSubjectMaxMarks.set(25);
  }

  removeSubject(index: number): void {
    const subjects = this.formData().subjects.filter((_, i) => i !== index);
    this.updateField('subjects', subjects);
  }

  getTotalMaxMarks(): number {
    return this.formData().subjects.reduce((sum, s) => sum + s.maxMarks, 0);
  }

  validate(): boolean {
    const data = this.formData();
    const newErrors: Record<string, string> = {};

    if (!data.sessionId) {
      newErrors['sessionId'] = 'Please select an admission session';
    }
    if (!data.testName.trim()) {
      newErrors['testName'] = 'Test name is required';
    }
    if (!data.testDate) {
      newErrors['testDate'] = 'Test date is required';
    }
    if (!data.startTime) {
      newErrors['startTime'] = 'Start time is required';
    }
    if (data.durationMinutes <= 0) {
      newErrors['durationMinutes'] = 'Duration must be greater than 0';
    }
    if (data.maxCandidates <= 0) {
      newErrors['maxCandidates'] = 'Max candidates must be greater than 0';
    }
    if (data.classNames.length === 0) {
      newErrors['classNames'] = 'Select at least one class';
    }
    if (data.subjects.length === 0) {
      newErrors['subjects'] = 'Add at least one subject';
    }

    this.errors.set(newErrors);
    return Object.keys(newErrors).length === 0;
  }

  onSubmit(): void {
    if (!this.validate()) return;

    const data = this.formData();
    this.save.emit({
      testName: data.testName,
      sessionId: data.sessionId,
      testDate: data.testDate,
      startTime: data.startTime,
      durationMinutes: data.durationMinutes,
      venue: data.venue || undefined,
      maxCandidates: data.maxCandidates,
      classNames: data.classNames,
      subjects: data.subjects,
    });
  }

  onCancel(): void {
    this.cancel.emit();
  }
}
