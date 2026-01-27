/**
 * MSLS Merit List Component
 *
 * Displays merit list with ranking, scores, and actions for admission decisions.
 * Supports filtering by cutoff score and bulk selection for decisions.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import {
  MslsButtonComponent,
  MslsCardComponent,
  MslsBadgeComponent,
  MslsSelectComponent,
  MslsInputComponent,
  MslsTableComponent,
  MslsTableCellDirective,
  MslsModalComponent,
  MslsFormFieldComponent,
  TableColumn,
  SelectOption,
} from '../../../shared/components';
import { ToastService } from '../../../shared/services/toast.service';
import {
  MeritList,
  MeritListEntry,
  DecisionType,
  getStatusConfig,
  APPLICATION_STATUS_CONFIG,
} from './merit.model';
import { MeritService } from './merit.service';
import { AdmissionSessionService } from '../sessions/admission-session.service';
import { DecisionFormComponent } from './decision-form.component';

@Component({
  selector: 'msls-merit-list',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsButtonComponent,
    MslsCardComponent,
    MslsBadgeComponent,
    MslsSelectComponent,
    MslsModalComponent,
    MslsFormFieldComponent,
    DecisionFormComponent,
  ],
  template: `
    <div class="merit-list-page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">
            <i class="fa-solid fa-ranking-star text-primary-500"></i>
            Merit Lists
          </h1>
          <p class="page-subtitle">Generate and manage admission merit lists based on test results</p>
        </div>
      </div>

      <!-- Filters Card -->
      <msls-card class="filters-card">
        <div class="filters-row">
          <!-- Session Select -->
          <msls-form-field label="Admission Session" class="filter-field">
            <msls-select
              [options]="sessionOptions()"
              [placeholder]="'Select session'"
              (valueChange)="onSessionChange($event)"
            />
          </msls-form-field>

          <!-- Class Select -->
          <msls-form-field label="Class" class="filter-field">
            <msls-select
              [options]="classOptions()"
              [placeholder]="'Select class'"
              [disabled]="!selectedSession()"
              (valueChange)="onClassChange($event)"
            />
          </msls-form-field>

          <!-- Cutoff Score -->
          <msls-form-field label="Cutoff Score (%)" class="filter-field filter-field--small">
            <input
              type="number"
              class="cutoff-input"
              [value]="cutoffScore()"
              placeholder="0"
              (input)="onCutoffChange($event)"
            />
          </msls-form-field>

          <!-- Generate Button -->
          <div class="filter-actions">
            <msls-button
              variant="primary"
              [loading]="generating()"
              [disabled]="!canGenerate()"
              (click)="generateMeritList()"
            >
              <i class="fa-solid fa-wand-magic-sparkles"></i>
              Generate Merit List
            </msls-button>
          </div>
        </div>
      </msls-card>

      <!-- Merit List Content -->
      @if (meritList()) {
        <!-- Summary Stats -->
        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-icon stat-icon--total">
              <i class="fa-solid fa-users"></i>
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ meritList()!.entries.length }}</span>
              <span class="stat-label">Total Candidates</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon stat-icon--selected">
              <i class="fa-solid fa-circle-check"></i>
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ selectedCount() }}</span>
              <span class="stat-label">Selected</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon stat-icon--waitlist">
              <i class="fa-solid fa-clock"></i>
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ waitlistedCount() }}</span>
              <span class="stat-label">Waitlisted</span>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon stat-icon--pending">
              <i class="fa-solid fa-hourglass-half"></i>
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ pendingCount() }}</span>
              <span class="stat-label">Pending Decision</span>
            </div>
          </div>
        </div>

        <!-- Bulk Actions -->
        @if (selectedEntries().length > 0) {
          <div class="bulk-actions-bar">
            <span class="selected-count">{{ selectedEntries().length }} selected</span>
            <div class="bulk-buttons">
              <msls-button variant="secondary" size="sm" (click)="clearSelection()">
                Clear Selection
              </msls-button>
              <msls-button variant="primary" size="sm" (click)="openBulkDecisionModal('selected')">
                <i class="fa-solid fa-check"></i>
                Approve Selected
              </msls-button>
              <msls-button variant="secondary" size="sm" (click)="openBulkDecisionModal('waitlisted')">
                <i class="fa-solid fa-clock"></i>
                Waitlist Selected
              </msls-button>
              <msls-button variant="danger" size="sm" (click)="openBulkDecisionModal('rejected')">
                <i class="fa-solid fa-xmark"></i>
                Reject Selected
              </msls-button>
            </div>
          </div>
        }

        <!-- Merit List Table -->
        <msls-card class="table-card">
          <div class="table-header">
            <h2 class="table-title">
              Merit List - {{ meritList()!.className }}
              @if (meritList()!.cutoffScore) {
                <span class="cutoff-badge">Cutoff: {{ meritList()!.cutoffScore }}%</span>
              }
            </h2>
            <div class="table-actions">
              <label class="select-all-checkbox">
                <input
                  type="checkbox"
                  [checked]="allSelected()"
                  [indeterminate]="someSelected()"
                  (change)="toggleSelectAll($any($event.target).checked)"
                />
                <span>Select All</span>
              </label>
            </div>
          </div>

          <div class="table-container">
            <table class="merit-table">
              <thead>
                <tr>
                  <th class="col-checkbox">
                    <input
                      type="checkbox"
                      class="table-checkbox"
                      [checked]="allSelected()"
                      [indeterminate]="someSelected()"
                      (change)="toggleSelectAll($any($event.target).checked)"
                    />
                  </th>
                  <th class="col-rank">Rank</th>
                  <th class="col-name">Student Name</th>
                  <th class="col-app">Application #</th>
                  <th class="col-score">Score</th>
                  <th class="col-percent">%</th>
                  <th class="col-status">Status</th>
                  <th class="col-actions">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (entry of meritList()!.entries; track entry.id) {
                  <tr [class.selected-row]="isEntrySelected(entry.applicationId)">
                    <td class="col-checkbox">
                      <input
                        type="checkbox"
                        class="table-checkbox"
                        [checked]="isEntrySelected(entry.applicationId)"
                        (change)="toggleEntrySelection(entry.applicationId, $any($event.target).checked)"
                      />
                    </td>
                    <td class="col-rank">
                      <div class="rank-badge" [class]="getRankClass(entry.rank)">
                        @if (entry.rank === 1) {
                          <i class="fa-solid fa-trophy"></i>
                        } @else if (entry.rank === 2) {
                          <i class="fa-solid fa-medal"></i>
                        } @else if (entry.rank === 3) {
                          <i class="fa-solid fa-award"></i>
                        }
                        {{ entry.rank }}
                      </div>
                    </td>
                    <td class="col-name">
                      <div class="student-info">
                        <span class="student-name">{{ entry.studentName }}</span>
                        <span class="parent-info">{{ entry.parentName }} | {{ entry.parentPhone }}</span>
                      </div>
                    </td>
                    <td class="col-app">
                      <span class="app-number">{{ entry.applicationNumber }}</span>
                    </td>
                    <td class="col-score">
                      <div class="score-display">
                        <span class="score-value">{{ entry.totalScore }}</span>
                        <span class="score-max">/ {{ entry.maxScore }}</span>
                      </div>
                    </td>
                    <td class="col-percent">
                      <div class="percentage-bar">
                        <div class="percentage-fill" [style.width.%]="entry.percentage"></div>
                        <span class="percentage-text">{{ entry.percentage | number:'1.1-1' }}%</span>
                      </div>
                    </td>
                    <td class="col-status">
                      <msls-badge [variant]="getStatusConfig(entry.status).variant">
                        <i [class]="getStatusConfig(entry.status).icon"></i>
                        {{ getStatusConfig(entry.status).label }}
                      </msls-badge>
                    </td>
                    <td class="col-actions">
                      <div class="action-buttons">
                        @if (canMakeDecision(entry)) {
                          <button
                            class="action-btn action-btn--approve"
                            title="Approve"
                            (click)="openDecisionModal(entry, 'selected')"
                          >
                            <i class="fa-solid fa-check"></i>
                          </button>
                          <button
                            class="action-btn action-btn--waitlist"
                            title="Waitlist"
                            (click)="openDecisionModal(entry, 'waitlisted')"
                          >
                            <i class="fa-solid fa-clock"></i>
                          </button>
                          <button
                            class="action-btn action-btn--reject"
                            title="Reject"
                            (click)="openDecisionModal(entry, 'rejected')"
                          >
                            <i class="fa-solid fa-xmark"></i>
                          </button>
                        }
                        @if (entry.status === 'offer_sent' || entry.status === 'selected') {
                          <button
                            class="action-btn action-btn--letter"
                            title="View Offer Letter"
                            (click)="viewOfferLetter(entry)"
                          >
                            <i class="fa-solid fa-file-pdf"></i>
                          </button>
                        }
                        @if (entry.status === 'offer_accepted') {
                          <button
                            class="action-btn action-btn--enroll"
                            title="Complete Enrollment"
                            (click)="goToEnrollment(entry)"
                          >
                            <i class="fa-solid fa-user-plus"></i>
                          </button>
                        }
                        <button
                          class="action-btn action-btn--view"
                          title="View Details"
                          (click)="viewDetails(entry)"
                        >
                          <i class="fa-solid fa-eye"></i>
                        </button>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>

          @if (meritList()!.entries.length === 0) {
            <div class="empty-state">
              <i class="fa-regular fa-folder-open"></i>
              <p>No candidates found matching the criteria</p>
            </div>
          }
        </msls-card>
      } @else if (!selectedSession() || !selectedClass()) {
        <!-- No Selection State -->
        <msls-card class="empty-state-card">
          <div class="empty-state">
            <i class="fa-solid fa-list-check"></i>
            <h3>Select Session and Class</h3>
            <p>Choose an admission session and class to generate or view the merit list</p>
          </div>
        </msls-card>
      } @else if (loading()) {
        <!-- Loading State -->
        <msls-card class="empty-state-card">
          <div class="loading-state">
            <div class="spinner"></div>
            <p>Loading merit list...</p>
          </div>
        </msls-card>
      } @else {
        <!-- No Merit List -->
        <msls-card class="empty-state-card">
          <div class="empty-state">
            <i class="fa-solid fa-wand-magic-sparkles"></i>
            <h3>No Merit List Generated</h3>
            <p>Click "Generate Merit List" to create a merit list for this class</p>
          </div>
        </msls-card>
      }

      <!-- Decision Modal -->
      <msls-modal
        [isOpen]="showDecisionModal()"
        [title]="decisionModalTitle()"
        size="md"
        (closed)="closeDecisionModal()"
      >
        <div modal-body>
          <msls-decision-form
            [entry]="selectedEntry()"
            [entries]="bulkEntries()"
            [decisionType]="currentDecision()"
            [isBulk]="isBulkDecision()"
            (save)="onDecisionSave($event)"
            (cancel)="closeDecisionModal()"
          />
        </div>
      </msls-modal>

      <!-- Details Modal -->
      <msls-modal
        [isOpen]="showDetailsModal()"
        title="Applicant Details"
        size="lg"
        (closed)="closeDetailsModal()"
      >
        <div modal-body>
          @if (selectedEntry()) {
            <div class="details-content">
              <div class="details-section">
                <h4>Student Information</h4>
                <div class="details-grid">
                  <div class="detail-item">
                    <span class="detail-label">Name</span>
                    <span class="detail-value">{{ selectedEntry()!.studentName }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">Application #</span>
                    <span class="detail-value">{{ selectedEntry()!.applicationNumber }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">Class Applied</span>
                    <span class="detail-value">{{ selectedEntry()!.classApplying }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">Parent/Guardian</span>
                    <span class="detail-value">{{ selectedEntry()!.parentName }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">Contact</span>
                    <span class="detail-value">{{ selectedEntry()!.parentPhone }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">Status</span>
                    <msls-badge [variant]="getStatusConfig(selectedEntry()!.status).variant">
                      {{ getStatusConfig(selectedEntry()!.status).label }}
                    </msls-badge>
                  </div>
                </div>
              </div>

              <div class="details-section">
                <h4>Test Scores</h4>
                <div class="scores-table">
                  <table>
                    <thead>
                      <tr>
                        <th>Subject</th>
                        <th>Score</th>
                        <th>Max</th>
                        <th>%</th>
                      </tr>
                    </thead>
                    <tbody>
                      @for (score of selectedEntry()!.subjectScores; track score.subjectName) {
                        <tr>
                          <td>{{ score.subjectName }}</td>
                          <td>{{ score.score }}</td>
                          <td>{{ score.maxScore }}</td>
                          <td>{{ (score.score / score.maxScore * 100) | number:'1.1-1' }}%</td>
                        </tr>
                      }
                      <tr class="total-row">
                        <td><strong>Total</strong></td>
                        <td><strong>{{ selectedEntry()!.totalScore }}</strong></td>
                        <td><strong>{{ selectedEntry()!.maxScore }}</strong></td>
                        <td><strong>{{ selectedEntry()!.percentage | number:'1.1-1' }}%</strong></td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              @if (selectedEntry()!.sectionAssigned) {
                <div class="details-section">
                  <h4>Admission Details</h4>
                  <div class="details-grid">
                    <div class="detail-item">
                      <span class="detail-label">Section Assigned</span>
                      <span class="detail-value">{{ selectedEntry()!.sectionAssigned }}</span>
                    </div>
                    @if (selectedEntry()!.waitlistPosition) {
                      <div class="detail-item">
                        <span class="detail-label">Waitlist Position</span>
                        <span class="detail-value">#{{ selectedEntry()!.waitlistPosition }}</span>
                      </div>
                    }
                    @if (selectedEntry()!.decisionDate) {
                      <div class="detail-item">
                        <span class="detail-label">Decision Date</span>
                        <span class="detail-value">{{ selectedEntry()!.decisionDate | date:'mediumDate' }}</span>
                      </div>
                    }
                  </div>
                </div>
              }
            </div>
          }
        </div>
        <div modal-footer>
          <msls-button variant="secondary" (click)="closeDetailsModal()">Close</msls-button>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .merit-list-page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .page-header {
      margin-bottom: 1.5rem;
    }

    .page-title {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-size: 1.5rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.25rem 0;
    }

    .page-subtitle {
      color: #64748b;
      margin: 0;
      font-size: 0.875rem;
    }

    .filters-card {
      margin-bottom: 1.5rem;
    }

    .filters-row {
      display: flex;
      flex-wrap: wrap;
      gap: 1rem;
      align-items: flex-end;
    }

    .filter-field {
      flex: 1;
      min-width: 180px;
    }

    .filter-field--small {
      flex: 0 0 120px;
      min-width: 100px;
    }

    .filter-actions {
      display: flex;
      align-items: flex-end;
      padding-bottom: 0.25rem;
    }

    .cutoff-input {
      width: 100%;
      padding: 0.625rem 0.75rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      transition: border-color 0.15s, box-shadow 0.15s;
    }

    .cutoff-input:focus {
      outline: none;
      border-color: #3b82f6;
      box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    .select-all-checkbox {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      cursor: pointer;
    }

    .table-checkbox {
      width: 1rem;
      height: 1rem;
      cursor: pointer;
      accent-color: #3b82f6;
    }

    .stats-row {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .stat-card {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem 1.25rem;
      background: white;
      border-radius: 0.75rem;
      border: 1px solid #e2e8f0;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    }

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .stat-icon--total {
      background: #dbeafe;
      color: #2563eb;
    }

    .stat-icon--selected {
      background: #dcfce7;
      color: #16a34a;
    }

    .stat-icon--waitlist {
      background: #fef3c7;
      color: #d97706;
    }

    .stat-icon--pending {
      background: #e0e7ff;
      color: #4f46e5;
    }

    .stat-content {
      display: flex;
      flex-direction: column;
    }

    .stat-value {
      font-size: 1.5rem;
      font-weight: 700;
      color: #0f172a;
      line-height: 1;
    }

    .stat-label {
      font-size: 0.75rem;
      color: #64748b;
      margin-top: 0.25rem;
    }

    .bulk-actions-bar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0.75rem 1rem;
      background: #f1f5f9;
      border-radius: 0.5rem;
      margin-bottom: 1rem;
    }

    .selected-count {
      font-size: 0.875rem;
      font-weight: 500;
      color: #475569;
    }

    .bulk-buttons {
      display: flex;
      gap: 0.5rem;
    }

    .table-card {
      overflow: hidden;
    }

    .table-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1rem 1.25rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .table-title {
      font-size: 1rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .cutoff-badge {
      font-size: 0.75rem;
      font-weight: 500;
      padding: 0.25rem 0.5rem;
      background: #dbeafe;
      color: #2563eb;
      border-radius: 0.375rem;
    }

    .table-container {
      overflow-x: auto;
    }

    .merit-table {
      width: 100%;
      border-collapse: collapse;
      font-size: 0.875rem;
    }

    .merit-table th {
      text-align: left;
      padding: 0.75rem 1rem;
      font-weight: 600;
      color: #64748b;
      font-size: 0.75rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
      white-space: nowrap;
    }

    .merit-table td {
      padding: 0.875rem 1rem;
      border-bottom: 1px solid #f1f5f9;
      vertical-align: middle;
    }

    .merit-table tbody tr:hover {
      background: #f8fafc;
    }

    .merit-table tbody tr.selected-row {
      background: #eff6ff;
    }

    .col-checkbox {
      width: 40px;
    }

    .col-rank {
      width: 80px;
    }

    .col-name {
      min-width: 200px;
    }

    .col-app {
      width: 140px;
    }

    .col-score {
      width: 100px;
    }

    .col-percent {
      width: 140px;
    }

    .col-status {
      width: 130px;
    }

    .col-actions {
      width: 160px;
    }

    .rank-badge {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 0.25rem;
      padding: 0.375rem 0.625rem;
      border-radius: 0.5rem;
      font-weight: 600;
      font-size: 0.8125rem;
      min-width: 48px;
    }

    .rank-badge.rank-1 {
      background: linear-gradient(135deg, #fef3c7, #fde68a);
      color: #92400e;
    }

    .rank-badge.rank-2 {
      background: linear-gradient(135deg, #e2e8f0, #cbd5e1);
      color: #475569;
    }

    .rank-badge.rank-3 {
      background: linear-gradient(135deg, #fed7aa, #fdba74);
      color: #9a3412;
    }

    .rank-badge.rank-other {
      background: #f1f5f9;
      color: #64748b;
    }

    .student-info {
      display: flex;
      flex-direction: column;
    }

    .student-name {
      font-weight: 500;
      color: #0f172a;
    }

    .parent-info {
      font-size: 0.75rem;
      color: #64748b;
      margin-top: 0.125rem;
    }

    .app-number {
      font-family: 'Monaco', 'Menlo', monospace;
      font-size: 0.8125rem;
      color: #64748b;
    }

    .score-display {
      display: flex;
      align-items: baseline;
      gap: 0.125rem;
    }

    .score-value {
      font-weight: 600;
      color: #0f172a;
    }

    .score-max {
      font-size: 0.75rem;
      color: #94a3b8;
    }

    .percentage-bar {
      position: relative;
      height: 24px;
      background: #f1f5f9;
      border-radius: 0.375rem;
      overflow: hidden;
      min-width: 100px;
    }

    .percentage-fill {
      position: absolute;
      top: 0;
      left: 0;
      height: 100%;
      background: linear-gradient(90deg, #3b82f6, #60a5fa);
      border-radius: 0.375rem;
      transition: width 0.3s ease;
    }

    .percentage-text {
      position: relative;
      z-index: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      font-size: 0.75rem;
      font-weight: 600;
      color: #1e3a8a;
    }

    .action-buttons {
      display: flex;
      gap: 0.375rem;
    }

    .action-btn {
      width: 32px;
      height: 32px;
      border: none;
      border-radius: 0.375rem;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: all 0.15s;
      font-size: 0.875rem;
    }

    .action-btn--approve {
      background: #dcfce7;
      color: #16a34a;
    }

    .action-btn--approve:hover {
      background: #bbf7d0;
    }

    .action-btn--waitlist {
      background: #fef3c7;
      color: #d97706;
    }

    .action-btn--waitlist:hover {
      background: #fde68a;
    }

    .action-btn--reject {
      background: #fee2e2;
      color: #dc2626;
    }

    .action-btn--reject:hover {
      background: #fecaca;
    }

    .action-btn--letter {
      background: #dbeafe;
      color: #2563eb;
    }

    .action-btn--letter:hover {
      background: #bfdbfe;
    }

    .action-btn--enroll {
      background: #d1fae5;
      color: #059669;
    }

    .action-btn--enroll:hover {
      background: #a7f3d0;
    }

    .action-btn--view {
      background: #f1f5f9;
      color: #64748b;
    }

    .action-btn--view:hover {
      background: #e2e8f0;
    }

    .empty-state-card {
      padding: 3rem;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;
      padding: 2rem;
    }

    .empty-state i {
      font-size: 3rem;
      color: #cbd5e1;
      margin-bottom: 1rem;
    }

    .empty-state h3 {
      margin: 0 0 0.5rem 0;
      font-size: 1.125rem;
      color: #334155;
    }

    .empty-state p {
      margin: 0;
      color: #64748b;
      font-size: 0.875rem;
    }

    .loading-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 3rem;
    }

    .spinner {
      width: 40px;
      height: 40px;
      border: 3px solid #e2e8f0;
      border-top-color: #3b82f6;
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin-bottom: 1rem;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .details-content {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    .details-section h4 {
      font-size: 0.875rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      margin: 0 0 1rem 0;
      padding-bottom: 0.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .details-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
    }

    .detail-item {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .detail-label {
      font-size: 0.75rem;
      color: #64748b;
    }

    .detail-value {
      font-size: 0.875rem;
      font-weight: 500;
      color: #0f172a;
    }

    .scores-table table {
      width: 100%;
      border-collapse: collapse;
      font-size: 0.875rem;
    }

    .scores-table th,
    .scores-table td {
      padding: 0.625rem 1rem;
      text-align: left;
      border-bottom: 1px solid #e2e8f0;
    }

    .scores-table th {
      font-weight: 600;
      color: #64748b;
      background: #f8fafc;
      font-size: 0.75rem;
      text-transform: uppercase;
    }

    .scores-table .total-row {
      background: #f1f5f9;
    }

    @media (max-width: 768px) {
      .merit-list-page {
        padding: 1rem;
      }

      .filters-row {
        flex-direction: column;
      }

      .filter-field,
      .filter-field--small {
        flex: none;
        width: 100%;
        min-width: auto;
      }

      .stats-row {
        grid-template-columns: repeat(2, 1fr);
      }

      .bulk-actions-bar {
        flex-direction: column;
        gap: 0.75rem;
        align-items: stretch;
      }

      .bulk-buttons {
        flex-wrap: wrap;
        justify-content: center;
      }

      .details-grid {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class MeritListComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly meritService = inject(MeritService);
  private readonly sessionService = inject(AdmissionSessionService);
  private readonly toastService = inject(ToastService);

  // State signals
  loading = signal(false);
  generating = signal(false);
  meritList = signal<MeritList | null>(null);

  // Filter state
  selectedSession = signal<string | null>(null);
  selectedClass = signal<string | null>(null);
  cutoffScore = signal<number>(0);

  // Selection state
  selectedEntries = signal<string[]>([]);

  // Modal state
  showDecisionModal = signal(false);
  showDetailsModal = signal(false);
  selectedEntry = signal<MeritListEntry | null>(null);
  bulkEntries = signal<MeritListEntry[]>([]);
  currentDecision = signal<DecisionType>('selected');
  isBulkDecision = signal(false);

  // Options
  sessionOptions = signal<SelectOption[]>([]);
  classOptions = signal<SelectOption[]>([]);

  // Computed values
  canGenerate = computed(() =>
    this.selectedSession() !== null &&
    this.selectedClass() !== null &&
    !this.generating()
  );

  selectedCount = computed(() =>
    this.meritList()?.entries.filter(e =>
      e.status === 'selected' || e.status === 'offer_sent' || e.status === 'offer_accepted' || e.status === 'enrolled'
    ).length || 0
  );

  waitlistedCount = computed(() =>
    this.meritList()?.entries.filter(e => e.status === 'waitlisted').length || 0
  );

  pendingCount = computed(() =>
    this.meritList()?.entries.filter(e => e.status === 'test_completed').length || 0
  );

  allSelected = computed(() => {
    const entries = this.meritList()?.entries || [];
    return entries.length > 0 && entries.every(e => this.selectedEntries().includes(e.applicationId));
  });

  someSelected = computed(() => {
    const entries = this.meritList()?.entries || [];
    const selected = this.selectedEntries();
    return entries.some(e => selected.includes(e.applicationId)) && !this.allSelected();
  });

  decisionModalTitle = computed(() => {
    if (this.isBulkDecision()) {
      const count = this.bulkEntries().length;
      const decision = this.currentDecision();
      const actionMap = { selected: 'Approve', waitlisted: 'Waitlist', rejected: 'Reject' };
      return `${actionMap[decision]} ${count} Candidates`;
    }
    const entry = this.selectedEntry();
    const decision = this.currentDecision();
    const actionMap = { selected: 'Approve', waitlisted: 'Waitlist', rejected: 'Reject' };
    return `${actionMap[decision]} - ${entry?.studentName || ''}`;
  });

  ngOnInit(): void {
    this.loadSessions();
  }

  private loadSessions(): void {
    this.sessionService.getSessions().subscribe({
      next: (sessions) => {
        this.sessionOptions.set(
          sessions.map(s => ({
            value: s.id,
            label: `${s.name} (${s.status})`,
          }))
        );
      },
      error: (err) => {
        console.error('Failed to load sessions:', err);
        this.toastService.error('Failed to load admission sessions');
      },
    });
  }

  onSessionChange(value: string | number | (string | number)[] | null): void {
    const sessionId = value as string;
    this.selectedSession.set(sessionId);
    this.selectedClass.set(null);
    this.meritList.set(null);
    this.selectedEntries.set([]);

    if (sessionId) {
      this.loadClasses(sessionId);
    } else {
      this.classOptions.set([]);
    }
  }

  private loadClasses(sessionId: string): void {
    this.meritService.getAvailableClasses(sessionId).subscribe({
      next: (classes) => {
        this.classOptions.set(
          classes.map(c => ({
            value: c,
            label: c,
          }))
        );
      },
      error: (err) => {
        console.error('Failed to load classes:', err);
        this.toastService.error('Failed to load available classes');
      },
    });
  }

  onClassChange(value: string | number | (string | number)[] | null): void {
    const className = value as string;
    this.selectedClass.set(className);
    this.meritList.set(null);
    this.selectedEntries.set([]);

    if (className && this.selectedSession()) {
      this.loadMeritList();
    }
  }

  onCutoffChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    const value = parseFloat(input.value) || 0;
    this.cutoffScore.set(Math.max(0, Math.min(100, value)));
  }

  private loadMeritList(): void {
    const sessionId = this.selectedSession();
    const className = this.selectedClass();
    if (!sessionId || !className) return;

    this.loading.set(true);
    this.meritService.getMeritList(sessionId, className).subscribe({
      next: (list) => {
        this.meritList.set(list);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load merit list:', err);
        this.loading.set(false);
      },
    });
  }

  generateMeritList(): void {
    const sessionId = this.selectedSession();
    const className = this.selectedClass();
    if (!sessionId || !className) return;

    this.generating.set(true);
    this.meritService.generateMeritList(sessionId, {
      className,
      cutoffScore: this.cutoffScore() > 0 ? this.cutoffScore() : undefined,
    }).subscribe({
      next: (list) => {
        this.meritList.set(list);
        this.generating.set(false);
        this.toastService.success(`Merit list generated with ${list.entries.length} candidates`);
      },
      error: (err: Error) => {
        console.error('Failed to generate merit list:', err);
        this.generating.set(false);
        this.toastService.error(err.message || 'Failed to generate merit list');
      },
    });
  }

  getStatusConfig(status: string) {
    return getStatusConfig(status as any) || { label: status, variant: 'neutral' as const, icon: 'fa-solid fa-circle' };
  }

  getRankClass(rank: number): string {
    if (rank === 1) return 'rank-1';
    if (rank === 2) return 'rank-2';
    if (rank === 3) return 'rank-3';
    return 'rank-other';
  }

  canMakeDecision(entry: MeritListEntry): boolean {
    return entry.status === 'test_completed';
  }

  // Selection methods
  isEntrySelected(applicationId: string): boolean {
    return this.selectedEntries().includes(applicationId);
  }

  toggleEntrySelection(applicationId: string, selected: boolean): void {
    if (selected) {
      this.selectedEntries.update(entries => [...entries, applicationId]);
    } else {
      this.selectedEntries.update(entries => entries.filter(id => id !== applicationId));
    }
  }

  toggleSelectAll(selected: boolean): void {
    if (selected) {
      const allIds = this.meritList()?.entries.map(e => e.applicationId) || [];
      this.selectedEntries.set(allIds);
    } else {
      this.selectedEntries.set([]);
    }
  }

  clearSelection(): void {
    this.selectedEntries.set([]);
  }

  // Modal methods
  openDecisionModal(entry: MeritListEntry, decision: DecisionType): void {
    this.selectedEntry.set(entry);
    this.bulkEntries.set([]);
    this.currentDecision.set(decision);
    this.isBulkDecision.set(false);
    this.showDecisionModal.set(true);
  }

  openBulkDecisionModal(decision: DecisionType): void {
    const selectedIds = this.selectedEntries();
    const entries = this.meritList()?.entries.filter(e => selectedIds.includes(e.applicationId)) || [];
    this.selectedEntry.set(null);
    this.bulkEntries.set(entries);
    this.currentDecision.set(decision);
    this.isBulkDecision.set(true);
    this.showDecisionModal.set(true);
  }

  closeDecisionModal(): void {
    this.showDecisionModal.set(false);
    this.selectedEntry.set(null);
    this.bulkEntries.set([]);
  }

  onDecisionSave(result: { success: boolean }): void {
    if (result.success) {
      this.closeDecisionModal();
      this.loadMeritList();
      this.clearSelection();
    }
  }

  viewDetails(entry: MeritListEntry): void {
    this.selectedEntry.set(entry);
    this.showDetailsModal.set(true);
  }

  closeDetailsModal(): void {
    this.showDetailsModal.set(false);
    this.selectedEntry.set(null);
  }

  viewOfferLetter(entry: MeritListEntry): void {
    this.meritService.generateOfferLetter(entry.applicationId).subscribe({
      next: (result) => {
        this.toastService.info('Offer letter available (placeholder)');
        // In production: window.open(result.url, '_blank');
      },
      error: () => {
        this.toastService.error('Failed to generate offer letter');
      },
    });
  }

  goToEnrollment(entry: MeritListEntry): void {
    this.router.navigate(['/admissions/enrollment', entry.applicationId]);
  }
}
