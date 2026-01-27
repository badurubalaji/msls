# 09 - Digital Classroom

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 2

---

## 1. Overview

The Digital Classroom module provides class recording capabilities (screen + audio), content management, knowledge replay, and a digital content library for enhanced learning experiences.

---

## 2. Class Recording System

### 2.1 Recording Session

**Entity: ClassRecording**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| class_id | UUID | Class reference |
| section_id | UUID | Section reference |
| subject_id | UUID | Subject reference |
| topic_title | VARCHAR(255) | Topic/lesson title |
| description | TEXT | Session description |
| teacher_id | UUID | Recording teacher |
| recording_date | DATE | Recording date |
| start_time | TIME | Start time |
| end_time | TIME | End time |
| duration_seconds | INT | Total duration |
| status | ENUM | recording, processing, ready, failed, archived |
| visibility | ENUM | class_only, branch, all_classes |
| created_at | TIMESTAMP | Creation time |

### 2.2 Recording Assets

**Entity: RecordingAsset**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| recording_id | UUID | Parent recording |
| asset_type | ENUM | screen_video, audio, thumbnail, transcript |
| file_url | VARCHAR(500) | Storage URL |
| file_size_bytes | BIGINT | File size |
| duration_seconds | INT | Asset duration |
| mime_type | VARCHAR(100) | File type |
| quality | VARCHAR(20) | 720p, 1080p, etc. |
| processing_status | ENUM | pending, processing, ready, failed |
| created_at | TIMESTAMP | Creation time |

### 2.3 Recording Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    RECORDING FLOW                            │
└─────────────────────────────────────────────────────────────┘

1. Teacher Initiates Recording
   ├── Select Class/Section/Subject
   ├── Enter Topic Title
   └── Start Screen + Audio Capture

2. During Recording
   ├── Screen capture (1080p recommended)
   ├── Audio capture (teacher's microphone)
   ├── Real-time preview
   ├── Pause/Resume capability
   └── Add timestamps/bookmarks

3. End Recording
   ├── Stop capture
   ├── Upload to server
   └── Start processing

4. Server Processing
   ├── Video encoding (multiple qualities)
   ├── Audio normalization
   ├── Generate thumbnail
   ├── Optional: Generate transcript
   └── Mark as ready

5. Available for Students
   └── Students can replay anytime
```

### 2.4 Recording Interface (Teacher)

```
┌─────────────────────────────────────────────────────────────┐
│  NEW CLASS RECORDING                                         │
├─────────────────────────────────────────────────────────────┤
│  Class: [10-A          ▼]   Subject: [Mathematics   ▼]      │
│  Topic: [Quadratic Equations - Completing the Square    ]   │
│  Description: [                                          ]   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                                                      │    │
│  │              SCREEN PREVIEW                          │    │
│  │                                                      │    │
│  │         (Shows what will be recorded)               │    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  🎤 Microphone: [Built-in Mic ▼]  Level: ████████░░ OK     │
│                                                              │
│  [🔴 START RECORDING]              Duration: 00:00:00       │
│                                                              │
│  Options:                                                    │
│  ☑ Record full screen  ☐ Record window only                 │
│  ☑ Include audio       ☐ Generate transcript                │
└─────────────────────────────────────────────────────────────┘
```

### 2.5 Recording During Class

```
┌─────────────────────────────────────────────────────────────┐
│  🔴 RECORDING IN PROGRESS                     ⏱ 00:23:45    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Topic: Quadratic Equations - Completing the Square         │
│  Class: 10-A | Subject: Mathematics                         │
│                                                              │
│  [⏸ PAUSE]  [🔖 ADD BOOKMARK]  [⏹ STOP RECORDING]          │
│                                                              │
│  Bookmarks:                                                  │
│  • 00:05:23 - Introduction                                  │
│  • 00:12:45 - Example 1                                     │
│  • 00:18:30 - Practice problem                              │
│                                                              │
│  Audio Level: ████████░░░ Good                              │
│  Storage: 245 MB used                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Content Bookmarks & Chapters

### 3.1 Bookmark Entity

**Entity: RecordingBookmark**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| recording_id | UUID | Recording reference |
| timestamp_seconds | INT | Position in video |
| title | VARCHAR(200) | Bookmark title |
| description | TEXT | Optional description |
| bookmark_type | ENUM | chapter, important, question, example |
| created_by | UUID | Creator |
| created_at | TIMESTAMP | Creation time |

### 3.2 Auto-Generated Chapters

```
Recording: Quadratic Equations | Duration: 45:30

Chapters (Auto-detected):
00:00 - Introduction
05:23 - What is a Quadratic Equation?
12:45 - Standard Form
18:30 - Solving by Factoring
25:15 - Completing the Square Method
32:00 - Quadratic Formula
40:00 - Practice Problems
44:00 - Summary
```

---

## 4. Student Playback

### 4.1 Playback Interface

```
┌─────────────────────────────────────────────────────────────┐
│  Quadratic Equations - Completing the Square                 │
│  Teacher: Mr. Rajesh Kumar | Date: 22-Jan-2026              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                                                      │    │
│  │                   VIDEO PLAYER                       │    │
│  │                                                      │    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  ▶ ━━━━━━━━━●━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 23:45 / 45:30  │
│                                                              │
│  [⏮] [⏪ 10s] [▶ Play] [⏩ 10s] [⏭]  🔊 ━━━●━━ [⚙]        │
│                                                              │
│  Playback Speed: [1x ▼]                                     │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│  📑 CHAPTERS                          📝 MY NOTES           │
│  ├ 00:00 Introduction                 ┌──────────────────┐  │
│  ├ 05:23 What is Quadratic Eq?       │ Remember:        │  │
│  ├ 12:45 Standard Form               │ a≠0 is important │  │
│  ├ 18:30 Solving by Factoring        │                  │  │
│  ► 25:15 Completing Square ◀         │ [Add Note]       │  │
│  ├ 32:00 Quadratic Formula           └──────────────────┘  │
│  └ 40:00 Practice Problems                                  │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Student Notes

**Entity: StudentNote**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| recording_id | UUID | Recording reference |
| timestamp_seconds | INT | Video position |
| note_content | TEXT | Note text |
| created_at | TIMESTAMP | Creation time |
| updated_at | TIMESTAMP | Last update |

### 4.3 Watch Progress

**Entity: WatchProgress**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| recording_id | UUID | Recording reference |
| last_position_seconds | INT | Last watched position |
| total_watched_seconds | INT | Total time watched |
| completion_percentage | DECIMAL | % completed |
| watch_count | INT | Number of views |
| first_watched_at | TIMESTAMP | First view |
| last_watched_at | TIMESTAMP | Last view |

---

## 5. Content Library

### 5.1 Content Entity

**Entity: DigitalContent**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| title | VARCHAR(255) | Content title |
| description | TEXT | Description |
| content_type | ENUM | video, pdf, document, presentation, link, interactive |
| subject_id | UUID | Subject reference |
| class_ids | JSONB | Applicable classes |
| topic_id | UUID | Topic reference |
| file_url | VARCHAR(500) | File/link URL |
| thumbnail_url | VARCHAR(500) | Thumbnail |
| file_size_bytes | BIGINT | File size |
| duration_seconds | INT | For video/audio |
| is_downloadable | BOOLEAN | Allow download |
| visibility | ENUM | public, class_restricted, teacher_only |
| uploaded_by | UUID | Uploader |
| uploaded_at | TIMESTAMP | Upload time |
| view_count | INT | Total views |
| is_featured | BOOLEAN | Featured content |

### 5.2 Content Categories

```
Digital Library
├── Class 10
│   ├── Mathematics
│   │   ├── Videos (25)
│   │   ├── PDFs (12)
│   │   ├── Practice Sheets (8)
│   │   └── External Links (5)
│   ├── Science
│   │   ├── Videos (30)
│   │   ├── Lab Manuals (10)
│   │   └── Simulations (15)
│   └── ...
├── Class 9
│   └── ...
└── Shared Resources
    ├── Career Guidance
    ├── Soft Skills
    └── General Knowledge
```

### 5.3 Content Upload Flow

```
Teacher uploads content
        │
        ▼
Select content type
        │
        ▼
Fill metadata (title, subject, class, topic)
        │
        ▼
Upload file / Enter URL
        │
        ▼
Processing (for videos: encoding, thumbnail)
        │
        ▼
Set visibility & permissions
        │
        ▼
Publish to library
```

---

## 6. Explained PDFs

### 6.1 PDF with Audio Explanation

**Entity: ExplainedPDF**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| content_id | UUID | Base content reference |
| pdf_url | VARCHAR(500) | PDF file URL |
| total_pages | INT | Number of pages |
| created_by | UUID | Creator |
| created_at | TIMESTAMP | Creation time |

**Entity: PDFPageExplanation**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| explained_pdf_id | UUID | Parent reference |
| page_number | INT | Page number |
| audio_url | VARCHAR(500) | Audio explanation URL |
| audio_duration_seconds | INT | Audio length |
| annotations | JSONB | Page annotations |
| created_at | TIMESTAMP | Creation time |

### 6.2 Explained PDF Viewer

```
┌─────────────────────────────────────────────────────────────┐
│  📄 Chapter 5: Quadratic Equations                          │
│  Teacher: Mr. Rajesh Kumar                                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                                                      │    │
│  │                    PDF PAGE                          │    │
│  │                   (Page 3/12)                        │    │
│  │                                                      │    │
│  │    [Annotations/highlights shown on page]           │    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  [◀ Prev] Page 3 of 12 [Next ▶]                            │
│                                                              │
│  🔊 Audio Explanation for this page:                        │
│  ▶ ━━━━━━━━━●━━━━━━━━━━━━━━━━━━━ 01:23 / 02:45             │
│  [▶ Play] [⏩ Next Page Audio]                              │
│                                                              │
│  📥 Download PDF                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Live Class Integration (Future)

### 7.1 Virtual Classroom

**Entity: LiveClass**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| title | VARCHAR(255) | Class title |
| class_id | UUID | Target class |
| section_id | UUID | Target section |
| subject_id | UUID | Subject |
| teacher_id | UUID | Host teacher |
| scheduled_at | TIMESTAMP | Scheduled time |
| duration_minutes | INT | Expected duration |
| meeting_url | VARCHAR(500) | Meeting link |
| meeting_provider | ENUM | jitsi, bbb, zoom, meet |
| recording_enabled | BOOLEAN | Record session |
| status | ENUM | scheduled, live, ended, cancelled |
| attendance_count | INT | Students joined |
| recording_id | UUID | Saved recording |

### 7.2 Integration Options

| Provider | Type | Features |
|----------|------|----------|
| Jitsi Meet | Self-hosted | Free, open-source, full control |
| BigBlueButton | Self-hosted | Education-focused, whiteboard |
| Zoom | Cloud | Reliable, but paid |
| Google Meet | Cloud | Easy integration |

---

## 8. Analytics

### 8.1 Content Engagement

```
Content Analytics | January 2026

Most Viewed Recordings:
| # | Title                          | Subject | Views | Avg Watch |
|---|--------------------------------|---------|-------|-----------|
| 1 | Quadratic Equations            | Math    | 245   | 85%       |
| 2 | Chemical Reactions             | Science | 198   | 78%       |
| 3 | Shakespeare's Macbeth          | English | 156   | 72%       |

Student Engagement:
- Total recordings available: 150
- Total watch time: 1,250 hours
- Average completion rate: 68%
- Students who watched: 420/450 (93%)
```

### 8.2 Teacher Dashboard

```
My Recordings | Mr. Rajesh Kumar

Total Recordings: 25
Total Duration: 18 hours 30 min
Total Views: 1,250
Average Rating: 4.5/5

Recent:
| Recording                    | Date    | Views | Completion |
|------------------------------|---------|-------|------------|
| Trigonometry - Heights       | 20-Jan  | 35    | 82%        |
| Quadratic Formula            | 18-Jan  | 42    | 88%        |
| Linear Equations Review      | 15-Jan  | 38    | 75%        |
```

---

## 9. Storage Management

### 9.1 Storage Quotas

| Plan | Storage Limit | Retention |
|------|---------------|-----------|
| Free | 10 GB | 3 months |
| Standard | 100 GB | 1 year |
| Premium | 500 GB | 2 years |
| Enterprise | Unlimited | Custom |

### 9.2 Storage Optimization

- Automatic quality selection based on bandwidth
- Compression for older recordings
- Archive to cold storage after semester
- Delete policy for expired content

---

## 10. API Endpoints

```
# Recordings
POST   /api/v1/recordings/start             # Start recording
POST   /api/v1/recordings/{id}/stop         # Stop recording
GET    /api/v1/recordings                   # List recordings
GET    /api/v1/recordings/{id}              # Get recording
DELETE /api/v1/recordings/{id}              # Delete recording
POST   /api/v1/recordings/{id}/bookmark     # Add bookmark

# Playback
GET    /api/v1/recordings/{id}/play         # Get playback URL
POST   /api/v1/recordings/{id}/progress     # Update watch progress
GET    /api/v1/recordings/{id}/notes        # Get student notes
POST   /api/v1/recordings/{id}/notes        # Add note

# Content Library
GET    /api/v1/content                      # List content
POST   /api/v1/content                      # Upload content
GET    /api/v1/content/{id}                 # Get content
DELETE /api/v1/content/{id}                 # Delete content

# Explained PDFs
POST   /api/v1/explained-pdf                # Create explained PDF
POST   /api/v1/explained-pdf/{id}/page      # Add page explanation
GET    /api/v1/explained-pdf/{id}           # Get explained PDF

# Analytics
GET    /api/v1/recordings/analytics         # Recording analytics
GET    /api/v1/content/analytics            # Content analytics
```

---

## 11. Business Rules

| Rule | Description |
|------|-------------|
| Recording Consent | Display recording notice at class start |
| Storage Quota | Alert at 80%, block at 100% |
| Download Rights | Teacher can control downloadability |
| Access Control | Content visible only to assigned classes |
| Auto-Archive | Archive recordings older than 1 year |
| Quality Options | Provide 360p, 720p, 1080p options |

---

## 12. Related Documents

- [04-academic-operations.md](./04-academic-operations.md) - Timetable
- [08-online-quiz-assessment.md](./08-online-quiz-assessment.md) - Assessments
- [07-homework-assignments.md](./07-homework-assignments.md) - Homework
- [index.md](./index.md) - Main PRD index

---

**Previous**: [08-online-quiz-assessment.md](./08-online-quiz-assessment.md)
**Next**: [10-staff-management.md](./10-staff-management.md)
