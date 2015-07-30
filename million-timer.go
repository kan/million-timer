package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/headzoo/surf"
)

const endpoint = "http://app.ip.bn765.com/app/index.php"
const userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12D508 Safari/600.1.4"

type appConfig struct {
	Email           string `toml:"email"`
	Password        string `toml:"password"`
	PushBulletToken string `toml:"pb-token"`
}

// package version
var VERSION string

func main() {
	var configFile = flag.String("config", "config.toml", "path to config")
	var silent = flag.Bool("silent", false, "don't output")
	var showVersion = flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s\n", VERSION)
		return
	}

	var config appConfig
	_, err := toml.DecodeFile(*configFile, &config)
	if err != nil {
		panic(err)
	}

	bw := surf.NewBrowser()
	bw.SetUserAgent(userAgent)
	err = bw.Open(endpoint + "/mypage")
	if err != nil {
		panic(err)
	}

	fm, _ := bw.Form("form#login")
	fm.Input("mail", config.Email)
	fm.Input("user_password", config.Password)
	err = fm.Submit()
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`var url = "([^"]+)";`)
	matchs := re.FindSubmatch([]byte(bw.Find("script").Text()))
	err = bw.Open(string(matchs[1]))
	if err != nil {
		panic(err)
	}

	fm, _ = bw.Form("form")
	err = fm.Submit()
	if err != nil {
		panic(err)
	}

	checker := NewChecker(config.PushBulletToken, *silent)
	defer checker.Close()

	err = checker.CheckElement(bw, "div.appeal-theater", "send theater notify",
		"ライブ開催可能", "劇場でライブ開催が可能になりました")
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckElement(bw, "div.appeal-caravan", "send caravan notify",
		"お仕事完了", "キャラバンのお仕事が完了しています")
	if err != nil {
		log.Fatal(err)
	}

	re = regexp.MustCompile(`(\d+)/(\d+)`)
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

	err = bw.Open(endpoint + "/event")
	if err != nil {
		panic(err)
	}

	fm, _ = bw.Form("form")
	err = fm.Submit()
	if err != nil {
		panic(err)
	}

	err = checker.CheckElement(bw, "div#mood-send-reward div.mood-send-btn a", "hitokoto can send",
		"ひとこと送信できます", "ひとこと送信して報酬ゲット")
	if err != nil {
		log.Fatal(err)
	}

	re = regexp.MustCompile(`本日の報酬 (\d+) / (\d+)`)
	err = checker.CheckTextAt(bw, re, 23, "div#daily_point_reward span.m-pl", "daily point unachieved",
		"デイリー報酬未達", "まだ今日のデイリー達成してないですよ。急いで!")
	if err != nil {
		log.Fatal(err)
	}
}
