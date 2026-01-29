/**
 * MSLS Stream Service
 *
 * HTTP service for stream management API calls.
 */

import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  Stream,
  StreamListResponse,
  CreateStreamRequest,
  UpdateStreamRequest,
} from '../academic.model';

/**
 * StreamService - Handles all stream-related API operations.
 */
@Injectable({ providedIn: 'root' })
export class StreamService {
  private readonly apiService = inject(ApiService);
  private readonly basePath = '/streams';

  /**
   * Get all streams
   */
  getStreams(isActive?: boolean, search?: string): Observable<Stream[]> {
    const params: Record<string, string> = {};
    if (isActive !== undefined) params['is_active'] = String(isActive);
    if (search) params['search'] = search;

    return this.apiService.get<StreamListResponse>(this.basePath, { params }).pipe(
      map(response => response.streams || [])
    );
  }

  /**
   * Get streams with total count
   */
  getStreamsWithTotal(isActive?: boolean, search?: string): Observable<StreamListResponse> {
    const params: Record<string, string> = {};
    if (isActive !== undefined) params['is_active'] = String(isActive);
    if (search) params['search'] = search;

    return this.apiService.get<StreamListResponse>(this.basePath, { params });
  }

  /**
   * Get a single stream by ID
   */
  getStream(id: string): Observable<Stream> {
    return this.apiService.get<Stream>(`${this.basePath}/${id}`);
  }

  /**
   * Create a new stream
   */
  createStream(data: CreateStreamRequest): Observable<Stream> {
    return this.apiService.post<Stream>(this.basePath, data);
  }

  /**
   * Update an existing stream
   */
  updateStream(id: string, data: UpdateStreamRequest): Observable<Stream> {
    return this.apiService.put<Stream>(`${this.basePath}/${id}`, data);
  }

  /**
   * Delete a stream
   */
  deleteStream(id: string): Observable<void> {
    return this.apiService.delete<void>(`${this.basePath}/${id}`);
  }
}
