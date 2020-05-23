package sqlfmt

import (
	"fmt"
	"io"
	"strings"
)

type Layout interface {
	Write(io.Writer, int) error
	Width() int
}

type Atom string

func (a Atom) Write(w io.Writer, i int) error {
	_, err := fmt.Fprintf(w, "%s", a)
	return err
}

func (a Atom) Width() int {
	return len(a)
}

type Stack struct {
	Up   Layout
	Down Layout
}

func (s Stack) Write(w io.Writer, i int) error {
	if err := s.Up.Write(w, 0); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "\n%s", strings.Repeat(" ", i)); err != nil {
		return err
	}
	if err := s.Down.Write(w, i); err != nil {
		return err
	}
	return nil
}

func (s Stack) Width() int {
	return s.Down.Width()
}

type Juxtaposition struct {
	Concentrated bool
	Left         Layout
	Right        Layout
}

func (j Juxtaposition) Write(w io.Writer, i int) error {
	if err := j.Left.Write(w, i); err != nil {
		return err
	}
	if !j.Concentrated {
		i++
		if _, err := fmt.Fprint(w, " "); err != nil {
			return err
		}
	}
	if err := j.Right.Write(w, i+j.Left.Width()); err != nil {
		return err
	}
	return nil
}

func (j Juxtaposition) Width() int {
	w := j.Left.Width() + j.Right.Width()
	if !j.Concentrated {
		w++
	}
	return w
}

type Choice struct {
	One     Layout
	Another Layout
}
