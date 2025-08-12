package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomModel(t *testing.T, user User) Model {

	arg := CreateModelParams{
		UserID:      pgtype.Int4{Int32: user.ID, Valid: true},
		Name:        "Test Model",
		Description: pgtype.Text{String: "This is a test model", Valid: true},
		FilePath:    "/tmp/model.bin",
	}

	model, err := testQueries.CreateModel(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, model)

	require.Equal(t, arg.UserID.Int32, model.UserID.Int32)
	require.Equal(t, arg.Name, model.Name)
	require.Equal(t, arg.Description.String, model.Description.String)
	require.Equal(t, arg.FilePath, model.FilePath)

	return model
}

func TestCreateModel(t *testing.T) {
	user := createRandomUser(t)
	createRandomModel(t, user)
}

func TestGetModelByID(t *testing.T) {
	user := createRandomUser(t)
	model := createRandomModel(t, user)

	fetchedModel, err := testQueries.GetModelByID(context.Background(), model.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedModel)
	require.Equal(t, model.ID, fetchedModel.ID)
}

func TestGetModelsByUserID(t *testing.T) {
	user := createRandomUser(t)
	model := createRandomModel(t, user)

	models, err := testQueries.GetModelsByUserID(context.Background(), model.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, models)

	found := false
	for _, m := range models {
		if m.ID == model.ID {
			found = true
			break
		}
	}
	require.True(t, found, "created model not found in list")
}

func TestUpdateModel(t *testing.T) {
	user := createRandomUser(t)
	model := createRandomModel(t, user)

	arg := UpdateModelParams{
		ID:          model.ID,
		Name:        "Updated Name",
		Description: pgtype.Text{String: "Updated description", Valid: true},
		FilePath:    "/tmp/updated_model.bin",
	}

	updatedModel, err := testQueries.UpdateModel(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedModel)
	require.Equal(t, arg.Name, updatedModel.Name)
	require.Equal(t, arg.Description.String, updatedModel.Description.String)
	require.Equal(t, arg.FilePath, updatedModel.FilePath)
}

func TestDeleteModel(t *testing.T) {
	user := createRandomUser(t)
	model := createRandomModel(t, user)

	err := testQueries.DeleteModel(context.Background(), model.ID)
	require.NoError(t, err)

	_, err = testQueries.GetModelByID(context.Background(), model.ID)
	require.Error(t, err)
}
