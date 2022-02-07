package eval_test

import (
	"testing"

	"github.com/EclesioMeloJunior/monkey-lang/eval"
	"github.com/EclesioMeloJunior/monkey-lang/lexer"
	"github.com/EclesioMeloJunior/monkey-lang/object"
	"github.com/EclesioMeloJunior/monkey-lang/parser"
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
		{"if (1) { 10 }", 10},
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

func testEval(input string) object.Representation {
	l := lexer.New(input)
	p := parser.New(l)

	prog := p.ParseProgram()
	return eval.Eval(prog)
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
