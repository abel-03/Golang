package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var formatted strings.Builder
	argIter := 0
	indexIter := 0

	for i := 0; i < len(format); i++ {
		if format[i] == '{' {
			i++
			if i < len(format) && format[i] == '}' {
				formatted.WriteString(fmt.Sprint(args[argIter]))
			} else {
				index := ""
				for ; i < len(format) && format[i] != '}'; i++ {
					index += string(format[i])
				}
				indexIter, _ = strconv.Atoi(index)
				formatted.WriteString(fmt.Sprint(args[indexIter]))
			}
			argIter++
		} else {
			formatted.WriteByte(format[i])
		}
	}
	return formatted.String()
}
