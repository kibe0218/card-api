package card

import (
	"card-api/firebase"
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"net/http"      //HTTPサーバやクライアントの機能を使うため

	"google.golang.org/api/iterator"
)

func GetCards(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userIdを指定してね", http.StatusBadRequest)
		return
	}

	// ① lists を取得
	listIter := firebase.FirestoreClient.
		Collection("users").
		Doc(userID).
		Collection("lists").
		Documents(ctx)
	defer listIter.Stop()

	var cards []Card

	for {
		listDoc, err := listIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listID := listDoc.Ref.ID

		// ② その list の cards を取得
		cardIter := firebase.FirestoreClient.
			Collection("users").
			Doc(userID).
			Collection("lists").
			Doc(listID).
			Collection("cards").
			Documents(ctx)

		for {
			doc, err := cardIter.Next()
			if err != nil {
				break
			}

			var c Card
			if err := doc.DataTo(&c); err != nil {
				continue
			}
			c.ID = doc.Ref.ID
			cards = append(cards, c)
		}

		cardIter.Stop()
	}

	w.Header().Set("Content-Type", "application/json")
	if cards == nil {
		cards = []Card{}
	}
	json.NewEncoder(w).Encode(cards)
}

func GetCardsBy(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background() //処理用のコンテキスト情報を入れる箱を作っている

	listID := r.URL.Query().Get("listId")
	if listID == "" {
		http.Error(w, "listIdを指定してね", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("userId")
	//URLから.Queryでクエリを解析、解析するのはURLの？yserId=以降、それをuserIDに代入
	if userID == "" {
		http.Error(w, "userIdを指定してね", http.StatusBadRequest)
		return
	}

	iter := firebase.FirestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Collection("cards").
		Documents(ctx)
	//usersコレクションのuserIDに対応するcardsというドキュメントを指定
	//Firestoreのドキュメントを順番に読み取るイテレーター
	//ctxでGoの並行処理でキャンセルやタイムアウトを伝える
	defer iter.Stop() //.Stopで上で作ったイテレーターを削除,deferにより関数終了時に自動実行

	var cards []Card
	for { //iterから一件ずつドキュメントを取り出す
		doc, err := iter.Next() //イテレータから次のドキュメントを取り出す,初回は最初のdocを読み取る
		//docには一件分のドキュメント情報が入る,errは取り出せなかった時のエラー情報が入る
		if err != nil {
			break
		}
		var c Card                             //上で作ったCard型のcを宣言
		if err := doc.DataTo(&c); err != nil { //代入と判定を同時にやる書き方
			//cの構造にdocを変化させて代入
			continue //errが何か入っていたら次のループに入り、次のイテレーターにとぶ
		}
		c.ID = doc.Ref.ID
		cards = append(cards, c)
	}

	w.Header().Set("Content-Type", "application/json") //httpレスポンスのヘッダーに返すデータの種類を教えている
	if cards == nil {
		cards = []Card{}
	}
	json.NewEncoder(w).Encode(cards) //jsonデータを書き込んで送信
	//EncodeはGoの構造体をJSONに変換して返す
}
