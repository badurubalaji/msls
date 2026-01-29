import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MslsModalComponent } from '../../../shared/components/modal/modal.component';
import { StreamService } from '../services/stream.service';
import { Stream, CreateStreamRequest, UpdateStreamRequest } from '../academic.model';
import { ToastService } from '../../../shared/services/toast.service';

@Component({
  selector: 'msls-streams',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsModalComponent],
  template: `
    <div class="page">
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-stream"></i>
          </div>
          <div class="header-text">
            <h1>Streams</h1>
            <p>Manage academic streams (Science, Commerce, Arts) for senior classes</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          Add Stream
        </button>
      </div>

      <div class="filters-bar">
        <div class="search-box">
          <i class="fa-solid fa-search search-icon"></i>
          <input type="text" placeholder="Search streams..." [ngModel]="searchTerm()" (ngModelChange)="searchTerm.set($event)" class="search-input" />
        </div>
        <div class="filter-group">
          <select class="filter-select" [ngModel]="statusFilter()" (ngModelChange)="statusFilter.set($event)">
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>

      <div class="content-card">
        @if (loading()) {
          <div class="loading-container"><div class="spinner"></div><span>Loading streams...</span></div>
        } @else if (error()) {
          <div class="error-container"><i class="fa-solid fa-circle-exclamation"></i><span>{{ error() }}</span></div>
        } @else {
          <table class="data-table">
            <thead>
              <tr>
                <th>Stream</th>
                <th>Code</th>
                <th>Description</th>
                <th>Order</th>
                <th>Status</th>
                <th style="width: 140px; text-align: right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              @for (stream of filteredStreams(); track stream.id) {
                <tr>
                  <td>
                    <div class="name-wrapper">
                      <div class="stream-icon" [class]="'stream-' + stream.code.toLowerCase()">
                        <i class="fa-solid fa-stream"></i>
                      </div>
                      <span class="name">{{ stream.name }}</span>
                    </div>
                  </td>
                  <td><span class="code-badge">{{ stream.code }}</span></td>
                  <td class="desc-cell">{{ stream.description || '-' }}</td>
                  <td>{{ stream.displayOrder }}</td>
                  <td>
                    <span class="badge" [class.badge-green]="stream.isActive" [class.badge-gray]="!stream.isActive">
                      {{ stream.isActive ? 'Active' : 'Inactive' }}
                    </span>
                  </td>
                  <td class="actions-cell">
                    <button class="action-btn" title="Edit" (click)="editStream(stream)"><i class="fa-regular fa-pen-to-square"></i></button>
                    <button class="action-btn" title="Toggle" (click)="toggleStatus(stream)">
                      <i class="fa-solid" [class.fa-toggle-on]="stream.isActive" [class.fa-toggle-off]="!stream.isActive"></i>
                    </button>
                    <button class="action-btn action-btn--danger" title="Delete" (click)="confirmDelete(stream)"><i class="fa-regular fa-trash-can"></i></button>
                  </td>
                </tr>
              } @empty {
                <tr><td colspan="6" class="empty-cell"><div class="empty-state"><i class="fa-regular fa-folder-open"></i><p>No streams found</p></div></td></tr>
              }
            </tbody>
          </table>
        }
      </div>

      <msls-modal [isOpen]="showStreamModal()" [title]="editingStream() ? 'Edit Stream' : 'Create Stream'" size="md" (closed)="closeStreamModal()">
        <form class="form" (ngSubmit)="saveStream()">
          <div class="form-group">
            <label for="streamName">Stream Name <span class="required">*</span></label>
            <input type="text" id="streamName" [(ngModel)]="formData.name" name="name" placeholder="e.g., Science" required />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label for="streamCode">Code <span class="required">*</span></label>
              <input type="text" id="streamCode" [(ngModel)]="formData.code" name="code" placeholder="e.g., SCI" required />
            </div>
            <div class="form-group">
              <label for="displayOrder">Display Order</label>
              <input type="number" id="displayOrder" [(ngModel)]="formData.displayOrder" name="displayOrder" min="0" />
            </div>
          </div>
          <div class="form-group">
            <label for="description">Description</label>
            <textarea id="description" [(ngModel)]="formData.description" name="description" rows="3" placeholder="Optional description..."></textarea>
          </div>
          <div class="form-actions">
            <button type="button" class="btn btn-secondary" (click)="closeStreamModal()">Cancel</button>
            <button type="submit" class="btn btn-primary" [disabled]="saving()">
              @if (saving()) { <div class="btn-spinner"></div> Saving... } @else { {{ editingStream() ? 'Update' : 'Create' }} }
            </button>
          </div>
        </form>
      </msls-modal>

      <msls-modal [isOpen]="showDeleteModal()" title="Delete Stream" size="sm" (closed)="closeDeleteModal()">
        <div class="delete-confirmation">
          <div class="delete-icon"><i class="fa-solid fa-triangle-exclamation"></i></div>
          <p>Are you sure you want to delete <strong>"{{ streamToDelete()?.name }}"</strong>?</p>
          <p class="delete-warning">Streams in use by classes cannot be deleted.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
            <button class="btn btn-danger" [disabled]="deleting()" (click)="deleteStream()">@if (deleting()) { Deleting... } @else { Delete }</button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1200px; margin: 0 auto; }
    .page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem; }
    .header-content { display: flex; align-items: center; gap: 1rem; }
    .header-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; background: #e0e7ff; color: #4338ca; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; }
    .header-text h1 { margin: 0; font-size: 1.5rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.875rem; }
    .filters-bar { display: flex; gap: 1rem; margin-bottom: 1rem; }
    .search-box { flex: 1; max-width: 400px; position: relative; }
    .search-icon { position: absolute; left: 0.875rem; top: 50%; transform: translateY(-50%); color: #9ca3af; }
    .search-input { width: 100%; padding: 0.625rem 2.5rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; }
    .filter-select { padding: 0.625rem 2rem 0.625rem 0.875rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; background: white; }
    .content-card { background: white; border: 1px solid #e2e8f0; border-radius: 1rem; overflow: hidden; }
    .loading-container { display: flex; align-items: center; justify-content: center; gap: 1rem; padding: 3rem; color: #64748b; }
    .spinner { width: 24px; height: 24px; border: 3px solid #e2e8f0; border-top-color: #4f46e5; border-radius: 50%; animation: spin 0.8s linear infinite; }
    @keyframes spin { to { transform: rotate(360deg); } }
    .data-table { width: 100%; border-collapse: collapse; }
    .data-table th { text-align: left; padding: 0.875rem 1rem; font-size: 0.75rem; font-weight: 600; text-transform: uppercase; color: #64748b; background: #f8fafc; border-bottom: 1px solid #e2e8f0; }
    .data-table td { padding: 1rem; border-bottom: 1px solid #f1f5f9; color: #374151; }
    .data-table tbody tr:hover { background: #f8fafc; }
    .name-wrapper { display: flex; align-items: center; gap: 0.75rem; }
    .stream-icon { width: 2.5rem; height: 2.5rem; border-radius: 0.5rem; display: flex; align-items: center; justify-content: center; }
    .stream-sci { background: #dbeafe; color: #1d4ed8; }
    .stream-com { background: #dcfce7; color: #15803d; }
    .stream-arts { background: #fae8ff; color: #a21caf; }
    .name { font-weight: 500; color: #1e293b; }
    .code-badge { display: inline-flex; padding: 0.25rem 0.5rem; background: #f1f5f9; border-radius: 0.25rem; font-family: monospace; font-size: 0.75rem; }
    .desc-cell { max-width: 300px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: #64748b; }
    .badge { display: inline-flex; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 500; }
    .badge-green { background: #dcfce7; color: #166534; }
    .badge-gray { background: #f1f5f9; color: #64748b; }
    .actions-cell { text-align: right; }
    .action-btn { display: inline-flex; align-items: center; justify-content: center; width: 2rem; height: 2rem; border: none; background: transparent; color: #64748b; border-radius: 0.375rem; cursor: pointer; }
    .action-btn:hover { background: #f1f5f9; color: #4f46e5; }
    .action-btn--danger:hover { background: #fef2f2; color: #dc2626; }
    .empty-cell { padding: 3rem !important; }
    .empty-state { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; color: #64748b; }
    .btn { display: inline-flex; align-items: center; gap: 0.5rem; padding: 0.625rem 1.25rem; border-radius: 0.5rem; font-size: 0.875rem; font-weight: 500; cursor: pointer; border: none; }
    .btn-primary { background: #4f46e5; color: white; }
    .btn-secondary { background: #f1f5f9; color: #475569; }
    .btn-danger { background: #dc2626; color: white; }
    .btn:disabled { opacity: 0.5; cursor: not-allowed; }
    .btn-spinner { width: 16px; height: 16px; border: 2px solid transparent; border-top-color: currentColor; border-radius: 50%; animation: spin 0.8s linear infinite; }
    .form { display: flex; flex-direction: column; gap: 1rem; }
    .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    .form-group { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-group label { font-size: 0.875rem; font-weight: 500; color: #374151; }
    .required { color: #dc2626; }
    .form-group input, .form-group select, .form-group textarea { padding: 0.625rem 0.875rem; border: 1px solid #e2e8f0; border-radius: 0.5rem; font-size: 0.875rem; }
    .form-actions { display: flex; justify-content: flex-end; gap: 0.75rem; padding-top: 1rem; border-top: 1px solid #e2e8f0; }
    .delete-confirmation { text-align: center; padding: 1rem; }
    .delete-icon { width: 4rem; height: 4rem; margin: 0 auto 1rem; border-radius: 50%; background: #fef2f2; color: #dc2626; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
    .delete-warning { font-size: 0.875rem; color: #64748b; }
    .delete-actions { display: flex; gap: 0.75rem; justify-content: center; margin-top: 1.5rem; }
  `],
})
export class StreamsComponent implements OnInit {
  private streamService = inject(StreamService);
  private toastService = inject(ToastService);

  streams = signal<Stream[]>([]);
  loading = signal(true);
  saving = signal(false);
  deleting = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  statusFilter = signal<'all' | 'active' | 'inactive'>('all');

  showStreamModal = signal(false);
  showDeleteModal = signal(false);
  editingStream = signal<Stream | null>(null);
  streamToDelete = signal<Stream | null>(null);

  formData = { name: '', code: '', description: '', displayOrder: 0 };

  filteredStreams = computed(() => {
    let result = this.streams();
    const term = this.searchTerm().toLowerCase();
    if (term) result = result.filter(s => s.name.toLowerCase().includes(term) || s.code.toLowerCase().includes(term));
    if (this.statusFilter() === 'active') result = result.filter(s => s.isActive);
    else if (this.statusFilter() === 'inactive') result = result.filter(s => !s.isActive);
    return result;
  });

  ngOnInit(): void { this.loadStreams(); }

  loadStreams(): void {
    this.loading.set(true);
    this.streamService.getStreams().subscribe({
      next: streams => { this.streams.set(streams); this.loading.set(false); },
      error: () => { this.error.set('Failed to load streams'); this.loading.set(false); },
    });
  }

  openCreateModal(): void {
    this.editingStream.set(null);
    this.formData = { name: '', code: '', description: '', displayOrder: 0 };
    this.showStreamModal.set(true);
  }

  editStream(stream: Stream): void {
    this.editingStream.set(stream);
    this.formData = { name: stream.name, code: stream.code, description: stream.description || '', displayOrder: stream.displayOrder };
    this.showStreamModal.set(true);
  }

  closeStreamModal(): void { this.showStreamModal.set(false); this.editingStream.set(null); }

  saveStream(): void {
    if (!this.formData.name || !this.formData.code) { this.toastService.error('Please fill required fields'); return; }
    this.saving.set(true);
    const editing = this.editingStream();
    const data: CreateStreamRequest | UpdateStreamRequest = {
      name: this.formData.name, code: this.formData.code, description: this.formData.description || undefined, displayOrder: this.formData.displayOrder,
    };
    const operation = editing ? this.streamService.updateStream(editing.id, data) : this.streamService.createStream(data as CreateStreamRequest);
    operation.subscribe({
      next: () => { this.toastService.success(editing ? 'Stream updated' : 'Stream created'); this.closeStreamModal(); this.loadStreams(); this.saving.set(false); },
      error: () => { this.toastService.error(editing ? 'Failed to update' : 'Failed to create'); this.saving.set(false); },
    });
  }

  toggleStatus(stream: Stream): void {
    this.streamService.updateStream(stream.id, { isActive: !stream.isActive }).subscribe({
      next: () => { this.toastService.success('Status updated'); this.loadStreams(); },
      error: () => this.toastService.error('Failed to update status'),
    });
  }

  confirmDelete(stream: Stream): void { this.streamToDelete.set(stream); this.showDeleteModal.set(true); }
  closeDeleteModal(): void { this.showDeleteModal.set(false); this.streamToDelete.set(null); }

  deleteStream(): void {
    const stream = this.streamToDelete();
    if (!stream) return;
    this.deleting.set(true);
    this.streamService.deleteStream(stream.id).subscribe({
      next: () => { this.toastService.success('Stream deleted'); this.closeDeleteModal(); this.loadStreams(); this.deleting.set(false); },
      error: () => { this.toastService.error('Failed to delete stream'); this.deleting.set(false); },
    });
  }
}
