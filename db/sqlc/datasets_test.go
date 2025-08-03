package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomDataset(t *testing.T, userID int32) Dataset {
	userIDPg := pgtype.Int4{Int32: userID, Valid: true}

	descriptionPg := pgtype.Text{String: "A description for the dataset", Valid: true}

	arg := CreateDatasetParams{
		UserID:      userIDPg,
		Name:        "Sample Dataset",
		Description: descriptionPg,
		FilePath:    "/path/to/file",
	}

	dataset, err := testQueries.CreateDataset(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dataset)
	require.Equal(t, arg.Name, dataset.Name)
	require.Equal(t, arg.Description.String, dataset.Description.String)

	return dataset
}

func TestCreateDataset(t *testing.T) {
	user := createRandomUser(t)
	createRandomDataset(t, user.ID)
}

func TestGetDataset(t *testing.T) {
	user := createRandomUser(t)
	dataset := createRandomDataset(t, user.ID)

	retrievedDataset, err := testQueries.GetDataset(context.Background(), dataset.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedDataset)

	require.Equal(t, dataset.ID, retrievedDataset.ID)
	require.Equal(t, dataset.Name, retrievedDataset.Name)
	require.Equal(t, dataset.Description.String, retrievedDataset.Description.String)
	require.Equal(t, dataset.FilePath, retrievedDataset.FilePath)
}

func TestUpdateDataset(t *testing.T) {
	user := createRandomUser(t)
	dataset := createRandomDataset(t, user.ID)

	arg := UpdateDatasetParams{
		ID:          dataset.ID,
		Name:        "Updated Dataset Name",
		Description: pgtype.Text{String: "Updated description", Valid: true},
		FilePath:    "/new/path/to/file",
	}

	updatedDataset, err := testQueries.UpdateDataset(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedDataset)

	require.Equal(t, arg.Name, updatedDataset.Name)
	require.Equal(t, arg.Description.String, updatedDataset.Description.String)
	require.Equal(t, arg.FilePath, updatedDataset.FilePath)
}

func TestDeleteDataset(t *testing.T) {
	user := createRandomUser(t)
	dataset := createRandomDataset(t, user.ID)

	err := testQueries.DeleteDataset(context.Background(), dataset.ID)
	require.NoError(t, err)

	retrievedDataset, err := testQueries.GetDataset(context.Background(), dataset.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, retrievedDataset)
}

func TestListDatasets(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomDataset(t, user.ID)
	}

	arg := ListDatasetsParams{
		UserID: pgtype.Int4{Int32: user.ID, Valid: true},
		Limit:  5,
		Offset: 5,
	}

	datasets, err := testQueries.ListDatasets(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, datasets, 5)

	for _, dataset := range datasets {
		require.NotEmpty(t, dataset)
	}
}
