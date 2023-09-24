package parse

import (
	"strconv"
	"strings"
)

func ParseIntFromCommaedDecimal(v string) (int, error) {
	aInt, err := strconv.ParseInt(strings.ReplaceAll(v, ",", ""), 10, 32)
	return int(aInt), err
}
