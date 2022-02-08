package eval

import (
	"fmt"

	"github.com/EclesioMeloJunior/ducklang/ast"
	"github.com/EclesioMeloJunior/ducklang/object"
	"github.com/EclesioMeloJunior/ducklang/token"
)

var (
	Null  *object.Null    = &object.Null{}
	True  *object.Boolean = &object.Boolean{Value: true}
	False *object.Boolean = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Env) object.Representation {
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
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		switch node.Operator {
		case token.BANG:
			return evalBangPrefixOperatorExpression(right)
		case token.MINUS:
			return evalMinusPrefixOperatorExpression(right)
		default:
			return errorF("unknow operator: %s%s", node.Operator, right.Type())
		}

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.ReturnStatement:
		returned := Eval(node.Value, env)
		return &object.Return{
			Value: returned,
		}

	case *ast.LetStatement:
		valueToBind := Eval(node.Value, env)
		if isError(valueToBind) {
			return valueToBind
		}

		env.Set(node.Name.Value, valueToBind)
		return nil

	case *ast.Identifier:
		stored, has := env.Get(node.Value)
		if !has {
			return errorF("identifier not found: %s", node.Value)
		}

		return stored

	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	default:
		return nil
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Env) object.Representation {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		// cannot evaluate since the Expression is not valid
		return condition
	}

	if condition.Type() != object.BOOLEAN_OBJ {
		return errorF("condition must evaluate to a boolean, got=%s", condition.Type())
	}

	switch condition {
	case Null:
		if node.Alternative != nil {
			return Eval(node.Alternative, env)
		}
		return Null

	case False:
		if node.Alternative != nil {
			return Eval(node.Alternative, env)
		}
		return Null

	case True:
		return Eval(node.Consequence, env)
	default:
		return Eval(node.Consequence, env)
	}
}

func evalBlockStatements(stmts []ast.Statement, env *object.Env) object.Representation {
	var rep object.Representation

	for _, stmt := range stmts {
		rep = Eval(stmt, env)

		switch rep := rep.(type) {
		case *object.Return:
			return rep
		case *object.Error:
			return rep
		}
	}

	return rep
}

func evalProgram(stmts []ast.Statement, env *object.Env) object.Representation {
	var rep object.Representation

	for _, stmt := range stmts {
		rep = Eval(stmt, env)

		switch rep := rep.(type) {
		case *object.Return:
			return rep.Value
		case *object.Error:
			return rep
		}
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
		return errorF("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(op string, left, right object.Representation) object.Representation {

	switch l := left.(type) {
	case *object.Integer:

		switch r := right.(type) {
		case *object.Integer:
			return evalIntegerInfixExpression(op, l, r)
		default:
			return errorF("type mismatch: %s %s %s", left.Type(), op, right.Type())
		}

	case *object.Boolean:

		switch r := right.(type) {
		case *object.Boolean:
			return evalBooleanInfixExpression(op, l, r)
		default:
			return errorF("type mismatch: %s %s %s", left.Type(), op, right.Type())
		}

	default:
		return errorF("unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return errorF("unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return errorF("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func isError(r object.Representation) bool {
	switch r.(type) {
	case *object.Error:
		return true
	default:
		return false
	}
}

func errorF(format string, o ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, o...),
	}
}
