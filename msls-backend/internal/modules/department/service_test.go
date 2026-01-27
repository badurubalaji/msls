package department

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"msls-backend/internal/pkg/database/models"
)

func TestToDepartmentResponse(t *testing.T) {
	tenantID := uuid.New()
	branchID := uuid.New()
	deptID := uuid.New()
	headID := uuid.New()
	description := "Test Description"
	now := time.Now()

	t.Run("with all fields", func(t *testing.T) {
		dept := &models.Department{
			ID:          deptID,
			TenantID:    tenantID,
			BranchID:    branchID,
			Name:        "Science Department",
			Code:        "SCI",
			Description: &description,
			HeadID:      &headID,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
			Branch: &models.Branch{
				Name: "Main Campus",
			},
			Head: &models.Staff{
				FirstName: "John",
				LastName:  "Doe",
			},
		}

		resp := ToDepartmentResponse(dept, 5)

		assert.Equal(t, deptID.String(), resp.ID)
		assert.Equal(t, "Science Department", resp.Name)
		assert.Equal(t, "SCI", resp.Code)
		assert.Equal(t, description, resp.Description)
		assert.Equal(t, branchID.String(), resp.BranchID)
		assert.Equal(t, "Main Campus", resp.BranchName)
		assert.Equal(t, headID.String(), resp.HeadID)
		assert.Equal(t, "John Doe", resp.HeadName)
		assert.True(t, resp.IsActive)
		assert.Equal(t, 5, resp.StaffCount)
	})

	t.Run("without optional fields", func(t *testing.T) {
		dept := &models.Department{
			ID:        deptID,
			TenantID:  tenantID,
			BranchID:  branchID,
			Name:      "Admin Department",
			Code:      "ADMIN",
			IsActive:  false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		resp := ToDepartmentResponse(dept, 0)

		assert.Equal(t, deptID.String(), resp.ID)
		assert.Equal(t, "Admin Department", resp.Name)
		assert.Equal(t, "ADMIN", resp.Code)
		assert.Empty(t, resp.Description)
		assert.Empty(t, resp.HeadID)
		assert.Empty(t, resp.HeadName)
		assert.Empty(t, resp.BranchName)
		assert.False(t, resp.IsActive)
		assert.Equal(t, 0, resp.StaffCount)
	})
}

func TestToDepartmentResponses(t *testing.T) {
	tenantID := uuid.New()
	branchID := uuid.New()
	dept1ID := uuid.New()
	dept2ID := uuid.New()
	now := time.Now()

	departments := []models.Department{
		{
			ID:        dept1ID,
			TenantID:  tenantID,
			BranchID:  branchID,
			Name:      "Science",
			Code:      "SCI",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        dept2ID,
			TenantID:  tenantID,
			BranchID:  branchID,
			Name:      "Arts",
			Code:      "ART",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	staffCounts := map[uuid.UUID]int{
		dept1ID: 10,
		dept2ID: 5,
	}

	responses := ToDepartmentResponses(departments, staffCounts)

	assert.Len(t, responses, 2)
	assert.Equal(t, "Science", responses[0].Name)
	assert.Equal(t, 10, responses[0].StaffCount)
	assert.Equal(t, "Arts", responses[1].Name)
	assert.Equal(t, 5, responses[1].StaffCount)
}

func TestToDropdownItem(t *testing.T) {
	deptID := uuid.New()
	dept := &models.Department{
		ID:   deptID,
		Name: "Science Department",
	}

	item := ToDropdownItem(dept)

	assert.Equal(t, deptID.String(), item.ID)
	assert.Equal(t, "Science Department", item.Name)
}

func TestToDropdownItems(t *testing.T) {
	departments := []models.Department{
		{ID: uuid.New(), Name: "Science"},
		{ID: uuid.New(), Name: "Arts"},
		{ID: uuid.New(), Name: "Admin"},
	}

	items := ToDropdownItems(departments)

	assert.Len(t, items, 3)
	assert.Equal(t, "Science", items[0].Name)
	assert.Equal(t, "Arts", items[1].Name)
	assert.Equal(t, "Admin", items[2].Name)
}

func TestCreateDepartmentDTO(t *testing.T) {
	tenantID := uuid.New()
	branchID := uuid.New()
	description := "Test department"

	dto := CreateDepartmentDTO{
		TenantID:    tenantID,
		BranchID:    branchID,
		Name:        "Science Department",
		Code:        "SCI",
		Description: &description,
		IsActive:    true,
	}

	assert.Equal(t, tenantID, dto.TenantID)
	assert.Equal(t, branchID, dto.BranchID)
	assert.Equal(t, "Science Department", dto.Name)
	assert.Equal(t, "SCI", dto.Code)
	assert.Equal(t, &description, dto.Description)
	assert.True(t, dto.IsActive)
}

func TestUpdateDepartmentDTO(t *testing.T) {
	name := "Updated Department"
	code := "UPD"
	isActive := false

	dto := UpdateDepartmentDTO{
		Name:     &name,
		Code:     &code,
		IsActive: &isActive,
	}

	assert.Equal(t, &name, dto.Name)
	assert.Equal(t, &code, dto.Code)
	assert.Equal(t, &isActive, dto.IsActive)
}

func TestListFilter(t *testing.T) {
	tenantID := uuid.New()
	branchID := uuid.New()
	isActive := true

	filter := ListFilter{
		TenantID: tenantID,
		BranchID: &branchID,
		IsActive: &isActive,
		Search:   "science",
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.Equal(t, &branchID, filter.BranchID)
	assert.Equal(t, &isActive, filter.IsActive)
	assert.Equal(t, "science", filter.Search)
}
