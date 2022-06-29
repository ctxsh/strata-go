package apex

import (
	"strconv"
	"strings"
)

func BuildProm(name string, help string, ctype string, labels map[string]string, value float64) string {
	var builder strings.Builder
	builder.WriteString("# HELP ")
	builder.WriteString(name + " ")
	builder.WriteString(help + "\n")
	builder.WriteString("# TYPE ")
	builder.WriteString(name + " ")
	builder.WriteString(ctype + "\n")
	builder.WriteString(name + "{")
	for k, v := range labels {
		builder.WriteString(k + "=\"")
		builder.WriteString(v + "\"} ")
	}

	val := strconv.FormatFloat(value, 'E', -1, 64)
	builder.WriteString(val + "\n")

	return builder.String()
}
