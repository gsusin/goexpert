package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFastestApicep(t *testing.T) {
	v, err := delayedQueryFastest("88010-100", viacepId)
	as := assert.New(t)
	as.Nil(err)
	as.NotNil(v)
	result := v.(Apicep)
	if !as.Truef(result.Ok, "API apicep.com returned error. Status: %d", result.Status) {
		return
	}
	as.Equal("Florianópolis", result.City)
}

func TestFastestViacep(t *testing.T) {
	v, err := delayedQueryFastest("88010-100", apicepId)
	as := assert.New(t)
	as.Nil(err)
	as.NotNil(v)
	result := v.(Viacep)
	if !as.False(result.Erro) {
		return
	}
	as.Equal("Florianópolis", result.Localidade)
}

func TestNonexistentApicep(t *testing.T) {
	v, err := delayedQueryFastest("88000-000", viacepId)
	as := assert.New(t)
	as.Nil(err)
	as.NotNil(v)
	result := v.(Apicep)
	as.Equal(404, result.Status)
}

func TestNonexistentViacep(t *testing.T) {
	v, err := delayedQueryFastest("88000-000", apicepId)
	as := assert.New(t)
	as.Nil(err)
	as.NotNil(v)
	result := v.(Viacep)
	as.True(result.Erro)
}

func TestInvalidViacep(t *testing.T) {
	_, err := delayedQueryFastest("88000", apicepId)
	as := assert.New(t)
	as.NotNil(err)
	as.Equal("StatusBadRequest", err.Error())
}

func delayedQueryFastest(cep string, delayedApiId uint8) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, delayContextKey("delayPlan"), queryDelay{5 * time.Second, delayedApiId})
	v, err := QueryFastest(ctx, cep)
	return v, err
}
