package yoshikei

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/mu"
)

type Receipt struct {
	TargetWeek  string
	TotalAmount int
}

func Filter(mail escol.Mail) bool {
	return mail.From() == "meisai@yoshikei-tokyo.co.jp" //TODO: 十分な絞り込みか
}

func Scrape(mail escol.Mail) (*Receipt, error) {
	week := regexp.MustCompile(`\d+/\d+週`).FindString(mail.Subject())
	body, err := joinBody(mail)
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

func joinBody(mail escol.Mail) (string, error) {
	var ss []string
	body, err := mu.DecodeBase64(mail.Body())
	if err != nil {
		return "", err
	}
	ss = append(ss, strings.TrimSpace(body))
	for _, p := range mail.Parts() {
		ss = append(ss, p.Body())
	}
	return strings.Join(ss, "\n"), nil
}
