package parser

import (
	"fmt"

	"github.com/EclesioMeloJunior/monkey-lang/ast"
	"github.com/EclesioMeloJunior/monkey-lang/lexer"
	"github.com/EclesioMeloJunior/monkey-lang/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// in the let statment, after the keyword `let`
	// the next token must be and IDENT
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// in the let statement, after the identifier
	// the next token must be the `=`
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: ignoring what comes after `=` for now.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// expectPeek advances the parser cursor to the next token if
// the given `t` is equals the next token, otherwise returns false
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.unexpectedTypeErr(t, p.peekToken.Type)
	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) unexpectedTypeErr(expected, got token.TokenType) {
	err := fmt.Errorf("expected next token type be %s. got type %s", expected, got)
	p.errors = append(p.errors, err)
}

// Errors return all errors faced by the parser
func (p *Parser) Errors() []error {
	return p.errors
}
