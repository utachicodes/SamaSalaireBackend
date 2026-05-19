package router

import (
	"net/http"

	"samasalaire-backend/internal/handlers"
	"samasalaire-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func New(db *mongo.Database) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authH := handlers.NewAuthHandler(db)
	empH := handlers.NewEmployeeHandler(db)
	salH := handlers.NewSalaryHandler(db)
	payH := handlers.NewPayrollHandler(db)
	leaveH := handlers.NewLeaveHandler(db)
	userH := handlers.NewUserHandler(db)
	repH := handlers.NewReportHandler(db)

	auth := middleware.AuthRequired()
	hrAdmin := middleware.RequireRole("hr", "admin")
	hrOnly := middleware.RequireRole("hr")
	adminOnly := middleware.RequireRole("admin")
	hrAdminMgr := middleware.RequireRole("hr", "admin", "manager")

	api := r.Group("/api")

	// Auth
	api.POST("/auth/login", authH.Login)
	api.POST("/auth/logout", auth, authH.Logout)

	// Employees
	emp := api.Group("/employees", auth)
	emp.GET("", hrAdminMgr, empH.List)
	emp.POST("", hrAdmin, middleware.AuditLog(db, "create_employee"), empH.Create)
	emp.GET("/:id", empH.Get)
	emp.PUT("/:id", hrAdmin, middleware.AuditLog(db, "update_employee"), empH.Update)
	emp.DELETE("/:id", adminOnly, middleware.AuditLog(db, "delete_employee"), empH.Delete)

	// Salary components
	sal := api.Group("/salary-components", auth)
	sal.GET("/:employeeId", hrAdmin, salH.ListByEmployee)
	sal.POST("", hrOnly, middleware.AuditLog(db, "create_salary_component"), salH.Create)
	sal.PUT("/:id", hrOnly, middleware.AuditLog(db, "update_salary_component"), salH.Update)
	sal.DELETE("/:id", hrOnly, middleware.AuditLog(db, "delete_salary_component"), salH.Delete)

	// Payroll periods
	pp := api.Group("/payroll-periods", auth, hrAdmin)
	pp.GET("", payH.ListPeriods)
	pp.POST("", middleware.AuditLog(db, "create_payroll_period"), payH.CreatePeriod)
	pp.POST("/:id/run", hrOnly, middleware.AuditLog(db, "run_payroll"), payH.RunPayroll)
	pp.POST("/:id/finalize", hrOnly, middleware.AuditLog(db, "finalize_payroll"), payH.FinalizePeriod)

	// Payslips
	ps := api.Group("/payslips", auth)
	ps.GET("", payH.ListPayslips)
	ps.GET("/:id", payH.GetPayslip)

	// Leave types
	lt := api.Group("/leave-types", auth)
	lt.GET("", leaveH.ListLeaveTypes)
	lt.POST("", adminOnly, middleware.AuditLog(db, "create_leave_type"), leaveH.CreateLeaveType)
	lt.PUT("/:id", adminOnly, middleware.AuditLog(db, "update_leave_type"), leaveH.UpdateLeaveType)

	// Leave balances
	api.GET("/leave-balances/:employeeId", auth, leaveH.GetBalance)

	// Leave requests
	lr := api.Group("/leave-requests", auth)
	lr.GET("", leaveH.ListRequests)
	lr.POST("", middleware.AuditLog(db, "create_leave_request"), leaveH.CreateRequest)
	lr.PUT("/:id/decide", hrAdminMgr, middleware.AuditLog(db, "decide_leave_request"), leaveH.DecideRequest)

	// Reports
	rep := api.Group("/reports", auth, hrAdmin)
	rep.GET("/payroll-summary", repH.PayrollSummary)
	rep.GET("/leave-summary", repH.LeaveSummary)

	// Users
	usr := api.Group("/users", auth, adminOnly)
	usr.GET("", userH.List)
	usr.POST("", middleware.AuditLog(db, "create_user"), userH.Create)
	usr.PUT("/:id", middleware.AuditLog(db, "update_user"), userH.Update)

	return r
}
