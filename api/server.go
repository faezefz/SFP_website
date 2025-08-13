package api

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	db "github.com/faezefz/SFP_website/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt" // اضافه کردن این خط برای استفاده از bcrypt
)

// Server struct
type Server struct {
	Db     *db.Queries // فیلد db هم می‌تواند exported باشد (اگر نیاز به دسترسی از خارج پکیج است)
	Router *gin.Engine // تغییر از router به Router (با حرف بزرگ)
}

// NewServer
func NewServer(dbPool *pgxpool.Pool) *Server {
	server := &Server{
		Db:     db.New(dbPool),
		Router: gin.Default(),
	}

	server.Routes()

	return server
}

// routes
func (s *Server) Routes() {
	// فعال‌سازی CORS برای همه روت‌ها
	s.Router.Use(cors.Default()) // استفاده از تنظیمات پیش‌فرض CORS

	// مسیرهایی که نیازی به احراز هویت ندارند:
	s.Router.GET("/", s.home)          // صفحه اصلی
	s.Router.POST("/login", s.login)   // صفحه ورود
	s.Router.POST("/signup", s.signup) // ثبت‌نام

	// این گروه فقط برای مسیرهایی که نیاز به احراز هویت دارند:
	auth := s.Router.Group("/")
	auth.Use(s.authMiddleware()) // فقط این گروه به احراز هویت نیاز دارد
	{
		auth.GET("/dashboard", s.userDashboard) // صفحه داشبورد
		auth.POST("/datasets", s.uploadDataset) // آپلود داده
		auth.GET("/datasets", s.listDatasets)
		auth.POST("/projects", s.createProject)                      // ایجاد پروژه
		auth.GET("/projects/:owner_user_id", s.getProjectsByOwnerID) // دریافت پروژه‌ها بر اساس owner_user_id
		auth.PUT("/projects/:project_id", s.updateProject)           // ویرایش پروژه
		auth.DELETE("/projects/:project_id", s.deleteProject)        // نمایش داده‌ها
	}
}

// Run
func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

// home
func (s *Server) home(c *gin.Context) {
	// ارسال یک پاسخ JSON به فلاتر
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the API, please use /login or /signup.",
	})
}

// signup
func (s *Server) signup(c *gin.Context) {
	type signupRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name"`
	}

	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// هش کردن پسورد قبل از ذخیره در دیتابیس
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// ذخیره‌سازی پسورد هش شده
	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword), // پسورد هش شده را ذخیره می‌کنیم
		FullName:     pgtype.Text{String: req.FullName, Valid: req.FullName != ""},
	}

	user, err := s.Db.CreateUser(context.Background(), arg)
	if err != nil {
		log.Printf("Error creating user: %v", arg.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": user.ID, "email": user.Email})
}

// login
func (s *Server) login(c *gin.Context) {
	type loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format or missing fields"})
		return
	}

	// جستجو برای کاربر در دیتابیس
	user, err := s.Db.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// مقایسه پسورد وارد شده با پسورد هش شده در دیتابیس
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// لاگین موفق
	c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
}

// authMiddleware
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// هیچ احراز هویتی برای این مسیرها لازم نیست
		c.Next()
	}
}

// uploadDataset
func (s *Server) uploadDataset(c *gin.Context) {
	// تعریف ساختار درخواست
	type uploadRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	var req uploadRequest

	// بررسی پارامترهای ورودی
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// دریافت فایل CSV از درخواست
	file, _, err := c.Request.FormFile("content")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// خواندن فایل به صورت بایت
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// گرفتن ID کاربر از سشن یا کانتکست
	userID, _ := c.Get("user_id")

	// آماده‌سازی داده‌ها برای ذخیره در پایگاه داده
	arg := db.CreateDatasetParams{
		UserID: pgtype.Int4{Int32: userID.(int32), Valid: true},
		Name:   req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		Content: fileContent, // ذخیره فایل به صورت بایت
	}

	// ایجاد دیتاست در پایگاه داده
	dataset, err := s.Db.CreateDataset(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dataset"})
		return
	}

	// ارسال پاسخ موفقیت‌آمیز
	c.JSON(http.StatusCreated, gin.H{"dataset_id": dataset.ID})
}

// listDatasets
func (s *Server) listDatasets(c *gin.Context) {
	// به سادگی داده‌ها را نمایش می‌دهیم
	userID, _ := c.Get("user_id")
	userIDParam := pgtype.Int4{Int32: userID.(int32), Valid: true}
	datasets, err := s.Db.GetDatasetsByUserID(context.Background(), userIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasets"})
		return
	}

	c.JSON(http.StatusOK, datasets)
}

// dashboard
type dashboardRequest struct {
	ID pgtype.Int4 `uri:"id" binding:"required"`
}

func (s *Server) userDashboard(c *gin.Context) {
	var req dashboardRequest

	// دریافت پارامترهای URI
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID in URL"})
		return
	}

	// پر کردن Int4 با مقدار id
	userID := req.ID.Int32
	if !req.ID.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// ارسال اطلاعات پروفایل یا داشبورد
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your dashboard",
		"user_id": userID,
	})
}

// createProject
func (s *Server) createProject(c *gin.Context) {
	type createProjectRequest struct {
		OwnerUserID int32  `json:"owner_user_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	print(c.Request.Body)
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ذخیره پروژه در دیتابیس
	arg := db.CreateProjectParams{
		OwnerUserID: req.OwnerUserID,
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	}

	project, err := s.Db.CreateProject(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// getProjectsByOwnerID
func (s *Server) getProjectsByOwnerID(c *gin.Context) {
	ownerUserID := c.Param("owner_user_id")

	// تبدیل شناسه کاربر از string به int32
	ownerUserIDInt, err := strconv.Atoi(ownerUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_user_id format"})
		return
	}

	// دریافت پروژه‌ها از دیتابیس با استفاده از شناسه کاربر
	projects, err := s.Db.GetProjectsByOwnerID(context.Background(), int32(ownerUserIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// updateProject
func (s *Server) updateProject(c *gin.Context) {
	type updateProjectRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	projectID := c.Param("project_id")
	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id format"})
		return
	}

	var req updateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ویرایش پروژه در دیتابیس
	arg := db.UpdateProjectParams{
		ID:          int32(projectIDInt),
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	}

	project, err := s.Db.UpdateProject(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// deleteProject
func (s *Server) deleteProject(c *gin.Context) {
	projectID := c.Param("project_id")
	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id format"})
		return
	}
	projectIDInt32 := int32(projectIDInt)

	// حذف پروژه از دیتابیس
	err = s.Db.DeleteProject(context.Background(), projectIDInt32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
