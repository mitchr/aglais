package lexer

import (
	"log"
	"unicode"
)

type TokenType int

//go:generate stringer -type=TokenType
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
	Tokens []Token

	state    stateFn
	input    []byte
	start    int
	position int
}

type stateFn func(*Lexer) stateFn

// Pushes a Token onto the Tokens channel
func (l *Lexer) push(t TokenType) {
	l.Tokens = append(l.Tokens, Token{t, string(l.input[l.start:l.position])})
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

func Lex(input []byte) *Lexer {
	l := &Lexer{
		state: lexAny,
		input: input,
	}

	for l.state != nil {
		l.state = l.state(l)
	}

	return l
}

// lexAny acts as a switchboard for all stateFn's
// first character is consumed
func lexAny(l *Lexer) stateFn {
	switch r := l.next(); {
	case isAlphaNumeric(r):
		// check for hexcharacters
		if p := l.next(); (r == '0' && (p == 'x' || p == 'X')) && isHexChar(l.peek()) {
			return lexHex
		}
		l.backup()
		return lexIdentifier
	case isOperator(r):
		if r == '.' {
			return lexDecimal
		}
		l.push(Operator)
		return lexAny
	case isQuoteChar(r):
		return lexQuote
	// this is ugly
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
	case unicode.IsDigit(r):
		if r == '0' {
			if c := l.next(); c == 'x' || c == 'X' {
				return lexHex
			}
		}
		l.backup()
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
		return nil
	default:
		log.Fatalf("Unrecognized character: %s\n", string(r))
	}
	log.Fatalf("Don't know what went wrong!\nlast rune scanned: %s\nposition:%v", string(l.input[l.position]), l.input[l.start:l.position])
	return nil
}

// has already captured the first character
func lexIdentifier(l *Lexer) stateFn {
	for isAlphaNumeric(l.peek()) {
		l.next()
	}

	for _, v := range l.input[l.start:l.position] {
		// if it's not a decimal character, then it must be an identifier
		if !isDecChar(rune(v)) {
			l.push(Identifier)
			return lexAny
		}
	}
	return lexDecimal
}

// first quoteChar has been consumed
// this functions still sucks oh my lord
func lexQuote(l *Lexer) stateFn {
	for r := l.next(); r != eof; r = l.next() {
		switch r {
		// the next character is an escaped one
		case '\\':
			// if l.next() is a character literal, don't follow the escape logic
			if isEscapeDelim(l.next()) {
				return lexQuote
			}
			// if it wasn't a characterDelim, backup so we're on the correct char
			l.backup()
			// backup so we're currently on the '\' char
			l.backup()
			// skip over the \ (remove it from the input entirely)
			l.input = append(l.input[:l.position], l.input[l.position+1:]...)
			return lexAny
		// r is the next '"' found
		case '"':
			if l.next() == '"' {
				// we have found the first part of a TriQuote, now we just need to skip the inside part
				for !isQuoteChar(l.peek()) {
					l.next()
				}
				if l.next() == '"' && l.next() == '"' && l.next() == '"' {
					l.push(TriQuote)
					return lexAny
				} else {
					log.Fatal("End quotes missing from TriQuote")
				}
			} else {
				l.backup()
				l.push(MonoQuote)
				return lexAny
			}
		}
	}
	l.push(MonoQuote)
	return nil
}

// terminator is included in the Comment token
func lexComment(l *Lexer) stateFn {
	for r := l.next(); r != eof; r = l.next() {
		if r == '\n' {
			l.push(Comment)
			return lexAny
		}
	}
	// push whatever we have before exiting
	l.push(Comment)
	return nil
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
	for isHexChar(l.peek()) {
		l.next()
	}
	l.push(HexNumber)
	return lexAny
}

func lexDecimal(l *Lexer) stateFn {
	for isDecChar(l.peek()) {
		l.next()
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
	return c == 'a' || c == 'b' || c == 'c' || c == 'd' || c == 'e' || c == 'f' || unicode.IsDigit(c)
}

func isDecChar(c rune) bool {
	return unicode.IsDigit(c) || c == 'e' || c == '.' || c == '-'
}
func isQuoteChar(c rune) bool {
	return c == '"' || c == '\''
}

func isEscapeDelim(c rune) bool {
	return c == 'n' || c == 'f' || c == 't' || c == 'r' || c == 'v'
}
