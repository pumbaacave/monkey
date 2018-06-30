package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 999;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program does not contain 3 statement, go to %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let'. get = %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not an *ast.Statement. got = %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got = '%s'", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got = '%s'", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParsrErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 11111;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program does not contain 3 statement, go to %d", len(program.Statements))
	}

	for _, s := range program.Statements {
		rs, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("s not an *ast.Statement. got = %T", s)
		}

		if rs.TokenLiteral() != "return" {
			t.Fatalf("s.TokenLiteral is not 'let'. get = %q", s.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program not enough statement. got=%d", len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	id, ok := s.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *asp.Expression. get=%T", s.Expression)
	}
	if id.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", id.TokenLiteral())
	}
	if id.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", id.TokenLiteral())
	}

}
func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program not enough statement. got=%d", len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	il, ok := s.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *asp.Expression. get=%T", s.Expression)
	}
	if il.Value != 5 {
		t.Errorf("il.Value not %s. got=%s", "foobar", il.TokenLiteral())
	}
	if il.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", il.TokenLiteral())
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParsrErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program not enough statement. got=%d", len(program.Statements))
		}
		s, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := s.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *asp.Expression. get=%T", s.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("operater not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	i, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if i.Value != value {
		t.Errorf("i.Value not %d,. got=%s", value, i.TokenLiteral())
		return false
	}

	if i.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("i.TokenLiteral() not %d. got=%s", value,
			i.TokenLiteral())
		return false
	}

	return true
}
