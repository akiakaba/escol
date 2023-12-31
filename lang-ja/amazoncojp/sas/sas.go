// SubscribeAndSave = 定期オトク便
package sas

import (
	"fmt"
	"regexp"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/internal/parse"
	"github.com/akiakaba/escol/internal/slices"
	"github.com/akiakaba/escol/mu"
)

type Receipt struct {
	Products    []string
	TotalAmount int
}

func Filter(mail escol.Mail) bool {
	return mail.From() == `"Amazon定期おトク便" <no-reply@amazon.co.jp>`
}

func Scrape(mail escol.Mail) (*Receipt, error) {
	plain, found := mail.FindPart("text/plain")
	if !found {
		return nil, fmt.Errorf("text/plain parts not found. from:%s", mail.From())
	}
	body, err := mu.DecodeBase64(plain.Body())
	if err != nil {
		return nil, err
	}

	products := regexp.MustCompile(`商品\d+\s+(.+?)\s+数量：\s+\d+\s+商品の価格:\s+￥[\d,]+`).FindAllStringSubmatch(body, -1)
	amountMatch := regexp.MustCompile(`合計金額\s+￥([\d,]+)`).FindStringSubmatch(body)
	if len(amountMatch) < 2 {
		return nil, fmt.Errorf("sas amount unmatched: %s", body)
	}
	totalAmount := parse.ParseIntFromCommaedDecimal(amountMatch[1])

	return &Receipt{
		Products:    slices.Map(products, func(p []string) string { return p[1] }),
		TotalAmount: totalAmount,
	}, nil
}
