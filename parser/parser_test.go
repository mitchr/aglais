package parser

import (
	"fmt"
	"github.com/Mitchell-Riley/aglais/lexer"
	"io/ioutil"
	"testing"
)

func TestParseFile(t *testing.T) {
	t.Skip()
	b, err := ioutil.ReadFile("..\\test.io")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(Parse(lexer.Lex(string(b))).Buff)

	// fmt.Printf("%#v\n",
}
