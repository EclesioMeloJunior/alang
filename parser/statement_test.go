package parser_test

import (
	"strings"
	"testing"

	"github.com/EclesioMeloJunior/alang/ast"
	"github.com/EclesioMeloJunior/alang/lexer"
	"github.com/EclesioMeloJunior/alang/parser"
	"github.com/EclesioMeloJunior/alang/token"
)

func TestLetStatements(t *testing.T) {
	testcases := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{
			`let x = 5;`, "x", 5,
		},
		{
			`let y = true;`, "y", true,
		},
		{
			`let foobar = y;`, "foobar", "y",
		},
	}

	for _, tt := range testcases {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectedIdentifier)

		letStatement := stmt.(*ast.LetStatement)
		testLiteralExpression(t, letStatement.Value, tt.expectedValue)
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedIdentifier string) bool {
	t.Helper()

	if stmt.TokenLiteral() != strings.ToLower(token.LET) {
		t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt is not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Token.Type != token.IDENT {
		t.Errorf("token should be of type IDENT. got=%s", letStmt.Token.Type)
		return false
	}

	if letStmt.Name.Value != expectedIdentifier {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", expectedIdentifier, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != expectedIdentifier {
		t.Errorf("s.Name not '%s'. got=%s", expectedIdentifier, letStmt.Name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()

	parserErrs := p.Errors()

	const colorReset = "\033[0m"

	if len(parserErrs) > 0 {
		const colorRed = "\033[31m"
		t.Logf("%s Parser Errors: %d %s\n", colorRed, len(parserErrs), colorReset)
	} else {
		const colorGreen = "\033[32m"
		t.Logf("%s Parser Errors: %d %s\n", colorGreen, len(parserErrs), colorReset)
	}

	for _, err := range parserErrs {
		t.Errorf("parser error: %q", err)
	}

	if len(parserErrs) > 0 {
		t.FailNow()
	}
}

func TestReturnStatement(t *testing.T) {
	testcases := []struct {
		input         string
		expectedValue interface{}
	}{
		{`return 5;`, 5},
		{`return y;`, "y"},
		{`return true;`, true},
	}

	for _, tt := range testcases {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement. got=%d",
				len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", program.Statements[0])
		}

		if returnStmt.Token.Type != token.RETURN {
			t.Errorf("returnStmt.TokenLiteral not return. got=%q",
				returnStmt.TokenLiteral())
		}

		testLiteralExpression(t, returnStmt.Value, tt.expectedValue)
	}
}
