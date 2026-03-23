package mocking

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/romangurevitch/go-training/internal/basics/mocking/calculator/mocks"
)

// TestExampleFunction_Times shows how to assert the number of times a method is called.
func TestExampleFunction_Times(t *testing.T) {
	a := new(mocks.Adder)

	// Expect exactly one call
	a.On("SingleDigitAdd", 1, 2).Return(3, nil).Once()

	got, err := ExampleFunction(a, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, got)

	a.AssertExpectations(t)
}

// TestExampleFunction_Run shows Run: execute custom logic when the mock is called.
// This is often used to set internal state or perform side effects.
func TestExampleFunction_Run(t *testing.T) {
	a := new(mocks.Adder)
	called := false

	a.On("SingleDigitAdd", 1, 2).
		Return(3, nil).
		Run(func(args mock.Arguments) {
			called = true
		})

	_, _ = ExampleFunction(a, 1, 2)
	assert.True(t, called)
	a.AssertExpectations(t)
}

// TestExampleFunction_Anything shows mock.Anything: allow any value for an argument.
func TestExampleFunction_Anything(t *testing.T) {
	a := new(mocks.Adder)

	a.On("SingleDigitAdd", mock.Anything, mock.Anything).Return(10, nil)

	got, _ := ExampleFunction(a, 5, 5)
	assert.Equal(t, 10, got)
	a.AssertExpectations(t)
}

func TestExampleFunction(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *mocks.Adder)
		x, y    int
		want    int
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *mocks.Adder) {
				m.On("SingleDigitAdd", 1, 2).Return(3, nil)
			},
			x:    1,
			y:    2,
			want: 3,
		},
		{
			name: "error from mock",
			setup: func(m *mocks.Adder) {
				m.On("SingleDigitAdd", 1, 2).Return(0, errors.New("mock error"))
			},
			x:       1,
			y:       2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mocks.Adder)
			if tt.setup != nil {
				tt.setup(m)
			}

			got, err := ExampleFunction(m, tt.x, tt.y)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			m.AssertExpectations(t)
		})
	}
}
