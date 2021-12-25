package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	PLUS   = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"let": LET,
	"fn":  FUNCTION,
}

// LookupLiteralType receives a word as argument and check if
// it is a current language keyword and if true returns the right
// TokenType to the keyword, otherwise return the type IDENT for all
// user-defined identifiers
func LookupLiteralType(word string) TokenType {
	if ttype, ok := keywords[word]; ok {
		return ttype
	}

	return IDENT
}
