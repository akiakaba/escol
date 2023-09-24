// Package gmu is "GMail Utility"
package gmu

import (
	"google.golang.org/api/gmail/v1"
)

func FindPartByMimeType(parent *gmail.MessagePart, mimeType string) (child *gmail.MessagePart, found bool) {
	for _, pp := range parent.Parts {
		if pp.MimeType == mimeType {
			return pp, true
		}
	}
	return nil, false
}
