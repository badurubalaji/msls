-- Drop exam types table and related objects

DROP TRIGGER IF EXISTS set_updated_at_exam_types ON exam_types;
DROP FUNCTION IF EXISTS update_exam_types_updated_at();

DROP POLICY IF EXISTS tenant_isolation_exam_types ON exam_types;

DROP TABLE IF EXISTS exam_types;

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'exam:type:view',
    'exam:type:create',
    'exam:type:update',
    'exam:type:delete'
);
