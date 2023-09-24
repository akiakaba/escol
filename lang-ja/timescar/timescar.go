package timescar

import (
	"regexp"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/internal/parse"
	"github.com/akiakaba/escol/mu"
)

type Receipt struct {
	OrderID     string
	TotalAmount int
}

func Filter(mail escol.Mail) bool {
	return mail.From() == "Times CAR <inquiry@share.timescar.jp>" //TODO: 十分な絞り込みか
}

func Scrape(mail escol.Mail) (*Receipt, error) {
	body, err := mu.DecodeBase64(mail.Body())
	if err != nil {
		return nil, err
	}

	amount := regexp.MustCompile(`合計金額[\s　]+([\d,]+)円`).FindStringSubmatch(body)
	totalAmount := parse.ParseIntFromCommaedDecimal(amount[1])
	return &Receipt{
		TotalAmount: totalAmount,
	}, nil
}
