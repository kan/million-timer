package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/headzoo/surf/browser"
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
func (c *Checker) CheckElement(bw *browser.Browser, s, msg, title, body string) error {
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

func (c *Checker) checkTextCore(bw *browser.Browser, r *regexp.Regexp, f func(m [][]byte) bool, s, msg, title, body string) error {
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
func (c *Checker) CheckText(bw *browser.Browser, r *regexp.Regexp, s, msg, title, body string) error {
	return c.checkTextCore(bw, r, func(m [][]byte) bool { return string(m[1]) == string(m[2]) },
		s, msg, title, body)
}

// CheckTextAt compares retrieves the value by regular expression from the text in the specified selector
func (c *Checker) CheckTextAt(bw *browser.Browser, r *regexp.Regexp, hour int, s, msg, title, body string) error {
	return c.checkTextCore(bw, r,
		func(m [][]byte) bool { return string(m[1]) != string(m[2]) && time.Now().Hour() == hour },
		s, msg, title, body)
}
