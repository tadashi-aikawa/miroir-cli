# miroir-cli

![GitHub Release](https://img.shields.io/github/release/tadashi-aikawa/miroir-cli.svg)

`miroir-cli` は Miroir のレポートサマリや詳細レポートを AWS から取得し、不要なサマリを整理するための CLI です。

## できること

- DynamoDB からサマリ一覧を取得する
- S3 から指定レポートを取得する
- S3 上に実体がないレポートに対応する DynamoDB サマリを削除する

## 前提条件

- Go 1.24 以上
- AWS の認証情報が利用可能であること
- 対象リソースが `ap-northeast-1` に存在すること

必要に応じて `--role-arn` で AssumeRole を指定できます。

> [!INFO]
> AWS の認証情報は、実行環境で利用可能な標準的な認証設定から解決されます。

## インストール

リリース済みバイナリを使う場合は GitHub Releases から取得します。

- <https://github.com/tadashi-aikawa/miroir-cli/releases>

ソースからビルドする場合:

```bash
go build -o miroir
```

## 設定ファイル

`~/.miroirconfig` が存在する場合、各コマンドのデフォルト値として読み込まれます。設定ファイルがなくても CLI は起動しますが、必須値はオプションで与える必要があります。

設定例:

```toml
table = "miroir-summaries"
bucket = "miroir-report-bucket"
bucket_prefix = "production"
role_arn = "arn:aws:iam::123456789012:role/miroir-readonly"
```

利用できるキー:

| キー | 説明 |
| --- | --- |
| `table` | DynamoDB のテーブル名 |
| `bucket` | S3 バケット名 |
| `bucket_prefix` | S3 上のプレフィックス |
| `role_arn` | AssumeRole に使う IAM Role ARN |

## 使い方

```text
Usage:
  miroir get summaries [--table=<table>] [--role-arn=<role_arn>]
  miroir get report <key> [--bucket=<bucket>] [--bucket-prefix=<bucket-prefix>] [--role-arn=<role_arn>]
  miroir prune [--table=<table>] [--bucket=<bucket>] [--bucket-prefix=<bucket-prefix>] [--dry] [--role-arn=<role_arn>]
  miroir --help
```

### サマリ一覧を取得する

```bash
./miroir get summaries
```

設定ファイルを使わず直接指定する例:

```bash
./miroir get summaries --table miroir-summaries
```

出力はタブ区切りで、以下の順に並びます。

| 順番 | 項目 |
| --- | --- |
| 1 | `begin_time` |
| 2 | `hashkey` の先頭 7 文字 |
| 3 | `same_count` |
| 4 | `different_count` |
| 5 | `failure_count` |
| 6 | `title` |

### レポートを取得する

```bash
./miroir get report 0123456789abcdef
```

このコマンドは指定したキーに対応するレポート JSON を標準出力に出します。

| ファイル | 用途 |
| --- | --- |
| `results/<key>/trials.json` | trial 情報 |
| `results/<key>/report-without-trials.json` | trial を除く本体レポート |

`bucket_prefix` が設定されている場合は、先頭に `<bucket_prefix>/` が付きます。

直接指定する例:

```bash
./miroir get report 0123456789abcdef \
  --bucket miroir-report-bucket \
  --bucket-prefix production
```

### 不整合なサマリを削除する

```bash
./miroir prune
```

`prune` は対応するレポート実体が存在しないサマリを削除します。

まず確認したい場合は `--dry` を使います。

```bash
./miroir prune --dry
```

> [!INFO]
> `prune` は S3 上の `results/<key>` 配下の有無を見て、対応する DynamoDB サマリを削除します。

## オプション

| オプション | 説明 |
| --- | --- |
| `--table`, `-t` | DynamoDB テーブル名 |
| `--bucket`, `-b` | S3 バケット名 |
| `--bucket-prefix`, `-B` | S3 プレフィックス |
| `--role-arn`, `-a` | AssumeRole 用 ARN |
| `--dry`, `-d` | `prune` を削除せず確認のみで実行する |
| `--help`, `-h` | ヘルプを表示する |
| `--version`, `-v` | バージョンを表示する |

## 注意点

- 対象リージョンは `ap-northeast-1` です
- `prune` は `--dry` を付けないと DynamoDB の項目を削除します

> [!INFO]
> 現在の実装では `get summaries` と `prune` は DynamoDB の全件取得を前提に動作します。

## 開発

依存関係の整理:

```bash
make init
```

各 OS 向けパッケージ作成:

```bash
make package-linux
make package-macos
make package-windows
```

生成物は `dist/` に出力されます。

## リリース

必要なツール:

| ツール | 用途 |
| --- | --- |
| `make` | リリース手順の実行 |
| `bash` | Makefile 内のシェル実行 |
| `go` | ビルド |
| `git` | バージョン更新コミットと push |
| `7z` | Windows 向け zip 作成 |
| `ghr` | GitHub Releases へのアップロード |

リリースブランチ名をそのままバージョンとして使います。`make release` は `args.go` の `version` をブランチ名に更新し、Linux / Windows 向けパッケージを作成してコミット・push まで行います。

```bash
make release
```

その後、Pull Request を作成してマージします。マージ後に GitHub Releases へアップロードします。

```bash
make deploy version=x.y.z
```

## ローカル検証用メモ

`moto/docker-compose.yml` には Moto Server の定義があります。

```bash
docker compose -f moto/docker-compose.yml up
```
