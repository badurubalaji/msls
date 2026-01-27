-- Student Health Records Migration
-- Manages student health profiles, allergies, conditions, medications, vaccinations, and medical incidents

-- ============================================================================
-- Student Health Profiles
-- ============================================================================
CREATE TABLE student_health_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,

    -- Basic health info
    blood_group VARCHAR(5),
    height_cm DECIMAL(5,2),
    weight_kg DECIMAL(5,2),
    vision_left VARCHAR(10),
    vision_right VARCHAR(10),
    hearing_status VARCHAR(20) DEFAULT 'normal',

    -- Medical history
    medical_notes TEXT,

    -- Insurance
    insurance_provider VARCHAR(100),
    insurance_policy_number VARCHAR(50),
    insurance_expiry DATE,

    -- Emergency medical info
    preferred_hospital VARCHAR(200),
    family_doctor_name VARCHAR(100),
    family_doctor_phone VARCHAR(15),

    -- Metadata
    last_checkup_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT uq_student_health_profile UNIQUE (tenant_id, student_id)
);

-- ============================================================================
-- Student Allergies
-- ============================================================================
CREATE TABLE student_allergies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,

    -- Allergy details
    allergen VARCHAR(100) NOT NULL,
    allergy_type VARCHAR(30) NOT NULL, -- food, medication, environmental, insect, other
    severity VARCHAR(20) NOT NULL, -- mild, moderate, severe, life_threatening

    -- Reaction and treatment
    reaction_description TEXT,
    treatment_instructions TEXT,

    -- Medication for emergency
    emergency_medication VARCHAR(100),

    -- Status
    diagnosed_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- ============================================================================
-- Student Chronic Conditions
-- ============================================================================
CREATE TABLE student_chronic_conditions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,

    -- Condition details
    condition_name VARCHAR(100) NOT NULL,
    condition_type VARCHAR(30) NOT NULL, -- respiratory, cardiac, neurological, endocrine, other
    severity VARCHAR(20) NOT NULL, -- mild, moderate, severe

    -- Management
    management_plan TEXT,
    restrictions TEXT,
    triggers TEXT,

    -- Status
    diagnosed_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- ============================================================================
-- Student Medications
-- ============================================================================
CREATE TABLE student_medications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,

    -- Medication details
    medication_name VARCHAR(100) NOT NULL,
    dosage VARCHAR(50) NOT NULL,
    frequency VARCHAR(50) NOT NULL, -- daily, twice_daily, as_needed, etc.
    route VARCHAR(30) NOT NULL, -- oral, injection, inhaler, topical, etc.

    -- Purpose and instructions
    purpose VARCHAR(200),
    special_instructions TEXT,

    -- Timing
    start_date DATE NOT NULL,
    end_date DATE,

    -- Administration at school
    administered_at_school BOOLEAN NOT NULL DEFAULT false,
    school_administration_time VARCHAR(50),

    -- Prescribing doctor
    prescribing_doctor VARCHAR(100),
    prescription_date DATE,

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- ============================================================================
-- Student Vaccinations
-- ============================================================================
CREATE TABLE student_vaccinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,

    -- Vaccination details
    vaccine_name VARCHAR(100) NOT NULL,
    vaccine_type VARCHAR(50), -- required, optional, booster
    dose_number INTEGER NOT NULL DEFAULT 1,

    -- Administration
    administered_date DATE NOT NULL,
    administered_by VARCHAR(100),
    administration_site VARCHAR(50), -- left_arm, right_arm, thigh, etc.
    batch_number VARCHAR(50),

    -- Next dose
    next_due_date DATE,

    -- Reaction
    had_reaction BOOLEAN NOT NULL DEFAULT false,
    reaction_description TEXT,

    -- Certificate
    certificate_url TEXT,

    -- Status
    is_verified BOOLEAN NOT NULL DEFAULT false,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    notes TEXT,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- ============================================================================
-- Student Medical Incidents
-- ============================================================================
CREATE TABLE student_medical_incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,

    -- Incident details
    incident_date DATE NOT NULL,
    incident_time TIME NOT NULL,
    location VARCHAR(100),
    incident_type VARCHAR(30) NOT NULL, -- illness, injury, emergency, other

    -- Description
    description TEXT NOT NULL,
    symptoms TEXT,

    -- Response
    first_aid_given BOOLEAN NOT NULL DEFAULT false,
    first_aid_description TEXT,
    action_taken TEXT NOT NULL,

    -- Follow-up
    parent_notified BOOLEAN NOT NULL DEFAULT false,
    parent_notified_at TIMESTAMPTZ,
    parent_notified_by UUID REFERENCES users(id),

    hospital_visit_required BOOLEAN NOT NULL DEFAULT false,
    hospital_name VARCHAR(200),
    hospital_visit_date DATE,

    -- Recovery
    student_sent_home BOOLEAN NOT NULL DEFAULT false,
    return_to_class_time TIME,
    follow_up_required BOOLEAN NOT NULL DEFAULT false,
    follow_up_date DATE,
    follow_up_notes TEXT,

    -- Outcome
    outcome TEXT,

    -- Reported by
    reported_by UUID NOT NULL REFERENCES users(id),
    notes TEXT,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES users(id)
);

-- ============================================================================
-- Indexes
-- ============================================================================

-- Health profiles
CREATE INDEX idx_health_profiles_tenant ON student_health_profiles(tenant_id);
CREATE INDEX idx_health_profiles_student ON student_health_profiles(student_id);

-- Allergies
CREATE INDEX idx_allergies_tenant ON student_allergies(tenant_id);
CREATE INDEX idx_allergies_student ON student_allergies(student_id);
CREATE INDEX idx_allergies_active ON student_allergies(student_id, is_active) WHERE is_active = true;
CREATE INDEX idx_allergies_severity ON student_allergies(severity);

-- Chronic conditions
CREATE INDEX idx_conditions_tenant ON student_chronic_conditions(tenant_id);
CREATE INDEX idx_conditions_student ON student_chronic_conditions(student_id);
CREATE INDEX idx_conditions_active ON student_chronic_conditions(student_id, is_active) WHERE is_active = true;

-- Medications
CREATE INDEX idx_medications_tenant ON student_medications(tenant_id);
CREATE INDEX idx_medications_student ON student_medications(student_id);
CREATE INDEX idx_medications_active ON student_medications(student_id, is_active) WHERE is_active = true;
CREATE INDEX idx_medications_school ON student_medications(administered_at_school) WHERE administered_at_school = true;

-- Vaccinations
CREATE INDEX idx_vaccinations_tenant ON student_vaccinations(tenant_id);
CREATE INDEX idx_vaccinations_student ON student_vaccinations(student_id);
CREATE INDEX idx_vaccinations_due_date ON student_vaccinations(next_due_date);
CREATE INDEX idx_vaccinations_verified ON student_vaccinations(is_verified);

-- Medical incidents
CREATE INDEX idx_incidents_tenant ON student_medical_incidents(tenant_id);
CREATE INDEX idx_incidents_student ON student_medical_incidents(student_id);
CREATE INDEX idx_incidents_date ON student_medical_incidents(incident_date DESC);
CREATE INDEX idx_incidents_type ON student_medical_incidents(incident_type);

-- ============================================================================
-- Row Level Security
-- ============================================================================

-- Health profiles
ALTER TABLE student_health_profiles ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_health_profiles ON student_health_profiles
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Allergies
ALTER TABLE student_allergies ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_allergies ON student_allergies
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Chronic conditions
ALTER TABLE student_chronic_conditions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_conditions ON student_chronic_conditions
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Medications
ALTER TABLE student_medications ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_medications ON student_medications
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Vaccinations
ALTER TABLE student_vaccinations ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_vaccinations ON student_vaccinations
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Medical incidents
ALTER TABLE student_medical_incidents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_incidents ON student_medical_incidents
    USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

