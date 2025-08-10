package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomDataset(t *testing.T) Dataset {
	user := createRandomUser(t)
	userID := user.ID

	arg := CreateDatasetParams{
		UserID:      pgtype.Int4{Int32: userID, Valid: true},
		Name:        "Test Dataset",
		Description: pgtype.Text{String: "A test dataset description", Valid: true},
		FilePath:    "/tmp/testfile.csv",
	}

	dataset, err := testQueries.CreateDataset(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dataset)

	require.Equal(t, arg.UserID.Int32, dataset.UserID.Int32)
	require.Equal(t, arg.Name, dataset.Name)
	require.Equal(t, arg.Description.String, dataset.Description.String)
	require.Equal(t, arg.FilePath, dataset.FilePath)

	require.NotZero(t, dataset.ID)
	require.NotZero(t, dataset.UploadedAt)

	return dataset
}

func TestCreateDataset(t *testing.T) {
	createRandomDataset(t)
}

func TestGetDatasetByID(t *testing.T) {
	dataset1 := createRandomDataset(t)

	dataset2, err := testQueries.GetDatasetByID(context.Background(), dataset1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dataset2)

	require.Equal(t, dataset1.ID, dataset2.ID)
	require.Equal(t, dataset1.Name, dataset2.Name)
	require.Equal(t, dataset1.Description.String, dataset2.Description.String)
	require.Equal(t, dataset1.FilePath, dataset2.FilePath)
	require.Equal(t, dataset1.UserID.Int32, dataset2.UserID.Int32)
}

func TestGetDatasetsByUserID(t *testing.T) {
	dataset := createRandomDataset(t)

	datasets, err := testQueries.GetDatasetsByUserID(context.Background(), dataset.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, datasets)

	found := false
	for _, d := range datasets {
		if d.ID == dataset.ID {
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestUpdateDataset(t *testing.T) {
	dataset1 := createRandomDataset(t)

	arg := UpdateDatasetParams{
		ID:          dataset1.ID,
		Name:        "Updated Dataset Name",
		Description: pgtype.Text{String: "Updated description", Valid: true},
		FilePath:    "/tmp/updatedfile.csv",
	}

	dataset2, err := testQueries.UpdateDataset(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dataset2)

	require.Equal(t, arg.Name, dataset2.Name)
	require.Equal(t, arg.Description.String, dataset2.Description.String)
	require.Equal(t, arg.FilePath, dataset2.FilePath)
	require.Equal(t, dataset1.UserID.Int32, dataset2.UserID.Int32)
}
