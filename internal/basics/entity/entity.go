package entity

import "fmt"

type Entity interface {
	Do(arg string) (string, error)
}

func New(externalField string) Entity {
	return &entity{
		ExternalField: externalField,
		internalField: "internal field",
	}
}

type entity struct {
	ExternalField string
	internalField string
}

func (e *entity) String() string {
	return fmt.Sprintf("external: %v, internal: %s", e.ExternalField, e.internalField)
}

func (e *entity) Do(arg string) (string, error) {
	fmt.Println(e, arg)
	return e.internalField, nil
}
