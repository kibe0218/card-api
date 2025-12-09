package list

import (
	"net/http" //HTTPサーバやクライアントの機能を使うため
	"time"     //現在時刻を取得
)

type List struct {
	Name      string    `firestore:"listname" json:"listname"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}

func ListsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetLists(w, r)
	case http.MethodPost:
		AddList(w, r)
	default:
		http.Error(w, "許可されていないメソッドっピ", http.StatusMethodNotAllowed)
	}
}
