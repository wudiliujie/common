package writer

import (
	"fmt"
	"strings"
)

type Writer struct {
	Content   string
	PrevCount int
}

func (writer *Writer) AddLine(msg string) {
	if strings.HasSuffix(msg, "}") {
		writer.PrevCount--
	}
	if strings.HasSuffix(msg, "});") {
		writer.PrevCount--
	}
	if writer.PrevCount < 0 {
		writer.PrevCount = 0
	}
	if writer.PrevCount > 0 {
		writer.Content += strings.Repeat("\t", writer.PrevCount)
	}

	writer.Content += msg
	writer.Content += "\r\n"
	if strings.HasSuffix(msg, "{") {
		writer.PrevCount++
	}
}
func (writer *Writer) AddLineFmt(format string, a ...interface{}) {
	writer.AddLine(fmt.Sprintf(format, a...))
}
