package embed

import (
	"fmt"
	"testing"
)

type embeddedStruct1 struct {
	embeddedStruct2
	num int
}

type embeddedStruct2 struct {
	num int
}

// Note function receiver is on the embedded struct
func (b embeddedStruct2) overshadowMethod() string {
	return fmt.Sprintf("base with num=%v", b.num)
}

func (b embeddedStruct1) overshadowMethod() string {
	return fmt.Sprintf("base with num=%v", b.num)
}

func (b containerStruct) overshadowMethod() string {
	return fmt.Sprintf("base with str=%v", b.str)
}

type baseStruct struct {
	embeddedStruct1
	str string
	num int
}

func Test_base_playground_describe(t *testing.T) {
	co := baseStruct{
		num: 0,
		embeddedStruct1: embeddedStruct1{
			num: 1,
			embeddedStruct2: embeddedStruct2{
				num: 2,
			},
		},

		str: "some name",
	}

	t.Logf("co.num: %v, co.str: %v}\n", co.num, co.str)
	t.Logf("also co.embeddedStruct1.num: %d", co.embeddedStruct1.num) // Now members of children members need to be accessed by name
	t.Logf("also co.embeddedStruct2.num: %d", co.embeddedStruct2.num)
	t.Logf("also co.embeddedStruct1.embeddedStruct2.num: %d", co.embeddedStruct1.embeddedStruct2.num) //nolint:staticcheck // explicit selector demo

	t.Logf("overshadowMethod: %s", co.overshadowMethod())
	t.Logf("also co.embeddedStruct1.overshadowMethod: %s", co.embeddedStruct1.overshadowMethod()) //nolint:staticcheck // explicit selector demo
	t.Logf("also co.embeddedStruct2.overshadowMethod: %s", co.embeddedStruct2.overshadowMethod())
	t.Logf("also co.embeddedStruct1.embeddedStruct2.overshadowMethod: %s", co.embeddedStruct1.embeddedStruct2.overshadowMethod()) //nolint:staticcheck // explicit selector demo

	type overshadowMethod interface {
		overshadowMethod() string
	}

	var d overshadowMethod = co
	t.Logf("interface co: %s", d.overshadowMethod())
	d = co.embeddedStruct1
	t.Logf("interface co.embeddedStruct1: %s", d.overshadowMethod())
	d = co.embeddedStruct2
	t.Logf("interface co.embeddedStruct2: %s", d.overshadowMethod())
	d = co.embeddedStruct1.embeddedStruct2 //nolint:staticcheck // explicit selector demo
	t.Logf("interface co.embeddedStruct1.embeddedStruct2: %s", d.overshadowMethod())

}
