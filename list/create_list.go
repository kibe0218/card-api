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

	if newList.Title == "" {
		newList.Title = "新しいリスト"
	}
	newList.CreatedAt = time.Now()
	newList.CardCount = 0

	// Firestoreのlistsコレクションに新しいドキュメントを追加
	newDocRef := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		NewDoc() // 自分で ID を作る
	newList.ID = newDocRef.ID
	_, err := newDocRef.Set(ctx, newList)
	if err != nil {
		http.Error(w, "リスト追加失敗っピ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "リスト追加完了っピ",
		"id":      newDocRef.ID,
	})
}
