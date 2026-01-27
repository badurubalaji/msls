/**
 * MSLS Application Form Component
 *
 * Multi-section form for creating/editing admission applications.
 * Includes sections for student details, parent/guardian details,
 * previous school information, and document uploads.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormArray, Validators, ReactiveFormsModule, FormsModule } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';

// Modal and Badge components available if needed for future enhancements
import { ToastService } from '../../../shared/services';
import {
  AdmissionApplication,
  CreateApplicationRequest,
  ParentRequest,
  Gender,
  BloodGroup,
  ParentRelation,
  DocumentType,
  ApplicationDocument,
  GENDER_LABELS,
  BLOOD_GROUP_OPTIONS,
  PARENT_RELATION_LABELS,
  DOCUMENT_TYPE_LABELS,
  CATEGORY_OPTIONS,
} from './application.model';
import { ApplicationService } from './application.service';
import { CLASS_NAMES, AdmissionSession } from '../sessions/admission-session.model';

type FormSection = 'student' | 'parents' | 'previousSchool' | 'documents';

@Component({
  selector: 'msls-application-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
  template: `
    <div class="form-page">
      <div class="form-card">
        <!-- Header -->
        <div class="form-header">
          <div class="form-header__left">
            <button class="back-btn" (click)="goBack()">
              <i class="fa-solid fa-arrow-left"></i>
            </button>
            <div>
              <h1 class="form-header__title">
                {{ isEditMode() ? 'Edit Application' : 'New Admission Application' }}
              </h1>
              <p class="form-header__subtitle">
                {{ isEditMode() ? 'Update application details' : 'Fill in the student and parent details' }}
              </p>
            </div>
          </div>
          @if (application()?.applicationNumber) {
            <div class="app-number">
              <span class="app-number__label">Application #</span>
              <span class="app-number__value">{{ application()?.applicationNumber }}</span>
            </div>
          }
        </div>

        <!-- Section Navigation -->
        <div class="section-nav">
          @for (section of sections; track section.key; let i = $index) {
            <button
              class="section-tab"
              [class.section-tab--active]="currentSection() === section.key"
              [class.section-tab--completed]="isSectionCompleted(section.key)"
              (click)="setSection(section.key)"
            >
              <span class="section-tab__number">{{ i + 1 }}</span>
              <span class="section-tab__label">{{ section.label }}</span>
              @if (isSectionCompleted(section.key)) {
                <i class="fa-solid fa-check section-tab__check"></i>
              }
            </button>
          }
        </div>

        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <p>Loading application...</p>
          </div>
        } @else {
          <form [formGroup]="applicationForm" (ngSubmit)="onSubmit()">
            <!-- Student Details Section -->
            @if (currentSection() === 'student') {
              <div class="form-section">
                <h2 class="section-title">
                  <i class="fa-solid fa-user-graduate"></i>
                  Student Details
                </h2>

                <div class="form-grid">
                  <!-- Session -->
                  <div class="form-group form-group--full">
                    <label class="form-label required">Admission Session</label>
                    <select class="form-select" formControlName="sessionId">
                      <option value="">Select Session</option>
                      @for (session of sessions(); track session.id) {
                        <option [value]="session.id">{{ session.name }}</option>
                      }
                    </select>
                    @if (showError('sessionId')) {
                      <span class="form-error">Session is required</span>
                    }
                  </div>

                  <!-- Class Applying -->
                  <div class="form-group">
                    <label class="form-label required">Class Applying For</label>
                    <select class="form-select" formControlName="classApplying">
                      <option value="">Select Class</option>
                      @for (className of classOptions; track className) {
                        <option [value]="className">{{ className }}</option>
                      }
                    </select>
                    @if (showError('classApplying')) {
                      <span class="form-error">Class is required</span>
                    }
                  </div>

                  <!-- First Name -->
                  <div class="form-group">
                    <label class="form-label required">First Name</label>
                    <input type="text" class="form-input" formControlName="firstName" placeholder="Enter first name" />
                    @if (showError('firstName')) {
                      <span class="form-error">First name is required</span>
                    }
                  </div>

                  <!-- Middle Name -->
                  <div class="form-group">
                    <label class="form-label">Middle Name</label>
                    <input
                      type="text"
                      class="form-input"
                      formControlName="middleName"
                      placeholder="Enter middle name"
                    />
                  </div>

                  <!-- Last Name -->
                  <div class="form-group">
                    <label class="form-label required">Last Name</label>
                    <input type="text" class="form-input" formControlName="lastName" placeholder="Enter last name" />
                    @if (showError('lastName')) {
                      <span class="form-error">Last name is required</span>
                    }
                  </div>

                  <!-- Date of Birth -->
                  <div class="form-group">
                    <label class="form-label required">Date of Birth</label>
                    <input type="date" class="form-input" formControlName="dateOfBirth" />
                    @if (showError('dateOfBirth')) {
                      <span class="form-error">Date of birth is required</span>
                    }
                  </div>

                  <!-- Gender -->
                  <div class="form-group">
                    <label class="form-label required">Gender</label>
                    <select class="form-select" formControlName="gender">
                      <option value="">Select Gender</option>
                      @for (gender of genderOptions; track gender.value) {
                        <option [value]="gender.value">{{ gender.label }}</option>
                      }
                    </select>
                    @if (showError('gender')) {
                      <span class="form-error">Gender is required</span>
                    }
                  </div>

                  <!-- Blood Group -->
                  <div class="form-group">
                    <label class="form-label">Blood Group</label>
                    <select class="form-select" formControlName="bloodGroup">
                      <option value="">Select Blood Group</option>
                      @for (bg of bloodGroupOptions; track bg) {
                        <option [value]="bg">{{ bg }}</option>
                      }
                    </select>
                  </div>

                  <!-- Nationality -->
                  <div class="form-group">
                    <label class="form-label">Nationality</label>
                    <input type="text" class="form-input" formControlName="nationality" placeholder="e.g., Indian" />
                  </div>

                  <!-- Religion -->
                  <div class="form-group">
                    <label class="form-label">Religion</label>
                    <input type="text" class="form-input" formControlName="religion" placeholder="Enter religion" />
                  </div>

                  <!-- Category -->
                  <div class="form-group">
                    <label class="form-label">Category</label>
                    <select class="form-select" formControlName="category">
                      <option value="">Select Category</option>
                      @for (cat of categoryOptions; track cat.value) {
                        <option [value]="cat.value">{{ cat.label }}</option>
                      }
                    </select>
                  </div>

                  <!-- Aadhaar Number -->
                  <div class="form-group">
                    <label class="form-label">Aadhaar Number</label>
                    <input
                      type="text"
                      class="form-input"
                      formControlName="aadhaarNumber"
                      placeholder="12-digit Aadhaar"
                      maxlength="12"
                    />
                  </div>
                </div>

                <!-- Address Section -->
                <h3 class="subsection-title">
                  <i class="fa-solid fa-location-dot"></i>
                  Address
                </h3>
                <div class="form-grid">
                  <div class="form-group form-group--full">
                    <label class="form-label">Address Line 1</label>
                    <input
                      type="text"
                      class="form-input"
                      formControlName="addressLine1"
                      placeholder="House/Flat No., Street"
                    />
                  </div>
                  <div class="form-group form-group--full">
                    <label class="form-label">Address Line 2</label>
                    <input
                      type="text"
                      class="form-input"
                      formControlName="addressLine2"
                      placeholder="Area, Landmark"
                    />
                  </div>
                  <div class="form-group">
                    <label class="form-label">City</label>
                    <input type="text" class="form-input" formControlName="city" placeholder="Enter city" />
                  </div>
                  <div class="form-group">
                    <label class="form-label">State</label>
                    <input type="text" class="form-input" formControlName="state" placeholder="Enter state" />
                  </div>
                  <div class="form-group">
                    <label class="form-label">Postal Code</label>
                    <input
                      type="text"
                      class="form-input"
                      formControlName="postalCode"
                      placeholder="PIN Code"
                      maxlength="6"
                    />
                  </div>
                </div>
              </div>
            }

            <!-- Parents Section -->
            @if (currentSection() === 'parents') {
              <div class="form-section">
                <div class="section-header">
                  <h2 class="section-title">
                    <i class="fa-solid fa-users"></i>
                    Parent/Guardian Details
                  </h2>
                  <button type="button" class="btn btn-secondary btn-sm" (click)="addParent()">
                    <i class="fa-solid fa-plus"></i>
                    Add Parent/Guardian
                  </button>
                </div>

                @if (parentsArray.length === 0) {
                  <div class="empty-parents">
                    <i class="fa-solid fa-user-plus"></i>
                    <p>No parent/guardian added yet</p>
                    <button type="button" class="btn btn-primary btn-sm" (click)="addParent()">
                      <i class="fa-solid fa-plus"></i>
                      Add First Parent/Guardian
                    </button>
                  </div>
                }

                <div formArrayName="parents">
                  @for (parent of parentsArray.controls; track i; let i = $index) {
                    <div class="parent-card" [formGroupName]="i">
                      <div class="parent-card__header">
                        <span class="parent-card__title">{{ getParentTitle(i) }}</span>
                        <button
                          type="button"
                          class="parent-card__remove"
                          (click)="removeParent(i)"
                          title="Remove"
                        >
                          <i class="fa-solid fa-times"></i>
                        </button>
                      </div>
                      <div class="form-grid">
                        <div class="form-group">
                          <label class="form-label required">Relation</label>
                          <select class="form-select" formControlName="relation">
                            <option value="">Select Relation</option>
                            @for (rel of relationOptions; track rel.value) {
                              <option [value]="rel.value">{{ rel.label }}</option>
                            }
                          </select>
                        </div>
                        <div class="form-group">
                          <label class="form-label required">Full Name</label>
                          <input type="text" class="form-input" formControlName="name" placeholder="Enter full name" />
                        </div>
                        <div class="form-group">
                          <label class="form-label">Phone Number</label>
                          <input
                            type="tel"
                            class="form-input"
                            formControlName="phone"
                            placeholder="10-digit mobile"
                          />
                        </div>
                        <div class="form-group">
                          <label class="form-label">Email</label>
                          <input type="email" class="form-input" formControlName="email" placeholder="Email address" />
                        </div>
                        <div class="form-group">
                          <label class="form-label">Occupation</label>
                          <input
                            type="text"
                            class="form-input"
                            formControlName="occupation"
                            placeholder="e.g., Engineer"
                          />
                        </div>
                        <div class="form-group">
                          <label class="form-label">Education</label>
                          <input
                            type="text"
                            class="form-input"
                            formControlName="education"
                            placeholder="e.g., Graduate"
                          />
                        </div>
                        <div class="form-group">
                          <label class="form-label">Annual Income</label>
                          <input
                            type="text"
                            class="form-input"
                            formControlName="annualIncome"
                            placeholder="e.g., 5-10 Lakhs"
                          />
                        </div>
                      </div>
                    </div>
                  }
                </div>
              </div>
            }

            <!-- Previous School Section -->
            @if (currentSection() === 'previousSchool') {
              <div class="form-section">
                <h2 class="section-title">
                  <i class="fa-solid fa-school"></i>
                  Previous School Details
                </h2>
                <p class="section-hint">
                  Fill this section if the student is transferring from another school. Skip if applying for entry-level
                  classes.
                </p>

                <div class="form-grid">
                  <div class="form-group form-group--full">
                    <label class="form-label">Previous School Name</label>
                    <input
                      type="text"
                      class="form-input"
                      formControlName="previousSchool"
                      placeholder="Name of the previous school"
                    />
                  </div>
                  <div class="form-group">
                    <label class="form-label">Previous Class</label>
                    <input
                      type="text"
                      class="form-input"
                      formControlName="previousClass"
                      placeholder="e.g., Class 5"
                    />
                  </div>
                  <div class="form-group">
                    <label class="form-label">Percentage/Grade</label>
                    <input
                      type="number"
                      class="form-input"
                      formControlName="previousPercentage"
                      placeholder="e.g., 85.5"
                      min="0"
                      max="100"
                      step="0.1"
                    />
                  </div>
                </div>
              </div>
            }

            <!-- Documents Section -->
            @if (currentSection() === 'documents') {
              <div class="form-section">
                <h2 class="section-title">
                  <i class="fa-solid fa-file-arrow-up"></i>
                  Documents Upload
                </h2>
                <p class="section-hint">
                  Upload required documents. Supported formats: PDF, JPG, PNG (Max 5MB each)
                </p>

                <!-- Upload Area -->
                <div class="upload-section">
                  <div class="upload-controls">
                    <select class="form-select upload-select" [(ngModel)]="selectedDocType" [ngModelOptions]="{standalone: true}">
                      <option value="">Select Document Type</option>
                      @for (docType of documentTypeOptions; track docType.value) {
                        <option [value]="docType.value">{{ docType.label }}</option>
                      }
                    </select>
                    <label class="upload-btn">
                      <i class="fa-solid fa-cloud-arrow-up"></i>
                      Choose File
                      <input
                        type="file"
                        accept=".pdf,.jpg,.jpeg,.png"
                        (change)="onFileSelected($event)"
                        #fileInput
                        hidden
                      />
                    </label>
                  </div>

                  @if (uploadingDocument()) {
                    <div class="upload-progress">
                      <div class="spinner-small"></div>
                      <span>Uploading {{ uploadingFileName() }}...</span>
                    </div>
                  }
                </div>

                <!-- Uploaded Documents -->
                <div class="documents-list">
                  <h3 class="documents-list__title">Uploaded Documents</h3>
                  @if (documents().length === 0) {
                    <div class="no-documents">
                      <i class="fa-regular fa-folder-open"></i>
                      <p>No documents uploaded yet</p>
                    </div>
                  } @else {
                    @for (doc of documents(); track doc.id) {
                      <div class="document-item">
                        <div class="document-item__icon">
                          <i [class]="getDocumentIcon(doc.fileName)"></i>
                        </div>
                        <div class="document-item__info">
                          <span class="document-item__name">{{ doc.fileName }}</span>
                          <span class="document-item__type">{{ getDocumentTypeLabel(doc.documentType) }}</span>
                        </div>
                        <div class="document-item__status">
                          @if (doc.isVerified) {
                            <span class="doc-status doc-status--verified">
                              <i class="fa-solid fa-circle-check"></i>
                              Verified
                            </span>
                          } @else {
                            <span class="doc-status doc-status--pending">
                              <i class="fa-solid fa-clock"></i>
                              Pending
                            </span>
                          }
                        </div>
                        <div class="document-item__actions">
                          <a [href]="doc.fileUrl" target="_blank" class="action-btn" title="View">
                            <i class="fa-regular fa-eye"></i>
                          </a>
                          <button
                            type="button"
                            class="action-btn action-btn--danger"
                            (click)="deleteDocument(doc)"
                            title="Delete"
                          >
                            <i class="fa-regular fa-trash-can"></i>
                          </button>
                        </div>
                      </div>
                    }
                  }
                </div>
              </div>
            }

            <!-- Form Actions -->
            <div class="form-actions">
              <div class="form-actions__left">
                @if (currentSectionIndex() > 0) {
                  <button type="button" class="btn btn-secondary" (click)="prevSection()">
                    <i class="fa-solid fa-arrow-left"></i>
                    Previous
                  </button>
                }
              </div>
              <div class="form-actions__right">
                @if (currentSectionIndex() < sections.length - 1) {
                  <button type="button" class="btn btn-primary" (click)="nextSection()">
                    Next
                    <i class="fa-solid fa-arrow-right"></i>
                  </button>
                } @else {
                  <button type="button" class="btn btn-secondary" (click)="saveDraft()" [disabled]="saving()">
                    @if (saving() && !submitting()) {
                      <div class="btn-spinner"></div>
                    }
                    Save as Draft
                  </button>
                  <button type="submit" class="btn btn-primary" [disabled]="submitting()">
                    @if (submitting()) {
                      <div class="btn-spinner"></div>
                      Submitting...
                    } @else {
                      <i class="fa-solid fa-paper-plane"></i>
                      Submit Application
                    }
                  </button>
                }
              </div>
            </div>
          </form>
        }
      </div>
    </div>
  `,
  styles: [
    `
      .form-page {
        padding: 1.5rem;
        max-width: 900px;
        margin: 0 auto;
      }

      .form-card {
        background: #ffffff;
        border: 1px solid #e2e8f0;
        border-radius: 1rem;
        padding: 1.5rem;
      }

      /* Header */
      .form-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1.5rem;
        padding-bottom: 1.5rem;
        border-bottom: 1px solid #e2e8f0;
        gap: 1rem;
        flex-wrap: wrap;
      }

      .form-header__left {
        display: flex;
        align-items: flex-start;
        gap: 1rem;
      }

      .back-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2.5rem;
        height: 2.5rem;
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        color: #64748b;
        cursor: pointer;
        transition: all 0.15s;
      }

      .back-btn:hover {
        background: #f1f5f9;
        color: #0f172a;
      }

      .form-header__title {
        font-size: 1.5rem;
        font-weight: 700;
        color: #0f172a;
        margin: 0 0 0.25rem 0;
      }

      .form-header__subtitle {
        font-size: 0.875rem;
        color: #64748b;
        margin: 0;
      }

      .app-number {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: 0.25rem;
      }

      .app-number__label {
        font-size: 0.75rem;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }

      .app-number__value {
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 1rem;
        font-weight: 600;
        color: #4f46e5;
      }

      /* Section Navigation */
      .section-nav {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1.5rem;
        padding: 0.5rem;
        background: #f8fafc;
        border-radius: 0.75rem;
        overflow-x: auto;
      }

      .section-tab {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: transparent;
        border: none;
        border-radius: 0.5rem;
        color: #64748b;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.15s;
        white-space: nowrap;
      }

      .section-tab:hover {
        background: #ffffff;
        color: #334155;
      }

      .section-tab--active {
        background: #ffffff;
        color: #4f46e5;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
      }

      .section-tab__number {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 1.5rem;
        height: 1.5rem;
        background: #e2e8f0;
        border-radius: 50%;
        font-size: 0.75rem;
        font-weight: 600;
      }

      .section-tab--active .section-tab__number {
        background: #4f46e5;
        color: #ffffff;
      }

      .section-tab--completed .section-tab__number {
        background: #10b981;
        color: #ffffff;
      }

      .section-tab__check {
        color: #10b981;
        font-size: 0.75rem;
      }

      /* Form Section */
      .form-section {
        margin-bottom: 1.5rem;
      }

      .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
        flex-wrap: wrap;
        gap: 1rem;
      }

      .section-title {
        display: flex;
        align-items: center;
        gap: 0.625rem;
        font-size: 1.125rem;
        font-weight: 600;
        color: #0f172a;
        margin: 0 0 1.5rem 0;
      }

      .section-header .section-title {
        margin: 0;
      }

      .section-title i {
        color: #4f46e5;
      }

      .section-hint {
        font-size: 0.875rem;
        color: #64748b;
        margin: -1rem 0 1.5rem 0;
        padding: 0.75rem;
        background: #f8fafc;
        border-radius: 0.5rem;
        border-left: 3px solid #4f46e5;
      }

      .subsection-title {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9375rem;
        font-weight: 600;
        color: #334155;
        margin: 2rem 0 1rem 0;
        padding-top: 1.5rem;
        border-top: 1px solid #e2e8f0;
      }

      .subsection-title i {
        color: #64748b;
      }

      /* Form Grid */
      .form-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 1rem;
      }

      .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
      }

      .form-group--full {
        grid-column: 1 / -1;
      }

      .form-label {
        font-size: 0.8125rem;
        font-weight: 500;
        color: #374151;
      }

      .form-label.required::after {
        content: ' *';
        color: #dc2626;
      }

      .form-input,
      .form-select {
        padding: 0.625rem 0.875rem;
        font-size: 0.875rem;
        border: 1px solid #e2e8f0;
        border-radius: 0.5rem;
        background: #ffffff;
        color: #0f172a;
        transition: all 0.15s;
      }

      .form-select {
        padding-right: 2.5rem;
        background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3E%3Cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3E%3C/svg%3E");
        background-position: right 0.5rem center;
        background-size: 1.5rem;
        background-repeat: no-repeat;
        appearance: none;
        cursor: pointer;
      }

      .form-input::placeholder {
        color: #9ca3af;
      }

      .form-input:focus,
      .form-select:focus {
        outline: none;
        border-color: #4f46e5;
        box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
      }

      .form-error {
        font-size: 0.75rem;
        color: #dc2626;
      }

      /* Parent Cards */
      .empty-parents {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        padding: 3rem;
        background: #f8fafc;
        border: 2px dashed #e2e8f0;
        border-radius: 0.75rem;
        color: #94a3b8;
      }

      .empty-parents i {
        font-size: 2.5rem;
      }

      .empty-parents p {
        margin: 0;
        font-size: 0.875rem;
      }

      .parent-card {
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        border-radius: 0.75rem;
        padding: 1rem;
        margin-bottom: 1rem;
      }

      .parent-card__header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
        padding-bottom: 0.75rem;
        border-bottom: 1px solid #e2e8f0;
      }

      .parent-card__title {
        font-size: 0.9375rem;
        font-weight: 600;
        color: #334155;
      }

      .parent-card__remove {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 1.75rem;
        height: 1.75rem;
        background: transparent;
        border: 1px solid #fecaca;
        border-radius: 0.375rem;
        color: #dc2626;
        cursor: pointer;
        transition: all 0.15s;
      }

      .parent-card__remove:hover {
        background: #fef2f2;
      }

      .parent-card .form-input,
      .parent-card .form-select {
        background: #ffffff;
      }

      /* Documents Section */
      .upload-section {
        margin-bottom: 1.5rem;
      }

      .upload-controls {
        display: flex;
        gap: 0.75rem;
        align-items: center;
        flex-wrap: wrap;
      }

      .upload-select {
        flex: 1;
        min-width: 200px;
      }

      .upload-btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.625rem 1rem;
        background: #4f46e5;
        color: #ffffff;
        border-radius: 0.5rem;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.15s;
      }

      .upload-btn:hover {
        background: #4338ca;
      }

      .upload-progress {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-top: 0.75rem;
        padding: 0.75rem;
        background: #eff6ff;
        border-radius: 0.5rem;
        color: #1e40af;
        font-size: 0.875rem;
      }

      .spinner-small {
        width: 1rem;
        height: 1rem;
        border: 2px solid rgba(30, 64, 175, 0.3);
        border-top-color: #1e40af;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
      }

      .documents-list {
        border: 1px solid #e2e8f0;
        border-radius: 0.75rem;
        overflow: hidden;
      }

      .documents-list__title {
        font-size: 0.8125rem;
        font-weight: 600;
        color: #64748b;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        padding: 0.75rem 1rem;
        background: #f8fafc;
        border-bottom: 1px solid #e2e8f0;
        margin: 0;
      }

      .no-documents {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
        padding: 2rem;
        color: #94a3b8;
      }

      .no-documents i {
        font-size: 2rem;
      }

      .no-documents p {
        margin: 0;
        font-size: 0.875rem;
      }

      .document-item {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.875rem 1rem;
        border-bottom: 1px solid #f1f5f9;
      }

      .document-item:last-child {
        border-bottom: none;
      }

      .document-item__icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2.5rem;
        height: 2.5rem;
        background: #f1f5f9;
        border-radius: 0.5rem;
        color: #64748b;
      }

      .document-item__info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
        min-width: 0;
      }

      .document-item__name {
        font-size: 0.875rem;
        font-weight: 500;
        color: #0f172a;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .document-item__type {
        font-size: 0.75rem;
        color: #64748b;
      }

      .document-item__status {
        flex-shrink: 0;
      }

      .doc-status {
        display: inline-flex;
        align-items: center;
        gap: 0.375rem;
        padding: 0.25rem 0.5rem;
        font-size: 0.75rem;
        font-weight: 500;
        border-radius: 9999px;
      }

      .doc-status--verified {
        background: #dcfce7;
        color: #166534;
      }

      .doc-status--pending {
        background: #fef3c7;
        color: #92400e;
      }

      .document-item__actions {
        display: flex;
        gap: 0.375rem;
        flex-shrink: 0;
      }

      /* Buttons */
      .btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 0.625rem 1rem;
        font-size: 0.875rem;
        font-weight: 500;
        border-radius: 0.5rem;
        border: none;
        cursor: pointer;
        transition: all 0.15s;
      }

      .btn-sm {
        padding: 0.5rem 0.875rem;
        font-size: 0.8125rem;
      }

      .btn-primary {
        background: #4f46e5;
        color: #ffffff;
      }

      .btn-primary:hover:not(:disabled) {
        background: #4338ca;
      }

      .btn-primary:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      .btn-secondary {
        background: #ffffff;
        color: #334155;
        border: 1px solid #e2e8f0;
      }

      .btn-secondary:hover:not(:disabled) {
        background: #f8fafc;
        border-color: #cbd5e1;
      }

      .btn-secondary:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      .btn-spinner {
        width: 1rem;
        height: 1rem;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-top-color: #ffffff;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
      }

      .btn-secondary .btn-spinner {
        border-color: rgba(0, 0, 0, 0.1);
        border-top-color: #334155;
      }

      .action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 1.75rem;
        height: 1.75rem;
        background: transparent;
        border: 1px solid #e2e8f0;
        border-radius: 0.375rem;
        color: #64748b;
        cursor: pointer;
        transition: all 0.15s;
        text-decoration: none;
      }

      .action-btn:hover {
        background: #f8fafc;
        border-color: #cbd5e1;
        color: #0f172a;
      }

      .action-btn--danger:hover {
        background: #fef2f2;
        border-color: #fecaca;
        color: #dc2626;
      }

      /* Form Actions */
      .form-actions {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 1.5rem;
        margin-top: 1.5rem;
        border-top: 1px solid #e2e8f0;
      }

      .form-actions__left,
      .form-actions__right {
        display: flex;
        gap: 0.75rem;
      }

      /* Loading */
      .loading-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem;
        gap: 1rem;
      }

      .spinner {
        width: 2rem;
        height: 2rem;
        border: 3px solid #e2e8f0;
        border-top-color: #4f46e5;
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
      }

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }

      .loading-container p {
        color: #64748b;
        font-size: 0.875rem;
        margin: 0;
      }

      /* Responsive */
      @media (max-width: 640px) {
        .form-grid {
          grid-template-columns: 1fr;
        }

        .form-actions {
          flex-direction: column;
          gap: 1rem;
        }

        .form-actions__left,
        .form-actions__right {
          width: 100%;
          flex-direction: column;
        }

        .btn {
          width: 100%;
          justify-content: center;
        }

        .section-nav {
          padding: 0.375rem;
        }

        .section-tab {
          padding: 0.5rem 0.75rem;
          font-size: 0.8125rem;
        }

        .section-tab__label {
          display: none;
        }

        .upload-controls {
          flex-direction: column;
        }

        .upload-select {
          width: 100%;
        }

        .upload-btn {
          width: 100%;
          justify-content: center;
        }
      }
    `,
  ],
})
export class ApplicationFormComponent implements OnInit {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private applicationService = inject(ApplicationService);
  private toastService = inject(ToastService);

  // State
  loading = signal(false);
  saving = signal(false);
  submitting = signal(false);
  uploadingDocument = signal(false);
  uploadingFileName = signal('');
  application = signal<AdmissionApplication | null>(null);
  documents = signal<ApplicationDocument[]>([]);
  sessions = signal<{ id: string; name: string; status: string }[]>([]);
  currentSection = signal<FormSection>('student');

  // Document upload state
  selectedDocType: DocumentType | '' = '';

  // Form
  applicationForm!: FormGroup;

  // Options
  classOptions = CLASS_NAMES;
  genderOptions = Object.entries(GENDER_LABELS).map(([value, label]) => ({ value, label }));
  bloodGroupOptions = BLOOD_GROUP_OPTIONS;
  relationOptions = Object.entries(PARENT_RELATION_LABELS).map(([value, label]) => ({ value, label }));
  categoryOptions = CATEGORY_OPTIONS;
  documentTypeOptions = Object.entries(DOCUMENT_TYPE_LABELS).map(([value, label]) => ({
    value: value as DocumentType,
    label,
  }));

  sections: { key: FormSection; label: string }[] = [
    { key: 'student', label: 'Student Details' },
    { key: 'parents', label: 'Parents' },
    { key: 'previousSchool', label: 'Previous School' },
    { key: 'documents', label: 'Documents' },
  ];

  // Computed
  /** True when editing an existing application */
  isEditMode = computed(() => !!this.route.snapshot.paramMap.get('id'));

  currentSectionIndex = computed(() => this.sections.findIndex((s) => s.key === this.currentSection()));

  get parentsArray(): FormArray {
    return this.applicationForm.get('parents') as FormArray;
  }

  ngOnInit(): void {
    this.initForm();
    this.loadSessions();

    const id = this.route.snapshot.paramMap.get('id');
    if (id && id !== 'new') {
      this.loadApplication(id);
    }
  }

  private initForm(): void {
    this.applicationForm = this.fb.group({
      sessionId: ['', Validators.required],
      classApplying: ['', Validators.required],
      firstName: ['', Validators.required],
      middleName: [''],
      lastName: ['', Validators.required],
      dateOfBirth: ['', Validators.required],
      gender: ['', Validators.required],
      bloodGroup: [''],
      nationality: ['Indian'],
      religion: [''],
      category: [''],
      aadhaarNumber: [''],
      addressLine1: [''],
      addressLine2: [''],
      city: [''],
      state: [''],
      postalCode: [''],
      previousSchool: [''],
      previousClass: [''],
      previousPercentage: [null],
      parents: this.fb.array([]),
    });
  }

  private loadSessions(): void {
    this.applicationService.getAvailableSessions().subscribe({
      next: (sessions) => this.sessions.set(sessions),
      error: (err) => console.error('Failed to load sessions:', err),
    });
  }

  private loadApplication(id: string): void {
    this.loading.set(true);
    this.applicationService.getApplication(id).subscribe({
      next: (app) => {
        this.application.set(app);
        this.documents.set(app.documents || []);
        this.patchForm(app);
        this.loading.set(false);
      },
      error: (err) => {
        this.toastService.error('Failed to load application');
        this.loading.set(false);
        console.error('Failed to load application:', err);
      },
    });
  }

  private patchForm(app: AdmissionApplication): void {
    this.applicationForm.patchValue({
      sessionId: app.sessionId,
      classApplying: app.classApplying,
      firstName: app.firstName,
      middleName: app.middleName,
      lastName: app.lastName,
      dateOfBirth: app.dateOfBirth,
      gender: app.gender,
      bloodGroup: app.bloodGroup,
      nationality: app.nationality,
      religion: app.religion,
      category: app.category,
      aadhaarNumber: app.aadhaarNumber,
      addressLine1: app.addressLine1,
      addressLine2: app.addressLine2,
      city: app.city,
      state: app.state,
      postalCode: app.postalCode,
      previousSchool: app.previousSchool,
      previousClass: app.previousClass,
      previousPercentage: app.previousPercentage,
    });

    // Populate parents
    this.parentsArray.clear();
    app.parents?.forEach((parent) => {
      this.parentsArray.push(
        this.fb.group({
          id: [parent.id],
          relation: [parent.relation, Validators.required],
          name: [parent.name, Validators.required],
          phone: [parent.phone],
          email: [parent.email],
          occupation: [parent.occupation],
          education: [parent.education],
          annualIncome: [parent.annualIncome],
        })
      );
    });
  }

  // Section navigation
  setSection(section: FormSection): void {
    this.currentSection.set(section);
  }

  nextSection(): void {
    const currentIndex = this.currentSectionIndex();
    if (currentIndex < this.sections.length - 1) {
      this.currentSection.set(this.sections[currentIndex + 1].key);
    }
  }

  prevSection(): void {
    const currentIndex = this.currentSectionIndex();
    if (currentIndex > 0) {
      this.currentSection.set(this.sections[currentIndex - 1].key);
    }
  }

  isSectionCompleted(section: FormSection): boolean {
    switch (section) {
      case 'student':
        return (
          this.applicationForm.get('firstName')?.valid === true &&
          this.applicationForm.get('lastName')?.valid === true &&
          this.applicationForm.get('dateOfBirth')?.valid === true &&
          this.applicationForm.get('gender')?.valid === true
        );
      case 'parents':
        return this.parentsArray.length > 0 && this.parentsArray.valid;
      case 'previousSchool':
        return true; // Optional section
      case 'documents':
        return this.documents().length > 0;
      default:
        return false;
    }
  }

  // Parent management
  addParent(): void {
    this.parentsArray.push(
      this.fb.group({
        relation: ['', Validators.required],
        name: ['', Validators.required],
        phone: [''],
        email: [''],
        occupation: [''],
        education: [''],
        annualIncome: [''],
      })
    );
  }

  removeParent(index: number): void {
    this.parentsArray.removeAt(index);
  }

  getParentTitle(index: number): string {
    const relation = this.parentsArray.at(index).get('relation')?.value;
    const name = this.parentsArray.at(index).get('name')?.value;
    if (relation && name) {
      return `${PARENT_RELATION_LABELS[relation as ParentRelation]} - ${name}`;
    }
    if (relation) {
      return PARENT_RELATION_LABELS[relation as ParentRelation];
    }
    return `Parent/Guardian ${index + 1}`;
  }

  // Document management
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];

    if (!file) return;

    if (!this.selectedDocType) {
      this.toastService.error('Please select a document type first');
      input.value = '';
      return;
    }

    // Validate file size (5MB max)
    if (file.size > 5 * 1024 * 1024) {
      this.toastService.error('File size exceeds 5MB limit');
      input.value = '';
      return;
    }

    const appId = this.application()?.id;
    if (!appId) {
      // If no application exists yet, we need to save first
      this.toastService.info('Please save the application first before uploading documents');
      input.value = '';
      return;
    }

    this.uploadingDocument.set(true);
    this.uploadingFileName.set(file.name);

    this.applicationService.uploadDocument(appId, this.selectedDocType, file).subscribe({
      next: (doc) => {
        this.documents.update((docs) => [...docs, doc]);
        this.toastService.success('Document uploaded successfully');
        this.uploadingDocument.set(false);
        this.uploadingFileName.set('');
        this.selectedDocType = '';
        input.value = '';
      },
      error: (err) => {
        this.toastService.error('Failed to upload document');
        this.uploadingDocument.set(false);
        this.uploadingFileName.set('');
        input.value = '';
        console.error('Failed to upload document:', err);
      },
    });
  }

  deleteDocument(doc: ApplicationDocument): void {
    const appId = this.application()?.id;
    if (!appId) return;

    this.applicationService.deleteDocument(appId, doc.id).subscribe({
      next: () => {
        this.documents.update((docs) => docs.filter((d) => d.id !== doc.id));
        this.toastService.success('Document deleted');
      },
      error: (err) => {
        this.toastService.error('Failed to delete document');
        console.error('Failed to delete document:', err);
      },
    });
  }

  getDocumentIcon(fileName: string): string {
    const ext = fileName.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return 'fa-solid fa-file-pdf';
      case 'jpg':
      case 'jpeg':
      case 'png':
        return 'fa-solid fa-file-image';
      default:
        return 'fa-solid fa-file';
    }
  }

  getDocumentTypeLabel(type: DocumentType): string {
    return DOCUMENT_TYPE_LABELS[type] || type;
  }

  // Helper to extract parent data from parents array
  private extractParentData(): Partial<CreateApplicationRequest> {
    const parentData: Partial<CreateApplicationRequest> = {};
    const parents = this.parentsArray.value;

    for (const parent of parents) {
      if (parent.relation === 'father') {
        parentData.fatherName = parent.name;
        parentData.fatherPhone = parent.phone;
        parentData.fatherEmail = parent.email;
        parentData.fatherOccupation = parent.occupation;
      } else if (parent.relation === 'mother') {
        parentData.motherName = parent.name;
        parentData.motherPhone = parent.phone;
        parentData.motherEmail = parent.email;
        parentData.motherOccupation = parent.occupation;
      } else if (parent.relation === 'guardian' || parent.relation === 'other') {
        parentData.guardianName = parent.name;
        parentData.guardianPhone = parent.phone;
        parentData.guardianEmail = parent.email;
        parentData.guardianRelation = parent.relation;
      }
    }

    return parentData;
  }

  // Form submission
  saveDraft(): void {
    if (this.applicationForm.invalid) {
      this.markFormGroupTouched(this.applicationForm);
      this.toastService.error('Please fill in all required fields');
      return;
    }

    this.saving.set(true);
    const formValue = this.applicationForm.value;
    const parentData = this.extractParentData();
    const data: CreateApplicationRequest = {
      sessionId: formValue.sessionId,
      classApplying: formValue.classApplying,
      firstName: formValue.firstName,
      middleName: formValue.middleName,
      lastName: formValue.lastName,
      dateOfBirth: formValue.dateOfBirth,
      gender: formValue.gender,
      bloodGroup: formValue.bloodGroup,
      nationality: formValue.nationality,
      religion: formValue.religion,
      category: formValue.category,
      aadhaarNumber: formValue.aadhaarNumber,
      addressLine1: formValue.addressLine1,
      addressLine2: formValue.addressLine2,
      city: formValue.city,
      state: formValue.state,
      postalCode: formValue.postalCode,
      previousSchool: formValue.previousSchool,
      previousClass: formValue.previousClass,
      previousPercentage: formValue.previousPercentage,
      ...parentData,
    };

    const existingApp = this.application();

    if (existingApp) {
      this.applicationService.updateApplication(existingApp.id, data).subscribe({
        next: (app) => {
          this.application.set(app);
          this.toastService.success('Application saved as draft');
          this.saving.set(false);
        },
        error: (err) => {
          this.toastService.error('Failed to save application');
          this.saving.set(false);
          console.error('Failed to save application:', err);
        },
      });
    } else {
      this.applicationService.createApplication(data).subscribe({
        next: (app) => {
          this.application.set(app);
          this.toastService.success('Application created as draft');
          this.saving.set(false);
          // Update URL to include the new ID
          this.router.navigate(['/admissions/applications', app.id, 'edit'], { replaceUrl: true });
        },
        error: (err) => {
          this.toastService.error('Failed to create application');
          this.saving.set(false);
          console.error('Failed to create application:', err);
        },
      });
    }
  }

  onSubmit(): void {
    if (this.applicationForm.invalid) {
      this.markFormGroupTouched(this.applicationForm);
      this.toastService.error('Please fill in all required fields');
      return;
    }

    if (this.parentsArray.length === 0) {
      this.toastService.error('Please add at least one parent/guardian');
      this.setSection('parents');
      return;
    }

    this.submitting.set(true);

    // First save the application
    this.saveDraftAndSubmit();
  }

  private saveDraftAndSubmit(): void {
    const formValue = this.applicationForm.value;
    const parentData = this.extractParentData();
    const data: CreateApplicationRequest = {
      sessionId: formValue.sessionId,
      classApplying: formValue.classApplying,
      firstName: formValue.firstName,
      middleName: formValue.middleName,
      lastName: formValue.lastName,
      dateOfBirth: formValue.dateOfBirth,
      gender: formValue.gender,
      bloodGroup: formValue.bloodGroup,
      nationality: formValue.nationality,
      religion: formValue.religion,
      category: formValue.category,
      aadhaarNumber: formValue.aadhaarNumber,
      addressLine1: formValue.addressLine1,
      addressLine2: formValue.addressLine2,
      city: formValue.city,
      state: formValue.state,
      postalCode: formValue.postalCode,
      previousSchool: formValue.previousSchool,
      previousClass: formValue.previousClass,
      previousPercentage: formValue.previousPercentage,
      ...parentData,
    };

    const existingApp = this.application();

    const saveObs = existingApp
      ? this.applicationService.updateApplication(existingApp.id, data)
      : this.applicationService.createApplication(data);

    saveObs.subscribe({
      next: (app) => {
        // Save parents
        this.saveParentsAndSubmit(app);
      },
      error: (err) => {
        this.toastService.error('Failed to save application');
        this.submitting.set(false);
        console.error('Failed to save application:', err);
      },
    });
  }

  private saveParentsAndSubmit(app: AdmissionApplication): void {
    // For simplicity, submit the application
    // In a real implementation, we would save parents first
    this.applicationService.submitApplication(app.id).subscribe({
      next: () => {
        this.toastService.success('Application submitted successfully!');
        this.submitting.set(false);
        this.router.navigate(['/admissions/applications', app.id]);
      },
      error: (err) => {
        this.toastService.error('Failed to submit application');
        this.submitting.set(false);
        console.error('Failed to submit application:', err);
      },
    });
  }

  // Helpers
  showError(field: string): boolean {
    const control = this.applicationForm.get(field);
    return control ? control.invalid && (control.dirty || control.touched) : false;
  }

  private markFormGroupTouched(formGroup: FormGroup): void {
    Object.values(formGroup.controls).forEach((control) => {
      control.markAsTouched();
      if (control instanceof FormGroup) {
        this.markFormGroupTouched(control);
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/admissions/applications']);
  }
}
