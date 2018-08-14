miroir-cli
==========

CLI for Miroir.
Not support for Windows.


Install
-------

```
$ wget https://github.com/tadashi-aikawa/miroir-cli/releases/download/v0.1.0/miroir && chmod +x miroir
$ ./miroir --help                                                                                                                                   Tue 14 Aug 2018 03:29:07 PM JST
Miroir CLI

Usage:
  miroir get summaries --table=<table>
  miroir get report <key-prefix> --bucket=<bucket> [--format]
  miroir create --table=<table> --bucket=<bucket>
  miroir prune  --table=<table> --bucket=<bucket> [--dry]
  miroir --help

Options:
  -h --help            Show this screen.
  -f --format          Pretty format
  -d --dry             Dry run
  --table=<table>      DynamoDB table name
  --bucket=<bucket>    S3 bucket name
```


For developer
-------------

### Release

```
$ make release
```

### Upload

Upload `target/release/miroir`

