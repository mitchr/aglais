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

	for _, m := range Lex(b).Tokens {
		fmt.Println(m, len(m.Value))
	}
}

func TestLexIdentifier(t *testing.T) {
	b := [][]byte{
		[]byte(`o setSlot`),
	}

	for _, v := range b {
		for _, m := range Lex(v).Tokens {
			if m.Type.String() != "Identifier" {
				fmt.Println(m, len(m.Value))
				t.Fail()
			} else {
				fmt.Println(m, len(m.Value))
			}
		}
	}
}

func TestLexQuote(t *testing.T) {
	t.Skip()
	b := [][]byte{
		[]byte(`'apostrophe quotes'`),
		[]byte(`"monoquote"`),
		[]byte(`"monoquote with whitespace and
		newlines"`),
		[]byte(`"""triplequote"""`),
		[]byte(`"""triplequote with
		newline and
		whitespace characters"""`),
		// escaped
		[]byte(`"\""`),
		[]byte(`"\"\"\""`),
		[]byte(`""`),
		[]byte(`""""""`),
	}

	for _, v := range b {
		for _, m := range Lex(v).Tokens {
			if g := m.Type.String(); g != "MonoQuote" || g != "TriQuote" {
				t.Fail()
			} else {
				fmt.Println(m, len(m.Value))
			}
		}
	}
}

func TestLexComment(t *testing.T) {
	b := [][]byte{
		[]byte(`/*star
		* comment
		 yeah*/`),
		[]byte("//singlecomment"),
		[]byte("#hashcomment"),
	}
	for _, v := range b {
		for _, m := range Lex(v).Tokens {
			if m.Type.String() != "Comment" {
				t.Fail()
			}
		}
	}
}

func TestLexHex(t *testing.T) {
	b := [][]byte{
		[]byte(`0x89`),
		[]byte(`0x3a67`),
		[]byte(`0X7afe78`),
	}

	// this should fail
	// c := []string{
	// 	`9xae`,
	// }

	for _, v := range b {
		for _, m := range Lex(v).Tokens {
			if m.Type.String() != "HexNumber" {
				t.Fail()
			}
		}
	}
}

func TestLexDecimal(t *testing.T) {
	b := [][]byte{
		[]byte(`54`),
		[]byte(`7.80`),
		[]byte(`3.141592654`),
		[]byte(`0.21`),
		[]byte(`5e3`),
		[]byte(`4e-65`),
	}

	for _, v := range b {
		for _, m := range Lex(v).Tokens {
			if m.Type.String() != "Decimal" {
				t.Fail()
			}
		}
	}
}
