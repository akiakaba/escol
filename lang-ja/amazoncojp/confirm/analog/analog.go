// Package analog is not digital purchases
package analog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/internal/parse"
	"github.com/akiakaba/escol/mu"
)

type Receipt struct {
	ShopName string `json:"shop_name"` // 販売者
	Name     string `json:"name"`      // あなたの名前
	// 領収書/購入明細書
	//
	// 例えばマーケットプレイスと Amazon.co.jp から同時に購入すると、メールは1通だが注文は複数の明細に分割されるように見受けられる。
	Details []Detail `json:"details"`

	Subject string `json:"-"`
	Body    string `json:"-"`
}

func (r *Receipt) OrderDetailURLs() []string {
	var urls []string
	for _, d := range r.Details {
		urls = append(urls, d.OrderDetailURL())
	}
	return urls
}

type Detail struct {
	OrderID       string   `json:"order_id"`       // 注文番号
	OrderDate     string   `json:"order_date"`     // 注文日
	TotalAmount   int      `json:"total_amount"`   // 注文合計
	PaymentMethod string   `json:"payment_method"` // 支払い方法
	Delivery      Delivery `json:"delivery"`
}

type Delivery struct {
	EstimatedDate string `json:"estimated_date"` // お届け予定日
	Option        string `json:"option"`         // 配送オプション
	Destination   string `json:"destination"`    // お届け先
}

func (d *Detail) OrderDetailURL() string {
	return fmt.Sprintf("https://www.amazon.co.jp/gp/your-account/order-details/ref=_or?ie=UTF8&orderID=%s", d.OrderID)
}

func Filter(mail escol.Mail) bool {
	return mail.From() == `"Amazon.co.jp" <auto-confirm@amazon.co.jp>` &&
		strings.HasPrefix(mail.Subject(), "Amazon.co.jpでのご注文")
}

func Scrape(mail escol.Mail) (*Receipt, error) {
	if !Filter(mail) {
		return nil, fmt.Errorf("filtered")
	}
	plain, found := mail.FindPart("text/plain")
	if !found {
		return nil, fmt.Errorf("text/plain parts not found. from:%s", mail.From())
	}
	body, err := mu.DecodeBase64(plain.Body())
	if err != nil {
		return nil, err
	}
	r := &Receipt{
		Subject: mail.Subject(),
		Body:    mu.NormalizeSpaces(body),
	}
	split := strings.Split(r.Body, " ================================================================================= ")
	if len(split) < 3 {
		return r, fmt.Errorf("splitBody split unmatched: %s", body)
	}
	headerPart := split[0]
	receiptParts := split[1 : len(split)-1]

	var gDelivery *Delivery
	{
		matches := regexp.MustCompile(`_+ (.+) 様 Amazon\.co\.jp をご利用いただき、ありがとうございます。(.+)がお客様のご注文を承ったことをお知らせいたします。`).FindStringSubmatch(headerPart)
		if len(matches) < 3 {
			return r, fmt.Errorf("match error") //todo:impl
		}
		r.Name = matches[1]
		// TODO: "および" 対応
		//matches = regexp.MustCompile(`^Amazonマーケットプレイス出品者「(.+)」$`).FindStringSubmatch(r.ShopName)
		//if len(matches) == 2 {
		//	r.MarketPlace = true
		//	r.ShopName = matches[1]
		//}
		r.ShopName = strings.TrimSpace(matches[2])
		if d, found := findDelivery(headerPart); found {
			gDelivery = d
		}
	}
	{
		for _, part := range receiptParts {
			head, tail, found := strings.Cut(part, ` _________________________________________________________________________________ `)
			if !found {
				return r, fmt.Errorf("receiptParts part cut match error. part:%s, body:%s", part, r.Body)
			}
			matches := regexp.MustCompile(`注文番号： ([\d\-]+) 注文日： ([\d/]+)`).FindStringSubmatch(head)
			if len(matches) < 3 {
				return r, fmt.Errorf("receiptParts head orderID/orderDate match error. head:%s, body:%s", head, r.Body)
			}
			orderID := matches[1]
			orderDate := matches[2]

			var delivery Delivery
			if gDelivery != nil {
				delivery = *gDelivery
			} else {
				if d, found := findDelivery(head); found {
					delivery = *d
				} else {
					delivery = Delivery{}
				}
			}

			matches = regexp.MustCompile(`注文合計： ￥ ([\d,]+) 支払い方法 (.+)$`).FindStringSubmatch(tail)
			if len(matches) < 3 {
				return r, fmt.Errorf("receiptParts tail totalAmount/paymentMethod match error. tail:%s, body:%s", tail, r.Body)
			}
			totalAmount, err := parse.ParseIntFromCommaedDecimal(matches[1])
			if err != nil {
				panic(err) // regexp is wrong
			}
			paymentMethod := matches[2]

			r.Details = append(r.Details, Detail{
				OrderID:       orderID,
				OrderDate:     orderDate,
				TotalAmount:   totalAmount,
				PaymentMethod: paymentMethod,
				Delivery:      delivery,
			})
		}
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

	return r, nil
}

func findDelivery(s string) (*Delivery, bool) {
	matches := regexp.MustCompile(` お届け予定日： (.+) 配送オプション： (.+) お届け先： (.+)$`).FindStringSubmatch(s)
	if len(matches) == 4 {
		return &Delivery{
			EstimatedDate: matches[1],
			Option:        matches[2],
			Destination:   matches[3],
		}, true
	}
	return nil, false
}
