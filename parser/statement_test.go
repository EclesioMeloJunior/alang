package parser_test

import (
	"strings"
	"testing"

	"github.com/EclesioMeloJunior/monkey-lang/ast"
	"github.com/EclesioMeloJunior/monkey-lang/lexer"
	"github.com/EclesioMeloJunior/monkey-lang/parser"
	"github.com/EclesioMeloJunior/monkey-lang/token"
)

func TestLetStatement(t *testing.T) {
	const prog = `
let x = 5;
let y = 10;
let foobar = 8989;	
`
	l := lexer.New(prog)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("parser.ParseProgram return nil, expected not nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contains 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"}, {"y"}, {"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
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
	const prog = `
return 5;
return 10;
return 1010101;`

	l := lexer.New(prog)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contains 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.Token.Type != token.RETURN {
			t.Errorf("returnStmt.TokenLiteral not return. got=%q",
				returnStmt.TokenLiteral())
		}

		const returnLiteral = "return"
		if returnStmt.TokenLiteral() != returnLiteral {
			t.Errorf("returnStmt.TokenLiteral not return. got=%q",
				returnStmt.TokenLiteral())
		}
	}
}
