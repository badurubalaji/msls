/**
 * MSLS Role Permissions Component
 *
 * Component for viewing and editing role permissions.
 */

import { Component, Input, Output, EventEmitter, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { ToastService } from '../../../shared/services';
import { RbacService } from '../../../core/services';
import { Role, RbacPermission } from '../../../core/models';

@Component({
  selector: 'msls-role-permissions',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="role-permissions">
      @if (loading()) {
        <div class="permissions-loading">
          <div class="spinner"></div>
          <p>Loading permissions...</p>
        </div>
      } @else if (allPermissions().length === 0) {
        <!-- Empty State -->
        <div class="permissions-empty">
          <i class="fa-solid fa-shield-halved"></i>
          <h3>No Permissions Available</h3>
          <p>No permissions have been configured in the system yet.</p>
        </div>
      } @else {
        <!-- Permissions Summary -->
        <div class="permissions-header">
          <p class="permissions-summary">
            <span class="count-badge">{{ selectedCount() }}</span> of
            <strong>{{ allPermissions().length }}</strong> permissions assigned
          </p>
          @if (canEdit) {
            <div class="permissions-actions">
              <button class="action-btn" (click)="selectAll()">
                <i class="fa-regular fa-square-check"></i>
                Select All
              </button>
              <button class="action-btn" (click)="deselectAll()">
                <i class="fa-regular fa-square"></i>
                Deselect All
              </button>
            </div>
          }
        </div>

        <!-- Permissions Grid -->
        <div class="permissions-grid">
          @for (module of modules(); track module) {
            <div class="module-card">
              <div class="module-header">
                <div class="module-title">
                  <i [class]="getModuleIcon(module)"></i>
                  <h3 class="module-name">{{ formatModuleName(module) }}</h3>
                </div>
                <span class="module-count">{{ getModulePermissionCount(module) }}</span>
              </div>
              <div class="module-permissions">
                @for (permission of getPermissionsByModule(module); track permission.id) {
                  <div class="permission-item" [class.permission-item--selected]="isPermissionSelected(permission.id)">
                    @if (canEdit) {
                      <label class="permission-label">
                        <input
                          type="checkbox"
                          [checked]="isPermissionSelected(permission.id)"
                          (change)="togglePermission(permission)"
                          class="permission-checkbox"
                        />
                        <span class="permission-name">{{ permission.name }}</span>
                      </label>
                    } @else {
                      <div class="permission-readonly">
                        <span
                          class="permission-indicator"
                          [class.permission-indicator--selected]="isPermissionSelected(permission.id)"
                        >
                          @if (isPermissionSelected(permission.id)) {
                            <i class="fa-solid fa-check"></i>
                          }
                        </span>
                        <span class="permission-name">{{ permission.name }}</span>
                      </div>
                    }
                    @if (permission.description) {
                      <p class="permission-description">{{ permission.description }}</p>
                    }
                  </div>
                }
              </div>
            </div>
          }
        </div>

        <!-- Footer Actions (only show when there are changes to save) -->
        @if (canEdit && hasChanges()) {
          <div class="permissions-footer">
            <button class="btn btn-secondary" (click)="resetChanges()">
              <i class="fa-solid fa-rotate-left"></i>
              Reset
            </button>
            <button class="btn btn-primary" [disabled]="saving()" (click)="saveChanges()">
              @if (saving()) {
                <span class="btn-spinner"></span>
                Saving...
              } @else {
                <i class="fa-solid fa-check"></i>
                Save Changes
              }
            </button>
          </div>
        }
      }
    </div>
  `,
  styles: [`
    .role-permissions {
      display: flex;
      flex-direction: column;
      background: #ffffff;
      overflow: hidden;
    }

    /* Loading State */
    .permissions-loading {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 3rem;
      gap: 1rem;
    }

    .permissions-loading p {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0;
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
      to { transform: rotate(360deg); }
    }

    /* Empty State */
    .permissions-empty {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 3rem;
      text-align: center;
    }

    .permissions-empty i {
      font-size: 3rem;
      color: #cbd5e1;
      margin-bottom: 1rem;
    }

    .permissions-empty h3 {
      font-size: 1rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0 0 0.5rem 0;
    }

    .permissions-empty p {
      font-size: 0.875rem;
      color: #64748b;
      margin: 0;
    }

    /* Permissions Header */
    .permissions-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 0.75rem;
      padding: 1rem 1.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .permissions-summary {
      color: #475569;
      font-size: 0.875rem;
      margin: 0;
    }

    .count-badge {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 1.5rem;
      height: 1.5rem;
      padding: 0 0.5rem;
      background: #4f46e5;
      color: #ffffff;
      font-size: 0.8125rem;
      font-weight: 600;
      border-radius: 9999px;
    }

    .permissions-actions {
      display: flex;
      gap: 0.5rem;
    }

    .action-btn {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.375rem 0.75rem;
      font-size: 0.8125rem;
      font-weight: 500;
      color: #475569;
      background: transparent;
      border: 1px solid #e2e8f0;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: all 0.15s;
    }

    .action-btn:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
      color: #0f172a;
    }

    /* Permissions Grid */
    .permissions-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
      gap: 1rem;
      padding: 1.5rem;
      overflow-y: auto;
      max-height: 400px;
    }

    .module-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      overflow: hidden;
    }

    .module-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.75rem 1rem;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .module-title {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .module-title i {
      font-size: 0.875rem;
      color: #4f46e5;
    }

    .module-name {
      font-size: 0.875rem;
      font-weight: 600;
      color: #0f172a;
      margin: 0;
      text-transform: capitalize;
    }

    .module-count {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
      background: #e2e8f0;
      padding: 0.25rem 0.5rem;
      border-radius: 0.25rem;
    }

    .module-permissions {
      display: flex;
      flex-direction: column;
      padding: 0.5rem;
    }

    /* Permission Items */
    .permission-item {
      padding: 0.5rem;
      border-radius: 0.375rem;
      transition: background-color 0.15s;
    }

    .permission-item:hover {
      background: #f8fafc;
    }

    .permission-item--selected {
      background: #f0fdf4;
    }

    .permission-item--selected:hover {
      background: #dcfce7;
    }

    .permission-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      cursor: pointer;
    }

    .permission-checkbox {
      width: 1rem;
      height: 1rem;
      accent-color: #4f46e5;
      cursor: pointer;
    }

    .permission-readonly {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .permission-indicator {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 1rem;
      height: 1rem;
      border-radius: 50%;
      background: #e2e8f0;
      font-size: 0.5rem;
      color: transparent;
      transition: all 0.15s;
    }

    .permission-indicator--selected {
      background: #16a34a;
      color: #ffffff;
    }

    .permission-name {
      font-size: 0.8125rem;
      color: #334155;
    }

    .permission-description {
      font-size: 0.75rem;
      color: #64748b;
      margin: 0.25rem 0 0 1.5rem;
      line-height: 1.4;
    }

    /* Footer */
    .permissions-footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding: 1rem 1.5rem;
      border-top: 1px solid #e2e8f0;
      background: #f8fafc;
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      font-size: 0.875rem;
      font-weight: 500;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: all 0.15s;
    }

    .btn:focus {
      outline: none;
    }

    .btn:focus-visible {
      outline: 2px solid #4f46e5;
      outline-offset: 2px;
    }

    .btn-primary {
      color: #ffffff;
      background: #4f46e5;
      border: 1px solid #4f46e5;
    }

    .btn-primary:hover:not(:disabled) {
      background: #4338ca;
      border-color: #4338ca;
    }

    .btn-primary:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }

    .btn-secondary {
      color: #334155;
      background: #ffffff;
      border: 1px solid #e2e8f0;
    }

    .btn-secondary:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
    }

    .btn-spinner {
      width: 0.875rem;
      height: 0.875rem;
      border: 2px solid rgba(255, 255, 255, 0.3);
      border-top-color: #ffffff;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    /* Responsive */
    @media (max-width: 640px) {
      .permissions-header {
        flex-direction: column;
        align-items: flex-start;
      }

      .permissions-grid {
        grid-template-columns: 1fr;
      }
    }
  `],
})
export class RolePermissionsComponent implements OnInit {
  private rbacService = inject(RbacService);
  private toastService = inject(ToastService);

  @Input({ required: true }) role!: Role;
  @Input() canEdit = false;

  @Output() permissionsChanged = new EventEmitter<void>();
  @Output() close = new EventEmitter<void>();

  // State
  allPermissions = signal<RbacPermission[]>([]);
  selectedPermissionIds = signal<Set<string>>(new Set());
  originalPermissionIds = signal<Set<string>>(new Set());
  loading = signal(true);
  saving = signal(false);

  // Computed
  modules = computed(() => {
    const moduleSet = new Set<string>();
    this.allPermissions().forEach(p => moduleSet.add(p.module));
    return Array.from(moduleSet).sort();
  });

  selectedCount = computed(() => this.selectedPermissionIds().size);

  hasChanges = computed(() => {
    const current = this.selectedPermissionIds();
    const original = this.originalPermissionIds();

    if (current.size !== original.size) return true;

    for (const id of current) {
      if (!original.has(id)) return true;
    }

    return false;
  });

  ngOnInit(): void {
    this.loadPermissions();
  }

  private loadPermissions(): void {
    this.loading.set(true);

    this.rbacService.getCachedPermissions().subscribe({
      next: permissions => {
        this.allPermissions.set(permissions);

        // Set initial selected permissions from the role
        const selectedIds = new Set(
          (this.role.permissions || []).map(p => p.id)
        );
        this.selectedPermissionIds.set(selectedIds);
        this.originalPermissionIds.set(new Set(selectedIds));

        this.loading.set(false);
      },
      error: err => {
        this.toastService.error('Failed to load permissions');
        this.loading.set(false);
        console.error('Failed to load permissions:', err);
      },
    });
  }

  getModuleIcon(module: string): string {
    const icons: Record<string, string> = {
      users: 'fa-solid fa-users',
      roles: 'fa-solid fa-user-shield',
      students: 'fa-solid fa-user-graduate',
      courses: 'fa-solid fa-book',
      attendance: 'fa-solid fa-calendar-check',
      grades: 'fa-solid fa-star',
      fees: 'fa-solid fa-dollar-sign',
      reports: 'fa-solid fa-chart-bar',
      settings: 'fa-solid fa-cog',
      dashboard: 'fa-solid fa-tachometer-alt',
      admin: 'fa-solid fa-user-tie',
    };
    return icons[module.toLowerCase()] || 'fa-solid fa-folder';
  }

  formatModuleName(module: string): string {
    return module.charAt(0).toUpperCase() + module.slice(1);
  }

  getPermissionsByModule(module: string): RbacPermission[] {
    return this.allPermissions().filter(p => p.module === module);
  }

  getModulePermissionCount(module: string): string {
    const modulePerms = this.getPermissionsByModule(module);
    const selected = modulePerms.filter(p =>
      this.selectedPermissionIds().has(p.id)
    ).length;
    return `${selected}/${modulePerms.length}`;
  }

  isPermissionSelected(permissionId: string): boolean {
    return this.selectedPermissionIds().has(permissionId);
  }

  togglePermission(permission: RbacPermission): void {
    const current = new Set(this.selectedPermissionIds());

    if (current.has(permission.id)) {
      current.delete(permission.id);
    } else {
      current.add(permission.id);
    }

    this.selectedPermissionIds.set(current);
  }

  selectAll(): void {
    const allIds = new Set(this.allPermissions().map(p => p.id));
    this.selectedPermissionIds.set(allIds);
  }

  deselectAll(): void {
    this.selectedPermissionIds.set(new Set());
  }

  resetChanges(): void {
    this.selectedPermissionIds.set(new Set(this.originalPermissionIds()));
  }

  saveChanges(): void {
    if (!this.hasChanges()) return;

    this.saving.set(true);

    const current = this.selectedPermissionIds();
    const original = this.originalPermissionIds();

    // Calculate permissions to add and remove
    const toAdd: string[] = [];
    const toRemove: string[] = [];

    for (const id of current) {
      if (!original.has(id)) toAdd.push(id);
    }

    for (const id of original) {
      if (!current.has(id)) toRemove.push(id);
    }

    // We'll use the set permissions approach for simplicity
    // In a production app, you might want to batch add/remove operations
    const permissionIds = Array.from(current);

    // For now, let's use a workaround: clear all and add new
    // This would be better handled by a setPermissions endpoint
    this.rbacService.getRole(this.role.id).subscribe({
      next: () => {
        // Remove all existing permissions first, then add the new ones
        this.performPermissionUpdates(toAdd, toRemove);
      },
      error: err => {
        this.toastService.error('Failed to save permissions');
        this.saving.set(false);
        console.error('Failed to save permissions:', err);
      },
    });
  }

  private performPermissionUpdates(toAdd: string[], toRemove: string[]): void {
    // Simple sequential approach - remove first, then add
    const removeOps = toRemove.length > 0
      ? this.rbacService.removePermissions(this.role.id, { permissionIds: toRemove })
      : null;

    const addOps = toAdd.length > 0
      ? this.rbacService.assignPermissions(this.role.id, { permissionIds: toAdd })
      : null;

    if (!removeOps && !addOps) {
      this.saving.set(false);
      return;
    }

    // Execute operations
    if (removeOps && addOps) {
      removeOps.subscribe({
        next: () => {
          addOps.subscribe({
            next: () => this.onSaveSuccess(),
            error: err => this.onSaveError(err),
          });
        },
        error: err => this.onSaveError(err),
      });
    } else if (removeOps) {
      removeOps.subscribe({
        next: () => this.onSaveSuccess(),
        error: err => this.onSaveError(err),
      });
    } else if (addOps) {
      addOps.subscribe({
        next: () => this.onSaveSuccess(),
        error: err => this.onSaveError(err),
      });
    }
  }

  private onSaveSuccess(): void {
    this.toastService.success('Permissions updated successfully');
    this.originalPermissionIds.set(new Set(this.selectedPermissionIds()));
    this.saving.set(false);
    this.permissionsChanged.emit();
  }

  private onSaveError(err: unknown): void {
    this.toastService.error('Failed to save permissions');
    this.saving.set(false);
    console.error('Failed to save permissions:', err);
  }

  onClose(): void {
    this.close.emit();
  }
}
