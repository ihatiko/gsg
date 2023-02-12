package utils

import (
	"fmt"
	"strings"
)

func RebuildQuery(query string, args ...any) string {
	for _, arg := range args {
		formattedData := ""
		switch arg.(type) {
		case uint, uint8, uint16, uint32, int8, int16, int32, int64, int:
			formattedData = fmt.Sprintf("%d", arg)
		default:
			formattedData = fmt.Sprintf("'%s'", arg)
		}

		query = strings.Replace(query, "?", formattedData, 1)
	}
	return query
}
