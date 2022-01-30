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

func testBooleanLiteral(t *testing.T, exp ast.Expression, expected bool) {
	boolLiteral, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("expected *ast.BooleanLiteral. got=%T", exp)
	}

	if expected &&
		boolLiteral.Token.Type != token.TRUE &&
		boolLiteral.TokenLiteral() != token.TRUE {
		t.Fatalf("expected token TRUE. got=%s", boolLiteral.TokenLiteral())
	}

	if !expected &&
		boolLiteral.Token.Type != token.FALSE &&
		boolLiteral.TokenLiteral() != token.FALSE {
		t.Fatalf("expected token FALSE. got=%s", boolLiteral.TokenLiteral())
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case string:
		testIdentifier(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
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

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, value interface{}) {
	prefixExpression, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("expected *ast.PrefixExpression. got=%T", exp)
	}

	if prefixExpression.Operator != operator {
		t.Fatalf("expected operator %s. got=%s", operator, prefixExpression.Operator)
	}

	testLiteralExpression(t, prefixExpression.Right, value)
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

func TestBooleanLiteralExpression(t *testing.T) {
	testcases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "true;",
			expected: true,
		},
		{
			input:    "false;",
			expected: false,
		},
	}

	for _, tt := range testcases {
		l := lexer.New(tt.input)
		p := parser.New(l)

		prog := p.ParseProgram()

		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("expected 1 statement. got=%d", len(prog.Statements))
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement. got=%T", prog.Statements[0])
		}

		testLiteralExpression(t, stmt.Expression, tt.expected)
	}
}

func TestParsingPrefixOperator(t *testing.T) {
	prefixTests := [...]struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!false;", "!", false},
		{"!true;", "!", true},
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

		testPrefixExpression(t, stmt.Expression, tt.operator, tt.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	testcases := []struct {
		input string

		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"false == false;", false, "==", false},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
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

func TestParsingInfixExpression(t *testing.T) {
	testcases := []struct {
		input              string
		expectedExpression string
	}{
		{
			"true;",
			"true",
		},
		{
			"false;",
			"false",
		},
		{
			"3 > 5 == false;",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true;",
			"((3 < 5) == true)",
		},
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
