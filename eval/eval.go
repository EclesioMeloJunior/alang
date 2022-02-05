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
			return evalBangOperatorExpression(right)
		default:
			return Null
		}

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

func evalBangOperatorExpression(right object.Representation) object.Representation {
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
