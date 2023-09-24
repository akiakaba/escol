// SubscribeAndSave = 定期オトク便
package sas

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
	return hints.From() == `"Amazon定期おトク便" <no-reply@amazon.co.jp>`
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

	products := regexp.MustCompile(`商品\d+\s+(.+?)\s+数量：\s+\d+\s+商品の価格:\s+￥[\d,]+`).FindAllStringSubmatch(body, -1)
	amount := regexp.MustCompile(`合計金額\s+￥([\d,]+)`).FindStringSubmatch(body)
	aInt, err := strconv.ParseInt(strings.ReplaceAll(amount[1], ",", ""), 10, 32)
	if err != nil {
		return nil, err
	}
	return &Receipt{
		Products:    slices.Map(products, func(p []string) string { return p[1] }),
		TotalAmount: int(aInt),
	}, nil
}
