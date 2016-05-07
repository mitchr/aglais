package lexer

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLexFile(t *testing.T) {
	b, err := ioutil.ReadFile("..\\test.io")
	if err != nil {
		fmt.Println(err)
	}
	Lex(string(b))
}

func TestLexIdentifier(t *testing.T) {
	b := []string{
		`o setSlot`,
	}

	for _, v := range b {
		for m := range Lex(v).Tokens {
			if m.Type.String() != "Identifier" {
				t.Fail()
			} else {
				fmt.Println(m, len(m.Value))
			}
		}
	}
}

func TestLexQuote(t *testing.T) {
	b := []string{
		`'single quote'`,
		`"monoquote"`,
		`"""triplequote"""`,
		`"""triplequote with
		newline and 
		 whitespace characters"""`,
		// `"escaped\""`,
		// `""`,
		// `""""""`,
	}
	for _, v := range b {
		switch m := <-Lex(v).Tokens; m.Type.String() {
		case "MonoQuote", "TriQuote":
			fmt.Println(m)
		default:
			t.Fail()
		}
	}
}

func TestLexComment(t *testing.T) {
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

func TestLexHex(t *testing.T) {
	b := []string{
		`0x89`,
		`0x3a67`,
		`0X7afe78`,
	}
	for _, v := range b {
		if m := <-Lex(v).Tokens; m.Type.String() != "HexNumber" {
			t.Fail()
		}
	}
}

func TestLexDecimal(t *testing.T) {
	b := []string{
		`54`,
		`7.80`,
		`3.141592654`,
		`0.21`,
		`5e3`,
		`4e-65`,
	}
	for _, v := range b {
		if m := <-Lex(v).Tokens; m.Type.String() != "Decimal" {
			t.Fail()
		}
	}
}
