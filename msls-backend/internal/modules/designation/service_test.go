package designation

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"msls-backend/internal/pkg/database/models"
)

func TestToDesignationResponse(t *testing.T) {
	tenantID := uuid.New()
	desigID := uuid.New()
	deptID := uuid.New()
	now := time.Now()

	t.Run("with department", func(t *testing.T) {
		desig := &models.Designation{
			ID:           desigID,
			TenantID:     tenantID,
			Name:         "Senior Teacher",
			Level:        3,
			DepartmentID: &deptID,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
			Department: &models.Department{
				Name: "Science Department",
			},
		}

		resp := ToDesignationResponse(desig, 10)

		assert.Equal(t, desigID.String(), resp.ID)
		assert.Equal(t, "Senior Teacher", resp.Name)
		assert.Equal(t, 3, resp.Level)
		assert.Equal(t, deptID.String(), resp.DepartmentID)
		assert.Equal(t, "Science Department", resp.DepartmentName)
		assert.True(t, resp.IsActive)
		assert.Equal(t, 10, resp.StaffCount)
	})

	t.Run("without department", func(t *testing.T) {
		desig := &models.Designation{
			ID:        desigID,
			TenantID:  tenantID,
			Name:      "Principal",
			Level:     1,
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		resp := ToDesignationResponse(desig, 1)

		assert.Equal(t, desigID.String(), resp.ID)
		assert.Equal(t, "Principal", resp.Name)
		assert.Equal(t, 1, resp.Level)
		assert.Empty(t, resp.DepartmentID)
		assert.Empty(t, resp.DepartmentName)
		assert.True(t, resp.IsActive)
		assert.Equal(t, 1, resp.StaffCount)
	})
}

func TestToDesignationResponses(t *testing.T) {
	tenantID := uuid.New()
	desig1ID := uuid.New()
	desig2ID := uuid.New()
	now := time.Now()

	designations := []models.Designation{
		{
			ID:        desig1ID,
			TenantID:  tenantID,
			Name:      "Principal",
			Level:     1,
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        desig2ID,
			TenantID:  tenantID,
			Name:      "Vice Principal",
			Level:     2,
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	staffCounts := map[uuid.UUID]int{
		desig1ID: 1,
		desig2ID: 2,
	}

	responses := ToDesignationResponses(designations, staffCounts)

	assert.Len(t, responses, 2)
	assert.Equal(t, "Principal", responses[0].Name)
	assert.Equal(t, 1, responses[0].Level)
	assert.Equal(t, 1, responses[0].StaffCount)
	assert.Equal(t, "Vice Principal", responses[1].Name)
	assert.Equal(t, 2, responses[1].Level)
	assert.Equal(t, 2, responses[1].StaffCount)
}

func TestToDropdownItem_Designation(t *testing.T) {
	desigID := uuid.New()
	desig := &models.Designation{
		ID:    desigID,
		Name:  "Senior Teacher",
		Level: 3,
	}

	item := ToDropdownItem(desig)

	assert.Equal(t, desigID.String(), item.ID)
	assert.Equal(t, "Senior Teacher", item.Name)
	assert.Equal(t, 3, item.Level)
}

func TestToDropdownItems_Designation(t *testing.T) {
	designations := []models.Designation{
		{ID: uuid.New(), Name: "Principal", Level: 1},
		{ID: uuid.New(), Name: "Vice Principal", Level: 2},
		{ID: uuid.New(), Name: "Teacher", Level: 5},
	}

	items := ToDropdownItems(designations)

	assert.Len(t, items, 3)
	assert.Equal(t, "Principal", items[0].Name)
	assert.Equal(t, 1, items[0].Level)
	assert.Equal(t, "Vice Principal", items[1].Name)
	assert.Equal(t, 2, items[1].Level)
	assert.Equal(t, "Teacher", items[2].Name)
	assert.Equal(t, 5, items[2].Level)
}

func TestCreateDesignationDTO(t *testing.T) {
	tenantID := uuid.New()
	deptID := uuid.New()

	dto := CreateDesignationDTO{
		TenantID:     tenantID,
		Name:         "Senior Teacher",
		Level:        3,
		DepartmentID: &deptID,
		IsActive:     true,
	}

	assert.Equal(t, tenantID, dto.TenantID)
	assert.Equal(t, "Senior Teacher", dto.Name)
	assert.Equal(t, 3, dto.Level)
	assert.Equal(t, &deptID, dto.DepartmentID)
	assert.True(t, dto.IsActive)
}

func TestUpdateDesignationDTO(t *testing.T) {
	name := "Updated Designation"
	level := 4
	isActive := false

	dto := UpdateDesignationDTO{
		Name:     &name,
		Level:    &level,
		IsActive: &isActive,
	}

	assert.Equal(t, &name, dto.Name)
	assert.Equal(t, &level, dto.Level)
	assert.Equal(t, &isActive, dto.IsActive)
}

func TestListFilter_Designation(t *testing.T) {
	tenantID := uuid.New()
	deptID := uuid.New()
	isActive := true

	filter := ListFilter{
		TenantID:     tenantID,
		DepartmentID: &deptID,
		IsActive:     &isActive,
		Search:       "teacher",
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.Equal(t, &deptID, filter.DepartmentID)
	assert.Equal(t, &isActive, filter.IsActive)
	assert.Equal(t, "teacher", filter.Search)
}

func TestDesignationLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected string
	}{
		{"Principal", 1, "highest"},
		{"Vice Principal", 2, "high"},
		{"HOD", 3, "mid-high"},
		{"Senior Teacher", 4, "mid"},
		{"Teacher", 5, "mid"},
		{"Junior Teacher", 6, "low"},
		{"Assistant", 10, "lowest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Lower level number = higher in hierarchy
			assert.GreaterOrEqual(t, tt.level, 1)
			assert.LessOrEqual(t, tt.level, 10)
		})
	}
}
