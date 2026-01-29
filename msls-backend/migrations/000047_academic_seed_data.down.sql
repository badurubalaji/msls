-- Migration: 000047_academic_seed_data.down.sql
-- Description: Remove seed data for classes, sections, and subjects

-- Remove class-subject mappings
DELETE FROM class_subjects
WHERE class_id IN (SELECT id FROM classes WHERE code IN ('LKG', 'UKG', 'I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X', 'XI', 'XII'));

-- Remove class-stream mappings
DELETE FROM class_streams
WHERE class_id IN (SELECT id FROM classes WHERE code IN ('XI', 'XII'));

-- Remove sections for seeded classes
DELETE FROM sections
WHERE class_id IN (SELECT id FROM classes WHERE code IN ('LKG', 'UKG', 'I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X', 'XI', 'XII'));

-- Remove seeded classes
DELETE FROM classes
WHERE code IN ('LKG', 'UKG', 'I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X', 'XI', 'XII');

-- Remove seeded subjects
DELETE FROM subjects
WHERE code IN (
    'ENG', 'HIN', 'MATH', 'SCI', 'SST',
    'PHY', 'CHEM', 'BIO',
    'ACC', 'BST', 'ECO',
    'HIST', 'GEO', 'POL', 'PSY', 'SOC',
    'CS', 'IT', 'PE',
    'SANS', 'FRE', 'GER',
    'ART', 'MUS', 'DAN',
    'EVS', 'MOR', 'GK'
);
