/**
 * MSLS Branch Service
 *
 * HTTP service for branch management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import { Branch, CreateBranchRequest, UpdateBranchRequest } from './branch.model';

/** Response format for branch list from backend */
interface BranchListResponse {
  branches: Branch[];
  total: number;
}

/**
 * BranchService - Handles all branch-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class BranchService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/branches';

  /**
   * Get all branches for the current tenant
   */
  getBranches(): Observable<Branch[]> {
    return this.apiService.get<BranchListResponse>(this.basePath).pipe(
      map(response => response.branches || [])
    );
  }

  /**
   * Get a single branch by ID
   */
  getBranch(id: string): Observable<Branch> {
    return this.apiService.get<Branch>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new branch
   */
  createBranch(data: CreateBranchRequest): Observable<Branch> {
    return this.apiService.post<Branch>(this.basePath, data);
  }

  /**
   * Update an existing branch
   */
  updateBranch(id: string, data: UpdateBranchRequest): Observable<Branch> {
    return this.apiService.put<Branch>(`${this.basePath}/${id}`, data);
  }

  /**
   * Set a branch as primary
   */
  setPrimary(id: string): Observable<Branch> {
    return this.apiService.patch<Branch>(`${this.basePath}/${id}/primary`, {});
  }

  /**
   * Set branch active status
   * @param id Branch ID
   * @param isActive New active status
   */
  setStatus(id: string, isActive: boolean): Observable<Branch> {
    return this.apiService.patch<Branch>(`${this.basePath}/${id}/status`, { isActive });
  }

  /**
   * Delete a branch
   */
  deleteBranch(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }
}
