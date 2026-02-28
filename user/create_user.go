package user

import (
	"bytes"
	"card-api/firebase"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("🟡 AddUser 入ったっピ")
	ctx := context.Background()

	// --- デバッグ用: raw body を読み込む ---
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("🟡 リクエストボディ読み取り失敗:", err)
		http.Error(w, "ボディ読み取り失敗っピ", http.StatusBadRequest)
		return
	}
	log.Println("🟡 raw body:", string(body)) // ←ここでbodyにIDがあるか確認

	// Bodyは一回読むと空になる
	r.Body = io.NopCloser(bytes.NewBuffer(body)) //@gmail.com

	// --- JSON デコード ---
	var newUser User
	log.Printf("%T %+v\n", newUser, newUser)
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println("🟡 JSON デコード失敗っピ:", err)
		http.Error(w, "JSONの形式が正しくありません", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Println("🟡 デコード後 newUser struct:", newUser) // ←ここでIDが正しく入っているか確認

	// --- デバッグ追加: ID が空か確認 ---
	if newUser.ID == "" {
		log.Println("🟡 デコード後 ID が空っピ")
		// ここで rawMap を作って確認することも可能
		var checkMap map[string]interface{} //jsonを受け取るための箱
		json.Unmarshal(body, &checkMap)
		log.Println("🟡 raw body を map にした場合:", checkMap)
		log.Println("🟡 map 内の id:", checkMap["id"])
		http.Error(w, "IDが空です（id フィールドを送ってください）", http.StatusBadRequest)
		return
	}

	newUser.CreatedAt = time.Now()

	docRef := firebase.FirestoreClient.Collection("users").Doc(newUser.ID)
	if _, err := docRef.Set(ctx, newUser); err != nil {
		log.Println("🟡 Firestore 保存失敗っピ:", err)
		http.Error(w, "ユーザー追加失敗っピ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "ユーザー追加完了っピ",
		"id":      docRef.ID,
	})
	log.Println("🟡 ユーザー追加完了:", newUser.ID)
}
