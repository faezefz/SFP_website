package db

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createRandomProject(t *testing.T, ownerID int32) Project {
	arg := CreateProjectParams{
		OwnerUserID: ownerID,
		Name:        randomString(10),
		Description: pgtype.Text{String: randomString(20), Valid: true},
	}

	project, err := testQueries.CreateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	require.Equal(t, arg.OwnerUserID, project.OwnerUserID)
	require.Equal(t, arg.Name, project.Name)
	require.Equal(t, arg.Description.String, project.Description.String)

	return project
}

func TestCreateProject(t *testing.T) {
	user := createRandomUser(t)
	project := createRandomProject(t, user.ID)
	require.NotEmpty(t, project)
}

// تست GetProjectByID
func TestGetProjectByID(t *testing.T) {
	user := createRandomUser(t)
	project1 := createRandomProject(t, user.ID)

	project2, err := testQueries.GetProjectByID(context.Background(), project1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, project2)
	require.Equal(t, project1.ID, project2.ID)
	require.Equal(t, project1.Name, project2.Name)
}

// تست GetProjectsByOwnerID
func TestGetProjectsByOwnerID(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 3; i++ {
		createRandomProject(t, user.ID)
	}

	projects, err := testQueries.GetProjectsByOwnerID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, projects)
	require.GreaterOrEqual(t, len(projects), 3)

	for _, p := range projects {
		require.Equal(t, user.ID, p.OwnerUserID)
	}
}

// تست UpdateProject
func TestUpdateProject(t *testing.T) {
	user := createRandomUser(t)
	project1 := createRandomProject(t, user.ID)

	arg := UpdateProjectParams{
		ID:          project1.ID,
		Name:        randomString(12),
		Description: pgtype.Text{String: randomString(25), Valid: true},
	}

	project2, err := testQueries.UpdateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project2)

	require.Equal(t, arg.Name, project2.Name)
	require.Equal(t, arg.Description.String, project2.Description.String)
}

// تست DeleteProject
func TestDeleteProject(t *testing.T) {
	user := createRandomUser(t)
	project := createRandomProject(t, user.ID)

	err := testQueries.DeleteProject(context.Background(), project.ID)
	require.NoError(t, err)

	// مطمئن شو که پروژه حذف شده
	project2, err := testQueries.GetProjectByID(context.Background(), project.ID)
	require.Error(t, err)
	require.Empty(t, project2)
}

func TestProjectDatasets(t *testing.T) {
	user := createRandomUser(t)
	project := createRandomProject(t, user.ID)

	// ایجاد دیتاست
	dataset := createRandomDataset(t)

	// افزودن دیتاست به پروژه
	err := testQueries.AddDatasetToProject(context.Background(), AddDatasetToProjectParams{
		ProjectID: project.ID,
		DatasetID: dataset.ID,
	})
	require.NoError(t, err)

	// بررسی گرفتن دیتاست‌ها
	datasets, err := testQueries.GetDatasetsByProjectID(context.Background(), project.ID)
	require.NoError(t, err)
	require.Len(t, datasets, 1)
	require.Equal(t, dataset.ID, datasets[0].ID)

	// حذف دیتاست از پروژه
	err = testQueries.RemoveDatasetFromProject(context.Background(), RemoveDatasetFromProjectParams{
		ProjectID: project.ID,
		DatasetID: dataset.ID,
	})
	require.NoError(t, err)

	// بررسی اینکه بعد از حذف خالی شده
	datasets, err = testQueries.GetDatasetsByProjectID(context.Background(), project.ID)
	require.NoError(t, err)
	require.Len(t, datasets, 0)
}

func TestProjectModels(t *testing.T) {
	user := createRandomUser(t)
	project := createRandomProject(t, user.ID)

	model := createRandomModel(t, user)

	err := testQueries.AddModelToProject(context.Background(), AddModelToProjectParams{
		ProjectID: project.ID,
		ModelID:   model.ID,
	})
	require.NoError(t, err)

	// بررسی گرفتن مدل‌ها
	models, err := testQueries.GetModelsByProjectID(context.Background(), project.ID)
	require.NoError(t, err)
	require.Len(t, models, 1)
	require.Equal(t, model.ID, models[0].ID)

	// حذف مدل از پروژه
	err = testQueries.RemoveModelFromProject(context.Background(), RemoveModelFromProjectParams{
		ProjectID: project.ID,
		ModelID:   model.ID,
	})
	require.NoError(t, err)

	// بررسی اینکه بعد از حذف خالی شده
	models, err = testQueries.GetModelsByProjectID(context.Background(), project.ID)
	require.NoError(t, err)
	require.Len(t, models, 0)
}
