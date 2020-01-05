package util

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// MySnsData 外部から参照するSNS情報の構造体
type MySnsData struct {
	Facebook  string
	Twitter   string
	Instagram string
	Line      string
}

func getFireBaseClient(ctx context.Context) (*firestore.Client, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")

	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIAL"))

	client, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// SaveUserItem Firestoreにデータを登録する
func SaveUserItem(m MySnsData) (MySnsData, error) {
	ctx := context.Background()

	client, err := getFireBaseClient(ctx)
	if err != nil {
		return m, err
	}

	_, err = client.Collection("users").Doc("1").Set(ctx, map[string]interface{}{
		"facebook":  m.Facebook,
		"twitter":   m.Twitter,
		"instagram": m.Instagram,
		"line":      m.Line,
	})
	if err != nil {
		return m, err
	}

	return m, nil
}

// GetUserItem Firestoreからデータを取得する
func GetUserItem() (MySnsData, error) {
	ctx := context.Background()

	m := MySnsData{}

	client, err := getFireBaseClient(ctx)
	if err != nil {
		return m, err
	}

	doc := client.Collection("users").Doc("1")

	field, err := doc.Get(ctx)
	if err != nil {
		return m, err
	}

	data := field.Data()

	m.Facebook = data["facebook"].(string)
	m.Twitter = data["twitter"].(string)
	m.Instagram = data["instagram"].(string)
	m.Line = data["line"].(string)

	return m, nil
}
