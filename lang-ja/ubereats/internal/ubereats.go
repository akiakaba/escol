package internal

import (
	"html"

	"github.com/akiakaba/escol/mu"
)

func ConvertBody(body string) (string, error) {
	body, err := mu.DecodeBase64(body)
	if err != nil {
		return "", err
	}
	return mu.NormalizeSpaces(html.UnescapeString(mu.RemoveHTMLTag(body))), nil
}
