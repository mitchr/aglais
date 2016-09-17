package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Mitchell-Riley/aglais/lexer"
	"github.com/Mitchell-Riley/aglais/parser"
)

// go:generate go build -o="debug.exe" -ldflags="-w" -gcflags="-N -l"
func main() {
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic("file not found")
	}

	for m := range lexer.Lex(b).Tokens {
		fmt.Println(m)
	}

	parser.Parse(lexer.Lex(b))

}
