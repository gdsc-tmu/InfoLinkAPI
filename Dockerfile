# ビルドステージ
FROM golang:1.21 AS builder

# 作業ディレクトリを設定
WORKDIR /app

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
ENV PATH="/go/bin:${PATH}"

COPY go.mod ./
COPY go.sum ./

# モジュールをダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# テストステージ
FROM golang:1.21 AS tester
WORKDIR /app
COPY --from=builder /go /go
COPY --from=builder /app /app
RUN go test -v ./...

# tidyステージ
FROM golang:1.21 AS tidy
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . .
RUN go mod tidy

# swaggerステージ
FROM golang:1.21 AS swagger
# 作業ディレクトリを設定
WORKDIR /app
# ソースコードをコピー
COPY . .
# Swaggerドキュメントを生成
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# 実行ステージ
FROM alpine:latest

# セキュリティ関連の更新を適用
RUN apk --no-cache add ca-certificates

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドステージから実行可能ファイルをコピー
COPY --from=builder /app/main .

# ビルドステージからドキュメントをコピー
COPY --from=builder /app/docs /app/docs

# アプリケーションを実行
CMD ["./main"]
