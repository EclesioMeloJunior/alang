package parser_test

import (
	"fmt"
	"testing"

	"github.com/EclesioMeloJunior/monkey-lang/ast"
	"github.com/EclesioMeloJunior/monkey-lang/lexer"
	"github.com/EclesioMeloJunior/monkey-lang/parser"
	"github.com/EclesioMeloJunior/monkey-lang/token"
)

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	identifier, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got=%T", exp)
	}

	if identifier.Value != value {
		t.Fatalf("expected identifier %s. got=%s", value, identifier.Value)
	}

	if identifier.TokenLiteral() != value {
		t.Fatalf("expected token literal %s. got=%s",
			value, identifier.TokenLiteral())
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int, int64:
		testIntegerLiteral(t, exp, v.(int64))
	case string:
		testIdentifier(t, exp, v)
	default:
		t.Fatalf("type of expected %T not handled", expected)
	}
}

func testInfixExpression(t *testing.T, expression ast.Expression,
	left, right interface{}, operator string) {
	infixExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("expected *ast.InfixExpression. got=%T", expression)
	}

	testLiteralExpression(t, infixExpression.Left, left)

	if infixExpression.Operator != operator {
		t.Fatalf("expected %s operator. got=%s",
			operator, infixExpression.Operator)
	}

	testLiteralExpression(t, infixExpression.Right, right)
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

	testIdentifier(t, identifierStmt.Expression, "foobar")
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

	testIntegerLiteral(t, expression.Expression, 5)
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
		t.Fatalf("expected right expression be *ast.IntegerLiteral. got=%T", rightExp)
	}

	if integerExpression.Token.Type != token.INT {
		t.Fatalf("expected %s. got=%s", token.INT, integerExpression.Token.Type)
	}

	if integerExpression.Value != expected {
		t.Fatalf("expected %d. got=%d", expected, integerExpression.Value)
	}

	if integerExpression.TokenLiteral() != fmt.Sprintf("%d", expected) {
		t.Fatalf("integer.TokenLiteral not %d. got=%s", expected, integerExpression.TokenLiteral())
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

		testInfixExpression(t, expressionStmt.Expression, tt.leftValue, tt.rightValue, tt.operator)
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
