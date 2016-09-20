package parser

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Mitchell-Riley/aglais/lexer"
)

func TestParseFile(t *testing.T) {
	t.Skip()
	b, err := ioutil.ReadFile("../test.io")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(Parse(lexer.Lex(b)))
}
