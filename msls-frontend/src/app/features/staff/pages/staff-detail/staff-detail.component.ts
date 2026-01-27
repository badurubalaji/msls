/**
 * Staff Detail Page Component
 *
 * Displays detailed information about a staff member.
 */

import {
  Component,
  ChangeDetectionStrategy,
  inject,
  OnInit,
  signal,
  computed,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';

import {
  MslsBadgeComponent,
  MslsAvatarComponent,
} from '../../../../shared/components';
import { StaffService } from '../../services/staff.service';
import {
  Staff,
  StaffStatus,
  getStatusBadgeVariant,
  getStatusLabel,
  getStaffTypeLabel,
  getGenderLabel,
  formatAddress,
  calculateAge,
  calculateTenure,
  StatusHistory,
} from '../../models/staff.model';
import { StaffSalaryComponent } from '../../components/staff-salary/staff-salary.component';

@Component({
  selector: 'msls-staff-detail',
  standalone: true,
  imports: [CommonModule, MslsBadgeComponent, MslsAvatarComponent, StaffSalaryComponent],
  templateUrl: './staff-detail.component.html',
  styleUrl: './staff-detail.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StaffDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private staffService = inject(StaffService);

  // State
  readonly staff = this.staffService.selectedStaff;
  readonly loading = this.staffService.loading;
  readonly error = this.staffService.error;
  readonly statusHistory = this.staffService.statusHistory;

  readonly staffId = signal<string>('');
  readonly activeTab = signal<'overview' | 'employment' | 'salary' | 'history'>('overview');

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.staffId.set(id);
      this.staffService.getStaff(id).subscribe();
      this.staffService.getStatusHistory(id).subscribe();
    }
  }

  // Navigation
  goBack(): void {
    this.router.navigate(['/staff']);
  }

  onEdit(): void {
    this.router.navigate(['/staff', this.staffId(), 'edit']);
  }

  // Tab management
  setTab(tab: 'overview' | 'employment' | 'salary' | 'history'): void {
    this.activeTab.set(tab);
  }

  // Helpers
  getStatusVariant(status: StaffStatus): 'success' | 'warning' | 'danger' | 'neutral' {
    return getStatusBadgeVariant(status);
  }

  getStatusLabel(status: StaffStatus): string {
    return getStatusLabel(status);
  }

  getStaffTypeLabel(type: string): string {
    return getStaffTypeLabel(type as any);
  }

  getGenderLabel(gender: string): string {
    return getGenderLabel(gender as any);
  }

  formatAddress(address: any): string {
    return formatAddress(address);
  }

  getAge(dateOfBirth: string): number {
    return calculateAge(dateOfBirth);
  }

  getTenure(joinDate: string): string {
    return calculateTenure(joinDate);
  }

  formatDate(date: string | undefined): string {
    if (!date) return '—';
    return new Date(date).toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  }

  formatDateTime(date: string): string {
    return new Date(date).toLocaleString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
}
