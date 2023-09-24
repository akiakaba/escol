package q

import (
	"fmt"
	"strings"
)

type Builder struct {
	After  string
	Before string
	From   string
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
	return strings.Join(ss, " ")
}
