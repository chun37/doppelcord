# Doppelcord - Discord Bot

Go言語で実装されたシンプルなDiscordボットです。

## 機能

- `/register` スラッシュコマンドでユーザーを登録
- `/test` スラッシュコマンドで、そのユーザーの発言履歴をもとにLLMが「なりきり」メッセージを生成
  - まずそのチャンネルでの発言履歴を取得（最大100件）
  - チャンネルに履歴がなければ全チャンネルの履歴を使用
- 登録済みユーザーからのメッセージには `[登録済]` プレフィックスを表示
- 登録済みユーザーのメッセージをDBに保存（月別パーティショニングで大規模対応）
- PostgreSQLによるデータ永続化
- メモリキャッシュによる高速な登録確認（起動時にDBから読み込み、以降はメモリ参照）
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

# LLM API設定（OpenAI互換API）
LLM_API_URL=https://api.openai.com/v1/chat/completions
LLM_API_KEY=your_api_key
LLM_MODEL=gpt-4o-mini
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
│   │   ├── user.go                  # ユーザードメインモデル
│   │   └── message.go               # メッセージドメインモデル
│   ├── llm/
│   │   ├── types.go                 # LLM API型定義
│   │   └── client.go                # OpenAI互換APIクライアント
│   ├── repository/
│   │   ├── user_repository.go       # UserRepositoryインターフェース
│   │   ├── message_repository.go    # MessageRepositoryインターフェース
│   │   ├── cached/
│   │   │   └── user_repository.go   # キャッシュ付きUserRepository
│   │   └── postgres/
│   │       ├── user_repository.go   # UserRepository PostgreSQL実装
│   │       └── message_repository.go # MessageRepository PostgreSQL実装
│   ├── handler/
│   │   ├── message_handler.go       # メッセージハンドラー
│   │   └── interaction_handler.go   # インタラクションハンドラー
│   └── database/
│       └── postgres.go              # DB接続管理
└── migrations/
    ├── 000001_create_users_table.up.sql
    ├── 000001_create_users_table.down.sql
    ├── 000002_create_messages_table.up.sql
    └── 000002_create_messages_table.down.sql
```

## 注意事項

- `.env`ファイルにはボットトークンやDBパスワードが含まれるため、絶対に公開リポジトリにコミットしないでください
- Discord Developer Portalで「MESSAGE CONTENT INTENT」を必ず有効化してください
- ボットは自分自身のメッセージには反応しません（無限ループ防止）

## デプロイ

### GitHub Secrets の設定

リポジトリの Settings > Secrets and variables > Actions で以下のSecretsを設定:

- `DISCORD_BOT_TOKEN` - Discord Bot トークン
- `GUILD_ID` - スラッシュコマンドを登録するサーバーID
- `DB_HOST` - PostgreSQLホスト
- `DB_PORT` - PostgreSQLポート
- `DB_USER` - PostgreSQLユーザー名
- `DB_PASSWORD` - PostgreSQLパスワード
- `DB_NAME` - データベース名
- `DB_SSLMODE` - SSL モード（通常は `disable`）
- `LLM_API_URL` - LLM APIエンドポイントURL
- `LLM_API_KEY` - LLM APIキー
- `LLM_MODEL` - 使用するモデル名

### 自動デプロイ

`master`ブランチにpushすると、Self-hosted Runner経由で自動的にデプロイされます。

- アプリは `/opt/doppelcord/` に配置されます
- systemdサービス `doppelcord` として実行されます

### 手動操作

```bash
# サービスの状態確認
sudo systemctl status doppelcord

# ログの確認
sudo journalctl -u doppelcord -f

# 再起動
sudo systemctl restart doppelcord
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
