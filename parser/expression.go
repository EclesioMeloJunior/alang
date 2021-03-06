package parser

import (
	"fmt"
	"strconv"

	"github.com/EclesioMeloJunior/alang/ast"
	"github.com/EclesioMeloJunior/alang/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunc(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NOT_EQ:    EQUALS,
	token.LT:        LESS_GREATER,
	token.GT:        LESS_GREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTHERISC: PRODUCT,
	token.LPAREN:    CALL,
}

func (p *Parser) peekPrecedence() int {
	prec, ok := precedences[p.peekToken.Type]
	if ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	prec, ok := precedences[p.curToken.Type]
	if ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParsers[p.curToken.Type]
	if prefixFn == nil {
		p.errors = append(p.errors,
			fmt.Errorf("no prefix parser found for %s found", p.curToken.Type))
		return nil
	}

	leftExp := prefixFn()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	intValue, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors,
			fmt.Errorf("cannot parse %s to int64", p.curToken.Literal))
		return nil
	}

	lit.Value = intValue
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	prefixExp.Right = p.parseExpression(PREFIX)
	return prefixExp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fnLiteral := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fnLiteral.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fnLiteral.Body = p.parseBlockStatement()

	return fnLiteral
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}

	p.nextToken()

	identifier := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	parameters = append(parameters, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // goes to comma token as it is the peek token

		if !p.expectPeek(token.IDENT) { // after the comma must exists an identifier
			return nil
		}

		identifier = &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		parameters = append(parameters, identifier)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{
		Token:    p.curToken,
		Function: left,
	}
	callExpression.Arguments = p.parseCallArguments()
	return callExpression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	firstArgument := p.parseExpression(LOWEST)
	args = append(args, firstArgument)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // advance to token.COMMA
		p.nextToken() // advance after token.COMMA to evaluate as an expression

		nextArgument := p.parseExpression(LOWEST)
		args = append(args, nextArgument)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
