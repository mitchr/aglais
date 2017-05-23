package lexer

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLexFile(t *testing.T) {
	b, err := ioutil.ReadFile("../test.io")
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
		[]byte(`name_with_underscores`),
	}

	for _, v := range b {
		for _, m := range Lex(v).Tokens {
			if m.Type != Identifier {
				fmt.Println(m, len(m.Value))
				t.Fail()
			}
		}
	}
}

func TestLexQuote(t *testing.T) {
	b := [][]byte{
		// monoquotes
		[]byte(`'apostrophe quotes'`),
		[]byte(`"monoquote"`),
		[]byte(`"monoquote with whitespace and
		newlines"`),
		[]byte(`'quote with operators _=-=-()'`),

		// small quote
		[]byte(`"p"`),

		// triples
		[]byte(`"""triplequote"""`),
		[]byte(`"""triplequote with
		newline and
		whitespace characters"""`),

		// escaped
		[]byte(`"\'"`),
		[]byte(`'\"'`),
		[]byte(`"\""`),
		[]byte(`"\"\"\""`),
		[]byte(`"\r"`),
		[]byte(`"\n\t"`),
		[]byte(`"\r\f\""`),
		[]byte(`"""\r\f"""`),

		// empty
		[]byte(`""`),
		[]byte(`""""""`),

		// []byte(`"""triquote with missing end quotes""`),
		// []byte(`"monoquote with trailing operator",`),
	}

	for _, v := range b {
		for _, m := range Lex(v).Tokens {
			switch m.Type {
			case TriQuote, MonoQuote:
				continue
			default:
				fmt.Println(m, len(m.Value))
				t.Fail()
			}
		}
	}
}

func TestLexComment(t *testing.T) {
	b := [][]byte{
		[]byte(`/*star
		* comment
		 with newlines*/`),
		[]byte("//singlecomment"),
		[]byte("#hashcomment"),
	}
	for _, v := range b {
		k := Lex(v).Tokens
		for _, m := range k {
			if m.Type != Comment || len(k) > 1 {
				fmt.Println(m, len(m.Value))
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
		k := Lex(v).Tokens
		for _, m := range k {
			if m.Type != HexNumber || len(k) > 1 {
				fmt.Println(m, len(m.Value))
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
		k := Lex(v).Tokens
		for _, m := range k {
			if m.Type != Decimal || len(k) > 1 {
				fmt.Println(m, len(m.Value))
				t.Fail()
			}
		}
	}
}
