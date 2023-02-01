package drivers

import (
	"fmt"
	"strings"
)

func DollarNParamBindWrapper(sql string) (result string) {
	var (
		parts    = strings.Split(sql, "?")
		partsLen = len(parts)
	)
	if partsLen == 1 {
		return sql
	}
	result = parts[0]
	for i := 1; i < partsLen; i++ {
		result = fmt.Sprintf("%s$%d%s", result, i, parts[i])
	}
	return
}
