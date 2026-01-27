/**
 * MSLS Entrance Tests Page Component
 *
 * Main component for managing entrance tests - displays a list of tests with CRUD operations.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

import { MslsModalComponent } from '../../../../../shared/components';
import { ToastService } from '../../../../../shared/services';
import {
  EntranceTest,
  TestStatus,
  TEST_STATUS_CONFIG,
  formatTime,
  formatDate,
} from '../../entrance-test.model';
import { EntranceTestService } from '../../entrance-test.service';
import { TestFormComponent } from '../../components/test-form/test-form';

@Component({
  selector: 'msls-tests',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent, TestFormComponent],
  templateUrl: './tests.html',
  styleUrl: './tests.scss',
})
export class TestsComponent implements OnInit {
  private readonly testService = inject(EntranceTestService);
  private readonly toastService = inject(ToastService);
  private readonly router = inject(Router);

  // State signals
  tests = signal<EntranceTest[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<TestStatus | ''>('');

  // Modal state
  showTestModal = signal(false);
  showDeleteModal = signal(false);
  editingTest = signal<EntranceTest | null>(null);
  testToDelete = signal<EntranceTest | null>(null);

  // Computed filtered and sorted tests
  filteredTests = computed(() => {
    let result = this.tests();
    const term = this.searchTerm().toLowerCase();
    const status = this.statusFilter();

    if (term) {
      result = result.filter(
        test =>
          test.testName.toLowerCase().includes(term) ||
          test.venue?.toLowerCase().includes(term) ||
          test.classNames.some(c => c.toLowerCase().includes(term))
      );
    }

    if (status) {
      result = result.filter(test => test.status === status);
    }

    // Sort by date (descending) and then by start time (ascending)
    // Upcoming tests first, then past tests
    return [...result].sort((a, b) => {
      // Parse dates for comparison
      const dateA = new Date(`${a.testDate}T${a.startTime}`);
      const dateB = new Date(`${b.testDate}T${b.startTime}`);
      const now = new Date();

      // Separate upcoming and past tests
      const aIsUpcoming = dateA >= now;
      const bIsUpcoming = dateB >= now;

      // Upcoming tests come first
      if (aIsUpcoming && !bIsUpcoming) return -1;
      if (!aIsUpcoming && bIsUpcoming) return 1;

      // For upcoming tests: sort ascending (nearest first)
      // For past tests: sort descending (most recent first)
      if (aIsUpcoming) {
        return dateA.getTime() - dateB.getTime();
      } else {
        return dateB.getTime() - dateA.getTime();
      }
    });
  });

  // Helper functions for template
  getStatusConfig(status: TestStatus) {
    return TEST_STATUS_CONFIG[status];
  }

  formatTime = formatTime;
  formatDate = formatDate;

  ngOnInit(): void {
    this.loadTests();
  }

  loadTests(): void {
    this.loading.set(true);
    this.error.set(null);

    this.testService.getTests().subscribe({
      next: tests => {
        this.tests.set(tests);
        this.loading.set(false);
      },
      error: err => {
        this.error.set('Failed to load entrance tests. Please try again.');
        this.loading.set(false);
        console.error('Failed to load tests:', err);
      },
    });
  }

  onSearchChange(term: string): void {
    this.searchTerm.set(term);
  }

  onStatusFilterChange(status: string): void {
    this.statusFilter.set(status as TestStatus | '');
  }

  // Create/Edit Modal
  openCreateModal(): void {
    this.editingTest.set(null);
    this.showTestModal.set(true);
  }

  editTest(test: EntranceTest): void {
    this.editingTest.set(test);
    this.showTestModal.set(true);
  }

  closeTestModal(): void {
    this.showTestModal.set(false);
    this.editingTest.set(null);
  }

  saveTest(data: any): void {
    this.saving.set(true);

    const editing = this.editingTest();
    const operation = editing
      ? this.testService.updateTest(editing.id, data)
      : this.testService.createTest(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Test updated successfully' : 'Test created successfully'
        );
        this.closeTestModal();
        this.loadTests();
        this.saving.set(false);
      },
      error: err => {
        this.toastService.error(
          editing ? 'Failed to update test' : 'Failed to create test'
        );
        this.saving.set(false);
        console.error('Failed to save test:', err);
      },
    });
  }

  // Delete Modal
  confirmDelete(test: EntranceTest): void {
    this.testToDelete.set(test);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.testToDelete.set(null);
  }

  deleteTest(): void {
    const test = this.testToDelete();
    if (!test) return;

    this.deleting.set(true);

    this.testService.deleteTest(test.id).subscribe({
      next: () => {
        this.toastService.success('Test deleted successfully');
        this.closeDeleteModal();
        this.loadTests();
        this.deleting.set(false);
      },
      error: err => {
        this.toastService.error('Failed to delete test');
        this.deleting.set(false);
        console.error('Failed to delete test:', err);
      },
    });
  }

  // View registrations/results
  viewRegistrations(test: EntranceTest): void {
    this.router.navigate(['/admissions/tests', test.id, 'registrations']);
  }

  viewResults(test: EntranceTest): void {
    this.router.navigate(['/admissions/tests', test.id, 'results']);
  }

  // Download hall tickets
  downloadHallTickets(test: EntranceTest): void {
    this.testService.generateHallTickets(test.id).subscribe({
      next: blob => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `hall-tickets-${test.testName.replace(/\s+/g, '-')}.pdf`;
        link.click();
        window.URL.revokeObjectURL(url);
        this.toastService.success('Hall tickets downloaded');
      },
      error: err => {
        this.toastService.error('Failed to download hall tickets');
        console.error('Failed to download hall tickets:', err);
      },
    });
  }

  // Status actions
  completeTest(test: EntranceTest): void {
    this.testService.completeTest(test.id).subscribe({
      next: () => {
        this.toastService.success('Test marked as completed');
        this.loadTests();
      },
      error: err => {
        this.toastService.error('Failed to update test status');
        console.error('Failed to complete test:', err);
      },
    });
  }

  cancelTest(test: EntranceTest): void {
    this.testService.cancelTest(test.id).subscribe({
      next: () => {
        this.toastService.success('Test cancelled');
        this.loadTests();
      },
      error: err => {
        this.toastService.error('Failed to cancel test');
        console.error('Failed to cancel test:', err);
      },
    });
  }

  getTotalMaxMarks(test: EntranceTest): number {
    return test.subjects.reduce((sum, s) => sum + s.maxMarks, 0);
  }
}
