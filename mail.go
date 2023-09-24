package escol

import "time"

type Mail interface {
	Time() time.Time
	From() string
	Subject() string
	Body() string
	Snippet() string
	Parts() []Part
	FindPart(mimetype string) (Part, bool)
}

type Part interface {
	Body() string
}
