package lexer

import "github.com/EclesioMeloJunior/monkey-lang/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
		return
	}

	l.char = l.input[l.readPosition]

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.char {
	case '=':
		tok = newToken(token.ASSIGN, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '!':
		tok = newToken(token.BANG, l.char)
	case '*':
		tok = newToken(token.ASTHERISC, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '<':
		tok = newToken(token.LT, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupLiteralType(tok.Literal)

			return tok
		}

		if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()

			return tok
		}

		tok = newToken(token.ILLEGAL, l.char)
	}

	l.readChar()
	return tok
}

// readIdentifier returns the name that belongs to a variable/function and is not a
// allowed a keyword, eg. `let name = "eclesio"` the name is the identifier
func (l *Lexer) readIdentifier() (ident string) {
	identStarts := l.position

	// read until the current character is not a letter
	// this allows things like: `let name="..."` or `let name = "..."`
	for isLetter(l.char) {
		l.readChar()
	}

	identEnds := l.position

	return l.input[identStarts:identEnds]
}

func (l *Lexer) readNumber() string {
	numberStarts := l.position

	for isDigit(l.char) {
		l.readChar()
	}

	numberEnds := l.position

	return l.input[numberStarts:numberEnds]
}

func (l *Lexer) skipWhitespace() {
	for isToIgnore(l.char) {
		l.readChar()
	}
}

func newToken(tokType token.TokenType, char byte) token.Token {
	return token.Token{
		Type:    tokType,
		Literal: string(char),
	}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' ||
		char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isToIgnore(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}
