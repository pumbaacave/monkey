package evaluator

import (
	"monkey/object"
	"testing"
)

func TestExpadMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`let exp = macro(){quote(1 + 2); };
			exp();`,
			`(1 + 2)`,
		},
		{
			`let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
			reverse(2 + 2, 10 - 5);`,
			`(10 - 5) - (2 + 2)`,
		},
		{
			// macro must be enclosd in quote()
			`
			let unless = macro(condition, consequence, alternative) {
				quote(if (!(unquote(condition))) {
					unquote(consequence);
				} else {
					unquote(alternative);
				});
			};

			unless(10 > 5, puts("not greater"), puts("greater"));
			`,
			`if (!(10 > 5)) {puts("not greater"); } else { puts("greater"); }`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.String() != expected.String() {
			t.Errorf("not equal. want=%q, got=%q", expected.String(), expanded.String())
		}
	}
}
