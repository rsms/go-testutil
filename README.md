Simple test helper for go

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/rsms/go-testutil.svg)][godoc]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/rsms/go-testutil)][godoc]
[![Go Report Card](https://goreportcard.com/badge/github.com/rsms/go-testutil)](https://goreportcard.com/report/github.com/rsms/go-testutil)

[godoc]: https://pkg.go.dev/github.com/rsms/go-testutil

Synopsis

```
type Assert struct { T *testing.T }
func NewAssert(t *testing.T) Assert
func (a Assert) Eq(assertionfmt string, value1, value2 interface{}, v ...interface{}) bool
func (a Assert) Err(assertionfmt, errsubstr string, err error, v ...interface{}) bool
func (a Assert) NoErr(assertionfmt string, err error, v ...interface{}) bool
func (a Assert) Ok(assertionfmt string, ok bool, v ...interface{}) bool
func (a Assert) Panic(expectedPanicRegExp string, f func()) bool
func RecoverAsFail(t *testing.T)
func Repr(v interface{}) string
```

## Examples

```go
func TestFoo(t *testing.T) {
  assert := testutil.NewAssert(t)
  assert.Eq("Foo does the expected thing", Foo(), "bar")
}
```

RecoverAsFail catches a panic and converts it into a test failure

```go
func TestThing(t *testing.T) {
  defer testutil.RecoverAsFail(t)
  somethingThatMayPanic()
}
```
