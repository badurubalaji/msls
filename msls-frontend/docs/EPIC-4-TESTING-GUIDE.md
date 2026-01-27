# Epic 4: Student Lifecycle Management - Manual Testing Guide

## Prerequisites

### 1. Start Backend Server
```bash
cd /home/ashulabs/workspace/msls/msls-backend
go run cmd/api/main.go
```

### 2. Start Frontend Server
```bash
cd /home/ashulabs/workspace/msls/msls-frontend
npm start
```

### 3. Login Credentials
```
URL:      http://localhost:4200
Email:    admin@myschoolsystem.com
Password: Admin@123
```

### 4. Required Setup Data
Before testing Epic 4, ensure you have:
- At least one **Academic Year** (Settings > Academic Years)
- Optionally: Classes and Sections (can use placeholder UUIDs for now)

---

## Story 4.1: Student Profile Creation

### Test Case 4.1.1: Create New Student (Minimum Required Fields)
1. Navigate to **Students** from sidebar
2. Click **"Add Student"** button
3. Fill in required fields:
   - First Name: `John`
   - Last Name: `Doe`
   - Date of Birth: `2015-05-15`
   - Gender: `Male`
   - Admission Date: `2024-01-15`
4. Click **"Save"**
5. **Expected**: Student created successfully, redirected to student list

### Test Case 4.1.2: Create Student with All Fields
1. Click **"Add Student"**
2. Fill all fields:
   - **Basic Info**: First Name, Middle Name, Last Name, DOB, Gender
   - **Contact**: Email, Phone, Address (Street, City, State, Postal Code, Country)
   - **Academic**: Admission Number (auto-generated), Admission Date
   - **Additional**: Blood Group, Religion, Nationality, Photo URL
3. Click **"Save"**
4. **Expected**: Student created with all details

### Test Case 4.1.3: Validation Errors
1. Click **"Add Student"**
2. Leave all fields empty
3. Click **"Save"**
4. **Expected**: Validation errors shown for required fields

### Test Case 4.1.4: Edit Student
1. From student list, click on a student row
2. Click **"Edit"** button
3. Modify some fields (e.g., change phone number)
4. Click **"Save"**
5. **Expected**: Changes saved, updated data visible

### Test Case 4.1.5: View Student Details
1. From student list, click on a student row
2. **Expected**: Student detail page shows all information in organized sections

---

## Story 4.2: Guardian Information Management

### Test Case 4.2.1: Add Primary Guardian
1. Navigate to a student's detail page
2. Go to **Guardians** section/tab
3. Click **"Add Guardian"**
4. Fill in:
   - Relationship: `Father`
   - First Name: `Robert`
   - Last Name: `Doe`
   - Phone: `+91-9876543210`
   - Email: `robert.doe@email.com`
   - Is Primary: `Yes`
5. Click **"Save"**
6. **Expected**: Guardian added, marked as primary

### Test Case 4.2.2: Add Secondary Guardian
1. Click **"Add Guardian"** again
2. Fill in:
   - Relationship: `Mother`
   - First Name: `Jane`
   - Last Name: `Doe`
   - Phone: `+91-9876543211`
   - Is Primary: `No`
3. Click **"Save"**
4. **Expected**: Second guardian added

### Test Case 4.2.3: Edit Guardian
1. Click edit icon on an existing guardian
2. Change phone number
3. Save
4. **Expected**: Guardian updated

### Test Case 4.2.4: Delete Guardian
1. Click delete icon on a guardian
2. Confirm deletion
3. **Expected**: Guardian removed from list

### Test Case 4.2.5: Change Primary Guardian
1. Edit secondary guardian
2. Set as Primary: `Yes`
3. Save
4. **Expected**: Previous primary guardian unmarked, new one marked as primary

---

## Story 4.3: Student Health Records

### Test Case 4.3.1: Create/Update Health Profile
1. Navigate to student detail page
2. Go to **Health** section/tab
3. Fill in health profile:
   - Blood Group: `O+`
   - Height: `140` cm
   - Weight: `35` kg
   - Vision (Left/Right): `6/6`
   - Hearing Status: `Normal`
   - Medical Notes: `No known issues`
4. Save
5. **Expected**: Health profile saved

### Test Case 4.3.2: Add Allergy
1. In Health section, find **Allergies** subsection
2. Click **"Add Allergy"**
3. Fill in:
   - Allergen: `Peanuts`
   - Severity: `Severe`
   - Symptoms: `Swelling, difficulty breathing`
   - Treatment: `Epinephrine injection`
4. Save
5. **Expected**: Allergy added to list

### Test Case 4.3.3: Add Chronic Condition
1. Find **Chronic Conditions** subsection
2. Click **"Add Condition"**
3. Fill in:
   - Condition Name: `Asthma`
   - Diagnosed Date: `2020-01-01`
   - Severity: `Mild`
   - Medications: `Inhaler as needed`
4. Save
5. **Expected**: Condition added

### Test Case 4.3.4: Add Vaccination Record
1. Find **Vaccinations** subsection
2. Click **"Add Vaccination"**
3. Fill in:
   - Vaccine Name: `MMR`
   - Date Administered: `2016-06-15`
   - Administered By: `City Hospital`
   - Next Due Date: `2026-06-15`
4. Save
5. **Expected**: Vaccination record added

### Test Case 4.3.5: Record Medical Incident
1. Find **Medical Incidents** subsection
2. Click **"Add Incident"**
3. Fill in:
   - Incident Date: Today
   - Type: `Injury`
   - Description: `Minor scrape on knee during PE`
   - Treatment Given: `Cleaned and bandaged`
   - Reported By: `PE Teacher`
4. Save
5. **Expected**: Incident recorded

---

## Story 4.4: Behavioral Incident Tracking

### Test Case 4.4.1: Create Positive Incident
1. Navigate to student detail page
2. Go to **Behavior** section/tab
3. Click **"Add Incident"**
4. Fill in:
   - Incident Type: `Positive`
   - Category: `Academic Achievement`
   - Title: `Math Competition Winner`
   - Description: `Won first place in inter-school math olympiad`
   - Incident Date: Today
   - Points: `+10`
5. Save
6. **Expected**: Positive incident recorded, points added

### Test Case 4.4.2: Create Negative Incident
1. Click **"Add Incident"**
2. Fill in:
   - Incident Type: `Negative`
   - Category: `Discipline`
   - Severity: `Minor`
   - Title: `Late to class`
   - Description: `Arrived 10 minutes late without valid reason`
   - Incident Date: Today
   - Points: `-2`
3. Save
4. **Expected**: Negative incident recorded, points deducted

### Test Case 4.4.3: Add Follow-up to Incident
1. Click on an existing incident
2. Click **"Add Follow-up"**
3. Fill in:
   - Follow-up Date: Today
   - Notes: `Spoke with student about punctuality`
   - Action Taken: `Verbal warning given`
   - Follow-up By: `Class Teacher`
4. Save
5. **Expected**: Follow-up added to incident

### Test Case 4.4.4: Resolve Incident
1. Click on an open incident
2. Click **"Resolve"** or change status to `Resolved`
3. Add resolution notes
4. Save
5. **Expected**: Incident marked as resolved

### Test Case 4.4.5: View Behavior Summary
1. Check the behavior summary section
2. **Expected**: Shows total points, positive count, negative count

---

## Story 4.5: Document Management

### Test Case 4.5.1: Upload Document
1. Navigate to student detail page
2. Go to **Documents** section/tab
3. Click **"Upload Document"**
4. Fill in:
   - Document Type: `Birth Certificate`
   - File: Select a PDF or image file
   - Notes: `Original copy`
5. Upload
6. **Expected**: Document uploaded, appears in list

### Test Case 4.5.2: View Document
1. Click on an uploaded document
2. **Expected**: Document preview or download starts

### Test Case 4.5.3: Verify Document
1. Find an unverified document
2. Click **"Verify"** button
3. **Expected**: Document marked as verified with timestamp

### Test Case 4.5.4: Delete Document
1. Click delete icon on a document
2. Confirm deletion
3. **Expected**: Document removed

### Test Case 4.5.5: Filter Documents by Type
1. Use document type filter/dropdown
2. Select a specific type
3. **Expected**: Only documents of that type shown

---

## Story 4.6: Enrollment History

### Test Case 4.6.1: View Enrollment History
1. Navigate to student detail page
2. Go to **Enrollment** section/tab
3. **Expected**: List of all enrollments with academic year, class, section, status

### Test Case 4.6.2: Create New Enrollment
1. Click **"Add Enrollment"**
2. Fill in:
   - Academic Year: Select from dropdown
   - Class: (Enter class ID or select)
   - Section: (Optional)
   - Roll Number: `15`
   - Enrollment Date: Today
3. Save
4. **Expected**: New enrollment created with `Active` status

### Test Case 4.6.3: View Enrollment Timeline
1. Check the enrollment timeline visualization
2. **Expected**: Shows chronological view of all enrollments

### Test Case 4.6.4: Process Transfer
1. Click on an active enrollment
2. Click **"Transfer"** action
3. Fill in:
   - Transfer Date: Today
   - Transfer Reason: `Family relocation`
4. Confirm
5. **Expected**: Enrollment status changed to `Transferred`

### Test Case 4.6.5: Process Dropout
1. Click on an active enrollment
2. Click **"Dropout"** action
3. Fill in:
   - Dropout Date: Today
   - Dropout Reason: `Financial constraints`
4. Confirm
5. **Expected**: Enrollment status changed to `Dropout`

---

## Story 4.7: Promotion/Retention Processing

### Test Case 4.7.1: Access Promotion Wizard
1. Navigate to **Students > Promotion** (or `/students/promotion`)
2. **Expected**: Multi-step wizard displayed

### Test Case 4.7.2: Step 1 - Select Source
1. Select **From Academic Year**: `2024-2025`
2. Select **To Academic Year**: `2025-2026`
3. Select **From Class**: Choose a class
4. Optionally select **From Section**
5. Select **To Class**: Next class level
6. Click **"Create Batch"**
7. **Expected**: Batch created, students loaded for review

### Test Case 4.7.3: Step 2 - Review Students
1. View list of students with pending decisions
2. Click **"Auto Decide"** button
3. **Expected**: Students automatically marked based on promotion rules

### Test Case 4.7.4: Manual Decision Override
1. Select a student
2. Click **Promote**, **Retain**, or **Transfer** action button
3. **Expected**: Decision updated for that student

### Test Case 4.7.5: Bulk Decision
1. Use checkboxes to select multiple students
2. Click **"Set Selected as Promote"** (or Retain)
3. **Expected**: All selected students' decisions updated

### Test Case 4.7.6: Step 3 - Assign Sections
1. Click **"Next"** to proceed to section assignment
2. For promoted students, assign target sections
3. **Expected**: Section assignments saved

### Test Case 4.7.7: Step 4 - Confirm & Process
1. Review the summary showing:
   - Total students
   - Promote count
   - Retain count
   - Transfer count
2. Toggle **"Generate Roll Numbers"** option
3. Click **"Process Promotions"**
4. **Expected**:
   - New enrollments created for next academic year
   - Old enrollments marked as completed
   - Success message displayed

### Test Case 4.7.8: Step 5 - Complete
1. View completion summary
2. Click **"Download Report"**
3. **Expected**: Promotion report downloaded

### Test Case 4.7.9: Verify New Enrollments
1. Go to a promoted student's detail page
2. Check Enrollment History
3. **Expected**: New enrollment for new academic year visible

---

## API Testing (Optional - via curl or Postman)

### Get Auth Token
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@myschoolsystem.com","password":"Admin@123","tenant_id":"61ef9fd2-2e9e-4b70-9f16-3b6ea73d4fa4"}' | jq -r '.data.access_token')
```

### Students API
```bash
# List students
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/students

# Get student by ID
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/students/{id}

# Create student
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"firstName":"Test","lastName":"Student","dateOfBirth":"2015-01-01","gender":"male","admissionDate":"2024-01-01"}' \
  http://localhost:8080/api/v1/students
```

### Guardians API
```bash
# List guardians for student
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/students/{id}/guardians

# Add guardian
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"relationship":"father","firstName":"Parent","lastName":"Name","phone":"+919876543210","isPrimary":true}' \
  http://localhost:8080/api/v1/students/{id}/guardians
```

### Health API
```bash
# Get health profile
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/students/{id}/health

# Update health profile
curl -X PUT -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"bloodGroup":"O+","heightCm":140,"weightKg":35}' \
  http://localhost:8080/api/v1/students/{id}/health
```

### Behavioral Incidents API
```bash
# List incidents
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/students/{id}/behavioral-incidents

# Create incident
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"incidentType":"positive","category":"achievement","title":"Good Work","description":"Excellent performance","incidentDate":"2024-01-15","points":5}' \
  http://localhost:8080/api/v1/students/{id}/behavioral-incidents
```

### Documents API
```bash
# List documents
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/students/{id}/documents

# Upload document (multipart form)
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/document.pdf" \
  -F "documentType=birth_certificate" \
  http://localhost:8080/api/v1/students/{id}/documents
```

### Enrollments API
```bash
# List enrollments
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/students/{id}/enrollments

# Create enrollment
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"academicYearId":"uuid","classId":"uuid","rollNumber":"15"}' \
  http://localhost:8080/api/v1/students/{id}/enrollments
```

### Promotion API
```bash
# List promotion batches
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/promotion-batches

# Create batch
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"fromAcademicYearId":"uuid","toAcademicYearId":"uuid","fromClassId":"uuid","toClassId":"uuid"}' \
  http://localhost:8080/api/v1/promotion-batches

# Auto-decide
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/promotion-batches/{id}/auto-decide

# Process batch
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"generateRollNumbers":true}' \
  http://localhost:8080/api/v1/promotion-batches/{id}/process
```

---

## Test Data Setup Script

Run this to create test data for Epic 4 testing:

```bash
# Get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@myschoolsystem.com","password":"Admin@123","tenant_id":"61ef9fd2-2e9e-4b70-9f16-3b6ea73d4fa4"}' | jq -r '.data.access_token')

# Create Academic Year (if not exists)
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"2024-2025","startDate":"2024-04-01","endDate":"2025-03-31","isCurrent":true}' \
  http://localhost:8080/api/v1/academic-years

# Create test students
for i in {1..5}; do
  curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d "{\"firstName\":\"Student\",\"lastName\":\"$i\",\"dateOfBirth\":\"201$i-01-15\",\"gender\":\"male\",\"admissionDate\":\"2024-01-01\"}" \
    http://localhost:8080/api/v1/students
done

echo "Test data created!"
```

---

## Known Limitations

1. **Story 4.8 (Student Search & Filters)**: Not yet implemented
2. **Classes/Sections**: Epic 6 not implemented - use placeholder UUIDs
3. **Attendance/Marks Integration**: Epic 7/8 not implemented - auto-decide uses placeholder logic
4. **Notifications**: Deferred to Epic 12

---

## Checklist Summary

| Story | Feature | Status |
|-------|---------|--------|
| 4.1 | Student Profile CRUD | Ready to Test |
| 4.2 | Guardian Management | Ready to Test |
| 4.3 | Health Records | Ready to Test |
| 4.4 | Behavioral Incidents | Ready to Test |
| 4.5 | Document Management | Ready to Test |
| 4.6 | Enrollment History | Ready to Test |
| 4.7 | Promotion Wizard | Ready to Test |
| 4.8 | Search & Filters | Not Implemented |
