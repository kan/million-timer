# million-timer

某ソシャゲの時間回復系をチェックしてくれるやつ。通知はPushbullet

## 使い方

* [releases](https://github.com/kan/million-timer/releases)から自分の環境向けのバイナリを取得
* `./million-timer`を実行すると設定ファイル"config.toml"の雛形が作られる
* config.tomlを編集して必要項目を埋める
* 必要に応じてcronとかで回すと良いんじゃないかな

## 実行オプション

詳細は`-help`オプションで

* `-config` 設定ファイルのパスを指定。cronで動かす時向け
* `-silent` 標準出力にエラー以外吐かなくなる。cronで動かす時向け

## ビルド方法

```sh
git clone https://github.com/kan/million-timer.git
cd million-timer
go get .
go build
```

## 機能

以下の項目について報告してくれます

* 劇場開催
* キャラバンのお仕事完了
* BP溢れ
* 元気溢れ
* イベントのひとこと送信
* イベントのデイリー報酬未達成
* 未読のお知らせ
* 終了間近の合同フェス
* 誕生日のアイドルがいる時に祝福できるか
 * 同じく、プレゼントを贈れるか
