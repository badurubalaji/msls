/**
 * MSLS Storage Service
 *
 * Provides a type-safe wrapper around localStorage for storing and retrieving
 * authentication tokens and other persistent data.
 */

import { Injectable } from '@angular/core';

/** Storage keys used throughout the application */
export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'msls_access_token',
  REFRESH_TOKEN: 'msls_refresh_token',
  TOKEN_EXPIRY: 'msls_token_expiry',
  CURRENT_USER: 'msls_current_user',
  TENANT_ID: 'msls_tenant_id',
  REMEMBER_ME: 'msls_remember_me',
} as const;

/** Type for storage key values */
export type StorageKey = (typeof STORAGE_KEYS)[keyof typeof STORAGE_KEYS];

/**
 * StorageService - Type-safe localStorage wrapper for token management.
 *
 * Usage:
 * constructor(private storageService: StorageService) {}
 *
 * saveToken() {
 *   this.storageService.setItem(STORAGE_KEYS.ACCESS_TOKEN, 'my-token');
 * }
 *
 * getToken() {
 *   return this.storageService.getItem(STORAGE_KEYS.ACCESS_TOKEN);
 * }
 */
@Injectable({ providedIn: 'root' })
export class StorageService {
  /**
   * Check if localStorage is available
   */
  private isLocalStorageAvailable(): boolean {
    try {
      const testKey = '__storage_test__';
      window.localStorage.setItem(testKey, testKey);
      window.localStorage.removeItem(testKey);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Set an item in localStorage
   * @param key - Storage key
   * @param value - Value to store (will be JSON serialized if object)
   */
  setItem<T>(key: StorageKey | string, value: T): void {
    if (!this.isLocalStorageAvailable()) {
      console.warn('localStorage is not available');
      return;
    }

    try {
      const serializedValue = typeof value === 'string' ? value : JSON.stringify(value);
      window.localStorage.setItem(key, serializedValue);
    } catch (error) {
      console.error(`Error saving to localStorage: ${key}`, error);
    }
  }

  /**
   * Get an item from localStorage
   * @param key - Storage key
   * @returns The stored value or null if not found
   */
  getItem<T = string>(key: StorageKey | string): T | null {
    if (!this.isLocalStorageAvailable()) {
      return null;
    }

    try {
      const item = window.localStorage.getItem(key);
      if (item === null) {
        return null;
      }

      // Try to parse as JSON, return as string if parsing fails
      try {
        return JSON.parse(item) as T;
      } catch {
        return item as T;
      }
    } catch (error) {
      console.error(`Error reading from localStorage: ${key}`, error);
      return null;
    }
  }

  /**
   * Remove an item from localStorage
   * @param key - Storage key
   */
  removeItem(key: StorageKey | string): void {
    if (!this.isLocalStorageAvailable()) {
      return;
    }

    try {
      window.localStorage.removeItem(key);
    } catch (error) {
      console.error(`Error removing from localStorage: ${key}`, error);
    }
  }

  /**
   * Clear all items from localStorage
   */
  clear(): void {
    if (!this.isLocalStorageAvailable()) {
      return;
    }

    try {
      window.localStorage.clear();
    } catch (error) {
      console.error('Error clearing localStorage', error);
    }
  }

  /**
   * Check if a key exists in localStorage
   * @param key - Storage key
   */
  hasItem(key: StorageKey | string): boolean {
    return this.getItem(key) !== null;
  }

  /**
   * Get all MSLS-related keys from localStorage
   */
  getMslsKeys(): string[] {
    if (!this.isLocalStorageAvailable()) {
      return [];
    }

    const keys: string[] = [];
    for (let i = 0; i < window.localStorage.length; i++) {
      const key = window.localStorage.key(i);
      if (key?.startsWith('msls_')) {
        keys.push(key);
      }
    }
    return keys;
  }

  /**
   * Clear only MSLS-related items from localStorage
   */
  clearMslsData(): void {
    const keys = this.getMslsKeys();
    keys.forEach(key => this.removeItem(key));
  }
}
