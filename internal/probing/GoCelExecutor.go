package probing

import (
	"github.com/google/cel-go/cel"
)

func GoCelEvaluate(expression string, celContext ProbeContext) (interface{}, error) {
	declarations, err := celContext.Declarations()
	if err != nil {
		return nil, err
	}
	env, err := cel.NewEnv(
		cel.Declarations(declarations...))
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}
	prg, err := env.Program(ast)
	if err != nil {
		return nil, issues.Err()
	}
	out, _, err := prg.Eval(celContext)
	if err != nil {
		return nil, issues.Err()
	}
	return out.Value(), nil
}
