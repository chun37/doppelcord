# Doppelcord - Discord Echo Bot

Go言語で実装されたシンプルなDiscordエコーボットです。ボット以外のすべてのメッセージを同じチャンネルにエコーします。

## 機能

- ボット以外のすべてのメッセージを自動的にエコー
- 環境変数を使用した安全なトークン管理
- グレースフルシャットダウン対応

## 前提条件

- Go 1.16以上
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

1. リポジトリをクローンまたはダウンロード

2. 依存関係をインストール
```bash
go mod download
```

3. `.env`ファイルにBotトークンを設定
```bash
# .envファイルを編集
DISCORD_BOT_TOKEN=あなたのボットトークン
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
├── main.go          # メインプログラム
├── go.mod           # Go modules設定
├── go.sum           # 依存関係チェックサム
├── .env             # 環境変数（gitignore対象）
├── .env.example     # 環境変数テンプレート
├── .gitignore       # Git除外設定
└── README.md        # このファイル
```

## 注意事項

- `.env`ファイルにはボットトークンが含まれるため、絶対に公開リポジトリにコミットしないでください
- Discord Developer Portalで「MESSAGE CONTENT INTENT」を必ず有効化してください
- ボットは自分自身のメッセージには反応しません（無限ループ防止）

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
