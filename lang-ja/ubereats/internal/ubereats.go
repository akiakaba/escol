package internal

import (
	"html"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/mu"
)

func ConvertBody(mail escol.Mail) (string, error) {
	body, err := mu.DecodeBase64(mail.Body())
	if err != nil {
		return "", err
	}
	return mu.NormalizeSpaces(html.UnescapeString(mu.RemoveHTMLTag(body))), nil
}
