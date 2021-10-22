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
)

func NewProbeContext() ProbeContext {
	return make(ProbeContext)
}

type ProbeContext map[string]interface{}

func (pctx ProbeContext) String() string {
	var buff *bytes.Buffer
	table := tablewriter.NewWriter(buff)
	table.SetHeader([]string{"NO", "KEY", "VALUE"})

	keys := make([]string, 0)
	for k := range pctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for k, v := range keys {
		table.Append([]string{strconv.Itoa(k + 1), v, fmt.Sprintf("%v", pctx[v])})
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
			return nil, fmt.Errorf("context value of key \"%s\" not supported : %s", k, vType.Name())
		}
		declarations = append(declarations, decls.NewVar(k, typ))
	}
	return declarations, nil
}
