package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	checks := []string{"SA1000", "SA1001"}

	data := []byte(`{"staticcheck":["SA1000","SA1001"]}`)
	var cfg ConfigData
	err := json.Unmarshal(data, &cfg)
	assert.NoError(t, err)

	if !reflect.DeepEqual(cfg.Staticcheck, checks) {
		t.Errorf("expected Staticcheck to be %v, but got %v", checks, cfg.Staticcheck)
	}
}
