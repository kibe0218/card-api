package card

import (
	"card-api/firebase"
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"net/http"      //HTTPサーバやクライアントの機能を使うため
	"time"          //現在時刻を取得
)

func AddCard(w http.ResponseWriter, r *http.Request) {
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

	var newCard Card
	if err := json.NewDecoder(r.Body).Decode(&newCard); err != nil { //r.Bodyはクライアントがhttpリクエストの本文に送ってきたデータ
		//Decodeは受け取ったJSONデータをGoの構造体に変換する処理
		http.Error(w, "JSONの形式が正しくないっピ", http.StatusBadRequest)
		return
	}
	newCard.CreatedAt = time.Now()

	_, _, err := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Collection("cards").
		Add(ctx, newCard) //追加されたドキュメントの参照情報はいらないから無視してる
	if err != nil {
		http.Error(w, "Firestoreへの追加失敗っピ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "カード追加完了っピ"})
	//キーも値もstringの辞書型を作る
}
