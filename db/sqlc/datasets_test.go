package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/faezefz/SFP_website/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomDataset(t *testing.T) Dataset {
	user := createRandomUser(t)
	userID := user.ID

	// خواندن محتوای فایل CSV به صورت بایت
	content, err := util.ReadCSV()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}

	arg := CreateDatasetParams{
		UserID:      pgtype.Int4{Int32: userID, Valid: true},
		Name:        "Test Dataset with CSV",
		Description: pgtype.Text{String: "A test dataset description", Valid: true},
		Content:     content, // ذخیره محتوای فایل به صورت بایت
	}

	dataset, err := testQueries.CreateDataset(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dataset)

	require.Equal(t, arg.UserID.Int32, dataset.UserID.Int32)
	require.Equal(t, arg.Name, dataset.Name)
	require.Equal(t, arg.Description.String, dataset.Description.String)
	require.Equal(t, arg.Content, dataset.Content) // مقایسه محتوای بایت

	require.NotZero(t, dataset.ID)
	require.NotZero(t, dataset.UploadedAt)

	return dataset
}

func TestCreateDataset(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	createRandomDataset(t) // تست ایجاد دیتاست با فایل CSV
}

func TestGetDatasetByID(t *testing.T) {
	dataset1 := createRandomDataset(t)

	dataset2, err := testQueries.GetDatasetByID(context.Background(), dataset1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dataset2)

	require.Equal(t, dataset1.ID, dataset2.ID)
	require.Equal(t, dataset1.Name, dataset2.Name)
	require.Equal(t, dataset1.Description.String, dataset2.Description.String)
	require.Equal(t, dataset1.Content, dataset2.Content)
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
	content, err := util.ReadCSV()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}
	arg := UpdateDatasetParams{
		ID:          dataset1.ID,
		Name:        "Updated Dataset Name",
		Description: pgtype.Text{String: "Updated description", Valid: true},
		Content:     content,
	}

	dataset2, err := testQueries.UpdateDataset(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, dataset2)

	require.Equal(t, arg.Name, dataset2.Name)
	require.Equal(t, arg.Description.String, dataset2.Description.String)
	require.Equal(t, arg.Content, dataset2.Content)
	require.Equal(t, dataset1.UserID.Int32, dataset2.UserID.Int32)
}
