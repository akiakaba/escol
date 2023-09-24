package receipt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/lang-ja/ubereats/internal"
)

type Receipt struct {
	ShopName    string
	TotalAmount int

	Subject string
	Body    string
}

func Filter(message *gmail.Message, hint *escol.Hint) bool {
	body, err := internal.ConvertBody(message)
	if err != nil {
		return false
	}
	return hint.From() == `"Uber の領収書" <noreply@uber.com>` &&
		strings.Contains(hint.Subject(), "ご注文") &&
		!strings.Contains(body, "領収書を変更いたしました。") &&
		!strings.Contains(body, "未払いのお支払いがあります。")
}

func Scrape(message *gmail.Message, hint *escol.Hint) (*Receipt, error) {
	r := &Receipt{
		Subject: hint.Subject(),
		Body:    message.Payload.Body.Data,
	}
	if !Filter(message, hint) {
		// fixme: ちゃんとしたエラー
		return r, fmt.Errorf("not target")
	}
	{
		body, err := internal.ConvertBody(message)
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
		amount, err := strconv.ParseInt(strings.ReplaceAll(amountMatches[1], ",", ""), 10, 32)
		if err != nil {
			return r, err
		}
		r.TotalAmount = int(amount)
	}
	return r, nil
}
