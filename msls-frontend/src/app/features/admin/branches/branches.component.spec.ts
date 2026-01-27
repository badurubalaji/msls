/**
 * MSLS Branches Component Tests
 *
 * Unit tests for the BranchesComponent.
 */

import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';

import { BranchesComponent } from './branches.component';
import { BranchService } from './branch.service';
import { ToastService } from '../../../shared/services';
import { Branch } from './branch.model';

describe('BranchesComponent', () => {
  let component: BranchesComponent;
  let fixture: ComponentFixture<BranchesComponent>;
  let branchServiceMock: {
    getBranches: ReturnType<typeof vi.fn>;
    createBranch: ReturnType<typeof vi.fn>;
    updateBranch: ReturnType<typeof vi.fn>;
    deleteBranch: ReturnType<typeof vi.fn>;
    setPrimary: ReturnType<typeof vi.fn>;
    setStatus: ReturnType<typeof vi.fn>;
  };
  let toastServiceMock: {
    success: ReturnType<typeof vi.fn>;
    error: ReturnType<typeof vi.fn>;
  };

  const mockBranches: Branch[] = [
    {
      id: '550e8400-e29b-41d4-a716-446655440001',
      code: 'MAIN',
      name: 'Main Campus',
      addressLine1: '123 Main Street',
      city: 'Mumbai',
      state: 'Maharashtra',
      country: 'India',
      timezone: 'Asia/Kolkata',
      isPrimary: true,
      isActive: true,
      createdAt: '2026-01-23T10:00:00Z',
      updatedAt: '2026-01-23T10:00:00Z',
    },
    {
      id: '550e8400-e29b-41d4-a716-446655440002',
      code: 'BR01',
      name: 'North Branch',
      city: 'Delhi',
      state: 'Delhi',
      country: 'India',
      timezone: 'Asia/Kolkata',
      isPrimary: false,
      isActive: true,
      createdAt: '2026-01-23T11:00:00Z',
      updatedAt: '2026-01-23T11:00:00Z',
    },
    {
      id: '550e8400-e29b-41d4-a716-446655440003',
      code: 'BR02',
      name: 'South Branch',
      city: 'Chennai',
      state: 'Tamil Nadu',
      country: 'India',
      timezone: 'Asia/Kolkata',
      isPrimary: false,
      isActive: false,
      createdAt: '2026-01-23T12:00:00Z',
      updatedAt: '2026-01-23T12:00:00Z',
    },
  ];

  beforeEach(async () => {
    branchServiceMock = {
      getBranches: vi.fn(),
      createBranch: vi.fn(),
      updateBranch: vi.fn(),
      deleteBranch: vi.fn(),
      setPrimary: vi.fn(),
      setStatus: vi.fn(),
    };
    toastServiceMock = {
      success: vi.fn(),
      error: vi.fn(),
    };

    // Default to returning mock branches
    branchServiceMock.getBranches.mockReturnValue(of(mockBranches));

    await TestBed.configureTestingModule({
      imports: [BranchesComponent],
      providers: [
        { provide: BranchService, useValue: branchServiceMock },
        { provide: ToastService, useValue: toastServiceMock },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(BranchesComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('initialization', () => {
    it('should load branches on init', fakeAsync(() => {
      fixture.detectChanges();
      tick();

      expect(branchServiceMock.getBranches).toHaveBeenCalled();
      expect(component.branches()).toEqual(mockBranches);
      expect(component.loading()).toBeFalsy();
    }));

    it('should show loading state initially', () => {
      expect(component.loading()).toBeTruthy();
    });

    it('should handle error when loading branches fails', fakeAsync(() => {
      branchServiceMock.getBranches.mockReturnValue(
        throwError(() => new Error('Network error'))
      );

      fixture.detectChanges();
      tick();

      expect(component.error()).toBe('Failed to load branches. Please try again.');
      expect(component.loading()).toBeFalsy();
    }));
  });

  describe('search functionality', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should filter branches by name', () => {
      component.onSearchChange('Main');

      expect(component.filteredBranches().length).toBe(1);
      expect(component.filteredBranches()[0].name).toBe('Main Campus');
    });

    it('should filter branches by code', () => {
      component.onSearchChange('BR01');

      expect(component.filteredBranches().length).toBe(1);
      expect(component.filteredBranches()[0].code).toBe('BR01');
    });

    it('should filter branches by city', () => {
      component.onSearchChange('Chennai');

      expect(component.filteredBranches().length).toBe(1);
      expect(component.filteredBranches()[0].city).toBe('Chennai');
    });

    it('should return all branches when search term is empty', () => {
      component.onSearchChange('');

      expect(component.filteredBranches().length).toBe(3);
    });

    it('should return empty array when no branches match', () => {
      component.onSearchChange('nonexistent');

      expect(component.filteredBranches().length).toBe(0);
    });

    it('should be case-insensitive', () => {
      component.onSearchChange('main');

      expect(component.filteredBranches().length).toBe(1);
      expect(component.filteredBranches()[0].name).toBe('Main Campus');
    });
  });

  describe('create branch modal', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should open create modal', () => {
      component.openCreateModal();

      expect(component.showBranchModal()).toBeTruthy();
      expect(component.editingBranch()).toBeNull();
    });

    it('should close create modal', () => {
      component.openCreateModal();
      component.closeBranchModal();

      expect(component.showBranchModal()).toBeFalsy();
    });

    it('should create branch successfully', fakeAsync(() => {
      const newBranch: Branch = {
        ...mockBranches[0],
        id: 'new-branch-id',
        code: 'NEW01',
        name: 'New Branch',
      };
      branchServiceMock.createBranch.mockReturnValue(of(newBranch));

      component.openCreateModal();
      component.saveBranch({ code: 'NEW01', name: 'New Branch' });
      tick();

      expect(branchServiceMock.createBranch).toHaveBeenCalledWith({
        code: 'NEW01',
        name: 'New Branch',
      });
      expect(toastServiceMock.success).toHaveBeenCalledWith('Branch created successfully');
      expect(component.showBranchModal()).toBeFalsy();
    }));

    it('should handle error when creating branch fails', fakeAsync(() => {
      branchServiceMock.createBranch.mockReturnValue(
        throwError(() => new Error('Create failed'))
      );

      component.openCreateModal();
      component.saveBranch({ code: 'NEW01', name: 'New Branch' });
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to create branch');
      expect(component.saving()).toBeFalsy();
    }));
  });

  describe('edit branch modal', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should open edit modal with branch data', () => {
      component.editBranch(mockBranches[0]);

      expect(component.showBranchModal()).toBeTruthy();
      expect(component.editingBranch()).toEqual(mockBranches[0]);
    });

    it('should update branch successfully', fakeAsync(() => {
      const updatedBranch: Branch = {
        ...mockBranches[0],
        name: 'Updated Main Campus',
      };
      branchServiceMock.updateBranch.mockReturnValue(of(updatedBranch));

      component.editBranch(mockBranches[0]);
      component.saveBranch({ code: 'MAIN', name: 'Updated Main Campus' });
      tick();

      expect(branchServiceMock.updateBranch).toHaveBeenCalledWith(
        mockBranches[0].id,
        { code: 'MAIN', name: 'Updated Main Campus' }
      );
      expect(toastServiceMock.success).toHaveBeenCalledWith('Branch updated successfully');
    }));

    it('should handle error when updating branch fails', fakeAsync(() => {
      branchServiceMock.updateBranch.mockReturnValue(
        throwError(() => new Error('Update failed'))
      );

      component.editBranch(mockBranches[0]);
      component.saveBranch({ code: 'MAIN', name: 'Updated Main Campus' });
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to update branch');
    }));
  });

  describe('toggle status', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should deactivate an active branch successfully', fakeAsync(() => {
      const activeBranch = mockBranches[0]; // isActive: true
      const deactivatedBranch: Branch = { ...activeBranch, isActive: false };
      branchServiceMock.setStatus.mockReturnValue(of(deactivatedBranch));

      component.toggleStatus(activeBranch);
      tick();

      expect(branchServiceMock.setStatus).toHaveBeenCalledWith(activeBranch.id, false);
      expect(toastServiceMock.success).toHaveBeenCalledWith('Branch deactivated successfully');
    }));

    it('should activate an inactive branch successfully', fakeAsync(() => {
      const inactiveBranch = mockBranches[2]; // isActive: false
      const activatedBranch: Branch = { ...inactiveBranch, isActive: true };
      branchServiceMock.setStatus.mockReturnValue(of(activatedBranch));

      component.toggleStatus(inactiveBranch);
      tick();

      expect(branchServiceMock.setStatus).toHaveBeenCalledWith(inactiveBranch.id, true);
      expect(toastServiceMock.success).toHaveBeenCalledWith('Branch activated successfully');
    }));

    it('should handle error when toggling status fails', fakeAsync(() => {
      branchServiceMock.setStatus.mockReturnValue(
        throwError(() => new Error('Toggle failed'))
      );

      component.toggleStatus(mockBranches[0]);
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to update branch status');
    }));
  });

  describe('set primary', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should open set primary confirmation modal', () => {
      component.setPrimary(mockBranches[1]);

      expect(component.showPrimaryModal()).toBeTruthy();
      expect(component.branchToSetPrimary()).toEqual(mockBranches[1]);
    });

    it('should close set primary modal', () => {
      component.setPrimary(mockBranches[1]);
      component.closePrimaryModal();

      expect(component.showPrimaryModal()).toBeFalsy();
      expect(component.branchToSetPrimary()).toBeNull();
    });

    it('should set branch as primary successfully', fakeAsync(() => {
      const updatedBranch: Branch = { ...mockBranches[1], isPrimary: true };
      branchServiceMock.setPrimary.mockReturnValue(of(updatedBranch));

      component.setPrimary(mockBranches[1]);
      component.confirmSetPrimary();
      tick();

      expect(branchServiceMock.setPrimary).toHaveBeenCalledWith(mockBranches[1].id);
      expect(toastServiceMock.success).toHaveBeenCalledWith('Primary branch updated successfully');
      expect(component.showPrimaryModal()).toBeFalsy();
    }));

    it('should handle error when setting primary fails', fakeAsync(() => {
      branchServiceMock.setPrimary.mockReturnValue(
        throwError(() => new Error('Set primary failed'))
      );

      component.setPrimary(mockBranches[1]);
      component.confirmSetPrimary();
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to set primary branch');
      expect(component.settingPrimary()).toBeFalsy();
    }));
  });

  describe('delete branch', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should open delete confirmation modal', () => {
      component.confirmDelete(mockBranches[1]);

      expect(component.showDeleteModal()).toBeTruthy();
      expect(component.branchToDelete()).toEqual(mockBranches[1]);
    });

    it('should close delete modal', () => {
      component.confirmDelete(mockBranches[1]);
      component.closeDeleteModal();

      expect(component.showDeleteModal()).toBeFalsy();
      expect(component.branchToDelete()).toBeNull();
    });

    it('should delete branch successfully', fakeAsync(() => {
      branchServiceMock.deleteBranch.mockReturnValue(of(void 0));

      component.confirmDelete(mockBranches[1]);
      component.deleteBranch();
      tick();

      expect(branchServiceMock.deleteBranch).toHaveBeenCalledWith(mockBranches[1].id);
      expect(toastServiceMock.success).toHaveBeenCalledWith('Branch deleted successfully');
      expect(component.showDeleteModal()).toBeFalsy();
    }));

    it('should handle error when deleting branch fails', fakeAsync(() => {
      branchServiceMock.deleteBranch.mockReturnValue(
        throwError(() => new Error('Delete failed'))
      );

      component.confirmDelete(mockBranches[1]);
      component.deleteBranch();
      tick();

      expect(toastServiceMock.error).toHaveBeenCalledWith('Failed to delete branch');
      expect(component.deleting()).toBeFalsy();
    }));

    it('should not delete if no branch is selected', fakeAsync(() => {
      component.deleteBranch();
      tick();

      expect(branchServiceMock.deleteBranch).not.toHaveBeenCalled();
    }));
  });

  describe('UI rendering', () => {
    beforeEach(fakeAsync(() => {
      fixture.detectChanges();
      tick();
    }));

    it('should render branch table with correct data', () => {
      fixture.detectChanges();

      const rows = fixture.nativeElement.querySelectorAll('tbody tr');
      expect(rows.length).toBe(3);

      const firstRowCells = rows[0].querySelectorAll('td');
      expect(firstRowCells[0].textContent).toContain('Main Campus');
      expect(firstRowCells[1].textContent).toContain('MAIN');
    });

    it('should show primary badge for primary branch', () => {
      fixture.detectChanges();

      const primaryBadge = fixture.nativeElement.querySelector('.badge-blue');
      expect(primaryBadge).toBeTruthy();
      expect(primaryBadge.textContent).toContain('Primary');
    });

    it('should show active status badge', () => {
      fixture.detectChanges();

      const activeBadges = fixture.nativeElement.querySelectorAll('.badge-green');
      expect(activeBadges.length).toBe(2); // Two active branches
    });

    it('should show inactive status badge', () => {
      fixture.detectChanges();

      const inactiveBadge = fixture.nativeElement.querySelector('.badge-gray');
      expect(inactiveBadge).toBeTruthy();
      expect(inactiveBadge.textContent).toContain('Inactive');
    });

    it('should not show delete button for primary branch', () => {
      fixture.detectChanges();

      const rows = fixture.nativeElement.querySelectorAll('tbody tr');
      const primaryRow = rows[0]; // First branch is primary
      const deleteButton = primaryRow.querySelector('.action-btn--danger');

      expect(deleteButton).toBeFalsy();
    });

    it('should not show set primary button for already primary branch', () => {
      fixture.detectChanges();

      const rows = fixture.nativeElement.querySelectorAll('tbody tr');
      const primaryRow = rows[0];
      const setPrimaryButton = primaryRow.querySelector('.action-btn--primary');

      expect(setPrimaryButton).toBeFalsy();
    });
  });

  describe('empty state', () => {
    it('should show empty state when no branches exist', fakeAsync(() => {
      branchServiceMock.getBranches.mockReturnValue(of([]));

      fixture.detectChanges();
      tick();
      fixture.detectChanges();

      const emptyState = fixture.nativeElement.querySelector('.empty-state');
      expect(emptyState).toBeTruthy();
      expect(emptyState.textContent).toContain('No branches found');
    }));

    it('should show empty state when search has no results', fakeAsync(() => {
      fixture.detectChanges();
      tick();

      component.onSearchChange('nonexistent');
      fixture.detectChanges();

      const emptyState = fixture.nativeElement.querySelector('.empty-state');
      expect(emptyState).toBeTruthy();
    }));
  });
});
