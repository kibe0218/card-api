package card

import (
	"card-api/firebase"
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"log"
	"net/http" //HTTPサーバやクライアントの機能を使うため
	"time"     //現在時刻を取得
)

func AddCard(w http.ResponseWriter, r *http.Request) {
	log.Println("🟡 addcard入る")
	ctx := context.Background()

	listID := r.URL.Query().Get("listId")
	if listID == "" {
		http.Error(w, "listIdを指定してね", http.StatusBadRequest)
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userIdを指定してね", http.StatusBadRequest)
		return
	}
	log.Println("🟡 指定したlistId:", listID)

	var newCard Card
	if err := json.NewDecoder(r.Body).Decode(&newCard); err != nil { //r.Bodyはクライアントがhttpリクエストの本文に送ってきたデータ
		//Decodeは受け取ったJSONデータをGoの構造体に変換する処理
		http.Error(w, "JSONの形式が正しくないっピ", http.StatusBadRequest)
		return
	}
	newCard.CreatedAt = time.Now()

	// Firestoreのcardsコレクションに新しいドキュメントを追加
	newDocRef := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Collection("cards").
		NewDoc() // 自分で ID を作る
	newCard.ID = newDocRef.ID
	newCard.ListID = listID
	_, err := newDocRef.Set(ctx, newCard)
	if err != nil {
		http.Error(w, "Firestoreへの追加失敗っピ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "カード追加完了っピ",
		"id":      newCard.ID,
		"listid":  newCard.ListID,
	})
	//キーも値もstringの辞書型を作る
}
