package object

import (
	"fmt"
	"strings"

	"github.com/EclesioMeloJunior/alang/ast"
)

var (
	_ Representation = (*Integer)(nil)
	_ Representation = (*Boolean)(nil)
	_ Representation = (*Null)(nil)
	_ Representation = (*Error)(nil)
	_ Representation = (*Function)(nil)
)

type Type string

const (
	INTEGER_OBJ         Type = "INTEGER"
	BOOLEAN_OBJ         Type = "BOOLEAN"
	NULL_OBJ            Type = "NULL"
	RETURN_VALUE_OBJECT Type = "RETURN_VALUE"
	ERROR               Type = "ERROR"
	FUNCTION_OBJ             = "FUNCTION_OBJ"
)

type Representation interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) Type() Type {
	return NULL_OBJ
}

type Return struct {
	Value Representation
}

func (r *Return) Type() Type {
	return RETURN_VALUE_OBJECT
}
func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) Type() Type {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	params := make([]string, len(f.Parameters))
	for i, param := range f.Parameters {
		params[i] = param.String()
	}

	return fmt.Sprintf("fn(%s){...}", strings.Join(params, ", "))
}
