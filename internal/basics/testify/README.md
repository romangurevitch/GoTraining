# Testify

[testify](https://github.com/stretchr/testify) is the standard Go assertion library.

## Assertions

```go
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.True(t, condition)
assert.Nil(t, value)
```

Use `require` to stop the test immediately on failure:

```go
require.NoError(t, err) // test stops here if err != nil
```

## Test Suites

```go
type MyTestSuite struct{ suite.Suite }
func (s *MyTestSuite) TestSomething() { s.Equal(1, 1) }
func TestMyTestSuite(t *testing.T) { suite.Run(t, new(MyTestSuite)) }
```

## Pitfalls

- `assert` continues after failure; `require` stops — use `require` for setup steps
