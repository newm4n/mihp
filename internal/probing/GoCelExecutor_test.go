package probing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestGoCelEvaluateExistBool(t *testing.T) {
	pc := NewProbeContext()
	pc["ref.existingBool"] = true
	expr := `ref.existingBool`
	out, err := GoCelEvaluate(context.Background(), expr, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.True(t, out.(bool))
}

func TestGoCelEvaluateNonExistBool(t *testing.T) {
	pc := NewProbeContext()
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
