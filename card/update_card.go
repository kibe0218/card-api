package card

import (
	"card-api/firebase"
	"context"
	"encoding/json"
	"net/http"
)

func UpdateCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	listID := r.URL.Query().Get("listId")
	if listID == "" {
		http.Error(w, "listIdを指定してね", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userIdを指定してね", http.StatusBadRequest)
		return
	}

	cardID := r.URL.Query().Get("cardId")
	if cardID == "" {
		http.Error(w, "cardIdを指定してね", http.StatusBadRequest)
		return
	}

	var updated Card
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "JSONの形式が正しくないっピ", http.StatusBadRequest)
		return
	}

	_, err := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Collection("cards").
		Doc(cardID).
		Set(ctx, updated)
	if err != nil {
		http.Error(w, "Firestoreの更新に失敗したっピ", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updated)
}
