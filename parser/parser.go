package parser

import (
	"fmt"
	"github.com/Mitchell-Riley/aglais/lexer"
)

type TokenType int

//go:generate stringer -type=TokenType
const (
	expression TokenType = iota
	message
	arguments
	argument
	symbol
)

type Node struct {
	Type TokenType
	// Val  string

	uToken lexer.Token //underlying Token type
}

type tree struct {
	Nodes []Node
}

type Parser struct {
	Tokens chan lexer.Token

	current lexer.Token
	Buff    []Node // unexport this when finished debugging
	state   stateFn
}

type stateFn func(*Parser) stateFn

func (p *Parser) next() lexer.Token {
	p.current = <-p.Tokens
	return p.current
}

func (p *Parser) push(t TokenType) {
	p.Buff = append(p.Buff, Node{t, lexer.Token{p.current.Type, p.current.Value}})
}

func Parse(l *lexer.Lexer) *Parser {
	tok := <-l.Tokens

	p := &Parser{
		Tokens:  l.Tokens,
		current: tok,
		state:   parseAny,
	}

	// go func() {
	for ; p.state != nil; p.state = p.state(p) {
	}
	// }()

	return p
}

func parseAny(p *Parser) stateFn {
	switch t := p.next(); t.Type {
	case lexer.Identifier, lexer.Operator, lexer.MonoQuote, lexer.TriQuote, lexer.Decimal, lexer.HexNumber:
		return parseExpression
	case lexer.Open:
		return parseArguments
	case lexer.Comment:
		return parseAny
	case -1:
		return nil
	default:
		fmt.Println("Unknown token:", t)
		return parseAny
	}
}

func parseExpression(p *Parser) stateFn {
	// for {
	// 	switch p.next().Type {
	// 	case lexer.Comment, lexer.Terminator:
	// 		p.push(expression)
	// 		return parseAny
	// 	}
	// }

	for t := p.next().Type; t != lexer.Comment || t != lexer.Terminator; {
		p.next()
	}
	p.push(expression)
	return parseAny
}

func parseArguments(p *Parser) stateFn {
	p.push(arguments)
	return parseAny
}

func parseArgument(p *Parser) stateFn {
	switch p.next().Type {
	case lexer.Comment:
		return parseExpression
	}
	return parseAny
}
