package list

import (
	"card-api/firebase"
	"context"
	"encoding/json"
	"net/http"
)

func DeleteList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	userID := r.URL.Query().Get("userId")
	listID := r.URL.Query().Get("listId")

	// ❌ ここで http.Error はやめる
	if userID == "" || listID == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]string{})
		return
	}
	_, err := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Delete(ctx)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode([]string{})
		return
	}

	// 成功レスポンスも JSON に
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]string{})
}
