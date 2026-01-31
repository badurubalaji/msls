// Package main is the entry point for the MSLS backend API server.
//
// @title MSLS API
// @version 1.0
// @description Multi-School Learning System API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"

	// _ "msls-backend/docs" // TODO: Generate docs with `swag init`
	academicyearhandler "msls-backend/internal/handlers/academicyear"
	admissionhandler "msls-backend/internal/handlers/admission"
	adminhandler "msls-backend/internal/handlers/admin"
	authhandler "msls-backend/internal/handlers/auth"
	branchhandler "msls-backend/internal/handlers/branch"
	profilehandler "msls-backend/internal/handlers/profile"
	rbachandler "msls-backend/internal/handlers/rbac"
	"msls-backend/internal/middleware"
	"msls-backend/internal/modules/assignment"
	"msls-backend/internal/modules/attendance"
	"msls-backend/internal/modules/behavioral"
	"msls-backend/internal/modules/bulk"
	"msls-backend/internal/modules/department"
	"msls-backend/internal/modules/designation"
	"msls-backend/internal/modules/document"
	"msls-backend/internal/modules/enrollment"
	"msls-backend/internal/modules/exam"
	"msls-backend/internal/modules/examination"
	"msls-backend/internal/modules/hallticket"
	"msls-backend/internal/modules/academic"
	"msls-backend/internal/modules/timetable"
	"msls-backend/internal/modules/guardian"
	"msls-backend/internal/modules/health"
	"msls-backend/internal/modules/payroll"
	"msls-backend/internal/modules/promotion"
	"msls-backend/internal/modules/salary"
	"msls-backend/internal/modules/staff"
	"msls-backend/internal/modules/staffdocument"
	"msls-backend/internal/modules/student"
	"msls-backend/internal/modules/studentattendance"
	"msls-backend/internal/pkg/config"
	"msls-backend/internal/pkg/database"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/logger"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/pkg/sms"
	"msls-backend/internal/services/academicyear"
	"msls-backend/internal/services/admission"
	"msls-backend/internal/services/auth"
	"msls-backend/internal/services/branch"
	"msls-backend/internal/services/featureflag"
	"msls-backend/internal/services/profile"
	"msls-backend/internal/services/rbac"
	"msls-backend/internal/pkg/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("starting server",
		zap.String("app", cfg.App.Name),
		zap.String("environment", cfg.App.Environment),
		zap.Int("port", cfg.Server.Port),
	)

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
	log.Info("database connected")

	// Set Gin mode based on environment
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router with middleware
	router := setupRouter(cfg, log, db)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("server listening", zap.String("address", addr))
		serverErrors <- srv.ListenAndServe()
	}()

	// Wait for interrupt signal or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-shutdown:
		log.Info("shutdown signal received", zap.String("signal", sig.String()))

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("graceful shutdown failed", zap.Error(err))
			if err := srv.Close(); err != nil {
				return fmt.Errorf("forced shutdown error: %w", err)
			}
		}

		log.Info("server stopped gracefully")
	}

	return nil
}

func setupRouter(cfg *config.Config, log *logger.Logger, db *gorm.DB) *gin.Engine {
	router := gin.New()

	// === Global Middleware (applied to all routes) ===
	// Order matters: these are executed in the order they are added

	// 1. CORS - Must be first to handle preflight requests
	corsConfig := middleware.DefaultCORSConfig()
	if cfg.App.IsProduction() {
		// In production, configure allowed origins from environment
		corsConfig = middleware.ProductionCORSConfig([]string{
			"https://msls.example.com",
			// Add production origins here
		})
	}
	router.Use(middleware.CORS(corsConfig))

	// 2. Request ID - Generate/propagate request ID for tracing
	router.Use(middleware.RequestIDDefault())

	// 3. Recovery - Catch panics and return 500 errors
	router.Use(middleware.RecoveryDefault(log))

	// 4. Logging - Log all requests (after request ID so it's available)
	router.Use(middleware.LoggingDefault(log))

	// 5. Error Handler - Convert errors to RFC 7807 responses
	router.Use(apperrors.Handler(log))

	// 6. Rate Limiting - Global rate limit (100 req/min by default)
	router.Use(middleware.RateLimitDefault())

	// === Static File Serving ===
	// Serve uploaded files (documents, avatars, etc.)
	router.Static("/uploads", "./uploads")

	// === Public Routes (no tenant required) ===
	// Health check endpoint (excluded from tenant middleware)
	router.GET("/health", healthHandler)
	router.GET("/ready", readyHandler)

	// Swagger documentation endpoint
	// TODO: Enable after running `swag init` to generate docs
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize services
	jwtService := auth.NewJWTService(auth.JWTConfig{
		Secret:     cfg.JWT.Secret,
		Issuer:     cfg.JWT.Issuer,
		AccessTTL:  cfg.JWT.AccessExpiresIn,
		RefreshTTL: cfg.JWT.RefreshExpiresIn,
	})
	authService := auth.NewAuthService(db, jwtService)

	// Initialize RBAC services
	permissionService := rbac.NewPermissionService(db)
	roleService := rbac.NewRoleService(db, permissionService)
	userRoleService := rbac.NewUserRoleService(db, roleService)

	// Initialize SMS provider (mock for development)
	smsProvider, err := sms.NewMockProvider("")
	if err != nil {
		log.Warn("failed to initialize SMS provider, OTP via SMS will not work", zap.Error(err))
	}

	// Initialize OTP service
	otpService := auth.NewOTPService(db, jwtService, auth.OTPConfig{
		SMSProvider: smsProvider,
	})

	// Initialize TOTP service for 2FA
	totpService, err := auth.NewTOTPService(db, cfg.JWT.Secret)
	if err != nil {
		log.Warn("failed to initialize TOTP service, 2FA will not work", zap.Error(err))
	}
	authService.SetTOTPService(totpService)

	// Initialize profile service
	profileService := profile.NewProfileService(db, profile.Config{
		UploadDir: "./uploads/avatars",
	})

	// Initialize feature flag service
	featureFlagService := featureflag.NewService(db, featureflag.DefaultConfig())

	// Initialize branch service
	branchService := branch.NewService(db)

	// Initialize academic year service
	academicYearService := academicyear.NewService(db)

	// Initialize admission services
	admissionSessionService := admission.NewSessionService(db)
	admissionReportService := admission.NewReportService(db)
	admissionExportService := admission.NewExportService(db, admissionReportService)
	enquiryService := admission.NewEnquiryService(db)
	applicationService := admission.NewApplicationService(db)
	testService := admission.NewTestService(db)
	reviewService := admission.NewReviewService(db)
	meritService := admission.NewMeritService(db)
	decisionService := admission.NewDecisionService(db)

	// Initialize student service
	studentService := student.NewService(db, branchService)

	// Initialize guardian service
	guardianRepo := guardian.NewRepository(db)
	guardianService := guardian.NewService(guardianRepo)

	// Initialize health service
	healthRepo := health.NewRepository(db)
	healthService := health.NewService(healthRepo)

	// Initialize behavioral service
	behavioralRepo := behavioral.NewRepository(db)
	behavioralService := behavioral.NewService(behavioralRepo)

	// Initialize document service
	documentRepo := document.NewRepository(db)
	documentService := document.NewService(documentRepo, "./uploads")

	// Initialize enrollment service (needs adapters for student and academic year)
	enrollmentRepo := enrollment.NewRepository(db)
	enrollmentService := enrollment.NewService(db, nil, nil) // Adapters will be added when needed

	// Initialize promotion service
	promotionRepo := promotion.NewRepository(db)
	promotionService := promotion.NewService(promotionRepo, enrollmentRepo)

	// Initialize bulk operation services
	bulkExportService := bulk.NewExportService(db, "./uploads")
	bulkService := bulk.NewService(db, bulkExportService)
	bulkImportService := bulk.NewImportService(db)

	// Initialize department service
	departmentRepo := department.NewRepository(db)
	departmentService := department.NewService(departmentRepo)

	// Initialize designation service
	designationRepo := designation.NewRepository(db)
	designationService := designation.NewService(designationRepo)

	// Initialize staff service
	staffService := staff.NewService(db, branchService)

	// Initialize salary service
	salaryRepo := salary.NewRepository(db)
	salaryService := salary.NewService(salaryRepo)

	// Initialize payroll service
	payrollRepo := payroll.NewRepository(db)
	payrollService := payroll.NewService(payrollRepo)

	// Initialize assignment service
	assignmentService := assignment.NewService(db)

	// Initialize academic service
	academicRepo := academic.NewRepository(db)
	academicService := academic.NewService(academicRepo)

	// Initialize timetable service
	timetableRepo := timetable.NewRepository(db)
	timetableService := timetable.NewService(timetableRepo)

	// Initialize exam service
	examRepo := exam.NewRepository(db)
	examService := exam.NewService(examRepo)

	// Initialize examination service
	examinationRepo := examination.NewRepository(db)
	examinationService := examination.NewService(examinationRepo)

	// Initialize hall ticket service
	hallTicketService := hallticket.NewService(db, cfg.JWT.Secret)

	// Initialize file storage for staff documents
	fileStorage, err := storage.NewLocalStorage("./uploads", "/uploads")
	if err != nil {
		log.Warn("failed to initialize file storage", zap.Error(err))
	}

	// Initialize staff document service
	staffDocumentRepo := staffdocument.NewRepository(db)
	staffDocumentService := staffdocument.NewService(staffDocumentRepo)

	// Initialize handlers
	authHandler := authhandler.NewHandler(authService)
	otpHandler := authhandler.NewOTPHandler(otpService)
	twoFactorHandler := authhandler.NewTwoFactorHandler(authService, totpService)
	profileHandler := profilehandler.NewHandler(profileService)
	roleHandler := rbachandler.NewRoleHandler(roleService)
	permissionHandler := rbachandler.NewPermissionHandler(permissionService)
	userRoleHandler := rbachandler.NewUserRoleHandler(userRoleService)
	featureFlagHandler := adminhandler.NewFeatureFlagHandler(featureFlagService)
	branchHandler := branchhandler.NewHandler(branchService)
	academicYearHandler := academicyearhandler.NewHandler(academicYearService)
	admissionSessionHandler := admissionhandler.NewSessionHandler(admissionSessionService)
	admissionReportHandler := admissionhandler.NewReportHandler(admissionReportService)
	admissionExportHandler := admissionhandler.NewExportHandler(admissionExportService)
	enquiryHandler := admissionhandler.NewEnquiryHandler(enquiryService)
	applicationHandler := admissionhandler.NewApplicationHandler(applicationService)
	testHandler := admissionhandler.NewTestHandler(testService)
	reviewHandler := admissionhandler.NewReviewHandler(reviewService)
	meritHandler := admissionhandler.NewMeritHandler(meritService)
	decisionHandler := admissionhandler.NewDecisionHandler(decisionService)
	studentHandler := student.NewHandler(studentService)
	guardianHandler := guardian.NewHandler(guardianService)
	healthHandler := health.NewHandler(healthService)
	behavioralHandler := behavioral.NewHandler(behavioralService)
	documentHandler := document.NewHandler(documentService)
	enrollmentHandler := enrollment.NewHandler(enrollmentService)
	promotionHandler := promotion.NewHandler(promotionService)
	bulkHandler := bulk.NewHandler(bulkService, bulkImportService)
	departmentHandler := department.NewHandler(departmentService)
	designationHandler := designation.NewHandler(designationService)
	staffHandler := staff.NewHandler(staffService)
	salaryHandler := salary.NewHandler(salaryService)
	payrollHandler := payroll.NewHandler(payrollService)
	assignmentHandler := assignment.NewHandler(assignmentService)
	academicHandler := academic.NewHandler(academicService)
	timetableHandler := timetable.NewHandler(timetableService)
	examHandler := exam.NewHandler(examService)
	examinationHandler := examination.NewHandler(examinationService)
	hallTicketHandler := hallticket.NewHandler(hallTicketService)
	staffDocumentHandler := staffdocument.NewHandler(staffDocumentService, fileStorage)

	// Initialize attendance service (wrapping staff service for lookup)
	attendanceService := attendance.NewService(db, staffService)
	attendanceHandler := attendance.NewHandler(attendanceService)

	// Initialize student attendance service
	studentAttendanceService := studentattendance.NewService(db)
	studentAttendanceHandler := studentattendance.NewHandler(studentAttendanceService)

	// === API v1 Routes ===
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication or tenant required)
		public := v1.Group("")
		{
			public.GET("/ping", pingHandler)
		}

		// Public routes that require tenant ID but no authentication
		publicTenant := v1.Group("/public")
		publicTenant.Use(middleware.TenantRequired())
		{
			// Application status check (for parents to check their application status)
			publicTenant.POST("/applications/status", applicationHandler.CheckStatus)
		}

		// Auth routes (public - no authentication required)
		authRoutes := v1.Group("/auth")
		{
			// Public auth endpoints
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
			authRoutes.POST("/verify-email", authHandler.VerifyEmail)
			authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
			authRoutes.POST("/reset-password", authHandler.ResetPassword)

			// OTP endpoints (public - for passwordless login)
			otpRoutes := authRoutes.Group("/otp")
			{
				otpRoutes.POST("/request", otpHandler.RequestOTP)
				otpRoutes.POST("/verify", otpHandler.VerifyOTP)
				otpRoutes.POST("/resend", otpHandler.ResendOTP)
			}

			// 2FA validation endpoint (public - uses partial token)
			authRoutes.POST("/2fa/validate", twoFactorHandler.Validate2FA)

			// Protected auth endpoints (require authentication)
			authProtected := authRoutes.Group("")
			authProtected.Use(middleware.AuthRequired(jwtService))
			{
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.GET("/me", authHandler.Me)

				// 2FA management endpoints (require authentication)
				twoFactorRoutes := authProtected.Group("/2fa")
				{
					twoFactorRoutes.POST("/setup", twoFactorHandler.Setup2FA)
					twoFactorRoutes.POST("/verify", twoFactorHandler.Verify2FA)
					twoFactorRoutes.POST("/disable", twoFactorHandler.Disable2FA)
					twoFactorRoutes.GET("/status", twoFactorHandler.GetStatus)
					twoFactorRoutes.GET("/backup-codes", twoFactorHandler.GetBackupCodes)
					twoFactorRoutes.POST("/regenerate-backup", twoFactorHandler.RegenerateBackupCodes)
				}
			}

			// Admin-only auth endpoints (require authentication + permission)
			authAdmin := authRoutes.Group("")
			authAdmin.Use(middleware.TenantRequired())
			authAdmin.Use(middleware.AuthRequired(jwtService))
			authAdmin.Use(middleware.PermissionRequired("users:write"))
			{
				authAdmin.POST("/register", authHandler.Register)
			}
		}

		// Profile routes (require authentication, tenant from token)
		profileRoutes := v1.Group("/profile")
		profileRoutes.Use(middleware.AuthRequired(jwtService))
		{
			profileRoutes.GET("", profileHandler.GetProfile)
			profileRoutes.PUT("", profileHandler.UpdateProfile)
			profileRoutes.DELETE("", profileHandler.RequestAccountDeletion)
			profileRoutes.POST("/avatar", profileHandler.UploadAvatar)
			profileRoutes.PUT("/password", profileHandler.ChangePassword)
			profileRoutes.GET("/preferences", profileHandler.GetPreferences)
			profileRoutes.PUT("/preferences", profileHandler.UpdatePreferences)
			profileRoutes.GET("/preferences/extended", profileHandler.GetUserPreferences)
			profileRoutes.POST("/preferences/extended", profileHandler.SetUserPreference)
			profileRoutes.DELETE("/preferences/extended", profileHandler.DeleteUserPreference)
		}

		// Protected routes (require tenant ID and authentication)
		protected := v1.Group("")
		protected.Use(middleware.TenantRequired())
		protected.Use(middleware.AuthRequired(jwtService))
		{
			// Role management routes
			roles := protected.Group("/roles")
			{
				// Read operations - require roles:read permission
				rolesRead := roles.Group("")
				rolesRead.Use(middleware.PermissionRequired("roles:read"))
				{
					rolesRead.GET("", roleHandler.List)
					rolesRead.GET("/:id", roleHandler.GetByID)
				}

				// Write operations - require roles:write permission
				rolesWrite := roles.Group("")
				rolesWrite.Use(middleware.PermissionRequired("roles:write"))
				{
					rolesWrite.POST("", roleHandler.Create)
					rolesWrite.PUT("/:id", roleHandler.Update)
					rolesWrite.POST("/:id/permissions", roleHandler.AssignPermissions)
					rolesWrite.DELETE("/:id/permissions", roleHandler.RemovePermissions)
				}

				// Delete operations - require roles:delete permission
				rolesDelete := roles.Group("")
				rolesDelete.Use(middleware.PermissionRequired("roles:delete"))
				{
					rolesDelete.DELETE("/:id", roleHandler.Delete)
				}
			}

			// Permission routes (read-only for most users)
			permissions := protected.Group("/permissions")
			permissions.Use(middleware.PermissionRequired("roles:read"))
			{
				permissions.GET("", permissionHandler.List)
				permissions.GET("/modules", permissionHandler.GetModules)
				permissions.GET("/modules/:module", permissionHandler.GetByModule)
			}

			// User role management routes
			users := protected.Group("/users")
			{
				// Current user can always view their own roles
				users.GET("/me/roles", userRoleHandler.GetMyRoles)

				// Admin operations on user roles - require users:write permission
				userRoles := users.Group("/:id/roles")
				userRoles.Use(middleware.PermissionRequired("users:write"))
				{
					userRoles.GET("", userRoleHandler.GetUserRoles)
					userRoles.POST("", userRoleHandler.AssignRoles)
					userRoles.DELETE("", userRoleHandler.RemoveRoles)
				}
			}

			// Branch management routes
			branches := protected.Group("/branches")
			{
				// Read operations - require branches:read permission
				branchesRead := branches.Group("")
				branchesRead.Use(middleware.PermissionRequired("branches:read"))
				{
					branchesRead.GET("", branchHandler.List)
					branchesRead.GET("/:id", branchHandler.GetByID)
				}

				// Write operations - require branches:update permission
				branchesWrite := branches.Group("")
				branchesWrite.Use(middleware.PermissionRequired("branches:update"))
				{
					branchesWrite.PUT("/:id", branchHandler.Update)
					branchesWrite.PATCH("/:id/primary", branchHandler.SetPrimary)
					branchesWrite.PATCH("/:id/status", branchHandler.SetStatus)
				}

				// Create operations - require branches:create permission
				branchesCreate := branches.Group("")
				branchesCreate.Use(middleware.PermissionRequired("branches:create"))
				{
					branchesCreate.POST("", branchHandler.Create)
				}

				// Delete operations - require branches:delete permission
				branchesDelete := branches.Group("")
				branchesDelete.Use(middleware.PermissionRequired("branches:delete"))
				{
					branchesDelete.DELETE("/:id", branchHandler.Delete)
				}
			}

			// Academic year management routes
			academicYears := protected.Group("/academic-years")
			{
				// Read operations - require academic-years:read permission
				academicYearsRead := academicYears.Group("")
				academicYearsRead.Use(middleware.PermissionRequired("academic-years:read"))
				{
					academicYearsRead.GET("", academicYearHandler.List)
					academicYearsRead.GET("/current", academicYearHandler.GetCurrent)
					academicYearsRead.GET("/:id", academicYearHandler.GetByID)
					academicYearsRead.GET("/:id/terms", academicYearHandler.ListTerms)
					academicYearsRead.GET("/:id/holidays", academicYearHandler.ListHolidays)
				}

				// Create operations - require academic-years:create permission
				academicYearsCreate := academicYears.Group("")
				academicYearsCreate.Use(middleware.PermissionRequired("academic-years:create"))
				{
					academicYearsCreate.POST("", academicYearHandler.Create)
				}

				// Update operations - require academic-years:update permission
				academicYearsUpdate := academicYears.Group("")
				academicYearsUpdate.Use(middleware.PermissionRequired("academic-years:update"))
				{
					academicYearsUpdate.PUT("/:id", academicYearHandler.Update)
					academicYearsUpdate.PATCH("/:id/current", academicYearHandler.SetCurrent)
					academicYearsUpdate.POST("/:id/terms", academicYearHandler.CreateTerm)
					academicYearsUpdate.PUT("/:id/terms/:termId", academicYearHandler.UpdateTerm)
					academicYearsUpdate.DELETE("/:id/terms/:termId", academicYearHandler.DeleteTerm)
					academicYearsUpdate.POST("/:id/holidays", academicYearHandler.CreateHoliday)
					academicYearsUpdate.PUT("/:id/holidays/:holidayId", academicYearHandler.UpdateHoliday)
					academicYearsUpdate.DELETE("/:id/holidays/:holidayId", academicYearHandler.DeleteHoliday)
				}

				// Delete operations - require academic-years:delete permission
				academicYearsDelete := academicYears.Group("")
				academicYearsDelete.Use(middleware.PermissionRequired("academic-years:delete"))
				{
					academicYearsDelete.DELETE("/:id", academicYearHandler.Delete)
				}
			}

			// Student management routes
			students := protected.Group("/students")
			{
				// Read operations - require students:read permission
				studentsRead := students.Group("")
				studentsRead.Use(middleware.PermissionRequired("students:read"))
				{
					studentsRead.GET("", studentHandler.List)
					studentsRead.GET("/next-admission-number", studentHandler.GetNextAdmissionNumber)
					studentsRead.GET("/:id", studentHandler.GetByID)
				}

				// Create operations - require students:write permission
				studentsCreate := students.Group("")
				studentsCreate.Use(middleware.PermissionRequired("students:write"))
				{
					studentsCreate.POST("", studentHandler.Create)
				}

				// Update operations - require students:write permission
				studentsUpdate := students.Group("")
				studentsUpdate.Use(middleware.PermissionRequired("students:write"))
				{
					studentsUpdate.PUT("/:id", studentHandler.Update)
					studentsUpdate.POST("/:id/photo", studentHandler.UpdatePhoto)
				}

				// Delete operations - require students:delete permission
				studentsDelete := students.Group("")
				studentsDelete.Use(middleware.PermissionRequired("students:delete"))
				{
					studentsDelete.DELETE("/:id", studentHandler.Delete)
				}

				// Bulk operations - require students:update permission
				studentsBulk := students.Group("/bulk")
				studentsBulk.Use(middleware.PermissionRequired("students:update"))
				{
					studentsBulk.POST("/status", bulkHandler.BulkStatusUpdate)
				}

				// Export operations - require students:export permission
				studentsExport := students.Group("")
				studentsExport.Use(middleware.PermissionRequired("students:export"))
				{
					studentsExport.POST("/export", bulkHandler.Export)
				}

				// Import operations - require students:create permission
				studentsImport := students.Group("/import")
				studentsImport.Use(middleware.PermissionRequired("students:create"))
				{
					studentsImport.GET("/template", bulkHandler.DownloadTemplate)
					studentsImport.POST("", bulkHandler.ImportStudents)
				}

				// Guardian management routes (nested under students)
				guardians := students.Group("/:id/guardians")
				{
					// Read operations - require guardians:read permission
					guardiansRead := guardians.Group("")
					guardiansRead.Use(middleware.PermissionRequired("guardians:read"))
					{
						guardiansRead.GET("", guardianHandler.ListGuardians)
						guardiansRead.GET("/:guardianId", guardianHandler.GetGuardian)
					}

					// Write operations - require guardians:write permission
					guardiansWrite := guardians.Group("")
					guardiansWrite.Use(middleware.PermissionRequired("guardians:write"))
					{
						guardiansWrite.POST("", guardianHandler.CreateGuardian)
						guardiansWrite.PUT("/:guardianId", guardianHandler.UpdateGuardian)
						guardiansWrite.DELETE("/:guardianId", guardianHandler.DeleteGuardian)
						guardiansWrite.POST("/:guardianId/set-primary", guardianHandler.SetPrimaryGuardian)
					}
				}

				// Emergency contact management routes (nested under students)
				emergencyContacts := students.Group("/:id/emergency-contacts")
				{
					// Read operations - require emergency_contacts:read permission
					contactsRead := emergencyContacts.Group("")
					contactsRead.Use(middleware.PermissionRequired("emergency_contacts:read"))
					{
						contactsRead.GET("", guardianHandler.ListEmergencyContacts)
						contactsRead.GET("/:contactId", guardianHandler.GetEmergencyContact)
					}

					// Write operations - require emergency_contacts:write permission
					contactsWrite := emergencyContacts.Group("")
					contactsWrite.Use(middleware.PermissionRequired("emergency_contacts:write"))
					{
						contactsWrite.POST("", guardianHandler.CreateEmergencyContact)
						contactsWrite.PUT("/:contactId", guardianHandler.UpdateEmergencyContact)
						contactsWrite.DELETE("/:contactId", guardianHandler.DeleteEmergencyContact)
					}
				}

				// Health records management routes (nested under students)
				healthRoutes := students.Group("/:id/health")
				{
					// Read operations - require health:read permission
					healthRead := healthRoutes.Group("")
					healthRead.Use(middleware.PermissionRequired("health:read"))
					{
						healthRead.GET("", healthHandler.GetHealthSummary)
						healthRead.GET("/profile", healthHandler.GetHealthProfile)
						healthRead.GET("/allergies", healthHandler.ListAllergies)
						healthRead.GET("/allergies/:allergyId", healthHandler.GetAllergy)
						healthRead.GET("/conditions", healthHandler.ListConditions)
						healthRead.GET("/conditions/:conditionId", healthHandler.GetCondition)
						healthRead.GET("/medications", healthHandler.ListMedications)
						healthRead.GET("/medications/:medicationId", healthHandler.GetMedication)
						healthRead.GET("/vaccinations", healthHandler.ListVaccinations)
						healthRead.GET("/vaccinations/:vaccinationId", healthHandler.GetVaccination)
						healthRead.GET("/incidents", healthHandler.ListIncidents)
						healthRead.GET("/incidents/:incidentId", healthHandler.GetIncident)
					}

					// Write operations - require health:write permission
					healthWrite := healthRoutes.Group("")
					healthWrite.Use(middleware.PermissionRequired("health:write"))
					{
						healthWrite.PUT("/profile", healthHandler.CreateOrUpdateHealthProfile)
						healthWrite.POST("/allergies", healthHandler.CreateAllergy)
						healthWrite.PUT("/allergies/:allergyId", healthHandler.UpdateAllergy)
						healthWrite.DELETE("/allergies/:allergyId", healthHandler.DeleteAllergy)
						healthWrite.POST("/conditions", healthHandler.CreateCondition)
						healthWrite.PUT("/conditions/:conditionId", healthHandler.UpdateCondition)
						healthWrite.DELETE("/conditions/:conditionId", healthHandler.DeleteCondition)
						healthWrite.POST("/medications", healthHandler.CreateMedication)
						healthWrite.PUT("/medications/:medicationId", healthHandler.UpdateMedication)
						healthWrite.DELETE("/medications/:medicationId", healthHandler.DeleteMedication)
						healthWrite.POST("/vaccinations", healthHandler.CreateVaccination)
						healthWrite.PUT("/vaccinations/:vaccinationId", healthHandler.UpdateVaccination)
						healthWrite.DELETE("/vaccinations/:vaccinationId", healthHandler.DeleteVaccination)
						healthWrite.POST("/incidents", healthHandler.CreateIncident)
						healthWrite.PUT("/incidents/:incidentId", healthHandler.UpdateIncident)
						healthWrite.DELETE("/incidents/:incidentId", healthHandler.DeleteIncident)
					}
				}

				// Behavioral incidents management routes (nested under students)
				behavioralRoutes := students.Group("/:id/behavioral-incidents")
				{
					// Read operations - require behavior:read permission
					behavioralRead := behavioralRoutes.Group("")
					behavioralRead.Use(middleware.PermissionRequired("behavior:read"))
					{
						behavioralRead.GET("", behavioralHandler.ListIncidents)
						behavioralRead.GET("/:incidentId", behavioralHandler.GetIncident)
					}

					// Write operations - require behavior:write permission
					behavioralWrite := behavioralRoutes.Group("")
					behavioralWrite.Use(middleware.PermissionRequired("behavior:write"))
					{
						behavioralWrite.POST("", behavioralHandler.CreateIncident)
						behavioralWrite.PUT("/:incidentId", behavioralHandler.UpdateIncident)
						behavioralWrite.DELETE("/:incidentId", behavioralHandler.DeleteIncident)
					}
				}

				// Behavioral summary route
				behavioralSummary := students.Group("/:id/behavioral-summary")
				behavioralSummary.Use(middleware.PermissionRequired("behavior:read"))
				{
					behavioralSummary.GET("", behavioralHandler.GetBehaviorSummary)
				}

				// Enrollment management routes (nested under students)
				enrollments := students.Group("/:id/enrollments")
				{
					// Read operations - require enrollments:read permission
					enrollmentsRead := enrollments.Group("")
					enrollmentsRead.Use(middleware.PermissionRequired("enrollments:read"))
					{
						enrollmentsRead.GET("", enrollmentHandler.ListEnrollments)
						enrollmentsRead.GET("/current", enrollmentHandler.GetCurrentEnrollment)
						enrollmentsRead.GET("/:eid/status-history", enrollmentHandler.GetStatusHistory)
					}

					// Create operations - require enrollments:create permission
					enrollmentsCreate := enrollments.Group("")
					enrollmentsCreate.Use(middleware.PermissionRequired("enrollments:create"))
					{
						enrollmentsCreate.POST("", enrollmentHandler.CreateEnrollment)
					}

					// Update operations - require enrollments:update permission
					enrollmentsUpdate := enrollments.Group("")
					enrollmentsUpdate.Use(middleware.PermissionRequired("enrollments:update"))
					{
						enrollmentsUpdate.PUT("/:eid", enrollmentHandler.UpdateEnrollment)
						enrollmentsUpdate.POST("/:eid/transfer", enrollmentHandler.ProcessTransfer)
						enrollmentsUpdate.POST("/:eid/dropout", enrollmentHandler.ProcessDropout)
					}
				}
			}

				// Document management routes (nested under students)
				documentRoutes := students.Group("/:id/documents")
				{
					// Read operations - require document:read permission
					documentRead := documentRoutes.Group("")
					documentRead.Use(middleware.PermissionRequired("document:read"))
					{
						documentRead.GET("", documentHandler.ListDocuments)
						documentRead.GET("/:docId", documentHandler.GetDocument)
					}

					// Write operations - require document:create permission
					documentWrite := documentRoutes.Group("")
					documentWrite.Use(middleware.PermissionRequired("document:create"))
					{
						documentWrite.POST("", documentHandler.UploadDocument)
					}

					// Update operations - require document:update permission
					documentUpdate := documentRoutes.Group("")
					documentUpdate.Use(middleware.PermissionRequired("document:update"))
					{
						documentUpdate.PUT("/:docId", documentHandler.UpdateDocument)
					}

					// Delete operations - require document:delete permission
					documentDelete := documentRoutes.Group("")
					documentDelete.Use(middleware.PermissionRequired("document:delete"))
					{
						documentDelete.DELETE("/:docId", documentHandler.DeleteDocument)
					}

					// Verify operations - require document:verify permission
					documentVerify := documentRoutes.Group("")
					documentVerify.Use(middleware.PermissionRequired("document:verify"))
					{
						documentVerify.POST("/:docId/verify", documentHandler.VerifyDocument)
						documentVerify.POST("/:docId/reject", documentHandler.RejectDocument)
					}
				}

				// Document checklist route
				documentChecklist := students.Group("/:id/document-checklist")
				documentChecklist.Use(middleware.PermissionRequired("document:read"))
				{
					documentChecklist.GET("", documentHandler.GetDocumentChecklist)
				}

			// Document type management routes
			documentTypes := protected.Group("/document-types")
			{
				// Read operations - require document:read permission
				documentTypesRead := documentTypes.Group("")
				documentTypesRead.Use(middleware.PermissionRequired("document:read"))
				{
					documentTypesRead.GET("", documentHandler.ListDocumentTypes)
				}

				// Write operations - require document_type:manage permission
				documentTypesWrite := documentTypes.Group("")
				documentTypesWrite.Use(middleware.PermissionRequired("document_type:manage"))
				{
					documentTypesWrite.POST("", documentHandler.CreateDocumentType)
					documentTypesWrite.PUT("/:id", documentHandler.UpdateDocumentType)
				}
			}

			// Bulk operations management routes
			bulkOperations := protected.Group("/bulk-operations")
			{
				// Read operations - require students:read permission
				bulkOpsRead := bulkOperations.Group("")
				bulkOpsRead.Use(middleware.PermissionRequired("students:read"))
				{
					bulkOpsRead.GET("", bulkHandler.ListOperations)
					bulkOpsRead.GET("/:id", bulkHandler.GetOperation)
					bulkOpsRead.GET("/:id/result", bulkHandler.DownloadResult)
				}
			}

			// Enrollment lookup routes (by class/section)
			enrollmentLookup := protected.Group("/enrollments")
			{
				enrollmentLookupRead := enrollmentLookup.Group("")
				enrollmentLookupRead.Use(middleware.PermissionRequired("enrollments:read"))
				{
					enrollmentLookupRead.GET("/by-class/:classId", enrollmentHandler.ListByClass)
					enrollmentLookupRead.GET("/by-section/:sectionId", enrollmentHandler.ListBySection)
				}
			}

			// Promotion rules management routes
			promotionRules := protected.Group("/promotion-rules")
			{
				// Read operations - require promotion_rules:read permission
				promotionRulesRead := promotionRules.Group("")
				promotionRulesRead.Use(middleware.PermissionRequired("promotion_rules:read"))
				{
					promotionRulesRead.GET("", promotionHandler.ListRules)
					promotionRulesRead.GET("/:id", promotionHandler.GetRule)
				}

				// Write operations - require promotion_rules:manage permission
				promotionRulesWrite := promotionRules.Group("")
				promotionRulesWrite.Use(middleware.PermissionRequired("promotion_rules:manage"))
				{
					promotionRulesWrite.POST("", promotionHandler.CreateOrUpdateRule)
					promotionRulesWrite.DELETE("/:id", promotionHandler.DeleteRule)
				}
			}

			// Promotion batch management routes
			promotionBatches := protected.Group("/promotion-batches")
			{
				// Read operations - require promotion:read permission
				promotionBatchesRead := promotionBatches.Group("")
				promotionBatchesRead.Use(middleware.PermissionRequired("promotion:read"))
				{
					promotionBatchesRead.GET("", promotionHandler.ListBatches)
					promotionBatchesRead.GET("/:id", promotionHandler.GetBatch)
					promotionBatchesRead.GET("/:id/records", promotionHandler.ListRecords)
					promotionBatchesRead.GET("/:id/report", promotionHandler.GetReport)
				}

				// Create operations - require promotion:create permission
				promotionBatchesCreate := promotionBatches.Group("")
				promotionBatchesCreate.Use(middleware.PermissionRequired("promotion:create"))
				{
					promotionBatchesCreate.POST("", promotionHandler.CreateBatch)
				}

				// Update operations - require promotion:update permission
				promotionBatchesUpdate := promotionBatches.Group("")
				promotionBatchesUpdate.Use(middleware.PermissionRequired("promotion:update"))
				{
					promotionBatchesUpdate.PUT("/:id/records/:recordId", promotionHandler.UpdateRecord)
					promotionBatchesUpdate.POST("/:id/records/bulk", promotionHandler.BulkUpdateRecords)
					promotionBatchesUpdate.POST("/:id/auto-decide", promotionHandler.AutoDecide)
				}

				// Process operations - require promotion:process permission
				promotionBatchesProcess := promotionBatches.Group("")
				promotionBatchesProcess.Use(middleware.PermissionRequired("promotion:process"))
				{
					promotionBatchesProcess.POST("/:id/process", promotionHandler.ProcessBatch)
				}

				// Cancel operations - require promotion:cancel permission
				promotionBatchesCancel := promotionBatches.Group("")
				promotionBatchesCancel.Use(middleware.PermissionRequired("promotion:cancel"))
				{
					promotionBatchesCancel.POST("/:id/cancel", promotionHandler.CancelBatch)
					promotionBatchesCancel.DELETE("/:id", promotionHandler.DeleteBatch)
				}
			}

			// Behavioral follow-ups routes (separate from students)
			behavioralIncidents := protected.Group("/behavioral-incidents")
			{
				// Follow-up routes
				followUps := behavioralIncidents.Group("/:incidentId/follow-ups")
				{
					followUpsWrite := followUps.Group("")
					followUpsWrite.Use(middleware.PermissionRequired("behavior:write"))
					{
						followUpsWrite.POST("", behavioralHandler.CreateFollowUp)
						followUpsWrite.PUT("/:followUpId", behavioralHandler.UpdateFollowUp)
						followUpsWrite.DELETE("/:followUpId", behavioralHandler.DeleteFollowUp)
					}
				}
			}

			// Pending follow-ups list
			followUpsRoute := protected.Group("/follow-ups")
			followUpsRoute.Use(middleware.PermissionRequired("behavior:read"))
			{
				followUpsRoute.GET("/pending", behavioralHandler.ListPendingFollowUps)
			}

			// Admission session management routes
			admissionSessions := protected.Group("/admission-sessions")
			{
				// Read operations - require admissions:read permission
				admissionsRead := admissionSessions.Group("")
				admissionsRead.Use(middleware.PermissionRequired("admissions:read"))
				{
					admissionsRead.GET("", admissionSessionHandler.List)
					admissionsRead.GET("/:id", admissionSessionHandler.GetByID)
					admissionsRead.GET("/:id/seats", admissionSessionHandler.ListSeats)
					admissionsRead.GET("/:id/stats", admissionSessionHandler.GetStats)
					admissionsRead.GET("/:id/merit-list", meritHandler.GetMeritList)
					admissionsRead.GET("/:id/merit-lists", meritHandler.ListMeritLists)
				}

				// Create operations - require admissions:create permission
				admissionsCreate := admissionSessions.Group("")
				admissionsCreate.Use(middleware.PermissionRequired("admissions:create"))
				{
					admissionsCreate.POST("", admissionSessionHandler.Create)
				}

				// Update operations - require admissions:update permission
				admissionsUpdate := admissionSessions.Group("")
				admissionsUpdate.Use(middleware.PermissionRequired("admissions:update"))
				{
					admissionsUpdate.PUT("/:id", admissionSessionHandler.Update)
					admissionsUpdate.PATCH("/:id/status", admissionSessionHandler.ChangeStatus)
					admissionsUpdate.PATCH("/:id/extend", admissionSessionHandler.ExtendDeadline)
					admissionsUpdate.POST("/:id/seats", admissionSessionHandler.CreateSeat)
					admissionsUpdate.PUT("/:id/seats/:seatId", admissionSessionHandler.UpdateSeat)
					admissionsUpdate.DELETE("/:id/seats/:seatId", admissionSessionHandler.DeleteSeat)
					admissionsUpdate.POST("/:id/merit-list", meritHandler.GenerateMeritList)
				}

				// Delete operations - require admissions:delete permission
				admissionsDelete := admissionSessions.Group("")
				admissionsDelete.Use(middleware.PermissionRequired("admissions:delete"))
				{
					admissionsDelete.DELETE("/:id", admissionSessionHandler.Delete)
				}
			}

			// Merit list management routes
			meritLists := protected.Group("/merit-lists")
			{
				// Read operations - require admissions:read permission
				meritListsRead := meritLists.Group("")
				meritListsRead.Use(middleware.PermissionRequired("admissions:read"))
				{
					// Note: No GET /:id here as merit lists are accessed via session
				}

				// Update operations - require admissions:update permission
				meritListsUpdate := meritLists.Group("")
				meritListsUpdate.Use(middleware.PermissionRequired("admissions:update"))
				{
					meritListsUpdate.POST("/:id/finalize", meritHandler.FinalizeMeritList)
					meritListsUpdate.PATCH("/:id/cutoff", meritHandler.UpdateCutoff)
				}

				// Delete operations - require admissions:delete permission
				meritListsDelete := meritLists.Group("")
				meritListsDelete.Use(middleware.PermissionRequired("admissions:delete"))
				{
					meritListsDelete.DELETE("/:id", meritHandler.DeleteMeritList)
				}
			}

			// Admission enquiry management routes
			enquiries := protected.Group("/enquiries")
			{
				// Read operations - require enquiries:read permission
				enquiriesRead := enquiries.Group("")
				enquiriesRead.Use(middleware.PermissionRequired("enquiries:read"))
				{
					enquiriesRead.GET("", enquiryHandler.List)
					enquiriesRead.GET("/:id", enquiryHandler.GetByID)
					enquiriesRead.GET("/:id/follow-ups", enquiryHandler.ListFollowUps)
				}

				// Create operations - require enquiries:create permission
				enquiriesCreate := enquiries.Group("")
				enquiriesCreate.Use(middleware.PermissionRequired("enquiries:create"))
				{
					enquiriesCreate.POST("", enquiryHandler.Create)
				}

				// Update operations - require enquiries:update permission
				enquiriesUpdate := enquiries.Group("")
				enquiriesUpdate.Use(middleware.PermissionRequired("enquiries:update"))
				{
					enquiriesUpdate.PUT("/:id", enquiryHandler.Update)
					enquiriesUpdate.POST("/:id/follow-ups", enquiryHandler.AddFollowUp)
					enquiriesUpdate.POST("/:id/convert", enquiryHandler.ConvertToApplication)
				}

				// Delete operations - require enquiries:delete permission
				enquiriesDelete := enquiries.Group("")
				enquiriesDelete.Use(middleware.PermissionRequired("enquiries:delete"))
				{
					enquiriesDelete.DELETE("/:id", enquiryHandler.Delete)
				}
			}

			// Admission reports and analytics routes
			admissions := protected.Group("/admissions")
			{
				// Report endpoints - require admissions:read permission
				admissionsReportsRead := admissions.Group("")
				admissionsReportsRead.Use(middleware.PermissionRequired("admissions:read"))
				{
					// Dashboard overview
					admissionsReportsRead.GET("/dashboard", admissionReportHandler.GetDashboard)

					// Report endpoints
					admissionsReportsRead.GET("/reports/funnel", admissionReportHandler.GetFunnel)
					admissionsReportsRead.GET("/reports/class-wise", admissionReportHandler.GetClassWise)
					admissionsReportsRead.GET("/reports/source-analysis", admissionReportHandler.GetSourceAnalysis)
					admissionsReportsRead.GET("/reports/daily-trend", admissionReportHandler.GetDailyTrend)

					// Export endpoint
					admissionsReportsRead.GET("/export", admissionExportHandler.Export)
				}
			}

			// Admission application management routes
			applications := protected.Group("/applications")
			{
				// Read operations - require applications:read permission
				applicationsRead := applications.Group("")
				applicationsRead.Use(middleware.PermissionRequired("applications:read"))
				{
					applicationsRead.GET("", applicationHandler.List)
					applicationsRead.GET("/:id", applicationHandler.GetByID)
					applicationsRead.GET("/:id/parents", applicationHandler.ListParents)
					applicationsRead.GET("/:id/documents", applicationHandler.ListDocuments)
					applicationsRead.GET("/:id/decision", decisionHandler.GetDecision)
				}

				// Create operations - require applications:create permission
				applicationsCreate := applications.Group("")
				applicationsCreate.Use(middleware.PermissionRequired("applications:create"))
				{
					applicationsCreate.POST("", applicationHandler.Create)
				}

				// Update operations - require applications:update permission
				applicationsUpdate := applications.Group("")
				applicationsUpdate.Use(middleware.PermissionRequired("applications:update"))
				{
					applicationsUpdate.PUT("/:id", applicationHandler.Update)
					applicationsUpdate.POST("/:id/submit", applicationHandler.Submit)
					applicationsUpdate.PATCH("/:id/stage", applicationHandler.UpdateStage)
					applicationsUpdate.POST("/:id/parents", applicationHandler.AddParent)
					applicationsUpdate.PUT("/:id/parents/:parentId", applicationHandler.UpdateParent)
					applicationsUpdate.DELETE("/:id/parents/:parentId", applicationHandler.DeleteParent)
					applicationsUpdate.POST("/:id/documents", applicationHandler.AddDocument)
					applicationsUpdate.PATCH("/:id/documents/:documentId/verify", applicationHandler.VerifyDocument)
					applicationsUpdate.DELETE("/:id/documents/:documentId", applicationHandler.DeleteDocument)
				}

				// Delete operations - require applications:delete permission
				applicationsDelete := applications.Group("")
				applicationsDelete.Use(middleware.PermissionRequired("applications:delete"))
				{
					applicationsDelete.DELETE("/:id", applicationHandler.Delete)
				}

				// Review operations - require applications:review permission
				applicationsReview := applications.Group("")
				applicationsReview.Use(middleware.PermissionRequired("applications:review"))
				{
					applicationsReview.GET("/:id/reviews", reviewHandler.GetReviews)
					applicationsReview.POST("/:id/review", reviewHandler.CreateReview)
					applicationsReview.PATCH("/:id/status", reviewHandler.UpdateStatus)
				}

				// Decision operations - require admissions:update permission
				applicationsDecision := applications.Group("")
				applicationsDecision.Use(middleware.PermissionRequired("admissions:update"))
				{
					applicationsDecision.POST("/:id/decision", decisionHandler.MakeDecision)
					applicationsDecision.POST("/:id/offer-letter", decisionHandler.GenerateOfferLetter)
					applicationsDecision.POST("/:id/accept-offer", decisionHandler.AcceptOffer)
					applicationsDecision.POST("/:id/enroll", decisionHandler.Enroll)
					applicationsDecision.POST("/:id/promote", decisionHandler.PromoteFromWaitlist)
					applicationsDecision.PATCH("/:id/waitlist-position", decisionHandler.UpdateWaitlistPosition)
				}

				// Bulk decision operations - require admissions:update permission
				applicationsBulk := applications.Group("")
				applicationsBulk.Use(middleware.PermissionRequired("admissions:update"))
				{
					applicationsBulk.POST("/bulk-decision", decisionHandler.MakeBulkDecision)
				}
			}

			// Entrance test management routes
			entranceTests := protected.Group("/entrance-tests")
			{
				// Read operations - require tests:read permission
				testsRead := entranceTests.Group("")
				testsRead.Use(middleware.PermissionRequired("tests:read"))
				{
					testsRead.GET("", testHandler.ListTests)
					testsRead.GET("/:id", testHandler.GetTest)
					testsRead.GET("/:id/registrations", testHandler.ListRegistrations)
					testsRead.GET("/:id/hall-tickets", testHandler.GetHallTickets)
					testsRead.GET("/:id/hall-tickets/:registrationId", testHandler.GetHallTicket)
				}

				// Create operations - require tests:create permission
				testsCreate := entranceTests.Group("")
				testsCreate.Use(middleware.PermissionRequired("tests:create"))
				{
					testsCreate.POST("", testHandler.CreateTest)
				}

				// Update operations - require tests:update permission
				testsUpdate := entranceTests.Group("")
				testsUpdate.Use(middleware.PermissionRequired("tests:update"))
				{
					testsUpdate.PUT("/:id", testHandler.UpdateTest)
					testsUpdate.POST("/:id/register", testHandler.RegisterCandidate)
					testsUpdate.DELETE("/:id/registrations/:registrationId", testHandler.CancelRegistration)
				}

				// Results management - require tests:manage permission
				testsManage := entranceTests.Group("")
				testsManage.Use(middleware.PermissionRequired("tests:manage"))
				{
					testsManage.POST("/:id/results", testHandler.SubmitResult)
					testsManage.POST("/:id/results/bulk", testHandler.BulkSubmitResults)
				}

				// Delete operations - require tests:delete permission
				testsDelete := entranceTests.Group("")
				testsDelete.Use(middleware.PermissionRequired("tests:delete"))
				{
					testsDelete.DELETE("/:id", testHandler.DeleteTest)
				}
			}

			// Department management routes
			departmentHandler.RegisterRoutes(protected, middleware.AuthRequired(jwtService))

			// Designation management routes
			designationHandler.RegisterRoutes(protected, middleware.AuthRequired(jwtService))

			// Staff management routes
			staffRoutes := protected.Group("/staff")
			{
				// Read operations - require staff:read permission
				staffRead := staffRoutes.Group("")
				staffRead.Use(middleware.PermissionRequired("staff:read"))
				{
					staffRead.GET("", staffHandler.List)
					staffRead.GET("/employee-id/preview", staffHandler.PreviewEmployeeID)
					staffRead.GET("/:id", staffHandler.Get)
					staffRead.GET("/:id/status-history", staffHandler.GetStatusHistory)
				}

				// Create operations - require staff:create permission
				staffCreate := staffRoutes.Group("")
				staffCreate.Use(middleware.PermissionRequired("staff:create"))
				{
					staffCreate.POST("", staffHandler.Create)
				}

				// Update operations - require staff:update permission
				staffUpdate := staffRoutes.Group("")
				staffUpdate.Use(middleware.PermissionRequired("staff:update"))
				{
					staffUpdate.PUT("/:id", staffHandler.Update)
					staffUpdate.PATCH("/:id/status", staffHandler.UpdateStatus)
					staffUpdate.POST("/:id/photo", staffHandler.UpdatePhoto)
				}

				// Delete operations - require staff:delete permission
				staffDelete := staffRoutes.Group("")
				staffDelete.Use(middleware.PermissionRequired("staff:delete"))
				{
					staffDelete.DELETE("/:id", staffHandler.Delete)
				}
			}

			// Attendance management routes
			attendanceRoutes := protected.Group("/attendance")
			{
				// Self-service attendance (requires attendance:mark_self permission)
				attendanceSelf := attendanceRoutes.Group("")
				attendanceSelf.Use(middleware.PermissionRequired("attendance:mark_self"))
				{
					attendanceSelf.POST("/check-in", attendanceHandler.CheckIn)
					attendanceSelf.POST("/check-out", attendanceHandler.CheckOut)
				}

				// View own attendance (requires attendance:view_self permission)
				attendanceViewSelf := attendanceRoutes.Group("/my")
				attendanceViewSelf.Use(middleware.PermissionRequired("attendance:view_self"))
				{
					attendanceViewSelf.GET("", attendanceHandler.GetMyAttendance)
					attendanceViewSelf.GET("/today", attendanceHandler.GetMyToday)
					attendanceViewSelf.GET("/summary", attendanceHandler.GetMySummary)
				}

				// View all attendance (requires attendance:view_all permission)
				attendanceViewAll := attendanceRoutes.Group("")
				attendanceViewAll.Use(middleware.PermissionRequired("attendance:view_all"))
				{
					attendanceViewAll.GET("", attendanceHandler.List)
				}

				// Mark attendance for others (requires attendance:mark_others permission)
				attendanceMark := attendanceRoutes.Group("")
				attendanceMark.Use(middleware.PermissionRequired("attendance:mark_others"))
				{
					attendanceMark.POST("/mark", attendanceHandler.MarkAttendance)
				}

				// Regularization - submit (requires attendance:regularize permission)
				attendanceRegularize := attendanceRoutes.Group("/regularization")
				attendanceRegularize.Use(middleware.PermissionRequired("attendance:regularize"))
				{
					attendanceRegularize.POST("", attendanceHandler.SubmitRegularization)
				}

				// Regularization - approve/reject (requires attendance:approve_regularization permission)
				attendanceApprove := attendanceRoutes.Group("/regularization")
				attendanceApprove.Use(middleware.PermissionRequired("attendance:approve_regularization"))
				{
					attendanceApprove.GET("", attendanceHandler.ListRegularizations)
					attendanceApprove.PUT("/:id/approve", attendanceHandler.ApproveRegularization)
					attendanceApprove.PUT("/:id/reject", attendanceHandler.RejectRegularization)
				}

				// Settings management (requires attendance:settings permission)
				attendanceSettings := attendanceRoutes.Group("/settings")
				attendanceSettings.Use(middleware.PermissionRequired("attendance:settings"))
				{
					attendanceSettings.GET("", attendanceHandler.GetSettings)
					attendanceSettings.PUT("", attendanceHandler.UpdateSettings)
				}
			}

			// Student Attendance management routes
			studentAttendanceRoutes := protected.Group("/student-attendance")
			{
				// Get teacher's assigned classes for attendance marking
				studentAttendanceRoutes.GET("/my-classes", studentAttendanceHandler.GetMyClasses)

				// Class attendance operations (requires student_attendance:mark_class permission)
				classAttendance := studentAttendanceRoutes.Group("/class")
				classAttendance.Use(middleware.PermissionRequired("student_attendance:mark_class"))
				{
					classAttendance.GET("/:id", studentAttendanceHandler.GetClassAttendance)
					classAttendance.POST("/:id", studentAttendanceHandler.MarkClassAttendance)
				}

				// Period-wise attendance operations (Story 7.2)
				periodAttendance := studentAttendanceRoutes.Group("")
				periodAttendance.Use(middleware.PermissionRequired("student_attendance:mark_class"))
				{
					// Get periods for a section on a date
					periodAttendance.GET("/periods", studentAttendanceHandler.GetPeriods)
					// Get attendance for a specific period
					periodAttendance.GET("/period/:id", studentAttendanceHandler.GetPeriodAttendance)
					// Mark attendance for a specific period
					periodAttendance.POST("/period/:id", studentAttendanceHandler.MarkPeriodAttendance)
					// Get daily summary (all periods aggregated)
					periodAttendance.GET("/daily-summary", studentAttendanceHandler.GetDailySummary)
				}

				// Subject-wise attendance analytics
				subjectAttendance := studentAttendanceRoutes.Group("/subject")
				subjectAttendance.Use(middleware.PermissionRequired("student_attendance:view_class"))
				{
					subjectAttendance.GET("/:id", studentAttendanceHandler.GetSubjectAttendance)
				}

				// View all student attendance (requires student_attendance:view_all permission)
				viewAll := studentAttendanceRoutes.Group("")
				viewAll.Use(middleware.PermissionRequired("student_attendance:view_all"))
				{
					viewAll.GET("", studentAttendanceHandler.ListAttendance)
				}

				// Settings management (requires student_attendance:manage_settings permission)
				studentAttendanceSettings := studentAttendanceRoutes.Group("/settings")
				studentAttendanceSettings.Use(middleware.PermissionRequired("student_attendance:manage_settings"))
				{
					studentAttendanceSettings.GET("", studentAttendanceHandler.GetSettings)
					studentAttendanceSettings.PUT("", studentAttendanceHandler.UpdateSettings)
				}

				// Edit and audit trail routes (Story 7.3)
				// Edit attendance - requires mark_class permission
				editRoutes := studentAttendanceRoutes.Group("")
				editRoutes.Use(middleware.PermissionRequired("student_attendance:mark_class"))
				{
					// Edit an attendance record (with reason)
					editRoutes.PUT("/:id", studentAttendanceHandler.EditAttendance)
					// Get edit window status
					editRoutes.GET("/:id/edit-status", studentAttendanceHandler.GetEditWindowStatus)
				}

				// View audit history - requires view_class permission
				historyRoutes := studentAttendanceRoutes.Group("")
				historyRoutes.Use(middleware.PermissionRequired("student_attendance:view_class"))
				{
					// Get audit trail for an attendance record
					historyRoutes.GET("/:id/history", studentAttendanceHandler.GetAttendanceHistory)
				}

				// Calendar and summary routes (Story 7.4) - for students/parents
				calendarRoutes := studentAttendanceRoutes.Group("")
				calendarRoutes.Use(middleware.PermissionRequired("student_attendance:view_self"))
				{
					calendarRoutes.GET("/calendar/:studentId", studentAttendanceHandler.GetStudentCalendar)
					calendarRoutes.GET("/summary/:studentId", studentAttendanceHandler.GetStudentSummary)
				}

				// Reports routes (Stories 7.5, 7.6) - for teachers/admins
				reportRoutes := studentAttendanceRoutes.Group("/reports")
				reportRoutes.Use(middleware.PermissionRequired("student_attendance:view_reports"))
				{
					reportRoutes.GET("/class/:sectionId", studentAttendanceHandler.GetClassReport)
					reportRoutes.GET("/class/:sectionId/monthly", studentAttendanceHandler.GetMonthlyClassReport)
					reportRoutes.GET("/daily", studentAttendanceHandler.GetDailyReport)
				}

				// Alerts routes (Stories 7.7, 7.8) - for admins
				alertRoutes := studentAttendanceRoutes.Group("/alerts")
				alertRoutes.Use(middleware.PermissionRequired("student_attendance:view_alerts"))
				{
					alertRoutes.GET("/low-attendance", studentAttendanceHandler.GetLowAttendanceDashboard)
					alertRoutes.GET("/unmarked", studentAttendanceHandler.GetUnmarkedAttendance)
				}
			}

			// Salary management routes
			salaryHandler.RegisterRoutes(protected, middleware.AuthRequired(jwtService))
			salaryHandler.RegisterStaffSalaryRoutes(staffRoutes)

			// Payroll management routes
			payrollHandler.RegisterRoutes(protected)
			payrollHandler.RegisterStaffPayslipRoutes(staffRoutes)

			// Academic structure routes (classes, sections, streams)
			academicHandler.RegisterRoutes(protected)

			// Timetable structure routes (shifts, day patterns, period slots)
			timetableHandler.RegisterRoutes(protected)

			// Substitution management routes
			timetableHandler.RegisterSubstitutionRoutes(protected)

			// Exam type management routes
			examHandler.RegisterRoutes(protected)

			// Examination management routes
			examinationHandler.RegisterRoutes(protected)

			// Hall ticket management routes
			hallTicketHandler.RegisterRoutes(protected, middleware.AuthRequired(jwtService))

			// Teacher assignment routes
			assignmentHandler.RegisterRoutes(protected)
			assignmentHandler.RegisterStaffRoutes(staffRoutes)
			assignmentHandler.RegisterClassRoutes(protected)

			// Staff document management routes
			staffDocumentHandler.RegisterDocumentTypeRoutes(protected)
			staffDocumentHandler.RegisterStaffDocumentRoutes(staffRoutes)
			staffDocumentHandler.RegisterGlobalDocumentRoutes(protected)
		}

		// Feature flags routes (authenticated - returns flags for current user)
		featureFlagsRoutes := v1.Group("/feature-flags")
		featureFlagsRoutes.Use(middleware.AuthRequired(jwtService))
		featureFlagsRoutes.Use(middleware.FeatureFlagDefault(featureFlagService))
		{
			featureFlagsRoutes.GET("", featureFlagHandler.GetCurrentFlags)
			featureFlagsRoutes.GET("/:key", featureFlagHandler.IsEnabled)
		}

		// Admin routes (require admin permissions)
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(middleware.TenantRequired())
		adminRoutes.Use(middleware.AuthRequired(jwtService))
		{
			// Feature flag management (admin only)
			adminFlags := adminRoutes.Group("/feature-flags")
			adminFlags.Use(middleware.PermissionRequired("settings:write"))
			{
				adminFlags.GET("", featureFlagHandler.ListFlags)
				adminFlags.GET("/:id", featureFlagHandler.GetFlag)
				adminFlags.POST("", featureFlagHandler.CreateFlag)
				adminFlags.PUT("/:id", featureFlagHandler.UpdateFlag)
				adminFlags.DELETE("/:id", featureFlagHandler.DeleteFlag)
			}

			// Tenant feature flag overrides (admin only)
			adminTenants := adminRoutes.Group("/tenants")
			adminTenants.Use(middleware.PermissionRequired("settings:write"))
			{
				adminTenants.GET("/:id/feature-flags", featureFlagHandler.GetTenantFlags)
				adminTenants.PUT("/:id/feature-flags", featureFlagHandler.SetTenantFlags)
			}

			// User feature flag overrides (admin only - for beta testing)
			adminUsers := adminRoutes.Group("/users")
			adminUsers.Use(middleware.PermissionRequired("settings:write"))
			{
				adminUsers.GET("/:id/feature-flags", featureFlagHandler.GetUserFlags)
				adminUsers.PUT("/:id/feature-flags", featureFlagHandler.SetUserFlags)
			}
		}
	}

	return router
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(c *gin.Context) {
	response.OK(c, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func readyHandler(c *gin.Context) {
	// TODO: Add database and cache connectivity checks
	response.OK(c, HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func pingHandler(c *gin.Context) {
	response.OK(c, gin.H{"message": "pong"})
}
