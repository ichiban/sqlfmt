package sqlfmt

import (
	"errors"
	"fmt"
)

func AlignGutter(l Layout) (Layout, error) {
	m := maxGutter(l)
	if m == 0 {
		return nil, errors.New("failed")
	}

	return applyGutter(l, m)
}

func maxGutter(l Layout) int {
	switch l := l.(type) {
	case Atom:
		return len(l)
	case Juxtaposition:
		return maxGutter(l.Left)
	case Stack:
		m, n := maxGutter(l.Up), maxGutter(l.Down)
		if m > n {
			return m
		} else {
			return n
		}
	default:
		return 0
	}
}

func applyGutter(l Layout, m int) (Layout, error) {
	switch l := l.(type) {
	case Juxtaposition:
		if a, ok := l.Left.(Atom); ok {
			l.Left = Atom(fmt.Sprintf("%*s", m, a))
			return l, nil
		} else {
			left, err := applyGutter(l.Left, m)
			if err != nil {
				return nil, err
			}
			l.Left = left
			return l, nil
		}
	case Stack:
		up, err := applyGutter(l.Up, m)
		if err != nil {
			return nil, err
		}
		down, err := applyGutter(l.Down, m)
		if err != nil {
			return nil, err
		}
		return Stack{
			Up:   up,
			Down: down,
		}, nil
	default:
		return nil, errors.New("failed")
	}
}
