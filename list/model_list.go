package list

import (
	"net/http" //HTTPサーバやクライアントの機能を使うため
	"time"     //現在時刻を取得
)

type List struct {
	ID        string    `firestore:"id" json:"id"`
	Title     string    `firestore:"title" json:"title"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}

func ListsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetLists(w, r)
		return
	case http.MethodPost:
		AddList(w, r)
		return
	case http.MethodDelete:
		DeleteList(w, r)
		return
	default:
		http.Error(w, "許可されていないメソッドっピ", http.StatusMethodNotAllowed)
	}
}
