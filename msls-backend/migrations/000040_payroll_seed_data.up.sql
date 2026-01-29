-- Migration: 000040_payroll_seed_data.up.sql
-- Description: Seed data for payroll - creates sample staff, salaries, and payslips for last two months
-- This migration creates test data for demonstrating payroll functionality

-- Get tenant_id and branch_id for seed data
DO $$
DECLARE
    v_tenant_id UUID;
    v_branch_id UUID;
    v_admin_user_id UUID;
    v_dept_teaching_id UUID;
    v_dept_admin_id UUID;
    v_desig_teacher_id UUID;
    v_desig_accountant_id UUID;
    v_comp_basic_id UUID;
    v_comp_hra_id UUID;
    v_comp_da_id UUID;
    v_comp_pf_id UUID;
    v_comp_tax_id UUID;
    v_structure_teacher_id UUID;
    v_structure_admin_id UUID;
    v_staff1_id UUID;
    v_staff2_id UUID;
    v_staff3_id UUID;
    v_staff4_id UUID;
    v_staff5_id UUID;
    v_salary1_id UUID;
    v_salary2_id UUID;
    v_salary3_id UUID;
    v_salary4_id UUID;
    v_salary5_id UUID;
    v_pay_run_dec_id UUID;
    v_pay_run_jan_id UUID;
    v_payslip_id UUID;
BEGIN
    -- Get tenant ID (assuming first tenant)
    SELECT id INTO v_tenant_id FROM tenants LIMIT 1;
    IF v_tenant_id IS NULL THEN
        RAISE EXCEPTION 'No tenant found';
    END IF;

    -- Get branch ID (primary branch)
    SELECT id INTO v_branch_id FROM branches WHERE tenant_id = v_tenant_id AND is_primary = true LIMIT 1;
    IF v_branch_id IS NULL THEN
        SELECT id INTO v_branch_id FROM branches WHERE tenant_id = v_tenant_id LIMIT 1;
    END IF;

    -- Get admin user ID
    SELECT id INTO v_admin_user_id FROM users WHERE tenant_id = v_tenant_id LIMIT 1;

    -- Create or get departments (unique constraint: tenant_id, branch_id, code)
    INSERT INTO departments (tenant_id, branch_id, name, code, description, is_active)
    VALUES (v_tenant_id, v_branch_id, 'Teaching', 'TEACH', 'Teaching Department', true)
    ON CONFLICT (tenant_id, branch_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_dept_teaching_id;

    INSERT INTO departments (tenant_id, branch_id, name, code, description, is_active)
    VALUES (v_tenant_id, v_branch_id, 'Administration', 'ADMIN', 'Administration Department', true)
    ON CONFLICT (tenant_id, branch_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_dept_admin_id;

    -- Create or get designations (unique constraint: tenant_id, name)
    INSERT INTO designations (tenant_id, name, level, is_active)
    VALUES (v_tenant_id, 'Senior Teacher', 3, true)
    ON CONFLICT (tenant_id, name) DO UPDATE SET level = EXCLUDED.level
    RETURNING id INTO v_desig_teacher_id;

    INSERT INTO designations (tenant_id, name, level, is_active)
    VALUES (v_tenant_id, 'Accountant', 2, true)
    ON CONFLICT (tenant_id, name) DO UPDATE SET level = EXCLUDED.level
    RETURNING id INTO v_desig_accountant_id;

    -- Create salary components (unique constraint: tenant_id, code)
    INSERT INTO salary_components (tenant_id, name, code, component_type, calculation_type, is_taxable, is_active, display_order)
    VALUES (v_tenant_id, 'Basic Salary', 'BASIC', 'earning', 'fixed', true, true, 1)
    ON CONFLICT (tenant_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_comp_basic_id;

    INSERT INTO salary_components (tenant_id, name, code, component_type, calculation_type, is_taxable, is_active, display_order)
    VALUES (v_tenant_id, 'House Rent Allowance', 'HRA', 'earning', 'fixed', false, true, 2)
    ON CONFLICT (tenant_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_comp_hra_id;

    INSERT INTO salary_components (tenant_id, name, code, component_type, calculation_type, is_taxable, is_active, display_order)
    VALUES (v_tenant_id, 'Dearness Allowance', 'DA', 'earning', 'fixed', true, true, 3)
    ON CONFLICT (tenant_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_comp_da_id;

    INSERT INTO salary_components (tenant_id, name, code, component_type, calculation_type, is_taxable, is_active, display_order)
    VALUES (v_tenant_id, 'Provident Fund', 'PF', 'deduction', 'fixed', false, true, 10)
    ON CONFLICT (tenant_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_comp_pf_id;

    INSERT INTO salary_components (tenant_id, name, code, component_type, calculation_type, is_taxable, is_active, display_order)
    VALUES (v_tenant_id, 'Professional Tax', 'PT', 'deduction', 'fixed', false, true, 11)
    ON CONFLICT (tenant_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_comp_tax_id;

    -- Create salary structures (unique constraint: tenant_id, code)
    INSERT INTO salary_structures (tenant_id, name, code, description, designation_id, is_active)
    VALUES (v_tenant_id, 'Teacher Salary Structure', 'TEACH-SAL', 'Standard salary structure for teachers', v_desig_teacher_id, true)
    ON CONFLICT (tenant_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_structure_teacher_id;

    INSERT INTO salary_structures (tenant_id, name, code, description, designation_id, is_active)
    VALUES (v_tenant_id, 'Admin Salary Structure', 'ADMIN-SAL', 'Standard salary structure for admin staff', v_desig_accountant_id, true)
    ON CONFLICT (tenant_id, code) DO UPDATE SET name = EXCLUDED.name
    RETURNING id INTO v_structure_admin_id;

    -- Add components to structures
    INSERT INTO salary_structure_components (structure_id, component_id, amount)
    VALUES
        (v_structure_teacher_id, v_comp_basic_id, 35000),
        (v_structure_teacher_id, v_comp_hra_id, 14000),
        (v_structure_teacher_id, v_comp_da_id, 7000),
        (v_structure_teacher_id, v_comp_pf_id, 4200),
        (v_structure_teacher_id, v_comp_tax_id, 200)
    ON CONFLICT (structure_id, component_id) DO UPDATE SET amount = EXCLUDED.amount;

    INSERT INTO salary_structure_components (structure_id, component_id, amount)
    VALUES
        (v_structure_admin_id, v_comp_basic_id, 28000),
        (v_structure_admin_id, v_comp_hra_id, 11200),
        (v_structure_admin_id, v_comp_da_id, 5600),
        (v_structure_admin_id, v_comp_pf_id, 3360),
        (v_structure_admin_id, v_comp_tax_id, 200)
    ON CONFLICT (structure_id, component_id) DO UPDATE SET amount = EXCLUDED.amount;

    -- Initialize employee sequence
    INSERT INTO staff_employee_sequences (tenant_id, prefix, last_sequence)
    VALUES (v_tenant_id, 'EMP', 5)
    ON CONFLICT (tenant_id, prefix) DO UPDATE SET last_sequence = GREATEST(staff_employee_sequences.last_sequence, 5);

    -- Create staff members (unique constraint: tenant_id, employee_id)
    INSERT INTO staff (tenant_id, branch_id, employee_id, employee_id_prefix, first_name, middle_name, last_name,
                       date_of_birth, gender, work_email, work_phone, staff_type, department_id, designation_id,
                       join_date, status, created_by)
    VALUES (v_tenant_id, v_branch_id, 'EMP00001', 'EMP', 'Rajesh', 'Kumar', 'Sharma',
            '1985-03-15', 'male', 'rajesh.sharma@school.com', '9876543201', 'teaching', v_dept_teaching_id, v_desig_teacher_id,
            '2022-06-01', 'active', v_admin_user_id)
    ON CONFLICT (tenant_id, employee_id) DO UPDATE SET first_name = EXCLUDED.first_name
    RETURNING id INTO v_staff1_id;

    INSERT INTO staff (tenant_id, branch_id, employee_id, employee_id_prefix, first_name, middle_name, last_name,
                       date_of_birth, gender, work_email, work_phone, staff_type, department_id, designation_id,
                       join_date, status, created_by)
    VALUES (v_tenant_id, v_branch_id, 'EMP00002', 'EMP', 'Priya', NULL, 'Patel',
            '1990-07-22', 'female', 'priya.patel@school.com', '9876543202', 'teaching', v_dept_teaching_id, v_desig_teacher_id,
            '2023-01-15', 'active', v_admin_user_id)
    ON CONFLICT (tenant_id, employee_id) DO UPDATE SET first_name = EXCLUDED.first_name
    RETURNING id INTO v_staff2_id;

    INSERT INTO staff (tenant_id, branch_id, employee_id, employee_id_prefix, first_name, middle_name, last_name,
                       date_of_birth, gender, work_email, work_phone, staff_type, department_id, designation_id,
                       join_date, status, created_by)
    VALUES (v_tenant_id, v_branch_id, 'EMP00003', 'EMP', 'Amit', 'Singh', 'Verma',
            '1988-11-05', 'male', 'amit.verma@school.com', '9876543203', 'teaching', v_dept_teaching_id, v_desig_teacher_id,
            '2021-08-01', 'active', v_admin_user_id)
    ON CONFLICT (tenant_id, employee_id) DO UPDATE SET first_name = EXCLUDED.first_name
    RETURNING id INTO v_staff3_id;

    INSERT INTO staff (tenant_id, branch_id, employee_id, employee_id_prefix, first_name, middle_name, last_name,
                       date_of_birth, gender, work_email, work_phone, staff_type, department_id, designation_id,
                       join_date, status, created_by)
    VALUES (v_tenant_id, v_branch_id, 'EMP00004', 'EMP', 'Sunita', NULL, 'Reddy',
            '1992-02-18', 'female', 'sunita.reddy@school.com', '9876543204', 'non_teaching', v_dept_admin_id, v_desig_accountant_id,
            '2023-04-01', 'active', v_admin_user_id)
    ON CONFLICT (tenant_id, employee_id) DO UPDATE SET first_name = EXCLUDED.first_name
    RETURNING id INTO v_staff4_id;

    INSERT INTO staff (tenant_id, branch_id, employee_id, employee_id_prefix, first_name, middle_name, last_name,
                       date_of_birth, gender, work_email, work_phone, staff_type, department_id, designation_id,
                       join_date, status, created_by)
    VALUES (v_tenant_id, v_branch_id, 'EMP00005', 'EMP', 'Mohammed', NULL, 'Khan',
            '1987-09-10', 'male', 'mohammed.khan@school.com', '9876543205', 'non_teaching', v_dept_admin_id, v_desig_accountant_id,
            '2022-11-01', 'active', v_admin_user_id)
    ON CONFLICT (tenant_id, employee_id) DO UPDATE SET first_name = EXCLUDED.first_name
    RETURNING id INTO v_staff5_id;

    -- Create staff salaries
    -- Staff 1: Teacher - 56000 gross
    INSERT INTO staff_salaries (tenant_id, staff_id, structure_id, effective_from, gross_salary, net_salary, ctc, is_current, created_by)
    VALUES (v_tenant_id, v_staff1_id, v_structure_teacher_id, '2024-01-01', 56000, 51600, 60200, true, v_admin_user_id)
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_salary1_id;

    IF v_salary1_id IS NOT NULL THEN
        INSERT INTO staff_salary_components (staff_salary_id, component_id, amount, is_overridden)
        VALUES
            (v_salary1_id, v_comp_basic_id, 35000, false),
            (v_salary1_id, v_comp_hra_id, 14000, false),
            (v_salary1_id, v_comp_da_id, 7000, false),
            (v_salary1_id, v_comp_pf_id, 4200, false),
            (v_salary1_id, v_comp_tax_id, 200, false)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Staff 2: Teacher - 56000 gross
    INSERT INTO staff_salaries (tenant_id, staff_id, structure_id, effective_from, gross_salary, net_salary, ctc, is_current, created_by)
    VALUES (v_tenant_id, v_staff2_id, v_structure_teacher_id, '2024-01-01', 56000, 51600, 60200, true, v_admin_user_id)
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_salary2_id;

    IF v_salary2_id IS NOT NULL THEN
        INSERT INTO staff_salary_components (staff_salary_id, component_id, amount, is_overridden)
        VALUES
            (v_salary2_id, v_comp_basic_id, 35000, false),
            (v_salary2_id, v_comp_hra_id, 14000, false),
            (v_salary2_id, v_comp_da_id, 7000, false),
            (v_salary2_id, v_comp_pf_id, 4200, false),
            (v_salary2_id, v_comp_tax_id, 200, false)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Staff 3: Senior Teacher - 62000 gross (higher salary)
    INSERT INTO staff_salaries (tenant_id, staff_id, structure_id, effective_from, gross_salary, net_salary, ctc, is_current, created_by)
    VALUES (v_tenant_id, v_staff3_id, v_structure_teacher_id, '2024-01-01', 62000, 56900, 66700, true, v_admin_user_id)
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_salary3_id;

    IF v_salary3_id IS NOT NULL THEN
        INSERT INTO staff_salary_components (staff_salary_id, component_id, amount, is_overridden)
        VALUES
            (v_salary3_id, v_comp_basic_id, 38000, true),
            (v_salary3_id, v_comp_hra_id, 15200, true),
            (v_salary3_id, v_comp_da_id, 8800, true),
            (v_salary3_id, v_comp_pf_id, 4560, true),
            (v_salary3_id, v_comp_tax_id, 540, true)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Staff 4: Accountant - 44800 gross
    INSERT INTO staff_salaries (tenant_id, staff_id, structure_id, effective_from, gross_salary, net_salary, ctc, is_current, created_by)
    VALUES (v_tenant_id, v_staff4_id, v_structure_admin_id, '2024-01-01', 44800, 41240, 48160, true, v_admin_user_id)
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_salary4_id;

    IF v_salary4_id IS NOT NULL THEN
        INSERT INTO staff_salary_components (staff_salary_id, component_id, amount, is_overridden)
        VALUES
            (v_salary4_id, v_comp_basic_id, 28000, false),
            (v_salary4_id, v_comp_hra_id, 11200, false),
            (v_salary4_id, v_comp_da_id, 5600, false),
            (v_salary4_id, v_comp_pf_id, 3360, false),
            (v_salary4_id, v_comp_tax_id, 200, false)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Staff 5: Admin - 48000 gross
    INSERT INTO staff_salaries (tenant_id, staff_id, structure_id, effective_from, gross_salary, net_salary, ctc, is_current, created_by)
    VALUES (v_tenant_id, v_staff5_id, v_structure_admin_id, '2024-01-01', 48000, 44040, 51600, true, v_admin_user_id)
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_salary5_id;

    IF v_salary5_id IS NOT NULL THEN
        INSERT INTO staff_salary_components (staff_salary_id, component_id, amount, is_overridden)
        VALUES
            (v_salary5_id, v_comp_basic_id, 30000, true),
            (v_salary5_id, v_comp_hra_id, 12000, true),
            (v_salary5_id, v_comp_da_id, 6000, true),
            (v_salary5_id, v_comp_pf_id, 3600, true),
            (v_salary5_id, v_comp_tax_id, 360, true)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Re-fetch salary IDs if they weren't returned (due to conflict)
    IF v_salary1_id IS NULL THEN
        SELECT id INTO v_salary1_id FROM staff_salaries WHERE staff_id = v_staff1_id AND is_current = true;
    END IF;
    IF v_salary2_id IS NULL THEN
        SELECT id INTO v_salary2_id FROM staff_salaries WHERE staff_id = v_staff2_id AND is_current = true;
    END IF;
    IF v_salary3_id IS NULL THEN
        SELECT id INTO v_salary3_id FROM staff_salaries WHERE staff_id = v_staff3_id AND is_current = true;
    END IF;
    IF v_salary4_id IS NULL THEN
        SELECT id INTO v_salary4_id FROM staff_salaries WHERE staff_id = v_staff4_id AND is_current = true;
    END IF;
    IF v_salary5_id IS NULL THEN
        SELECT id INTO v_salary5_id FROM staff_salaries WHERE staff_id = v_staff5_id AND is_current = true;
    END IF;

    -- ==========================================
    -- Create Pay Run for December 2025 (Finalized)
    -- ==========================================
    INSERT INTO pay_runs (tenant_id, pay_period_month, pay_period_year, status, total_staff, total_gross, total_deductions, total_net,
                          calculated_at, approved_at, approved_by, finalized_at, finalized_by, notes, created_by)
    VALUES (v_tenant_id, 12, 2025, 'finalized', 5, 266800, 16660, 250140,
            '2025-12-28 10:00:00+05:30', '2025-12-29 14:00:00+05:30', v_admin_user_id, '2025-12-30 10:00:00+05:30', v_admin_user_id,
            'December 2025 payroll - processed on time', v_admin_user_id)
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_pay_run_dec_id;

    IF v_pay_run_dec_id IS NOT NULL THEN
        -- Payslip 1: Rajesh Sharma (full attendance)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status, payment_date, payment_reference)
        VALUES (v_tenant_id, v_pay_run_dec_id, v_staff1_id, v_salary1_id, 26, 26, 0, 0, 0,
                56000, 56000, 4400, 51600, 0, 'paid', '2025-12-31', 'DEC-2025-001')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 35000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 14000, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 7000, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 4200, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 200, false);

        -- Payslip 2: Priya Patel (2 leave days)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status, payment_date, payment_reference)
        VALUES (v_tenant_id, v_pay_run_dec_id, v_staff2_id, v_salary2_id, 26, 24, 2, 0, 0,
                56000, 56000, 4400, 51600, 0, 'paid', '2025-12-31', 'DEC-2025-002')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 35000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 14000, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 7000, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 4200, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 200, false);

        -- Payslip 3: Amit Verma (1 LOP day)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status, payment_date, payment_reference)
        VALUES (v_tenant_id, v_pay_run_dec_id, v_staff3_id, v_salary3_id, 26, 24, 0, 1, 1,
                62000, 59615, 5100, 52130, 2385, 'paid', '2025-12-31', 'DEC-2025-003')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 36538, true),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 14615, true),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 8462, true),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 4385, true),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 540, false);

        -- Payslip 4: Sunita Reddy (full attendance)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status, payment_date, payment_reference)
        VALUES (v_tenant_id, v_pay_run_dec_id, v_staff4_id, v_salary4_id, 26, 26, 0, 0, 0,
                44800, 44800, 3560, 41240, 0, 'paid', '2025-12-31', 'DEC-2025-004')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 28000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 11200, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 5600, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 3360, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 200, false);

        -- Payslip 5: Mohammed Khan (1 leave day)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status, payment_date, payment_reference)
        VALUES (v_tenant_id, v_pay_run_dec_id, v_staff5_id, v_salary5_id, 26, 25, 1, 0, 0,
                48000, 48000, 3960, 44040, 0, 'paid', '2025-12-31', 'DEC-2025-005')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 30000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 12000, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 6000, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 3600, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 360, false);
    END IF;

    -- ==========================================
    -- Create Pay Run for January 2026 (Approved, ready to finalize)
    -- ==========================================
    INSERT INTO pay_runs (tenant_id, pay_period_month, pay_period_year, status, total_staff, total_gross, total_deductions, total_net,
                          calculated_at, approved_at, approved_by, notes, created_by)
    VALUES (v_tenant_id, 1, 2026, 'approved', 5, 266800, 16660, 250140,
            '2026-01-26 10:00:00+05:30', '2026-01-27 14:00:00+05:30', v_admin_user_id,
            'January 2026 payroll - awaiting finalization', v_admin_user_id)
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_pay_run_jan_id;

    IF v_pay_run_jan_id IS NOT NULL THEN
        -- Payslip 1: Rajesh Sharma (full attendance)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status)
        VALUES (v_tenant_id, v_pay_run_jan_id, v_staff1_id, v_salary1_id, 27, 27, 0, 0, 0,
                56000, 56000, 4400, 51600, 0, 'approved')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 35000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 14000, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 7000, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 4200, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 200, false);

        -- Payslip 2: Priya Patel (1 leave day)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status)
        VALUES (v_tenant_id, v_pay_run_jan_id, v_staff2_id, v_salary2_id, 27, 26, 1, 0, 0,
                56000, 56000, 4400, 51600, 0, 'approved')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 35000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 14000, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 7000, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 4200, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 200, false);

        -- Payslip 3: Amit Verma (full attendance)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status)
        VALUES (v_tenant_id, v_pay_run_jan_id, v_staff3_id, v_salary3_id, 27, 27, 0, 0, 0,
                62000, 62000, 5100, 56900, 0, 'approved')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 38000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 15200, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 8800, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 4560, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 540, false);

        -- Payslip 4: Sunita Reddy (2 LOP days)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status)
        VALUES (v_tenant_id, v_pay_run_jan_id, v_staff4_id, v_salary4_id, 27, 24, 0, 3, 2,
                44800, 41481, 3560, 34604, 3317, 'approved')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 25926, true),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 10370, true),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 5185, true),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 3111, true),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 200, false);

        -- Payslip 5: Mohammed Khan (full attendance)
        INSERT INTO payslips (tenant_id, pay_run_id, staff_id, staff_salary_id, working_days, present_days, leave_days, absent_days, lop_days,
                              gross_salary, total_earnings, total_deductions, net_salary, lop_deduction, status)
        VALUES (v_tenant_id, v_pay_run_jan_id, v_staff5_id, v_salary5_id, 27, 27, 0, 0, 0,
                48000, 48000, 3960, 44040, 0, 'approved')
        RETURNING id INTO v_payslip_id;

        INSERT INTO payslip_components (payslip_id, component_id, component_name, component_code, component_type, amount, is_prorated)
        VALUES
            (v_payslip_id, v_comp_basic_id, 'Basic Salary', 'BASIC', 'earning', 30000, false),
            (v_payslip_id, v_comp_hra_id, 'House Rent Allowance', 'HRA', 'earning', 12000, false),
            (v_payslip_id, v_comp_da_id, 'Dearness Allowance', 'DA', 'earning', 6000, false),
            (v_payslip_id, v_comp_pf_id, 'Provident Fund', 'PF', 'deduction', 3600, false),
            (v_payslip_id, v_comp_tax_id, 'Professional Tax', 'PT', 'deduction', 360, false);
    END IF;

    RAISE NOTICE 'Payroll seed data created successfully';
    RAISE NOTICE 'Created 5 staff members with salary assignments';
    RAISE NOTICE 'Created pay run for December 2025 (Finalized)';
    RAISE NOTICE 'Created pay run for January 2026 (Approved)';
END $$;
