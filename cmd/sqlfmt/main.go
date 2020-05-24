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
		s, err := sql(p)
		if err != nil {
			panic(err)
		}

		f, err := sqlfmt.Format(s)
		if err != nil {
			panic(err)
		}

		if _, err = fmt.Print(f); err != nil {
			panic(err)
		}
	}
}

func sql(p string) (string, error) {
	var r io.Reader
	if p == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(p)
		if err != nil {
			return "", err
		}
		defer f.Close()
		r = f
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}
