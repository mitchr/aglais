package parser

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Mitchell-Riley/aglais/lexer"
)

func TestParseFile(t *testing.T) {
	files := []string{
		// "../test.io",
		"../fizzbuzz.io",
		// "../grammar",
	}

	for _, v := range files {
		b, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(Parse(lexer.Lex(b)))
	}
}
