package parser

import (
	"log"

	"github.com/Mitchell-Riley/aglais/lexer"
)

var eof = lexer.Token{Type: -1}

type parser struct {
	lexer *lexer.Lexer

	current  lexer.Token
	position int
}

func (p *parser) next() lexer.Token {
	if p.position == len(p.lexer.Tokens) {
		p.current = eof
		return eof
	}

	p.current = p.lexer.Tokens[p.position]
	p.position++
	return p.current
}

func Parse(l *lexer.Lexer) *parser {
	p := &parser{
		lexer: l,
	}

	p.next()
	p.expression()

	return p
}

func (p *parser) accept(tok lexer.TokenType) bool {
	if p.current.Type == tok {
		p.next()
		return true
	}
	return false
}

func (p *parser) expect(tok lexer.TokenType) bool {
	if p.accept(tok) {
		return true
	}
	log.Fatalf("unexpected token: %s", tok)
	return false
}

// expression ::= { symbol [arguments] | [Terminator | Comment] }
func (p *parser) expression() {
	for {
		if isSymbol(p.current.Type) {
			// check if we are at eof?
			p.next()
			if p.accept(lexer.Open) {
				p.arguments()
			}
		} else if p.accept(lexer.Terminator) || p.accept(lexer.Comment) {
			continue
		} else {
			break
		}
	}
}

// arguments ::= Open [expression [ { Comma expression } ] ] Close
// Open Close
// Open expression Close
// Open expression {Comma expression} Close
func (p *parser) arguments() {
	// no expression inside parens
	if p.accept(lexer.Close) {
		return
	} else if isSymbol(p.current.Type) {
		p.expression()
		for p.accept(lexer.Comma) {
			p.expression()
		}
		// if p.accept(lexer.Comma) {
		// 	for {
		// 		p.expression()
		// 		if !p.accept(lexer.Comma) {
		// 			break
		// 		}
		// 	}
		// }
		p.expect(lexer.Close)
	} else {
		log.Fatal("unbalanced parens?")
	}
}

func isSymbol(t lexer.TokenType) bool {
	return t == lexer.Identifier || t == lexer.HexNumber || t == lexer.Decimal || t == lexer.Operator || t == lexer.MonoQuote || t == lexer.TriQuote || t == lexer.Terminator
}
