package eval_test

import (
	"testing"

	"github.com/EclesioMeloJunior/ducklang/eval"
	"github.com/EclesioMeloJunior/ducklang/lexer"
	"github.com/EclesioMeloJunior/ducklang/object"
	"github.com/EclesioMeloJunior/ducklang/parser"
)

func TestEvaluationLiteralObjects(t *testing.T) {
	testcases := []struct {
		input    string
		expected interface{}
	}{
		{"5;", 5},
		{"10;", 10},
		{"-5;", -5},
		{"-10;", -10},

		{"5 + 5 + 5 + 5 - 10;", 10},
		{"2 * 2 * 2 * 2 * 2;", 32},
		{"-50 + 100 + -50;", 0},
		{"5 * 2 + 10;", 20},
		{"5 + 2 * 10;", 25},
		{"20 + 2 * -10;", 0},
		{"50 / 2 * 2 + 10;", 60},
		{"2 * (5 + 10);", 30},
		{"3 * 3 * 3 + 10;", 37},
		{"3 * (3 * 3) + 10;", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},

		{"3 > 5;", false},
		{"1 < 2;", true},
		{"1 < 1;", false},
		{"1 < 1;", false},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"1 == 2;", false},
		{"1 != 2;", true},

		{"true;", true},
		{"false;", false},
		{"!true;", false},
		{"!false;", true},
		{"!!true;", true},
		{"!!false;", false},
		{"!5;", false},
		{"!!5;", true},

		{"true == true;", true},
		{"true == false;", false},
		{"true != false;", true},
		{"false == false;", true},
		{"false != false;", false},
		{"(1 < 2) == true;", true},
		{"(1 > 2) == true;", false},
		{"(1 > 2) == false;", true},
	}

	for _, tt := range testcases {
		evaluated := testEval(tt.input)
		testEvaluatedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestEvaluatesConditions(t *testing.T) {
	testcases := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", &object.Error{Message: "condition must evaluate to a boolean, got=INTEGER"}},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 10 + 10 }", 20},
		{"if (!5) { 10 } else { 2 + 3 }", 5},
	}

	for _, tt := range testcases {
		evaluated := testEval(tt.input)
		testEvaluatedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestEvaluatesReturns(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"7; return 2 * 5; 9;", 10},
		{`if (true) {
			if (true) {
				return 10;
			}
			return 1;
		}`, 10},
	}

	for _, tt := range tests {
		evalulated := testEval(tt.input)
		testEvaluatedObject(t, tt.input, evalulated, tt.expected)
	}
}

func TestErrorWhileEvaluating(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.Error
	}{
		{"-(10 + -(true + false));", &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"5 + true;", &object.Error{Message: "type mismatch: INTEGER + BOOLEAN"}},
		{"5 + true; 5;", &object.Error{Message: "type mismatch: INTEGER + BOOLEAN"}},
		{"-true;", &object.Error{Message: "unknown operator: -BOOLEAN"}},
		{"true + false;", &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"5; true + false; 5", &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"if( 10 > 1) { true + false }", &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"if( true > 1 ) { true + false }", &object.Error{Message: "type mismatch: BOOLEAN > INTEGER"}},
		{`if( 10 ) { true + false }`, &object.Error{Message: "condition must evaluate to a boolean, got=INTEGER"}},
		{`if (true) {
			if (true) {
				return true + false;
			}
			return 1;
		}`, &object.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"}},
		{"x;", &object.Error{Message: "identifier not found: x"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testEvaluatedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestEvalutaionLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let x = 5; x;`, 5},
		{`let x = (2 / 3) + 1; x;`, 1},
		{`let x = 10; let b = x; b;`, 10},
		{`let x = 5; let b = x * 5; b + 1;`, 26},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testEvaluatedObject(t, tt.input, evaluated, tt.expected)
	}
}

func testEval(input string) object.Representation {
	l := lexer.New(input)
	p := parser.New(l)

	prog := p.ParseProgram()
	return eval.Eval(prog, object.NewEnv())
}

func testEvaluatedObject(t *testing.T, input string, r object.Representation, expected interface{}) {
	switch exp := expected.(type) {
	case nil:
		testNullObject(t, input, r)
	case int:
		testIntegerObject(t, input, r, int64(exp))
	case int64:
		testIntegerObject(t, input, r, int64(exp))
	case bool:
		testBooleanObject(t, input, r, exp)
	case *object.Error:
		testErrorObject(t, input, r, exp)
	}
}

func testErrorObject(t *testing.T, input string, r object.Representation, expected *object.Error) {
	switch eval := r.(type) {
	case *object.Error:
		if eval.Message != expected.Message {
			t.Fatalf("%s\n\texpected=%s. got=%s",
				input, expected.Message, eval.Message)
		}
	default:
		t.Fatalf("%s\n\texpected *object.Error. got=%T (%+v)",
			input, eval, eval)
	}
}

func testIntegerObject(t *testing.T, input string, r object.Representation, expected int64) {
	result, ok := r.(*object.Integer)
	if !ok {
		t.Fatalf("%s\n\texpected *object.Integer. got=%T (%+v)", input, r, r)
	}

	if result.Value != expected {
		t.Fatalf("%s\n\texpected %d. got=%d", input, expected, result.Value)
	}
}

func testBooleanObject(t *testing.T, input string, r object.Representation, expected bool) {
	result, ok := r.(*object.Boolean)
	if !ok {
		t.Fatalf("%s\n\texpected *object.Boolean. got=%T (%+v)", input, r, r)
	}

	if result.Value != expected {
		t.Fatalf("%s\n\texpected %t. got=%t", input, expected, result.Value)
	}
}

func testNullObject(t *testing.T, input string, r object.Representation) {
	_, ok := r.(*object.Null)
	if !ok {
		t.Fatalf("%s\n\texpected *object.Null. got=%T (%+v)", input, r, r)
	}
}
