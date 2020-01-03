package graph

import (
	"fmt"
	"strings"
)

type ToStringer interface {
	ToString() string
}

func NewWriter() *Writer {
	var sb strings.Builder
	return &Writer{
		sb: &sb,
	}
}

type Writer struct {
	sb *strings.Builder
}

func (w *Writer) WriteLn(depth int, ts ToStringer) {
	w.Write(depth, ts)
	w.sb.WriteString("\n")
}

func (w *Writer) Write(depth int, ts ToStringer) {
	tabs := strings.Repeat(" ", 4*depth)
	w.sb.WriteString(tabs)
	w.sb.WriteString(ts.ToString())
}

func (w *Writer) WriteString(depth int, s string) {
	tabs := strings.Repeat(" ", 4*depth)
	w.sb.WriteString(tabs)
	w.sb.WriteString(s)
}

func (w *Writer) WriteStringf(depth int, format string, args ...interface{}) {
	tabs := strings.Repeat(" ", 4*depth)
	w.sb.WriteString(tabs)
	s := fmt.Sprintf(format, args...)
	w.sb.WriteString(s)
}

func (w *Writer) String() string {
	return w.sb.String()
}
