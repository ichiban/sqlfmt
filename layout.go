package sqlfmt

import (
	"fmt"
	"io"
	"strings"
)

type Layout interface {
	Write(io.Writer, int) error
	Width() int
	Gutter() int
}

type Atom string

func (a Atom) Write(w io.Writer, indent int) error {
	_, err := fmt.Fprintf(w, "%s", a)
	return err
}

func (a Atom) Width() int {
	return len(a)
}

func (a Atom) Gutter() int {
	return len(a)
}

type Stack []Layout

func (s Stack) Write(w io.Writer, indent int) error {
	var m int
	for _, l := range s {
		g := l.Gutter()
		if g > m {
			m = g
		}
	}

	for i, l := range s {
		indent := indent + m - l.Gutter()
		if i != 0 {
			if _, err := fmt.Fprintf(w, "\n%s", strings.Repeat(" ", indent)); err != nil {
				return err
			}
		}
		if err := l.Write(w, indent); err != nil {
			return err
		}
	}
	return nil
}

func (s Stack) Width() int {
	if len(s) == 0 {
		return 0
	}
	return s[len(s)-1].Width()
}

func (s Stack) Gutter() int {
	return s[0].Gutter()
}

func (s *Stack) Append(ls ...Layout) {
	for _, l := range ls {
		if ls, ok := l.(Stack); ok {
			*s = append(*s, ls...)
		} else {
			*s = append(*s, l)
		}
	}
}

func (s Stack) Strip() Layout {
	if len(s) == 1 {
		return s[0]
	}
	return s
}

type Juxtaposition []Layout

func (j Juxtaposition) Write(w io.Writer, indent int) error {
	for i, l := range j {
		if i != 0 {
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
			indent++
		}
		if err := l.Write(w, indent); err != nil {
			return err
		}
		indent += l.Width()
	}
	return nil
}

func (j Juxtaposition) Width() int {
	var w int
	for i, l := range j {
		if i != 0 {
			w++
		}
		w += l.Width()
	}
	return w
}

func (j Juxtaposition) Gutter() int {
	return j[0].Gutter()
}

func (j *Juxtaposition) Append(ls ...Layout) {
	for _, l := range ls {
		if ls, ok := l.(Juxtaposition); ok {
			*j = append(*j, ls...)
		} else {
			*j = append(*j, l)
		}
	}
}

func (j Juxtaposition) Strip() Layout {
	if len(j) == 1 {
		return j[0]
	}
	return j
}

type Concatenation []Layout

func (c Concatenation) Write(w io.Writer, indent int) error {
	for _, l := range c {
		if err := l.Write(w, indent); err != nil {
			return err
		}
		indent += l.Width()
	}
	return nil
}

func (c Concatenation) Width() int {
	var w int
	for _, l := range c {
		w += l.Width()
	}
	return w
}

func (c Concatenation) Gutter() int {
	var g int
	for _, l := range c {
		g += l.Gutter()
	}
	return g
}

func (c *Concatenation) Append(ls ...Layout) {
	for _, l := range ls {
		if ls, ok := l.(Concatenation); ok {
			*c = append(*c, ls...)
		} else {
			*c = append(*c, l)
		}
	}
}

func (c Concatenation) Strip() Layout {
	if len(c) == 1 {
		return c[0]
	}
	return c
}

type Choice struct {
	One     Layout
	Another Layout
}
