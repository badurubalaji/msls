-- Seed Data for Testing Stories 8.2 (Examinations) and 8.3 (Hall Tickets)
-- Run: PGPASSWORD=postgres psql -h localhost -U postgres -d msls -f scripts/seed_exam_data.sql

-- Variables (update if your tenant/data differs)
-- Tenant: Demo School = 61ef9fd2-2e9e-4b70-9f16-3b6ea73d4fa4
-- Academic Year: 2025-26 = 5384ce1c-d1c0-4a22-a1f9-7dd43fe8d353
-- Exam Type: Mid Term = 019c0ce0-0551-7d3d-8872-50e37e724e6a
-- Class: LKG = 019c0948-5e5f-7993-be66-360deafd167d

-- Set tenant context for RLS
SET app.tenant_id = '61ef9fd2-2e9e-4b70-9f16-3b6ea73d4fa4';

-- ============================================================
-- 1. Create a Hall Ticket Template
-- ============================================================
INSERT INTO hall_ticket_templates (
    id, tenant_id, name, school_name, school_address, instructions, is_default
) VALUES (
    '019d1000-0001-7000-8000-000000000001',
    '61ef9fd2-2e9e-4b70-9f16-3b6ea73d4fa4',
    'Default Template',
    'Demo School',
    '123 Education Lane, Knowledge City',
    '1. Bring this hall ticket to the examination hall.
2. No electronic devices allowed.
3. Report 30 minutes before the exam.
4. Carry valid ID proof.',
    true
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 2. Create an Examination (Status: scheduled)
-- ============================================================
INSERT INTO examinations (
    id, tenant_id, name, exam_type_id, academic_year_id,
    start_date, end_date, status, description
) VALUES (
    '019d1000-0002-7000-8000-000000000001',
    '61ef9fd2-2e9e-4b70-9f16-3b6ea73d4fa4',
    'Mid Term Examination - LKG',
    '019c0ce0-0551-7d3d-8872-50e37e724e6a',
    '5384ce1c-d1c0-4a22-a1f9-7dd43fe8d353',
    '2026-02-15',
    '2026-02-20',
    'scheduled',
    'Mid term examination for LKG students'
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 3. Link Examination to Class
-- ============================================================
INSERT INTO examination_classes (examination_id, class_id)
VALUES (
    '019d1000-0002-7000-8000-000000000001',
    '019c0948-5e5f-7993-be66-360deafd167d'
) ON CONFLICT DO NOTHING;

-- ============================================================
-- 4. Create Exam Schedules (Subject-wise)
-- ============================================================
-- English
INSERT INTO exam_schedules (
    id, examination_id, subject_id, exam_date, start_time, end_time, max_marks, passing_marks, venue
) VALUES (
    '019d1000-0003-7000-8000-000000000001',
    '019d1000-0002-7000-8000-000000000001',
    '019c095a-1505-7353-94a3-3c4e3bddce2d',
    '2026-02-15',
    '09:00',
    '11:00',
    100,
    35,
    'Hall A'
) ON CONFLICT (id) DO NOTHING;

-- Hindi
INSERT INTO exam_schedules (
    id, examination_id, subject_id, exam_date, start_time, end_time, max_marks, passing_marks, venue
) VALUES (
    '019d1000-0003-7000-8000-000000000002',
    '019d1000-0002-7000-8000-000000000001',
    '019c095a-1506-72e7-9287-bc90ecd190a0',
    '2026-02-16',
    '09:00',
    '11:00',
    100,
    35,
    'Hall A'
) ON CONFLICT (id) DO NOTHING;

-- Mathematics
INSERT INTO exam_schedules (
    id, examination_id, subject_id, exam_date, start_time, end_time, max_marks, passing_marks, venue
) VALUES (
    '019d1000-0003-7000-8000-000000000003',
    '019d1000-0002-7000-8000-000000000001',
    '019c095a-1506-7c92-ac79-134fadb03013',
    '2026-02-17',
    '09:00',
    '12:00',
    100,
    35,
    'Hall B'
) ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Summary
-- ============================================================
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '===========================================';
    RAISE NOTICE 'SEED DATA CREATED SUCCESSFULLY';
    RAISE NOTICE '===========================================';
    RAISE NOTICE 'Examination: Mid Term Examination - LKG';
    RAISE NOTICE 'Exam ID: 019d1000-0002-7000-8000-000000000001';
    RAISE NOTICE 'Status: scheduled';
    RAISE NOTICE 'Dates: Feb 15-20, 2026';
    RAISE NOTICE 'Class: LKG';
    RAISE NOTICE 'Subjects: English, Hindi, Mathematics';
    RAISE NOTICE '';
    RAISE NOTICE 'Hall Ticket Template: Default Template';
    RAISE NOTICE '';
    RAISE NOTICE 'Now you can:';
    RAISE NOTICE '1. View exam at /exams/list';
    RAISE NOTICE '2. View schedules at /exams/{id}/schedules';
    RAISE NOTICE '3. Generate hall tickets at /exams/{id}/hall-tickets';
    RAISE NOTICE '4. Manage templates at /exams/hall-ticket-templates';
    RAISE NOTICE '===========================================';
END $$;
