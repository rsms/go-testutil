package testutil

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

// RecoverAsFail catches a panic and converts it into a test failure.
// Example:
//   func TestThing(t *testing.T) {
//     defer testutil.RecoverAsFail(t)
//     somethingThatMayPanic()
//   }
func RecoverAsFail(t *testing.T) {
	if v := recover(); v != nil {
		t.Log(v)
		t.Log(string(debug.Stack()))
		t.Fail()
	}
}

type Assert struct {
	T *testing.T

	// note: T instead of implicit since testing.T has function like Error that may easily cause
	// mistakes, i.e. type assert.Error("thing", "msg", err) would not do what you think.
}

func NewAssert(t *testing.T) Assert {
	return Assert{t}
}

func (a Assert) Ok(assertionfmt string, ok bool, v ...interface{}) bool {
	a.T.Helper() // mark this function as a helper (for stack traces)
	if !ok {
		a.T.Errorf("FAIL: (assertion) "+assertionfmt, v...)
	}
	return ok
}

func (a Assert) Err(assertionfmt, errsubstr string, err error, v ...interface{}) bool {
	a.T.Helper()
	assertionfmt = "FAIL: " + assertionfmt
	if err == nil {
		a.T.Errorf(assertionfmt+"\nexpected error with substring %q\n(no error)",
			append(v, errsubstr)...)
		return false
	}
	if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(errsubstr)) {
		a.T.Errorf(
			assertionfmt+"\nExpected error to contain substring %q (case-insensitive)\nGot %q",
			append(v, errsubstr, err.Error())...)
		return false
	}
	return true
}

func (a Assert) NoErr(assertionfmt string, err error, v ...interface{}) bool {
	a.T.Helper()
	if err != nil {
		a.T.Errorf("FAIL: "+assertionfmt+"; error: %v", append(v, err)...)
	}
	return err == nil
}

// Eq checks that value1 and value2 are equal.
//
// Accepts the following types of values:
// - Anything type that you can compare with "==" in Go (bool, int, float, string, etc.)
// - Pointer types (e.g. *struct)
// - []byte
//
func (a Assert) Eq(assertionfmt string, value1, value2 interface{}, v ...interface{}) bool {
	a.T.Helper()
	if value1 == nil && value2 == nil {
		return true
	}

	t1, t2 := reflect.TypeOf(value1), reflect.TypeOf(value2)

	if value1 == nil || value2 == nil {
		if value1 == nil && reflect.ValueOf(value2).IsZero() {
			// e.g. var value1 []byte, var value2 []byte = []byte{}
			return true
		}
		if value2 == nil && reflect.ValueOf(value1).IsZero() {
			// e.g. var value2 []byte, var value1 []byte = []byte{}
			return true
		}
		a.T.Errorf("FAIL: "+assertionfmt+"\nvalue1: %s\nvalue2: %s",
			append(v, Repr(value1), Repr(value2))...)
		return false
	}

	if t1 != t2 {
		a.T.Errorf("FAIL: "+assertionfmt+"; types differ:\nvalue1: %T\nvalue2: %T",
			append(v, value1, value2)...)
		return false
	}

	v1, v2 := reflect.ValueOf(value1), reflect.ValueOf(value2)
	if v1.CanAddr() && v1.UnsafeAddr() == v2.UnsafeAddr() {
		return true
	}

	eq := false

	if t1.Comparable() {
		eq = value1 == value2
	} else {
		switch v1 := value1.(type) {
		case []byte:
			eq = bytes.Equal(v1, value2.([]byte))
		default:
			a.T.Errorf("FAIL: "+assertionfmt+"; assert.Eq can not compare values of type %T",
				append(v, value1)...)
			return false
		}
	}

	if !eq {
		a.T.Errorf("FAIL: "+assertionfmt+"\nvalue1: %s\nvalue2: %s",
			append(v, Repr(value1), Repr(value2))...)
		return false
	}

	return true
}

func (a Assert) Panic(expectedPanicRegExp string, f func()) bool {
	a.T.Helper()
	ok := false
	// Note: (?i) makes it case-insensitive
	expected := regexp.MustCompile("(?i)" + expectedPanicRegExp)
	defer func() {
		a.T.Helper()
		if v := recover(); v != nil {
			panicMsg := fmt.Sprint(v)
			if ok = expected.MatchString(panicMsg); !ok {
				a.T.Log(string(debug.Stack()))
				a.T.Errorf("expected panic to match %q but got %q", expectedPanicRegExp, panicMsg)
			}
		} else {
			a.T.Log(string(debug.Stack()))
			a.T.Errorf("expected panic (but there was no panic)")
		}
	}()
	f()
	return ok
}

func Repr(v interface{}) string {
	switch v.(type) {
	case []byte, string:
		return fmt.Sprintf("%q", v)
	}
	return pretty.Sprintf("%# v", v)
}
