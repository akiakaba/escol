package parse

import (
	"strconv"
	"strings"
)

func ParseIntFromCommaedDecimal(v string) int {
	i, err := strconv.ParseInt(strings.ReplaceAll(v, ",", ""), 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}
