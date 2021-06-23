# linQ

linQ backend  
frontend: https://github.com/hackathon21spring-05/linq-frontend

## Develop environment

#### Requirements

- docker
- docker-compose

1. 以下をプロジェクトルートで実行して起動
```
docker compose up
```

2. その他
- `http://localhost:3001` バックエンドサーバー
- データベースへの接続
    お好きな方法でどうぞ
    - `mariadb -h 127.0.0.1 -u user -p` mariadb
        - password: `password`
        - database: `linq`
    - `docker compose exec db mariadb -u user -p`
        - password `password`
        - database: `linq`

## 環境変数について
### STAGING
docker-compose.yml 内で指定している環境変数 `STAGING: development`を削除または変更すると，oauth認証が必要になる．  
- oauth認証を有効にする場合  
    プロジェクトルートに`.env`ファイルを設置し，その中でoauth認証に用いる環境変数 `CLIENT_SECRET` と `CLIENT_ID` を設定する
- oauth認証が無効な場合  
    `STAGING: development`が設定されている場合，oauth認証は無効化される．  
    ダミーのユーザがログイン済みであるとして動作する．
