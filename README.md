# asamit-backend

## Overview

第 34 回 高専プロコン作品 Asamit！のバックエンド

## Requirement

### OS

- Mac OS Ventura 13.0(動作確認済み)

### Library

- Go
  - Gin
  - GORM
- Docker
- docker-compose

## Installation(local)

1. Clone this repository

```
git clone git@github.com:Asamit-NITTC/asamit-backend.git
```

2. Change directory

```
cd asamit-backend
```

3. Build docker image

```
docker-compose up -d
```

4. Create database

ボリュームマウントでは何故か docker-compose.yml ファイルで DB 構築出来なかった為仕方なく...

```
docker-compose exec -it mysql bash
```

```
mysql -u root -p
```

```
password
```

```
CREATE DATABASE asamit;
```

5. Create env file
```
touch authenv.env
```

```
#LINE関連
CLIENT_ID=hogehoge
TEST_SUB=hogehoge

CLOUD_SQL_USER_NAME=hogehoge
CLOUD_SQL_PASSWORD=hogehoge
CLOUD_SQL_IP=hogehoge
CLOUD_SQL_PORT=3306


PROJECT_ID=hogehoge

PORT=8080

MODE=DEBUG

#ローカルでCloudStorageにアクセスする際はdocker内のファイルパスを指定
GOOGLE_APPLICATION_CREDENTIALS=/go/src/api/credential/asami-gorugo.json
```

6.Create credential directory

IAMキーファイルを保存する

```
mkdir credential
```

## Usage(local)

1. Build & start container

```
docker-compose up -d
```

2.

```
docker-compose exec api go run main.go
```

## Terraformについて

実証実験用(重要なユーザー情報を保存しない)で用いた。

リソースを立ち上げる時
```
terraform apply
```
リソースを削除する時
```
terraform destroy
```

## Author

- [Yuta Ito](https://github.com/GoRuGoo)
- [Shuntaro Nozaki](https://github.com/shun-harutaro)
