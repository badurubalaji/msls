-- Migration: 000055_attendance_audit.up.sql
-- Story 7.3: Attendance Edit & History - Audit Trail

-- ============================================================================
-- Student Attendance Audit Table
-- Tracks all changes to student attendance records for complete audit trail
-- ============================================================================

CREATE TABLE student_attendance_audit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    attendance_id UUID NOT NULL REFERENCES student_attendance(id) ON DELETE CASCADE,

    -- Previous values (NULL for initial creation)
    previous_status VARCHAR(20),
    previous_remarks TEXT,
    previous_late_arrival_time TIME,

    -- New values
    new_status VARCHAR(20) NOT NULL,
    new_remarks TEXT,
    new_late_arrival_time TIME,

    -- Change metadata
    change_type VARCHAR(20) NOT NULL DEFAULT 'edit', -- 'create', 'edit'
    change_reason TEXT NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT valid_change_type CHECK (change_type IN ('create', 'edit')),
    CONSTRAINT valid_previous_status CHECK (previous_status IS NULL OR previous_status IN ('present', 'absent', 'late', 'half_day')),
    CONSTRAINT valid_new_status CHECK (new_status IN ('present', 'absent', 'late', 'half_day'))
);

-- Comments
COMMENT ON TABLE student_attendance_audit IS 'Audit trail for student attendance changes';
COMMENT ON COLUMN student_attendance_audit.change_type IS 'Type of change: create (initial marking) or edit (modification)';
COMMENT ON COLUMN student_attendance_audit.change_reason IS 'Reason for the change - required for all edits';

-- Indexes
CREATE INDEX idx_attendance_audit_attendance ON student_attendance_audit(attendance_id);
CREATE INDEX idx_attendance_audit_tenant ON student_attendance_audit(tenant_id);
CREATE INDEX idx_attendance_audit_changed_at ON student_attendance_audit(changed_at DESC);
CREATE INDEX idx_attendance_audit_changed_by ON student_attendance_audit(changed_by);

-- Row Level Security
ALTER TABLE student_attendance_audit ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON student_attendance_audit
    FOR ALL
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Grant permissions
GRANT SELECT, INSERT ON student_attendance_audit TO msls;

-- ============================================================================
-- Trigger to auto-create audit record on attendance update
-- ============================================================================

CREATE OR REPLACE FUNCTION fn_audit_student_attendance()
RETURNS TRIGGER AS $$
BEGIN
    -- Only log if status or remarks changed
    IF (OLD.status IS DISTINCT FROM NEW.status) OR
       (OLD.remarks IS DISTINCT FROM NEW.remarks) OR
       (OLD.late_arrival_time IS DISTINCT FROM NEW.late_arrival_time) THEN

        INSERT INTO student_attendance_audit (
            tenant_id,
            attendance_id,
            previous_status,
            previous_remarks,
            previous_late_arrival_time,
            new_status,
            new_remarks,
            new_late_arrival_time,
            change_type,
            change_reason,
            changed_by,
            changed_at
        ) VALUES (
            NEW.tenant_id,
            NEW.id,
            OLD.status,
            OLD.remarks,
            OLD.late_arrival_time,
            NEW.status,
            NEW.remarks,
            NEW.late_arrival_time,
            'edit',
            COALESCE(current_setting('app.audit_reason', true), 'System update'),
            COALESCE(current_setting('app.current_user_id', true)::UUID, NEW.marked_by),
            NOW()
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_audit_student_attendance
    AFTER UPDATE ON student_attendance
    FOR EACH ROW
    EXECUTE FUNCTION fn_audit_student_attendance();

-- ============================================================================
-- Update student_attendance_settings for edit window configuration
-- The edit_window_minutes column already exists from Story 7.1
-- ============================================================================
