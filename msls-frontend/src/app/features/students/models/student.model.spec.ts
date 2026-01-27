/**
 * Student Model Tests
 */

import { describe, it, expect } from 'vitest';
import {
  StudentAddress,
  getStatusBadgeVariant,
  getStatusLabel,
  getGenderLabel,
  calculateAge,
  formatAddress,
} from './student.model';

describe('Student Model Helpers', () => {
  describe('getStatusBadgeVariant', () => {
    it('should return success for active status', () => {
      expect(getStatusBadgeVariant('active')).toBe('success');
    });

    it('should return neutral for inactive status', () => {
      expect(getStatusBadgeVariant('inactive')).toBe('neutral');
    });

    it('should return warning for transferred status', () => {
      expect(getStatusBadgeVariant('transferred')).toBe('warning');
    });

    it('should return info for graduated status', () => {
      expect(getStatusBadgeVariant('graduated')).toBe('info');
    });
  });

  describe('getStatusLabel', () => {
    it('should return proper labels for each status', () => {
      expect(getStatusLabel('active')).toBe('Active');
      expect(getStatusLabel('inactive')).toBe('Inactive');
      expect(getStatusLabel('transferred')).toBe('Transferred');
      expect(getStatusLabel('graduated')).toBe('Graduated');
    });
  });

  describe('getGenderLabel', () => {
    it('should return proper labels for each gender', () => {
      expect(getGenderLabel('male')).toBe('Male');
      expect(getGenderLabel('female')).toBe('Female');
      expect(getGenderLabel('other')).toBe('Other');
    });
  });

  describe('calculateAge', () => {
    it('should calculate age correctly', () => {
      // Use a date that is definitely more than 10 years ago
      const tenYearsAgo = new Date();
      tenYearsAgo.setFullYear(tenYearsAgo.getFullYear() - 10);
      tenYearsAgo.setMonth(0); // January
      tenYearsAgo.setDate(1);

      const age = calculateAge(tenYearsAgo.toISOString().split('T')[0]);
      expect(age).toBe(10);
    });

    it('should handle birthday not yet occurred this year', () => {
      const nextMonth = new Date();
      nextMonth.setFullYear(nextMonth.getFullYear() - 10);
      nextMonth.setMonth(nextMonth.getMonth() + 1);

      const age = calculateAge(nextMonth.toISOString().split('T')[0]);
      expect(age).toBe(9);
    });
  });

  describe('formatAddress', () => {
    it('should format address correctly', () => {
      const address: StudentAddress = {
        id: '1',
        addressType: 'current',
        addressLine1: '123 Main St',
        addressLine2: 'Apt 4B',
        city: 'Mumbai',
        state: 'Maharashtra',
        postalCode: '400001',
        country: 'India',
      };

      const formatted = formatAddress(address);
      expect(formatted).toBe('123 Main St, Apt 4B, Mumbai, Maharashtra, 400001, India');
    });

    it('should handle missing address line 2', () => {
      const address: StudentAddress = {
        id: '1',
        addressType: 'current',
        addressLine1: '123 Main St',
        city: 'Mumbai',
        state: 'Maharashtra',
        postalCode: '400001',
        country: 'India',
      };

      const formatted = formatAddress(address);
      expect(formatted).toBe('123 Main St, Mumbai, Maharashtra, 400001, India');
    });

    it('should return empty string for undefined address', () => {
      expect(formatAddress(undefined)).toBe('');
    });
  });
});
