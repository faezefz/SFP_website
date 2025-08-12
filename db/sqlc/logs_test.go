package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomLog(t *testing.T) Log {
	user := createRandomUser(t)
	project := createRandomProject(t, user.ID)

	arg := CreateLogParams{
		UserID:    pgtype.Int4{Int32: user.ID, Valid: true},
		ProjectID: pgtype.Int4{Int32: project.ID, Valid: true},
		Action:    pgtype.Text{String: "Created", Valid: true},
		Details:   pgtype.Text{String: "Created a new project", Valid: true},
	}

	logEntry, err := testQueries.CreateLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, logEntry)

	require.Equal(t, arg.UserID.Int32, logEntry.UserID.Int32)
	require.Equal(t, arg.ProjectID.Int32, logEntry.ProjectID.Int32)
	require.Equal(t, arg.Action.String, logEntry.Action.String)
	require.Equal(t, arg.Details.String, logEntry.Details.String)

	return logEntry
}

func TestCreateLog(t *testing.T) {
	createRandomLog(t)
}

func TestGetLogByID(t *testing.T) {
	logEntry := createRandomLog(t)

	fetched, err := testQueries.GetLogByID(context.Background(), logEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetched)
	require.Equal(t, logEntry.ID, fetched.ID)
}

func TestGetLogsByProjectOrUser(t *testing.T) {
	logEntry := createRandomLog(t)

	params := GetLogsByProjectOrUserParams{
		ProjectID: logEntry.ProjectID,
		UserID:    logEntry.UserID,
	}

	logs, err := testQueries.GetLogsByProjectOrUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, logs)

	found := false
	for _, l := range logs {
		if l.ID == logEntry.ID {
			found = true
			break
		}
	}
	require.True(t, found, "created log not found in list")
}

func TestDeleteLog(t *testing.T) {
	logEntry := createRandomLog(t)

	err := testQueries.DeleteLog(context.Background(), logEntry.ID)
	require.NoError(t, err)

	_, err = testQueries.GetLogByID(context.Background(), logEntry.ID)
	require.Error(t, err)
}
