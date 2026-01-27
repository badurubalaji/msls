/**
 * Check-In/Out Component
 * Displays the check-in/check-out button and today's status
 */

import { Component, Input, Output, EventEmitter, inject, signal, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MslsIconComponent } from '../../../../shared/components/icon/icon.component';
import { AttendanceService } from '../../services/attendance.service';
import { TodayAttendance, HalfDayType, HALF_DAY_TYPE_LABELS } from '../../models/attendance.model';

@Component({
  selector: 'app-check-in-out',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsIconComponent,
  ],
  template: `
    <div class="check-in-wrapper">
      <!-- Main Card with Glass Effect -->
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-primary-600 via-primary-700 to-primary-900"
           style="box-shadow: 0 25px 50px -12px rgba(99, 102, 241, 0.4), 0 0 0 1px rgba(255,255,255,0.1) inset;">

        <!-- Animated Background Elements -->
        <div class="absolute inset-0 overflow-hidden">
          <div class="absolute -top-24 -right-24 w-96 h-96 rounded-full bg-white/10 blur-3xl"></div>
          <div class="absolute -bottom-32 -left-32 w-80 h-80 rounded-full bg-primary-400/20 blur-3xl"></div>
          <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] rounded-full bg-gradient-radial from-white/5 to-transparent"></div>
        </div>

        <!-- Grid Pattern Overlay -->
        <div class="absolute inset-0 opacity-[0.03]"
             style="background-image: url('data:image/svg+xml,%3Csvg width=%2260%22 height=%2260%22 viewBox=%220 0 60 60%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cg fill=%22none%22 fill-rule=%22evenodd%22%3E%3Cg fill=%22%23ffffff%22 fill-opacity=%221%22%3E%3Cpath d=%22M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z%22/%3E%3C/g%3E%3C/g%3E%3C/svg%3E');">
        </div>

        <div class="relative p-6 sm:p-8 lg:p-10">
          <!-- Header Row -->
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
            <div class="flex items-center gap-4">
              <div class="relative">
                <div class="w-16 h-16 rounded-2xl bg-white/20 backdrop-blur-xl flex items-center justify-center"
                     style="box-shadow: 0 8px 32px rgba(0,0,0,0.12), 0 0 0 1px rgba(255,255,255,0.2) inset;">
                  <msls-icon name="clock" size="lg" class="text-white"></msls-icon>
                </div>
                <div class="absolute -bottom-1 -right-1 w-5 h-5 rounded-full bg-emerald-400 border-2 border-primary-700 flex items-center justify-center">
                  <div class="w-2 h-2 rounded-full bg-white animate-pulse"></div>
                </div>
              </div>
              <div>
                <h2 class="text-2xl font-bold text-white tracking-tight">Today's Attendance</h2>
                <p class="text-primary-200/80 text-sm font-medium">{{ currentDate }}</p>
              </div>
            </div>

            <!-- Live Clock -->
            <div class="inline-flex items-center gap-3 px-5 py-3 rounded-2xl bg-white/10 backdrop-blur-xl border border-white/20"
                 style="box-shadow: 0 4px 24px rgba(0,0,0,0.1);">
              <div class="relative flex items-center justify-center w-3 h-3">
                <div class="absolute w-3 h-3 rounded-full bg-emerald-400 animate-ping opacity-75"></div>
                <div class="w-2 h-2 rounded-full bg-emerald-400"></div>
              </div>
              <span class="text-white font-mono text-xl font-semibold tracking-wider">{{ currentTime }}</span>
            </div>
          </div>

          @if (error()) {
            <div class="mb-6 p-4 rounded-2xl bg-red-500/20 backdrop-blur-sm border border-red-400/30"
                 style="box-shadow: 0 4px 24px rgba(239, 68, 68, 0.2);">
              <div class="flex items-center gap-3 text-white">
                <msls-icon name="exclamation-circle" size="sm"></msls-icon>
                <span class="text-sm font-medium">{{ error() }}</span>
              </div>
            </div>
          }

          @if (loading()) {
            <div class="flex items-center justify-center py-16">
              <div class="relative">
                <div class="w-20 h-20 rounded-full border-4 border-white/20"></div>
                <div class="absolute inset-0 w-20 h-20 rounded-full border-4 border-transparent border-t-white animate-spin"></div>
                <div class="absolute inset-2 w-16 h-16 rounded-full border-4 border-transparent border-b-primary-300 animate-spin" style="animation-direction: reverse; animation-duration: 1.5s;"></div>
              </div>
            </div>
          } @else {
            <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
              <!-- Status Card -->
              <div class="lg:col-span-4">
                <div class="h-full p-6 rounded-2xl bg-white/10 backdrop-blur-xl border border-white/20 flex flex-col items-center justify-center text-center"
                     style="box-shadow: 0 8px 32px rgba(0,0,0,0.1), 0 0 0 1px rgba(255,255,255,0.1) inset;">
                  @switch (todayStatus()?.status) {
                    @case ('not_marked') {
                      <div class="relative mb-4">
                        <div class="w-24 h-24 rounded-full bg-white/10 flex items-center justify-center border-2 border-white/20">
                          <msls-icon name="clock" class="text-white/80 text-5xl"></msls-icon>
                        </div>
                      </div>
                      <h3 class="text-2xl font-bold text-white mb-1">Not Checked In</h3>
                      <p class="text-primary-200/70 text-sm">Ready to start your day?</p>
                    }
                    @case ('checked_in') {
                      <div class="relative mb-4">
                        <div class="absolute inset-0 w-24 h-24 rounded-full bg-emerald-400/30 animate-ping"></div>
                        <div class="relative w-24 h-24 rounded-full bg-gradient-to-br from-emerald-400 to-emerald-600 flex items-center justify-center"
                             style="box-shadow: 0 0 40px rgba(52, 211, 153, 0.5);">
                          <msls-icon name="check-circle" class="text-white text-5xl"></msls-icon>
                        </div>
                      </div>
                      <h3 class="text-2xl font-bold text-white mb-1">Checked In</h3>
                      <p class="text-emerald-300 text-sm font-medium">You're on the clock</p>
                    }
                    @case ('checked_out') {
                      <div class="relative mb-4">
                        <div class="w-24 h-24 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center"
                             style="box-shadow: 0 0 40px rgba(59, 130, 246, 0.4);">
                          <msls-icon name="check-circle" class="text-white text-5xl"></msls-icon>
                        </div>
                      </div>
                      <h3 class="text-2xl font-bold text-white mb-1">Day Complete</h3>
                      <p class="text-blue-300 text-sm font-medium">Great work today!</p>
                    }
                  }
                </div>
              </div>

              <!-- Time Details -->
              <div class="lg:col-span-4 grid grid-cols-2 gap-4">
                <!-- Check In Time -->
                <div class="p-5 rounded-2xl bg-white/10 backdrop-blur-xl border border-white/20 flex flex-col items-center justify-center"
                     style="box-shadow: 0 8px 32px rgba(0,0,0,0.1);">
                  <div class="w-12 h-12 mb-3 rounded-xl bg-emerald-500/20 flex items-center justify-center">
                    <msls-icon name="login" size="md" class="text-emerald-300"></msls-icon>
                  </div>
                  <span class="text-xs text-primary-200/60 uppercase tracking-widest font-semibold mb-1">Check In</span>
                  <span class="text-2xl font-bold text-white font-mono">
                    {{ formatTime(todayStatus()?.attendance?.checkInTime) || '--:--' }}
                  </span>
                </div>

                <!-- Check Out Time -->
                <div class="p-5 rounded-2xl bg-white/10 backdrop-blur-xl border border-white/20 flex flex-col items-center justify-center"
                     style="box-shadow: 0 8px 32px rgba(0,0,0,0.1);">
                  <div class="w-12 h-12 mb-3 rounded-xl bg-blue-500/20 flex items-center justify-center">
                    <msls-icon name="logout" size="md" class="text-blue-300"></msls-icon>
                  </div>
                  <span class="text-xs text-primary-200/60 uppercase tracking-widest font-semibold mb-1">Check Out</span>
                  <span class="text-2xl font-bold text-white font-mono">
                    {{ formatTime(todayStatus()?.attendance?.checkOutTime) || '--:--' }}
                  </span>
                </div>

                <!-- Late Indicator -->
                @if (todayStatus()?.attendance?.isLate) {
                  <div class="col-span-2 p-4 rounded-2xl bg-amber-500/20 backdrop-blur-xl border border-amber-400/30 flex items-center justify-center gap-3"
                       style="box-shadow: 0 4px 24px rgba(245, 158, 11, 0.2);">
                    <msls-icon name="exclamation-triangle" size="sm" class="text-amber-300"></msls-icon>
                    <span class="text-amber-200 text-sm font-semibold">
                      Late by {{ todayStatus()?.attendance?.lateMinutes }} minutes
                    </span>
                  </div>
                }

                <!-- Working Hours (if checked in) -->
                @if (todayStatus()?.status === 'checked_in' || todayStatus()?.status === 'checked_out') {
                  <div class="col-span-2 p-4 rounded-2xl bg-white/5 backdrop-blur-xl border border-white/10 flex items-center justify-center gap-3">
                    <msls-icon name="clock" size="sm" class="text-primary-200"></msls-icon>
                    <span class="text-primary-200 text-sm">
                      Working: <span class="font-bold text-white">{{ getWorkingHours() }}</span>
                    </span>
                  </div>
                }
              </div>

              <!-- Action Section -->
              <div class="lg:col-span-4">
                <div class="h-full p-6 rounded-2xl bg-white/10 backdrop-blur-xl border border-white/20 flex flex-col items-center justify-center"
                     style="box-shadow: 0 8px 32px rgba(0,0,0,0.1);">

                  @if (showHalfDayOption && todayStatus()?.canCheckIn) {
                    <div class="w-full mb-5">
                      <label class="block text-xs text-primary-200/60 uppercase tracking-widest font-semibold mb-2 text-center">
                        Working Hours
                      </label>
                      <select
                        [(ngModel)]="selectedHalfDay"
                        class="w-full px-4 py-3 rounded-xl border-0 bg-white/10 text-white placeholder-primary-200/50 focus:ring-2 focus:ring-white/30 text-center font-medium appearance-none cursor-pointer transition-all hover:bg-white/20"
                        style="background-image: url('data:image/svg+xml,%3Csvg xmlns=%22http://www.w3.org/2000/svg%22 fill=%22none%22 viewBox=%220 0 20 20%22%3E%3Cpath stroke=%22%23ffffff%22 stroke-linecap=%22round%22 stroke-linejoin=%22round%22 stroke-width=%222%22 d=%22M6 8l4 4 4-4%22/%3E%3C/svg%3E'); background-position: right 0.75rem center; background-repeat: no-repeat; background-size: 1.25rem;"
                      >
                        <option value="" class="text-gray-900 bg-white">Full Day</option>
                        @for (type of halfDayTypes; track type.value) {
                          <option [value]="type.value" class="text-gray-900 bg-white">{{ type.label }}</option>
                        }
                      </select>
                    </div>
                  }

                  @if (todayStatus()?.canCheckIn) {
                    <button
                      type="button"
                      [disabled]="submitting()"
                      (click)="onCheckIn()"
                      class="group w-full py-4 px-8 rounded-2xl bg-white text-primary-700 font-bold text-lg transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 hover:scale-[1.02] active:scale-[0.98]"
                      style="box-shadow: 0 10px 40px rgba(255,255,255,0.3), 0 0 0 1px rgba(255,255,255,0.5) inset;"
                    >
                      @if (submitting()) {
                        <div class="w-6 h-6 rounded-full border-2 border-primary-300 border-t-primary-700 animate-spin"></div>
                      } @else {
                        <msls-icon name="login" size="sm" class="group-hover:translate-x-1 transition-transform"></msls-icon>
                      }
                      <span>Check In Now</span>
                    </button>
                  }

                  @if (todayStatus()?.canCheckOut) {
                    <button
                      type="button"
                      [disabled]="submitting()"
                      (click)="onCheckOut()"
                      class="group w-full py-4 px-8 rounded-2xl bg-white/20 text-white font-bold text-lg border-2 border-white/30 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 hover:bg-white/30 hover:scale-[1.02] active:scale-[0.98]"
                    >
                      @if (submitting()) {
                        <div class="w-6 h-6 rounded-full border-2 border-white/30 border-t-white animate-spin"></div>
                      } @else {
                        <msls-icon name="logout" size="sm" class="group-hover:translate-x-1 transition-transform"></msls-icon>
                      }
                      <span>Check Out</span>
                    </button>
                  }

                  @if (!todayStatus()?.canCheckIn && !todayStatus()?.canCheckOut) {
                    <div class="text-center py-4">
                      <div class="w-16 h-16 mx-auto mb-3 rounded-full bg-emerald-500/20 flex items-center justify-center">
                        <msls-icon name="check-circle" size="lg" class="text-emerald-300"></msls-icon>
                      </div>
                      <p class="text-primary-200 text-sm font-medium">All done for today!</p>
                      <p class="text-primary-300/50 text-xs mt-1">See you tomorrow</p>
                    </div>
                  }
                </div>
              </div>
            </div>
          }
        </div>
      </div>
    </div>
  `,
  styles: [`
    .check-in-wrapper {
      perspective: 1000px;
    }
  `]
})
export class CheckInOutComponent implements OnInit, OnDestroy {
  @Input({ required: true }) staffId!: string;
  @Input() showHalfDayOption = true;
  @Output() statusChanged = new EventEmitter<TodayAttendance>();

  private attendanceService = inject(AttendanceService);
  private timeInterval: ReturnType<typeof setInterval> | null = null;

  loading = signal(false);
  submitting = signal(false);
  error = signal<string | null>(null);
  todayStatus = signal<TodayAttendance | null>(null);
  currentTime = '';

  selectedHalfDay: HalfDayType | '' = '';

  halfDayTypes = Object.entries(HALF_DAY_TYPE_LABELS).map(([value, label]) => ({
    value: value as HalfDayType,
    label,
  }));

  get currentDate(): string {
    return new Date().toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  }

  ngOnInit(): void {
    this.loadTodayStatus();
    this.updateTime();
    this.timeInterval = setInterval(() => this.updateTime(), 1000);
  }

  ngOnDestroy(): void {
    if (this.timeInterval) {
      clearInterval(this.timeInterval);
    }
  }

  private updateTime(): void {
    this.currentTime = new Date().toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: true,
    });
  }

  loadTodayStatus(): void {
    this.loading.set(true);
    this.error.set(null);

    this.attendanceService.getTodayAttendance(this.staffId).subscribe({
      next: (status) => {
        this.todayStatus.set(status);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set(err.message || 'Failed to load today\'s status');
        this.loading.set(false);
      },
    });
  }

  getWorkingHours(): string {
    const attendance = this.todayStatus()?.attendance;
    if (!attendance?.checkInTime) return '0h 0m';

    const checkIn = new Date(attendance.checkInTime);
    const checkOut = attendance.checkOutTime ? new Date(attendance.checkOutTime) : new Date();

    const diff = checkOut.getTime() - checkIn.getTime();
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

    return `${hours}h ${minutes}m`;
  }

  onCheckIn(): void {
    this.submitting.set(true);
    this.error.set(null);

    const request = {
      staffId: this.staffId,
      halfDayType: this.selectedHalfDay || undefined,
    };

    this.attendanceService.checkIn(request).subscribe({
      next: () => {
        this.submitting.set(false);
        this.loadTodayStatus();
        this.statusChanged.emit(this.todayStatus()!);
      },
      error: (err) => {
        this.error.set(err.message || 'Failed to check in');
        this.submitting.set(false);
      },
    });
  }

  onCheckOut(): void {
    this.submitting.set(true);
    this.error.set(null);

    const request = {
      staffId: this.staffId,
    };

    this.attendanceService.checkOut(request).subscribe({
      next: () => {
        this.submitting.set(false);
        this.loadTodayStatus();
        this.statusChanged.emit(this.todayStatus()!);
      },
      error: (err) => {
        this.error.set(err.message || 'Failed to check out');
        this.submitting.set(false);
      },
    });
  }

  formatTime(isoString?: string): string | null {
    if (!isoString) return null;
    const date = new Date(isoString);
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  }
}
