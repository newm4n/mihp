package probing

import (
	"context"
	"github.com/newm4n/mihp/internal"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestGoCelEvaluateExistBool(t *testing.T) {
	pc := internal.NewProbeContext()
	pc["ref.existingBool"] = true
	expr := `ref.existingBool`
	out, err := GoCelEvaluate(context.Background(), expr, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.True(t, out.(bool))
}

func TestGoCelEvaluateArrayString(t *testing.T) {
	pc := internal.NewProbeContext()
	pc["ref.strArray"] = []string{"one", "two", "three"}
	pc["ref.intArray"] = []int{10, 11, 12}
	pc["ref.uintArray"] = []uint{10, 11, 12}
	pc["ref.floatArray"] = []float64{10, 11, 12}
	pc["ref.boolArray"] = []bool{true, false, true}

	expr := `GetStringElem("ref.strArray",0)`
	out, err := GoCelEvaluate(context.Background(), expr, pc, reflect.String)
	assert.NoError(t, err)
	assert.Equal(t, "one", out.(string))

	expr = `GetIntElem("ref.intArray",1)`
	out, err = GoCelEvaluate(context.Background(), expr, pc, reflect.Int64)
	assert.NoError(t, err)
	assert.Equal(t, int64(11), out.(int64))

	expr = `GetUintElem("ref.uintArray",1)`
	out, err = GoCelEvaluate(context.Background(), expr, pc, reflect.Uint64)
	assert.NoError(t, err)
	assert.Equal(t, uint64(11), out.(uint64))

	expr = `GetFloatElem("ref.floatArray",1)`
	out, err = GoCelEvaluate(context.Background(), expr, pc, reflect.Float64)
	assert.NoError(t, err)
	assert.Equal(t, float64(11), out.(float64))

	expr = `GetBoolElem("ref.boolArray",1)`
	out, err = GoCelEvaluate(context.Background(), expr, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.False(t, out.(bool))

}

func TestGoCelEvaluateNonExistBool(t *testing.T) {
	pc := internal.NewProbeContext()
	pc["ref.existingBool"] = true
	pc["ref.existingTime"] = time.Now()

	out, err := GoCelEvaluate(context.Background(), `IsDefined("ref.existingBool")`, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.True(t, out.(bool))

	out, err = GoCelEvaluate(context.Background(), `IsDefined("ref.nonExistingBool") `, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.False(t, out.(bool))

	out, err = GoCelEvaluate(context.Background(), `IsDefined("ref_nonExistingBool") && GetBool("ref_nonExistingBool")`, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.False(t, out.(bool))

	out, err = GoCelEvaluate(context.Background(), `IsDefined("ref.existingTime") && GetTime("ref.existingTime") > timestamp("2020-01-01T10:00:20.021-05:00")`, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.True(t, out.(bool))
}
