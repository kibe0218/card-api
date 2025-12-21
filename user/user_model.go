package user

import(
	"time"
	"net/http" //HTTPサーバやクライアントの機能を使うため
)

type User struct {
	ID        string    `firestore:"userid" json:"userid"`
	Name      string    `firestore:"name" json:"name"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUsers(w, r)
		return
	case http.MethodPost:
		AddUser(w, r)
		return
	default:
		http.Error(w, "許可されていないメソッドっピ", http.StatusMethodNotAllowed)
	}
}
