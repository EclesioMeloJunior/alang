package object

import "fmt"

var (
	_ Representation = (*Integer)(nil)
	_ Representation = (*Boolean)(nil)
)

type Type string

const (
	INTEGER_OBJ Type = "INTEGER"
	BOOLEAN_OBJ Type = "BOOLEAN"
	NULL_OBJ    Type = "NULL"
)

type Representation interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.value)
}
func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

type Boolean struct {
	value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.value)
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
