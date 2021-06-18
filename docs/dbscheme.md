# ハッカソン21春 5班 データベース

データベース(わからん)

## user table

traQのユーザーを格納するテーブル。最初のログイン時にレコードがなければ、追加。
**users**

| Name | Type        | Nullable | Key     | Default | 説明   |
| ---- | ----------- | -------- | ------- | ------- | ------ |
| id   | varchar(36) | false    | PRIMARY |         | uuid   |
| name | varchar(36) | false    |         |         | UNIQUE |

## entry table

ユーザーがブックマークした記事．  bookmark解除によって登録者がゼロ人になったら，削除する
**entrys**

| Name       | Type        | Nullable | Key              | Default | 説明      |
| ---------- | ----------- | -------- | ---------------- | ------- | --------- |
| id         | varchar(36) | false    | PRIMARY          |         | uuid      |
| url        | text        | false    |                  |         | unique    |
| title      | text        | false    |                  |         |           |
| caption    | text        |          |                  |         |           |
| thumbnail  | text        |          |                  |         | urlいれる |
| created_at | datetime    | false    | CURRENT_TIMESTMP |         |           |

:::warning
text型（固定長でない）にindexを貼ることができなかったので，idを追加しました
:::

## bookmark table

誰がブックマークをしたのか
**bookmarks**

| Name       | Type        | Nullable | Key              | Default | 説明           |
| ---------- | ----------- | -------- | ---------------- | ------- | -------------- |
| id         | int         | false    | PRIMARY          |         | AUTO_INCREMENT |
| user_id    | varchar(36) | false    |                  |         |                |
| entry_id   | varchar(32) | false    |                  |         |                |
| created_at | datetime    | false    | CURRENT_TIMESTMP |         |                |

## tag table

記事についているブックマーク一覧
**tags**
| Name       | Type        | Nullable | Key              | Default | 説明           |
| ---------- | ----------- | -------- | ---------------- | ------- | -------------- |
| id         | int         | false    | PRIMARY          |         | AUTO_INCREMENT |
| tag        | varchar(32) | false    |                  |         |                |
| entry_id   | varchar(36) | false    |                  |         |                |
| created_at | datetime    | false    | CURRENT_TIMESTMP |         |                |

## favorite table

余裕があれば
**favorites**
| Name       | Type        | Nullable | Key     | Default | 説明           |
| ---------- | ----------- | -------- | ------- | ------- | -------------- |
| id         | int         | false    | PRIMARY |         | AUTO_INCREMENT |
| user_id    | varchar(36) | false    |         |         |                |
| entry_id   | varchar(36) | false    |         |         |                |
| created_at | datetime    | false    |         |         |                |
