# dockerのセットアップ

## 基本操作

### コンテナ起動
```
docker-compose up --build -d
```

### データベースに入る
```
docker-conpose up -d
docker exec -it mysql-container bash
mysql -u root -p demo
```

## コンテナに入る
```
docker exec -it コンテナ名 sh
```

## 手順(初期設定)

### dumpデータの流し込み(初回のみ)
```
docker-conpose up -d
docker cp ./schema/dump.sql mysql-container:/tmp/dump.sql
docker exec -it mysql-container bash
mysql -u root -p demo < /tmp/dump.sql
```

### 手順（開発時）
コードを書いたら
```
swag init
docker-compose up --build -d
```

## [Swagger UI](http://localhost:8080/swagger/index.html)

## テーブル名は必ず複数形
GORMのドキュメントをよく読むこと