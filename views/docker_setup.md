# 開発手順

## 手順(初期設定)

Dockerをインストールして起動しておきます。(Docker Desktopが良いと思います)

Dockerコンテナを起動します。
```
docker-compose up --build -d
```

データベースのdumpファイルを`schema/dump.sql`に配置してコンテナ上にコピーします。
```
docker cp ./schema/dump.sql mysql-container:/tmp/dump.sql
```

mysqlのコンテナに入ります
```
docker exec -it mysql-container bash
```

dumpされたデータをデータベースに流し込みます。
```
mysql -u root -p demo < /tmp/dump.sql
```

コンテナを再起動します。
```
docker-compose up -d
```
[Swagger UI](http://localhost:8080/swagger/index.html)にアクセスして正常に動作しているか確認します。

### 手順（開発時）
コードを書いたら、コンテナを再起動します。(ホットリロードはされません)
```
docker-compose up --build -d
```

# その他

## [Swagger UI](http://localhost:8080/swagger/index.html)

## テーブル名は必ず複数形
GORMのドキュメントをよく読むこと