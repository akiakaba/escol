package timescar

import (
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/mu"
)

type Receipt struct {
	OrderID     string
	TotalAmount int
}

func Filter(message *gmail.Message, hints *escol.Hint) bool {
	return hints.From() == "Times CAR <inquiry@share.timescar.jp>" //TODO: 十分な絞り込みか
}

func Scrape(message *gmail.Message, hints *escol.Hint) (*Receipt, error) {
	body, err := mu.DecodeBase64(message.Payload.Body.Data)
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
