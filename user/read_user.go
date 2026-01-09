package user

import (
	"card-api/firebase"
	"context"       //
	"encoding/json" //Encode/Decodeのため
	"net/http"      //HTTPサーバやクライアントの機能を使うため
)

type UserResponse struct {
	ID string `json:"id"`
	User
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userId := r.URL.Query().Get("userId")

	if userId != "" {
		doc, err := firebase.FirestoreClient.Collection("users").Doc(userId).Get(ctx)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		var u User
		if err := doc.DataTo(&u); err != nil {
			http.Error(w, "Error decoding user data", http.StatusInternalServerError)
			return
		}
		user := UserResponse{
			ID:   doc.Ref.ID,
			User: u,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}
}
