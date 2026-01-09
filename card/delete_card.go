package card

import (
	"card-api/firebase"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	log.Println("🟡DeleteCardに入る")

	ctx := context.Background()

	userID := r.URL.Query().Get("userId")
	listID := r.URL.Query().Get("listId")
	cardID := r.URL.Query().Get("cardId")

	if userID == "" || listID == "" || cardID == "" {
		http.Error(w, "userId / listId/ cardId が必要っピ", http.StatusBadRequest)
		return
	}

	_, err := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Collection("cards").
		Doc(cardID).
		Delete(ctx)

	if err != nil {
		http.Error(w, "削除に失敗したっピ", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "カード削除完了っピ",
	})

}
