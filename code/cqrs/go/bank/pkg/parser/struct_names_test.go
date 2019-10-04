package parser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pdt256/talks/code/cqrs/go/bank/pkg/parser"
)

func TestGetStructNames_ReturnsEmptyList(t *testing.T) {
	// Given
	code := `package test`
	content := strings.NewReader(code)

	// When
	structNames, err := parser.GetStructNames(content)

	// Then
	require.NoError(t, err)
	assert.Equal(t, 0, len(structNames))
}

func TestGetStructNames_ReturnsOneStructNameWithOtherKeywordsPresent(t *testing.T) {
	// Given
	code := `package test
type Struct1 struct {}
type Interface1 interface {}
func Func1() {}`
	content := strings.NewReader(code)

	// When
	structNames, err := parser.GetStructNames(content)

	// Then
	require.NoError(t, err)
	assert.Equal(t, 1, len(structNames))
	assert.Equal(t, []string{"Struct1"}, structNames)
}

func TestGetStructNames_ReturnsTwoStructName(t *testing.T) {
	// Given
	code := `package test
type Struct1 struct {}
type Struct2 struct{}`
	content := strings.NewReader(code)

	// When
	structNames, err := parser.GetStructNames(content)

	// Then
	require.NoError(t, err)
	assert.Equal(t, []string{"Struct1", "Struct2"}, structNames)
}
