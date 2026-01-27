/**
 * MSLS Branch Service Tests
 *
 * Unit tests for the BranchService.
 */

import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { of, throwError } from 'rxjs';
import { vi } from 'vitest';

import { ApiService } from '../../../core/services/api.service';
import { BranchService } from './branch.service';
import { Branch, CreateBranchRequest, UpdateBranchRequest } from './branch.model';

describe('BranchService', () => {
  let service: BranchService;
  let apiServiceMock: {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    patch: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
  };

  const mockBranch: Branch = {
    id: '550e8400-e29b-41d4-a716-446655440001',
    code: 'MAIN',
    name: 'Main Campus',
    addressLine1: '123 Main Street',
    addressLine2: 'Building A',
    city: 'Mumbai',
    state: 'Maharashtra',
    postalCode: '400001',
    country: 'India',
    phone: '+91 22 1234 5678',
    email: 'main@school.edu',
    logoUrl: undefined,
    timezone: 'Asia/Kolkata',
    isPrimary: true,
    isActive: true,
    createdAt: '2026-01-23T10:00:00Z',
    updatedAt: '2026-01-23T10:00:00Z',
  };

  const mockBranches: Branch[] = [
    mockBranch,
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
  ];

  beforeEach(() => {
    apiServiceMock = {
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      patch: vi.fn(),
      delete: vi.fn(),
    };

    TestBed.configureTestingModule({
      providers: [
        BranchService,
        { provide: ApiService, useValue: apiServiceMock },
      ],
    });

    service = TestBed.inject(BranchService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getBranches', () => {
    it('should fetch all branches', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockBranches));

      let result: Branch[] | undefined;
      service.getBranches().subscribe((branches) => {
        result = branches;
      });
      tick();

      expect(result).toEqual(mockBranches);
      expect(result?.length).toBe(2);
      expect(apiServiceMock.get).toHaveBeenCalledWith('/v1/branches');
    }));

    it('should handle error when fetching branches fails', fakeAsync(() => {
      const error = new Error('Network error');
      apiServiceMock.get.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.getBranches().subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));
  });

  describe('getBranch', () => {
    it('should fetch a single branch by ID', fakeAsync(() => {
      apiServiceMock.get.mockReturnValue(of(mockBranch));

      let result: Branch | undefined;
      service.getBranch(mockBranch.id).subscribe((branch) => {
        result = branch;
      });
      tick();

      expect(result).toEqual(mockBranch);
      expect(apiServiceMock.get).toHaveBeenCalledWith(
        `/v1/branches/${mockBranch.id}`
      );
    }));

    it('should handle error when branch not found', fakeAsync(() => {
      const error = new Error('Branch not found');
      apiServiceMock.get.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.getBranch('non-existent-id').subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));
  });

  describe('createBranch', () => {
    it('should create a new branch', fakeAsync(() => {
      const createRequest: CreateBranchRequest = {
        code: 'NEW01',
        name: 'New Branch',
        city: 'Pune',
        state: 'Maharashtra',
        country: 'India',
        timezone: 'Asia/Kolkata',
        isPrimary: false,
      };

      const createdBranch: Branch = {
        ...mockBranch,
        id: 'new-branch-id',
        code: createRequest.code,
        name: createRequest.name,
        city: createRequest.city,
        isPrimary: false,
      };

      apiServiceMock.post.mockReturnValue(of(createdBranch));

      let result: Branch | undefined;
      service.createBranch(createRequest).subscribe((branch) => {
        result = branch;
      });
      tick();

      expect(result).toEqual(createdBranch);
      expect(apiServiceMock.post).toHaveBeenCalledWith(
        '/v1/branches',
        createRequest
      );
    }));

    it('should handle validation errors on create', fakeAsync(() => {
      const createRequest: CreateBranchRequest = {
        code: '',
        name: '',
      };

      const error = new Error('Validation failed');
      apiServiceMock.post.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.createBranch(createRequest).subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));
  });

  describe('updateBranch', () => {
    it('should update an existing branch', fakeAsync(() => {
      const updateRequest: UpdateBranchRequest = {
        name: 'Updated Main Campus',
        city: 'New Mumbai',
      };

      const updatedBranch: Branch = {
        ...mockBranch,
        name: updateRequest.name!,
        city: updateRequest.city,
      };

      apiServiceMock.put.mockReturnValue(of(updatedBranch));

      let result: Branch | undefined;
      service.updateBranch(mockBranch.id, updateRequest).subscribe((branch) => {
        result = branch;
      });
      tick();

      expect(result?.name).toBe(updateRequest.name);
      expect(result?.city).toBe(updateRequest.city);
      expect(apiServiceMock.put).toHaveBeenCalledWith(
        `/v1/branches/${mockBranch.id}`,
        updateRequest
      );
    }));

    it('should handle error when updating non-existent branch', fakeAsync(() => {
      const error = new Error('Branch not found');
      apiServiceMock.put.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.updateBranch('non-existent-id', { name: 'Test' }).subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));
  });

  describe('setPrimary', () => {
    it('should set a branch as primary', fakeAsync(() => {
      const nonPrimaryBranch = mockBranches[1];
      const updatedBranch: Branch = {
        ...nonPrimaryBranch,
        isPrimary: true,
      };

      apiServiceMock.patch.mockReturnValue(of(updatedBranch));

      let result: Branch | undefined;
      service.setPrimary(nonPrimaryBranch.id).subscribe((branch) => {
        result = branch;
      });
      tick();

      expect(result?.isPrimary).toBe(true);
      expect(apiServiceMock.patch).toHaveBeenCalledWith(
        `/v1/branches/${nonPrimaryBranch.id}/primary`,
        {}
      );
    }));

    it('should handle error when setting primary fails', fakeAsync(() => {
      const error = new Error('Cannot set primary');
      apiServiceMock.patch.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.setPrimary('some-id').subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));
  });

  describe('setStatus', () => {
    it('should set branch status to inactive', fakeAsync(() => {
      const activeBranch = mockBranch;
      const deactivatedBranch: Branch = {
        ...activeBranch,
        isActive: false,
      };

      apiServiceMock.patch.mockReturnValue(of(deactivatedBranch));

      let result: Branch | undefined;
      service.setStatus(activeBranch.id, false).subscribe((branch) => {
        result = branch;
      });
      tick();

      expect(result?.isActive).toBe(false);
      expect(apiServiceMock.patch).toHaveBeenCalledWith(
        `/v1/branches/${activeBranch.id}/status`,
        { isActive: false }
      );
    }));

    it('should set branch status to active', fakeAsync(() => {
      const inactiveBranch: Branch = { ...mockBranch, isActive: false };
      const activatedBranch: Branch = { ...inactiveBranch, isActive: true };

      apiServiceMock.patch.mockReturnValue(of(activatedBranch));

      let result: Branch | undefined;
      service.setStatus(inactiveBranch.id, true).subscribe((branch) => {
        result = branch;
      });
      tick();

      expect(result?.isActive).toBe(true);
      expect(apiServiceMock.patch).toHaveBeenCalledWith(
        `/v1/branches/${inactiveBranch.id}/status`,
        { isActive: true }
      );
    }));

    it('should handle error when setting status fails', fakeAsync(() => {
      const error = new Error('Cannot set status');
      apiServiceMock.patch.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.setStatus('some-id', true).subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));
  });

  describe('deleteBranch', () => {
    it('should delete a branch', fakeAsync(() => {
      const branchToDelete = mockBranches[1];
      apiServiceMock.delete.mockReturnValue(of(void 0));

      let completed = false;
      service.deleteBranch(branchToDelete.id).subscribe({
        next: () => {
          completed = true;
        },
      });
      tick();

      expect(completed).toBe(true);
      expect(apiServiceMock.delete).toHaveBeenCalledWith(
        `/v1/branches/${branchToDelete.id}`
      );
    }));

    it('should handle error when deleting primary branch', fakeAsync(() => {
      const error = new Error('Cannot delete primary branch');
      apiServiceMock.delete.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.deleteBranch(mockBranch.id).subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));

    it('should handle error when branch not found', fakeAsync(() => {
      const error = new Error('Branch not found');
      apiServiceMock.delete.mockReturnValue(throwError(() => error));

      let caughtError: Error | undefined;
      service.deleteBranch('non-existent-id').subscribe({
        next: () => {},
        error: (err) => {
          caughtError = err;
        },
      });
      tick();

      expect(caughtError).toBe(error);
    }));
  });
});
