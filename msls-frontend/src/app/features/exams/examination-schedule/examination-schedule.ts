import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ExamService } from '../exam.service';
import {
  Examination,
  ExamSchedule,
  CreateScheduleRequest,
  UpdateScheduleRequest,
  EXAM_STATUSES,
} from '../exam.model';
import { ToastService } from '../../../shared/services/toast.service';
import { ApiService } from '../../../core/services/api.service';

interface SubjectOption {
  id: string;
  name: string;
  code: string;
}

@Component({
  selector: 'app-examination-schedule',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './examination-schedule.html',
  styleUrl: './examination-schedule.scss',
})
export class ExaminationSchedule implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private examService = inject(ExamService);
  private apiService = inject(ApiService);
  private toastService = inject(ToastService);

  // Data signals
  examination = signal<Examination | null>(null);
  schedules = signal<ExamSchedule[]>([]);
  subjects = signal<SubjectOption[]>([]);

  // State signals
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);

  // Modal state
  showFormModal = signal(false);
  showDeleteModal = signal(false);
  editingSchedule = signal<ExamSchedule | null>(null);
  scheduleToDelete = signal<ExamSchedule | null>(null);

  // Form data
  formData = this.getEmptyFormData();

  // Computed values
  modalTitle = computed(() => {
    const editing = this.editingSchedule();
    return editing ? `Edit Schedule: ${editing.subjectName}` : 'Add Schedule';
  });

  examName = computed(() => this.examination()?.name || 'Examination');

  isDraft = computed(() => this.examination()?.status === 'draft');

  canEdit = computed(() => this.isDraft());

  sortedSchedules = computed(() => {
    return [...this.schedules()].sort((a, b) => {
      const dateCompare = a.examDate.localeCompare(b.examDate);
      if (dateCompare !== 0) return dateCompare;
      return a.startTime.localeCompare(b.startTime);
    });
  });

  // Get available subjects (not already scheduled)
  availableSubjects = computed(() => {
    const scheduled = new Set(this.schedules().map(s => s.subjectId));
    const editing = this.editingSchedule();
    return this.subjects().filter(s => {
      // Include if not scheduled, or if it's the subject being edited
      return !scheduled.has(s.id) || (editing && editing.subjectId === s.id);
    });
  });

  ngOnInit(): void {
    const examId = this.route.snapshot.paramMap.get('id');
    if (examId) {
      this.loadExamination(examId);
      this.loadSubjects();
    } else {
      this.error.set('Invalid examination ID');
      this.loading.set(false);
    }
  }

  loadExamination(id: string): void {
    this.loading.set(true);
    this.error.set(null);

    this.examService.getExamination(id).subscribe({
      next: (exam) => {
        this.examination.set(exam);
        this.schedules.set(exam.schedules || []);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load examination');
        this.loading.set(false);
      },
    });
  }

  loadSubjects(): void {
    this.apiService.get<SubjectOption[]>('/subjects').subscribe({
      next: (subjects) => this.subjects.set(subjects),
      error: () => console.error('Failed to load subjects'),
    });
  }

  goBack(): void {
    this.router.navigate(['/exams/list']);
  }

  openAddModal(): void {
    if (!this.canEdit()) {
      this.toastService.error('Cannot modify schedules for published examinations');
      return;
    }

    this.editingSchedule.set(null);
    this.formData = this.getEmptyFormData();

    // Pre-fill with exam date range defaults
    const exam = this.examination();
    if (exam) {
      this.formData.examDate = exam.startDate;
      this.formData.maxMarks = 100;
    }

    this.showFormModal.set(true);
  }

  editSchedule(schedule: ExamSchedule): void {
    if (!this.canEdit()) {
      this.toastService.error('Cannot modify schedules for published examinations');
      return;
    }

    this.editingSchedule.set(schedule);
    this.formData = {
      subjectId: schedule.subjectId,
      examDate: schedule.examDate,
      startTime: schedule.startTime,
      endTime: schedule.endTime,
      maxMarks: schedule.maxMarks,
      passingMarks: schedule.passingMarks ?? null,
      venue: schedule.venue || '',
      notes: schedule.notes || '',
    };
    this.showFormModal.set(true);
  }

  closeFormModal(): void {
    this.showFormModal.set(false);
    this.editingSchedule.set(null);
  }

  isFormValid(): boolean {
    return !!(
      this.formData.subjectId &&
      this.formData.examDate &&
      this.formData.startTime &&
      this.formData.endTime &&
      this.formData.maxMarks > 0 &&
      this.formData.endTime > this.formData.startTime
    );
  }

  isDateValid(): boolean {
    const exam = this.examination();
    if (!exam || !this.formData.examDate) return true;
    return this.formData.examDate >= exam.startDate && this.formData.examDate <= exam.endDate;
  }

  saveSchedule(): void {
    if (!this.isFormValid()) {
      this.toastService.error('Please fill in all required fields correctly');
      return;
    }

    if (!this.isDateValid()) {
      this.toastService.error('Schedule date must be within examination period');
      return;
    }

    const exam = this.examination();
    if (!exam) return;

    this.saving.set(true);
    const editing = this.editingSchedule();

    const payload: CreateScheduleRequest | UpdateScheduleRequest = {
      subjectId: this.formData.subjectId,
      examDate: this.formData.examDate,
      startTime: this.formData.startTime,
      endTime: this.formData.endTime,
      maxMarks: this.formData.maxMarks,
      passingMarks: this.formData.passingMarks ?? undefined,
      venue: this.formData.venue || undefined,
      notes: this.formData.notes || undefined,
    };

    const operation = editing
      ? this.examService.updateSchedule(exam.id, editing.id, payload)
      : this.examService.createSchedule(exam.id, payload as CreateScheduleRequest);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Schedule updated successfully' : 'Schedule added successfully'
        );
        this.closeFormModal();
        this.loadExamination(exam.id);
        this.saving.set(false);
      },
      error: (err) => {
        const message = err?.error?.detail || (editing ? 'Failed to update schedule' : 'Failed to add schedule');
        this.toastService.error(message);
        this.saving.set(false);
      },
    });
  }

  confirmDelete(schedule: ExamSchedule): void {
    if (!this.canEdit()) {
      this.toastService.error('Cannot modify schedules for published examinations');
      return;
    }

    this.scheduleToDelete.set(schedule);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.scheduleToDelete.set(null);
  }

  deleteSchedule(): void {
    const schedule = this.scheduleToDelete();
    const exam = this.examination();
    if (!schedule || !exam) return;

    this.deleting.set(true);

    this.examService.deleteSchedule(exam.id, schedule.id).subscribe({
      next: () => {
        this.toastService.success('Schedule deleted successfully');
        this.closeDeleteModal();
        this.loadExamination(exam.id);
        this.deleting.set(false);
      },
      error: (err) => {
        const message = err?.error?.detail || 'Failed to delete schedule';
        this.toastService.error(message);
        this.deleting.set(false);
      },
    });
  }

  getStatusLabel(status: string): string {
    return EXAM_STATUSES.find(s => s.value === status)?.label || status;
  }

  getStatusClass(status: string): string {
    return `status-${status}`;
  }

  formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-IN', {
      weekday: 'short',
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  formatTime(timeStr: string): string {
    const [hours, minutes] = timeStr.split(':').map(Number);
    const period = hours >= 12 ? 'PM' : 'AM';
    const displayHours = hours % 12 || 12;
    return `${displayHours}:${minutes.toString().padStart(2, '0')} ${period}`;
  }

  private getEmptyFormData() {
    return {
      subjectId: '',
      examDate: '',
      startTime: '09:00',
      endTime: '12:00',
      maxMarks: 100,
      passingMarks: null as number | null,
      venue: '',
      notes: '',
    };
  }
}
