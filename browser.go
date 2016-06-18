package main

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

const endpoint = "http://app.ip.bn765.com/app/index.php"
const userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12D508 Safari/600.1.4"

// Browser is million-live browser wrapper for surf.Browser
type Browser struct {
	Browser  *browser.Browser
	Email    string
	Password string
}

// NewBrowser is generator for Browser
func NewBrowser(email, password string) *Browser {
	bw := surf.NewBrowser()
	bw.SetUserAgent(userAgent)

	return &Browser{Browser: bw, Email: email, Password: password}
}

func (b *Browser) login() error {
	fm, _ := b.Browser.Form("form#login")
	fm.Input("mail", b.Email)
	fm.Input("user_password", b.Password)
	err := fm.Submit()
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`var url = "([^"]+)";`)
	if matchs := re.FindSubmatch([]byte(b.Find("script").Text())); matchs != nil {
		err = b.Browser.Open(string(matchs[1]))
		if err != nil {
			return err
		}
	} else {
		fm, err = b.Browser.Form(`form[name="redirect"]`)
		if err != nil {
			return err
		}
		err = fm.Submit()
		if err != nil {
			return err
		}
	}

	return nil
}

// Open to open a path for million-live
func (b *Browser) Open(path string) error {
	err := b.Browser.Open(endpoint + path)
	if err != nil {
		return err
	}

	if m, _ := regexp.MatchString(`^https://id.gree.net/`, b.Browser.Url().String()); m {
		err = b.login()
		if err != nil {
			return err
		}
	}

	fm, err := b.Browser.Form("form")
	if err != nil {
		return nil
	}
	err = fm.Submit()
	if err != nil {
		return err
	}

	return nil
}

// Find is a wrapper for browser.Browser.Find
func (b *Browser) Find(expr string) *goquery.Selection {
	return b.Browser.Find(expr)
}
