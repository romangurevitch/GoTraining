# Struct Embedding

Go uses embedding instead of inheritance. Embedding promotes methods and fields to the outer struct.

## Basic Embedding

```go
type Animal struct{ Name string }
func (a Animal) Speak() string { return a.Name }

type Dog struct {
    Animal       // embedded — Dog gets Speak() for free
    Breed string
}
```

## Pitfalls

- Embedding is not inheritance — the outer type is not a subtype of the embedded type
- If both outer and embedded define the same method, the outer wins (shadowing)
