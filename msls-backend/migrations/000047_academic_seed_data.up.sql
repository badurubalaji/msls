-- Migration: 000047_academic_seed_data.up.sql
-- Description: Seed data for classes, sections, and subjects

-- First, get the branch_id from the first branch of each tenant
-- Insert Classes (LKG to Class 12) for each tenant
INSERT INTO classes (tenant_id, branch_id, name, code, level, display_order, has_streams, description, created_at, updated_at)
SELECT
    t.id as tenant_id,
    b.id as branch_id,
    c.name,
    c.code,
    c.level::VARCHAR(30),
    c.display_order,
    c.has_streams,
    c.description,
    NOW(),
    NOW()
FROM tenants t
JOIN LATERAL (
    SELECT id FROM branches WHERE tenant_id = t.id LIMIT 1
) b ON true
CROSS JOIN (
    VALUES
        ('LKG', 'LKG', 'nursery', 1, false, 'Lower Kindergarten'),
        ('UKG', 'UKG', 'nursery', 2, false, 'Upper Kindergarten'),
        ('Class 1', 'I', 'primary', 3, false, 'First standard'),
        ('Class 2', 'II', 'primary', 4, false, 'Second standard'),
        ('Class 3', 'III', 'primary', 5, false, 'Third standard'),
        ('Class 4', 'IV', 'primary', 6, false, 'Fourth standard'),
        ('Class 5', 'V', 'primary', 7, false, 'Fifth standard'),
        ('Class 6', 'VI', 'middle', 8, false, 'Sixth standard'),
        ('Class 7', 'VII', 'middle', 9, false, 'Seventh standard'),
        ('Class 8', 'VIII', 'middle', 10, false, 'Eighth standard'),
        ('Class 9', 'IX', 'secondary', 11, false, 'Ninth standard'),
        ('Class 10', 'X', 'secondary', 12, false, 'Tenth standard (Board exams)'),
        ('Class 11', 'XI', 'senior_secondary', 13, true, 'Eleventh standard with streams'),
        ('Class 12', 'XII', 'senior_secondary', 14, true, 'Twelfth standard (Board exams)')
) AS c(name, code, level, display_order, has_streams, description)
ON CONFLICT (tenant_id, branch_id, code) DO NOTHING;

-- Insert Sections (A, B, C) for each class
INSERT INTO sections (tenant_id, class_id, name, code, capacity, display_order, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    s.name,
    s.code,
    40,
    s.display_order,
    NOW(),
    NOW()
FROM classes c
CROSS JOIN (
    VALUES
        ('Section A', 'A', 1),
        ('Section B', 'B', 2),
        ('Section C', 'C', 3)
) AS s(name, code, display_order)
ON CONFLICT (tenant_id, class_id, code) DO NOTHING;

-- Insert Subjects
INSERT INTO subjects (tenant_id, name, code, short_name, subject_type, max_marks, passing_marks, credit_hours, display_order, description, created_at, updated_at)
SELECT
    t.id as tenant_id,
    s.name,
    s.code,
    s.short_name,
    s.subject_type,
    s.max_marks,
    s.passing_marks,
    s.credit_hours,
    s.display_order,
    s.description,
    NOW(),
    NOW()
FROM tenants t
CROSS JOIN (
    VALUES
        -- Core Subjects
        ('English', 'ENG', 'Eng', 'core', 100, 35, 5.0, 1, 'English Language and Literature'),
        ('Hindi', 'HIN', 'Hin', 'language', 100, 35, 5.0, 2, 'Hindi Language and Literature'),
        ('Mathematics', 'MATH', 'Math', 'core', 100, 35, 6.0, 3, 'Mathematics'),
        ('Science', 'SCI', 'Sci', 'core', 100, 35, 5.0, 4, 'General Science (for primary/middle)'),
        ('Social Studies', 'SST', 'SST', 'core', 100, 35, 5.0, 5, 'Social Studies'),

        -- Science Stream Subjects
        ('Physics', 'PHY', 'Phy', 'core', 100, 35, 5.0, 6, 'Physics'),
        ('Chemistry', 'CHEM', 'Chem', 'core', 100, 35, 5.0, 7, 'Chemistry'),
        ('Biology', 'BIO', 'Bio', 'elective', 100, 35, 5.0, 8, 'Biology'),

        -- Commerce Stream Subjects
        ('Accountancy', 'ACC', 'Acc', 'core', 100, 35, 5.0, 9, 'Accountancy'),
        ('Business Studies', 'BST', 'BusSt', 'core', 100, 35, 5.0, 10, 'Business Studies'),
        ('Economics', 'ECO', 'Eco', 'core', 100, 35, 5.0, 11, 'Economics'),

        -- Arts Stream Subjects
        ('History', 'HIST', 'Hist', 'core', 100, 35, 5.0, 12, 'History'),
        ('Geography', 'GEO', 'Geo', 'core', 100, 35, 5.0, 13, 'Geography'),
        ('Political Science', 'POL', 'Pol', 'core', 100, 35, 5.0, 14, 'Political Science'),
        ('Psychology', 'PSY', 'Psy', 'elective', 100, 35, 4.0, 15, 'Psychology'),
        ('Sociology', 'SOC', 'Soc', 'elective', 100, 35, 4.0, 16, 'Sociology'),

        -- Electives
        ('Computer Science', 'CS', 'CS', 'elective', 100, 35, 5.0, 17, 'Computer Science'),
        ('Information Technology', 'IT', 'IT', 'elective', 100, 35, 4.0, 18, 'Information Technology'),
        ('Physical Education', 'PE', 'PE', 'elective', 100, 35, 3.0, 19, 'Physical Education'),

        -- Additional Languages
        ('Sanskrit', 'SANS', 'Sans', 'language', 100, 35, 4.0, 20, 'Sanskrit'),
        ('French', 'FRE', 'Fr', 'language', 100, 35, 4.0, 21, 'French'),
        ('German', 'GER', 'Ger', 'language', 100, 35, 4.0, 22, 'German'),

        -- Co-curricular
        ('Art & Craft', 'ART', 'Art', 'co_curricular', 50, 20, 2.0, 23, 'Art and Craft'),
        ('Music', 'MUS', 'Mus', 'co_curricular', 50, 20, 2.0, 24, 'Music'),
        ('Dance', 'DAN', 'Dan', 'co_curricular', 50, 20, 2.0, 25, 'Dance'),

        -- Environmental Studies (for primary)
        ('Environmental Studies', 'EVS', 'EVS', 'core', 100, 35, 4.0, 26, 'Environmental Studies'),

        -- Moral Science
        ('Moral Science', 'MOR', 'Mor', 'co_curricular', 50, 20, 1.0, 27, 'Moral Science and Value Education'),

        -- General Knowledge
        ('General Knowledge', 'GK', 'GK', 'co_curricular', 50, 20, 1.0, 28, 'General Knowledge')
) AS s(name, code, short_name, subject_type, max_marks, passing_marks, credit_hours, display_order, description)
ON CONFLICT (tenant_id, code) DO NOTHING;

-- Link Class-Streams for Class 11 and 12
INSERT INTO class_streams (tenant_id, class_id, stream_id, is_active, created_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    s.id as stream_id,
    true,
    NOW()
FROM classes c
JOIN streams s ON s.tenant_id = c.tenant_id
WHERE c.code IN ('XI', 'XII') AND c.has_streams = true
ON CONFLICT (class_id, stream_id) DO NOTHING;

-- Assign subjects to classes
-- Primary classes (1-5): English, Hindi, Mathematics, EVS, Art, Music, GK, Moral Science
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    CASE WHEN sub.code IN ('ART', 'MUS', 'GK', 'MOR') THEN false ELSE true END as is_mandatory,
    CASE
        WHEN sub.code IN ('ENG', 'MATH') THEN 8
        WHEN sub.code IN ('HIN', 'EVS') THEN 6
        ELSE 2
    END as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'primary'
AND sub.code IN ('ENG', 'HIN', 'MATH', 'EVS', 'ART', 'MUS', 'GK', 'MOR')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;

-- Middle classes (6-8): English, Hindi, Mathematics, Science, Social Studies, Sanskrit/French, Computer, Art, PE, GK
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    CASE WHEN sub.code IN ('ART', 'PE', 'GK', 'SANS', 'FRE', 'CS') THEN false ELSE true END as is_mandatory,
    CASE
        WHEN sub.code IN ('ENG', 'MATH') THEN 7
        WHEN sub.code IN ('HIN', 'SCI', 'SST') THEN 6
        WHEN sub.code IN ('SANS', 'FRE') THEN 4
        ELSE 2
    END as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'middle'
AND sub.code IN ('ENG', 'HIN', 'MATH', 'SCI', 'SST', 'SANS', 'FRE', 'CS', 'ART', 'PE', 'GK')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;

-- Secondary classes (9-10): English, Hindi, Mathematics, Science, Social Studies, IT/Computer, PE
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    CASE WHEN sub.code IN ('IT', 'PE') THEN false ELSE true END as is_mandatory,
    CASE
        WHEN sub.code IN ('ENG', 'MATH', 'SCI') THEN 7
        WHEN sub.code IN ('HIN', 'SST') THEN 6
        ELSE 3
    END as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'secondary'
AND sub.code IN ('ENG', 'HIN', 'MATH', 'SCI', 'SST', 'IT', 'PE')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;

-- Senior Secondary (11-12): Will vary by stream - adding common + stream-specific
-- Common for all streams
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    true as is_mandatory,
    5 as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'senior_secondary'
AND sub.code IN ('ENG')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;

-- Science stream subjects for 11-12
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    CASE WHEN sub.code IN ('BIO', 'CS', 'PE') THEN false ELSE true END as is_mandatory,
    6 as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'senior_secondary'
AND sub.code IN ('PHY', 'CHEM', 'MATH', 'BIO', 'CS', 'PE')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;

-- Commerce stream subjects for 11-12
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    CASE WHEN sub.code IN ('IT', 'PE') THEN false ELSE true END as is_mandatory,
    6 as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'senior_secondary'
AND sub.code IN ('ACC', 'BST', 'ECO', 'MATH', 'IT', 'PE')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;

-- Arts stream subjects for 11-12
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    CASE WHEN sub.code IN ('PSY', 'SOC', 'PE') THEN false ELSE true END as is_mandatory,
    6 as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'senior_secondary'
AND sub.code IN ('HIST', 'GEO', 'POL', 'ECO', 'PSY', 'SOC', 'PE')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;

-- Nursery classes (LKG, UKG): English, Hindi, Math basics, Art, Music, GK
INSERT INTO class_subjects (tenant_id, class_id, subject_id, is_mandatory, periods_per_week, created_at, updated_at)
SELECT
    c.tenant_id,
    c.id as class_id,
    sub.id as subject_id,
    true as is_mandatory,
    CASE
        WHEN sub.code IN ('ENG', 'HIN', 'MATH') THEN 5
        ELSE 3
    END as periods_per_week,
    NOW(),
    NOW()
FROM classes c
JOIN subjects sub ON sub.tenant_id = c.tenant_id
WHERE c.level = 'nursery'
AND sub.code IN ('ENG', 'HIN', 'MATH', 'ART', 'MUS', 'GK')
ON CONFLICT (tenant_id, class_id, subject_id) DO NOTHING;
