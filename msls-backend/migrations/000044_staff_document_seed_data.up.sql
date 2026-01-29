-- Migration: 000044_staff_document_seed_data.up.sql
-- Seed default document types for all tenants

-- Seed default document types for all existing tenants
INSERT INTO staff_document_types (id, tenant_id, name, code, category, is_mandatory, has_expiry, default_validity_months, applicable_to, display_order)
SELECT
    uuid_generate_v7(),
    t.id,
    dt.name,
    dt.code,
    dt.category,
    dt.is_mandatory,
    dt.has_expiry,
    dt.default_validity_months,
    dt.applicable_to::VARCHAR(20)[],
    dt.display_order
FROM tenants t
CROSS JOIN (
    VALUES
        ('Aadhaar Card', 'aadhaar', 'identity', true, false, NULL, '{teaching,non_teaching}', 1),
        ('PAN Card', 'pan', 'identity', true, false, NULL, '{teaching,non_teaching}', 2),
        ('Passport', 'passport', 'identity', false, true, 120, '{teaching,non_teaching}', 3),
        ('Driving License', 'driving_license', 'identity', false, true, 60, '{teaching,non_teaching}', 4),
        ('Voter ID', 'voter_id', 'identity', false, false, NULL, '{teaching,non_teaching}', 5),
        ('10th Marksheet', 'ssc_marksheet', 'education', true, false, NULL, '{teaching,non_teaching}', 10),
        ('12th Marksheet', 'hsc_marksheet', 'education', true, false, NULL, '{teaching,non_teaching}', 11),
        ('Degree Certificate', 'degree_certificate', 'education', false, false, NULL, '{teaching}', 12),
        ('B.Ed Certificate', 'bed_certificate', 'education', false, false, NULL, '{teaching}', 13),
        ('Post Graduate Certificate', 'pg_certificate', 'education', false, false, NULL, '{teaching}', 14),
        ('Professional Certification', 'professional_cert', 'education', false, true, 36, '{teaching,non_teaching}', 15),
        ('Offer Letter', 'offer_letter', 'employment', true, false, NULL, '{teaching,non_teaching}', 20),
        ('Appointment Letter', 'appointment_letter', 'employment', true, false, NULL, '{teaching,non_teaching}', 21),
        ('Experience Letter', 'experience_letter', 'employment', false, false, NULL, '{teaching,non_teaching}', 22),
        ('Relieving Letter', 'relieving_letter', 'employment', false, false, NULL, '{teaching,non_teaching}', 23),
        ('Salary Slip', 'salary_slip', 'employment', false, false, NULL, '{teaching,non_teaching}', 24),
        ('Police Verification', 'police_verification', 'compliance', true, true, 36, '{teaching,non_teaching}', 30),
        ('Medical Fitness Certificate', 'medical_fitness', 'compliance', false, true, 12, '{teaching,non_teaching}', 31),
        ('Background Check Report', 'background_check', 'compliance', false, true, 24, '{teaching,non_teaching}', 32),
        ('Bank Account Details', 'bank_details', 'other', true, false, NULL, '{teaching,non_teaching}', 40),
        ('Photograph', 'photograph', 'other', true, false, NULL, '{teaching,non_teaching}', 41),
        ('Address Proof', 'address_proof', 'other', false, false, NULL, '{teaching,non_teaching}', 42)
) AS dt(name, code, category, is_mandatory, has_expiry, default_validity_months, applicable_to, display_order)
ON CONFLICT (tenant_id, code) DO NOTHING;
