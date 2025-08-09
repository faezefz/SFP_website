package api

import (
	"context"
	_ "fmt"
	"log"
	"net/http"

	db "github.com/faezefz/SFP_website/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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
		auth.GET("/dashboard/:id", s.userDashboard) // صفحه داشبورد
		auth.POST("/datasets", s.uploadDataset)     // آپلود داده
		auth.GET("/datasets", s.listDatasets)       // نمایش داده‌ها
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

	// ذخیره‌سازی پسورد بدون هش کردن
	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: req.Password, // ذخیره‌سازی پسورد بدون هش کردن
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

	// مقایسه پسورد وارد شده با پسورد ذخیره شده (بدون هش)
	if user.PasswordHash != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// لاگین موفق، فقط پیام تایید باز می‌گردد
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
	// برای سادگی، ما نیازی به احراز هویت نداریم
	type uploadRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		FilePath    string `json:"file_path" binding:"required"`
	}

	var req uploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id") // به سادگی، از اطلاعات کاربر موجود استفاده می‌کنیم
	arg := db.CreateDatasetParams{
		UserID: pgtype.Int4{Int32: userID.(int32), Valid: true},
		Name:   req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		FilePath: req.FilePath,
	}

	dataset, err := s.Db.CreateDataset(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dataset"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"dataset_id": dataset.ID})
}

// listDatasets
func (s *Server) listDatasets(c *gin.Context) {
	// به سادگی داده‌ها را نمایش می‌دهیم
	userID, _ := c.Get("user_id")
	userIDParam := pgtype.Int4{Int32: userID.(int32), Valid: true}
	datasets, err := s.Db.ListDatasetsByUserID(context.Background(), userIDParam)
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

	datasets, err := s.Db.ListDatasetsByUserID(context.Background(), pgtype.Int4{Int32: userID, Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasets"})
		return
	}

	// ارسال اطلاعات پروفایل یا داشبورد
	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome to your dashboard",
		"user_id":  userID,
		"datasets": datasets,
	})
}
