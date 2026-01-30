import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';

/**
 * AdmissionsComponent - Admissions module landing page.
 * Shows admission workflow steps and navigation to sub-modules.
 */
@Component({
  selector: 'msls-admissions',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-user-plus"></i>
          </div>
          <div class="header-text">
            <h1>Admissions Management</h1>
            <p>Manage the complete admission lifecycle from enquiry to enrollment</p>
          </div>
        </div>
        <a routerLink="dashboard" class="btn btn-primary">
          <i class="fa-solid fa-chart-pie"></i>
          View Dashboard
        </a>
      </div>

      <!-- Quick Actions Grid -->
      <div class="config-grid">
        <!-- Sessions Card -->
        <a routerLink="sessions" class="config-card">
          <div class="card-icon sessions">
            <i class="fa-solid fa-calendar-alt"></i>
          </div>
          <div class="card-content">
            <h3>Admission Sessions</h3>
            <p>Create and manage admission sessions with dates, fees, and seat allocation</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Enquiries Card -->
        <a routerLink="enquiries" class="config-card">
          <div class="card-icon enquiries">
            <i class="fa-solid fa-comments"></i>
          </div>
          <div class="card-content">
            <h3>Enquiries</h3>
            <p>Track and follow up on admission enquiries from parents</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Applications Card -->
        <a routerLink="applications" class="config-card featured">
          <div class="card-icon applications">
            <i class="fa-solid fa-file-pen"></i>
          </div>
          <div class="card-content">
            <h3>Applications</h3>
            <p>Review and process student admission applications</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Tests Card -->
        <a routerLink="tests" class="config-card">
          <div class="card-icon tests">
            <i class="fa-solid fa-clipboard-check"></i>
          </div>
          <div class="card-content">
            <h3>Entrance Tests</h3>
            <p>Schedule tests, enter results, and manage test registrations</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Merit List Card -->
        <a routerLink="merit-list" class="config-card">
          <div class="card-icon merit">
            <i class="fa-solid fa-trophy"></i>
          </div>
          <div class="card-content">
            <h3>Merit List</h3>
            <p>Generate merit lists and make admission decisions</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Enrollment Card -->
        <a routerLink="enrollment" class="config-card">
          <div class="card-icon enrollment">
            <i class="fa-solid fa-user-check"></i>
          </div>
          <div class="card-content">
            <h3>Enrollment</h3>
            <p>Convert admitted students to enrolled students</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>
      </div>

      <!-- How It Works Section -->
      <div class="info-section">
        <h2>How It Works</h2>
        <div class="workflow-steps">
          <div class="workflow-step">
            <div class="step-number">1</div>
            <div class="step-icon"><i class="fa-solid fa-calendar-plus"></i></div>
            <div class="step-content">
              <h4>Create Session</h4>
              <p>Set up admission session with academic year, dates, application fees, and class-wise seat allocation</p>
            </div>
          </div>
          <div class="workflow-arrow"><i class="fa-solid fa-arrow-right"></i></div>

          <div class="workflow-step">
            <div class="step-number">2</div>
            <div class="step-icon"><i class="fa-solid fa-phone"></i></div>
            <div class="step-content">
              <h4>Receive Enquiries</h4>
              <p>Capture walk-in and phone enquiries, schedule follow-ups, and convert to applications</p>
            </div>
          </div>
          <div class="workflow-arrow"><i class="fa-solid fa-arrow-right"></i></div>

          <div class="workflow-step">
            <div class="step-number">3</div>
            <div class="step-icon"><i class="fa-solid fa-file-import"></i></div>
            <div class="step-content">
              <h4>Collect Applications</h4>
              <p>Accept online/offline applications with student details, documents, and application fee payment</p>
            </div>
          </div>
          <div class="workflow-arrow"><i class="fa-solid fa-arrow-right"></i></div>

          <div class="workflow-step">
            <div class="step-number">4</div>
            <div class="step-icon"><i class="fa-solid fa-pen-to-square"></i></div>
            <div class="step-content">
              <h4>Conduct Tests</h4>
              <p>Schedule entrance tests, register candidates, conduct tests, and enter results</p>
            </div>
          </div>
          <div class="workflow-arrow"><i class="fa-solid fa-arrow-right"></i></div>

          <div class="workflow-step">
            <div class="step-number">5</div>
            <div class="step-icon"><i class="fa-solid fa-list-ol"></i></div>
            <div class="step-content">
              <h4>Generate Merit</h4>
              <p>Auto-generate merit list based on test scores, review applications, and make admission decisions</p>
            </div>
          </div>
          <div class="workflow-arrow"><i class="fa-solid fa-arrow-right"></i></div>

          <div class="workflow-step">
            <div class="step-number">6</div>
            <div class="step-icon"><i class="fa-solid fa-graduation-cap"></i></div>
            <div class="step-content">
              <h4>Enroll Students</h4>
              <p>Collect admission fees, complete enrollment, and create student records</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Application Status Flow -->
      <div class="status-flow-section">
        <h2>Application Status Flow</h2>
        <div class="status-flow">
          <div class="status-item draft">
            <div class="status-icon"><i class="fa-solid fa-file"></i></div>
            <span class="status-label">Draft</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item submitted">
            <div class="status-icon"><i class="fa-solid fa-paper-plane"></i></div>
            <span class="status-label">Submitted</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item review">
            <div class="status-icon"><i class="fa-solid fa-magnifying-glass"></i></div>
            <span class="status-label">Under Review</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item test">
            <div class="status-icon"><i class="fa-solid fa-clipboard-check"></i></div>
            <span class="status-label">Test Scheduled</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item admitted">
            <div class="status-icon"><i class="fa-solid fa-check-circle"></i></div>
            <span class="status-label">Admitted</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item enrolled">
            <div class="status-icon"><i class="fa-solid fa-user-graduate"></i></div>
            <span class="status-label">Enrolled</span>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1400px; margin: 0 auto; }

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 2rem;
      gap: 1rem;
    }

    .header-content { display: flex; align-items: center; gap: 1rem; }
    .header-icon {
      width: 3.5rem;
      height: 3.5rem;
      border-radius: 1rem;
      background: linear-gradient(135deg, #8b5cf6, #7c3aed);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
    }
    .header-text h1 { margin: 0; font-size: 1.75rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.9375rem; }

    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1.5rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      text-decoration: none;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-primary { background: #8b5cf6; color: white; }
    .btn-primary:hover { background: #7c3aed; }

    .config-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 1.25rem;
      margin-bottom: 2.5rem;
    }

    .config-card {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1.25rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      text-decoration: none;
      transition: all 0.2s;
    }

    .config-card:hover {
      border-color: #c4b5fd;
      box-shadow: 0 4px 12px rgba(139, 92, 246, 0.08);
      transform: translateY(-2px);
    }

    .config-card.featured {
      border: 2px solid #c4b5fd;
      background: linear-gradient(135deg, #faf5ff, #f3e8ff);
    }

    .config-card.featured:hover {
      border-color: #a78bfa;
      box-shadow: 0 4px 16px rgba(139, 92, 246, 0.15);
    }

    .card-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      flex-shrink: 0;
    }

    .card-icon.sessions { background: #dbeafe; color: #2563eb; }
    .card-icon.enquiries { background: #fef3c7; color: #d97706; }
    .card-icon.applications { background: linear-gradient(135deg, #8b5cf6, #7c3aed); color: white; }
    .card-icon.tests { background: #e0e7ff; color: #4f46e5; }
    .card-icon.merit { background: #fef3c7; color: #b45309; }
    .card-icon.enrollment { background: #dcfce7; color: #16a34a; }

    .card-content { flex: 1; }
    .card-content h3 { margin: 0 0 0.25rem; font-size: 1rem; font-weight: 600; color: #1e293b; }
    .card-content p { margin: 0; font-size: 0.8125rem; color: #64748b; line-height: 1.4; }

    .card-arrow { color: #94a3b8; transition: transform 0.2s; }
    .config-card:hover .card-arrow { transform: translateX(4px); color: #8b5cf6; }

    .info-section {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
      margin-bottom: 1.5rem;
    }

    .info-section h2 {
      margin: 0 0 1.5rem;
      font-size: 1.125rem;
      font-weight: 600;
      color: #1e293b;
    }

    .workflow-steps {
      display: flex;
      align-items: flex-start;
      gap: 0.5rem;
      overflow-x: auto;
      padding: 0.75rem 0.5rem 0.5rem;
    }

    .workflow-step {
      flex: 1;
      min-width: 160px;
      max-width: 200px;
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;
      position: relative;
    }

    .step-number {
      position: absolute;
      top: -4px;
      right: calc(50% - 2.5rem);
      width: 1.5rem;
      height: 1.5rem;
      border-radius: 50%;
      background: #8b5cf6;
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.75rem;
      font-weight: 600;
      z-index: 2;
      box-shadow: 0 2px 4px rgba(139, 92, 246, 0.3);
    }

    .step-icon {
      width: 3.5rem;
      height: 3.5rem;
      border-radius: 1rem;
      background: white;
      border: 2px solid #e2e8f0;
      color: #8b5cf6;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      margin-bottom: 0.75rem;
    }

    .step-content h4 { margin: 0 0 0.25rem; font-size: 0.875rem; font-weight: 600; color: #1e293b; }
    .step-content p { margin: 0; font-size: 0.6875rem; color: #64748b; line-height: 1.4; }

    .workflow-arrow {
      display: flex;
      align-items: center;
      color: #cbd5e1;
      font-size: 1rem;
      margin-top: 1.5rem;
    }

    .status-flow-section {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    .status-flow-section h2 {
      margin: 0 0 1.25rem;
      font-size: 1.125rem;
      font-weight: 600;
      color: #1e293b;
    }

    .status-flow {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      flex-wrap: wrap;
    }

    .status-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.375rem;
      padding: 0.75rem;
      border-radius: 0.75rem;
      min-width: 90px;
    }

    .status-item.draft { background: #f1f5f9; }
    .status-item.submitted { background: #dbeafe; }
    .status-item.review { background: #fef3c7; }
    .status-item.test { background: #e0e7ff; }
    .status-item.admitted { background: #dcfce7; }
    .status-item.enrolled { background: #d1fae5; }

    .status-icon {
      width: 2.25rem;
      height: 2.25rem;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.875rem;
    }

    .status-item.draft .status-icon { background: #e2e8f0; color: #64748b; }
    .status-item.submitted .status-icon { background: #93c5fd; color: #1e40af; }
    .status-item.review .status-icon { background: #fde68a; color: #92400e; }
    .status-item.test .status-icon { background: #a5b4fc; color: #3730a3; }
    .status-item.admitted .status-icon { background: #86efac; color: #166534; }
    .status-item.enrolled .status-icon { background: #6ee7b7; color: #065f46; }

    .status-label { font-weight: 500; font-size: 0.6875rem; color: #1e293b; }
    .status-arrow { color: #cbd5e1; font-size: 0.875rem; }

    @media (max-width: 768px) {
      .page-header { flex-direction: column; align-items: flex-start; }
      .workflow-steps { flex-direction: column; align-items: center; }
      .workflow-arrow { transform: rotate(90deg); margin: 0.5rem 0; }
      .workflow-step { max-width: 100%; width: 100%; }
      .status-flow { flex-direction: column; }
      .status-arrow { transform: rotate(90deg); }
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdmissionsComponent {}
