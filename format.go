package sqlfmt

import (
	"bytes"
	"fmt"
)

func Format(sql string) (string, error) {
	var b bytes.Buffer
	p := NewParser(sql)

	l, err := p.DirectSQLStatement()
	if err != nil {
		return "", err
	}

	if err := l.Write(&b, 0); err != nil {
		return "", err
	}

	if _, err := fmt.Fprintln(&b); err != nil {
		return "", err
	}

	return b.String(), nil
}
