package testutil

import (
	"errors"
	"testing"
)

func ExampleErr() error {
	return errors.New("An error abc123!")
}

func ExampleNoErr() error {
	return nil
}

func ExampleEq() string {
	return "foo"
}

func ExamplePanic() {
	panic("Oh no, a lolcat!")
}

func TestTestutil(t *testing.T) {
	assert := NewAssert(t)
	assert.Ok("1+2=3", 1 + 2 == 3)
	assert.Err("ExampleErr should return an error", "abc123", ExampleErr())
	assert.NoErr("ExampleNoErr should not return an error", ExampleNoErr())
	assert.Eq("ExampleEq", ExampleEq(), "foo")
	assert.Panic("lolcats?", ExamplePanic)
	assert.Eq("Repr", Repr(map[string]int{ "A": 1, "B": 2 }), `map[string]int{"A":1, "B":2}`)
}
