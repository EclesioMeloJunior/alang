package parser

import (
	"github.com/EclesioMeloJunior/monkey-lang/ast"
	"github.com/EclesioMeloJunior/monkey-lang/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	// load the next token
	p.nextToken()

	// TODO: We're skiping the expressions until we encounter
	//  a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
