package yoshikei

import (
	"regexp"
	"strconv"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/gmu"
)

type Receipt struct {
	TargetWeek  string
	TotalAmount int
}

func Filter(message *gmail.Message, hint *escol.Hint) bool {
	return hint.From() == "meisai@yoshikei-tokyo.co.jp" //TODO: 十分な絞り込みか
}

func Scrape(message *gmail.Message, hint *escol.Hint) (*Receipt, error) {
	week := regexp.MustCompile(`\d+/\d+週`).FindString(hint.Subject())
	body, err := gmu.JoinBody(message.Payload)
	if err != nil {
		return nil, err
	}
	amount := regexp.MustCompile(`合計金額[\s　]+(\d+)[\s　]+円`).FindStringSubmatch(body)
	aInt, err := strconv.ParseInt(amount[1], 10, 32)
	if err != nil {
		return nil, err
	}
	return &Receipt{
		TargetWeek:  week,
		TotalAmount: int(aInt),
	}, nil
}
