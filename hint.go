package escol

import "google.golang.org/api/gmail/v1"

type Hint struct {
	headers map[string]string
}

func (h *Hint) From() string {
	return h.headers["From"]
}

func (h *Hint) Subject() string {
	return h.headers["Subject"]
}

func ParseHint(m *gmail.Message) *Hint {
	hints := Hint{headers: make(map[string]string, len(m.Payload.Headers))}
	for _, v := range m.Payload.Headers {
		// fixme: 1key 1value の保証はない
		hints.headers[v.Name] = v.Value
	}
	return &hints
}
