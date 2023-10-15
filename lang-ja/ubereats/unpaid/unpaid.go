package unpaid

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/lang-ja/ubereats/internal"
)

type Unpaid struct {
	ShopName string `json:"shop_name"`

	Subject string `json:"-"`
	Body    string `json:"-"`
}

func Filter(mail escol.Mail) bool {
	body, err := internal.ConvertBody(mail.Body())
	if err != nil {
		return false
	}
	return mail.From() == `"Uber の領収書" <noreply@uber.com>` &&
		strings.Contains(body, "未払いのお支払いがあります。")
}

func Scrape(mail escol.Mail) (*Unpaid, error) {
	r := &Unpaid{
		Subject: mail.Subject(),
		Body:    mail.Body(),
	}
	if !Filter(mail) {
		// fixme: ちゃんとしたエラー
		return r, fmt.Errorf("not target")
	}
	{
		body, err := internal.ConvertBody(mail.Body())
		if err != nil {
			return r, err
		}
		r.Body = body
	}
	{
		shopMatches := regexp.MustCompile(`ありがとうございます。?[\s　]+(.+?)の領収書をお受け取りください。`).FindStringSubmatch(r.Body)
		if len(shopMatches) < 2 {
			return r, fmt.Errorf("len(shopMatches): %v, body: %s", len(shopMatches), r.Body)
		}
		r.ShopName = shopMatches[1]
	}
	return r, nil
}
