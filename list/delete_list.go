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

	if userID == "" || listID == "" {
		http.Error(w, "userId / listIdが必要っピ", http.StatusBadRequest)
		return
	}

	_, err := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Delete(ctx)

	if err != nil {
		http.Error(w, "削除に失敗したっピ", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "リスト削除完了っピ",
	})

}
