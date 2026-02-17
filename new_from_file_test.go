package jman_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
)

func TestNewFromFile_Obj(t *testing.T) {
	path := filepath.Join(t.TempDir(), "obj.json")
	err := os.WriteFile(path, []byte(`{"key":"value","num":2}`), 0o600)
	assert.NoError(t, err)

	obj := jman.NewFromFile[jman.Obj](t, path)
	assert.Equal(t, jman.Obj{
		"key": "value",
		"num": float64(2),
	}, obj)
}

func TestNewFromFile_Arr(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arr.json")
	err := os.WriteFile(path, []byte(`["a",2,true]`), 0o600)
	assert.NoError(t, err)

	arr := jman.NewFromFile[jman.Arr](t, path)
	assert.Equal(t, jman.Arr{
		"a",
		float64(2),
		true,
	}, arr)
}

func TestNewFromFile_FileNotFound(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing.json")
	mt := newMockT(fmt.Sprintf("error reading JSON file %s:", missingPath))
	defer mt.AssertExpectations(t)

	assert.Panics(t, func() { jman.NewFromFile[jman.Obj](mt, missingPath) })
}

func TestNewFromFile_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	err := os.WriteFile(path, []byte(`{"a":`), 0o600)
	assert.NoError(t, err)

	mt := newMockT(`error parsing JSON data {"a"::`)
	defer mt.AssertExpectations(t)

	assert.Panics(t, func() { jman.NewFromFile[jman.Obj](mt, path) })
}
