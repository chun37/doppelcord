# Doppelcord - Discord Bot

Go言語で実装されたシンプルなDiscordボットです。

## 機能

- `/register` スラッシュコマンドでユーザーを登録
- 登録済みユーザーからのメッセージには `[登録済]` プレフィックスを表示
- PostgreSQLによるデータ永続化
- 環境変数を使用した安全なトークン管理
- グレースフルシャットダウン対応

## 前提条件

- Go 1.16以上
- PostgreSQL
- Discord Bot アカウントとトークン

## Discord Botの作成方法

1. [Discord Developer Portal](https://discord.com/developers/applications)にアクセス
2. 「New Application」をクリックして新しいアプリケーションを作成
3. 左側メニューから「Bot」を選択し、「Add Bot」をクリック
4. 「MESSAGE CONTENT INTENT」を有効化（重要！）
5. 「TOKEN」セクションから「Copy」をクリックしてトークンをコピー
6. 左側メニューから「OAuth2」→「URL Generator」を選択
   - SCOPES: `bot`を選択
   - BOT PERMISSIONS: `Send Messages`, `Read Message History`, `View Channels`を選択
7. 生成されたURLからボットをサーバーに招待

## セットアップ

### 1. リポジトリをクローンまたはダウンロード

### 2. 依存関係をインストール
```bash
go mod download
```

### 3. PostgreSQLのセットアップ
```bash
# データベースとユーザーを作成
createdb doppelcord
createuser doppelcord -P
```

### 4. golang-migrateのインストールとマイグレーション実行
```bash
# golang-migrateをインストール
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# マイグレーション実行（.envから接続情報を読み込む）
./scripts/migrate.sh up

# その他のコマンド
./scripts/migrate.sh version  # 現在のバージョン確認
./scripts/migrate.sh down     # ロールバック
./scripts/migrate.sh force 1  # バージョン強制設定
```

### 5. `.env`ファイルを設定
```bash
# Discord設定
DISCORD_BOT_TOKEN=あなたのボットトークン
GUILD_ID=スラッシュコマンドを登録するサーバーID

# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=doppelcord
DB_PASSWORD=your_password
DB_NAME=doppelcord
DB_SSLMODE=disable
```

## 実行方法

```bash
go run main.go
```

または、ビルドしてから実行

```bash
go build
./doppelcord
```

## 停止方法

`Ctrl+C`を押すとグレースフルにシャットダウンします。

## プロジェクト構造

```
doppelcord/
├── main.go                          # エントリーポイント
├── go.mod                           # Go modules設定
├── go.sum                           # 依存関係チェックサム
├── .env                             # 環境変数（gitignore対象）
├── .env.example                     # 環境変数テンプレート
├── .gitignore                       # Git除外設定
├── README.md                        # このファイル
├── scripts/
│   └── migrate.sh                   # マイグレーション実行スクリプト
├── internal/
│   ├── domain/
│   │   └── user.go                  # ドメインモデル
│   ├── repository/
│   │   ├── user_repository.go       # Repositoryインターフェース
│   │   └── postgres/
│   │       └── user_repository.go   # PostgreSQL実装
│   ├── handler/
│   │   ├── message_handler.go       # メッセージハンドラー
│   │   └── interaction_handler.go   # インタラクションハンドラー
│   └── database/
│       └── postgres.go              # DB接続管理
└── migrations/
    ├── 000001_create_users_table.up.sql
    └── 000001_create_users_table.down.sql
```

## 注意事項

- `.env`ファイルにはボットトークンやDBパスワードが含まれるため、絶対に公開リポジトリにコミットしないでください
- Discord Developer Portalで「MESSAGE CONTENT INTENT」を必ず有効化してください
- ボットは自分自身のメッセージには反応しません（無限ループ防止）

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
