package main

import (
	"fmt"
	"io/ioutil"
	"lexer"
	"os"
)

// go:generate go build -o="debug.exe" -ldflags="-w" -gcflags="-N -l"
func main() {
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic("file not found")
	}

	for m := range lexer.Lex(string(b)).Tokens {
		fmt.Println(m)
	}
}
