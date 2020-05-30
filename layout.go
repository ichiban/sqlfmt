package sqlfmt

import (
	"fmt"
	"io"
	"strings"
)

type Layout interface {
	Write(io.Writer, int) error
	Offset() int
	Gutter() int
}

type Atom string

func (a Atom) Write(w io.Writer, indent int) error {
	_, err := fmt.Fprintf(w, "%s", a)
	return err
}

func (a Atom) Offset() int {
	return len(a)
}

func (a Atom) Gutter() int {
	return len(a)
}

type Stack []Layout

func (s Stack) Write(w io.Writer, indent int) error {
	for i, l := range s {
		if err := l.Write(w, indent); err != nil {
			return err
		}
		if i < len(s)-1 {
			if _, err := fmt.Fprintf(w, "\n%s", strings.Repeat(" ", indent)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s Stack) Offset() int {
	if len(s) == 0 {
		return 0
	}
	return s[len(s)-1].Offset()
}

func (s Stack) Gutter() int {
	return s[0].Gutter()
}

func (s Stack) AlignGutter() {
	var m int
	for _, l := range s {
		g := l.Gutter()
		if g > m {
			m = g
		}
	}

	for i := range s {
		s[i] = gutterAligned(s[i], m)
	}
}

func gutterAligned(l Layout, g int) Layout {
	switch l := l.(type) {
	case Atom:
		return Atom(fmt.Sprintf("%*s", g, l))
	case Stack:
		r := make(Stack, len(l))
		copy(r, l)
		for i := range r {
			r[i] = gutterAligned(r[i], g)
		}
		return r
	case Juxtaposition:
		r := make(Juxtaposition, len(l))
		copy(r, l)
		r[0] = gutterAligned(r[0], g)
		return r
	case Concatenation:
		r := make(Concatenation, len(l))
		copy(r, l)

		var d int
		for _, e := range r[1:] {
			d += e.Gutter()
			switch e.(type) {
			case Atom, Concatenation:
			default:
				break
			}
		}
		r[0] = gutterAligned(r[0], g-d)
		return r
	default:
		panic(l)
	}
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
		indent += l.Offset()
	}
	return nil
}

func (j Juxtaposition) Offset() int {
	var o int
	for i, l := range j {
		if i != 0 {
			o++
		}
		o += l.Offset()
	}
	return o
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
		indent += l.Offset()
	}
	return nil
}

func (c Concatenation) Offset() int {
	var o int
	for _, l := range c {
		o += l.Offset()
	}
	return o
}

func (c Concatenation) Gutter() int {
	var g int
	for _, l := range c {
		switch l := l.(type) {
		case Atom, Concatenation:
			g += l.Gutter()
		default:
			return g + l.Gutter()
		}
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
