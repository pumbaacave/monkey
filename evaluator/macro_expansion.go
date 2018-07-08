package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, s := range program.Statements {
		if isMarcoDefinition(s) {
			addMacro(s, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func isMarcoDefinition(node ast.Statement) bool {
	letS, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letS.Value.(*ast.MacroLiteral)
	if !ok {
		return false
	}

	return true
}

func addMacro(s ast.Statement, env *object.Environment) {
	letS, _ := s.(*ast.LetStatement)
	macroLiteral, _ := letS.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}

	env.Set(letS.Name.Value, macro)
}
