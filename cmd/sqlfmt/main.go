package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ichiban/sqlfmt"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	p := sqlfmt.NewParser(string(b))

	l, err := p.DirectSQLStatement()
	if err != nil {
		panic(err)
	}

	l, err = sqlfmt.AlignGutter(l)
	if err != nil {
		panic(err)
	}

	if err := l.Write(os.Stdout, 0); err != nil {
		panic(err)
	}

	if _, err := fmt.Println(); err != nil {
		panic(err)
	}
}
