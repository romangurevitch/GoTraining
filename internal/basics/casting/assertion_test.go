package casting

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Read more at:
// https://golang.org/ref/spec#Type_assertions
func TestAssertions(t *testing.T) {
	entityInterfaceInstance := New()
	err := entityInterfaceInstance.Do("the first thing")
	assert.NoError(t, err)

	// I want to run 'DoSomethingElse'.
	entityStruct, ok := entityInterfaceInstance.(*entity)
	assert.True(t, ok)

	err = entityStruct.DoSomethingElse("the second thing")
	assert.NoError(t, err)

	// What will happen in this case?
	//entity2Struct, ok := entityInterfaceInstance.(*entity2)
	//assert.True(t, ok)
	//err = entity2Struct.DoAThirdThing("do something else")
	//assert.NoError(t, err)
}

type Entity interface {
	Do(string) error
}

func New() Entity {
	return &entity{}
}

type entity struct {
}

func (e *entity) Do(s string) error {
	fmt.Printf("Run `Do` with, %v\n", s)
	return nil
}

func (e *entity) DoSomethingElse(s string) error {
	fmt.Printf("Run `DoSomethingElse` with, %v\n", s)
	return nil
}

//nolint
type entity2 struct {
}

//nolint
func (e *entity2) Do(s string) error {
	panic("implement me")
}

//nolint
func (e *entity2) DoAThirdThing(s string) error {
	fmt.Printf("Run `DoAThirdThing` with, %v\n", s)
	return nil
}
