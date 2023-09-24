package gmu

import (
	"time"

	"google.golang.org/api/gmail/v1"

	"github.com/akiakaba/escol"
	"github.com/akiakaba/escol/internal/slices"
)

func MailAs(m *gmail.Message) escol.Mail {
	return &mail{
		message: m,
		hint:    parseHint(m),
	}
}

type mail struct {
	message *gmail.Message
	hint    *hint
}

func (m *mail) Time() time.Time {
	return time.UnixMilli(m.message.InternalDate)
}

func (m *mail) From() string {
	return m.hint.from()
}

func (m *mail) Subject() string {
	return m.hint.subject()
}

func (m *mail) Body() string {
	return m.message.Payload.Body.Data
}

func (m *mail) Snippet() string {
	return m.message.Snippet
}

func (m *mail) Parts() []escol.Part {
	return slices.Map(m.message.Payload.Parts, func(mp *gmail.MessagePart) escol.Part {
		return &part{part: mp}
	})
}

func (m *mail) FindPart(mimetype string) (escol.Part, bool) {
	if p, found := FindPartByMimeType(m.message.Payload, mimetype); found {
		return &part{part: p}, true
	}
	return nil, false
}

type part struct {
	part *gmail.MessagePart
}

func (p *part) Body() string {
	return p.part.Body.Data
}

type hint struct {
	headers map[string]string
}

func (h *hint) from() string {
	return h.headers["From"]
}

func (h *hint) subject() string {
	return h.headers["Subject"]
}

func parseHint(m *gmail.Message) *hint {
	hints := hint{headers: make(map[string]string, len(m.Payload.Headers))}
	for _, v := range m.Payload.Headers {
		// fixme: 1key 1value の保証はない
		hints.headers[v.Name] = v.Value
	}
	return &hints
}
