# 🐳 軽量なGo環境
FROM golang:1.25

# 📁 作業ディレクトリ
WORKDIR /app

# 📦 依存ファイルだけ先にコピー（キャッシュ効かせるため）
COPY go.mod go.sum ./
RUN go mod download

# 📂 ソースコピー
COPY . .

# ⚙️ ビルド
RUN go build -o app

# 🚀 実行
CMD ["./app"]