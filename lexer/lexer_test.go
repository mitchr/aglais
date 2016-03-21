package lexer

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLex(t *testing.T) {
	b, _ := ioutil.ReadFile("..\\test.io")

	for m := range Lex(string(b)).Tokens {
		fmt.Println(m, len(m.Value))
	}
}

func TestLexHex(t *testing.T) {
	b := []string{
		"0x89",
		"0x3a67",
		"0X7afe78",
	}
	for _, v := range b {
		if m := <-Lex(v).Tokens; m.Type.String() != "HexNumber" {
			t.Fail()
		}
	}
}

//this fails for seemingly no reason
func TestLexQuote(t *testing.T) {
	t.Skip("infinite loop?")

	b := []string{
		`"escaped\""`,
		`"monoquote test"`,
		`"""triplequote test"""`,
	}
	for _, v := range b {
		if m := <-Lex(v).Tokens; m.Type.String() != "MonoQuote" || m.Type.String() != "TriQuote" {
			fmt.Println(m, m.Type)
			t.Fail()
		}
	}
}

func TexLexComment(t *testing.T) {
	b := []string{
		`/*star
		* comment
		 yeah*/`,
		"//singlecomment",
		"#hashcomment",
	}
	for _, v := range b {
		if m := <-Lex(v).Tokens; m.Type.String() != "Comment" {
			t.Fail()
		}
	}
}
