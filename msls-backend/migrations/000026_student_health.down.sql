-- Rollback Student Health Records Migration

-- Drop triggers
DROP TRIGGER IF EXISTS update_incidents_updated_at ON student_medical_incidents;
DROP TRIGGER IF EXISTS update_vaccinations_updated_at ON student_vaccinations;
DROP TRIGGER IF EXISTS update_medications_updated_at ON student_medications;
DROP TRIGGER IF EXISTS update_conditions_updated_at ON student_chronic_conditions;
DROP TRIGGER IF EXISTS update_allergies_updated_at ON student_allergies;
DROP TRIGGER IF EXISTS update_health_profiles_updated_at ON student_health_profiles;

-- Drop policies
DROP POLICY IF EXISTS tenant_isolation_incidents ON student_medical_incidents;
DROP POLICY IF EXISTS tenant_isolation_vaccinations ON student_vaccinations;
DROP POLICY IF EXISTS tenant_isolation_medications ON student_medications;
DROP POLICY IF EXISTS tenant_isolation_conditions ON student_chronic_conditions;
DROP POLICY IF EXISTS tenant_isolation_allergies ON student_allergies;
DROP POLICY IF EXISTS tenant_isolation_health_profiles ON student_health_profiles;

-- Drop tables
DROP TABLE IF EXISTS student_medical_incidents;
DROP TABLE IF EXISTS student_vaccinations;
DROP TABLE IF EXISTS student_medications;
DROP TABLE IF EXISTS student_chronic_conditions;
DROP TABLE IF EXISTS student_allergies;
DROP TABLE IF EXISTS student_health_profiles;
