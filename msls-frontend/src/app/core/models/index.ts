/**
 * MSLS Core Models - Barrel Export
 *
 * This file exports all core models from the core/models directory.
 */

// User Model
export type { User, UserStatus, UserRole, Permission } from './user.model';
export { getUserFullName, hasPermission, hasRole } from './user.model';

// Auth Models
export type {
  AuthTokens,
  LoginRequest,
  RegisterRequest,
  PasswordResetRequest,
  PasswordChangeRequest,
  RefreshTokenRequest,
  AuthResponse,
  TokenVerifyResponse,
} from './auth.model';

// API Response Models
export type {
  PaginationMeta,
  PaginationParams,
  ApiError,
  ApiResponse,
  PaginatedResponse,
} from './api-response.model';
export {
  isSuccessResponse,
  isErrorResponse,
  createSuccessResponse,
  createErrorResponse,
} from './api-response.model';

// Profile Models
export type {
  UserProfile,
  ProfileRole,
  NotificationPreferences,
  UpdateProfileRequest,
  ChangePasswordRequest,
  UpdatePreferencesRequest,
  AvatarUploadResponse,
  UserPreference,
  UserPreferencesResponse,
  SetPreferenceRequest,
} from './profile.model';
export { AVAILABLE_TIMEZONES, AVAILABLE_LOCALES } from './profile.model';

// RBAC Models
export type {
  Permission as RbacPermission,
  Role,
  UserRoles,
  PermissionModules as PermissionModulesResponse,
  CreateRoleRequest,
  UpdateRoleRequest,
  PermissionsRequest,
  RolesRequest,
  SystemRoleName,
  PermissionModule,
  PermissionAction,
} from './rbac.model';
export {
  SystemRoles,
  PermissionModules,
  PermissionActions,
  buildPermissionCode,
  parsePermissionCode,
  isSystemRole,
  canModifyRole,
  canDeleteRole,
} from './rbac.model';
