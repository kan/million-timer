package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/xconstruct/go-pushbullet"
)

// Checker holds the settings for the checker of the millon-timer
type Checker struct {
	PBClient *pushbullet.Client
	Target   *pushbullet.User
	Silent   bool
	Cache    *simplejson.Json
}

type pbLink struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Body  string `json:"body,omitempty"`
}

// NewChecker is generator for Checker
func NewChecker(token string, s bool) *Checker {
	pb := pushbullet.New(token)
	user, _ := pb.Me()

	var r *os.File
	_, err := os.Stat(".million-timer")
	if err == nil {
		r, _ = os.Open(".million-timer")
	} else {
		r, _ = os.Create(".million-timer")
	}
	c, err := simplejson.NewFromReader(r)
	if err != nil {
		c = simplejson.New()
	}

	checker := &Checker{PBClient: pb, Target: user, Silent: s, Cache: c}

	return checker
}

// Close closes the Checker
func (c *Checker) Close() {
	w, _ := os.Create(".million-timer")
	defer w.Close()
	b, _ := c.Cache.EncodePretty()
	w.Write(b)
}

func (c *Checker) pushNotify(title string, body string) error {
	link := pbLink{
		Email: c.Target.Email,
		Type:  "link",
		Title: title,
		URL:   endpoint + "/mypage",
		Body:  body,
	}
	return c.PBClient.Push("/pushes", link)
}

// CheckElement is checks whether the specified selector is present in the corresponding page
func (c *Checker) CheckElement(bw *Browser, s, msg, title, body string) error {
	html, _ := bw.Find(s).Html()
	if html != "" {
		flg := c.Cache.Get("CheckElement:" + s).MustBool(false)
		if !flg {
			if !c.Silent {
				fmt.Println(msg)
			}
			c.Cache.Set("CheckElement:"+s, true)
			return c.pushNotify(title, body)
		}
	} else {
		c.Cache.Del("CheckElement:" + s)
	}
	return nil
}

// CheckPopup is checks popup list
func (c *Checker) CheckPopup(bw *Browser) error {
	n := false
	info := ""
	m := c.Cache.Get("CheckPopup").MustMap(make(map[string]interface{}))

	bw.Find("div#main-img div#popup ul li a").Each(func(_ int, s *goquery.Selection) {
		t := s.Text()
		if t == "合同フェスへの参加要請が届いています" {
			return
		}
		n = true
		_, ok := m[t]
		if ok {
			return
		}
		if !c.Silent {
			fmt.Println("new info: " + t)
		}
		info = info + t + "\n"
		m[t] = 1
	})

	if n {
		c.Cache.Set("CheckPopup", m)
		if info != "" {
			c.pushNotify("未読のお知らせがあります", info)
		}
	} else {
		c.Cache.Del("CheckPopup")
	}

	return nil
}

// CheckFes is checks fes
func (c *Checker) CheckFes(bw *Browser) error {
	flg := false

	checker := func() {
		bw.Find("ul.list-bg li table dd.txt-ngtv").Each(func(_ int, s *goquery.Selection) {
			t, _ := time.Parse("15:04:05", s.Text())
			if t.Minute() <= 10 {
				flg = true
			}
		})
	}

	for _, path := range []string{"/fes/event_multi_list","/fes/event_list","/fes"} {
		err := bw.Open(path)
		if err != nil {
			return err
		}
		checker()
	}

	if flg {
		if !c.Silent {
			fmt.Println("exist near the end of fes")
		}
		c.pushNotify("終了目前のフェスがあります", "勿体ないので処理しましょう")
	}

	return nil
}

func (c *Checker) checkTextCore(bw *Browser, r *regexp.Regexp, f func(m [][]byte) bool, s, msg, title, body string) error {
	matchs := r.FindSubmatch([]byte(bw.Find(s).Text()))
	if len(matchs) == 3 {
		if f(matchs) {
			flg := c.Cache.Get("CheckText:" + s + r.String()).MustBool(false)
			if !flg {
				if !c.Silent {
					fmt.Println(msg)
				}
				c.Cache.Set("CheckText:"+s+r.String(), true)
				return c.pushNotify(title, body)
			}
		} else {
			c.Cache.Del("CheckText:" + s + r.String())
		}
	}
	return nil
}

// CheckText compares retrieves the value by regular expression from the text in the specified selector
func (c *Checker) CheckText(bw *Browser, r *regexp.Regexp, s, msg, title, body string) error {
	return c.checkTextCore(bw, r, func(m [][]byte) bool { return string(m[1]) == string(m[2]) },
		s, msg, title, body)
}

// CheckTextAt compares retrieves the value by regular expression from the text in the specified selector
func (c *Checker) CheckTextAt(bw *Browser, r *regexp.Regexp, hour int, s, msg, title, body string) error {
	return c.checkTextCore(bw, r,
		func(m [][]byte) bool { return string(m[1]) != string(m[2]) && time.Now().Hour() == hour },
		s, msg, title, body)
}
