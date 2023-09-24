package q

import (
	"fmt"
	"strings"
)

type Builder struct {
	After   string
	Before  string
	From    string
	Subject string
}

var _ fmt.Stringer = Builder{}

func (b Builder) String() string {
	var ss []string
	if b.After != "" {
		ss = append(ss, "after:"+b.After)
	}
	if b.Before != "" {
		ss = append(ss, "before:"+b.Before)
	}
	if b.From != "" {
		ss = append(ss, fmt.Sprintf("from:(%s)", b.From))
	}
	if b.Subject != "" {
		ss = append(ss, fmt.Sprintf("subject:(%s)", b.Subject))
	}
	return strings.Join(ss, " ")
}
