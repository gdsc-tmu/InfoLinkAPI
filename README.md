# InfoLinkAPI
都立大のさまざまな情報をAPIで提供するプロジェクトです。

開発手順は等は[こちら](https://www.notion.so/928680d3ca9f4679ac28d4aae5c53dcf?pvs=4)を見てください

notion: https://www.notion.so/InfoLinkAPI-9780acffd1fd45ca8a7ca3ff0cb3ef82

このプロジェクトでは、データベース上にある大学のシラバス情報や授業教室情報にアクセスできるAPIを開発し、都立大生に対して提供します。
また、このプロジェクトでは開発初心者を広く受け入れたいと考えています。そのための開発手法を取り入れ、GDSCの開発者をサポートしていきます。

### 得られるもの

### 開発参加者

開発参加者は、この開発に関わることで一般的なWebバックエンド開発の経験を積むことが出来ます。また、REST APIを開発することでAPIを利用するサービスの開発もやりやすくなるでしょう。

### 一般に

大学のシラバス情報をAPI化し提供することで多くの便利なアプリケーションの開発に役立ちます。
特に、以下の様なAPI主体のアプリケーションはこのプロジェクト内でともに開発します。

### 展望

APIの開発にとどまらず、APIを使用したアプリケーションの開発も進めていきます。(現時点で予定しているものを以下に示しました。)
また、REST APIで開発したAPIのGraphQL化も視野に入れています。

### シラバス検索システム

CampusSquareでのシラバス検索は手順が多く非常に煩雑で使いづらいものとなっています。これを授業名だけ、教員名だけ、学部名だけなど様々な検索クエリや絞り込み機能を有した優れた検索システムを開発します。

### 教室検索システム

一時的に石池上で公開していた教室検索システムは様々な事情があり封鎖しました。PDFでは探しづらい教室情報が検索できるシステムは多くの都立大生に評価されており、これをGDSCの責任の元で開発し公開します。

### 都立大生向け時間割アプリ

このアプリケーションに関してはAPIを叩くだけでは実現できないので、別プロジェクトでの開発になる可能性があります。

## 開発手法

チームでの開発はタスクの振り分けが難しい問題となる。特にGDSCのような緩いサークルだと持続的な開発が難しくなる。そこで、このプロジェクトでは**セルフアサイン方式**で開発を進める。

### セルフアサイン方式とは

セルフアサイン方式では、チームメンバーに自分でタスクを選択し、アサイン(担当)する自由を与えます。タスクはリポジトリのIssueベースで管理します。ただし、Issueを立てるのは基本プロジェクトリーダーです。選択するタスクがスキルセットに合っていること、チームの目標に対してバランスが取れていることを確認し、担当してもらいます。
**アサインするかどうかは完全に自由**でノルマや期限などは基本的には儲けません。なので、1つだけアサインして終わりでも構いません。もちろん継続的にプロジェクトに参加してくれるプロジェクトメンバーも募集中です。
自身がない方や初心者にはベストエフォートで開発メンバーでサポートしていきます。また、仕様技術に関する定期的な勉強会・ドキュメント読み会等を実施し、チーム全体として学習する機会を儲けます。

## 使用ツール

- GitHub : ソース管理・テスト・デプロイ・タスク割り振りに使用
- Slack : コミュニケーションツールとして使用

## 技術スタック

### プログラミング言語

- Go言語:REST APIの開発に使用。
- Python:スクレイピングに使用。

### フレームワーク

- Gin : API開発フレームワーク
- GORM : ORM

### データベース

- MySQL

### コンテナ化

- Docker : 開発とテストの環境一貫性のために使用。本番環境での使用はリソースと相談して検討。

### CI/CD

- GitHub Actions : 自動化されたテストとデプロイメントパイプラインを実装
- デプロイ方法はインフラによって検討

### テスト

- ユニットテスト : GoのtestingでAPIをテスト

### スクレイピング

- スクレイピング先サイトの構造の定期的なチェックを含むスクレイピングを実行。

### 認証

- APIキー認証

### ログ・エラー管理

- Slackでセントリーbotのようなものを実装
- logを閲覧する管理画面を実装

## 開発に参加する
プログラミングはできるけど開発はしたことないようなweb開発初心者向けのチュートリアル的要素を含んでいるWebバックエンドのチーム開発のプロジェクトを立ち上げました。Web開発してみたいといった理由でサークルに入ったような人には是非参加してほしいプロジェクトです！
このプロジェクトでは常に開発参加者を募集しています。少しでも興味のある人は気軽に連絡してください！[あゆ](https://twitter.com/aya172957)にDMを送るか、GDSCのSlackのInfoLinkAPI開発チームのチャンネルに参加してきてください！

