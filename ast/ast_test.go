package ast_test

import (
	"testing"

	"github.com/EclesioMeloJunior/monkey-lang/ast"
	"github.com/EclesioMeloJunior/monkey-lang/token"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "my_var",
					},
					Value: "my_var",
				},
				Value: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "another_var",
					},
					Value: "another_var",
				},
			},
		},
	}

	program_str := program.String()
	const expected_prog_str = "let my_var = another_var;"

	if program_str != expected_prog_str {
		t.Errorf("program.String() wrong.\n\texpected = %q\n\tgot = %q", expected_prog_str, program_str)
	}
}
