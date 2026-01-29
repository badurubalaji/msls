import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ClassService } from '../services/class.service';
import { ClassWithSections, SectionStructure, CLASS_LEVELS } from '../academic.model';
import { BranchService } from '../../admin/branches/branch.service';
import { Branch } from '../../admin/branches/branch.model';

@Component({
  selector: 'msls-structure',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-sitemap"></i>
          </div>
          <div class="header-text">
            <h1>Academic Structure</h1>
            <p>Hierarchical view of classes and sections</p>
          </div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" (click)="toggleAllExpanded()">
            <i class="fa-solid" [class.fa-expand]="!allExpanded()" [class.fa-compress]="allExpanded()"></i>
            {{ allExpanded() ? 'Collapse All' : 'Expand All' }}
          </button>
        </div>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="filter-group">
          <label>Branch</label>
          <select
            class="filter-select"
            [ngModel]="selectedBranch()"
            (ngModelChange)="selectedBranch.set($event); loadStructure()"
          >
            <option value="">All Branches</option>
            @for (branch of branches(); track branch.id) {
              <option [value]="branch.id">{{ branch.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label>Status</label>
          <select
            class="filter-select"
            [ngModel]="statusFilter()"
            (ngModelChange)="statusFilter.set($event)"
          >
            <option value="all">All Classes</option>
            <option value="active">Active Only</option>
            <option value="inactive">Inactive Only</option>
          </select>
        </div>
      </div>

      <!-- Summary Stats -->
      <div class="summary-cards">
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fa-solid fa-graduation-cap"></i>
          </div>
          <div class="stat-content">
            <span class="stat-value">{{ totalClasses() }}</span>
            <span class="stat-label">Classes</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fa-solid fa-layer-group"></i>
          </div>
          <div class="stat-content">
            <span class="stat-value">{{ totalSections() }}</span>
            <span class="stat-label">Sections</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fa-solid fa-users"></i>
          </div>
          <div class="stat-content">
            <span class="stat-value">{{ totalStudents() }}</span>
            <span class="stat-label">Students</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fa-solid fa-chair"></i>
          </div>
          <div class="stat-content">
            <span class="stat-value">{{ overallCapacity() }}%</span>
            <span class="stat-label">Capacity Used</span>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading structure...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadStructure()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <div class="structure-tree">
            @for (cls of filteredClasses(); track cls.id; let i = $index) {
              <div class="class-node" [class.expanded]="expandedClasses().has(cls.id)">
                <div class="class-header" (click)="toggleClass(cls.id)">
                  <div class="expand-icon">
                    <i class="fa-solid" [class.fa-chevron-right]="!expandedClasses().has(cls.id)" [class.fa-chevron-down]="expandedClasses().has(cls.id)"></i>
                  </div>
                  <div class="class-icon">
                    <i class="fa-solid fa-graduation-cap"></i>
                  </div>
                  <div class="class-info">
                    <span class="class-name">{{ cls.name }}</span>
                    <span class="class-code">{{ cls.code }}</span>
                    @if (cls.hasStreams) {
                      <span class="stream-badge">Streams</span>
                    }
                  </div>
                  <div class="class-stats">
                    <div class="stat-item">
                      <i class="fa-solid fa-layer-group"></i>
                      <span>{{ cls.sections.length }} sections</span>
                    </div>
                    <div class="stat-item">
                      <i class="fa-solid fa-users"></i>
                      <span>{{ cls.totalStudents }} students</span>
                    </div>
                    <div class="capacity-bar-wrapper">
                      <div class="capacity-bar">
                        <div
                          class="capacity-fill"
                          [class.low]="getClassCapacityPercent(cls) < 50"
                          [class.medium]="getClassCapacityPercent(cls) >= 50 && getClassCapacityPercent(cls) < 80"
                          [class.high]="getClassCapacityPercent(cls) >= 80"
                          [style.width.%]="getClassCapacityPercent(cls)"
                        ></div>
                      </div>
                      <span class="capacity-text">{{ getClassCapacityPercent(cls) }}%</span>
                    </div>
                  </div>
                  <div class="class-status">
                    <span class="status-badge" [class.active]="cls.isActive" [class.inactive]="!cls.isActive">
                      {{ cls.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </div>
                </div>

                @if (expandedClasses().has(cls.id)) {
                  <div class="sections-container">
                    @if (cls.sections.length === 0) {
                      <div class="no-sections">
                        <i class="fa-regular fa-folder-open"></i>
                        <span>No sections created for this class</span>
                      </div>
                    } @else {
                      @for (section of cls.sections; track section.id) {
                        <div class="section-card">
                          <div class="section-header">
                            <div class="section-icon">
                              <i class="fa-solid fa-users-rectangle"></i>
                            </div>
                            <div class="section-info">
                              <span class="section-name">Section {{ section.name }}</span>
                              <span class="section-code">{{ section.code }}</span>
                            </div>
                            @if (section.streamName) {
                              <span class="stream-tag">{{ section.streamName }}</span>
                            }
                          </div>

                          <div class="section-details">
                            <div class="detail-row">
                              <div class="detail-item">
                                <i class="fa-solid fa-chalkboard-user"></i>
                                <span class="detail-label">Class Teacher:</span>
                                <span class="detail-value">
                                  {{ section.classTeacherName || 'Not Assigned' }}
                                </span>
                              </div>
                              @if (section.roomNumber) {
                                <div class="detail-item">
                                  <i class="fa-solid fa-door-open"></i>
                                  <span class="detail-label">Room:</span>
                                  <span class="detail-value">{{ section.roomNumber }}</span>
                                </div>
                              }
                            </div>

                            <div class="capacity-section">
                              <div class="capacity-info">
                                <span class="student-count">{{ section.studentCount }}</span>
                                <span class="capacity-separator">/</span>
                                <span class="total-capacity">{{ section.capacity }}</span>
                                <span class="capacity-label">students</span>
                              </div>
                              <div class="capacity-progress">
                                <div class="capacity-bar large">
                                  <div
                                    class="capacity-fill"
                                    [class.low]="section.capacityUsage < 50"
                                    [class.medium]="section.capacityUsage >= 50 && section.capacityUsage < 80"
                                    [class.high]="section.capacityUsage >= 80"
                                    [style.width.%]="section.capacityUsage"
                                  ></div>
                                </div>
                                <span class="capacity-percent" [class.warning]="section.capacityUsage >= 90">
                                  {{ section.capacityUsage }}% utilized
                                </span>
                              </div>
                            </div>
                          </div>
                        </div>
                      }
                    }
                  </div>
                }
              </div>
            } @empty {
              <div class="empty-state">
                <i class="fa-regular fa-folder-open"></i>
                <p>No classes found</p>
                <p class="empty-hint">Create classes in the Classes page to see them here</p>
              </div>
            }
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      background: #f0fdf4;
      color: #16a34a;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .header-text h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .header-text p {
      margin: 0.25rem 0 0;
      color: #64748b;
      font-size: 0.875rem;
    }

    .filters-bar {
      display: flex;
      gap: 1rem;
      margin-bottom: 1rem;
      flex-wrap: wrap;
    }

    .filter-group {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .filter-group label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
    }

    .filter-select {
      padding: 0.5rem 2rem 0.5rem 0.75rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
      cursor: pointer;
      min-width: 180px;
    }

    .summary-cards {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .stat-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      padding: 1rem 1.25rem;
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .stat-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #f1f5f9;
      color: #64748b;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1rem;
    }

    .stat-content {
      display: flex;
      flex-direction: column;
    }

    .stat-value {
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .stat-label {
      font-size: 0.75rem;
      color: #64748b;
    }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .loading-container, .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-container {
      color: #dc2626;
      flex-direction: column;
    }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .structure-tree {
      padding: 1rem;
    }

    .class-node {
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      margin-bottom: 0.75rem;
      overflow: hidden;
      transition: box-shadow 0.2s;
    }

    .class-node:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    }

    .class-node.expanded {
      border-color: #c7d2fe;
    }

    .class-header {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 1rem;
      cursor: pointer;
      background: #f8fafc;
      transition: background-color 0.2s;
    }

    .class-header:hover {
      background: #f1f5f9;
    }

    .expand-icon {
      width: 1.5rem;
      height: 1.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #64748b;
      font-size: 0.75rem;
    }

    .class-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #eef2ff;
      color: #4f46e5;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .class-info {
      flex: 1;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .class-name {
      font-weight: 600;
      color: #1e293b;
    }

    .class-code {
      padding: 0.125rem 0.5rem;
      background: #f1f5f9;
      border-radius: 0.25rem;
      font-size: 0.75rem;
      font-family: monospace;
      color: #475569;
    }

    .stream-badge {
      padding: 0.125rem 0.5rem;
      background: #fef3c7;
      color: #92400e;
      border-radius: 0.25rem;
      font-size: 0.625rem;
      font-weight: 500;
      text-transform: uppercase;
    }

    .class-stats {
      display: flex;
      align-items: center;
      gap: 1.5rem;
    }

    .stat-item {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      font-size: 0.75rem;
      color: #64748b;
    }

    .stat-item i {
      font-size: 0.75rem;
    }

    .capacity-bar-wrapper {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .capacity-bar {
      width: 60px;
      height: 6px;
      background: #e2e8f0;
      border-radius: 3px;
      overflow: hidden;
    }

    .capacity-bar.large {
      width: 100px;
      height: 8px;
      border-radius: 4px;
    }

    .capacity-fill {
      height: 100%;
      border-radius: inherit;
      transition: width 0.3s ease;
    }

    .capacity-fill.low { background: #22c55e; }
    .capacity-fill.medium { background: #f59e0b; }
    .capacity-fill.high { background: #ef4444; }

    .capacity-text {
      font-size: 0.625rem;
      font-weight: 500;
      color: #64748b;
      min-width: 30px;
    }

    .class-status {
      margin-left: 1rem;
    }

    .status-badge {
      padding: 0.25rem 0.625rem;
      border-radius: 9999px;
      font-size: 0.625rem;
      font-weight: 500;
      text-transform: uppercase;
    }

    .status-badge.active {
      background: #dcfce7;
      color: #166534;
    }

    .status-badge.inactive {
      background: #f1f5f9;
      color: #64748b;
    }

    .sections-container {
      padding: 1rem 1rem 1rem 3.5rem;
      background: white;
      border-top: 1px solid #e2e8f0;
    }

    .no-sections {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 1rem;
      color: #9ca3af;
      font-size: 0.875rem;
    }

    .section-card {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      padding: 1rem;
      margin-bottom: 0.5rem;
    }

    .section-card:last-child {
      margin-bottom: 0;
    }

    .section-header {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      margin-bottom: 0.75rem;
    }

    .section-icon {
      width: 2rem;
      height: 2rem;
      border-radius: 0.375rem;
      background: #dbeafe;
      color: #2563eb;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.875rem;
    }

    .section-info {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      flex: 1;
    }

    .section-name {
      font-weight: 500;
      color: #1e293b;
    }

    .section-code {
      padding: 0.125rem 0.375rem;
      background: #e2e8f0;
      border-radius: 0.25rem;
      font-size: 0.625rem;
      font-family: monospace;
      color: #475569;
    }

    .stream-tag {
      padding: 0.125rem 0.5rem;
      background: #e0e7ff;
      color: #3730a3;
      border-radius: 0.25rem;
      font-size: 0.625rem;
      font-weight: 500;
    }

    .section-details {
      padding-left: 2.75rem;
    }

    .detail-row {
      display: flex;
      gap: 2rem;
      margin-bottom: 0.75rem;
    }

    .detail-item {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      font-size: 0.75rem;
    }

    .detail-item i {
      color: #9ca3af;
    }

    .detail-label {
      color: #64748b;
    }

    .detail-value {
      color: #1e293b;
      font-weight: 500;
    }

    .capacity-section {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0.75rem;
      background: white;
      border-radius: 0.375rem;
      border: 1px solid #e2e8f0;
    }

    .capacity-info {
      display: flex;
      align-items: baseline;
      gap: 0.25rem;
    }

    .student-count {
      font-size: 1.25rem;
      font-weight: 600;
      color: #1e293b;
    }

    .capacity-separator {
      color: #9ca3af;
    }

    .total-capacity {
      font-size: 0.875rem;
      color: #64748b;
    }

    .capacity-label {
      font-size: 0.75rem;
      color: #9ca3af;
      margin-left: 0.25rem;
    }

    .capacity-progress {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .capacity-percent {
      font-size: 0.75rem;
      color: #64748b;
      font-weight: 500;
    }

    .capacity-percent.warning {
      color: #dc2626;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.5rem;
      padding: 3rem;
      color: #64748b;
    }

    .empty-state i {
      font-size: 2.5rem;
      color: #cbd5e1;
    }

    .empty-state p {
      margin: 0;
    }

    .empty-hint {
      font-size: 0.875rem;
      color: #9ca3af;
    }

    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    @media (max-width: 1024px) {
      .summary-cards {
        grid-template-columns: repeat(2, 1fr);
      }
    }

    @media (max-width: 768px) {
      .page-header {
        flex-direction: column;
        gap: 1rem;
      }

      .summary-cards {
        grid-template-columns: 1fr;
      }

      .class-header {
        flex-wrap: wrap;
      }

      .class-stats {
        width: 100%;
        margin-top: 0.5rem;
        padding-left: 4rem;
      }

      .sections-container {
        padding-left: 1rem;
      }

      .section-details {
        padding-left: 0;
      }

      .detail-row {
        flex-direction: column;
        gap: 0.5rem;
      }

      .capacity-section {
        flex-direction: column;
        gap: 0.75rem;
        align-items: flex-start;
      }

      .capacity-progress {
        width: 100%;
      }

      .capacity-bar.large {
        flex: 1;
      }
    }
  `],
})
export class StructureComponent implements OnInit {
  private classService = inject(ClassService);
  private branchService = inject(BranchService);

  // Data signals
  classes = signal<ClassWithSections[]>([]);
  branches = signal<Branch[]>([]);
  expandedClasses = signal<Set<string>>(new Set());

  // State signals
  loading = signal(true);
  error = signal<string | null>(null);
  selectedBranch = signal<string>('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');

  // Computed values
  filteredClasses = computed(() => {
    let result = this.classes();
    const status = this.statusFilter();

    if (status === 'active') {
      result = result.filter(cls => cls.isActive);
    } else if (status === 'inactive') {
      result = result.filter(cls => !cls.isActive);
    }

    return result.sort((a, b) => a.displayOrder - b.displayOrder);
  });

  totalClasses = computed(() => this.filteredClasses().length);

  totalSections = computed(() =>
    this.filteredClasses().reduce((sum, cls) => sum + cls.sections.length, 0)
  );

  totalStudents = computed(() =>
    this.filteredClasses().reduce((sum, cls) => sum + cls.totalStudents, 0)
  );

  totalCapacity = computed(() =>
    this.filteredClasses().reduce((sum, cls) => sum + cls.totalCapacity, 0)
  );

  overallCapacity = computed(() => {
    const total = this.totalCapacity();
    const students = this.totalStudents();
    return total > 0 ? Math.round((students / total) * 100) : 0;
  });

  allExpanded = computed(() => {
    const expanded = this.expandedClasses();
    const allClasses = this.filteredClasses();
    return allClasses.length > 0 && allClasses.every(cls => expanded.has(cls.id));
  });

  ngOnInit(): void {
    this.loadBranches();
    this.loadStructure();
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: branches => this.branches.set(branches),
      error: () => console.error('Failed to load branches'),
    });
  }

  loadStructure(): void {
    this.loading.set(true);
    this.error.set(null);

    const branchId = this.selectedBranch() || undefined;

    this.classService.getClassStructure(branchId).subscribe({
      next: response => {
        this.classes.set(response.classes || []);
        this.loading.set(false);

        // Auto-expand first few classes if not many
        if (response.classes && response.classes.length <= 5) {
          const expanded = new Set<string>();
          response.classes.forEach(cls => expanded.add(cls.id));
          this.expandedClasses.set(expanded);
        }
      },
      error: () => {
        this.error.set('Failed to load class structure. Please try again.');
        this.loading.set(false);
      },
    });
  }

  toggleClass(classId: string): void {
    const expanded = new Set(this.expandedClasses());
    if (expanded.has(classId)) {
      expanded.delete(classId);
    } else {
      expanded.add(classId);
    }
    this.expandedClasses.set(expanded);
  }

  toggleAllExpanded(): void {
    if (this.allExpanded()) {
      this.expandedClasses.set(new Set());
    } else {
      const expanded = new Set<string>();
      this.filteredClasses().forEach(cls => expanded.add(cls.id));
      this.expandedClasses.set(expanded);
    }
  }

  getClassCapacityPercent(cls: ClassWithSections): number {
    return cls.totalCapacity > 0
      ? Math.round((cls.totalStudents / cls.totalCapacity) * 100)
      : 0;
  }
}
