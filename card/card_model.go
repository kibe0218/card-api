package card

import (
	"net/http" //HTTPサーバやクライアントの機能を使うため
	"time"
)

type Card struct {
	EN        string    `firestore:"en" json:"en"`
	JP        string    `firestore:"jp" json:"jp"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}

func CardsHandler(w http.ResponseWriter, r *http.Request) { //rは受け取るものwは返すもの
	switch r.Method {
	case http.MethodGet: //リクエストのメソッドがGETなら・・
		GetCards(w, r)
		return //この関数の処理をここで終わらせる
	case http.MethodPost: //リクエストのメソッドがPOSTなら・・
		AddCard(w, r)
		return
	default:
		http.Error(w, "許可されていないメソッドっピ", http.StatusMethodNotAllowed)
	}
}
