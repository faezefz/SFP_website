package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	db "github.com/faezefz/SFP_website/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Server struct
type Server struct {
	Db     *db.Queries
	Router *gin.Engine
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
	s.Router.Use(cors.Default())

	s.Router.GET("/", s.home)
	s.Router.POST("/login", s.login)
	s.Router.POST("/signup", s.signup)

	auth := s.Router.Group("/")
	auth.Use(s.authMiddleware())
	{
		auth.GET("/dashboard", s.userDashboard)
		auth.POST("/import-dataset", s.uploadDataset)
		auth.GET("/datasets/:owner_user_id", s.listDatasets)
		auth.POST("/projects", s.createProject)
		auth.GET("/projects/:owner_user_id", s.getProjectsByOwnerID)
		auth.PUT("/projects/:project_id", s.updateProject)
		auth.DELETE("/projects/:project_id", s.deleteProject)
		auth.GET("/projects/project/:project_id", s.getProjectByID)
		auth.GET("/dataset/:dataset_id", s.getDatasetByID)
	}
}

// Run
func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

// home
func (s *Server) home(c *gin.Context) {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
}

// authMiddleware
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Place holder
		c.Next()
	}
}

// uploadDataset
func (s *Server) uploadDataset(c *gin.Context) {
	for key, values := range c.Request.Form {
		fmt.Printf("%s: %v\n", key, values)
	}

	file, _, err := c.Request.FormFile("content")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	name := c.PostForm("name")
	description := c.PostForm("description")
	userID := c.PostForm("user_id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userIDInt32 := int32(userIDInt)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	arg := db.CreateDatasetParams{
		UserID:      pgtype.Int4{Int32: userIDInt32, Valid: true},
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
		Content:     fileContent,
	}

	dataset, err := s.Db.CreateDataset(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store dataset in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dataset_id": dataset.ID,
	})
}

// listDatasets
func (s *Server) listDatasets(c *gin.Context) {
	userID := c.Param("owner_user_id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	userIDParam := pgtype.Int4{Int32: int32(userIDInt), Valid: true}

	datasets, err := s.Db.GetDatasetsByUserID(context.Background(), userIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch datasets"})
		return
	}

	if len(datasets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No datasets found for this user"})
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

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID in URL"})
		return
	}

	userID := req.ID.Int32
	if !req.ID.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

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

	ownerUserIDInt, err := strconv.Atoi(ownerUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_user_id format"})
		return
	}

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

	err = s.Db.DeleteProject(context.Background(), projectIDInt32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// get dataset by dataset id
func (s *Server) getDatasetByID(c *gin.Context) {
	datasetID := c.Param("dataset_id")
	datasetIDInt, err := strconv.Atoi(datasetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dataset id format"})
		return
	}

	content, err := s.Db.GetDatasetContentByID(context.Background(), int32(datasetIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dataset content"})
		return
	}
	fmt.Println("Content:", content)
	contentBase64 := base64.StdEncoding.EncodeToString(content)

	c.JSON(http.StatusOK, gin.H{"content": contentBase64})
}

func (s *Server) getProjectByID(c *gin.Context) {
	projectID := c.Param("project_id")
	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project id format"})
		return
	}

	project, err := s.Db.GetProjectByID(context.Background(), int32(projectIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	c.JSON(http.StatusOK, project)
}
