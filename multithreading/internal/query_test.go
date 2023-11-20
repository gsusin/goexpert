package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// https://pkg.go.dev/context@go1.21.3#WithValue

func TestFastestApicep(t *testing.T) {
	ctx, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()
	//	d := delay{}
	ctx = context.WithValue(ctx, delayContextKey("delayPlan"), viacepDelay(5*time.Second))
	a, err := QueryFastest(ctx, "88037-500")
	assert.Nil(t, err)
	assert.NotNil(t, a)
	aTyped := a.(Apicep)
	assert.Equal(t, aTyped.City, "Florianópolis")
}
