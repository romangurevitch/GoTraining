package receivers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmutable(t *testing.T) {
	immutable := &immutable{
		Value: "will this be changed?",
	}

	newImmutable := immutable.SetString("nope, it will not change!")
	assert.Equal(t, "will this be changed?", immutable.String())
	assert.Equal(t, "nope, it will not change!", newImmutable.String())
	assert.NotEqual(t, immutable, newImmutable)
}

func TestMutable(t *testing.T) {
	mutable := &mutable{
		Value: "will this be changed?",
	}

	newMutable := mutable.SetString("yes, it will be changed!")
	assert.Equal(t, "yes, it will be changed!", mutable.String())
	assert.Equal(t, "yes, it will be changed!", newMutable.String())
	assert.Equal(t, mutable, newMutable)
}

func TestImmutableJSON(t *testing.T) {
	immutable := immutable{
		Value: "will this be changed?",
	}

	bytes, err := json.Marshal(immutable)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":"yes it will be changed!"}`, string(bytes))
}

func TestMutableJSON(t *testing.T) {
	mutable := mutable{
		Value: "will this be changed?",
	}

	bytes, err := json.Marshal(mutable)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":"will this be changed?"}`, string(bytes))
}
