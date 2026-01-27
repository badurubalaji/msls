// Package main provides a database seeder for development and testing.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/config"
	"msls-backend/internal/pkg/database"
	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/services/auth"
)

// Seed credentials - ONLY FOR DEVELOPMENT
const (
	DefaultTenantName     = "Demo School"
	DefaultTenantSlug     = "demo-school"
	SuperAdminEmail       = "admin@myschoolsystem.com"
	SuperAdminPassword    = "Admin@123"
	SuperAdminFirstName   = "Super"
	SuperAdminLastName    = "Admin"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "seed error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("=== MSLS Database Seeder ===")
	fmt.Println()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize database connection
	dbConfig := database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.Name,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}
	conn, err := database.New(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	db := conn.DB()

	fmt.Println("Connected to database")

	// Bypass RLS for seeding
	if err := db.Exec("SET app.bypass_rls = 'true'").Error; err != nil {
		return fmt.Errorf("failed to bypass RLS: %w", err)
	}

	ctx := context.Background()

	// Check if tenant already exists
	var existingTenantID string
	err = db.Raw("SELECT id::text FROM tenants WHERE slug = ?", DefaultTenantSlug).Scan(&existingTenantID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing tenant: %w", err)
	}

	var tenantID uuid.UUID
	if existingTenantID != "" {
		tenantID, _ = uuid.Parse(existingTenantID)
		fmt.Printf("Tenant '%s' already exists (ID: %s)\n", DefaultTenantName, tenantID)
	} else {
		// Create tenant
		tenantID = uuid.New()
		err = db.Exec(`
			INSERT INTO tenants (id, name, slug, settings, status, created_at, updated_at)
			VALUES (?, ?, ?, '{}', 'active', NOW(), NOW())
		`, tenantID, DefaultTenantName, DefaultTenantSlug).Error
		if err != nil {
			return fmt.Errorf("failed to create tenant: %w", err)
		}
		fmt.Printf("Created tenant: %s (ID: %s)\n", DefaultTenantName, tenantID)
	}

	// Check if super admin already exists
	var existingUserID string
	err = db.Raw("SELECT id::text FROM users WHERE tenant_id = ? AND email = ?", tenantID, SuperAdminEmail).Scan(&existingUserID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	var userID uuid.UUID
	if existingUserID != "" {
		userID, _ = uuid.Parse(existingUserID)
		fmt.Printf("Super admin '%s' already exists (ID: %s)\n", SuperAdminEmail, userID)
	} else {
		// Hash password using Argon2id
		passwordService := auth.NewPasswordService()
		hashedPassword, err := passwordService.HashPassword(SuperAdminPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Create super admin user
		userID = uuid.New()
		now := time.Now()
		err = db.Exec(`
			INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, status, email_verified_at, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, 'active', ?, NOW(), NOW())
		`, userID, tenantID, SuperAdminEmail, hashedPassword, SuperAdminFirstName, SuperAdminLastName, now).Error
		if err != nil {
			return fmt.Errorf("failed to create super admin user: %w", err)
		}
		fmt.Printf("Created super admin user: %s (ID: %s)\n", SuperAdminEmail, userID)
	}

	// Get super_admin role ID
	var roleIDStr string
	err = db.Raw("SELECT id::text FROM roles WHERE name = 'super_admin' AND tenant_id IS NULL").Scan(&roleIDStr).Error
	if err != nil {
		return fmt.Errorf("failed to get super_admin role: %w", err)
	}
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse role ID: %w", err)
	}

	// Check if user already has the role
	var existingRoleID string
	err = db.Raw("SELECT id::text FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).Scan(&existingRoleID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing user role: %w", err)
	}

	if existingRoleID != "" {
		fmt.Println("Super admin role already assigned")
	} else {
		// Assign super_admin role to user
		err = db.Exec(`
			INSERT INTO user_roles (id, user_id, role_id, created_at)
			VALUES (?, ?, ?, NOW())
		`, uuid.New(), userID, roleID).Error
		if err != nil {
			return fmt.Errorf("failed to assign role to user: %w", err)
		}
		fmt.Println("Assigned super_admin role to user")
	}

	// Assign staff permissions to super_admin
	if err := assignStaffPermissions(db, roleID); err != nil {
		fmt.Printf("Warning: failed to assign staff permissions: %v\n", err)
	}

	// Get or create branch for staff
	var branchIDStr string
	err = db.Raw("SELECT id::text FROM branches WHERE tenant_id = ? LIMIT 1", tenantID).Scan(&branchIDStr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing branch: %w", err)
	}

	var branchID uuid.UUID
	if branchIDStr != "" {
		branchID, _ = uuid.Parse(branchIDStr)
		fmt.Printf("Branch already exists (ID: %s)\n", branchID)
	} else {
		// Create a branch
		branchID = uuid.New()
		err = db.Exec(`
			INSERT INTO branches (id, tenant_id, name, code, address, phone, status, created_at, updated_at)
			VALUES (?, ?, 'Main Campus', 'MAIN', '123 School Street, Mumbai', '022-12345678', 'active', NOW(), NOW())
		`, branchID, tenantID).Error
		if err != nil {
			return fmt.Errorf("failed to create branch: %w", err)
		}
		fmt.Printf("Created branch: Main Campus (ID: %s)\n", branchID)
	}

	// Seed staff members
	if err := seedStaff(db, tenantID, branchID); err != nil {
		return fmt.Errorf("failed to seed staff: %w", err)
	}

	// Reset RLS bypass
	if err := db.Exec("SET app.bypass_rls = 'false'").Error; err != nil {
		return fmt.Errorf("failed to reset RLS: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Seed Completed Successfully ===")
	fmt.Println()
	fmt.Println("Super Admin Credentials:")
	fmt.Println("-------------------------")
	fmt.Printf("Email:    %s\n", SuperAdminEmail)
	fmt.Printf("Password: %s\n", SuperAdminPassword)
	fmt.Printf("Tenant:   %s (%s)\n", DefaultTenantName, DefaultTenantSlug)
	fmt.Println()
	fmt.Println("WARNING: Change these credentials in production!")
	fmt.Println()

	_ = ctx // silence unused warning
	return nil
}

// assignStaffPermissions assigns staff permissions to a role.
func assignStaffPermissions(db *gorm.DB, roleID uuid.UUID) error {
	staffPermissions := []string{
		"staff:read",
		"staff:create",
		"staff:update",
		"staff:delete",
		"staff:export",
	}

	for _, code := range staffPermissions {
		// Get permission ID
		var permissionIDStr string
		err := db.Raw("SELECT id::text FROM permissions WHERE code = ?", code).Scan(&permissionIDStr).Error
		if err != nil {
			continue // Permission might not exist
		}
		permissionID, _ := uuid.Parse(permissionIDStr)

		// Check if already assigned
		var exists int
		db.Raw("SELECT 1 FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permissionID).Scan(&exists)
		if exists == 1 {
			continue
		}

		// Assign permission
		err = db.Exec(`
			INSERT INTO role_permissions (id, role_id, permission_id, created_at)
			VALUES (?, ?, ?, NOW())
			ON CONFLICT DO NOTHING
		`, uuid.New(), roleID, permissionID).Error
		if err != nil {
			return err
		}
	}

	fmt.Println("Staff permissions assigned to super_admin")
	return nil
}

// seedStaff creates sample staff members for testing.
func seedStaff(db *gorm.DB, tenantID, branchID uuid.UUID) error {
	// Check if staff already exist
	var staffCount int64
	db.Raw("SELECT COUNT(*) FROM staff WHERE tenant_id = ?", tenantID).Scan(&staffCount)
	if staffCount > 0 {
		fmt.Printf("Staff members already exist (%d total)\n", staffCount)
		return nil
	}

	// Initialize or update employee sequence
	err := db.Exec(`
		INSERT INTO staff_employee_sequences (id, tenant_id, prefix, last_sequence, created_at, updated_at)
		VALUES (?, ?, 'EMP', 0, NOW(), NOW())
		ON CONFLICT (tenant_id, prefix) DO NOTHING
	`, uuid.New(), tenantID).Error
	if err != nil {
		return fmt.Errorf("failed to initialize employee sequence: %w", err)
	}

	// Staff data for seeding
	staffMembers := []models.Staff{
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			BranchID:         branchID,
			EmployeeID:       "EMP00001",
			EmployeeIDPrefix: "EMP",
			FirstName:        "Rajesh",
			MiddleName:       "Kumar",
			LastName:         "Sharma",
			DateOfBirth:      time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
			Gender:           models.GenderMale,
			BloodGroup:       "O+",
			Nationality:      "Indian",
			Religion:         "Hindu",
			MaritalStatus:    "married",
			PersonalEmail:    "rajesh.sharma@gmail.com",
			WorkEmail:        "rajesh.sharma@school.com",
			PersonalPhone:    "9876543210",
			WorkPhone:        "9123456789",
			StaffType:        models.StaffTypeTeaching,
			JoinDate:         time.Date(2015, 7, 1, 0, 0, 0, 0, time.UTC),
			Status:           models.StaffStatusActive,
			Bio:              "Senior Mathematics Teacher with 15+ years of experience.",
			CurrentAddressLine1: "123 Main Street",
			CurrentAddressLine2: "Andheri West",
			CurrentCity:      "Mumbai",
			CurrentState:     "Maharashtra",
			CurrentPincode:   "400058",
			CurrentCountry:   "India",
			SameAsCurrent:    true,
			Version:          1,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			BranchID:         branchID,
			EmployeeID:       "EMP00002",
			EmployeeIDPrefix: "EMP",
			FirstName:        "Priya",
			LastName:         "Gupta",
			DateOfBirth:      time.Date(1990, 8, 22, 0, 0, 0, 0, time.UTC),
			Gender:           models.GenderFemale,
			BloodGroup:       "B+",
			Nationality:      "Indian",
			Religion:         "Hindu",
			MaritalStatus:    "single",
			PersonalEmail:    "priya.gupta@gmail.com",
			WorkEmail:        "priya.gupta@school.com",
			PersonalPhone:    "9876543211",
			WorkPhone:        "9123456780",
			StaffType:        models.StaffTypeTeaching,
			JoinDate:         time.Date(2019, 4, 15, 0, 0, 0, 0, time.UTC),
			Status:           models.StaffStatusActive,
			Bio:              "English Literature Teacher passionate about creative writing.",
			CurrentAddressLine1: "456 Park Avenue",
			CurrentCity:      "Mumbai",
			CurrentState:     "Maharashtra",
			CurrentPincode:   "400092",
			CurrentCountry:   "India",
			SameAsCurrent:    true,
			Version:          1,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			BranchID:         branchID,
			EmployeeID:       "EMP00003",
			EmployeeIDPrefix: "EMP",
			FirstName:        "Amit",
			LastName:         "Patel",
			DateOfBirth:      time.Date(1988, 12, 3, 0, 0, 0, 0, time.UTC),
			Gender:           models.GenderMale,
			BloodGroup:       "A+",
			Nationality:      "Indian",
			PersonalEmail:    "amit.patel@gmail.com",
			WorkEmail:        "amit.patel@school.com",
			PersonalPhone:    "9876543212",
			WorkPhone:        "9123456781",
			StaffType:        models.StaffTypeTeaching,
			JoinDate:         time.Date(2018, 8, 1, 0, 0, 0, 0, time.UTC),
			Status:           models.StaffStatusActive,
			Bio:              "Physics teacher and science lab coordinator.",
			CurrentAddressLine1: "789 Science Lane",
			CurrentCity:      "Pune",
			CurrentState:     "Maharashtra",
			CurrentPincode:   "411001",
			CurrentCountry:   "India",
			SameAsCurrent:    false,
			PermanentAddressLine1: "12 Heritage Road",
			PermanentCity:    "Ahmedabad",
			PermanentState:   "Gujarat",
			PermanentPincode: "380001",
			PermanentCountry: "India",
			Version:          1,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			BranchID:         branchID,
			EmployeeID:       "EMP00004",
			EmployeeIDPrefix: "EMP",
			FirstName:        "Sunita",
			MiddleName:       "Devi",
			LastName:         "Verma",
			DateOfBirth:      time.Date(1982, 3, 18, 0, 0, 0, 0, time.UTC),
			Gender:           models.GenderFemale,
			BloodGroup:       "AB+",
			Nationality:      "Indian",
			WorkEmail:        "sunita.verma@school.com",
			WorkPhone:        "9123456782",
			StaffType:        models.StaffTypeNonTeaching,
			JoinDate:         time.Date(2010, 1, 15, 0, 0, 0, 0, time.UTC),
			Status:           models.StaffStatusActive,
			Bio:              "Administrative Officer handling admissions and records.",
			CurrentAddressLine1: "321 Admin Block",
			CurrentCity:      "Mumbai",
			CurrentState:     "Maharashtra",
			CurrentPincode:   "400001",
			CurrentCountry:   "India",
			SameAsCurrent:    true,
			Version:          1,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			BranchID:         branchID,
			EmployeeID:       "EMP00005",
			EmployeeIDPrefix: "EMP",
			FirstName:        "Mohammed",
			LastName:         "Khan",
			DateOfBirth:      time.Date(1992, 7, 25, 0, 0, 0, 0, time.UTC),
			Gender:           models.GenderMale,
			BloodGroup:       "O-",
			Nationality:      "Indian",
			Religion:         "Islam",
			MaritalStatus:    "married",
			PersonalEmail:    "mohammed.khan@gmail.com",
			WorkEmail:        "mohammed.khan@school.com",
			PersonalPhone:    "9876543215",
			WorkPhone:        "9123456785",
			EmergencyContactName: "Fatima Khan",
			EmergencyContactPhone: "9876543299",
			EmergencyContactRelation: "Spouse",
			StaffType:        models.StaffTypeTeaching,
			JoinDate:         time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC),
			Status:           models.StaffStatusActive,
			Bio:              "Computer Science teacher specializing in programming.",
			CurrentAddressLine1: "567 Tech Park",
			CurrentCity:      "Bangalore",
			CurrentState:     "Karnataka",
			CurrentPincode:   "560001",
			CurrentCountry:   "India",
			SameAsCurrent:    true,
			Version:          1,
		},
		{
			ID:               uuid.New(),
			TenantID:         tenantID,
			BranchID:         branchID,
			EmployeeID:       "EMP00006",
			EmployeeIDPrefix: "EMP",
			FirstName:        "Lakshmi",
			LastName:         "Iyer",
			DateOfBirth:      time.Date(1987, 11, 9, 0, 0, 0, 0, time.UTC),
			Gender:           models.GenderFemale,
			BloodGroup:       "B-",
			Nationality:      "Indian",
			WorkEmail:        "lakshmi.iyer@school.com",
			WorkPhone:        "9123456786",
			StaffType:        models.StaffTypeTeaching,
			JoinDate:         time.Date(2017, 3, 1, 0, 0, 0, 0, time.UTC),
			Status:           models.StaffStatusOnLeave,
			StatusReason:     "Maternity leave",
			Bio:              "Chemistry teacher with research background.",
			CurrentAddressLine1: "890 Chemical Valley",
			CurrentCity:      "Chennai",
			CurrentState:     "Tamil Nadu",
			CurrentPincode:   "600001",
			CurrentCountry:   "India",
			SameAsCurrent:    true,
			Version:          1,
		},
	}

	// Insert staff members
	for _, staff := range staffMembers {
		staff.CreatedAt = time.Now()
		staff.UpdatedAt = time.Now()

		err := db.Exec(`
			INSERT INTO staff (
				id, tenant_id, branch_id, employee_id, employee_id_prefix,
				first_name, middle_name, last_name, date_of_birth, gender,
				blood_group, nationality, religion, marital_status,
				personal_email, work_email, personal_phone, work_phone,
				emergency_contact_name, emergency_contact_phone, emergency_contact_relation,
				current_address_line1, current_address_line2, current_city, current_state, current_pincode, current_country,
				permanent_address_line1, permanent_address_line2, permanent_city, permanent_state, permanent_pincode, permanent_country,
				same_as_current, staff_type, join_date, status, status_reason, bio,
				version, created_at, updated_at
			) VALUES (
				?, ?, ?, ?, ?,
				?, ?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?,
				?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?,
				?, NOW(), NOW()
			)
		`,
			staff.ID, staff.TenantID, staff.BranchID, staff.EmployeeID, staff.EmployeeIDPrefix,
			staff.FirstName, staff.MiddleName, staff.LastName, staff.DateOfBirth, staff.Gender,
			staff.BloodGroup, staff.Nationality, staff.Religion, staff.MaritalStatus,
			staff.PersonalEmail, staff.WorkEmail, staff.PersonalPhone, staff.WorkPhone,
			staff.EmergencyContactName, staff.EmergencyContactPhone, staff.EmergencyContactRelation,
			staff.CurrentAddressLine1, staff.CurrentAddressLine2, staff.CurrentCity, staff.CurrentState, staff.CurrentPincode, staff.CurrentCountry,
			staff.PermanentAddressLine1, staff.PermanentAddressLine2, staff.PermanentCity, staff.PermanentState, staff.PermanentPincode, staff.PermanentCountry,
			staff.SameAsCurrent, staff.StaffType, staff.JoinDate, staff.Status, staff.StatusReason, staff.Bio,
			staff.Version,
		).Error
		if err != nil {
			return fmt.Errorf("failed to insert staff %s: %w", staff.EmployeeID, err)
		}
	}

	// Update the sequence
	err = db.Exec(`
		UPDATE staff_employee_sequences
		SET last_sequence = 6, updated_at = NOW()
		WHERE tenant_id = ? AND prefix = 'EMP'
	`, tenantID).Error
	if err != nil {
		return fmt.Errorf("failed to update employee sequence: %w", err)
	}

	fmt.Printf("Created %d staff members\n", len(staffMembers))
	return nil
}
