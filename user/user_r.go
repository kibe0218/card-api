package user

import (
	"card-api/firebase"
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"net/http"      //HTTPサーバやクライアントの機能を使うため
)

type UserResponse struct {
	ID string `json:"id"`
	User
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	iter := firebase.FirestoreClient.Collection("users").
		Documents(ctx)
	defer iter.Stop()

	users := []UserResponse{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var u User
		if err := doc.DataTo(&u); err != nil {
			continue
		}
		users = append(users, UserResponse{
			ID:   doc.Ref.ID,
			User: u,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
