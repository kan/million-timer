バイナリ類の構築を忘れた時のメモ

## 普段の実行

```sh
go run *.go
```

でok

## 手元でビルド

```sh
go build
```

で、million-timerファイルが出来る

## 新しいバージョンのバイナリをリリース

* .goxc.json の PackageVersion を書き換える。`goxc bump`で良い感じにやってくれる気がする。
* .goxc.local.json が無ければ作っておく
* `goxc`を叩けばクロスコンパイルからgithubへのアップロードまで全部済ませてくれる
* github上でリリースノートを適宜書き換える

## .goxc.local.jsonの作り方

```sh
goxc -wlc default publish-github -apikey=XXXXXXXXXXX
```

API Keyの値は https://github.com/settings/tokens を参照のこと

