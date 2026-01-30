-- Examinations Table
-- Story 8.2: Examination Creation & Scheduling

CREATE TABLE examinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(200) NOT NULL,
    exam_type_id UUID NOT NULL REFERENCES exam_types(id),
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT chk_exam_dates CHECK (end_date >= start_date),
    CONSTRAINT chk_exam_status CHECK (status IN ('draft', 'scheduled', 'ongoing', 'completed', 'cancelled'))
);

-- Enable RLS
ALTER TABLE examinations ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_examinations ON examinations
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_examinations_tenant ON examinations(tenant_id);
CREATE INDEX idx_examinations_type ON examinations(tenant_id, exam_type_id);
CREATE INDEX idx_examinations_year ON examinations(tenant_id, academic_year_id);
CREATE INDEX idx_examinations_status ON examinations(tenant_id, status);
CREATE INDEX idx_examinations_dates ON examinations(tenant_id, start_date, end_date);

-- Updated at trigger
CREATE TRIGGER set_updated_at_examinations
    BEFORE UPDATE ON examinations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- Examination Classes (Many-to-Many)
-- ============================================================

CREATE TABLE examination_classes (
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (examination_id, class_id)
);

-- Enable RLS (inherits from parent)
ALTER TABLE examination_classes ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_examination_classes ON examination_classes
    USING (examination_id IN (
        SELECT id FROM examinations
        WHERE tenant_id = current_setting('app.tenant_id', true)::UUID
    ));

-- Index for reverse lookup
CREATE INDEX idx_examination_classes_class ON examination_classes(class_id);

-- ============================================================
-- Exam Schedules
-- ============================================================

CREATE TABLE exam_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id),
    exam_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_marks INTEGER NOT NULL DEFAULT 100,
    passing_marks INTEGER,
    venue VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_schedule_times CHECK (end_time > start_time),
    CONSTRAINT chk_schedule_marks CHECK (max_marks > 0),
    CONSTRAINT chk_schedule_passing CHECK (passing_marks IS NULL OR (passing_marks >= 0 AND passing_marks <= max_marks))
);

-- Enable RLS
ALTER TABLE exam_schedules ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_exam_schedules ON exam_schedules
    USING (examination_id IN (
        SELECT id FROM examinations
        WHERE tenant_id = current_setting('app.tenant_id', true)::UUID
    ));

-- Indexes
CREATE INDEX idx_exam_schedules_examination ON exam_schedules(examination_id);
CREATE INDEX idx_exam_schedules_subject ON exam_schedules(subject_id);
CREATE INDEX idx_exam_schedules_date ON exam_schedules(exam_date);

-- Unique constraint: one subject per examination
CREATE UNIQUE INDEX idx_exam_schedules_unique_subject ON exam_schedules(examination_id, subject_id);

-- Updated at trigger
CREATE TRIGGER set_updated_at_exam_schedules
    BEFORE UPDATE ON exam_schedules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- Permissions
-- ============================================================

INSERT INTO permissions (id, code, name, description, module, created_at, updated_at)
VALUES
    (uuid_generate_v7(), 'exam:view', 'View Examinations', 'Permission to view examinations and schedules', 'exam', NOW(), NOW()),
    (uuid_generate_v7(), 'exam:create', 'Create Examinations', 'Permission to create new examinations', 'exam', NOW(), NOW()),
    (uuid_generate_v7(), 'exam:update', 'Update Examinations', 'Permission to update examination details', 'exam', NOW(), NOW()),
    (uuid_generate_v7(), 'exam:delete', 'Delete Examinations', 'Permission to delete examinations', 'exam', NOW(), NOW()),
    (uuid_generate_v7(), 'exam:publish', 'Publish Examinations', 'Permission to publish/schedule examinations', 'exam', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign permissions to roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('super_admin', 'admin', 'principal')
AND p.code IN ('exam:view', 'exam:create', 'exam:update', 'exam:delete', 'exam:publish')
ON CONFLICT DO NOTHING;

-- Coordinators can view, create, update, and publish
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'coordinator'
AND p.code IN ('exam:view', 'exam:create', 'exam:update', 'exam:publish')
ON CONFLICT DO NOTHING;

-- Teachers can view
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code = 'exam:view'
ON CONFLICT DO NOTHING;
