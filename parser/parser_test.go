package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParsrErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() return nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain 1 statement, go to %d", len(program.Statements))
		}
		s := program.Statements[0]
		if !testLetStatement(t, s, tt.expectedIdentifier) {
			return
		}

		val := s.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	id, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not *ast.Identifier. got=%T", exp)
		return false
	}

	if id.Value != value {
		t.Errorf("id.Value not %s. got=%s", value, id.Value)
		return false
	}

	if id.TokenLiteral() != value {
		t.Errorf("id.TokenLiteral not %s. got=%s", value, id.TokenLiteral())
		return false
	}

	return true
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
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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
			t.Fatalf("exp not *asp.PrefixExpression. get=%T", s.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("operater not %s. got=%s", tt.operator, exp.Operator)
		}

		testLiteralExpression(t, exp.Right, tt.value)
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s', got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	prefixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		exp, ok := s.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp not *asp.InfixExpression. get=%T", s.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("operater not %s. got=%s", tt.operator, exp.Operator)
		}

		if !testInfixExpression(t, s.Expression, tt.leftValue, tt.operator, tt.rightValue) {
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b -c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{
			" a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])),(b[1]),(2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParsrErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected = %q, got = %q", tt.expected, actual)
		}
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not * ast.Bool. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t, got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statement. got=%d", 1, len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := s.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *asp.IfExpression. get=%T", s.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequense is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statement. got=%d", 1, len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := s.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *asp.IfExpression. get=%T", s.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequense is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequense.Statement[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statement[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `fn(x, y) { x + y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statement. got=%d", 1, len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	fun, ok := s.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("exp not *asp.FunctionLiteral. get=%T", s.Expression)
	}

	if len(fun.Parameters) != 2 {
		t.Fatalf("function literal param wrong. what 2, get=%d, %v", len(fun.Parameters), fun.Parameters)
	}

	testLiteralExpression(t, fun.Parameters[0], "x")
	testLiteralExpression(t, fun.Parameters[1], "y")
	if len(fun.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements hat not 1 statements. got=%d", len(fun.Body.Statements))
	}

	b, ok := fun.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			fun.Body.Statements[0])
	}

	testInfixExpression(t, b.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn(){}", expectedParams: []string{}},
		{input: "fn(x){}", expectedParams: []string{"x"}},
		{input: "fn(x, y, z){}", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParsrErrors(t, p)

		s := program.Statements[0].(*ast.ExpressionStatement)
		fun := s.Expression.(*ast.FunctionLiteral)

		if len(fun.Parameters) != len(tt.expectedParams) {
			t.Errorf("length of params wrong. want %d, got=%d",
				len(tt.expectedParams), len(fun.Parameters))
		}

		for i, id := range tt.expectedParams {
			testLiteralExpression(t, fun.Parameters[i], id)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statement. got=%d", 1, len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ce, ok := s.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp not *asp.CallExpression. get=%T", s.Expression)
	}

	if len(ce.Arguments) != 3 {
		t.Fatalf("function literal param wrong. what 3, get=%d, %v", len(ce.Arguments), ce.Arguments)
	}

	testLiteralExpression(t, ce.Arguments[0], 1)
	testInfixExpression(t, ce.Arguments[1], 2, "*", 3)
	testInfixExpression(t, ce.Arguments[2], 4, "+", 5)
}

func TestStingLiteralExpression(t *testing.T) {
	input := `"Hello World!";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	s := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := s.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", s.Expression)
	}

	if literal.Value != "Hello World!" {
		t.Fatalf("literal.Value not %q. got=%q", "Hello World", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := s.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", s.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", s.Expression)
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressioins(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	index, ok := s.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", s.Expression)
	}
	if !testIdentifier(t, index.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, index.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	s := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := s.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", s.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}

}

func TestParsingEmptyHashLiterals(t *testing.T) {
	input := "{}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	s := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := s.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", s.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 -8, "three": 15 /5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	s := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := s.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", s.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one":   func(e ast.Expression) { testInfixExpression(t, e, 0, "+", 1) },
		"two":   func(e ast.Expression) { testInfixExpression(t, e, 10, "-", 8) },
		"three": func(e ast.Expression) { testInfixExpression(t, e, 15, "/", 5) },
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key%q found", literal.String())
			continue
		}
		testFunc(value)
	}

}

func TestMacroLiteralExpression(t *testing.T) {
	input := `macro(x, y) { x + y ;}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsrErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statement. got=%d", 1, len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prgram.Statement[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	macro, ok := s.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("exp not *asp.MacroLiteral. get=%T", s.Expression)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("macro literal param wrong. what 2, get=%d, %v", len(macro.Parameters), macro.Parameters)
	}

	testLiteralExpression(t, macro.Parameters[0], "x")
	testLiteralExpression(t, macro.Parameters[1], "y")
	if len(macro.Body.Statements) != 1 {
		t.Fatalf("macro.Body.Statements hat not 1 statements. got=%d", len(macro.Body.Statements))
	}

	b, ok := macro.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("macro body stmt is not ast.ExpressionStatement. got=%T",
			macro.Body.Statements[0])
	}

	testInfixExpression(t, b.Expression, "x", "+", "y")
}
