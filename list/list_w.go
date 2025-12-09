package list

import (
	"card-api/firebase"
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"net/http"      //HTTPサーバやクライアントの機能を使うため
	"time"          //現在時刻を取得
)

func AddList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userIdを指定してね", http.StatusBadRequest)
		return
	}

	var newList List
	if err := json.NewDecoder(r.Body).Decode(&newList); err != nil {
		http.Error(w, "JSONの形式が正しくないっピ", http.StatusBadRequest)
		return
	}

	if newList.Name == "" {
		newList.Name = "新しいリスト"
	}
	newList.CreatedAt = time.Now()

	// Firestoreのlistsコレクションに新しいドキュメントを追加
	_, _, err := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Add(ctx, newList)
	if err != nil {
		http.Error(w, "リスト追加失敗っピ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "リスト追加完了っピ"})
}
