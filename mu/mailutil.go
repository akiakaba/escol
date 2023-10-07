// Package mu is "Mail Utility"
package mu

import (
	"encoding/base64"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

func DecodeBase64(data string) (string, error) {
	bs, err := base64.URLEncoding.DecodeString(data)
	return string(bs), err
}

func EncodeBase64(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

var bmPolicy = bluemonday.NewPolicy().AddSpaceWhenStrippingTag(true)

func RemoveHTMLTag(body string) string {
	return bmPolicy.Sanitize(body)
}

func NormalizeSpaces(s string) string {
	var rs []rune
	beforeSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !beforeSpace {
				rs = append(rs, ' ')
				beforeSpace = true
			}
		} else {
			rs = append(rs, r)
			beforeSpace = false
		}
	}
	return string(rs)
}
