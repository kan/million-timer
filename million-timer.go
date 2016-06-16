package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/alecthomas/kingpin.v2"
)

type appConfig struct {
	Email          string        `toml:"email"`
	Password       string        `toml:"password"`
	CheckerSetting CheckerConfig `toml:"checker"`
}

// package version
var VERSION string

func genConfig(file string) {
	config := &appConfig{
		Email: "your gree email", Password: "your gree password", CheckerSetting: CheckerConfig{
			PushBulletToken: "your pushbullet token", DailyRewardHour: 23, FesTimeLeftMin: 10}}
	f, _ := os.Create(file)
	encoder := toml.NewEncoder(f)
	err := encoder.Encode(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("can't find %s. generate sample config.\n", file)
	fmt.Printf("please edit %s. And rerun\n", file)
}

var (
	app         = kingpin.New("million-timer", "checker for IDOL M@STER MillionLIVE")

	web     = app.Command("web", "web server mode")
	webPort = web.Flag("port", "httpd port").Default("5000").OverrideDefaultFromEnvar("PORT").Short('p').Int()

	checker    = app.Command("check", "checker mode")
	configFile = checker.Flag("config", "path to config").Default("config.toml").String()
	silent     = checker.Flag("silent", "don't output").Short('s').Bool()
)

func main() {
	app.Version(VERSION)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case web.FullCommand():
		http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello")
		})
		err := http.ListenAndServe(fmt.Sprintf(":%d", *webPort), nil)
		if err != nil {
			panic(err)
		}
	case checker.FullCommand():
		check();
	}
}

func check() {
	var config appConfig
	_, err := toml.DecodeFile(*configFile, &config)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			genConfig(*configFile)
			return
		}
		panic(err)
	}

	bw := NewBrowser(config.Email, config.Password)
	bw.Open("/mypage")

	checker := NewChecker(config.CheckerSetting, *silent)
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
