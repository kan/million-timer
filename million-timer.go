package main

import (
	"fmt"
	"regexp"

	"github.com/headzoo/surf"
	"github.com/xconstruct/go-pushbullet"
	"github.com/BurntSushi/toml"
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
    Email string `toml:"email"`
    Password string `toml:"password"`
	PushBulletToken string `toml:"pb-token"`
}

func pushNotify(token string, title string, body string) error {
	pb := pushbullet.New(token)
	user, _ := pb.Me()

	link := MyLink{
		Email: user.Email,
		Type: "link",
		Title: title,
		URL: Endpoint,
		Body: body,
	}
	return pb.Push("/pushes", link)
}

func main() {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
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

	re, _ := regexp.Compile("var url = \"([^\"]+)\";")
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
		fmt.Println("send theater notify")
		err = pushNotify(config.PushBulletToken,
		                 "ライブ開催可能", "劇場でライブ開催が可能になりました") 
		if err != nil {
			panic(err)
		}
	}

	caravan, _ := bw.Find("div.appeal-caravan").Html()
	if caravan != "" {
		fmt.Println("send caravan notify")
		err = pushNotify(config.PushBulletToken, "お仕事完了", "キャラバンのお仕事が完了しています")
		if err != nil {
			panic(err)
		}
	}

	re, _ = regexp.Compile("(\\d+)/5")
	matchs = re.FindSubmatch([]byte(bw.Find("li.bp-container div").Text()))
	if string(matchs[1]) == "5" {
		fmt.Println("BP is full tank")
		err = pushNotify(config.PushBulletToken,
		                 "BP回復完了", "BPが満タン(5)になりました。フェス回しましょう")
		if err != nil {
			panic(err)
		}
	}
}
