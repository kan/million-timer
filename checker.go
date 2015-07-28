package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/headzoo/surf/browser"
	"github.com/xconstruct/go-pushbullet"
)

type Checker struct {
	PBClient *pushbullet.Client
	Target   *pushbullet.User
	Silent   bool
}

type pbLink struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Body  string `json:"body,omitempty"`
}

func NewChecker(token string, s bool) *Checker {
	pb := pushbullet.New(token)
	user, _ := pb.Me()
	checker := &Checker{PBClient: pb, Target: user, Silent: s}
	return checker
}

func (self *Checker) pushNotify(title string, body string) error {
	link := pbLink{
		Email: self.Target.Email,
		Type:  "link",
		Title: title,
		URL:   Endpoint + "/mypage",
		Body:  body,
	}
	return self.PBClient.Push("/pushes", link)
}

func (self *Checker) CheckElement(bw *browser.Browser, s, msg, title, body string) error {
	html, _ := bw.Find(s).Html()
	if html != "" {
		if !self.Silent {
			fmt.Println(msg)
		}
		return self.pushNotify(title, body)
	}
	return nil
}

func (self *Checker) checkTextCore(bw *browser.Browser, r *regexp.Regexp, f func(m [][]byte) bool, s, msg, title, body string) error {
	matchs := r.FindSubmatch([]byte(bw.Find(s).Text()))
	if len(matchs) == 3 {
		if f(matchs) {
			if !self.Silent {
				fmt.Println(msg)
			}
			return self.pushNotify(title, body)
		}
	}
	return nil
}

func (self *Checker) CheckText(bw *browser.Browser, r *regexp.Regexp, s, msg, title, body string) error {
	return self.checkTextCore(bw, r, func(m [][]byte) bool { return string(m[1]) == string(m[2]) },
		s, msg, title, body)
}

func (self *Checker) CheckTextAt(bw *browser.Browser, r *regexp.Regexp, hour int, s, msg, title, body string) error {
	return self.checkTextCore(bw, r,
		func(m [][]byte) bool { return string(m[1]) != string(m[2]) && time.Now().Hour() == hour },
		s, msg, title, body)
}
