# 13 - Communication System

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Communication System handles all school-to-parent, school-to-staff, and internal communications including notices, announcements, SMS, email, and push notifications.

---

## 2. Notice & Circular Management

### 2.1 Notice Entity

**Entity: Notice**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch (NULL = all) |
| title | VARCHAR(255) | Notice title |
| content | TEXT | Notice content (HTML/Markdown) |
| notice_type | ENUM | circular, announcement, event, holiday, exam, emergency |
| priority | ENUM | low, normal, high, urgent |
| target_audience | JSONB | Who should receive |
| attachments | JSONB | Attached files |
| publish_date | DATE | When to publish |
| expiry_date | DATE | When to hide |
| requires_acknowledgement | BOOLEAN | Needs confirmation |
| status | ENUM | draft, scheduled, published, archived |
| published_by | UUID | Publisher |
| published_at | TIMESTAMP | Publish time |
| created_by | UUID | Creator |
| created_at | TIMESTAMP | Creation time |

### 2.2 Target Audience Configuration

```json
{
  "target_type": "specific",
  "recipients": {
    "all_parents": false,
    "all_staff": false,
    "classes": ["5-A", "5-B", "6-A"],
    "sections": [],
    "specific_students": [],
    "specific_staff": [],
    "roles": ["teacher", "class_teacher"]
  }
}
```

### 2.3 Notice Acknowledgement

**Entity: NoticeAcknowledgement**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| notice_id | UUID | Notice reference |
| user_id | UUID | User who acknowledged |
| user_type | ENUM | parent, staff, student |
| acknowledged_at | TIMESTAMP | When acknowledged |
| response | TEXT | Optional response |

### 2.4 Notice Types

| Type | Description | Priority |
|------|-------------|----------|
| circular | General school circulars | Normal |
| announcement | Important announcements | Normal-High |
| event | Event invitations | Normal |
| holiday | Holiday notifications | Normal |
| exam | Exam schedules/results | High |
| emergency | Emergency alerts | Urgent |

---

## 3. SMS Integration

### 3.1 SMS Configuration

**Entity: SMSConfig**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| provider | ENUM | msg91, textlocal, twilio, custom |
| api_key | VARCHAR(255) | API key (encrypted) |
| sender_id | VARCHAR(20) | Sender ID |
| is_active | BOOLEAN | Active status |
| daily_limit | INT | Daily SMS limit |
| monthly_limit | INT | Monthly limit |

**Supported Providers**:
- MSG91 (India)
- TextLocal (India)
- Twilio (Global)
- Custom HTTP API

### 3.2 SMS Template

**Entity: SMSTemplate**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Template name |
| code | VARCHAR(50) | Template code |
| content | TEXT | Message template |
| variables | JSONB | Available variables |
| dlt_template_id | VARCHAR(50) | DLT ID (India) |
| is_active | BOOLEAN | Active status |

**Example Templates**:
```yaml
templates:
  - code: ATTENDANCE_ABSENT
    content: "Dear Parent, {student_name} of class {class} was absent on {date}. Regards, {school_name}"
    variables: [student_name, class, date, school_name]

  - code: FEE_REMINDER
    content: "Dear Parent, Fee of Rs.{amount} for {student_name} is due on {due_date}. Pay now to avoid late fee. {school_name}"
    variables: [amount, student_name, due_date, school_name]

  - code: EXAM_RESULT
    content: "Dear Parent, Results for {exam_name} are now available. {student_name} scored {percentage}%. View details in parent app. {school_name}"
    variables: [exam_name, student_name, percentage, school_name]
```

### 3.3 SMS Log

**Entity: SMSLog**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| template_id | UUID | Template used |
| recipient_phone | VARCHAR(20) | Phone number |
| recipient_user_id | UUID | User reference |
| message | TEXT | Actual message sent |
| status | ENUM | queued, sent, delivered, failed |
| provider_message_id | VARCHAR(100) | Provider reference |
| sent_at | TIMESTAMP | When sent |
| delivered_at | TIMESTAMP | When delivered |
| error_message | TEXT | If failed |
| cost | DECIMAL | SMS cost |
| triggered_by | VARCHAR(50) | What triggered |

### 3.4 SMS Usage Dashboard

```
SMS Usage Report | January 2026

| Category        | Sent   | Delivered | Failed | Cost    |
|-----------------|--------|-----------|--------|---------|
| Attendance      | 1,250  | 1,220     | 30     | ₹625    |
| Fee Reminders   | 450    | 445       | 5      | ₹225    |
| Circulars       | 3,200  | 3,150     | 50     | ₹1,600  |
| Emergency       | 15     | 15        | 0      | ₹7.50   |
| Others          | 180    | 175       | 5      | ₹90     |
|-----------------|--------|-----------|--------|---------|
| Total           | 5,095  | 5,005     | 90     | ₹2,547  |

Delivery Rate: 98.2%
Monthly Limit: 10,000 | Used: 5,095 (50.95%)
```

---

## 4. Email Integration

### 4.1 Email Configuration

**Entity: EmailConfig**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| provider | ENUM | smtp, sendgrid, ses, mailgun |
| from_email | VARCHAR(255) | From address |
| from_name | VARCHAR(100) | From name |
| smtp_host | VARCHAR(255) | SMTP host |
| smtp_port | INT | SMTP port |
| smtp_username | VARCHAR(255) | Username |
| smtp_password | VARCHAR(255) | Password (encrypted) |
| api_key | VARCHAR(255) | API key (for cloud) |
| is_active | BOOLEAN | Active status |

### 4.2 Email Template

**Entity: EmailTemplate**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Template name |
| code | VARCHAR(50) | Template code |
| subject | VARCHAR(255) | Email subject |
| body_html | TEXT | HTML body |
| body_text | TEXT | Plain text body |
| variables | JSONB | Available variables |
| is_active | BOOLEAN | Active status |

### 4.3 Email Log

**Entity: EmailLog**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| template_id | UUID | Template used |
| recipient_email | VARCHAR(255) | Email address |
| recipient_user_id | UUID | User reference |
| subject | VARCHAR(255) | Actual subject |
| status | ENUM | queued, sent, delivered, bounced, complained |
| provider_message_id | VARCHAR(100) | Provider reference |
| sent_at | TIMESTAMP | When sent |
| opened_at | TIMESTAMP | When opened |
| clicked_at | TIMESTAMP | When link clicked |
| error_message | TEXT | If failed |

---

## 5. Push Notifications

### 5.1 Device Registration

**Entity: UserDevice**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | User reference |
| device_token | VARCHAR(500) | FCM/APNs token |
| platform | ENUM | android, ios, web |
| device_name | VARCHAR(100) | Device name |
| app_version | VARCHAR(20) | App version |
| is_active | BOOLEAN | Active status |
| last_used_at | TIMESTAMP | Last activity |
| registered_at | TIMESTAMP | Registration time |

### 5.2 Push Notification

**Entity: PushNotification**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| title | VARCHAR(200) | Notification title |
| body | TEXT | Notification body |
| data | JSONB | Custom data payload |
| target_type | ENUM | all, segment, user |
| target_users | JSONB | Target user list |
| status | ENUM | draft, scheduled, sent |
| scheduled_at | TIMESTAMP | When to send |
| sent_at | TIMESTAMP | When sent |
| sent_count | INT | Successful sends |
| failed_count | INT | Failed sends |

### 5.3 Notification Triggers

| Trigger | Notification | Priority |
|---------|--------------|----------|
| Student marked absent | "Your child was absent today" | High |
| Fee due reminder | "Fee payment due in 3 days" | Normal |
| Fee received | "Payment of ₹X received" | Normal |
| Exam result published | "Results are now available" | High |
| New circular | "New circular from school" | Normal |
| Homework assigned | "New homework in {subject}" | Normal |
| Emergency alert | "URGENT: {message}" | Urgent |

---

## 6. In-App Messaging

### 6.1 Conversation Thread

**Entity: Conversation**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| type | ENUM | parent_teacher, group, support |
| title | VARCHAR(200) | Conversation title |
| participants | JSONB | User IDs |
| student_id | UUID | Related student (if P-T chat) |
| is_active | BOOLEAN | Active status |
| created_at | TIMESTAMP | Creation time |
| last_message_at | TIMESTAMP | Last activity |

### 6.2 Message

**Entity: Message**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| conversation_id | UUID | Conversation reference |
| sender_id | UUID | Sender user |
| content | TEXT | Message content |
| message_type | ENUM | text, image, file, system |
| attachments | JSONB | Attached files |
| is_read | BOOLEAN | Read status |
| read_at | TIMESTAMP | When read |
| created_at | TIMESTAMP | Sent time |

### 6.3 Chat Rules

| Rule | Description |
|------|-------------|
| Office Hours | Chat available 8 AM - 6 PM only |
| Response Time | Teachers should respond within 24 hours |
| Moderation | Admin can view all conversations |
| Archive | Auto-archive after 30 days of inactivity |
| Block | Teachers can mute notifications from specific parents |

---

## 7. Communication Workflow

### 7.1 Automated Notifications

```yaml
automated_triggers:
  - event: student_absent
    channels: [sms, push, email]
    template: ATTENDANCE_ABSENT
    timing: immediate
    recipients: primary_guardian

  - event: fee_due_3_days
    channels: [sms, push]
    template: FEE_REMINDER
    timing: 9:00 AM
    recipients: all_guardians

  - event: exam_result_published
    channels: [push, email]
    template: EXAM_RESULT
    timing: immediate
    recipients: all_guardians

  - event: circular_published
    channels: [push]
    template: NEW_CIRCULAR
    timing: immediate
    recipients: target_audience
```

### 7.2 Broadcast Flow

```
Create Notice
     │
     ▼
Select Audience
     │
     ▼
Choose Channels (SMS/Email/Push/In-App)
     │
     ▼
Schedule or Send Now
     │
     ▼
Track Delivery Status
     │
     ▼
Monitor Acknowledgements
```

---

## 8. WhatsApp Integration (Optional)

### 8.1 WhatsApp Business API

**Entity: WhatsAppConfig**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| provider | ENUM | twilio, gupshup, interakt |
| api_key | VARCHAR(255) | API key |
| phone_number_id | VARCHAR(50) | Business phone ID |
| is_active | BOOLEAN | Active status |

### 8.2 WhatsApp Templates

Pre-approved message templates:
- Fee reminder
- Attendance alert
- Exam schedule
- Holiday notice
- Emergency alert

---

## 9. API Endpoints

```
# Notices
GET    /api/v1/notices                      # List notices
POST   /api/v1/notices                      # Create notice
GET    /api/v1/notices/{id}                 # Get notice
PUT    /api/v1/notices/{id}                 # Update notice
POST   /api/v1/notices/{id}/publish         # Publish notice
GET    /api/v1/notices/{id}/acknowledgements # Get acks

# SMS
POST   /api/v1/sms/send                     # Send SMS
POST   /api/v1/sms/send-bulk                # Bulk SMS
GET    /api/v1/sms/logs                     # SMS logs
GET    /api/v1/sms/usage                    # Usage stats

# Email
POST   /api/v1/email/send                   # Send email
POST   /api/v1/email/send-bulk              # Bulk email
GET    /api/v1/email/logs                   # Email logs

# Push Notifications
POST   /api/v1/push/send                    # Send push
GET    /api/v1/push/logs                    # Push logs

# In-App Chat
GET    /api/v1/conversations                # List conversations
POST   /api/v1/conversations                # Start conversation
GET    /api/v1/conversations/{id}/messages  # Get messages
POST   /api/v1/conversations/{id}/messages  # Send message
```

---

## 10. Business Rules

| Rule | Description |
|------|-------------|
| SMS Timing | No SMS between 9 PM - 8 AM |
| Emergency Override | Emergency alerts can be sent anytime |
| Consent Required | Parents must opt-in for promotional messages |
| Template Approval | All SMS templates need DLT approval (India) |
| Character Limit | SMS limited to 160 chars (or 2 SMS max) |
| Attachment Size | Email attachments max 10 MB |
| Chat Response | Auto-reply if teacher doesn't respond in 24h |

---

## 11. Related Documents

- [14-parent-portal.md](./14-parent-portal.md) - Parent access
- [15-student-portal.md](./15-student-portal.md) - Student access
- [index.md](./index.md) - Main PRD index

---

**Previous**: [12-fees-payments.md](./12-fees-payments.md)
**Next**: [14-parent-portal.md](./14-parent-portal.md)
