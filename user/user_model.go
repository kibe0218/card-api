package user

// import(
// 	"time"
// 	"net/http"      //HTTPサーバやクライアントの機能を使うため

// )

// type User struct {
// 	Name      string    `firestore:"listname" json:"listname"`
// 	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
// }

// func ListsHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		GetUsers(w, r)
// 	case http.MethodPost:
// 		AddUser(w, r)
// 	default:
// 		http.Error(w, "許可されていないメソッドっピ", http.StatusMethodNotAllowed)
// 	}
// }
