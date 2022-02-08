package parser

import (
	"fmt"

	"github.com/EclesioMeloJunior/ducklang/ast"
	"github.com/EclesioMeloJunior/ducklang/lexer"
	"github.com/EclesioMeloJunior/ducklang/token"
)

type (
	prefixParserFn func() ast.Expression

	// the argument represents the left part of the infix expression that's being parsed
	infixParserFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []error

	prefixParsers map[token.TokenType]prefixParserFn
	infixParsers  map[token.TokenType]infixParserFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	p.prefixParsers = make(map[token.TokenType]prefixParserFn)

	p.addPrefixParserFn(token.IDENT, p.parseIdentifier)
	p.addPrefixParserFn(token.INT, p.parseIntegerLiteral)
	p.addPrefixParserFn(token.TRUE, p.parseBooleanLiteral)
	p.addPrefixParserFn(token.FALSE, p.parseBooleanLiteral)
	p.addPrefixParserFn(token.BANG, p.parsePrefixExpression)
	p.addPrefixParserFn(token.MINUS, p.parsePrefixExpression)
	p.addPrefixParserFn(token.LPAREN, p.parseGroupedExpression)
	p.addPrefixParserFn(token.IF, p.parseIfExpression)
	p.addPrefixParserFn(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParsers = make(map[token.TokenType]infixParserFn)
	p.addInfixParserFn(token.PLUS, p.parseInfixExpression)
	p.addInfixParserFn(token.MINUS, p.parseInfixExpression)
	p.addInfixParserFn(token.SLASH, p.parseInfixExpression)
	p.addInfixParserFn(token.ASTHERISC, p.parseInfixExpression)
	p.addInfixParserFn(token.EQ, p.parseInfixExpression)
	p.addInfixParserFn(token.NOT_EQ, p.parseInfixExpression)
	p.addInfixParserFn(token.LT, p.parseInfixExpression)
	p.addInfixParserFn(token.GT, p.parseInfixExpression)
	p.addInfixParserFn(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// expectPeek advances the parser cursor to the next token if
// the given `t` is equals the next token, otherwise returns false
// and add an error to the parser errors field
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

func (p *Parser) addPrefixParserFn(token token.TokenType, f prefixParserFn) {
	p.prefixParsers[token] = f
}

func (p *Parser) addInfixParserFn(token token.TokenType, f infixParserFn) {
	p.infixParsers[token] = f
}
