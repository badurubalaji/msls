-- Migration: 000044_staff_document_seed_data.down.sql
-- Remove seeded document types

DELETE FROM staff_document_types WHERE code IN (
    'aadhaar', 'pan', 'passport', 'driving_license', 'voter_id',
    'ssc_marksheet', 'hsc_marksheet', 'degree_certificate', 'bed_certificate',
    'pg_certificate', 'professional_cert',
    'offer_letter', 'appointment_letter', 'experience_letter', 'relieving_letter', 'salary_slip',
    'police_verification', 'medical_fitness', 'background_check',
    'bank_details', 'photograph', 'address_proof'
);
