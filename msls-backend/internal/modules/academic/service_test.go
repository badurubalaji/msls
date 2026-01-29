// Package academic provides academic structure management functionality.
package academic

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"msls-backend/internal/pkg/database/models"
)

// ========================================
// Error Definitions Tests
// ========================================

func TestErrorDefinitions(t *testing.T) {
	t.Run("class errors are defined", func(t *testing.T) {
		assert.NotNil(t, ErrClassNotFound)
		assert.NotNil(t, ErrClassCodeExists)
		assert.NotNil(t, ErrClassHasSections)
		assert.NotNil(t, ErrClassHasStudents)

		assert.Equal(t, "class not found", ErrClassNotFound.Error())
		assert.Equal(t, "class code already exists", ErrClassCodeExists.Error())
		assert.Equal(t, "cannot delete class with sections", ErrClassHasSections.Error())
	})

	t.Run("section errors are defined", func(t *testing.T) {
		assert.NotNil(t, ErrSectionNotFound)
		assert.NotNil(t, ErrSectionCodeExists)
		assert.NotNil(t, ErrSectionHasStudents)

		assert.Equal(t, "section not found", ErrSectionNotFound.Error())
		assert.Equal(t, "section code already exists for this class", ErrSectionCodeExists.Error())
		assert.Equal(t, "cannot delete section with enrolled students", ErrSectionHasStudents.Error())
	})

	t.Run("stream errors are defined", func(t *testing.T) {
		assert.NotNil(t, ErrStreamNotFound)
		assert.NotNil(t, ErrStreamCodeExists)
		assert.NotNil(t, ErrStreamInUse)

		assert.Equal(t, "stream not found", ErrStreamNotFound.Error())
		assert.Equal(t, "stream code already exists", ErrStreamCodeExists.Error())
		assert.Equal(t, "cannot delete stream that is in use", ErrStreamInUse.Error())
	})

	t.Run("general errors are defined", func(t *testing.T) {
		assert.NotNil(t, ErrInvalidBranch)
		assert.NotNil(t, ErrInvalidAcademicYear)
		assert.NotNil(t, ErrInvalidClassTeacher)
	})
}

// ========================================
// Class DTO Tests
// ========================================

func TestClassToResponse(t *testing.T) {
	classID := uuid.New()
	branchID := uuid.New()
	now := time.Now()

	class := &models.Class{
		ID:           classID,
		TenantID:     uuid.New(),
		BranchID:     branchID,
		Name:         "Class 10",
		Code:         "X",
		DisplayOrder: 10,
		Description:  "Senior secondary class",
		HasStreams:   true,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		Branch: models.Branch{
			Name: "Main Branch",
		},
	}
	class.Branch.ID = branchID

	resp := ClassToResponse(class)

	assert.Equal(t, classID, resp.ID)
	assert.Equal(t, branchID, resp.BranchID)
	assert.Equal(t, "Main Branch", resp.BranchName)
	assert.Equal(t, "Class 10", resp.Name)
	assert.Equal(t, "X", resp.Code)
	assert.Equal(t, 10, resp.DisplayOrder)
	assert.Equal(t, "Senior secondary class", resp.Description)
	assert.True(t, resp.HasStreams)
	assert.True(t, resp.IsActive)
}

func TestClassToResponse_WithSections(t *testing.T) {
	classID := uuid.New()
	sectionID := uuid.New()
	now := time.Now()

	class := &models.Class{
		ID:        classID,
		BranchID:  uuid.New(),
		Name:      "Class 10",
		Code:      "X",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
		Sections: []models.Section{
			{
				ID:        sectionID,
				ClassID:   classID,
				Name:      "A",
				Code:      "A",
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	resp := ClassToResponse(class)

	assert.Len(t, resp.Sections, 1)
	assert.Equal(t, sectionID, resp.Sections[0].ID)
	assert.Equal(t, "A", resp.Sections[0].Name)
}

func TestClassToResponse_WithStreams(t *testing.T) {
	now := time.Now()
	streamID := uuid.New()

	class := &models.Class{
		ID:         uuid.New(),
		BranchID:   uuid.New(),
		Name:       "Class 11",
		Code:       "XI",
		HasStreams: true,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
		Streams: []models.Stream{
			{
				ID:        streamID,
				Name:      "Science",
				Code:      "SCI",
				IsActive:  true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	resp := ClassToResponse(class)

	assert.Len(t, resp.Streams, 1)
	assert.Equal(t, streamID, resp.Streams[0].ID)
	assert.Equal(t, "Science", resp.Streams[0].Name)
}

func TestClassToResponse_EmptyBranch(t *testing.T) {
	now := time.Now()

	class := &models.Class{
		ID:        uuid.New(),
		BranchID:  uuid.New(),
		Name:      "Class 10",
		Code:      "X",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := ClassToResponse(class)

	assert.Empty(t, resp.BranchName)
}

// ========================================
// Section DTO Tests
// ========================================

func TestSectionToResponse(t *testing.T) {
	sectionID := uuid.New()
	classID := uuid.New()
	academicYearID := uuid.New()
	streamID := uuid.New()
	teacherID := uuid.New()
	capacity := 35
	now := time.Now()

	section := &models.Section{
		ID:             sectionID,
		TenantID:       uuid.New(),
		ClassID:        classID,
		AcademicYearID: &academicYearID,
		StreamID:       &streamID,
		ClassTeacherID: &teacherID,
		Name:           "A",
		Code:           "XA",
		Capacity:       &capacity,
		RoomNumber:     "101",
		DisplayOrder:   1,
		IsActive:       true,
		StudentCount:   30,
		CreatedAt:      now,
		UpdatedAt:      now,
		Class: models.Class{
			ID:   classID,
			Name: "Class 10",
		},
		AcademicYear: models.AcademicYear{
			Name: "2024-2025",
		},
		Stream: &models.Stream{
			ID:   streamID,
			Name: "Science",
		},
		ClassTeacher: &models.Staff{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	section.AcademicYear.ID = academicYearID

	resp := SectionToResponse(section)

	assert.Equal(t, sectionID, resp.ID)
	assert.Equal(t, classID, resp.ClassID)
	assert.Equal(t, "Class 10", resp.ClassName)
	require.NotNil(t, resp.AcademicYearID)
	assert.Equal(t, academicYearID.String(), *resp.AcademicYearID)
	assert.Equal(t, "2024-2025", resp.AcademicYearName)
	require.NotNil(t, resp.StreamID)
	assert.Equal(t, streamID.String(), *resp.StreamID)
	assert.Equal(t, "Science", resp.StreamName)
	require.NotNil(t, resp.ClassTeacherID)
	assert.Equal(t, teacherID.String(), *resp.ClassTeacherID)
	assert.Equal(t, "John Doe", resp.ClassTeacherName)
	assert.Equal(t, "A", resp.Name)
	assert.Equal(t, "XA", resp.Code)
	assert.Equal(t, 35, resp.Capacity)
	assert.Equal(t, "101", resp.RoomNumber)
	assert.Equal(t, 30, resp.StudentCount)
	assert.True(t, resp.IsActive)
}

func TestSectionToResponse_DefaultCapacity(t *testing.T) {
	now := time.Now()

	section := &models.Section{
		ID:        uuid.New(),
		ClassID:   uuid.New(),
		Name:      "A",
		Code:      "A",
		Capacity:  nil, // No capacity set
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := SectionToResponse(section)

	assert.Equal(t, 40, resp.Capacity, "default capacity should be 40")
}

func TestSectionToResponse_WithoutOptionalFields(t *testing.T) {
	now := time.Now()

	section := &models.Section{
		ID:        uuid.New(),
		ClassID:   uuid.New(),
		Name:      "A",
		Code:      "A",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := SectionToResponse(section)

	assert.Nil(t, resp.AcademicYearID)
	assert.Empty(t, resp.AcademicYearName)
	assert.Nil(t, resp.StreamID)
	assert.Empty(t, resp.StreamName)
	assert.Nil(t, resp.ClassTeacherID)
	assert.Empty(t, resp.ClassTeacherName)
}

func TestSectionToResponse_ClassTeacherNoLastName(t *testing.T) {
	now := time.Now()
	teacherID := uuid.New()

	section := &models.Section{
		ID:             uuid.New(),
		ClassID:        uuid.New(),
		Name:           "A",
		Code:           "A",
		ClassTeacherID: &teacherID,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
		ClassTeacher: &models.Staff{
			FirstName: "John",
			LastName:  "", // No last name
		},
	}

	resp := SectionToResponse(section)

	assert.Equal(t, "John", resp.ClassTeacherName)
}

// ========================================
// Stream DTO Tests
// ========================================

func TestStreamToResponse(t *testing.T) {
	streamID := uuid.New()
	now := time.Now()

	stream := &models.Stream{
		ID:           streamID,
		TenantID:     uuid.New(),
		Name:         "Science",
		Code:         "SCI",
		Description:  "Science stream for senior classes",
		DisplayOrder: 1,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	resp := StreamToResponse(stream)

	assert.Equal(t, streamID, resp.ID)
	assert.Equal(t, "Science", resp.Name)
	assert.Equal(t, "SCI", resp.Code)
	assert.Equal(t, "Science stream for senior classes", resp.Description)
	assert.Equal(t, 1, resp.DisplayOrder)
	assert.True(t, resp.IsActive)
}

func TestStreamToResponse_EmptyDescription(t *testing.T) {
	now := time.Now()

	stream := &models.Stream{
		ID:        uuid.New(),
		Name:      "Commerce",
		Code:      "COM",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := StreamToResponse(stream)

	assert.Empty(t, resp.Description)
}

// ========================================
// CreateClassRequest Tests
// ========================================

func TestCreateClassRequest_Fields(t *testing.T) {
	branchID := uuid.New()
	streamIDs := []uuid.UUID{uuid.New(), uuid.New()}

	req := CreateClassRequest{
		BranchID:     branchID,
		Name:         "Class 10",
		Code:         "X",
		DisplayOrder: 10,
		Description:  "Senior class",
		HasStreams:   true,
		StreamIDs:    streamIDs,
	}

	assert.Equal(t, branchID, req.BranchID)
	assert.Equal(t, "Class 10", req.Name)
	assert.Equal(t, "X", req.Code)
	assert.Equal(t, 10, req.DisplayOrder)
	assert.Equal(t, "Senior class", req.Description)
	assert.True(t, req.HasStreams)
	assert.Len(t, req.StreamIDs, 2)
}

func TestCreateClassRequest_MinimalFields(t *testing.T) {
	branchID := uuid.New()

	req := CreateClassRequest{
		BranchID: branchID,
		Name:     "Class 1",
		Code:     "I",
	}

	assert.Equal(t, branchID, req.BranchID)
	assert.Equal(t, "Class 1", req.Name)
	assert.Equal(t, "I", req.Code)
	assert.False(t, req.HasStreams)
	assert.Nil(t, req.StreamIDs)
}

// ========================================
// UpdateClassRequest Tests
// ========================================

func TestUpdateClassRequest_AllFields(t *testing.T) {
	name := "Updated Class"
	code := "UC"
	displayOrder := 5
	description := "Updated description"
	hasStreams := true
	isActive := false
	streamIDs := []uuid.UUID{uuid.New()}

	req := UpdateClassRequest{
		Name:         &name,
		Code:         &code,
		DisplayOrder: &displayOrder,
		Description:  &description,
		HasStreams:   &hasStreams,
		IsActive:     &isActive,
		StreamIDs:    streamIDs,
	}

	assert.Equal(t, "Updated Class", *req.Name)
	assert.Equal(t, "UC", *req.Code)
	assert.Equal(t, 5, *req.DisplayOrder)
	assert.Equal(t, "Updated description", *req.Description)
	assert.True(t, *req.HasStreams)
	assert.False(t, *req.IsActive)
	assert.Len(t, req.StreamIDs, 1)
}

func TestUpdateClassRequest_PartialUpdate(t *testing.T) {
	name := "Updated Name Only"

	req := UpdateClassRequest{
		Name: &name,
	}

	assert.NotNil(t, req.Name)
	assert.Nil(t, req.Code)
	assert.Nil(t, req.DisplayOrder)
	assert.Nil(t, req.Description)
	assert.Nil(t, req.HasStreams)
	assert.Nil(t, req.IsActive)
}

// ========================================
// CreateSectionRequest Tests
// ========================================

func TestCreateSectionRequest_AllFields(t *testing.T) {
	classID := uuid.New()
	academicYearID := uuid.New()
	streamID := uuid.New()
	teacherID := uuid.New()

	req := CreateSectionRequest{
		ClassID:        classID,
		AcademicYearID: &academicYearID,
		StreamID:       &streamID,
		ClassTeacherID: &teacherID,
		Name:           "A",
		Code:           "XA",
		Capacity:       35,
		RoomNumber:     "101",
	}

	assert.Equal(t, classID, req.ClassID)
	assert.Equal(t, academicYearID, *req.AcademicYearID)
	assert.Equal(t, streamID, *req.StreamID)
	assert.Equal(t, teacherID, *req.ClassTeacherID)
	assert.Equal(t, "A", req.Name)
	assert.Equal(t, "XA", req.Code)
	assert.Equal(t, 35, req.Capacity)
	assert.Equal(t, "101", req.RoomNumber)
}

func TestCreateSectionRequest_MinimalFields(t *testing.T) {
	classID := uuid.New()

	req := CreateSectionRequest{
		ClassID: classID,
		Name:    "A",
		Code:    "A",
	}

	assert.Equal(t, classID, req.ClassID)
	assert.Equal(t, "A", req.Name)
	assert.Equal(t, "A", req.Code)
	assert.Nil(t, req.AcademicYearID)
	assert.Nil(t, req.StreamID)
	assert.Nil(t, req.ClassTeacherID)
	assert.Equal(t, 0, req.Capacity)
}

// ========================================
// UpdateSectionRequest Tests
// ========================================

func TestUpdateSectionRequest_AllFields(t *testing.T) {
	academicYearID := uuid.New()
	streamID := uuid.New()
	teacherID := uuid.New()
	name := "B"
	code := "XB"
	capacity := 40
	roomNumber := "102"
	displayOrder := 2
	isActive := true

	req := UpdateSectionRequest{
		AcademicYearID: &academicYearID,
		StreamID:       &streamID,
		ClassTeacherID: &teacherID,
		Name:           &name,
		Code:           &code,
		Capacity:       &capacity,
		RoomNumber:     &roomNumber,
		DisplayOrder:   &displayOrder,
		IsActive:       &isActive,
	}

	assert.Equal(t, academicYearID, *req.AcademicYearID)
	assert.Equal(t, streamID, *req.StreamID)
	assert.Equal(t, teacherID, *req.ClassTeacherID)
	assert.Equal(t, "B", *req.Name)
	assert.Equal(t, "XB", *req.Code)
	assert.Equal(t, 40, *req.Capacity)
	assert.Equal(t, "102", *req.RoomNumber)
	assert.Equal(t, 2, *req.DisplayOrder)
	assert.True(t, *req.IsActive)
}

// ========================================
// CreateStreamRequest Tests
// ========================================

func TestCreateStreamRequest_AllFields(t *testing.T) {
	req := CreateStreamRequest{
		Name:         "Science",
		Code:         "SCI",
		Description:  "Science stream",
		DisplayOrder: 1,
	}

	assert.Equal(t, "Science", req.Name)
	assert.Equal(t, "SCI", req.Code)
	assert.Equal(t, "Science stream", req.Description)
	assert.Equal(t, 1, req.DisplayOrder)
}

func TestCreateStreamRequest_MinimalFields(t *testing.T) {
	req := CreateStreamRequest{
		Name: "Arts",
		Code: "ART",
	}

	assert.Equal(t, "Arts", req.Name)
	assert.Equal(t, "ART", req.Code)
	assert.Empty(t, req.Description)
	assert.Equal(t, 0, req.DisplayOrder)
}

// ========================================
// UpdateStreamRequest Tests
// ========================================

func TestUpdateStreamRequest_AllFields(t *testing.T) {
	name := "Updated Science"
	code := "USCI"
	description := "Updated description"
	displayOrder := 2
	isActive := false

	req := UpdateStreamRequest{
		Name:         &name,
		Code:         &code,
		Description:  &description,
		DisplayOrder: &displayOrder,
		IsActive:     &isActive,
	}

	assert.Equal(t, "Updated Science", *req.Name)
	assert.Equal(t, "USCI", *req.Code)
	assert.Equal(t, "Updated description", *req.Description)
	assert.Equal(t, 2, *req.DisplayOrder)
	assert.False(t, *req.IsActive)
}

// ========================================
// ClassFilter Tests
// ========================================

func TestClassFilter_AllFields(t *testing.T) {
	tenantID := uuid.New()
	branchID := uuid.New()
	isActive := true
	hasStreams := false

	filter := ClassFilter{
		TenantID:   tenantID,
		BranchID:   &branchID,
		IsActive:   &isActive,
		HasStreams: &hasStreams,
		Search:     "Class",
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.Equal(t, branchID, *filter.BranchID)
	assert.True(t, *filter.IsActive)
	assert.False(t, *filter.HasStreams)
	assert.Equal(t, "Class", filter.Search)
}

func TestClassFilter_MinimalFields(t *testing.T) {
	tenantID := uuid.New()

	filter := ClassFilter{
		TenantID: tenantID,
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.Nil(t, filter.BranchID)
	assert.Nil(t, filter.IsActive)
	assert.Nil(t, filter.HasStreams)
	assert.Empty(t, filter.Search)
}

// ========================================
// SectionFilter Tests
// ========================================

func TestSectionFilter_AllFields(t *testing.T) {
	tenantID := uuid.New()
	classID := uuid.New()
	academicYearID := uuid.New()
	streamID := uuid.New()
	isActive := true

	filter := SectionFilter{
		TenantID:       tenantID,
		ClassID:        &classID,
		AcademicYearID: &academicYearID,
		StreamID:       &streamID,
		IsActive:       &isActive,
		Search:         "A",
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.Equal(t, classID, *filter.ClassID)
	assert.Equal(t, academicYearID, *filter.AcademicYearID)
	assert.Equal(t, streamID, *filter.StreamID)
	assert.True(t, *filter.IsActive)
	assert.Equal(t, "A", filter.Search)
}

func TestSectionFilter_MinimalFields(t *testing.T) {
	tenantID := uuid.New()

	filter := SectionFilter{
		TenantID: tenantID,
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.Nil(t, filter.ClassID)
	assert.Nil(t, filter.AcademicYearID)
	assert.Nil(t, filter.StreamID)
	assert.Nil(t, filter.IsActive)
	assert.Empty(t, filter.Search)
}

// ========================================
// StreamFilter Tests
// ========================================

func TestStreamFilter_AllFields(t *testing.T) {
	tenantID := uuid.New()
	isActive := true

	filter := StreamFilter{
		TenantID: tenantID,
		IsActive: &isActive,
		Search:   "Science",
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.True(t, *filter.IsActive)
	assert.Equal(t, "Science", filter.Search)
}

func TestStreamFilter_MinimalFields(t *testing.T) {
	tenantID := uuid.New()

	filter := StreamFilter{
		TenantID: tenantID,
	}

	assert.Equal(t, tenantID, filter.TenantID)
	assert.Nil(t, filter.IsActive)
	assert.Empty(t, filter.Search)
}

// ========================================
// ClassStructure Response Tests
// ========================================

func TestClassWithSectionsResponse_Fields(t *testing.T) {
	classID := uuid.New()
	sectionID := uuid.New()
	teacherID := uuid.New().String()

	sections := []SectionStructureResponse{
		{
			ID:               sectionID,
			Name:             "A",
			Code:             "XA",
			Capacity:         40,
			StudentCount:     30,
			ClassTeacherID:   &teacherID,
			ClassTeacherName: "John Doe",
			StreamName:       "Science",
			RoomNumber:       "101",
			CapacityUsage:    75.0,
		},
	}

	resp := ClassWithSectionsResponse{
		ID:            classID,
		Name:          "Class 10",
		Code:          "X",
		DisplayOrder:  10,
		HasStreams:    false,
		IsActive:      true,
		Sections:      sections,
		TotalStudents: 30,
		TotalCapacity: 40,
	}

	assert.Equal(t, classID, resp.ID)
	assert.Equal(t, "Class 10", resp.Name)
	assert.Equal(t, "X", resp.Code)
	assert.Equal(t, 10, resp.DisplayOrder)
	assert.False(t, resp.HasStreams)
	assert.True(t, resp.IsActive)
	assert.Len(t, resp.Sections, 1)
	assert.Equal(t, 30, resp.TotalStudents)
	assert.Equal(t, 40, resp.TotalCapacity)
}

func TestSectionStructureResponse_Fields(t *testing.T) {
	sectionID := uuid.New()
	teacherID := uuid.New().String()

	resp := SectionStructureResponse{
		ID:               sectionID,
		Name:             "A",
		Code:             "XA",
		Capacity:         40,
		StudentCount:     30,
		ClassTeacherID:   &teacherID,
		ClassTeacherName: "John Doe",
		StreamName:       "Science",
		RoomNumber:       "101",
		CapacityUsage:    75.0,
	}

	assert.Equal(t, sectionID, resp.ID)
	assert.Equal(t, "A", resp.Name)
	assert.Equal(t, "XA", resp.Code)
	assert.Equal(t, 40, resp.Capacity)
	assert.Equal(t, 30, resp.StudentCount)
	require.NotNil(t, resp.ClassTeacherID)
	assert.Equal(t, teacherID, *resp.ClassTeacherID)
	assert.Equal(t, "John Doe", resp.ClassTeacherName)
	assert.Equal(t, "Science", resp.StreamName)
	assert.Equal(t, "101", resp.RoomNumber)
	assert.Equal(t, 75.0, resp.CapacityUsage)
}

func TestSectionStructureResponse_WithoutOptionalFields(t *testing.T) {
	sectionID := uuid.New()

	resp := SectionStructureResponse{
		ID:           sectionID,
		Name:         "A",
		Code:         "A",
		Capacity:     40,
		StudentCount: 0,
	}

	assert.Nil(t, resp.ClassTeacherID)
	assert.Empty(t, resp.ClassTeacherName)
	assert.Empty(t, resp.StreamName)
	assert.Empty(t, resp.RoomNumber)
}

// ========================================
// ClassStructureResponse Tests
// ========================================

func TestClassStructureResponse_Fields(t *testing.T) {
	classID := uuid.New()

	resp := ClassStructureResponse{
		Classes: []ClassWithSectionsResponse{
			{
				ID:            classID,
				Name:          "Class 10",
				Code:          "X",
				TotalStudents: 100,
				TotalCapacity: 120,
				Sections:      []SectionStructureResponse{},
			},
		},
	}

	assert.Len(t, resp.Classes, 1)
	assert.Equal(t, classID, resp.Classes[0].ID)
}

// ========================================
// ClassListResponse Tests
// ========================================

func TestClassListResponse_Fields(t *testing.T) {
	classID := uuid.New()
	now := time.Now()

	classes := []ClassResponse{
		{
			ID:           classID,
			BranchID:     uuid.New(),
			BranchName:   "Main Branch",
			Name:         "Class 10",
			Code:         "X",
			DisplayOrder: 10,
			HasStreams:   false,
			IsActive:     true,
			CreatedAt:    now.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    now.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	resp := ClassListResponse{
		Classes: classes,
		Total:   1,
	}

	assert.Len(t, resp.Classes, 1)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, classID, resp.Classes[0].ID)
}

// ========================================
// SectionListResponse Tests
// ========================================

func TestSectionListResponse_Fields(t *testing.T) {
	sectionID := uuid.New()
	now := time.Now()

	sections := []SectionResponse{
		{
			ID:       sectionID,
			ClassID:  uuid.New(),
			Name:     "A",
			Code:     "XA",
			Capacity: 40,
			IsActive: true,
			CreatedAt: now.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: now.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	resp := SectionListResponse{
		Sections: sections,
		Total:    1,
	}

	assert.Len(t, resp.Sections, 1)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, sectionID, resp.Sections[0].ID)
}

// ========================================
// StreamListResponse Tests
// ========================================

func TestStreamListResponse_Fields(t *testing.T) {
	streamID := uuid.New()
	now := time.Now()

	streams := []StreamResponse{
		{
			ID:           streamID,
			Name:         "Science",
			Code:         "SCI",
			DisplayOrder: 1,
			IsActive:     true,
			CreatedAt:    now.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    now.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	resp := StreamListResponse{
		Streams: streams,
		Total:   1,
	}

	assert.Len(t, resp.Streams, 1)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, streamID, resp.Streams[0].ID)
}

// ========================================
// Timestamp Format Tests
// ========================================

func TestClassToResponse_TimestampFormat(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	class := &models.Class{
		ID:        uuid.New(),
		BranchID:  uuid.New(),
		Name:      "Class 10",
		Code:      "X",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := ClassToResponse(class)

	assert.Equal(t, "2024-06-15T10:30:00Z", resp.CreatedAt)
	assert.Equal(t, "2024-06-15T10:30:00Z", resp.UpdatedAt)
}

func TestSectionToResponse_TimestampFormat(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	section := &models.Section{
		ID:        uuid.New(),
		ClassID:   uuid.New(),
		Name:      "A",
		Code:      "A",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := SectionToResponse(section)

	assert.Equal(t, "2024-06-15T10:30:00Z", resp.CreatedAt)
	assert.Equal(t, "2024-06-15T10:30:00Z", resp.UpdatedAt)
}

func TestStreamToResponse_TimestampFormat(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	stream := &models.Stream{
		ID:        uuid.New(),
		Name:      "Science",
		Code:      "SCI",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := StreamToResponse(stream)

	assert.Equal(t, "2024-06-15T10:30:00Z", resp.CreatedAt)
	assert.Equal(t, "2024-06-15T10:30:00Z", resp.UpdatedAt)
}
