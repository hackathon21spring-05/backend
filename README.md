# linQ

linQ backend

開発環境しか作ってないです……🙏

## Develop environment

#### Requirements

- docker
- docker-compose
- mariadb（←どうしよ？とりあえずはいらないか）

1. 以下をプロジェクトルートで実行して起動
```
docker compose up
```

2. なんかいろいろ
- `http://localhost:7777` バックエンドサーバー
- `mariadb -h 127.0.0.1 -u user -p` mariadb
    - password: `password`
    - database: `linq`
- `docker compose exec db mariadb -u user -p` でもmariadbに繋がるはず
    - password `password`
    - database: `linq`
