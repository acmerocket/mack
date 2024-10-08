package mack

type match struct {
	path  string
	lines []line
}

type line struct {
	num     int
	column  int
	text    any
	matched bool
}

func (m *match) add(num int, column int, text any, matched bool) {
	m.lines = append(m.lines, line{
		num:     num,
		column:  column,
		text:    text,
		matched: matched,
	})
}

func (m match) size() int {
	return len(m.lines)
}
