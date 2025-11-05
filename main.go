package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type Card struct {
    ID    int    `json:"id"`
    Word  string `json:"word"`
    Mean  string `json:"mean"`
}

func main() {//f
    http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
        // テスト用のダミーデータ
        cards := []Card{
            {ID: 1, Word: "apple", Mean: "りんご"},
            {ID: 2, Word: "sky", Mean: "空"},
            {ID: 3, Word: "book", Mean: "本"},
        }

        // JSONで返す
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cards)
    })

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}