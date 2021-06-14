# linQ

linQ backend

開発環境しか作ってないです……🙏

## Develop environment

#### Requirements

- docker
- docker-compose
- mariadb（←どうしよ？とりあえずはいらないか）

1. Run the following command in the project root
```
docker compose up
```

- `http://localhost:7777` for backend server
- `mariadb -h 127.0.0.1 -u user -p` mariadb
    - password: `password`
    - database: `linq`
- `docker compose exec db mariadb -u user -p` でもmariadbに繋がるはず
    - password `password`
    - database: `linq`
