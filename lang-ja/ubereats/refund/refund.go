package refund

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/lang-ja/ubereats/internal"
)

type Refund struct {
	ShopName       string
	TotalAmountInt int

	Subject string
	Body    string
}

func Filter(mail escol.Mail) bool {
	body, err := internal.ConvertBody(mail)
	if err != nil {
		return false
	}
	return mail.From() == `"Uber の領収書" <noreply@uber.com>` &&
		strings.Contains(body, "領収書を変更いたしました。")
}

func Scrape(mail escol.Mail) (*Refund, error) {
	r := &Refund{
		Subject: mail.Subject(),
		Body:    mail.Body(),
	}
	if !Filter(mail) {
		// fixme: ちゃんとしたエラー
		return r, fmt.Errorf("not target")
	}
	{
		body, err := internal.ConvertBody(mail)
		if err != nil {
			return r, err
		}
		r.Body = body
	}
	{
		shopMatches := regexp.MustCompile(`ありがとうございます。?[\s　]+(.+?)の領収書を変更いたしました。`).FindStringSubmatch(r.Body)
		if len(shopMatches) < 2 {
			return r, fmt.Errorf("len(shopMatches): %v, body: %s", len(shopMatches), r.Body)
		}
		r.ShopName = shopMatches[1]
	}
	{
		amountMatches := regexp.MustCompile(`-￥([\d,]+)\s+返金`).FindStringSubmatch(r.Body)
		if len(amountMatches) < 2 {
			return r, fmt.Errorf("len(amountMatches): %v, body: %s", len(amountMatches), r.Body)
		}
		amount, err := strconv.ParseInt(strings.ReplaceAll(amountMatches[1], ",", ""), 10, 32)
		if err != nil {
			return r, err
		}
		r.TotalAmountInt = -int(amount)
	}
	return r, nil
}
