package lexer_test

import (
	"testing"

	"github.com/EclesioMeloJunior/ducklang/lexer"
	"github.com/EclesioMeloJunior/ducklang/token"
)

func Test_BasicTokens_NextToken(t *testing.T) {
	input := "=+(){},;!-/*5<>"

	tests := []struct {
		exepextedType   token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTHERISC, "*"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.GT, ">"},
	}

	l := lexer.New(input)

	for idx, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.exepextedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				idx, tt.exepextedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong. expected=%q, got=%q",
				idx, tt.expectedLiteral, tok.Literal)
		}
	}
}

func Test_ProgTokens_NextToken(t *testing.T) {
	const prog = `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
}

let return = add(five, ten);

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
5 != 6;
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LET, "let"},
		{token.RETURN, "return"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := lexer.New(prog)

	for idx, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				idx, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong. expected=%q, got=%q",
				idx, tt.expectedLiteral, tok.Literal)
		}
	}
}
