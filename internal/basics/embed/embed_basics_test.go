package embed

import (
	"fmt"
	"testing"
)

type embeddedStruct struct {
	embeddedStruct1
	num int
}

type containerStruct struct {
	embeddedStruct
	str string
}

// Note function receiver is on the embedded struct
func (b embeddedStruct) describe() string {
	return fmt.Sprintf("base with num=%v", b.num)
}

func Test_base_describe(t *testing.T) {

	co := containerStruct{ // A containerStruct embeds an embeddedStruct. An embedding looks like a field without a name.
		embeddedStruct: embeddedStruct{ // When creating structs with literals, we have to initialize the embedding explicitly; here the embedded type serves as the field name.
			num: 1,
		},

		str: "some name",
	}

	t.Logf("co={num: %v, str: %v}\n", co.embeddedStruct.num, co.str) // We can access the base's fields directly on co, e.g. co.num
	t.Logf("also num: %d", co.embeddedStruct.embeddedStruct1.num)    // Alternatively, we can spell out the full path using the embedded type name.
	t.Logf("describe: %s", co.describe())

	// What happens with the describe method if we embed a struct within a struct
	t.Logf("describe: %s", co.describe())

	type describer interface {
		describe() string
	}

	var d describer = co                  // Since containerStruct embeds embeddedStruct, the methods of embeddedStruct also become methods of a container.
	t.Logf("describer: %s", d.describe()) // Here we invoke a method that was embedded from base directly on co.

	// Do you think the following would be a valid call?
	//t.Logf("str: %s", d.str)
}
