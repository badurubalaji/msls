/**
 * MSLS Enquiries List Component
 *
 * Main component for viewing and managing admission enquiries.
 * Features filtering, searching, and actions for each enquiry.
 */

import {
  Component,
  OnInit,
  inject,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { ToastService } from '../../../shared/services/toast.service';

import { EnquiryService } from './enquiry.service';
import { EnquiryFormComponent } from './enquiry-form.component';
import { FollowUpFormComponent } from './follow-up-form.component';
import {
  Enquiry,
  EnquiryStatus,
  EnquiryFilterParams,
  ENQUIRY_STATUS_CONFIG,
} from './enquiry.model';

/**
 * EnquiriesComponent - List and manage admission enquiries.
 */
@Component({
  selector: 'msls-enquiries',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsModalComponent,
    EnquiryFormComponent,
    FollowUpFormComponent,
  ],
  template: `
    <div class="enquiries-page">
      <!-- Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-left">
            <div class="header-icon">
              <i class="fa-solid fa-clipboard-question"></i>
            </div>
            <div>
              <h1 class="header-title">Admission Enquiries</h1>
              <p class="header-subtitle">Track and manage prospective student enquiries</p>
            </div>
          </div>
          <button class="btn-new" (click)="openEnquiryForm()">
            <i class="fa-solid fa-plus"></i>
            New Enquiry
          </button>
        </div>

        <!-- Stats Row -->
        <div class="stats-row">
          <button
            *ngFor="let stat of statusStats()"
            class="stat-card"
            [class.active]="filters().status === stat.status"
            (click)="filterByStatus(stat.status)"
          >
            <div class="stat-icon">
              <i [class]="stat.icon"></i>
            </div>
            <div class="stat-info">
              <div class="stat-count">{{ stat.count }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </button>
        </div>
      </div>

      <!-- Main Content -->
      <div class="main-content">
        <div class="content-card">
          <!-- Search Bar -->
          <div class="search-section">
            <div class="search-row">
              <div class="search-input-wrapper">
                <i class="fa-solid fa-search search-icon"></i>
                <input
                  type="text"
                  [value]="searchTerm()"
                  (input)="onSearchInput($event)"
                  (keyup.enter)="applyFilters()"
                  placeholder="Search by student name, parent name, or phone..."
                  class="search-input"
                />
              </div>
              <div class="search-actions">
                <button class="btn-filter" (click)="toggleFilters()">
                  <i class="fa-solid fa-sliders"></i>
                  <span>Filters</span>
                  <span *ngIf="hasActiveFilters()" class="filter-badge">{{ activeFilterCount() }}</span>
                </button>
                <button class="btn-search" (click)="applyFilters()">
                  <i class="fa-solid fa-search"></i>
                  Search
                </button>
              </div>
            </div>
          </div>

          <!-- Expandable Filters -->
          <div *ngIf="showFilters()" class="filters-section">
            <div class="filters-grid">
              <div class="filter-field">
                <label>Status</label>
                <select [value]="filters().status || ''" (change)="onStatusSelectChange($event)">
                  <option value="">All Status</option>
                  <option value="new">New</option>
                  <option value="contacted">Contacted</option>
                  <option value="interested">Interested</option>
                  <option value="converted">Converted</option>
                  <option value="closed">Closed</option>
                </select>
              </div>
              <div class="filter-field">
                <label>From Date</label>
                <input type="date" [value]="filters().fromDate || ''" (change)="onFromDateChange($event)" />
              </div>
              <div class="filter-field">
                <label>To Date</label>
                <input type="date" [value]="filters().toDate || ''" (change)="onToDateChange($event)" />
              </div>
              <div class="filter-field filter-actions">
                <button class="btn-clear" (click)="clearFilters()">
                  <i class="fa-solid fa-times"></i>
                  Clear All
                </button>
              </div>
            </div>
          </div>

          <!-- Results Info -->
          <div class="results-info">
            <div class="results-count">
              <ng-container *ngIf="enquiryService.loading(); else showCount">
                <i class="fa-solid fa-spinner fa-spin"></i>
                Loading enquiries...
              </ng-container>
              <ng-template #showCount>
                Showing <strong>{{ tableData().length }}</strong> of
                <strong>{{ enquiryService.totalItems() }}</strong> enquiries
              </ng-template>
            </div>
            <button *ngIf="filters().status" class="btn-clear-filter" (click)="clearStatusFilter()">
              <i class="fa-solid fa-times"></i>
              Clear "{{ getStatusLabel(filters().status!) }}" filter
            </button>
          </div>

          <!-- Enquiries List -->
          <div class="enquiries-list">
            <!-- Loading State -->
            <ng-container *ngIf="enquiryService.loading()">
              <div *ngFor="let i of [1,2,3,4,5]" class="skeleton-item">
                <div class="skeleton-avatar"></div>
                <div class="skeleton-content">
                  <div class="skeleton-line skeleton-title"></div>
                  <div class="skeleton-line skeleton-subtitle"></div>
                </div>
                <div class="skeleton-badge"></div>
              </div>
            </ng-container>

            <!-- Error State -->
            <div *ngIf="!enquiryService.loading() && enquiryService.error()" class="empty-state">
              <div class="empty-icon error">
                <i class="fa-solid fa-triangle-exclamation"></i>
              </div>
              <h3>Failed to load enquiries</h3>
              <p>{{ enquiryService.error() }}</p>
              <button class="btn-primary" (click)="loadEnquiries()">
                <i class="fa-solid fa-rotate"></i>
                Retry
              </button>
            </div>

            <!-- Empty State -->
            <div *ngIf="!enquiryService.loading() && !enquiryService.error() && tableData().length === 0" class="empty-state">
              <div class="empty-icon">
                <i class="fa-solid fa-clipboard-question"></i>
              </div>
              <h3>No enquiries found</h3>
              <p *ngIf="hasActiveFilters()">Try adjusting your filters or search terms</p>
              <p *ngIf="!hasActiveFilters()">Get started by creating your first enquiry</p>
              <button *ngIf="hasActiveFilters()" class="btn-secondary" (click)="clearFilters()">
                <i class="fa-solid fa-times"></i>
                Clear Filters
              </button>
              <button *ngIf="!hasActiveFilters()" class="btn-primary" (click)="openEnquiryForm()">
                <i class="fa-solid fa-plus"></i>
                New Enquiry
              </button>
            </div>

            <!-- Enquiry Cards -->
            <ng-container *ngIf="!enquiryService.loading() && !enquiryService.error() && tableData().length > 0">
              <div
                *ngFor="let enquiry of tableData()"
                class="enquiry-card"
                (click)="editEnquiry(enquiry)"
              >
                <div class="enquiry-avatar" [style.background]="getAvatarColor(enquiry.studentName)">
                  {{ getInitials(enquiry.studentName) }}
                </div>
                <div class="enquiry-content">
                  <div class="enquiry-header">
                    <div class="enquiry-main">
                      <h3 class="enquiry-name">{{ enquiry.studentName }}</h3>
                      <span class="enquiry-number">{{ enquiry.enquiryNumber }}</span>
                    </div>
                    <span
                      class="status-badge"
                      [style.background]="getStatusBgColor(enquiry.status)"
                      [style.color]="getStatusTextColor(enquiry.status)"
                    >
                      <span class="status-dot" [style.background]="getStatusTextColor(enquiry.status)"></span>
                      {{ getStatusLabel(enquiry.status) }}
                    </span>
                  </div>
                  <div class="enquiry-details">
                    <span><i class="fa-solid fa-user"></i> {{ enquiry.parentName }}</span>
                    <span><i class="fa-solid fa-phone"></i> {{ enquiry.parentPhone }}</span>
                    <span><i class="fa-solid fa-graduation-cap"></i> {{ enquiry.classApplying }}</span>
                  </div>
                  <div class="enquiry-footer">
                    <div class="enquiry-meta">
                      <span><i class="fa-regular fa-calendar"></i> {{ formatDate(enquiry.createdAt) }}</span>
                      <span *ngIf="enquiry.followUpDate" class="follow-up" [class.overdue]="isOverdue(enquiry.followUpDate)">
                        <i class="fa-solid fa-bell"></i>
                        Follow-up: {{ formatDate(enquiry.followUpDate) }}
                        <strong *ngIf="isOverdue(enquiry.followUpDate)">(Overdue)</strong>
                      </span>
                      <span *ngIf="enquiry.source"><i class="fa-solid fa-location-dot"></i> {{ enquiry.source }}</span>
                    </div>
                    <div class="enquiry-actions">
                      <button class="action-btn" (click)="openFollowUpForm(enquiry, $event)" title="Add Follow-up">
                        <i class="fa-solid fa-phone-volume"></i>
                      </button>
                      <button class="action-btn" (click)="editEnquiry(enquiry, $event)" title="Edit">
                        <i class="fa-solid fa-pen"></i>
                      </button>
                      <button
                        *ngIf="enquiry.status !== 'converted' && enquiry.status !== 'closed'"
                        class="action-btn convert"
                        (click)="convertToApplication(enquiry, $event)"
                        title="Convert to Application"
                      >
                        <i class="fa-solid fa-arrow-right-to-bracket"></i>
                      </button>
                      <button class="action-btn delete" (click)="confirmDelete(enquiry, $event)" title="Delete">
                        <i class="fa-solid fa-trash"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </ng-container>
          </div>

          <!-- Pagination -->
          <div *ngIf="enquiryService.totalPages() > 1" class="pagination">
            <div class="pagination-info">
              Page <strong>{{ enquiryService.currentPage() }}</strong> of <strong>{{ enquiryService.totalPages() }}</strong>
            </div>
            <div class="pagination-controls">
              <button
                class="page-btn"
                [disabled]="!enquiryService.hasPreviousPage()"
                (click)="goToPage(enquiryService.currentPage() - 1)"
              >
                <i class="fa-solid fa-chevron-left"></i>
              </button>
              <ng-container *ngFor="let page of getVisiblePages()">
                <span *ngIf="page === '...'" class="page-ellipsis">...</span>
                <button
                  *ngIf="page !== '...'"
                  class="page-btn"
                  [class.active]="enquiryService.currentPage() === +page"
                  (click)="goToPage(+page)"
                >
                  {{ page }}
                </button>
              </ng-container>
              <button
                class="page-btn"
                [disabled]="!enquiryService.hasNextPage()"
                (click)="goToPage(enquiryService.currentPage() + 1)"
              >
                <i class="fa-solid fa-chevron-right"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Enquiry Form Modal -->
      <msls-modal
        [isOpen]="showEnquiryForm()"
        [title]="editingEnquiry() ? 'Edit Enquiry' : 'New Enquiry'"
        size="lg"
        (closed)="closeEnquiryForm()"
      >
        <ng-container modal-body>
          <msls-enquiry-form
            [enquiry]="editingEnquiry()"
            (submitted)="onEnquirySubmitted($event)"
            (cancelled)="closeEnquiryForm()"
          />
        </ng-container>
      </msls-modal>

      <!-- Follow-up Form Modal -->
      <msls-modal
        [isOpen]="showFollowUpForm()"
        title="Add Follow-up"
        size="md"
        (closed)="closeFollowUpForm()"
      >
        <ng-container modal-body>
          <msls-follow-up-form
            *ngIf="selectedEnquiryForFollowUp()"
            [enquiryId]="selectedEnquiryForFollowUp()!.id"
            (submitted)="onFollowUpSubmitted($event)"
            (cancelled)="closeFollowUpForm()"
          />
        </ng-container>
      </msls-modal>

      <!-- Delete Confirmation Modal -->
      <msls-modal
        [isOpen]="showDeleteConfirm()"
        title="Delete Enquiry"
        size="sm"
        (closed)="closeDeleteConfirm()"
      >
        <ng-container modal-body>
          <div class="delete-confirm">
            <div class="delete-icon">
              <i class="fa-solid fa-trash"></i>
            </div>
            <h3>Delete Enquiry?</h3>
            <p>
              Are you sure you want to delete the enquiry for
              <strong>{{ enquiryToDelete()?.studentName }}</strong>?
              This action cannot be undone.
            </p>
          </div>
        </ng-container>
        <ng-container modal-footer>
          <div class="delete-actions">
            <button class="btn-secondary" (click)="closeDeleteConfirm()">Cancel</button>
            <button class="btn-danger" [disabled]="deleting()" (click)="deleteEnquiry()">
              <i *ngIf="deleting()" class="fa-solid fa-spinner fa-spin"></i>
              Delete
            </button>
          </div>
        </ng-container>
      </msls-modal>

      <!-- Convert to Application Modal -->
      <msls-modal
        [isOpen]="showConvertConfirm()"
        title="Convert to Application"
        size="sm"
        (closed)="closeConvertConfirm()"
      >
        <ng-container modal-body>
          <div class="convert-confirm">
            <div class="convert-icon">
              <i class="fa-solid fa-arrow-right-to-bracket"></i>
            </div>
            <h3>Convert Enquiry to Application</h3>
            <p>
              Convert the enquiry for <strong>{{ enquiryToConvert()?.studentName }}</strong>
              into an admission application.
            </p>
            <div class="session-select" *ngIf="!enquiryToConvert()?.sessionId">
              <label for="sessionSelect">Select Admission Session <span class="required">*</span></label>
              <div *ngIf="loadingSessions()" class="loading-sessions">
                <i class="fa-solid fa-spinner fa-spin"></i> Loading sessions...
              </div>
              <select
                *ngIf="!loadingSessions()"
                id="sessionSelect"
                class="form-select"
                [value]="selectedSessionId()"
                (change)="onSessionSelect($event)"
              >
                <option value="">-- Select Session --</option>
                <option *ngFor="let session of availableSessions()" [value]="session.id">
                  {{ session.name }}
                </option>
              </select>
              <small *ngIf="!loadingSessions() && availableSessions().length === 0" class="no-sessions">
                No open admission sessions available. Please create one first.
              </small>
              <small *ngIf="!loadingSessions() && availableSessions().length > 0" class="hint">
                A session is required to create an application
              </small>
            </div>
          </div>
        </ng-container>
        <ng-container modal-footer>
          <div class="convert-actions">
            <button class="btn-secondary" (click)="closeConvertConfirm()">Cancel</button>
            <button
              class="btn-success"
              [disabled]="converting() || (!enquiryToConvert()?.sessionId && !selectedSessionId())"
              (click)="confirmConvert()"
            >
              <i *ngIf="converting()" class="fa-solid fa-spinner fa-spin"></i>
              Convert
            </button>
          </div>
        </ng-container>
      </msls-modal>
    </div>
  `,
  styles: [`
    .enquiries-page {
      min-height: 100vh;
      background: #f1f5f9;
    }

    /* Header */
    .page-header {
      background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
      padding: 1.5rem;
    }

    .header-content {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .header-left {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      background: rgba(255,255,255,0.2);
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 1.25rem;
    }

    .header-title {
      font-size: 1.5rem;
      font-weight: 700;
      color: white;
      margin: 0;
    }

    .header-subtitle {
      font-size: 0.875rem;
      color: rgba(255,255,255,0.8);
      margin: 0.25rem 0 0 0;
    }

    .btn-new {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      background: white;
      color: #6366f1;
      font-weight: 600;
      border: none;
      border-radius: 0.75rem;
      cursor: pointer;
      transition: transform 0.15s;
    }

    .btn-new:hover {
      transform: scale(1.05);
    }

    /* Stats Row */
    .stats-row {
      max-width: 1200px;
      margin: 1.5rem auto 0;
      display: grid;
      grid-template-columns: repeat(5, 1fr);
      gap: 0.75rem;
    }

    @media (max-width: 768px) {
      .stats-row {
        grid-template-columns: repeat(2, 1fr);
      }
    }

    .stat-card {
      background: rgba(255,255,255,0.15);
      border: none;
      border-radius: 0.75rem;
      padding: 1rem;
      display: flex;
      align-items: center;
      gap: 0.75rem;
      cursor: pointer;
      transition: all 0.15s;
      text-align: left;
    }

    .stat-card:hover {
      background: rgba(255,255,255,0.25);
      transform: scale(1.02);
    }

    .stat-card.active {
      background: rgba(255,255,255,0.3);
      box-shadow: 0 0 0 2px white;
    }

    .stat-icon {
      width: 2.5rem;
      height: 2.5rem;
      background: rgba(255,255,255,0.2);
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
    }

    .stat-count {
      font-size: 1.5rem;
      font-weight: 700;
      color: white;
    }

    .stat-label {
      font-size: 0.75rem;
      color: rgba(255,255,255,0.8);
    }

    /* Main Content */
    .main-content {
      max-width: 1200px;
      margin: 0 auto;
      padding: 1.5rem;
      margin-top: -1rem;
    }

    .content-card {
      background: white;
      border-radius: 1rem;
      box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1);
      overflow: hidden;
    }

    /* Search Section */
    .search-section {
      padding: 1rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .search-row {
      display: flex;
      gap: 1rem;
      flex-wrap: wrap;
    }

    .search-input-wrapper {
      flex: 1;
      min-width: 200px;
      position: relative;
    }

    .search-icon {
      position: absolute;
      left: 1rem;
      top: 50%;
      transform: translateY(-50%);
      color: #94a3b8;
    }

    .search-input {
      width: 100%;
      height: 2.75rem;
      padding: 0 1rem 0 2.75rem;
      border: 2px solid #e2e8f0;
      border-radius: 0.75rem;
      font-size: 0.875rem;
      transition: border-color 0.15s;
    }

    .search-input:focus {
      outline: none;
      border-color: #6366f1;
    }

    .search-actions {
      display: flex;
      gap: 0.5rem;
    }

    .btn-filter {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      height: 2.75rem;
      padding: 0 1rem;
      background: white;
      border: 2px solid #e2e8f0;
      border-radius: 0.75rem;
      cursor: pointer;
      transition: all 0.15s;
    }

    .btn-filter:hover {
      border-color: #6366f1;
      color: #6366f1;
    }

    .filter-badge {
      width: 1.25rem;
      height: 1.25rem;
      background: #6366f1;
      color: white;
      font-size: 0.75rem;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .btn-search {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      height: 2.75rem;
      padding: 0 1.5rem;
      background: #6366f1;
      color: white;
      border: none;
      border-radius: 0.75rem;
      font-weight: 500;
      cursor: pointer;
      transition: background 0.15s;
    }

    .btn-search:hover {
      background: #4f46e5;
    }

    /* Filters Section */
    .filters-section {
      padding: 1rem;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .filters-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 1rem;
    }

    @media (max-width: 768px) {
      .filters-grid {
        grid-template-columns: 1fr;
      }
    }

    .filter-field label {
      display: block;
      font-size: 0.75rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      margin-bottom: 0.5rem;
    }

    .filter-field select,
    .filter-field input {
      width: 100%;
      height: 2.5rem;
      padding: 0 0.75rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
    }

    .filter-field select:focus,
    .filter-field input:focus {
      outline: none;
      border-color: #6366f1;
    }

    .filter-actions {
      display: flex;
      align-items: flex-end;
    }

    .btn-clear {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      height: 2.5rem;
      padding: 0 1rem;
      background: transparent;
      border: none;
      color: #64748b;
      cursor: pointer;
      transition: color 0.15s;
    }

    .btn-clear:hover {
      color: #ef4444;
    }

    /* Results Info */
    .results-info {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.75rem 1rem;
      background: #f8fafc;
      font-size: 0.875rem;
      color: #64748b;
    }

    .results-count strong {
      color: #0f172a;
    }

    .btn-clear-filter {
      display: flex;
      align-items: center;
      gap: 0.25rem;
      background: transparent;
      border: none;
      color: #6366f1;
      cursor: pointer;
    }

    .btn-clear-filter:hover {
      text-decoration: underline;
    }

    /* Enquiries List */
    .enquiries-list {
      min-height: 200px;
    }

    /* Skeleton Loading */
    .skeleton-item {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
    }

    .skeleton-avatar {
      width: 3rem;
      height: 3rem;
      background: #e2e8f0;
      border-radius: 50%;
      animation: pulse 1.5s infinite;
    }

    .skeleton-content {
      flex: 1;
    }

    .skeleton-line {
      height: 0.75rem;
      background: #e2e8f0;
      border-radius: 0.25rem;
      animation: pulse 1.5s infinite;
    }

    .skeleton-title {
      width: 40%;
      margin-bottom: 0.5rem;
    }

    .skeleton-subtitle {
      width: 25%;
    }

    .skeleton-badge {
      width: 5rem;
      height: 1.5rem;
      background: #e2e8f0;
      border-radius: 1rem;
      animation: pulse 1.5s infinite;
    }

    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.5; }
    }

    /* Empty State */
    .empty-state {
      padding: 4rem 2rem;
      text-align: center;
    }

    .empty-icon {
      width: 5rem;
      height: 5rem;
      margin: 0 auto 1rem;
      background: linear-gradient(135deg, #e0e7ff 0%, #f3e8ff 100%);
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 2rem;
      color: #6366f1;
    }

    .empty-icon.error {
      background: #fee2e2;
      color: #ef4444;
    }

    .empty-state h3 {
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
    }

    .empty-state p {
      color: #64748b;
      margin: 0 0 1rem 0;
    }

    .btn-primary,
    .btn-secondary {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.15s;
    }

    .btn-primary {
      background: #6366f1;
      color: white;
      border: none;
    }

    .btn-primary:hover {
      background: #4f46e5;
    }

    .btn-secondary {
      background: white;
      color: #64748b;
      border: 1px solid #e2e8f0;
    }

    .btn-secondary:hover {
      background: #f8fafc;
    }

    /* Enquiry Card */
    .enquiry-card {
      display: flex;
      gap: 1rem;
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      cursor: pointer;
      transition: background 0.15s;
    }

    .enquiry-card:hover {
      background: #f8fafc;
    }

    .enquiry-card:hover .enquiry-actions {
      opacity: 1;
    }

    .enquiry-avatar {
      width: 3rem;
      height: 3rem;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-weight: 700;
      font-size: 1rem;
      flex-shrink: 0;
    }

    .enquiry-content {
      flex: 1;
      min-width: 0;
    }

    .enquiry-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 1rem;
      margin-bottom: 0.5rem;
    }

    .enquiry-main {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      flex-wrap: wrap;
    }

    .enquiry-name {
      font-weight: 600;
      color: #0f172a;
      margin: 0;
    }

    .enquiry-number {
      font-size: 0.75rem;
      padding: 0.125rem 0.5rem;
      background: #f1f5f9;
      color: #64748b;
      border-radius: 1rem;
    }

    .status-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.25rem 0.75rem;
      border-radius: 1rem;
      font-size: 0.75rem;
      font-weight: 500;
      flex-shrink: 0;
    }

    .status-dot {
      width: 0.375rem;
      height: 0.375rem;
      border-radius: 50%;
    }

    .enquiry-details {
      display: flex;
      flex-wrap: wrap;
      gap: 1rem;
      font-size: 0.875rem;
      color: #64748b;
      margin-bottom: 0.75rem;
    }

    .enquiry-details span {
      display: flex;
      align-items: center;
      gap: 0.375rem;
    }

    .enquiry-details i {
      font-size: 0.75rem;
      width: 1rem;
    }

    .enquiry-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-top: 0.75rem;
      border-top: 1px solid #f1f5f9;
    }

    .enquiry-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 1rem;
      font-size: 0.75rem;
      color: #94a3b8;
    }

    .enquiry-meta span {
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .follow-up {
      padding: 0.125rem 0.5rem;
      background: #fef3c7;
      color: #d97706;
      border-radius: 0.25rem;
    }

    .follow-up.overdue {
      background: #fee2e2;
      color: #dc2626;
    }

    .enquiry-actions {
      display: flex;
      gap: 0.25rem;
      opacity: 0;
      transition: opacity 0.15s;
    }

    .action-btn {
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: transparent;
      border: none;
      border-radius: 0.5rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.15s;
    }

    .action-btn:hover {
      background: #f1f5f9;
      color: #0f172a;
    }

    .action-btn.convert:hover {
      background: #d1fae5;
      color: #059669;
    }

    .action-btn.delete:hover {
      background: #fee2e2;
      color: #dc2626;
    }

    /* Pagination */
    .pagination {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem;
      border-top: 1px solid #e2e8f0;
    }

    .pagination-info {
      font-size: 0.875rem;
      color: #64748b;
    }

    .pagination-controls {
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }

    .page-btn {
      width: 2.25rem;
      height: 2.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.15s;
    }

    .page-btn:hover:not(:disabled) {
      background: #f8fafc;
      border-color: #cbd5e1;
    }

    .page-btn.active {
      background: #6366f1;
      border-color: #6366f1;
      color: white;
    }

    .page-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .page-ellipsis {
      padding: 0 0.5rem;
      color: #94a3b8;
    }

    /* Delete Modal */
    .delete-confirm {
      text-align: center;
      padding: 1rem;
    }

    .delete-icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      background: #fee2e2;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
      color: #ef4444;
    }

    .delete-confirm h3 {
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
    }

    .delete-confirm p {
      color: #64748b;
      margin: 0;
    }

    .delete-confirm strong {
      color: #0f172a;
    }

    .delete-actions {
      display: flex;
      justify-content: center;
      gap: 0.75rem;
    }

    .btn-danger {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      background: #ef4444;
      color: white;
      border: none;
      border-radius: 0.5rem;
      font-weight: 500;
      cursor: pointer;
      transition: background 0.15s;
    }

    .btn-danger:hover:not(:disabled) {
      background: #dc2626;
    }

    .btn-danger:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    /* Convert Modal Styles */
    .convert-confirm {
      text-align: center;
      padding: 1rem;
    }

    .convert-icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      background: #d1fae5;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
      color: #059669;
    }

    .convert-confirm h3 {
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
    }

    .convert-confirm p {
      color: #64748b;
      margin: 0 0 1rem 0;
    }

    .convert-confirm strong {
      color: #0f172a;
    }

    .session-select {
      text-align: left;
      margin-top: 1rem;
      padding: 1rem;
      background: #f8fafc;
      border-radius: 0.5rem;
    }

    .session-select label {
      display: block;
      font-weight: 500;
      color: #374151;
      margin-bottom: 0.5rem;
    }

    .session-select .required {
      color: #ef4444;
    }

    .session-select .form-select {
      width: 100%;
      padding: 0.625rem 0.75rem;
      border: 1px solid #d1d5db;
      border-radius: 0.375rem;
      font-size: 0.875rem;
      background: white;
    }

    .session-select .form-select:focus {
      outline: none;
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }

    .session-select .hint {
      display: block;
      margin-top: 0.5rem;
      color: #6b7280;
      font-size: 0.75rem;
    }

    .session-select .loading-sessions {
      padding: 0.75rem;
      text-align: center;
      color: #6b7280;
    }

    .session-select .no-sessions {
      display: block;
      margin-top: 0.5rem;
      color: #dc2626;
      font-size: 0.75rem;
    }

    .convert-actions {
      display: flex;
      justify-content: center;
      gap: 0.75rem;
    }

    .btn-success {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      background: #059669;
      color: white;
      border: none;
      border-radius: 0.5rem;
      font-weight: 500;
      cursor: pointer;
      transition: background 0.15s;
    }

    .btn-success:hover:not(:disabled) {
      background: #047857;
    }

    .btn-success:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EnquiriesComponent implements OnInit {
  readonly enquiryService = inject(EnquiryService);
  private readonly toast = inject(ToastService);
  private readonly router = inject(Router);
  private readonly http = inject(HttpClient);

  // Local state
  readonly searchTerm = signal('');
  readonly filters = signal<EnquiryFilterParams>({});
  readonly showFilters = signal(false);
  readonly showEnquiryForm = signal(false);
  readonly showFollowUpForm = signal(false);
  readonly showDeleteConfirm = signal(false);
  readonly showConvertConfirm = signal(false);
  readonly editingEnquiry = signal<Enquiry | null>(null);
  readonly selectedEnquiryForFollowUp = signal<Enquiry | null>(null);
  readonly enquiryToDelete = signal<Enquiry | null>(null);
  readonly enquiryToConvert = signal<Enquiry | null>(null);
  readonly deleting = signal(false);
  readonly converting = signal(false);
  readonly selectedSessionId = signal<string>('');
  readonly availableSessions = signal<Array<{ id: string; name: string }>>([]);
  readonly loadingSessions = signal(false);

  // Computed table data
  readonly tableData = computed(() => this.enquiryService.enquiries());

  // Check if any filters are active
  readonly hasActiveFilters = computed(() => {
    const f = this.filters();
    return !!(f.status || f.fromDate || f.toDate || this.searchTerm());
  });

  // Count active filters
  readonly activeFilterCount = computed(() => {
    const f = this.filters();
    let count = 0;
    if (f.status) count++;
    if (f.fromDate) count++;
    if (f.toDate) count++;
    if (this.searchTerm()) count++;
    return count;
  });

  // Status statistics
  readonly statusStats = computed(() => {
    const enquiries = this.enquiryService.enquiries();
    const stats: Record<EnquiryStatus, number> = {
      new: 0,
      contacted: 0,
      interested: 0,
      converted: 0,
      closed: 0,
    };

    if (Array.isArray(enquiries)) {
      enquiries.forEach((e) => {
        if (e && e.status && stats[e.status] !== undefined) {
          stats[e.status]++;
        }
      });
    }

    return [
      { status: 'new' as const, label: 'New', count: stats.new, icon: 'fa-solid fa-sparkles' },
      { status: 'contacted' as const, label: 'Contacted', count: stats.contacted, icon: 'fa-solid fa-phone' },
      { status: 'interested' as const, label: 'Interested', count: stats.interested, icon: 'fa-solid fa-heart' },
      { status: 'converted' as const, label: 'Converted', count: stats.converted, icon: 'fa-solid fa-check-circle' },
      { status: 'closed' as const, label: 'Closed', count: stats.closed, icon: 'fa-solid fa-times-circle' },
    ];
  });

  // Avatar colors
  private readonly avatarColors = [
    '#6366f1', '#8b5cf6', '#a855f7', '#ec4899', '#f43f5e',
    '#ef4444', '#f97316', '#eab308', '#22c55e', '#14b8a6',
    '#06b6d4', '#0ea5e9', '#3b82f6',
  ];

  ngOnInit(): void {
    this.loadEnquiries();
  }

  loadEnquiries(): void {
    this.enquiryService
      .loadEnquiries(this.filters(), {
        page: this.enquiryService.currentPage(),
        pageSize: this.enquiryService.pageSize(),
      })
      .subscribe();
  }

  toggleFilters(): void {
    this.showFilters.set(!this.showFilters());
  }

  onSearchInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.searchTerm.set(target.value);
  }

  onStatusSelectChange(event: Event): void {
    const target = event.target as HTMLSelectElement;
    this.filters.update((f) => ({
      ...f,
      status: target.value ? (target.value as EnquiryStatus) : undefined,
    }));
  }

  onFromDateChange(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.filters.update((f) => ({ ...f, fromDate: target.value || undefined }));
  }

  onToDateChange(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.filters.update((f) => ({ ...f, toDate: target.value || undefined }));
  }

  applyFilters(): void {
    this.filters.update((f) => ({ ...f, search: this.searchTerm() || undefined }));
    this.loadEnquiries();
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.filters.set({});
    this.loadEnquiries();
  }

  clearStatusFilter(): void {
    this.filters.update((f) => ({ ...f, status: undefined }));
    this.loadEnquiries();
  }

  filterByStatus(status: EnquiryStatus): void {
    if (this.filters().status === status) {
      this.clearStatusFilter();
    } else {
      this.filters.update((f) => ({ ...f, status }));
      this.loadEnquiries();
    }
  }

  goToPage(page: number): void {
    this.enquiryService
      .loadEnquiries(this.filters(), { page, pageSize: this.enquiryService.pageSize() })
      .subscribe();
  }

  getVisiblePages(): (number | string)[] {
    const total = this.enquiryService.totalPages();
    const current = this.enquiryService.currentPage();
    const pages: (number | string)[] = [];

    if (total <= 7) {
      for (let i = 1; i <= total; i++) pages.push(i);
    } else {
      pages.push(1);
      if (current > 3) pages.push('...');
      for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
        pages.push(i);
      }
      if (current < total - 2) pages.push('...');
      pages.push(total);
    }

    return pages;
  }

  openEnquiryForm(): void {
    this.editingEnquiry.set(null);
    this.showEnquiryForm.set(true);
  }

  editEnquiry(enquiry: Enquiry, event?: Event): void {
    event?.stopPropagation();
    this.editingEnquiry.set(enquiry);
    this.showEnquiryForm.set(true);
  }

  closeEnquiryForm(): void {
    this.showEnquiryForm.set(false);
    this.editingEnquiry.set(null);
  }

  onEnquirySubmitted(enquiry: Enquiry): void {
    this.closeEnquiryForm();
    this.toast.success(this.editingEnquiry() ? 'Enquiry updated' : 'Enquiry created');
    this.loadEnquiries();
  }

  openFollowUpForm(enquiry: Enquiry, event?: Event): void {
    event?.stopPropagation();
    this.selectedEnquiryForFollowUp.set(enquiry);
    this.showFollowUpForm.set(true);
  }

  closeFollowUpForm(): void {
    this.showFollowUpForm.set(false);
    this.selectedEnquiryForFollowUp.set(null);
  }

  onFollowUpSubmitted(followUp: unknown): void {
    this.closeFollowUpForm();
    this.toast.success('Follow-up added');
    this.loadEnquiries();
  }

  convertToApplication(enquiry: Enquiry, event?: Event): void {
    event?.stopPropagation();
    this.enquiryToConvert.set(enquiry);
    this.selectedSessionId.set(enquiry.sessionId || '');

    // Load available sessions if enquiry doesn't have one
    if (!enquiry.sessionId) {
      this.loadAvailableSessions();
    }

    this.showConvertConfirm.set(true);
  }

  closeConvertConfirm(): void {
    this.showConvertConfirm.set(false);
    this.enquiryToConvert.set(null);
    this.selectedSessionId.set('');
  }

  onSessionSelect(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.selectedSessionId.set(select.value);
  }

  confirmConvert(): void {
    const enquiry = this.enquiryToConvert();
    if (!enquiry) return;

    const sessionId = enquiry.sessionId || this.selectedSessionId();
    if (!sessionId) {
      this.toast.error('Please select an admission session');
      return;
    }

    this.converting.set(true);
    this.enquiryService.convertToApplication(enquiry.id, { sessionId }).subscribe({
      next: () => {
        this.toast.success('Enquiry converted to application successfully!');
        this.closeConvertConfirm();
        this.loadEnquiries();
      },
      error: (err) => {
        const message = err?.error?.message || err?.message || 'Failed to convert enquiry';
        this.toast.error(message);
      },
      complete: () => this.converting.set(false),
    });
  }

  private loadAvailableSessions(): void {
    interface SessionResponse {
      id: string;
      name: string;
      status: string;
    }
    interface ApiResponse {
      success: boolean;
      data: {
        sessions: SessionResponse[];
        total: number;
      };
    }

    this.loadingSessions.set(true);
    this.http.get<ApiResponse>('/api/v1/admission-sessions')
      .subscribe({
        next: (response) => {
          // Filter to open or draft sessions (sessions that can accept applications)
          const sessions = response.data?.sessions || [];
          const availableSessions = sessions
            .filter(s => s.status === 'open' || s.status === 'draft')
            .map(s => ({ id: s.id, name: `${s.name} (${s.status})` }));
          this.availableSessions.set(availableSessions);
          this.loadingSessions.set(false);
        },
        error: (err) => {
          console.error('Failed to load sessions:', err);
          this.availableSessions.set([]);
          this.loadingSessions.set(false);
        }
      });
  }

  confirmDelete(enquiry: Enquiry, event?: Event): void {
    event?.stopPropagation();
    this.enquiryToDelete.set(enquiry);
    this.showDeleteConfirm.set(true);
  }

  closeDeleteConfirm(): void {
    this.showDeleteConfirm.set(false);
    this.enquiryToDelete.set(null);
  }

  deleteEnquiry(): void {
    const enquiry = this.enquiryToDelete();
    if (!enquiry) return;

    this.deleting.set(true);
    this.enquiryService.deleteEnquiry(enquiry.id).subscribe({
      next: () => {
        this.closeDeleteConfirm();
        this.toast.success('Enquiry deleted');
        this.deleting.set(false);
      },
      error: () => {
        this.toast.error('Failed to delete enquiry');
        this.deleting.set(false);
      },
    });
  }

  getStatusLabel(status: string): string {
    return ENQUIRY_STATUS_CONFIG[status as EnquiryStatus]?.label || status;
  }

  getStatusBgColor(status: string): string {
    const colors: Record<string, string> = {
      new: '#eff6ff',
      contacted: '#f5f3ff',
      interested: '#fef3c7',
      converted: '#d1fae5',
      closed: '#f1f5f9',
    };
    return colors[status] || '#f1f5f9';
  }

  getStatusTextColor(status: string): string {
    const colors: Record<string, string> = {
      new: '#3b82f6',
      contacted: '#7c3aed',
      interested: '#d97706',
      converted: '#059669',
      closed: '#64748b',
    };
    return colors[status] || '#64748b';
  }

  formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-IN', { day: '2-digit', month: 'short', year: 'numeric' });
  }

  isOverdue(dateStr: string): boolean {
    if (!dateStr) return false;
    const date = new Date(dateStr);
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    return date < today;
  }

  getInitials(name: string): string {
    if (!name) return '?';
    const parts = name.split(' ').filter(Boolean);
    if (parts.length >= 2) {
      return (parts[0][0] + parts[1][0]).toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  }

  getAvatarColor(name: string): string {
    if (!name) return this.avatarColors[0];
    const index = name.charCodeAt(0) % this.avatarColors.length;
    return this.avatarColors[index];
  }
}
