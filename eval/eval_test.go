package eval_test

import (
	"testing"

	"github.com/EclesioMeloJunior/monkey-lang/eval"
	"github.com/EclesioMeloJunior/monkey-lang/lexer"
	"github.com/EclesioMeloJunior/monkey-lang/object"
	"github.com/EclesioMeloJunior/monkey-lang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	testcases := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"10;", 10},
	}

	for _, tt := range testcases {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Representation {
	l := lexer.New(input)
	p := parser.New(l)

	prog := p.ParseProgram()
	return eval.Eval(prog)
}

func testIntegerObject(t *testing.T, r object.Representation, expected int64) {
	result, ok := r.(*object.Integer)
	if !ok {
		t.Fatalf("expected *object.Integer. got=%T (%+v)", r, r)
	}

	if result.Value != expected {
		t.Fatalf("expected %d. got=%d", expected, result.Value)
	}
}
