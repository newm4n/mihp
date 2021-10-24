package probing

import (
	"context"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/newm4n/mihp/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
)

func GoCelEvaluate(ctx context.Context, expression string, celContext ProbeContext, expectReturnKind reflect.Kind) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if len(expression) == 0 {
		return nil, nil
	}
	declarations, err := celContext.Declarations()
	if err != nil {
		logrus.Errorf("error while creating declaration got %s", err)
		return nil, err
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	env, err := cel.NewEnv(
		cel.Declarations(declarations...))
	if err != nil {
		logrus.Errorf("error while creating go-cel environment got %s", err)
		return nil, err
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		logrus.Errorf("error while compiling expression [%s] got %s", expression, issues.Err())
		return nil, fmt.Errorf("%w : %s", issues.Err(), issues.String())
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	prg, err := env.Program(ast)
	if err != nil {
		logrus.Errorf("error while creating program for expression [%s] got %s", expression, err)
		return nil, issues.Err()
	}
	toEval := make(map[string]interface{})
	for k, v := range celContext {
		toEval[k] = v
	}
	out, _, err := prg.Eval(toEval)
	if err != nil {
		logrus.Errorf("error while valuating program for expression [%s] got %s", expression, err)
		return nil, issues.Err()
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if reflect.TypeOf(out.Value()).Kind() != expectReturnKind {
		return nil, fmt.Errorf("%w : expression \"%s\" expect returns %s but %s", errors.ErrEvalReturnInvalid, expression, expectReturnKind.String(), reflect.TypeOf(out.Value()).Kind().String())
	}
	return out.Value(), nil
}
