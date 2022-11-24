package utils_test

import (
	"moonspace/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Foo string
	Bar string
	Val int64
}

func Test_Inspect(t *testing.T) {
	x := Foo{
		Foo: "1",
		Bar: "2",
		Val: 3,
	}

	res := utils.Inspect(&x)
	expected := map[string]any{
		"Foo": "1",
		"Bar": "2",
		"Val": int64(3),
	}

	assert.Equal(t, expected, res)
}
