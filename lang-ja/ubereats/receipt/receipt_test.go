package receipt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/mu"
)

func TestScrape(t *testing.T) {
	type receipt struct {
		shopName    string
		totalAmount int
		payments    []Payment
	}
	tests := []struct {
		name string
		mail escol.Mail
		want *receipt
	}{
		{
			name: "Usage",
			mail: mail{
				subject: "金曜日 午前の Uber Eats のご注文",
				body:    `合計 ￥1,154 2022年7月15日 akiakaba 様、Uber One をご利用いただきありがとうございます。 デザインホームセンター 福岡店の領収書をお受け取りください。 注文を評価 注文を評価 合計 ￥1,154 利用明細を確認するには次へ移動 Uber Eats , または この PDF をダウンロードしてください お支払い Uber Cash 2022/07/15 8:07 ￥500 クレカA ••••8956 2022/07/15 18:55 ￥654 注文情報のページ にアクセスして、請求書 (利用可能な場合) などの詳細をご覧いただけます。`,
			},
			want: &receipt{
				shopName:    "デザインホームセンター 福岡店",
				totalAmount: 1154,
				payments: []Payment{
					{Method: "Uber Cash", Date: "2022/07/15 8:07", Amount: 500},
					{Method: "クレカA ••••8956", Date: "2022/07/15 18:55", Amount: 654},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Scrape(tt.mail)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.shopName, got.ShopName)
			assert.Equal(t, tt.want.totalAmount, got.TotalAmount)
			assert.Equal(t, tt.want.payments, got.Payments)
		})
	}
}

var _ escol.Mail = (*mail)(nil)

type mail struct {
	subject string
	body    string
}

func (m mail) Time() time.Time {
	panic("will not be called")
}

func (m mail) From() string {
	return `"Uber の領収書" <noreply@uber.com>`
}

func (m mail) Subject() string {
	return m.subject
}

func (m mail) Body() string {
	return mu.EncodeBase64(m.body)
}

func (m mail) Snippet() string {
	panic("will not be called")
}

func (m mail) Parts() []escol.Part {
	panic("will not be called")
}

func (m mail) FindPart(mimetype string) (escol.Part, bool) {
	panic("will not be called")
}
