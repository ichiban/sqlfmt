package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ichiban/sqlfmt"
)

func main() {
	for _, p := range os.Args[1:] {
		if err := format(p); err != nil {
			panic(err)
		}
	}
}

func format(s string) error {
	f, err := open(s)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	p := sqlfmt.NewParser(string(b))

	l, err := p.DirectSQLStatement()
	if err != nil {
		return err
	}

	if err := l.Write(os.Stdout, 0); err != nil {
		return err
	}

	if _, err := fmt.Println(); err != nil {
		return err
	}

	return nil
}

func open(s string) (io.ReadCloser, error) {
	if s == "-" {
		return ioutil.NopCloser(os.Stdin), nil
	}
	return os.Open(s)
}
