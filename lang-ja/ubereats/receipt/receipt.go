package receipt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/internal/parse"
	"github.com/akiakaba/escol/lang-ja/ubereats/internal"
)

type Receipt struct {
	ShopName    string
	TotalAmount int
	Payments    []Payment

	Subject string
	Body    string
}

type Payment struct {
	Method string
	Date   string
	Amount int
}

func Filter(mail escol.Mail) bool {
	body, err := internal.ConvertBody(mail.Body())
	if err != nil {
		return false
	}
	return mail.From() == `"Uber の領収書" <noreply@uber.com>` &&
		strings.Contains(mail.Subject(), "ご注文") &&
		!strings.Contains(body, "領収書を変更いたしました。") &&
		!strings.Contains(body, "未払いのお支払いがあります。")
}

func Scrape(mail escol.Mail) (*Receipt, error) {
	r := &Receipt{
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
	{
		amountMatches := regexp.MustCompile(`合計\s+￥([\d,]+)`).FindStringSubmatch(r.Body)
		if len(amountMatches) < 2 {
			return r, fmt.Errorf("len(amountMatches): %v, body: %s", len(amountMatches), r.Body)
		}
		r.TotalAmount = parse.ParseIntFromCommaedDecimal(amountMatches[1])
	}
	{
		paymentAreaMatches := regexp.MustCompile(`お支払い (.+) 注文情報のページ にアクセスして、`).FindStringSubmatch(r.Body)
		if len(paymentAreaMatches) < 2 {
			return r, fmt.Errorf("len(paymentAreaMatches): %v, body: %s", len(paymentAreaMatches), r.Body)
		}
		paymentMatches := regexp.MustCompile(`\s*(.+?)\s+(\d+/\d+/\d+ \d+:\d+)\s+￥([\d,]+)\s*`).FindAllStringSubmatch(paymentAreaMatches[1], -1)
		if len(paymentMatches) < 1 {
			return r, fmt.Errorf("len(paymentMatches): %v, body: %s", len(paymentMatches), r.Body)
		}
		for _, p := range paymentMatches {
			r.Payments = append(r.Payments, Payment{
				Method: p[1],
				Date:   p[2],
				Amount: parse.ParseIntFromCommaedDecimal(p[3]),
			})
		}

	}
	return r, nil
}
