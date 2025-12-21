package main

import (
	"card-api/card"
	"card-api/firebase"
	"card-api/list"
	"card-api/user"
	"log"      //logを出力
	"net/http" //HTTPサーバやクライアントの機能を使うため
)

func main() {
	firebase.InitFirebase()                //Firebase初期化
	defer firebase.FirestoreClient.Close() //終了時にFirebaseを終わる予約

	http.HandleFunc("/cards", card.CardsHandler) //cardsHandlerはrとwの処理を切り替える
	http.HandleFunc("/lists", list.ListsHandler)
	http.HandleFunc("/users", user.UsersHandler)
	log.Println("Server running on http://localhost:8080") //logはログとして出力を残す
	log.Fatal(http.ListenAndServe(":8080", nil))
	//listenandserveでサーバー起動。終了を待ち続けて終わったらlog.Fatalがerrorを読み取りプログラムを強制終了
}
