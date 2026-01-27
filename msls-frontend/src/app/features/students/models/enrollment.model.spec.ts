/**
 * Enrollment Model Tests
 */

import { describe, it, expect } from 'vitest';
import {
  getEnrollmentStatusBadgeVariant,
  getEnrollmentStatusLabel,
  getEnrollmentTimelineColor,
} from './enrollment.model';

describe('Enrollment Model Helpers', () => {
  describe('getEnrollmentStatusBadgeVariant', () => {
    it('should return success for active status', () => {
      expect(getEnrollmentStatusBadgeVariant('active')).toBe('success');
    });

    it('should return neutral for completed status', () => {
      expect(getEnrollmentStatusBadgeVariant('completed')).toBe('neutral');
    });

    it('should return warning for transferred status', () => {
      expect(getEnrollmentStatusBadgeVariant('transferred')).toBe('warning');
    });

    it('should return error for dropout status', () => {
      expect(getEnrollmentStatusBadgeVariant('dropout')).toBe('error');
    });
  });

  describe('getEnrollmentStatusLabel', () => {
    it('should return proper labels for each status', () => {
      expect(getEnrollmentStatusLabel('active')).toBe('Active');
      expect(getEnrollmentStatusLabel('completed')).toBe('Completed');
      expect(getEnrollmentStatusLabel('transferred')).toBe('Transferred');
      expect(getEnrollmentStatusLabel('dropout')).toBe('Dropout');
    });

    it('should return raw status for unknown status', () => {
      expect(getEnrollmentStatusLabel('unknown' as any)).toBe('unknown');
    });
  });

  describe('getEnrollmentTimelineColor', () => {
    it('should return green for active status', () => {
      expect(getEnrollmentTimelineColor('active')).toBe('bg-green-500');
    });

    it('should return gray for completed status', () => {
      expect(getEnrollmentTimelineColor('completed')).toBe('bg-gray-400');
    });

    it('should return orange for transferred status', () => {
      expect(getEnrollmentTimelineColor('transferred')).toBe('bg-orange-500');
    });

    it('should return red for dropout status', () => {
      expect(getEnrollmentTimelineColor('dropout')).toBe('bg-red-500');
    });

    it('should return gray for unknown status', () => {
      expect(getEnrollmentTimelineColor('unknown' as any)).toBe('bg-gray-400');
    });
  });
});
