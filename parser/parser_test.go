package parser_test

import (
	"fmt"
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

func TestIdentifierExpression(t *testing.T) {
	const input = "foobar;"

	l := lexer.New(input)
	p := parser.New(l)

	prog := p.ParseProgram()
	checkParserErrors(t, p)

	if len(prog.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(prog.Statements))
	}

	stmt := prog.Statements[0]

	identifierStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] should be *ast.ExpressionStatement. got=%T", stmt)
	}

	identifierExpression, ok := identifierStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression should be an *ast.Identifier. got=%T", identifierExpression)
	}

	const declaredVariableName = "foobar"
	if identifierExpression.Value != declaredVariableName {
		t.Fatalf("expected token value=%q. got=%q", declaredVariableName, identifierExpression.Value)
	}

	if identifierExpression.TokenLiteral() != declaredVariableName {
		t.Fatalf("expected token literal=%q. got=%q", declaredVariableName, identifierExpression.Value)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	const input = "5;"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	expression, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt => expected a *ast.ExpressionStatement. got=%T", stmt)
	}

	if expression.Token.Type != token.INT {
		t.Fatalf("expression.Token.Type => expected token type INT. got=%s", expression.Token.Type)
	}

	intLiteral, ok := expression.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression.Expression => expected *ast.IntegerLiteral. got=%T", expression.Expression)
	}

	if intLiteral.TokenLiteral() != "5" {
		t.Fatalf(`intLiteral.Value => expected "5". got=%s`, expression.TokenLiteral())
	}

	if intLiteral.Value != 5 {
		t.Fatalf("intLiteral.Value => expected 5. got=%d", intLiteral.Value)
	}
}

func TestParsingPrefixOperator(t *testing.T) {
	prefixTests := [...]struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement => expected *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression => expected *ast.PrefixExpression. got=%T",
				stmt.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expected operator %s. got=%s", tt.operator, expression.Operator)
		}

		testIntegerLiteral(t, expression.Right, tt.integerValue)
	}
}

func testIntegerLiteral(t *testing.T, rightExp ast.Expression, expected int64) {
	integerExpression, ok := rightExp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expected right expression be *ast.IntegerLiteral. got=%T", rightExp)
	}

	if integerExpression.Token.Type != token.INT {
		t.Errorf("expected %s. got=%s", token.INT, integerExpression.Token.Type)
	}

	if integerExpression.Value != expected {
		t.Errorf("expected %d. got=%d", expected, integerExpression.Value)
	}

	if integerExpression.TokenLiteral() != fmt.Sprintf("%d", expected) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", expected, integerExpression.TokenLiteral())
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	testcases := []struct {
		input string

		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range testcases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		prog := p.ParseProgram()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("expected 1 statement. got=%d", len(prog.Statements))
		}

		stmt := prog.Statements[0]
		expressionStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement. got=%T", stmt)
		}

		infixExpression, ok := expressionStmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expected *ast.InfixExpression. got=%T", expressionStmt.Expression)
		}

		testIntegerLiteral(t, infixExpression.Left, tt.leftValue)
		if infixExpression.Operator != tt.operator {
			t.Fatalf("expected operator %s. got=%s", tt.operator, infixExpression.Operator)
		}
		testIntegerLiteral(t, infixExpression.Right, tt.rightValue)
	}
}

func TestParsingLongInfixExpression(t *testing.T) {
	testcases := []struct {
		input              string
		expectedExpression string
	}{

		{
			"-a * b;",
			"((-a) * b)",
		},
		{
			"!-a;",
			"(!(-a))",
		},
		{
			"a + b + c;",
			"((a + b) + c)",
		},
		{
			"a + b - c;",
			"((a + b) - c)",
		},
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"5 + 2 * 3;",
			"(5 + (2 * 3))",
		},
		{
			"3 + 4; -5 * 5;",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range testcases {
		l := lexer.New(tt.input)
		p := parser.New(l)

		prog := p.ParseProgram()
		checkParserErrors(t, p)

		if prog.String() != tt.expectedExpression {
			t.Fatalf("expected expression string %s. got=%s",
				tt.expectedExpression, prog.String())
		}
	}
}
