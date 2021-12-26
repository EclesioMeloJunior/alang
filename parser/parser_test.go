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
