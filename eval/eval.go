package eval

import (
	"github.com/EclesioMeloJunior/monkey-lang/ast"
	"github.com/EclesioMeloJunior/monkey-lang/object"
)

func Eval(node ast.Node) object.Representation {
	switch node := node.(type) {

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

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
