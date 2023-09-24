// Package gmu is "GMail Utility"
package gmu

import (
	"strings"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol/mu"
)

func FindPartByMimeType(parent *gmail.MessagePart, mimeType string) (child *gmail.MessagePart, found bool) {
	for _, pp := range parent.Parts {
		if pp.MimeType == mimeType {
			return pp, true
		}
	}
	return nil, false
}

func JoinBody(part *gmail.MessagePart) (string, error) {
	var ss []string
	body, err := mu.DecodeBase64(part.Body.Data)
	if err != nil {
		return "", err
	}
	ss = append(ss, strings.TrimSpace(body))
	for _, p := range part.Parts {
		partBody, err := JoinBody(p)
		if err != nil {
			return "", err
		}
		ss = append(ss, partBody)
	}
	return strings.Join(ss, "\n"), nil
}
