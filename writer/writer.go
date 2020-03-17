package writer

import (
	"bytes"
	"fmt"
	"strings"
)

type Writer struct {
	Content   bytes.Buffer
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
		writer.Content.WriteString(strings.Repeat("\t", writer.PrevCount))
	}

	writer.Content.WriteString(msg)
	writer.Content.WriteString("\r\n")
	if strings.HasSuffix(msg, "{") {
		writer.PrevCount++
	}
}
func (writer *Writer) AddLineFmt(format string, a ...interface{}) {
	writer.AddLine(fmt.Sprintf(format, a...))
}
