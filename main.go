package main

import (
	"context"       //処理のキャンセル・タイムアウトを司る
	"encoding/json" //Encode/Decodeのため
	"log"           //logを出力
	"net/http"      //HTTPサーバやクライアントの機能を使うため
	"time"          //現在時刻を取得

	"cloud.google.com/go/firestore"   //firestoreにアクセスするための公式ライブラリ
	firebase "firebase.google.com/go" //firebase全体を使うためのパッケージ
	"google.golang.org/api/option"    //認証キーなどの設定を渡すときに使う
)

type Card struct {
	ID        int       `firestore:"ID" json:"ID"`
	EN        string    `firestore:"en" json:"en"`
	JP        string    `firestore:"jp" json:"jp"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}

var firestoreClient *firestore.Client

func initFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	//別ファイルに保存したキーを読み込む,:=は宣言して代入
	config := &firebase.Config{ProjectID: "memorize-db-23637"}
	//firebaseの設定を作る,＆はポインタ（アドレス）を作る演算子

	app, err := firebase.NewApp(context.Background(), config, opt)
	//firebaseを初期化,appにその情報を代入,さっき代入した値を送る
	//context.Background()は処理の文脈（キャンセルやタイムアウト）を渡すための仕組み
	if err != nil {
		log.Fatalf("Firebase初期化失敗: %v", err) //fatalfはログを出してすぐ終了,%vはエラー内容が入る
		//%v は 「どんな型でも、それなりにいい感じで文字列化する」 汎用の指定子
	}
	client, err := app.Firestore(context.Background())
	//clientはfirestoreへの接続情報、errはエラーが起きた時の情報
	if err != nil {
		log.Fatalf("Firestore接続失敗: %v", err)
	}
	firestoreClient = client
	//上で定義したfirestoreCliantにcliant情報を代入（他の関数でも使えるようにグローバル関数に代入）
}

func cardsHandler(w http.ResponseWriter, r *http.Request) { //rは受け取るものwは返すもの
	switch r.Method {
	case http.MethodGet: //リクエストのメソッドがGETなら・・
		getCards(w, r)
		return //この関数の処理をここで終わらせる
	case http.MethodPost: //リクエストのメソッドがPOSTなら・・
		addCard(w, r)
		return
	default:
		http.Error(w, "許可されていないメソッドっピ", http.StatusMethodNotAllowed)
	}
}

func getCards(w http.ResponseWriter, r *http.Request) {
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

	iter := firestoreClient.Collection("users").
		Doc(userID).
		Collection("lists").
		Doc(listID).
		Collection("cards").
		Documents(ctx)
	//usersコレクションのuserIDに対応するcardsというドキュメントを指定
	//Firestoreのドキュメントを順番に読み取るイテレーター
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
		cards = append(cards, c) //cardsスライスにcを追加したものをcardsに代入
	}

	if len(cards) == 0 {
		http.Error(w, "カードが見つからないっピ", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json") //httpレスポンスのヘッダーに返すデータの種類を教えている
	json.NewEncoder(w).Encode(cards)                   //jsonデータを書き込んで送信
	//EncodeはGoの構造体をJSONに変換して返す
}

func addCard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userIdを指定してね", http.StatusBadRequest)
		return
	}

	var newCard Card
	if err := json.NewDecoder(r.Body).Decode(&newCard); err != nil { //r.Bodyはクライアントがhttpリクエストの本文に送ってきたデータ
		//Decodeは受け取ったJSONデータをGoの構造体に変換する処理
		http.Error(w, "JSONの形式が正しくないっピ", http.StatusBadRequest)
		return
	}
	newCard.CreatedAt = time.Now()

	_, _, err := firestoreClient.Collection("users").Doc(userID).Collection("cards").Add(ctx, newCard)
	//追加されたドキュメントの参照情報はいらないから無視してる
	if err != nil {
		http.Error(w, "Firestoreへの追加失敗っピ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "カード追加完了っピ"})
	//キーも値もstringの辞書型を作る
}

func main() {
	initFirebase()                //Firebase初期化
	defer firestoreClient.Close() //終了時にFirebaseを終わる予約

	http.HandleFunc("/cards", cardsHandler)                //cardsHandlerはrとwの処理を切り替える
	log.Println("Server running on http://localhost:8080") //logはログとして出力を残す
	log.Fatal(http.ListenAndServe(":8080", nil))
	//listenandserveでサーバー起動。終了を待ち続けて終わったらlog.Fatalがerrorを読み取りプログラムを強制終了
}
