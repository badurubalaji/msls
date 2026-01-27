/**
 * MSLS Shared Services - Barrel Export
 *
 * This file exports all shared services from the shared/services directory.
 * Import from '@shared/services' in your feature modules.
 */

// Toast Service
export { ToastService } from './toast.service';
export type { Toast, ToastVariant, ToastConfig } from './toast.service';

// Modal Service
// Note: ModalSize is exported from components/modal to avoid duplicate exports
export { ModalService } from './modal.service';
export type { ModalConfig, ModalRef } from './modal.service';
