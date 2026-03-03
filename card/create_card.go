package card

import (
	"card-api/firebase"
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"log"
	"net/http" //HTTPサーバやクライアントの機能を使うため
	"time"     //現在時刻を取得

	"cloud.google.com/go/firestore"
)

func AddCard(w http.ResponseWriter, r *http.Request) {
	log.Println("🟡 addcard入る")
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
	log.Println("🟡 指定したlistId:", listID)

	var newCard Card
	if err := json.NewDecoder(r.Body).Decode(&newCard); err != nil {
		http.Error(w, "JSONの形式が正しくないっピ", http.StatusBadRequest)
		return
	}
	err := firebase.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		listRef := firebase.FirestoreClient.
			Collection("users").
			Doc(userID).
			Collection("lists").
			Doc(listID)
		listDoc, err := tx.Get(listRef)
		if err != nil {
			return err
		}
		var listData struct {
			CardCount int `firestore:"cardCount"`
		}
		if err := listDoc.DataTo(&listData); err != nil {
			return err
		}
		nextOrder := listData.CardCount
		newDocRef := listRef.Collection("cards").NewDoc()
		newCard.ID = newDocRef.ID
		newCard.ListID = listID
		newCard.Order = nextOrder
		newCard.CreatedAt = time.Now()
		if err := tx.Update(listRef, []firestore.Update{
			{Path: "cardCount", Value: nextOrder + 1},
		}); err != nil {
			return err
		}
		if err := tx.Set(newDocRef, newCard); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		http.Error(w, "Firestoreへの追加失敗っピ", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCard)
}
