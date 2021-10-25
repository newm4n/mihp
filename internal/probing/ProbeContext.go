package probing

import (
	"bytes"
	"fmt"
	"github.com/google/cel-go/checker/decls"
	"github.com/newm4n/mihp/pkg/helper"
	"github.com/olekukonko/tablewriter"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewProbeContext() ProbeContext {
	return make(ProbeContext)
}

type ProbeContext map[string]interface{}

func (pctx ProbeContext) String() string {
	return pctx.ToString(true)
}

func ToPrint(val reflect.Value) string {
	switch val.Type().Kind() {
	case reflect.String:
		return fmt.Sprintf("\"%s\"", strings.Replace(val.String(), `"`, `\"`, -1))
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", val.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", val.Float())
	case reflect.Bool:
		return fmt.Sprintf("%v", val.Bool())
	case reflect.Array, reflect.Slice:
		ell := make([]string, val.Len())
		for i := 0; i < len(ell); i++ {
			ell[i] = ToPrint(val.Index(i))
		}
		return fmt.Sprintf("[%s]", strings.Join(ell, ","))
	case reflect.Map:
		ell := make([]string, 0)
		for _, key := range val.MapKeys() {
			valueVal := val.MapIndex(key)
			ell = append(ell, fmt.Sprintf("%s:%s", ToPrint(key), ToPrint(valueVal)))
		}
		return fmt.Sprintf("{%s}", strings.Join(ell, ","))
	case reflect.Struct:
		switch val.Type().String() {
		case "time.Time":
			t := val.Interface().(time.Time)
			return fmt.Sprintf("\"%s\"", t.String())
		case "time.Duration":
			d := val.Interface().(time.Duration)
			return fmt.Sprintf("\"%s\"", d.String())
		default:
			return fmt.Sprintf("unprintable_%s", val.Type().String())
		}
	default:
		return fmt.Sprintf("unprintable_%s", val.Type().String())
	}
}

func (pctx ProbeContext) ToString(short bool) string {
	var buff = &bytes.Buffer{}
	buff.WriteString("\n")
	table := tablewriter.NewWriter(buff)
	table.SetHeader([]string{"NO", "KEY", "VALUE"})

	keys := make([]string, 0)
	for k := range pctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for k, v := range keys {
		if err, ok := pctx[v].(error); ok {
			table.Append([]string{strconv.Itoa(k + 1), v, err.Error()})
		} else {
			toPrint := ToPrint(reflect.ValueOf(pctx[v]))
			if reflect.TypeOf(pctx[v]).Kind() == reflect.String && len(toPrint) > 20 && short {
				if len(toPrint) > 20 {
					toPrint = fmt.Sprintf("%s...(%d bytes more)", toPrint[:20], len(toPrint)-20)
				}
			}
			table.Append([]string{strconv.Itoa(k + 1), v, toPrint})
		}
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return buff.String()
}

func (pctx ProbeContext) Declarations() ([]*exprpb.Decl, error) {
	declarations := make([]*exprpb.Decl, 0)
	for k, v := range pctx {
		vType := reflect.TypeOf(v)
		var typ *exprpb.Type
		switch helper.GetBaseKindOfType(vType) {
		case helper.BaseKindFloat:
			typ = decls.Double
		case helper.BaseKindBool:
			typ = decls.Bool
		case helper.BaseKindUint:
			typ = decls.Uint
		case helper.BaseKindInt:
			typ = decls.Int
		case helper.BaseKindString:
			typ = decls.String
		case helper.BaseKindTime:
			typ = decls.Timestamp
		default:
			continue
		}
		declarations = append(declarations, decls.NewVar(k, typ))
	}
	declarations = append(declarations, decls.NewFunction("IsDefined",
		decls.NewOverload("KeyExist_string_boolean",
			[]*exprpb.Type{decls.String},
			decls.Bool)))
	return declarations, nil
}
