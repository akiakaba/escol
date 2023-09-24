package yodobashi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/gmu"
	"github.com/akiakaba/escol/internal/slices"
	"github.com/akiakaba/escol/mu"
)

type Receipt struct {
	Products    []string
	TotalAmount int
}

func Filter(message *gmail.Message, hints *escol.Hint) bool {
	return hints.From() == `"ヨドバシ・ドット・コム" <thanks_gochuumon@yodobashi.com>`
}

func Scrape(message *gmail.Message, hints *escol.Hint) (*Receipt, error) {
	plain, found := gmu.FindPartByMimeType(message.Payload, "text/plain")
	if !found {
		return nil, fmt.Errorf("text/plain parts not found. from:%s", hints.From())
	}
	body, err := mu.DecodeBase64(plain.Body.Data)
	if err != nil {
		return nil, err
	}

	_, text, _ := strings.Cut(body, "【ご注文商品】")
	text = regexp.MustCompile(`[\s　]*[\r\n]+[\s　]*`).ReplaceAllString(text, "")
	product := regexp.MustCompile(`「(.+?)」`).FindAllStringSubmatch(text, -1)
	//FIXME: panic: runtime error: index out of range [1] with length 0
	amount := regexp.MustCompile(`【ご注文金額】今回のお買い物合計金額[\s　]*([\d,]+) 円`).FindStringSubmatch(body)
	aInt, err := strconv.ParseInt(strings.ReplaceAll(amount[1], ",", ""), 10, 32)
	if err != nil {
		return nil, err
	}
	return &Receipt{
		Products:    slices.Map(product, func(p []string) string { return p[1] }),
		TotalAmount: int(aInt),
	}, nil
}
