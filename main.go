package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/Mitchell-Riley/aglais/lexer"
	"github.com/Mitchell-Riley/aglais/parser"
)

// go:generate go build -o="debug.exe" -ldflags="-w" -gcflags="-N -l"
func main() {
	file := flag.Arg(0)

	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic("file not found")
	}

	for m := range lexer.Lex(b).Tokens {
		fmt.Println(m)
	}

	parser.Parse(lexer.Lex(b))

}
