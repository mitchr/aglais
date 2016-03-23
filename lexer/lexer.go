package lexer

import (
	"fmt"
	"unicode"
)

type TokenType int

//go:generate stringer -type=TokenType
// calling a print func on a struct calls the print func on all of it's
// fields; cool
const (
	Identifier TokenType = iota

	Operator

	MonoQuote
	TriQuote

	Terminator

	Comment
	HexNumber
	Decimal

	Comma
	Open
	Close

	eof = -1
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	Tokens chan Token

	state    stateFn
	input    string
	start    int
	position int
}

type stateFn func(*Lexer) stateFn

// Pushes a Token onto the Tokens channel
func (l *Lexer) push(t TokenType) {
	l.Tokens <- Token{t, l.input[l.start:l.position]}
	l.start = l.position
}

//next returns the next rune in the input
//calling this, even in a Println, advances the lexer
func (l *Lexer) next() rune {
	if l.position == len(l.input) {
		return eof
	}
	r := rune(l.input[l.position:][0])
	l.position++
	return r
}

//check if moving the lexer backwards will break the EOF condition
func (l *Lexer) peek() rune {
	r := l.next()

	if l.position == len(l.input) {
		return eof
	}

	l.backup()
	return r
}

func (l *Lexer) backup() {
	l.position--
}

func Lex(input string) *Lexer {
	l := &Lexer{
		Tokens: make(chan Token),
		state:  lexAny,
		input:  input,
	}

	go func() {
		for ; l.state != nil; l.state = l.state(l) {
		}
		close(l.Tokens)
	}()

	return l
}

// lexAny acts as a switchboard for all stateFn's
// first character is consumed
func lexAny(l *Lexer) stateFn {
	switch r := l.next(); {
	case isAlphaNumeric(r):
		//check for hexcharacters
		if p := l.next(); (p == 'x' || p == 'X') && isHexChar(l.next()) {
			return lexHex
		}
		l.backup()
		l.backup()
		return lexIdentifier
	case isOperator(r):
		if r == '.' {
			return lexDecimal
		}
		l.push(Operator)
		return lexAny
	case r == '"', r == '\'':
		//advance until another '"' is found
		for ; l.peek() == '"' || l.peek() == '\''; l.next() {
		}
		return lexQuote
	//this is ugly, but practical
	case isSeparator(r), r == '\r', r == '\n':
		if r == '\n' {
			l.push(Terminator)
			return lexAny
		} else if p := l.next(); p == ';' || isSeparator(p) {
			l.push(Terminator)
		} else {
			l.backup()
			l.start = l.position
		}
		return lexAny
	case isWhitespace(r):
		fmt.Println("Found whitespace!")
		l.next()
		return lexAny
	case r == '#', r == '/':
		if r == '#' {
			return lexComment
		}
		switch p := l.next(); {
		case p == '/':
			return lexComment
		case p == '*':
			return lexStarComment
		}
	case r == '0':
		if c := l.next(); c == 'x' || c == 'X' {
			return lexHex
		}
		l.backup()
		fallthrough
	case unicode.IsDigit(r):
		return lexDecimal
	case r == ',':
		l.push(Comma)
		return lexAny
	case r == '(', r == '[', r == '{':
		l.push(Open)
		return lexAny
	case r == ')', r == ']', r == '}':
		l.push(Close)
		return lexAny
	case r == eof:
		// fmt.Println("EOF reached!")
		return nil
	default:
		defer close(l.Tokens)
		panic(fmt.Sprintf("Unrecognized character: %s\n", string(r)))
	}
	panic(fmt.Sprintf("Don't know what went wrong!\nlast rune scanned: %s\nposition:%v", string(l.input[l.position]), l.input[l.start:l.position]))
}

// has already captured the first character
func lexIdentifier(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case !isAlphaNumeric(r):
			//this is a bottleneck; find a way to to this faster
			//iterate over current lexeme and check if it contains any characters that aren't numbers
			b := make([]bool, 0)
			for _, v := range l.input[l.start:l.position] {
				if isDecChar(v) {
					b = append(b, true)
				} else {
					b = append(b, false)
				}
			}
			for _, v := range b {
				if v == false {
					l.push(Identifier)
					return lexAny
				}
			}
			return lexDecimal
		}
	}
}

func lexQuote(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == '"':
			if l.next() == '"' && l.next() == '"' {
				l.push(TriQuote)
				return lexAny
			}
			//should we be backing up here?
			// l.backup()
			l.push(MonoQuote)
			return lexAny
		case r == '\'':
			l.push(MonoQuote)
			return lexAny
		}
		//here's why this didn't work: if you put the variable in the beginning of the for loop, it captures it so you just keep switching on the same rune over and over again
		// for l.next() != '"' {
		// 	l.next()
		// }
	}
}

//terminator is included in the Comment token
func lexComment(l *Lexer) stateFn {
	for {
		switch l.next() {
		case eof:
			//push whatever we have before exiting
			l.push(Comment)
			return nil
		case '\n':
			l.push(Comment)
			return lexAny
		}
	}
}

func lexStarComment(l *Lexer) stateFn {
	for {
		switch l.next() {
		case '*':
			if l.next() == '/' {
				l.push(Comment)
				return lexAny
			}
		}
	}
}

func lexHex(l *Lexer) stateFn {
	for ; isHexChar(l.peek()); l.next() {
	}
	l.push(HexNumber)
	return lexAny
}

func lexDecimal(l *Lexer) stateFn {
	for ; isDecChar(l.peek()); l.next() {
	}
	l.push(Decimal)
	return lexAny
}

func isAlphaNumeric(c rune) bool {
	return unicode.IsDigit(c) || unicode.IsLetter(c) || c == '_'
}

// '/' is an operator and the symbol for a comment; handle that here currently commented out
func isOperator(c rune) bool {
	// c == '/'
	// c == '\''
	return c == ':' || c == '.' || c == '~' || c == '!' || c == '@' || c == '$' || c == '%' || c == '^' || c == '&' || c == '*' || c == '-' || c == '+' || c == '=' || c == '{' || c == '}' || c == '[' || c == ']' || c == '|' || c == '\\' || c == '<' || c == '>' || c == '?'
}

func isSeparator(c rune) bool {
	return c == ' ' || c == '\f' || c == '\t' || c == '\v'
}

func isWhitespace(c rune) bool {
	return isSeparator(c) || c == '\r' || c == '\n'
}

func isHexChar(c rune) bool {
	return c == 'a' || c == 'b' || c == 'c' || c == 'd' || c == 'e' || unicode.IsDigit(c)
}

func isDecChar(c rune) bool {
	return unicode.IsDigit(c) || c == 'e' || c == '.' || c == '-'
}
