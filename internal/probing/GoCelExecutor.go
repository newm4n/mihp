package probing

import (
	"context"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/newm4n/mihp/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
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

	funcs := cel.Functions(
		&functions.Overload{
			Operator: "KeyExist_string_boolean",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to IsDefined", s1.Type())
				}
				if _, ok := celContext[s1.Value().(string)]; ok {
					return types.Bool(true)
				}
				return types.Bool(false)
			},
		},
		&functions.Overload{
			Operator: "GetString_string_string",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetString", s1.Type())
				}
				if strItv, ok := celContext[s1.Value().(string)]; ok {
					return types.String(strItv.(string))
				}
				return types.String("")
			},
		},
		&functions.Overload{
			Operator: "GetInt_string_int",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetInt", s1.Type())
				}
				if intItv, ok := celContext[s1.Value().(string)]; ok {
					return types.Int(intItv.(int))
				}
				return types.Int(0)
			},
		},
		&functions.Overload{
			Operator: "GetUint_string_uint",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetUint", s1.Type())
				}
				if uintItv, ok := celContext[s1.Value().(string)]; ok {
					return types.Uint(uintItv.(uint))
				}
				return types.Uint(0)
			},
		},
		&functions.Overload{
			Operator: "GetFloat_string_float",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetFloat", s1.Type())
				}
				if floatItv, ok := celContext[s1.Value().(string)]; ok {
					return types.Double(floatItv.(float64))
				}
				return types.Double(0)
			},
		},
		&functions.Overload{
			Operator: "GetBool_string_bool",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetBool", s1.Type())
				}
				if boolItv, ok := celContext[s1.Value().(string)]; ok {
					return types.Bool(boolItv.(bool))
				}
				return types.Bool(false)
			},
		},
		&functions.Overload{
			Operator: "GetTime_string_time",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetTime", s1.Type())
				}
				if timeItv, ok := celContext[s1.Value().(string)]; ok {
					return types.Timestamp{timeItv.(time.Time)}
				}
				return types.Timestamp{time.Now()}
			},
		},
		&functions.Overload{
			Operator: "GetDuration_string_duration",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetDuration", s1.Type())
				}
				if durItv, ok := celContext[s1.Value().(string)]; ok {
					return types.Duration{durItv.(time.Duration)}
				}
				return types.Duration{0}
			},
		},
	)

	prg, err := env.Program(ast, funcs)
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
