package eval

import (
	"github.com/EclesioMeloJunior/monkey-lang/ast"
	"github.com/EclesioMeloJunior/monkey-lang/object"
	"github.com/EclesioMeloJunior/monkey-lang/token"
)

var (
	Null  *object.Null    = &object.Null{}
	True  *object.Boolean = &object.Boolean{Value: true}
	False *object.Boolean = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Representation {
	switch node := node.(type) {

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BooleanLiteral:
		// avoid to create new instances
		// every time we encounter a bool
		if node.Value {
			return True
		}

		return False

	case *ast.PrefixExpression:
		right := Eval(node.Right)

		switch node.Operator {
		case token.BANG:
			return evalBangPrefixOperatorExpression(right)
		case token.MINUS:
			return evalMinusPrefixOperatorExpression(right)
		default:
			return Null
		}

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)

		return evalInfixExpression(node.Operator, left, right)

	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	default:
		return nil
	}
}

func evalStatements(stmts []ast.Statement) object.Representation {
	var rep object.Representation

	for _, stmt := range stmts {
		rep = Eval(stmt)
	}

	return rep
}

func evalBangPrefixOperatorExpression(right object.Representation) object.Representation {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right object.Representation) object.Representation {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{
			Value: -right.Value,
		}
	default:
		return Null
	}
}

func evalInfixExpression(op string, left, right object.Representation) object.Representation {
	switch l := left.(type) {
	case *object.Integer:

		switch r := right.(type) {
		case *object.Integer:
			return evalIntegerInfixExpression(op, l, r)
		default:
			return Null
		}

	case *object.Boolean:

		switch r := right.(type) {
		case *object.Boolean:
			return evalBooleanInfixExpression(op, l, r)
		default:
			return Null
		}

	default:
		return Null
	}

}

func evalIntegerInfixExpression(op string, left, right *object.Integer) object.Representation {
	switch op {
	case token.PLUS:
		return &object.Integer{
			Value: left.Value + right.Value,
		}
	case token.MINUS:
		return &object.Integer{
			Value: left.Value - right.Value,
		}
	case token.ASTHERISC:
		return &object.Integer{
			Value: left.Value * right.Value,
		}
	case token.SLASH:
		return &object.Integer{
			Value: int64(left.Value / right.Value),
		}
	case token.GT:
		if left.Value > right.Value {
			return True
		}
		return False
	case token.LT:
		if left.Value < right.Value {
			return True
		}
		return False
	case token.NOT_EQ:
		if left.Value != right.Value {
			return True
		}
		return False
	case token.EQ:
		if left.Value == right.Value {
			return True
		}
		return False
	default:
		return Null
	}
}

func evalBooleanInfixExpression(op string, left, right *object.Boolean) object.Representation {
	switch op {
	case token.NOT_EQ:
		if left.Value != right.Value {
			return True
		}
		return False
	case token.EQ:
		if left.Value == right.Value {
			return True
		}
		return False
	default:
		return Null
	}
}
