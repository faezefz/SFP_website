package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	db "github.com/faezefz/SFP_website/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Server struct
type Server struct {
	Db     *db.Queries // فیلد db هم می‌تواند exported باشد (اگر نیاز به دسترسی از خارج پکیج است)
	Router *gin.Engine // تغییر از router به Router (با حرف بزرگ)
	Secret string
}

// NewServer
func NewServer(dbPool *pgxpool.Pool, secret string) *Server {
	server := &Server{
		Db:     db.New(dbPool),
		Router: gin.Default(),
		Secret: secret,
	}

	server.Routes()

	return server
}

// routes
func (s *Server) Routes() {
	s.Router.GET("/", s.home)

	s.Router.POST("/signup", s.signup)
	s.Router.POST("/login", s.login)

	auth := s.Router.Group("/")
	auth.Use(s.authMiddleware())
	{
		auth.GET("/dashboard/:id", s.userDashboard) // تغییر برای دریافت :id از URI
		auth.POST("/datasets", s.uploadDataset)
		auth.GET("/datasets", s.listDatasets)
		// سایر روت‌ها: مدل‌ها، پیش‌بینی‌ها و غیره می‌توان اضافه کرد
	}
}

// Run
func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

// home
func (s *Server) home(c *gin.Context) {
	tokenString := c.GetHeader("Authorization") // گرفتن توکن از هدر درخواست

	if tokenString == "" {
		// اگر توکن موجود نباشد، کاربر را به صفحه ورود هدایت می‌کنیم
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please log in first."})
		return
	}

	// بررسی صحت توکن
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token. Please log in again."})
		return
	}

	// اگر توکن معتبر باشد، به داشبورد هدایت می‌شود
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims."})
		return
	}

	userID := claims["user_id"].(float64)
	c.Redirect(http.StatusFound, fmt.Sprintf("/dashboard/%d", int(userID)))
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     pgtype.Text{String: req.FullName, Valid: req.FullName != ""},
	}

	user, err := s.Db.CreateUser(context.Background(), arg)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.Db.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// authMiddleware
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := authHeader[len(prefix):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
			return
		}

		c.Set("user_id", int32(userID))

		c.Next()
	}
}

// uploadDataset
func (s *Server) uploadDataset(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDInterface.(int32)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}

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

	arg := db.CreateDatasetParams{
		UserID: pgtype.Int4{Int32: userID, Valid: true},
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
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDInterface.(int32)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
		return
	}
	userIDParam := pgtype.Int4{Int32: userID, Valid: true}
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
