# Epic 12: Communication Hub

**Phase:** 1 (MVP)
**Priority:** High - Essential for school-parent communication

## Epic Goal

Enable multi-channel communication between school and stakeholders.

## User Value

Admins can send notices, parents receive SMS/email/push notifications, and teachers can message parents directly.

## FRs Covered

FR-CM-01 to FR-CM-07

---

## Stories

### Story 12.1: Notice Creation & Publishing

As an **administrator**,
I want **to create and publish notices**,
So that **important information reaches stakeholders**.

**Acceptance Criteria:**

**Given** admin is creating a notice
**When** filling notice details
**Then** they can enter: title, content (rich text editor)
**And** they can set: priority (normal, important, urgent)
**And** they can attach: files (PDF, images)
**And** they can set: publish date (immediate or scheduled)

**Given** target audience selection
**When** choosing recipients
**Then** they can select: all, by class, by section
**And** they can select: parents, students, or both
**And** they can select: specific categories (transport users, etc.)
**And** recipient count is shown

**Given** notice is published
**When** distributed
**Then** notice appears in portal for recipients
**And** notification is sent via configured channels
**And** publish timestamp is recorded

---

### Story 12.2: SMS Notification System

As an **administrator**,
I want **to send SMS notifications**,
So that **urgent messages reach parents quickly**.

**Acceptance Criteria:**

**Given** SMS provider is configured
**When** setting up integration
**Then** they can configure: MSG91, TextLocal, or Twilio
**And** they can enter: API credentials
**And** they can set: sender ID
**And** they can configure: DLT template IDs (for India)

**Given** SMS is being sent
**When** composing message
**Then** they can use: predefined templates
**And** they can use: variables (student name, class, etc.)
**And** character count is shown
**And** preview before sending

**Given** bulk SMS is sent
**When** tracking delivery
**Then** delivery status is tracked (sent, delivered, failed)
**And** failure reasons are logged
**And** SMS credits usage is shown
**And** delivery report is available

---

### Story 12.3: Email Notification System

As an **administrator**,
I want **to send email notifications**,
So that **detailed communications can be sent**.

**Acceptance Criteria:**

**Given** email is configured
**When** setting up
**Then** SMTP settings can be configured
**And** or SendGrid/SES integration available
**And** from address and reply-to configured
**And** test email can be sent

**Given** email is being composed
**When** creating email
**Then** they can use: rich text editor
**And** they can use: email templates
**And** they can attach: files (within size limit)
**And** they can add: inline images

**Given** bulk email is sent
**When** tracking
**Then** sent/delivered/opened/bounced stats shown
**And** unsubscribe handling available
**And** email history is maintained

---

### Story 12.4: Push Notification System

As an **administrator**,
I want **to send push notifications**,
So that **mobile app users get instant alerts**.

**Acceptance Criteria:**

**Given** push notification service is configured
**When** setting up
**Then** Firebase Cloud Messaging configured
**And** notification icons and sounds set
**And** test notification can be sent

**Given** push notification is being sent
**When** composing
**Then** they can enter: title (short)
**And** they can enter: body (notification content)
**And** they can add: deep link (open specific screen)
**And** they can target: specific user segments

**Given** notification is delivered
**When** user taps
**Then** app opens to relevant screen
**And** notification action is logged

---

### Story 12.5: Notice Acknowledgement Tracking

As an **administrator**,
I want **to track notice acknowledgements**,
So that **I know who has read important notices**.

**Acceptance Criteria:**

**Given** a notice requires acknowledgement
**When** creating notice
**Then** they can enable: acknowledgement required
**And** they can set: acknowledgement deadline

**Given** recipient views the notice
**When** acknowledging
**Then** they must click: "I have read this"
**And** timestamp is recorded
**And** acknowledgement is linked to parent account

**Given** tracking acknowledgements
**When** viewing report
**Then** they see: acknowledged count, pending count
**And** they can see: list of non-acknowledgers
**And** reminder can be sent to pending
**And** export available

---

### Story 12.6: Parent-Teacher Messaging

As a **parent or teacher**,
I want **to message each other directly**,
So that **communication is easy and documented**.

**Acceptance Criteria:**

**Given** a parent wants to message teacher
**When** initiating conversation
**Then** they select: teacher from their child's teachers
**And** they can type: message
**And** they can attach: files (assignment queries, etc.)
**And** message is sent

**Given** a teacher receives message
**When** viewing inbox
**Then** they see: list of conversations
**And** they see: unread count
**And** they can reply to messages
**And** conversation history is maintained

**Given** messaging policies
**When** enforced
**Then** teachers can set: available hours
**And** auto-reply for out-of-hours
**And** block/report option for inappropriate messages

---

### Story 12.7: Automated Event Notifications

As a **system administrator**,
I want **automated notifications for key events**,
So that **stakeholders are informed without manual effort**.

**Acceptance Criteria:**

**Given** automation rules are configured
**When** setting up triggers
**Then** they can configure notifications for:
- Student marked absent
- Fee payment received
- Assignment published
- Exam result published
- Leave approved/rejected
**And** channel preference per event type

**Given** an event occurs
**When** trigger fires
**Then** notification is automatically sent
**And** correct template is used
**And** recipient is determined by event context
**And** notification is logged

**Given** notification templates exist
**When** managing templates
**Then** they can edit: message content
**And** they can use: variables specific to event
**And** they can enable/disable: specific notifications
**And** preview shows sample output

---

### Story 12.8: Communication History & Logs

As an **administrator**,
I want **to view communication history**,
So that **all communications are documented**.

**Acceptance Criteria:**

**Given** admin is on communication logs
**When** viewing history
**Then** they see: all communications (notices, SMS, email)
**And** they can filter: by type, date, recipient
**And** they can search: by content

**Given** individual communication is selected
**When** viewing details
**Then** they see: full content, recipients, delivery status
**And** they see: acknowledgements (if applicable)
**And** they see: sent timestamp and sender

**Given** analytics are needed
**When** viewing reports
**Then** they see: volume by channel (SMS, email, push)
**And** they see: delivery success rate
**And** they see: cost (SMS credits used)
