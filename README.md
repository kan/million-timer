# million-timer

某ソシャゲの時間回復系をチェックしてくれるやつ。通知はPushbullet

## 使い方

* go buildしてバイナリを作るか`go run million-timer.go`で実行
* `cp config.toml.sample config.toml`して、config.tomlの必要項目を埋める
* 必要に応じてcronとかで回すと良いんじゃないかな

## 機能

以下の項目について報告してくれます

* 劇場開催
* キャラバンのお仕事完了
* BP溢れ
* 元気溢れ
* イベントのひとこと送信
* イベントのデイリー報酬未達成
* 未読のお知らせ
