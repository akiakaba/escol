// Package digital is digital purchases
package digital

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/gmu"
	"github.com/akiakaba/escol/mu"
)

type Receipt struct {
	OrderID     string `json:"order_id"`     // 注文番号
	Subtotal    int    `json:"subtotal"`     // 商品の小計
	UsedPoint   int    `json:"used_point"`   // 消費したAmazonポイント
	TotalAmount int    `json:"total_amount"` // 注文合計
	Title       string `json:"title"`        // 商品のタイトル
	Publisher   string `json:"publisher"`    // 商品の販売者
	Name        string `json:"name"`         // あなたの名前
	Email       string `json:"email"`        // あなたのメールアドレス

	Subject string `json:"-"`
	Body    string `json:"-"`
}

func (r *Receipt) OrderDetailURL() string {
	return fmt.Sprintf("https://www.amazon.co.jp/gp/digital/your-account/order-summary.html?ie=UTF8&orderID=%s", r.OrderID)
}

func Filter(message *gmail.Message, hint *escol.Hint) bool {
	return hint.From() == `"Amazon.co.jp" <digital-no-reply@amazon.co.jp>` &&
		strings.HasPrefix(hint.Subject(), "Amazon.co.jpでのご注文")
}

func Scrape(message *gmail.Message, hint *escol.Hint) (*Receipt, error) {
	if !Filter(message, hint) {
		return nil, fmt.Errorf("filtered")
	}
	// in case of errors
	r := &Receipt{
		Subject: hint.Subject(),
	}
	plain, found := gmu.FindPartByMimeType(message.Payload, "text/plain")
	if !found {
		return r, fmt.Errorf("text/plain parts not found. from:%s", hint.From())
	}
	r.Body = plain.Body.Data

	body, err := mu.DecodeBase64(plain.Body.Data)
	if err != nil {
		return r, err
	}
	body = mu.NormalizeSpaces(body)
	r.Body = body
	re := regexp.MustCompile(
		`^(.+) 様 ご購入ありがとうございます。購入したKindle本はクラウドに保存され、「コンテンツと端末の管理」から確認できます。\s+` +
			`コンテンツと端末の管理: \S+ 領収書をご希望の場合は、お客様ご自身で印刷することができます。詳しくは、ヘルプページの「領収書を印刷する」をご覧ください。\s+` +
			`領収書を印刷する: \S+ \.+ 注文情報 Eメールアドレス: (\S+) 注文合計: ￥ [\d,]+ \.+ 注文内容 注文番号: (D[\d\-]+) 商品の小計: ￥ ([\d,]+)( Amazonポイント： -￥ ([\d,]+) (</tr>)?)? \.+ 注文合計: ￥ ([\d,]+) \.+ (.+) 販売: (.+) \.+ 注文内容は、「アカウントサービス」からご確認いただけます。`,
	)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 10 {
		return r, fmt.Errorf("match error")
	}

	//// debug
	//var m []string
	//for i, mm := range matches {
	//	if i == 0 {
	//		continue
	//	}
	//	m = append(m, fmt.Sprintf("[%v:%s]", i, mm))
	//}
	//fmt.Println(hint.Subject())
	//fmt.Println(strings.Join(m, ""))

	return &Receipt{
		Name:    matches[1],
		Email:   matches[2],
		OrderID: matches[3],
		Subtotal: func() int {
			i, _ := strconv.ParseInt(strings.ReplaceAll(matches[4], ",", ""), 10, 32)
			return int(i)
		}(),
		UsedPoint: func() int {
			if matches[6] == "" {
				return 0
			}
			i, _ := strconv.ParseInt(strings.ReplaceAll(matches[6], ",", ""), 10, 32)
			return int(i)
		}(),
		TotalAmount: func() int {
			i, _ := strconv.ParseInt(strings.ReplaceAll(matches[8], ",", ""), 10, 32)
			return int(i)
		}(),
		Title:     matches[9],
		Publisher: matches[10],
		Subject:   hint.Subject(),
		Body:      body,
	}, nil
}
