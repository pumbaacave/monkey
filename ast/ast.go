package ast

// Node basic term
type Node interface {
	TokenLiteral() string
}

// Statement describe action
type Statement interface {
	Node
	statementNode()
}

// Expression create some values
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
