package user

import (
    "card-api/firebase"
    "context"
    "encoding/json"
    "net/http"
    "time"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

    var newUser User

    defer r.Body.Close()
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
	http.Error(w, "JSONの形式が正しくないっピ", http.StatusBadRequest)
	return
    }

    newUser.CreatedAt = time.Now()

    newDocRef := firebase.FirestoreClient.Collection("users").NewDoc()
    newUser.ID = newDocRef.ID
    _, err := newDocRef.Set(ctx, newUser)
    if err != nil {
	http.Error(w, "ユーザー追加失敗っピ", http.StatusInternalServerError)
	return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
	"message": "ユーザー追加完了っピ",
	"id":      newDocRef.ID,
    })
}
