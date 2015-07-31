package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/BurntSushi/toml"
)

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

	bw := NewBrowser(config.Email, config.Password)
	bw.Open("/mypage")

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
	err = checker.CheckTextAt(bw, re, 23, "div#daily_point_reward span.m-pl", "daily point unachieved",
		"デイリー報酬未達", "まだ今日のデイリー達成してないですよ。急いで!")
	if err != nil {
		log.Fatal(err)
	}
}
