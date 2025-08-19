package db

import (
	"context"
	"testing"

	"github.com/faezefz/SFP_website/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomPrediction(t *testing.T, userID, datasetID, modelID, projectID int32) Prediction {
	arg := CreatePredictionParams{
		UserID:         pgtype.Int4{Int32: userID, Valid: true},
		DatasetID:      pgtype.Int4{Int32: datasetID, Valid: true},
		ModelID:        pgtype.Int4{Int32: modelID, Valid: true},
		ProjectID:      pgtype.Int4{Int32: projectID, Valid: true},
		ResultFilePath: pgtype.Text{String: util.RandomString(10) + ".csv", Valid: true},
	}

	prediction, err := testQueries.CreatePrediction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, prediction)

	return prediction
}

func TestPredictions(t *testing.T) {
	user := createRandomUser(t)
	project := createRandomProject(t, user.ID)
	dataset := createRandomDataset(t)
	model := createRandomModel(t, user)

	pred := createRandomPrediction(t, user.ID, dataset.ID, model.ID, project.ID)
	got, err := testQueries.GetPredictionByID(context.Background(), pred.ID)
	require.NoError(t, err)
	require.Equal(t, pred.ID, got.ID)

	preds, err := testQueries.GetPredictionsByUserID(context.Background(), pgtype.Int4{Int32: user.ID, Valid: true})
	require.NoError(t, err)
	require.True(t, len(preds) >= 1)

	newPath := pgtype.Text{String: util.RandomString(8) + ".csv", Valid: true}
	updated, err := testQueries.UpdatePrediction(context.Background(), UpdatePredictionParams{
		ID:             pred.ID,
		ResultFilePath: newPath,
	})
	require.NoError(t, err)
	require.Equal(t, newPath.String, updated.ResultFilePath.String)

	err = testQueries.DeletePrediction(context.Background(), pred.ID)
	require.NoError(t, err)

	_, err = testQueries.GetPredictionByID(context.Background(), pred.ID)
	require.Error(t, err)
}
