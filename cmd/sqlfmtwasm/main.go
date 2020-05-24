package main

import (
	"syscall/js"

	"github.com/ichiban/sqlfmt"
)

func main() {
	c := make(chan struct{})
	js.Global().Set("formatSql", js.FuncOf(formatSQL))
	<-c
}

func formatSQL(_ js.Value, args []js.Value) interface{} {
	f, err := sqlfmt.Format(args[0].String())
	if err != nil {
		panic(err)
	}
	return js.ValueOf(f)
}
