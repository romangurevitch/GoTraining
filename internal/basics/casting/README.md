# Type Assertions and Conversions

Go distinguishes between **type assertions** (interface → concrete type) and **type conversions** (compatible concrete types).

## Type Assertion

```go
var i interface{} = "hello"
s, ok := i.(string) // safe: ok is false if assertion fails
s := i.(string)     // unsafe: panics if i is not a string
```

## Type Switch

```go
switch v := i.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
}
```

## Pitfalls

- A failed type assertion without `ok` will **panic** at runtime
- You cannot convert between incompatible types without parsing
