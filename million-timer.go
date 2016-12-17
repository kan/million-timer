package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"gopkg.in/alecthomas/kingpin.v2"
)

// package version
var VERSION string

var (
	app = kingpin.New("million-timer", "checker for IDOL M@STER MillionLIVE")

	web     = app.Command("web", "web server mode")
	webPort = web.Flag("port", "httpd port").Default("5000").OverrideDefaultFromEnvar("PORT").Short('p').Int()

	chk                = app.Command("check", "checker mode")
	chkEmail           = chk.Flag("email", "your gree email").OverrideDefaultFromEnvar("MT_EMAIL").Required().String()
	chkPassword        = chk.Flag("password", "your gree password").OverrideDefaultFromEnvar("MT_PASSWORD").Required().String()
	chkPBToken         = chk.Flag("token", "your pushbullet token").OverrideDefaultFromEnvar("MT_PB_TOKEN").Required().String()
	chkDailyRewardHour = chk.Flag("daily-reward-hour", "forget daily reward report hour").Default("23").OverrideDefaultFromEnvar("MT_DAILY_REWARD_HOUR").Int()
	chkFesTimeLeftMin  = chk.Flag("fes-time-left-min", "report unfinish fes").Default("10").OverrideDefaultFromEnvar("MT_FES_TIME_LEFT_MIN").Int()
	chkRedisURL        = chk.Flag("redis-url", "redis address").OverrideDefaultFromEnvar("REDISCLOUD_URL").String()
	chkSilent          = chk.Flag("silent", "don't output").OverrideDefaultFromEnvar("MT_CHECK_SILENT").Short('s').Bool()
)

func main() {
	app.Version(VERSION)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case web.FullCommand():
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello")
		})
		err := http.ListenAndServe(fmt.Sprintf(":%d", *webPort), nil)
		if err != nil {
			panic(err)
		}
	case chk.FullCommand():
		check()
	}
}

func check() {
	config := CheckerConfig{
		PushBulletToken: *chkPBToken,
		DailyRewardHour: *chkDailyRewardHour,
		FesTimeLeftMin:  *chkFesTimeLeftMin,
		RedisURL:        *chkRedisURL,
	}
	bw := NewBrowser(*chkEmail, *chkPassword)
	err := bw.Open("/mypage")
	if err != nil {
		log.Fatal(err)
	}

	checker := NewChecker(config, *chkSilent)
	defer checker.Close()

	err = checker.CheckElement(bw, "div.appeal-theater", "send theater notify",
		"ライブ開催可能", "劇場でライブ開催が可能になりました")
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckElement(bw, `img[src="http://cdn.bn765.com/66f/ed0f9a38c2289ef0c3e44ac330a6dc3b20df133431a565312290c6b44fddb083?8a9a0804c3458dc7898adb5c5686d52f"]`, "send budokan notify",
		"PR可能", "BrandNewStageのPRが可能です")
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckElement(bw, `img[src="http://cdn.bn765.com/740/4c0451a7eef6abf31efb2c54ec72bc6845c5a93401c0519f608e1e199cf62b6b?81baa2ea232a2db783aa17ae01358cc7"]`, "send budokan notify",
		"PR活動完了", "BrandNewStageのPRが完了しています")
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`(\d+)/(\d+)`)
	err = checker.CheckText(bw, re, "li.bp-container div", "BP is full tank",
		"BP回復完了", "BPが満タン(5)になりました。フェス回しましょう")
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckText(bw, re, "li.ap-container div", "AP is full tank",
		"元気回復完了", "元気が全快しました。営業しましょう")
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckPopup(bw)
	if err != nil {
		log.Fatal(err)
	}

	err = bw.Open("/event")
	if err != nil {
		panic(err)
	}

	err = checker.CheckElement(bw, "div#mood-send-reward div.mood-send-btn a", "hitokoto can send",
		"ひとこと送信できます", "ひとこと送信して報酬ゲット")
	if err != nil {
		log.Fatal(err)
	}

	re = regexp.MustCompile(`本日の報酬 (\d+) / (\d+)`)
	err = checker.CheckTextDailyReward(bw, re,
		"div#daily_point_reward span.m-pl", "daily point unachieved",
		"デイリー報酬未達", "まだ今日のデイリー達成してないですよ。急いで!")
	if err != nil {
		log.Fatal(err)
	}

	re = regexp.MustCompile(`フィーバーライブ開催中!!`)
	err = checker.CheckText(bw, re, "div.txt-caution", "fiver live",
		"フィーバーライブ開催中", "フィーバーライブ開催中です。回しましょう!")
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckFes(bw)
	if err != nil {
		log.Fatal(err)
	}

	err = bw.Open("/birthday")
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckBirthday(bw)
	if err != nil {
		log.Fatal(err)
	}
}
