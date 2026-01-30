-- Attendance Seed Data
-- Creates sample students and attendance records for testing Epic 7 features

DO $$
DECLARE
    v_tenant_id UUID;
    v_branch_id UUID;
    v_user_id UUID;
    v_section_id UUID;
    v_section_id_2 UUID;
    v_academic_year_id UUID;
    v_student_ids UUID[] := ARRAY[]::UUID[];
    v_student_id UUID;
    v_current_date DATE;
    v_attendance_date DATE;
    v_status VARCHAR(20);
    v_random_val FLOAT;
    i INTEGER;
BEGIN
    -- Get existing tenant
    SELECT id INTO v_tenant_id FROM tenants LIMIT 1;
    IF v_tenant_id IS NULL THEN
        RAISE NOTICE 'No tenant found, skipping seed data';
        RETURN;
    END IF;

    -- Get existing branch
    SELECT id INTO v_branch_id FROM branches WHERE tenant_id = v_tenant_id LIMIT 1;

    -- Get existing user (teacher/admin)
    SELECT id INTO v_user_id FROM users WHERE tenant_id = v_tenant_id LIMIT 1;

    -- Get two sections for variety
    SELECT id INTO v_section_id FROM sections WHERE tenant_id = v_tenant_id LIMIT 1;
    SELECT id INTO v_section_id_2 FROM sections WHERE tenant_id = v_tenant_id OFFSET 1 LIMIT 1;

    -- Get academic year
    SELECT id INTO v_academic_year_id FROM academic_years WHERE tenant_id = v_tenant_id LIMIT 1;

    -- Create 10 sample students
    FOR i IN 1..10 LOOP
        -- Check if student already exists
        IF NOT EXISTS (SELECT 1 FROM students WHERE tenant_id = v_tenant_id AND admission_number = 'STU-2026-' || LPAD(i::TEXT, 4, '0')) THEN
            INSERT INTO students (
                id, tenant_id, branch_id, admission_number, first_name, last_name,
                date_of_birth, gender, blood_group, status, admission_date,
                created_at, updated_at
            ) VALUES (
                uuid_generate_v7(),
                v_tenant_id,
                v_branch_id,
                'STU-2026-' || LPAD(i::TEXT, 4, '0'),
                CASE i % 5
                    WHEN 0 THEN 'Rahul'
                    WHEN 1 THEN 'Priya'
                    WHEN 2 THEN 'Amit'
                    WHEN 3 THEN 'Sneha'
                    WHEN 4 THEN 'Vikram'
                END,
                CASE i % 5
                    WHEN 0 THEN 'Sharma'
                    WHEN 1 THEN 'Patel'
                    WHEN 2 THEN 'Kumar'
                    WHEN 3 THEN 'Singh'
                    WHEN 4 THEN 'Verma'
                END,
                DATE '2015-01-01' + (i * 30),
                CASE WHEN i % 2 = 0 THEN 'male' ELSE 'female' END,
                CASE i % 4
                    WHEN 0 THEN 'A+'
                    WHEN 1 THEN 'B+'
                    WHEN 2 THEN 'O+'
                    WHEN 3 THEN 'AB+'
                END,
                'active',
                CURRENT_DATE - INTERVAL '6 months',
                NOW(),
                NOW()
            )
            RETURNING id INTO v_student_id;

            -- Create enrollment for this student
            IF v_student_id IS NOT NULL AND v_academic_year_id IS NOT NULL THEN
                INSERT INTO student_enrollments (
                    id, tenant_id, student_id, section_id, academic_year_id,
                    enrollment_date, status, created_at, updated_at
                ) VALUES (
                    uuid_generate_v7(),
                    v_tenant_id,
                    v_student_id,
                    CASE WHEN i <= 5 THEN v_section_id ELSE COALESCE(v_section_id_2, v_section_id) END,
                    v_academic_year_id,
                    CURRENT_DATE - INTERVAL '6 months',
                    'active',
                    NOW(),
                    NOW()
                );
            END IF;
        END IF;
    END LOOP;

    -- Get all active students with enrollments
    SELECT ARRAY_AGG(DISTINCT s.id) INTO v_student_ids
    FROM students s
    JOIN student_enrollments se ON se.student_id = s.id AND se.status = 'active'
    WHERE s.tenant_id = v_tenant_id AND s.status = 'active';

    IF v_student_ids IS NULL OR array_length(v_student_ids, 1) IS NULL THEN
        RAISE NOTICE 'No students found with active enrollments, skipping attendance seed';
        RETURN;
    END IF;

    -- Current date for calculations
    v_current_date := CURRENT_DATE;

    -- Create attendance records for the past 30 days
    FOR i IN 0..29 LOOP
        v_attendance_date := v_current_date - i;

        -- Skip weekends
        IF EXTRACT(DOW FROM v_attendance_date) NOT IN (0, 6) THEN
            -- Create attendance for each student
            FOREACH v_student_id IN ARRAY v_student_ids LOOP
                -- Check if attendance already exists
                IF NOT EXISTS (SELECT 1 FROM student_attendance WHERE tenant_id = v_tenant_id AND student_id = v_student_id AND attendance_date = v_attendance_date) THEN
                    -- Random value for status distribution
                    v_random_val := random();

                    -- Determine status: 85% present, 10% absent, 5% late
                    IF v_random_val < 0.85 THEN
                        v_status := 'present';
                    ELSIF v_random_val < 0.95 THEN
                        v_status := 'absent';
                    ELSE
                        v_status := 'late';
                    END IF;

                    -- Get the section for this student from enrollment
                    SELECT se.section_id INTO v_section_id
                    FROM student_enrollments se
                    WHERE se.student_id = v_student_id AND se.status = 'active'
                    LIMIT 1;

                    INSERT INTO student_attendance (
                        id, tenant_id, student_id, section_id, attendance_date,
                        status, late_arrival_time, remarks, marked_by, marked_at,
                        created_at, updated_at
                    ) VALUES (
                        uuid_generate_v7(),
                        v_tenant_id,
                        v_student_id,
                        v_section_id,
                        v_attendance_date,
                        v_status,
                        CASE WHEN v_status = 'late' THEN TIME '09:15:00' ELSE NULL END,
                        CASE
                            WHEN v_status = 'absent' THEN 'Sick leave'
                            WHEN v_status = 'late' THEN 'Traffic delay'
                            ELSE NULL
                        END,
                        v_user_id,
                        v_attendance_date + TIME '09:00:00',
                        NOW(),
                        NOW()
                    );
                END IF;
            END LOOP;
        END IF;
    END LOOP;

    -- Create attendance settings for the branch if not exists
    IF NOT EXISTS (SELECT 1 FROM student_attendance_settings WHERE tenant_id = v_tenant_id AND branch_id = v_branch_id) THEN
        INSERT INTO student_attendance_settings (
            id, tenant_id, branch_id, edit_window_minutes, late_threshold_minutes, sms_on_absent,
            created_at, updated_at
        ) VALUES (
            uuid_generate_v7(),
            v_tenant_id,
            v_branch_id,
            120,
            15,
            false,
            NOW(),
            NOW()
        );
    END IF;

    RAISE NOTICE 'Attendance seed data created successfully';
END $$;
