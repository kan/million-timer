# million-timer

某ソシャゲの時間回復系をチェックしてくれるやつ。通知はPushbullet

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

## 使い方

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

1. 上記Heroku Buttonを押す
1. herokuのアプリケーション立ち上げ画面が表示されるので、MT\_EMAIL(GREEのアカウントメールアドレス)、MT\_PASSWORD(GREEのパスワード)、MT\_PB\_TOKEN(PushbulletのAPI Token)を入力し、「Deploy for Free」ボタンを押す
1. デプロイが完了したらHerokuの管理画面でheroku schedulerの設定を行なう
1. TASKには`million-timer check`と入力する。実行間隔は10分おきか1時間の好きなほうを指定
1. Enjoy ミリオンライブ!

