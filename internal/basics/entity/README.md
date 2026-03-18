# Structs and Entities

Go uses structs to model domain entities. Unlike Python dicts, structs are typed and statically defined.

## Defining a Struct

```go
type Account struct {
    ID        int       `json:"id"`
    FirstName string    `json:"first_name"`
    Balance   float64   `json:"balance"`
}
```

## Struct Tags

- `json:"field_name"` — JSON key
- `yaml:"field_name"` — YAML key

## Pitfalls

- Unexported fields (lowercase) are invisible to JSON/YAML encoders
- Zero values: structs are valid even when not initialised
