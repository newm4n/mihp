package helper

import "reflect"

const (
	BaseKindInt BaseKind = iota
	BaseKindUint
	BaseKindFloat
	BaseKindBool
	BaseKindString
	BaseKindTime
	BaseKindDuration
	BaseKindArray
	BaseKindMap
	BaseKindOther
)

type BaseKind int

func GetBaseKindOfType(typ reflect.Type) BaseKind {
	switch typ.Kind() {
	case reflect.Bool:
		return BaseKindBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return BaseKindInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return BaseKindUint
	case reflect.Float32, reflect.Float64:
		return BaseKindFloat
	case reflect.String:
		return BaseKindString
	case reflect.Array, reflect.Slice:
		return BaseKindArray
	case reflect.Map:
		return BaseKindMap
	case reflect.Struct:
		if typ.String() == "time.Time" {
			return BaseKindTime
		}
		if typ.String() == "time.Deadline" {
			return BaseKindDuration
		}
		return BaseKindOther
	default:
		return BaseKindOther
	}
}
