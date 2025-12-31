package firebase

import (
	"context" //処理のキャンセル・タイムアウトを司る
	"log"     //logを出力

	"cloud.google.com/go/firestore"   //firestoreにアクセスするための公式ライブラリ
	firebase "firebase.google.com/go" //firebase全体を使うためのパッケージ
	"google.golang.org/api/option"    //認証キーなどの設定を渡すときに使う
)

var FirestoreClient *firestore.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
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
	FirestoreClient = client
	//上で定義したfirestoreCliantにcliant情報を代入（他の関数でも使えるようにグローバル関数に代入）
}
