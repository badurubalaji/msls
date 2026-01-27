# Epic 2: User Authentication & Access Control

**Phase:** 1 (MVP)
**Priority:** Critical - Required for all user access

## Epic Goal

Enable users to securely register, login, and access the system based on their roles with complete RBAC implementation.

## User Value

Super Admins can create tenants, Admins can manage users and roles, all users can login securely with appropriate access levels.

## FRs Covered

FR-CF-04, FR-CF-05, FR-CF-07, FR-CF-10, ARCH-AUTH

---

## Stories

### Story 2.1: User Registration & Email Verification

As a **school administrator**,
I want **to register new staff accounts with email verification**,
So that **only verified users can access the system**.

**Acceptance Criteria:**

**Given** an admin is on the user management page
**When** they create a new user with email, name, and role
**Then** the user record is created with status "pending_verification"
**And** a verification email is sent with a secure token (valid 24 hours)
**And** the email contains a link to set password

**Given** a user clicks the verification link
**When** the token is valid
**Then** they are shown a password setup form
**And** password must meet requirements (min 8 chars, upper, lower, number, special)
**And** after setting password, account is activated
**And** user is redirected to login page

**Given** a user clicks an expired or invalid verification link
**When** they try to access the page
**Then** they see an error message
**And** option to request new verification email

---

### Story 2.2: User Login with Email/Password

As a **user**,
I want **to login with my email and password**,
So that **I can access the system based on my role**.

**Acceptance Criteria:**

**Given** a user is on the login page
**When** they enter valid credentials
**Then** they receive a JWT access token (15 min expiry)
**And** they receive a refresh token (7 day expiry, httpOnly cookie)
**And** they are redirected to their role-appropriate dashboard
**And** audit log records the login event

**Given** a user enters invalid credentials
**When** they submit the login form
**Then** they see a generic "Invalid credentials" error (no indication of which field)
**And** failed attempt is logged

**Given** a user fails login 5 times within 15 minutes
**When** they try again
**Then** account is temporarily locked for 30 minutes
**And** email notification is sent about suspicious activity

---

### Story 2.3: JWT Token Refresh & Session Management

As a **logged-in user**,
I want **my session to stay active without frequent re-login**,
So that **I can work uninterrupted**.

**Acceptance Criteria:**

**Given** a user has a valid refresh token
**When** their access token expires
**Then** the frontend automatically requests a new access token
**And** the new access token is used for subsequent requests
**And** user experience is uninterrupted

**Given** a user's refresh token is expired
**When** an API request is made
**Then** user is redirected to login page
**And** previous route is saved for redirect after login

**Given** a user clicks logout
**When** logout is processed
**Then** refresh token is invalidated on server
**And** all tokens are removed from client
**And** user is redirected to login page

---

### Story 2.4: Two-Factor Authentication (TOTP)

As an **administrator**,
I want **to enable 2FA for my account**,
So that **my account has additional security protection**.

**Acceptance Criteria:**

**Given** an admin navigates to security settings
**When** they enable 2FA
**Then** they see a QR code for authenticator app (Google Authenticator, Authy)
**And** they must enter a valid TOTP code to confirm setup
**And** backup codes (10) are generated and displayed once
**And** 2FA is enabled on their account

**Given** 2FA is enabled on an account
**When** the user logs in with correct password
**Then** they are prompted for TOTP code
**And** only after valid code entry are they fully authenticated

**Given** a user has lost their authenticator
**When** they enter a backup code
**Then** they can access their account
**And** backup code is marked as used
**And** they are prompted to set up 2FA again

---

### Story 2.5: Role & Permission Management

As a **super admin**,
I want **to create and manage roles with specific permissions**,
So that **users only have access to features they need**.

**Acceptance Criteria:**

**Given** a super admin is on role management
**When** they create a new role
**Then** they can select from available permissions grouped by module
**And** permissions follow format: module:action (e.g., students:read, students:write)
**And** role is saved and available for assignment

**Given** default roles exist
**When** viewing the system
**Then** the following roles are available:
- Super Admin (all permissions)
- Admin (tenant-level management)
- Principal (read all, manage academics)
- Teacher (manage assigned classes/subjects)
- Accountant (fee management)
- Librarian (library management)
- Parent (view child data)
- Student (view own data)

**Given** a user attempts an action
**When** they lack the required permission
**Then** they receive a 403 Forbidden response
**And** frontend shows appropriate access denied message

---

### Story 2.6: User Profile Management

As a **user**,
I want **to view and update my profile information**,
So that **my account details are accurate**.

**Acceptance Criteria:**

**Given** a user navigates to their profile
**When** viewing the profile page
**Then** they see their name, email, phone, avatar
**And** they see their assigned roles and last login time

**Given** a user wants to update their profile
**When** they edit allowed fields (name, phone, avatar)
**Then** changes are saved successfully
**And** confirmation message is shown
**And** audit log records the update

**Given** a user wants to change their password
**When** they enter current password and new password
**Then** current password is verified
**And** new password meets complexity requirements
**And** password is updated and all sessions are invalidated
**And** user must login again with new password

---

### Story 2.7: Parent/Student OTP Login

As a **parent or student**,
I want **to login using my phone number and OTP**,
So that **I can access the portal without remembering complex passwords**.

**Acceptance Criteria:**

**Given** a parent/student is on the login page
**When** they enter their registered phone number
**Then** an OTP is sent via SMS (6 digits, valid 5 minutes)
**And** they are shown an OTP entry screen

**Given** a valid OTP is entered
**When** submitted within validity period
**Then** user is authenticated and receives tokens
**And** they are redirected to parent/student portal

**Given** OTP is entered incorrectly 3 times
**When** they try again
**Then** OTP is invalidated
**And** they must request a new OTP
**And** there is a 1-minute cooldown between requests

---

### Story 2.8: Feature Flags & Configuration Engine

As an **administrator**,
I want **to enable/disable features per tenant/branch**,
So that **I can control feature rollout and customization**.

**Acceptance Criteria:**

**Given** an admin navigates to configuration settings
**When** they view feature flags
**Then** they see a list of toggleable features
**And** each feature shows current status (enabled/disabled)
**And** features can be enabled at tenant or branch level

**Given** a feature is disabled
**When** a user tries to access that feature
**Then** the menu item is hidden
**And** direct URL access shows "Feature not available" message

**Given** custom fields are configured for an entity
**When** that entity form is displayed
**Then** custom fields appear in designated section
**And** custom field values are saved with the entity
**And** custom fields appear in list views as configured
