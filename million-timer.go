package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/headzoo/surf"
	"github.com/xconstruct/go-pushbullet"
)

const Endpoint = "http://app.ip.bn765.com/app/index.php/mypage"
const UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12D508 Safari/600.1.4"

type MyLink struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Body  string `json:"body,omitempty"`
}

type Config struct {
	Email           string `toml:"email"`
	Password        string `toml:"password"`
	PushBulletToken string `toml:"pb-token"`
}

func pushNotify(token string, title string, body string) error {
	pb := pushbullet.New(token)
	user, _ := pb.Me()

	link := MyLink{
		Email: user.Email,
		Type:  "link",
		Title: title,
		URL:   Endpoint,
		Body:  body,
	}
	return pb.Push("/pushes", link)
}

var version string

func main() {
	var configFile = flag.String("config", "config.toml", "path to config")
	var silent = flag.Bool("silent", false, "don't output")
	var showVersion = flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s\n", version)
		return
	}

	var config Config
	_, err := toml.DecodeFile(*configFile, &config)
	if err != nil {
		panic(err)
	}

	bw := surf.NewBrowser()
	bw.SetUserAgent(UserAgent)
	err = bw.Open(Endpoint)
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

	theater, _ := bw.Find("div.appeal-theater").Html()
	if theater != "" {
		if !*silent {
			fmt.Println("send theater notify")
		}
		err = pushNotify(config.PushBulletToken,
			"ライブ開催可能", "劇場でライブ開催が可能になりました")
		if err != nil {
			panic(err)
		}
	}

	caravan, _ := bw.Find("div.appeal-caravan").Html()
	if caravan != "" {
		if !*silent {
			fmt.Println("send caravan notify")
		}
		err = pushNotify(config.PushBulletToken, "お仕事完了", "キャラバンのお仕事が完了しています")
		if err != nil {
			panic(err)
		}
	}

	re = regexp.MustCompile(`(\d+)/5`)
	matchs = re.FindSubmatch([]byte(bw.Find("li.bp-container div").Text()))
	if len(matchs) == 2 {
		if string(matchs[1]) == "5" {
			if !*silent {
				fmt.Println("BP is full tank")
			}
			err = pushNotify(config.PushBulletToken,
				"BP回復完了", "BPが満タン(5)になりました。フェス回しましょう")
			if err != nil {
				panic(err)
			}
		}
	}

	re = regexp.MustCompile(`(\d+)/(\d+)`)
	matchs = re.FindSubmatch([]byte(bw.Find("li.ap-container div").Text()))
	if len(matchs) == 3 {
		if string(matchs[1]) == string(matchs[2]) {
			if !*silent {
				fmt.Println("AP is full tank")
			}
			err = pushNotify(config.PushBulletToken,
				"元気回復完了", "元気が全快しました。営業しましょう")
			if err != nil {
				panic(err)
			}
		}
	}
}
