package timescar

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/akiakaba/escol"
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
	aInt, err := strconv.ParseInt(strings.ReplaceAll(amount[1], ",", ""), 10, 32)
	if err != nil {
		return nil, err
	}
	return &Receipt{
		TotalAmount: int(aInt),
	}, nil
}
