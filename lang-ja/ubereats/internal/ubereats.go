package internal

import (
	"html"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol/mu"
)

func ConvertBody(message *gmail.Message) (string, error) {
	body, err := mu.DecodeBase64(message.Payload.Body.Data)
	if err != nil {
		return "", err
	}
	return mu.NormalizeSpaces(html.UnescapeString(mu.RemoveHTMLTag(body))), nil
}
