all:  wasm_exec.js sqlfmt.wasm

wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/$@" .

sqlfmt.wasm: ../go.mod ../go.sum $(wildcard ../*.go) ../cmd/sqlfmtwasm/main.go
	GOOS=js GOARCH=wasm go build -o $@ github.com/ichiban/sqlfmt/cmd/sqlfmtwasm
