package list

import (
	"card-api/firebase"
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"net/http"      //HTTPサーバやクライアントの機能を使うため
)

func GetLists(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userIdを指定してね", http.StatusBadRequest)
		return
	}

	iter := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Documents(ctx)
	defer iter.Stop()

	var lists []List
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var l List
		if err := doc.DataTo(&l); err != nil {
			continue
		}
		lists = append(lists, l)
	}

	if len(lists) == 0 {
		http.Error(w, "リストが見つからないっピ", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lists)
}
