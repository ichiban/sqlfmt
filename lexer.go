package sqlfmt

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input  string
	tokens chan token

	pos   int
	width int
}

func NewLexer(input string) *Lexer {
	tokens := make(chan token)

	l := Lexer{
		input:  input,
		tokens: tokens,
	}

	go func() {
		defer close(tokens)

		s := l.start()
		for s != nil {
			s = s(l.next())
		}
	}()

	return &l
}

func (l *Lexer) Next() token {
	t := <-l.tokens
	return t
}

func (l *Lexer) next() (rune, int) {
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	pos := l.pos
	l.pos += l.width
	return r, pos
}

func (l *Lexer) peek() rune {
	r, _ := l.next()
	l.backup()
	return r
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) emit(t token) {
	l.tokens <- t
}

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	return fmt.Sprintf("<%s %v>", t.typ, t.val)
}

type tokenType int

// https://jakewheat.github.io/sql-overview/sql-2011-foundation-grammar.html

const (
	eos tokenType = iota
	errToken
	identifier
	unsignedNumeric
	characterString
	keyword
	leftParen
	rightParen
	asterisk
	plusSign
	comma
	minusSign
	period
	semicolon
	lessThanOperator
	equalsOperator
	greaterThanOperator
)

func (t tokenType) String() string {
	switch t {
	case eos:
		return "eos"
	case errToken:
		return "error"
	case identifier:
		return "identifier"
	case unsignedNumeric:
		return "unsigned numeric"
	case characterString:
		return "character string"
	case keyword:
		return "keyword"
	case leftParen:
		return "left paren"
	case rightParen:
		return "right paren"
	case asterisk:
		return "asterisk"
	case plusSign:
		return "plus sign"
	case comma:
		return "comma"
	case minusSign:
		return "minus sign"
	case period:
		return "period"
	case semicolon:
		return "semicolon"
	case lessThanOperator:
		return "less than operator"
	case equalsOperator:
		return "equals operator"
	case greaterThanOperator:
		return "greater than operator"
	default:
		return "unknown"
	}
}

type state func(rune, int) state

func (l *Lexer) start() state {
	return func(r rune, pos int) state {
		switch {
		case r == 0xfffd: // replacement character
			return nil
		case unicode.IsSpace(r):
			return l.start()
		case unicode.IsLetter(r), r == '_':
			l.backup()
			return l.regularIdent(pos)
		case unicode.IsDigit(r):
			return l.unsignedNumericLiteral(pos)
		case r == '\'':
			return l.characterStringLiteral(pos)
		case unicode.IsPunct(r), unicode.IsSymbol(r):
			l.backup()
			return l.specialChar()
		default:
			return nil
		}
	}
}

func (l *Lexer) regularIdent(start int) state {
	return func(r rune, pos int) state {
		switch {
		case unicode.IsLetter(r), r == '_', pos != start && unicode.IsNumber(r):
			return l.regularIdent(start)
		default:
			l.backup()
			val := l.input[start:pos]
			u := strings.ToUpper(val)
			if _, ok := keywords[u]; ok {
				l.emit(token{
					typ: keyword,
					val: u,
				})
				return l.start()
			}
			l.emit(token{
				typ: identifier,
				val: val,
			})
			return l.start()
		}
	}
}

func (l *Lexer) unsignedNumericLiteral(start int) state {
	return func(r rune, pos int) state {
		switch {
		case unicode.IsDigit(r):
			return l.unsignedNumericLiteral(start)
		case r == '.':
			return l.unsignedFloatLiteral(start)
		default:
			l.backup()
			l.emit(token{
				typ: unsignedNumeric,
				val: l.input[start:pos],
			})
			return l.start()
		}
	}
}

func (l *Lexer) unsignedFloatLiteral(start int) state {
	return func(r rune, pos int) state {
		switch {
		case unicode.IsDigit(r):
			return l.unsignedFloatLiteral(start)
		default:
			l.backup()
			l.emit(token{
				typ: unsignedNumeric,
				val: l.input[start:pos],
			})
			return l.start()
		}
	}
}

func (l *Lexer) characterStringLiteral(start int) state {
	return func(r rune, pos int) state {
		switch r {
		case '\'':
			return l.quoteSymbol(start)
		default:
			return l.characterStringLiteral(start)
		}
	}
}

func (l *Lexer) quoteSymbol(start int) state {
	return func(r rune, pos int) state {
		switch r {
		case '\'':
			return l.characterStringLiteral(start)
		default:
			l.backup()
			l.emit(token{
				typ: characterString,
				val: l.input[start:pos],
			})
			return l.start()
		}
	}
}

func (l *Lexer) specialChar() state {
	return func(r rune, pos int) state {
		v := l.input[pos : pos+1]
		switch r {
		case '(':
			l.emit(token{typ: leftParen, val: v})
			return l.start()
		case ')':
			l.emit(token{typ: rightParen, val: v})
			return l.start()
		case '*':
			l.emit(token{typ: asterisk, val: v})
			return l.start()
		case '+':
			l.emit(token{typ: plusSign, val: v})
			return l.start()
		case ',':
			l.emit(token{typ: comma, val: v})
			return l.start()
		case '-':
			l.emit(token{typ: minusSign, val: v})
			return l.start()
		case '.':
			l.emit(token{typ: period, val: v})
			return l.start()
		case ';':
			l.emit(token{typ: semicolon, val: v})
			return l.start()
		case '<':
			l.emit(token{typ: lessThanOperator, val: v})
			return l.start()
		case '=':
			l.emit(token{typ: equalsOperator, val: v})
			return l.start()
		case '>':
			l.emit(token{typ: greaterThanOperator, val: v})
			return l.start()
		default:
			l.emit(token{typ: errToken})
			return nil
		}
	}
}
