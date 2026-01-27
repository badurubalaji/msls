import { Component, ChangeDetectionStrategy, inject, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../core/services';
import { StudentService } from '../students/services/student.service';
import { Student } from '../students/models/student.model';

/**
 * DashboardComponent - Main dashboard with Mintlify-inspired clean design.
 * Uses Font Awesome icons and proper CSS styling.
 */
@Component({
  selector: 'msls-dashboard',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="dashboard">
      <!-- Hero Section -->
      <div class="hero">
        <div class="hero-content">
          <div class="hero-left">
            <div class="status-badge">
              <span class="status-dot"></span>
              System Online
            </div>
            <h1 class="hero-title">Welcome back, {{ userName }}</h1>
            <p class="hero-subtitle">Here's what's happening with your school today.</p>
          </div>
          <div class="hero-actions">
            <button class="btn btn-secondary">
              <i class="fa-regular fa-calendar"></i>
              View Schedule
            </button>
            <button class="btn btn-primary">
              <i class="fa-solid fa-plus"></i>
              Add Student
            </button>
          </div>
        </div>
      </div>

      <div class="dashboard-content">
        <!-- Stats Grid -->
        <div class="stats-grid">
          @for (stat of stats(); track stat.label) {
            <div class="stat-card">
              <div class="stat-header">
                <div class="stat-info">
                  <p class="stat-label">{{ stat.label }}</p>
                  <p class="stat-value">{{ stat.value }}</p>
                </div>
                <div class="stat-icon" [style.background]="stat.bgColor">
                  <i [class]="stat.icon" [style.color]="stat.iconColor"></i>
                </div>
              </div>
              <div class="stat-trend">
                <span class="trend-value" [class.trend-up]="stat.trendUp" [class.trend-down]="!stat.trendUp">
                  <i [class]="stat.trendUp ? 'fa-solid fa-arrow-up' : 'fa-solid fa-arrow-down'"></i>
                  {{ stat.trend }}
                </span>
                <span class="trend-label">vs last month</span>
              </div>
            </div>
          }
        </div>

        <!-- Main Grid -->
        <div class="main-grid">
          <!-- Quick Actions -->
          <div class="card">
            <div class="card-header">
              <div class="card-icon" style="background: #fef3c7;">
                <i class="fa-solid fa-bolt" style="color: #d97706;"></i>
              </div>
              <h2 class="card-title">Quick Actions</h2>
            </div>
            <div class="actions-list">
              @for (action of quickActions; track action.label) {
                <button class="action-item">
                  <div class="action-icon">
                    <i [class]="action.icon"></i>
                  </div>
                  <span class="action-label">{{ action.label }}</span>
                  <i class="fa-solid fa-chevron-right action-arrow"></i>
                </button>
              }
            </div>
          </div>

          <!-- Recent Activity -->
          <div class="card card-wide">
            <div class="card-header">
              <div class="card-icon" style="background: #eef2ff;">
                <i class="fa-regular fa-clock" style="color: #4f46e5;"></i>
              </div>
              <h2 class="card-title">Recent Activity</h2>
              <button class="view-all-btn">
                View all
                <i class="fa-solid fa-arrow-right"></i>
              </button>
            </div>

            <!-- Activity Table -->
            <div class="table-container">
              <table class="activity-table">
                <thead>
                  <tr>
                    <th>Activity</th>
                    <th>Type</th>
                    <th>User</th>
                    <th style="text-align: right;">Time</th>
                  </tr>
                </thead>
                <tbody>
                  @for (activity of recentActivity; track activity.activity) {
                    <tr>
                      <td>{{ activity.activity }}</td>
                      <td>
                        <span class="type-badge" [style.background]="getTypeBadgeBg(activity.type)" [style.color]="getTypeBadgeColor(activity.type)">
                          {{ activity.type }}
                        </span>
                      </td>
                      <td class="user-cell">{{ activity.user }}</td>
                      <td class="time-cell">{{ activity.time }}</td>
                    </tr>
                  }
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- Bottom Grid -->
        <div class="bottom-grid">
          <!-- Upcoming Events -->
          <div class="card">
            <div class="card-header">
              <div class="card-icon" style="background: #dcfce7;">
                <i class="fa-regular fa-calendar-days" style="color: #16a34a;"></i>
              </div>
              <h2 class="card-title">Upcoming Events</h2>
            </div>
            <div class="events-list">
              @for (event of upcomingEvents; track event.title) {
                <div class="event-item">
                  <div class="event-date" [style.background]="event.bgGradient">
                    <span class="event-day">{{ event.day }}</span>
                    <span class="event-month">{{ event.month }}</span>
                  </div>
                  <div class="event-info">
                    <h4 class="event-title">{{ event.title }}</h4>
                    <p class="event-description">{{ event.description }}</p>
                  </div>
                  <i class="fa-solid fa-chevron-right event-arrow"></i>
                </div>
              }
            </div>
          </div>

          <!-- Announcements -->
          <div class="card">
            <div class="card-header">
              <div class="card-icon" style="background: #fee2e2;">
                <i class="fa-solid fa-bullhorn" style="color: #dc2626;"></i>
              </div>
              <h2 class="card-title">Announcements</h2>
            </div>
            <div class="announcements-list">
              @for (announcement of announcements; track announcement.title) {
                <div class="announcement-item">
                  <div class="announcement-header">
                    <h4 class="announcement-title">{{ announcement.title }}</h4>
                    <span class="priority-badge" [class.priority-high]="announcement.priority === 'high'">
                      {{ announcement.priority }}
                    </span>
                  </div>
                  <p class="announcement-message">{{ announcement.message }}</p>
                  <span class="announcement-time">{{ announcement.time }}</span>
                </div>
              }
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .dashboard {
      min-height: 100vh;
      background: #f8fafc;
    }

    /* Hero Section */
    .hero {
      position: relative;
      overflow: hidden;
      background: white;
      border-bottom: 1px solid #e2e8f0;
    }

    .hero::before {
      content: '';
      position: absolute;
      inset: 0;
      background: linear-gradient(135deg, #eef2ff 0%, #ffffff 50%, #ecfdf5 100%);
    }

    .hero-content {
      position: relative;
      max-width: 1280px;
      margin: 0 auto;
      padding: 3rem 1.5rem;
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    @media (min-width: 768px) {
      .hero-content {
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
      }
    }

    .status-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.25rem 0.75rem;
      background: #eef2ff;
      border: 1px solid #c7d2fe;
      border-radius: 9999px;
      color: #4f46e5;
      font-size: 0.875rem;
      font-weight: 500;
      margin-bottom: 1rem;
    }

    .status-dot {
      width: 0.5rem;
      height: 0.5rem;
      background: #10b981;
      border-radius: 50%;
      animation: pulse 2s infinite;
    }

    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.5; }
    }

    .hero-title {
      font-size: 1.875rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
      letter-spacing: -0.025em;
    }

    @media (min-width: 768px) {
      .hero-title {
        font-size: 2.25rem;
      }
    }

    .hero-subtitle {
      font-size: 1.125rem;
      color: #64748b;
      margin: 0;
    }

    .hero-actions {
      display: flex;
      gap: 0.75rem;
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1rem;
      font-size: 0.875rem;
      font-weight: 500;
      border-radius: 9999px;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
      box-shadow: 0 2px 8px rgba(79, 70, 229, 0.25);
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    .btn-secondary {
      background: white;
      color: #334155;
      border: 1px solid #e2e8f0;
    }

    .btn-secondary:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
    }

    /* Dashboard Content */
    .dashboard-content {
      max-width: 1280px;
      margin: 0 auto;
      padding: 2rem 1.5rem;
    }

    /* Stats Grid */
    .stats-grid {
      display: grid;
      grid-template-columns: repeat(1, 1fr);
      gap: 1rem;
      margin-bottom: 2rem;
    }

    @media (min-width: 640px) {
      .stats-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }

    @media (min-width: 1024px) {
      .stats-grid {
        grid-template-columns: repeat(4, 1fr);
      }
    }

    .stat-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
      transition: all 0.3s;
    }

    .stat-card:hover {
      border-color: #cbd5e1;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
    }

    .stat-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
    }

    .stat-label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #64748b;
      margin: 0;
    }

    .stat-value {
      font-size: 1.875rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0.5rem 0 0 0;
    }

    .stat-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .stat-trend {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-top: 1rem;
    }

    .trend-value {
      display: flex;
      align-items: center;
      gap: 0.25rem;
      font-size: 0.875rem;
      font-weight: 500;
    }

    .trend-value i {
      font-size: 0.75rem;
    }

    .trend-up {
      color: #10b981;
    }

    .trend-down {
      color: #f43f5e;
    }

    .trend-label {
      font-size: 0.875rem;
      color: #94a3b8;
    }

    /* Main Grid */
    .main-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 1.5rem;
      margin-bottom: 2rem;
    }

    @media (min-width: 1024px) {
      .main-grid {
        grid-template-columns: 1fr 2fr;
      }
    }

    /* Card */
    .card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    .card-wide {
      grid-column: span 1;
    }

    @media (min-width: 1024px) {
      .card-wide {
        grid-column: span 1;
      }
    }

    .card-header {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      margin-bottom: 1.5rem;
    }

    .card-icon {
      width: 2rem;
      height: 2rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .card-title {
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
      flex: 1;
    }

    .view-all-btn {
      display: flex;
      align-items: center;
      gap: 0.25rem;
      font-size: 0.875rem;
      font-weight: 500;
      color: #4f46e5;
      background: none;
      border: none;
      cursor: pointer;
    }

    .view-all-btn:hover {
      color: #4338ca;
    }

    /* Actions List */
    .actions-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .action-item {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.75rem 1rem;
      background: none;
      border: 1px solid transparent;
      border-radius: 0.75rem;
      cursor: pointer;
      transition: all 0.2s;
      text-align: left;
      width: 100%;
    }

    .action-item:hover {
      background: #f8fafc;
      border-color: #e2e8f0;
    }

    .action-item:hover .action-icon {
      background: #eef2ff;
    }

    .action-item:hover .action-icon i {
      color: #4f46e5;
    }

    .action-icon {
      width: 2.25rem;
      height: 2.25rem;
      background: #f1f5f9;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.2s;
    }

    .action-icon i {
      color: #64748b;
      transition: color 0.2s;
    }

    .action-label {
      flex: 1;
      font-size: 0.875rem;
      font-weight: 500;
      color: #334155;
    }

    .action-arrow {
      font-size: 0.75rem;
      color: #cbd5e1;
    }

    .action-item:hover .action-arrow {
      color: #4f46e5;
    }

    /* Table */
    .table-container {
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      overflow: hidden;
    }

    .activity-table {
      width: 100%;
      border-collapse: collapse;
    }

    .activity-table thead {
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .activity-table th {
      padding: 0.75rem 1rem;
      text-align: left;
      font-size: 0.75rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .activity-table td {
      padding: 0.875rem 1rem;
      font-size: 0.875rem;
      border-bottom: 1px solid #f1f5f9;
    }

    .activity-table tbody tr:hover {
      background: #f8fafc;
    }

    .activity-table tbody tr:last-child td {
      border-bottom: none;
    }

    .type-badge {
      display: inline-flex;
      padding: 0.25rem 0.625rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .user-cell {
      color: #475569;
    }

    .time-cell {
      color: #94a3b8;
      text-align: right;
    }

    /* Bottom Grid */
    .bottom-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 1.5rem;
    }

    @media (min-width: 768px) {
      .bottom-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }

    /* Events List */
    .events-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }

    .event-item {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 0.75rem;
      border-radius: 0.75rem;
      cursor: pointer;
      transition: background 0.2s;
    }

    .event-item:hover {
      background: #f8fafc;
    }

    .event-date {
      width: 3.5rem;
      height: 3.5rem;
      border-radius: 0.75rem;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      color: white;
    }

    .event-day {
      font-size: 1.125rem;
      font-weight: 700;
      line-height: 1;
    }

    .event-month {
      font-size: 0.625rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      opacity: 0.9;
    }

    .event-info {
      flex: 1;
      min-width: 0;
    }

    .event-title {
      font-size: 0.875rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.25rem 0;
    }

    .event-item:hover .event-title {
      color: #4f46e5;
    }

    .event-description {
      font-size: 0.875rem;
      color: #64748b;
      margin: 0;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .event-arrow {
      color: #cbd5e1;
    }

    .event-item:hover .event-arrow {
      color: #4f46e5;
    }

    /* Announcements List */
    .announcements-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .announcement-item {
      padding: 1rem;
      border: 1px solid #f1f5f9;
      border-radius: 0.75rem;
      transition: all 0.2s;
    }

    .announcement-item:hover {
      border-color: #e2e8f0;
      background: rgba(248, 250, 252, 0.5);
    }

    .announcement-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 0.75rem;
      margin-bottom: 0.5rem;
    }

    .announcement-title {
      font-size: 0.875rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
    }

    .priority-badge {
      padding: 0.125rem 0.5rem;
      font-size: 0.625rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      border-radius: 9999px;
      background: #f1f5f9;
      color: #64748b;
    }

    .priority-badge.priority-high {
      background: #fee2e2;
      color: #dc2626;
    }

    .announcement-message {
      font-size: 0.875rem;
      color: #64748b;
      line-height: 1.5;
      margin: 0 0 0.75rem 0;
    }

    .announcement-time {
      font-size: 0.75rem;
      color: #94a3b8;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DashboardComponent implements OnInit {
  private authService = inject(AuthService);
  private studentService = inject(StudentService);

  // Student statistics signals
  private allStudents = signal<Student[]>([]);
  private totalStudentCount = signal<number>(0);
  private statsLoading = signal<boolean>(false);

  get userName(): string {
    const user = this.authService.currentUser();
    return user?.firstName || user?.email?.split('@')[0] || 'User';
  }

  // Computed stats from real data
  readonly activeStudents = computed(() =>
    this.allStudents().filter(s => s.status === 'active').length
  );

  readonly graduatedStudents = computed(() =>
    this.allStudents().filter(s => s.status === 'graduated').length
  );

  readonly newThisMonth = computed(() => {
    const now = new Date();
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
    return this.allStudents().filter(s => new Date(s.createdAt) >= startOfMonth).length;
  });

  // Dynamic stats computed from student data
  readonly stats = computed(() => [
    {
      label: 'Total Students',
      value: this.formatNumber(this.totalStudentCount()),
      icon: 'fa-solid fa-users',
      bgColor: '#eef2ff',
      iconColor: '#4f46e5',
      trend: '—',
      trendUp: true,
    },
    {
      label: 'Active Students',
      value: this.formatNumber(this.activeStudents()),
      icon: 'fa-solid fa-user-check',
      bgColor: '#dcfce7',
      iconColor: '#16a34a',
      trend: '—',
      trendUp: true,
    },
    {
      label: 'New This Month',
      value: this.formatNumber(this.newThisMonth()),
      icon: 'fa-solid fa-user-plus',
      bgColor: '#fef3c7',
      iconColor: '#d97706',
      trend: '—',
      trendUp: true,
    },
    {
      label: 'Graduated',
      value: this.formatNumber(this.graduatedStudents()),
      icon: 'fa-solid fa-graduation-cap',
      bgColor: '#e0f2fe',
      iconColor: '#0284c7',
      trend: '—',
      trendUp: true,
    },
  ]);

  ngOnInit(): void {
    this.loadStudentStats();
  }

  private loadStudentStats(): void {
    this.statsLoading.set(true);
    // Load all students to calculate stats (limit high enough to get all)
    this.studentService.loadStudents({ limit: 1000 }).subscribe({
      next: (response) => {
        this.allStudents.set(response.students);
        this.totalStudentCount.set(response.total);
        this.statsLoading.set(false);
      },
      error: () => {
        this.statsLoading.set(false);
      },
    });
  }

  private formatNumber(num: number): string {
    if (num === 0) return '0';
    return num.toLocaleString('en-IN');
  }

  quickActions = [
    { label: 'Add New Student', icon: 'fa-solid fa-user-plus' },
    { label: 'Record Attendance', icon: 'fa-solid fa-clipboard-check' },
    { label: 'View Reports', icon: 'fa-solid fa-chart-bar' },
    { label: 'Send Notification', icon: 'fa-solid fa-bell' },
    { label: 'Manage Fees', icon: 'fa-solid fa-wallet' },
  ];

  recentActivity = [
    { activity: 'New student enrolled', type: 'Enrollment', user: 'Admin', time: '2 min ago' },
    { activity: 'Attendance marked for Class 10A', type: 'Attendance', user: 'Mr. Smith', time: '15 min ago' },
    { activity: 'Fee payment received - $1,200', type: 'Finance', user: 'System', time: '1 hour ago' },
    { activity: 'Exam schedule published', type: 'Academic', user: 'Principal', time: '2 hours ago' },
    { activity: 'Parent meeting scheduled', type: 'Event', user: 'Mrs. Johnson', time: '3 hours ago' },
  ];

  upcomingEvents = [
    { day: '15', month: 'Feb', title: 'Parent-Teacher Meeting', description: 'Annual PTM for all grades', bgGradient: 'linear-gradient(135deg, #6366f1, #4f46e5)' },
    { day: '20', month: 'Feb', title: 'Science Exhibition', description: 'Student projects showcase', bgGradient: 'linear-gradient(135deg, #10b981, #059669)' },
    { day: '28', month: 'Feb', title: 'Term Exams Begin', description: 'Final term examinations', bgGradient: 'linear-gradient(135deg, #f59e0b, #d97706)' },
  ];

  announcements = [
    { title: 'Holiday Notice', message: 'School will remain closed on Feb 26th for Republic Day celebrations.', time: 'Posted 2 hours ago', priority: 'high' },
    { title: 'Fee Reminder', message: 'Last date for fee submission is Feb 28th. Late fee will be applicable.', time: 'Posted yesterday', priority: 'normal' },
  ];

  getTypeBadgeBg(type: string): string {
    const colors: Record<string, string> = {
      Enrollment: '#eef2ff',
      Attendance: '#dcfce7',
      Finance: '#fef3c7',
      Academic: '#e0f2fe',
      Event: '#f1f5f9',
    };
    return colors[type] || '#f1f5f9';
  }

  getTypeBadgeColor(type: string): string {
    const colors: Record<string, string> = {
      Enrollment: '#4338ca',
      Attendance: '#047857',
      Finance: '#b45309',
      Academic: '#0369a1',
      Event: '#475569',
    };
    return colors[type] || '#475569';
  }
}
