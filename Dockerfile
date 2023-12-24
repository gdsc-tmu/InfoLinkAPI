# ビルドステージ
FROM golang:1.19 AS builder

# 作業ディレクトリを設定
WORKDIR /app

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# モジュールファイルをコピー
COPY go.mod ./
COPY go.sum ./

# モジュールをダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

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
