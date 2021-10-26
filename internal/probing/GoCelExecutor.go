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

func GoCelEvaluate(ctx context.Context, expression string, celContext ProbeContext, expectReturnKind reflect.Kind) (output interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			err = fmt.Errorf("%w : panic occured during evaluating expression [%s] : got %v", errors.ErrEvalError, expression, err)
			switch expectReturnKind {
			case reflect.String:
				output = ""
			case reflect.Int:
				output = int(0)
			case reflect.Int8:
				output = int8(0)
			case reflect.Int16:
				output = int16(0)
			case reflect.Int32:
				output = int32(0)
			case reflect.Int64:
				output = int64(0)
			case reflect.Uint:
				output = uint(0)
			case reflect.Uint8:
				output = uint8(0)
			case reflect.Uint16:
				output = uint16(0)
			case reflect.Uint32:
				output = uint32(0)
			case reflect.Uint64:
				output = uint64(0)
			case reflect.Float32:
				output = float32(0)
			case reflect.Float64:
				output = float64(0)
			case reflect.Bool:
				output = false
			default:
				output = nil
			}
		}
	}()

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
		&functions.Overload{
			Operator: "GetLength_string_int",
			Unary: func(value ref.Val) ref.Val {
				s1, ok := value.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetLength", s1.Type())
				}
				if durItv, ok := celContext[s1.Value().(string)]; ok {
					val := reflect.ValueOf(durItv)
					if val.Type().Kind() == reflect.Slice || val.Type().Kind() == reflect.Array {
						return types.Int(val.Len())
					}
				}
				return types.Int(0)
			},
		},
		&functions.Overload{
			Operator: "GetStringElem_string_int_string",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetStringElem 1st Argument", s1.Type())
				}
				s2, ok := rhs.(types.Int)
				if !ok {
					return types.ValOrErr(s2, "unexpected type '%v' passed to GetStringElem 2nd Argument", s2.Type())
				}
				if strArrItv, ok := celContext[s1.Value().(string)]; ok {
					strArrTyp := reflect.TypeOf(strArrItv)
					strArrVal := reflect.ValueOf(strArrItv)
					if strArrTyp.Elem().Kind() == reflect.String {
						return types.String(strArrVal.Index(int(s2.Value().(int64))).String())
					}
				}
				return types.String("")
			},
		},

		&functions.Overload{
			Operator: "GetIntElem_string_int_int",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetIntElem 1st Argument", s1.Type())
				}
				s2, ok := rhs.(types.Int)
				if !ok {
					return types.ValOrErr(s2, "unexpected type '%v' passed to GetIntElem 2nd Argument", s2.Type())
				}
				if strArrItv, ok := celContext[s1.Value().(string)]; ok {
					strArrTyp := reflect.TypeOf(strArrItv)
					strArrVal := reflect.ValueOf(strArrItv)
					switch strArrTyp.Elem().Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						return types.Int(strArrVal.Index(int(s2.Value().(int64))).Int())
					}
				}
				return types.Int(0)
			},
		},
		&functions.Overload{
			Operator: "GetUintElem_string_int_uint",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetUintElem 1st Argument", s1.Type())
				}
				s2, ok := rhs.(types.Int)
				if !ok {
					return types.ValOrErr(s2, "unexpected type '%v' passed to GetUintElem 2nd Argument", s2.Type())
				}
				if strArrItv, ok := celContext[s1.Value().(string)]; ok {
					strArrTyp := reflect.TypeOf(strArrItv)
					strArrVal := reflect.ValueOf(strArrItv)
					switch strArrTyp.Elem().Kind() {
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						return types.Uint(strArrVal.Index(int(s2.Value().(int64))).Uint())
					}
				}
				return types.Uint(0)
			},
		},
		&functions.Overload{
			Operator: "GetFloatElem_string_int_float",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetFloatElem 1st Argument", s1.Type())
				}
				s2, ok := rhs.(types.Int)
				if !ok {
					return types.ValOrErr(s2, "unexpected type '%v' passed to GetFloatElem 2nd Argument", s2.Type())
				}
				if strArrItv, ok := celContext[s1.Value().(string)]; ok {
					strArrTyp := reflect.TypeOf(strArrItv)
					strArrVal := reflect.ValueOf(strArrItv)
					switch strArrTyp.Elem().Kind() {
					case reflect.Float32, reflect.Float64:
						return types.Double(strArrVal.Index(int(s2.Value().(int64))).Float())
					}
				}
				return types.Double(0)
			},
		},
		&functions.Overload{
			Operator: "GetBoolElem_string_int_bool",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetBoolElem 1st Argument", s1.Type())
				}
				s2, ok := rhs.(types.Int)
				if !ok {
					return types.ValOrErr(s2, "unexpected type '%v' passed to GetBoolElem 2nd Argument", s2.Type())
				}
				if strArrItv, ok := celContext[s1.Value().(string)]; ok {
					strArrTyp := reflect.TypeOf(strArrItv)
					strArrVal := reflect.ValueOf(strArrItv)
					if strArrTyp.Elem().Kind() == reflect.Bool {
						return types.Bool(strArrVal.Index(int(s2.Value().(int64))).Bool())
					}
				}
				return types.Bool(false)
			},
		},
		&functions.Overload{
			Operator: "GetTimeElem_string_int_time",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetTimeElem 1st Argument", s1.Type())
				}
				s2, ok := rhs.(types.Int)
				if !ok {
					return types.ValOrErr(s2, "unexpected type '%v' passed to GetTimeElem 2nd Argument", s2.Type())
				}
				if strArrItv, ok := celContext[s1.Value().(string)]; ok {
					strArrTyp := reflect.TypeOf(strArrItv)
					strArrVal := reflect.ValueOf(strArrItv)
					if strArrTyp.Elem().Kind() == reflect.Struct && strArrTyp.Elem().String() == "time.Time" {
						return types.Timestamp{strArrVal.Index(int(s2.Value().(int64))).Interface().(time.Time)}
					}
				}
				return types.Timestamp{time.Now()}
			},
		},
		&functions.Overload{
			Operator: "GetDurationElem_string_int_duration",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.ValOrErr(s1, "unexpected type '%v' passed to GetDurationElem 1st Argument", s1.Type())
				}
				s2, ok := rhs.(types.Int)
				if !ok {
					return types.ValOrErr(s2, "unexpected type '%v' passed to GetDurationElem 2nd Argument", s2.Type())
				}
				if strArrItv, ok := celContext[s1.Value().(string)]; ok {
					strArrTyp := reflect.TypeOf(strArrItv)
					strArrVal := reflect.ValueOf(strArrItv)
					if strArrTyp.Elem().Kind() == reflect.Struct && strArrTyp.Elem().String() == "time.Duration" {
						return types.Duration{strArrVal.Index(int(s2.Value().(int64))).Interface().(time.Duration)}
					}
				}
				return types.Duration{time.Duration(0)}
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
